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

type IOReader struct {
	// Read method
	reader     io.Reader
	binary     binary.ByteOrder
	offsetByte uint64
	offsetBit  uint64
	buff       []byte
	lastError  error
}

func NewIOReader(r io.Reader, b binary.ByteOrder) *IOReader {
	return &IOReader{reader: r, binary: b,
		offsetByte: 0, offsetBit: 0, buff: make([]byte, 8),
		lastError: nil}
}

func (r *IOReader) Read(buff []byte) (int, error) {
	r.AlignByte()
	if r.lastError != nil {
		return 0, r.lastError
	}
	var n int
	n, r.lastError = r.reader.Read(buff)
	r.offsetByte += uint64(n)
	return n, r.lastError
}

func (r *IOReader) ReadAll() ([]byte, error) {
	r.AlignByte()
	if r.lastError != nil {
		return nil, r.lastError
	}
	var buff []byte
	buff, r.lastError = ioutil.ReadAll(r)
	r.offsetByte += uint64(len(buff))
	return buff, nil
}

func (r *IOReader) GetOffset() (uint64, uint64) {
	return r.offsetByte, r.offsetBit
}

func (r *IOReader) AlignByte() {
	if r.lastError != nil {
		return
	}
	if r.offsetBit > 0 {
		r.offsetByte += 1
		r.offsetBit = 0
	}
}

func (r *IOReader) GetUInt8() uint8 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	var n int
	n, r.lastError = r.reader.Read(r.buff[:1])
	r.offsetByte += uint64(n)
	if r.lastError != nil {
		return 0
	}
	return uint8(r.buff[0])
}

func (r *IOReader) GetUInt16() uint16 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	var n int
	n, r.lastError = r.reader.Read(r.buff[:2])
	r.offsetByte += uint64(n)
	if r.lastError != nil {
		return 0
	}
	return r.binary.Uint16(r.buff[:2])
}

func (r *IOReader) GetUInt24() uint32 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	var n int
	n, r.lastError = r.reader.Read(r.buff[:3])
	r.offsetByte += uint64(n)
	if r.lastError != nil {
		return 0
	}
	var v uint32
	switch r.binary {
	case BigEndian:
		v = uint32(r.buff[0]) << 16
		v += uint32(r.buff[1]) << 8
		v += uint32(r.buff[2])
	case LittleEndian:
		v = uint32(r.buff[2]) << 16
		v += uint32(r.buff[1]) << 8
		v += uint32(r.buff[0])
	default:
		r.lastError = fmt.Errorf("GetUInt24 unsupported binary:%#v", r.binary)
		v = 0
	}
	return v
}

func (r *IOReader) GetUInt32() uint32 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	var n int
	n, r.lastError = r.reader.Read(r.buff[:4])
	r.offsetByte += uint64(n)
	if r.lastError != nil {
		return 0
	}
	return r.binary.Uint32(r.buff[:4])
}

func (r *IOReader) GetUIn64() uint64 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	var n int
	n, r.lastError = r.reader.Read(r.buff[:8])
	r.offsetByte += uint64(n)
	if r.lastError != nil {
		return 0
	}
	return r.binary.Uint64(r.buff[:8])
}

func (r *IOReader) GetUIBit() uint8 {
	if r.lastError != nil {
		return 0
	}
	if r.offsetBit == 0 {
		var n int
		n, r.lastError = r.reader.Read(r.buff[:1])
		r.offsetByte += uint64(n)
		if r.lastError != nil {
			return 0
		}
	}
	v := (uint8(r.buff[0]) >> (7 - r.offsetBit)) & 1
	r.offsetBit += 1
	if r.offsetBit > 7 {
		r.offsetByte += 1
		r.offsetBit -= 8
	}
	return v
}

func (r *IOReader) GetUIBits_uint8(n int) uint8 {
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

func (r *IOReader) GetUIBits_uint16(n int) uint16 {
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

func (r *IOReader) GetUIBits_uint32(n int) uint32 {
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

func (r *IOReader) GetUIBits_uint64(n int) uint64 {
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

func (r *IOReader) GetLastError() error {
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
