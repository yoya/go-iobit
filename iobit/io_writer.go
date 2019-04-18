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
	writer     io.Writer
	binary     binary.ByteOrder
	offsetByte uint64
	offsetBit  uint64
	buff       []byte
	lastError  error
}

func NewIOWriter(w io.Writer, b binary.ByteOrder) *IOWriter {
	return &IOWriter{writer: w, binary: b,
		offsetByte: 0, offsetBit: 0, buff: make([]byte, 8)}
}

func (w *IOWriter) Write(buff []byte) (int, error) {
	w.AlignByte()
	if w.lastError != nil {
		return 0, w.lastError
	}
	var n int
	n, w.lastError = w.writer.Write(buff)
	w.offsetByte += uint64(n)
	return n, w.lastError
}

func (w *IOWriter) GetOffset() (uint64, uint64) {
	return w.offsetByte, w.offsetBit
}

func (w *IOWriter) AlignByte() {
	if w.lastError != nil {
		return
	}
	if w.offsetBit > 0 {
		var n int
		n, w.lastError = w.writer.Write(w.buff[:1])
		w.offsetByte += uint64(n)
		w.offsetBit = 0
		w.buff[0] = 0
	}
}

func (w *IOWriter) PutUInt8(v uint8) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.buff[0] = v
	var n int
	n, w.lastError = w.writer.Write(w.buff[:1])
	w.offsetByte += uint64(n)
	if w.lastError != nil {
		return
	}
}

func (w *IOWriter) PutUInt16(v uint16) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.binary.PutUint16(w.buff[:2], v)
	var n int
	n, w.lastError = w.writer.Write(w.buff[:2])
	w.offsetByte += uint64(n)
}

func (w *IOWriter) PutUInt24(v uint32) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	switch w.binary {
	case BigEndian:
		w.buff[0] = uint8((v << 16) & 0xff)
		w.buff[1] = uint8((v << 8) & 0xff)
		w.buff[2] = uint8((v) & 0xff)
	case LittleEndian:
		w.buff[2] = uint8((v << 16) & 0xff)
		w.buff[1] = uint8((v << 8) & 0xff)
		w.buff[0] = uint8((v) & 0xff)
	default:
		w.lastError = fmt.Errorf("PutUInt24 unsupported binary:%#v", w.binary)
		return
	}
	var n int
	n, w.lastError = w.writer.Write(w.buff[:3])
	w.offsetByte += uint64(n)
}

func (w *IOWriter) PutUInt32(v uint32) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.binary.PutUint32(w.buff[:4], v)
	var n int
	n, w.lastError = w.writer.Write(w.buff[:4])
	w.offsetByte += uint64(n)
}

func (w *IOWriter) PutUInt64(v uint64) {
	w.AlignByte()
	if w.lastError != nil {
		return
	}
	w.binary.PutUint64(w.buff[:8], v)
	var n int
	n, w.lastError = w.writer.Write(w.buff[:8])
	w.offsetByte += uint64(n)
}

func (w *IOWriter) PutUIBit(v uint8) {
	if w.lastError != nil {
		return
	}
	if w.offsetBit == 0 {
		w.buff[0] = 0
	}
	if v == 1 {
		w.buff[0] = w.buff[0] | (1 << uint8(7-w.offsetBit))
	}
	w.offsetBit += 1
	if w.offsetBit > 7 {
		var n int
		n, w.lastError = w.writer.Write(w.buff[:1])
		if n > 0 {
			w.offsetByte += uint64(n)
			w.offsetBit -= 8
		}
		if w.lastError != nil {
			return
		}
		w.buff[0] = 0
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

func (w *IOWriter) PutBytes(bytes []byte) {
	if w.lastError != nil {
		return
	}
	var n int
	n, w.lastError = w.writer.Write(bytes)
	w.offsetByte += uint64(n)
}
func (w *IOWriter) PutString(str string) {
	if w.lastError != nil {
		return
	}
	w.PutBytes([]byte(str))
}

func (r *IOWriter) GetLastError() error {
	if r.lastError == nil {
		return nil
	}
	if r.lastError == io.EOF {
		return r.lastError
	}
	err := fmt.Errorf("%s in offset(%d:%d)",
		r.lastError, r.offsetByte, r.offsetBit)
	return err
}
