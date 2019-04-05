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
		OffsetByte: 0, OffsetBit: 0, Buff: make([]byte, 1)}
}

func (iob *IOBitReader) Read(buff []byte) (int, error) {
	return iob.Reader.Read(buff)
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
	var v uint8 = 0
	for i := 0; i < n; i++ {
		v <<= 1
		v1, err := iob.GetUIBit()
		if err != nil { // include io.EOF
			return 0, err
		}
		v |= uint8(v1)
	}
	return v, nil
}

func (iob *IOBitReader) GetUIBits_uint16(n int) (uint16, error) {
	if n > 16 {
		return 0, fmt.Errorf("GetUIBits_uint16 n:%d > 16", n)
	}
	var v uint16 = 0
	for i := 0; i < n; i++ {
		v <<= 1
		v1, err := iob.GetUIBit()
		if err != nil { // include io.EOF
			return 0, err
		}
		v |= uint16(v1)
	}
	return v, nil
}

func (iob *IOBitReader) GetUIBits_uint32(n int) (uint32, error) {
	if n > 32 {
		return 0, fmt.Errorf("GetUIBits_uint32 n:%d > 32", n)
	}
	var v uint32 = 0
	for i := 0; i < n; i++ {
		v <<= 1
		v1, err := iob.GetUIBit()
		if err != nil { // include io.EOF
			return 0, err
		}
		v |= uint32(v1)
	}
	return v, nil
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
