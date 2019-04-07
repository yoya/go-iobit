package iobit

/*
 * Copyright 2019/04/07- yoya@awm.jp. All rights reserved.
 */

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Writer struct {
	// Write method
	Writer     io.Writer
	Binary     binary.ByteOrder
	OffsetByte uint64
	OffsetBit  uint64
	Buff       []byte
}

func NewWriter(w io.Writer, b binary.ByteOrder) *Writer {
	return &Writer{Writer: w, Binary: b,
		OffsetByte: 0, OffsetBit: 0, Buff: make([]byte, 8)}
}

func (w *Writer) Write(buff []byte) (int, error) {
	w.AlignByte()
	w.OffsetByte += uint64(len(buff))
	return w.Writer.Write(buff)
}

func (w *Writer) GetOffset() (uint64, uint64) {
	return w.OffsetByte, w.OffsetBit
}

func (w *Writer) AlignByte() error {
	if w.OffsetBit > 0 {
		_, err := w.Writer.Write(w.Buff[:1])
		if err != nil {
			return err
		}
		w.OffsetByte += 1
		w.OffsetBit = 0
		w.Buff[0] = 0
	}
	return nil
}

func (w *Writer) PutUInt8(v uint8) error {
	w.AlignByte()
	w.Buff[0] = v
	_, err := w.Writer.Write(w.Buff[:1])
	w.OffsetByte += 1
	return err
}

func (w *Writer) PutUInt16(v uint16) error {
	err := w.AlignByte()
	if err != nil {
		return err
	}
	w.Binary.PutUint16(w.Buff[:2], v)
	_, err = w.Writer.Write(w.Buff[:2])
	w.OffsetByte += 2
	return err
}

func (w *Writer) PutUInt24(v uint32) error {
	err := w.AlignByte()
	if err != nil {
		return err
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
		return fmt.Errorf("PutUInt24 unsupported binary:%#v", w.Binary)
	}
	_, err = w.Writer.Write(w.Buff[:3])
	w.OffsetByte += 3
	return err
}

func (w *Writer) PutUInt32(v uint32) error {
	err := w.AlignByte()
	if err != nil {
		return err
	}
	w.Binary.PutUint32(w.Buff[:4], v)
	_, err = w.Writer.Write(w.Buff[:4])
	w.OffsetByte += 4
	return err
}

func (w *Writer) PutUInt64(v uint64) error {
	err := w.AlignByte()
	if err != nil {
		return err
	}
	w.Binary.PutUint64(w.Buff[:8], v)
	_, err = w.Writer.Write(w.Buff[:8])
	w.OffsetByte += 8
	return err
}

func (w *Writer) PutUIBit(v uint8) error {
	if v == 1 {
		w.Buff[0] = w.Buff[0] | (1 << uint8(7-w.OffsetBit))
	}
	w.OffsetBit += 1
	if w.OffsetBit > 7 {
		_, err := w.Writer.Write(w.Buff[:1])
		w.Buff[0] = 0
		if err != nil {
			return err
		}
		w.OffsetByte += 1
		w.OffsetBit -= 8
	}
	return nil
}

func (w *Writer) PutUIBits_uint8(v uint8, n int) error {
	if n > 8 {
		return fmt.Errorf("PutUIBits_uint8 n:%d > 8", n)
	}
	return w.PutUIBits_uint64(uint64(v), n)
}

func (w *Writer) PutUIBits_uint16(v uint16, n int) error {
	if n > 16 {
		return fmt.Errorf("PutUIBits_uint16 n:%d > 16", n)
	}
	return w.PutUIBits_uint64(uint64(v), n)
}

func (w *Writer) PutUIBits_uint32(v uint32, n int) error {
	if n > 32 {
		return fmt.Errorf("PutUIBits_uint32 n:%d > 32", n)
	}
	return w.PutUIBits_uint64(uint64(v), n)
}

func (w *Writer) PutUIBits_uint64(v uint64, n int) error {
	if n > 64 {
		return fmt.Errorf("PutUIBits_uint64 n:%d > 64", n)
	}
	for i := 0; i < n; i++ {
		b := (v >> uint8(n-1-i)) & 1
		err := w.PutUIBit(uint8(b))
		if err != nil { // include io.EOF
			return err
		}
	}
	return nil
}
