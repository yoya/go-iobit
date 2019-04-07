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

func NewWriter(writer io.Writer, binary binary.ByteOrder) *Writer {
	return &Writer{Writer: writer, Binary: binary,
		OffsetByte: 0, OffsetBit: 0, Buff: make([]byte, 8)}
}

func (iob *Writer) Write(buff []byte) (int, error) {
	iob.AlignByte()
	iob.OffsetByte += uint64(len(buff))
	return iob.Writer.Write(buff)
}

func (iob *Writer) GetOffset() (uint64, uint64) {
	return iob.OffsetByte, iob.OffsetBit
}

func (iob *Writer) AlignByte() error {
	if iob.OffsetBit > 0 {
		_, err := iob.Writer.Write(iob.Buff[:1])
		if err != nil {
			return err
		}
		iob.OffsetByte += 1
		iob.OffsetBit = 0
		iob.Buff[0] = 0
	}
	return nil
}

func (iob *Writer) PutUInt8(v uint8) error {
	iob.AlignByte()
	iob.Buff[0] = v
	_, err := iob.Writer.Write(iob.Buff[:1])
	iob.OffsetByte += 1
	return err
}

func (iob *Writer) PutUInt16(v uint16) error {
	err := iob.AlignByte()
	if err != nil {
		return err
	}
	iob.Binary.PutUint16(iob.Buff[:2], v)
	_, err = iob.Writer.Write(iob.Buff[:2])
	iob.OffsetByte += 2
	return err
}

func (iob *Writer) PutUInt24(v uint32) error {
	err := iob.AlignByte()
	if err != nil {
		return err
	}
	switch iob.Binary {
	case BigEndian:
		iob.Buff[0] = uint8((v << 16) & 0xff)
		iob.Buff[1] = uint8((v << 8) & 0xff)
		iob.Buff[2] = uint8((v) & 0xff)
	case LittleEndian:
		iob.Buff[2] = uint8((v << 16) & 0xff)
		iob.Buff[1] = uint8((v << 8) & 0xff)
		iob.Buff[0] = uint8((v) & 0xff)
	default:
		return fmt.Errorf("PutUInt24 unsupported binary:%#v", iob.Binary)
	}
	_, err = iob.Writer.Write(iob.Buff[:3])
	iob.OffsetByte += 3
	return err
}

func (iob *Writer) PutUInt32(v uint32) error {
	err := iob.AlignByte()
	if err != nil {
		return err
	}
	iob.Binary.PutUint32(iob.Buff[:4], v)
	_, err = iob.Writer.Write(iob.Buff[:4])
	iob.OffsetByte += 4
	return err
}

func (iob *Writer) PutUInt64(v uint64) error {
	err := iob.AlignByte()
	if err != nil {
		return err
	}
	iob.Binary.PutUint64(iob.Buff[:8], v)
	_, err = iob.Writer.Write(iob.Buff[:8])
	iob.OffsetByte += 8
	return err
}

func (iob *Writer) PutUIBit(v uint8) error {
	if v == 1 {
		iob.Buff[0] = iob.Buff[0] | (1 << uint8(7-iob.OffsetBit))
	}
	iob.OffsetBit += 1
	if iob.OffsetBit > 7 {
		_, err := iob.Writer.Write(iob.Buff[:1])
		iob.Buff[0] = 0
		if err != nil {
			return err
		}
		iob.OffsetByte += 1
		iob.OffsetBit -= 8
	}
	return nil
}

func (iob *Writer) PutUIBits_uint8(v uint8, n int) error {
	if n > 8 {
		return fmt.Errorf("PutUIBits_uint8 n:%d > 8", n)
	}
	return iob.PutUIBits_uint64(uint64(v), n)
}

func (iob *Writer) PutUIBits_uint16(v uint16, n int) error {
	if n > 16 {
		return fmt.Errorf("PutUIBits_uint16 n:%d > 16", n)
	}
	return iob.PutUIBits_uint64(uint64(v), n)
}

func (iob *Writer) PutUIBits_uint32(v uint32, n int) error {
	if n > 32 {
		return fmt.Errorf("PutUIBits_uint32 n:%d > 32", n)
	}
	return iob.PutUIBits_uint64(uint64(v), n)
}

func (iob *Writer) PutUIBits_uint64(v uint64, n int) error {
	if n > 64 {
		return fmt.Errorf("PutUIBits_uint64 n:%d > 64", n)
	}
	for i := 0; i < n; i++ {
		b := (v >> uint8(n-1-i)) & 1
		err := iob.PutUIBit(uint8(b))
		if err != nil { // include io.EOF
			return err
		}
	}
	return nil
}
