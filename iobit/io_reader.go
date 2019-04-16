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
	lastError  error
}

func NewReader(r io.Reader, b binary.ByteOrder) *IOBitReader {
	return &IOBitReader{Reader: r, Binary: b,
		OffsetByte: 0, OffsetBit: 0, Buff: make([]byte, 8),
		lastError: nil}
}

func (r *IOBitReader) Read(buff []byte) (int, error) {
	if r.lastError != nil {
		return 0, r.lastError
	}
	r.AlignByte()
	r.OffsetByte += uint64(len(buff))
	var n int
	n, r.lastError = r.Reader.Read(buff)
	return n, r.lastError
}

func (r *IOBitReader) ReadAll() ([]byte, error) {
	if r.lastError != nil {
		return nil, r.lastError
	}
	r.AlignByte()
	var buff []byte
	buff, r.lastError = ioutil.ReadAll(r)
	r.OffsetByte += uint64(len(buff))
	return buff, nil
}

func (r *IOBitReader) GetOffset() (uint64, uint64) {
	return r.OffsetByte, r.OffsetBit
}

func (r *IOBitReader) AlignByte() {
	if r.lastError != nil {
		return
	}
	if r.OffsetBit > 0 {
		r.OffsetByte += 1
		r.OffsetBit = 0
	}
}

func (r *IOBitReader) GetUInt8() uint8 {
	if r.lastError != nil {
		return 0
	}
	r.AlignByte()
	_, r.lastError = r.Reader.Read(r.Buff[:1])
	if r.lastError != nil {
		return 0
	}
	r.OffsetByte += 1
	return uint8(r.Buff[0])
}

func (r *IOBitReader) GetUInt16() uint16 {
	if r.lastError != nil {
		return 0
	}
	r.AlignByte()
	_, r.lastError = r.Reader.Read(r.Buff[:2])
	if r.lastError != nil {
		return 0
	}
	r.OffsetByte += 2
	return r.Binary.Uint16(r.Buff[:2])
}

func (r *IOBitReader) GetUInt24() uint32 {
	if r.lastError != nil {
		return 0
	}
	r.AlignByte()
	_, r.lastError = r.Reader.Read(r.Buff[:3])
	if r.lastError != nil {
		return 0
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
		r.lastError = fmt.Errorf("GetUInt24 unsupported binary:%#v", r.Binary)
		v = 0
	}
	r.OffsetByte += 3
	return v
}

func (r *IOBitReader) GetUInt32() uint32 {
	if r.lastError != nil {
		return 0
	}
	r.AlignByte()
	_, r.lastError = r.Reader.Read(r.Buff[:4])
	if r.lastError != nil {
		return 0
	}
	r.OffsetByte += 4
	return r.Binary.Uint32(r.Buff[:4])
}

func (r *IOBitReader) GetUIn64() uint64 {
	if r.lastError != nil {
		return 0
	}
	r.AlignByte()
	_, r.lastError = r.Reader.Read(r.Buff[:8])
	if r.lastError != nil {
		return 0
	}
	r.OffsetByte += 8
	return r.Binary.Uint64(r.Buff[:8])
}

func (r *IOBitReader) GetUIBit() uint8 {
	if r.lastError != nil {
		return 0
	}
	if r.OffsetBit == 0 {
		_, r.lastError = r.Reader.Read(r.Buff[:1])
		if r.lastError != nil {
			return 0
		}
	}
	v := (uint8(r.Buff[0]) >> (7 - r.OffsetBit)) & 1
	r.OffsetBit += 1
	if r.OffsetBit > 7 {
		r.OffsetByte += 1
		r.OffsetBit -= 8
	}
	return v
}

func (r *IOBitReader) GetUIBits_uint8(n int) uint8 {
	if r.lastError != nil {
		return 0
	}
	if n > 8 {
		r.lastError = fmt.Errorf("GetUIBits_uint8 n:%d > 8", n)
		return 0
	}
	v := r.GetUIBits_uint64(n)
	return uint8(v)
}

func (r *IOBitReader) GetUIBits_uint16(n int) uint16 {
	if r.lastError != nil {
		return 0
	}
	if n > 16 {
		r.lastError = fmt.Errorf("GetUIBits_uint16 n:%d > 16", n)
		return 0
	}
	v := r.GetUIBits_uint64(n)
	return uint16(v)
}

func (r *IOBitReader) GetUIBits_uint32(n int) uint32 {
	if r.lastError != nil {
		return 0
	}
	if n > 32 {
		r.lastError = fmt.Errorf("GetUIBits_uint32 n:%d > 32", n)
		return 0
	}
	v := r.GetUIBits_uint64(n)
	return uint32(v)
}

func (r *IOBitReader) GetUIBits_uint64(n int) uint64 {
	if r.lastError != nil {
		return 0
	}
	if n > 64 {
		r.lastError = fmt.Errorf("GetUIBits_uint32 n:%d > 64", n)
		return 0
	}
	var v uint64 = 0
	for i := 0; i < n; i++ {
		v <<= 1
		var v1 uint8
		v1 = r.GetUIBit()
		if r.lastError != nil { // include io.EOF
			return 0
		}
		v |= uint64(v1)
	}
	return v
}

func (r *IOBitReader) GetLastError() error {
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
