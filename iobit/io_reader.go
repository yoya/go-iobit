package iobit

/*
 * Copyright 2019/04/06- yoya@awm.jp. All rights reserved.
 */

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

type IOBitReader struct {
	// Read method
	Reader     io.Reader
	Binary     binary.ByteOrder
	OffsetByte uint64
	OffsetBit  uint64
	Buff       []byte
}

func NewReader(r io.Reader, b binary.ByteOrder) *IOBitReader {
	return &IOBitReader{Reader: r, Binary: b,
		OffsetByte: 0, OffsetBit: 0, Buff: make([]byte, 8)}
}

func (r *IOBitReader) Read(buff []byte) (int, error) {
	r.AlignByte()
	r.OffsetByte += uint64(len(buff))
	return r.Reader.Read(buff)
}

func (r *IOBitReader) ReadAll() ([]byte, error) {
	r.AlignByte()
	buff, err := ioutil.ReadAll(r)
	r.OffsetByte += uint64(len(buff))
	return buff, err
}

func (r *IOBitReader) GetOffset() (uint64, uint64) {
	return r.OffsetByte, r.OffsetBit
}

func (r *IOBitReader) AlignByte() error {
	if r.OffsetBit > 0 {
		r.OffsetByte += 1
		r.OffsetBit = 0
	}
	return nil
}

func (r *IOBitReader) GetUInt8() (uint8, error) {
	r.AlignByte()
	_, err := r.Reader.Read(r.Buff[:1])
	if err != nil {
		return 0, err
	}
	r.OffsetByte += 1
	return uint8(r.Buff[0]), nil
}

func (r *IOBitReader) GetUInt16() (uint16, error) {
	r.AlignByte()
	_, err := r.Reader.Read(r.Buff[:2])
	if err != nil {
		return 0, err
	}
	r.OffsetByte += 2
	return r.Binary.Uint16(r.Buff[:2]), nil
}

func (r *IOBitReader) GetUInt24() (uint32, error) {
	r.AlignByte()
	_, err := r.Reader.Read(r.Buff[:3])
	if err != nil {
		return 0, err
	}
	var v uint32
	switch r.Binary {
	case BigEndian:
		v = uint32(r.Buff[0]) << 16
		v += uint32(r.Buff[1]) << 8
		v += uint32(r.Buff[2])
	case LittleEndian:
		v = uint32(r.Buff[2]) << 16
		v += uint32(r.Buff[1]) << 8
		v += uint32(r.Buff[0])
	default:
		return 0, fmt.Errorf("GetUInt24 unsupported binary:%#v", r.Binary)
	}
	r.OffsetByte += 3
	return v, nil
}

func (r *IOBitReader) GetUInt32() (uint32, error) {
	r.AlignByte()
	_, err := r.Reader.Read(r.Buff[:4])
	if err != nil {
		return 0, err
	}
	r.OffsetByte += 4
	return r.Binary.Uint32(r.Buff[:4]), nil
}

func (r *IOBitReader) GetUIn64() (uint64, error) {
	r.AlignByte()
	_, err := r.Reader.Read(r.Buff[:8])
	if err != nil {
		return 0, err
	}
	r.OffsetByte += 8
	return r.Binary.Uint64(r.Buff[:8]), nil
}

func (r *IOBitReader) GetUIBit() (uint8, error) {
	if r.OffsetBit == 0 {
		_, err := r.Reader.Read(r.Buff[:1])
		if err != nil {
			return 0, err
		}
	}
	v := (uint8(r.Buff[0]) >> (7 - r.OffsetBit)) & 1
	r.OffsetBit += 1
	if r.OffsetBit > 7 {
		r.OffsetByte += 1
		r.OffsetBit -= 8
	}
	return v, nil
}

func (r *IOBitReader) GetUIBits_uint8(n int) (uint8, error) {
	if n > 8 {
		return 0, fmt.Errorf("GetUIBits_uint8 n:%d > 8", n)
	}
	v, err := r.GetUIBits_uint64(n)
	return uint8(v), err
}

func (r *IOBitReader) GetUIBits_uint16(n int) (uint16, error) {
	if n > 16 {
		return 0, fmt.Errorf("GetUIBits_uint16 n:%d > 16", n)
	}
	v, err := r.GetUIBits_uint64(n)
	return uint16(v), err
}

func (r *IOBitReader) GetUIBits_uint32(n int) (uint32, error) {
	if n > 32 {
		return 0, fmt.Errorf("GetUIBits_uint32 n:%d > 32", n)
	}
	v, err := r.GetUIBits_uint64(n)
	return uint32(v), err
}

func (r *IOBitReader) GetUIBits_uint64(n int) (uint64, error) {
	if n > 64 {
		return 0, fmt.Errorf("GetUIBits_uint32 n:%d > 64", n)
	}
	var v uint64 = 0
	for i := 0; i < n; i++ {
		v <<= 1
		v1, err := r.GetUIBit()
		if err != nil { // include io.EOF
			return 0, err
		}
		v |= uint64(v1)
	}
	return v, nil
}
