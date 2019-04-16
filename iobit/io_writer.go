package iobit

/*
 * Copyright 2019/04/07- yoya@awm.jp. All rights reserved.
 */

import (
	"encoding/binary"
	"fmt"
	"io"
)

type IOWriter struct {
	// Write method
	Writer     io.Writer
	Binary     binary.ByteOrder
	OffsetByte uint64
	OffsetBit  uint64
	Buff       []byte
	lastError  error
}

func NewIOWriter(w io.Writer, b binary.ByteOrder) *IOWriter {
	return &IOWriter{Writer: w, Binary: b,
		OffsetByte: 0, OffsetBit: 0, Buff: make([]byte, 8)}
}

func (w *IOWriter) Write(buff []byte) (int, error) {
	w.AlignByte()
	if w.lastError != nil {
		return 0, w.lastError
	}
	var n int
	n, w.lastError = w.Writer.Write(buff)
	w.OffsetByte += uint64(n)
	return n, w.lastError
}

func (w *IOWriter) GetOffset() (uint64, uint64) {
	return w.OffsetByte, w.OffsetBit
}

func (w *IOWriter) AlignByte() {
	if w.lastError != nil {
		return
	}
	if w.OffsetBit > 0 {
		var n int
		n, w.lastError = w.Writer.Write(w.Buff[:1])
		w.OffsetByte += uint64(n)
		w.OffsetBit = 0
		w.Buff[0] = 0
	}
}

func (w *IOWriter) PutUInt8(v uint8) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.Buff[0] = v
	var n int
	n, w.lastError = w.Writer.Write(w.Buff[:1])
	w.OffsetByte += uint64(n)
	if w.lastError != nil {
		return
	}
}

func (w *IOWriter) PutUInt16(v uint16) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.Binary.PutUint16(w.Buff[:2], v)
	var n int
	n, w.lastError = w.Writer.Write(w.Buff[:2])
	w.OffsetByte += uint64(n)
}

func (w *IOWriter) PutUInt24(v uint32) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	switch w.Binary {
	case BigEndian:
		w.Buff[0] = uint8((v << 16) & 0xff)
		w.Buff[1] = uint8((v << 8) & 0xff)
		w.Buff[2] = uint8((v) & 0xff)
	case LittleEndian:
		w.Buff[2] = uint8((v << 16) & 0xff)
		w.Buff[1] = uint8((v << 8) & 0xff)
		w.Buff[0] = uint8((v) & 0xff)
	default:
		w.lastError = fmt.Errorf("PutUInt24 unsupported binary:%#v", w.Binary)
		return
	}
	var n int
	n, w.lastError = w.Writer.Write(w.Buff[:3])
	w.OffsetByte += uint64(n)
}

func (w *IOWriter) PutUInt32(v uint32) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.Binary.PutUint32(w.Buff[:4], v)
	var n int
	n, w.lastError = w.Writer.Write(w.Buff[:4])
	w.OffsetByte += uint64(n)
}

func (w *IOWriter) PutUInt64(v uint64) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.Binary.PutUint64(w.Buff[:8], v)
	var n int
	n, w.lastError = w.Writer.Write(w.Buff[:8])
	w.OffsetByte += uint64(n)
}

func (w *IOWriter) PutUIBit(v uint8) {
	if w.lastError != nil {
		return
	}
	if v == 1 {
		w.Buff[0] = w.Buff[0] | (1 << uint8(7-w.OffsetBit))
	}
	w.OffsetBit += 1
	if w.OffsetBit > 7 {
		var n int
		n, w.lastError = w.Writer.Write(w.Buff[:1])
		if n > 0 {
			w.OffsetByte += uint64(n)
			w.OffsetBit -= 8
		}
		if w.lastError != nil {
			return
		}
		w.Buff[0] = 0
	}
}

func (w *IOWriter) PutUIBits_uint8(v uint8, n int) {
	if w.lastError != nil {
		return
	}
	if n > 8 {
		w.lastError = fmt.Errorf("PutUIBits_uint8 n:%d > 8", n)
		return
	}
	w.PutUIBits_uint64(uint64(v), n)
}

func (w *IOWriter) PutUIBits_uint16(v uint16, n int) {
	if w.lastError != nil {
		return
	}
	if n > 16 {
		w.lastError = fmt.Errorf("PutUIBits_uint16 n:%d > 16", n)
		return
	}
	w.PutUIBits_uint64(uint64(v), n)
}

func (w *IOWriter) PutUIBits_uint32(v uint32, n int) {
	if w.lastError != nil {
		return
	}
	if n > 32 {
		w.lastError = fmt.Errorf("PutUIBits_uint32 n:%d > 32", n)
		return
	}
	w.PutUIBits_uint64(uint64(v), n)
}

func (w *IOWriter) PutUIBits_uint64(v uint64, n int) {
	if w.lastError != nil {
		return
	}
	if n > 64 {
		w.lastError = fmt.Errorf("PutUIBits_uint64 n:%d > 64", n)
		return
	}
	for i := 0; i < n; i++ {
		b := (v >> uint8(n-1-i)) & 1
		w.PutUIBit(uint8(b))
		if w.lastError != nil { // include io.EOF
			return
		}
	}
}

func (r *IOWriter) GetLastError() error {
	if r.lastError == nil {
		return nil
	}
	if r.lastError == io.EOF {
		return r.lastError
	}
	err := fmt.Errorf("%s in offset(%d:%d)",
		r.lastError, r.OffsetByte, r.OffsetBit)
	return err
}
