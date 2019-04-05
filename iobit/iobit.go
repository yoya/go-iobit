package iobit

/*
 * Copyright 2019/04/6- yoya@awm.jp. All rights reserved.
 */

import (
	"encoding/binary"
	"io"
)

var BigEndian binary.ByteOrder = binary.BigEndian
var LittleEndian binary.ByteOrder = binary.LittleEndian

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

func (iob *IOBitReader) GetUIBits(n int) (uint64, error) {
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
