package iobit

/*
 * Copyright 2019/04/06- yoya@awm.jp. All rights reserved.
 */

import (
	"encoding/binary"
	"fmt"
	"io"
)

type IOBitReader struct {
	// Read method
	Reader     io.Reader
	Binary     binary.ByteOrder
	OffsetByte uint64
	OffsetBit  uint64
	Buff       []byte
}

func Reader(reader io.Reader, binary binary.ByteOrder) *IOBitReader {
	return &IOBitReader{Reader: reader, Binary: binary,
		OffsetByte: 0, OffsetBit: 0, Buff: make([]byte, 8)}
}

func (iob *IOBitReader) Read(buff []byte) (int, error) {
	iob.OffsetByte += uint64(len(buff))
	return iob.Reader.Read(buff)
}

func (iob *IOBitReader) GetOffset() (uint64, uint64) {
	return iob.OffsetByte, iob.OffsetBit
}

func (iob *IOBitReader) AlignByte() error {
	if iob.OffsetBit > 0 {
		iob.OffsetByte += 1
		iob.OffsetBit = 0
	}
	return nil
}

func (iob *IOBitReader) GetUInt8() (uint8, error) {
	iob.AlignByte()
	_, err := iob.Reader.Read(iob.Buff[:1])
	if err != nil {
		return 0, err
	}
	iob.OffsetByte += 1
	return uint8(iob.Buff[0]), nil
}

func (iob *IOBitReader) GetUInt16() (uint16, error) {
	iob.AlignByte()
	_, err := iob.Reader.Read(iob.Buff[:2])
	if err != nil {
		return 0, err
	}
	iob.OffsetByte += 2
	return iob.Binary.Uint16(iob.Buff[:2]), nil
}

func (iob *IOBitReader) GetUInt24() (uint32, error) {
	iob.AlignByte()
	_, err := iob.Reader.Read(iob.Buff[:3])
	if err != nil {
		return 0, err
	}
	var v uint32
	switch iob.Binary {
	case BigEndian:
		v = uint32(iob.Buff[0]) << 16
		v += uint32(iob.Buff[1]) << 8
		v += uint32(iob.Buff[2])
	case LittleEndian:
		v = uint32(iob.Buff[2]) << 16
		v += uint32(iob.Buff[1]) << 8
		v += uint32(iob.Buff[0])
	default:
		return 0, fmt.Errorf("GetUInt24 unsupported binary:%#v", iob.Binary)
	}
	iob.OffsetByte += 3
	return v, nil
}

func (iob *IOBitReader) GetUInt32() (uint32, error) {
	iob.AlignByte()
	_, err := iob.Reader.Read(iob.Buff[:4])
	if err != nil {
		return 0, err
	}
	iob.OffsetByte += 4
	return iob.Binary.Uint32(iob.Buff[:4]), nil
}

func (iob *IOBitReader) GetUIn64() (uint64, error) {
	iob.AlignByte()
	_, err := iob.Reader.Read(iob.Buff[:8])
	if err != nil {
		return 0, err
	}
	iob.OffsetByte += 8
	return iob.Binary.Uint64(iob.Buff[:8]), nil
}

func (iob *IOBitReader) GetUIBit() (uint8, error) {
	if iob.OffsetBit == 0 {
		_, err := iob.Reader.Read(iob.Buff[:1])
		if err != nil {
			return 0, err
		}
	}
	v := (uint8(iob.Buff[0]) >> (7 - iob.OffsetBit)) & 1
	iob.OffsetBit += 1
	if iob.OffsetBit > 7 {
		iob.OffsetByte += 1
		iob.OffsetBit -= 8
	}
	return v, nil
}

func (iob *IOBitReader) GetUIBits_uint8(n int) (uint8, error) {
	if n > 8 {
		return 0, fmt.Errorf("GetUIBits_uint8 n:%d > 8", n)
	}
	v, err := iob.GetUIBits_uint64(n)
	return uint8(v), err
}

func (iob *IOBitReader) GetUIBits_uint16(n int) (uint16, error) {
	if n > 16 {
		return 0, fmt.Errorf("GetUIBits_uint16 n:%d > 16", n)
	}
	v, err := iob.GetUIBits_uint64(n)
	return uint16(v), err
}

func (iob *IOBitReader) GetUIBits_uint32(n int) (uint32, error) {
	if n > 32 {
		return 0, fmt.Errorf("GetUIBits_uint32 n:%d > 32", n)
	}
	v, err := iob.GetUIBits_uint64(n)
	return uint32(v), err
}

func (iob *IOBitReader) GetUIBits_uint64(n int) (uint64, error) {
	if n > 64 {
		return 0, fmt.Errorf("GetUIBits_uint32 n:%d > 64", n)
	}
	var v uint64 = 0
	for i := 0; i < n; i++ {
		v <<= 1
		v1, err := iob.GetUIBit()
		if err != nil { // include io.EOF
			return 0, err
		}
		v |= uint64(v1)
	}
	return v, nil
}
