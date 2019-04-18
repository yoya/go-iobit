package iobit

/*
 * Copyright 2019/04/06- yoya@awm.jp. All rights reserved.
 */

import (
	"encoding/binary"
	"fmt"
)

type BytesReader struct {
	// Read method
	bytes      []byte
	binary     binary.ByteOrder
	offsetByte uint64
	offsetBit  uint64
	lastError  error
}

func NewBytesReader(bytes []byte, b binary.ByteOrder) *BytesReader {
	return &BytesReader{bytes: bytes, binary: b,
		offsetByte: 0, offsetBit: 0,
		lastError: nil}
}

func (r *BytesReader) hasNext(n int) bool {
	remain_len := int(len(r.bytes)) - int(r.offsetByte)
	if remain_len < n {
		return false
	}
	return true
}

func (r *BytesReader) Read(buff []byte) (int, error) {
	r.AlignByte()
	if r.lastError != nil {
		return 0, r.GetLastError()
	}
	buff_len := int(len(buff))
	remain_len := int(len(r.bytes)) - int(r.offsetByte)
	if remain_len == 0 {
		r.lastError = EOF
		return 0, r.GetLastError()
	}

	var n int
	if remain_len < buff_len {
		r.lastError = ErrUnexpectedEOF
		n = remain_len
	} else {
		n = buff_len
	}
	copy(buff, r.bytes[r.offsetByte:int(r.offsetByte)+n])
	r.offsetByte += uint64(n)
	return n, r.GetLastError()
}

func (r *BytesReader) ReadAll() ([]byte, error) {
	r.AlignByte()
	if r.lastError != nil {
		return nil, r.GetLastError()
	}
	remain_len := int(len(r.bytes)) - int(r.offsetByte)
	buff := make([]byte, remain_len)
	_, r.lastError = r.Read(buff)
	return buff, r.GetLastError()
}

func (r *BytesReader) ReadUntil(elim byte, return_include_elim bool) ([]byte, error) {
	r.lastError = fmt.Errorf("%s", "ReadUntil: Not implemented yet")
	return make([]byte, 0), r.GetLastError()
}

func (r *BytesReader) GetOffset() (uint64, uint64) {
	return r.offsetByte, r.offsetBit
}

func (r *BytesReader) AlignByte() {
	if r.lastError != nil {
		return
	}
	if r.offsetBit > 0 {
		r.offsetByte += 1
		r.offsetBit = 0
	}
}

func (r *BytesReader) GetUInt8() uint8 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	if r.hasNext(1) == false {
		r.lastError = EOF
		return 0
	}
	v := r.bytes[r.offsetByte]
	r.offsetByte += uint64(1)
	return uint8(v)
}

func (r *BytesReader) GetUInt16() uint16 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	if r.hasNext(2) == false {
		if r.hasNext(1) == false {
			r.lastError = EOF
		} else {
			r.lastError = ErrUnexpectedEOF
		}
		return 0
	}
	v := r.binary.Uint16(r.bytes[r.offsetByte : r.offsetByte+2])
	r.offsetByte += uint64(2)
	return v
}

func (r *BytesReader) GetUInt24() uint32 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	if r.hasNext(3) == false {
		if r.hasNext(1) == false {
			r.lastError = EOF
		} else {
			r.lastError = ErrUnexpectedEOF
		}
		return 0
	}
	buff := r.bytes[r.offsetByte : r.offsetByte+3]
	r.offsetByte += uint64(3)
	var v uint32
	switch r.binary {
	case BigEndian:
		v = uint32(buff[0]) << 16
		v += uint32(buff[1]) << 8
		v += uint32(buff[2])
	case LittleEndian:
		v = uint32(buff[2]) << 16
		v += uint32(buff[1]) << 8
		v += uint32(buff[0])
	default:
		r.lastError = fmt.Errorf("GetUInt24 unsupported binary:%#v", r.binary)
		v = 0
	}
	return v
}

func (r *BytesReader) GetUInt32() uint32 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	if r.hasNext(4) == false {
		if r.hasNext(1) == false {
			r.lastError = EOF
		} else {
			r.lastError = ErrUnexpectedEOF
		}
		return 0
	}
	v := r.binary.Uint32(r.bytes[r.offsetByte : r.offsetByte+4])
	r.offsetByte += uint64(4)
	return v
}

func (r *BytesReader) GetUInt64() uint64 {
	r.AlignByte()
	if r.lastError != nil {
		return 0
	}
	if r.hasNext(8) == false {
		if r.hasNext(1) == false {
			r.lastError = EOF
		} else {
			r.lastError = ErrUnexpectedEOF
		}
		return 0
	}
	v := r.binary.Uint64(r.bytes[r.offsetByte : r.offsetByte+8])
	r.offsetByte += uint64(4)
	return v
}

func (r *BytesReader) GetUIBit() uint8 {
	if r.lastError != nil {
		return 0
	}
	if r.offsetBit == 0 {
		if r.hasNext(1) == false {
			r.lastError = EOF
		}
	} else {
		if r.hasNext(0) == false {
			r.lastError = ErrUnexpectedEOF
		}
	}
	if r.lastError != nil {
		return 0
	}
	v := (uint8(r.bytes[r.offsetByte]) >> (7 - r.offsetBit)) & 1
	r.offsetBit += 1
	if r.offsetBit > 7 {
		r.offsetByte += 1
		r.offsetBit -= 8
	}
	return v
}

func (r *BytesReader) GetUIBits_uint8(n int) uint8 {
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

func (r *BytesReader) GetUIBits_uint16(n int) uint16 {
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

func (r *BytesReader) GetUIBits_uint32(n int) uint32 {
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

func (r *BytesReader) GetUIBits_uint64(n int) uint64 {
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
		if r.lastError != nil { // include EOF
			return 0
		}
		v |= uint64(v1)
	}
	return v
}

func (r *BytesReader) GetBytes(n int) []byte {
	r.AlignByte()
	if r.lastError != nil {
		return nil
	}
	if r.hasNext(n) == false {
		if r.hasNext(1) == false {
			r.lastError = EOF
			return nil
		}
		r.lastError = ErrUnexpectedEOF
	}
	buff := r.bytes[r.offsetByte : int(r.offsetByte)+n]
	r.offsetByte += uint64(len(buff))
	return buff
}
func (r *BytesReader) GetString(n int) string {
	r.AlignByte()
	if r.lastError != nil {
		return ""
	}
	buff := r.GetBytes(n)
	if r.lastError != nil {
		return ""
	}
	return string(buff)
}

func (r *BytesReader) GetLastError() error {
	if r.lastError == nil {
		return nil
	}
	if r.lastError == EOF {
		return r.lastError
	}
	if r.lastError == ErrUnexpectedEOF {
		return ErrUnexpectedEOF
	}
	err := fmt.Errorf("%s in offset(%d:%d)",
		r.lastError, r.offsetByte, r.offsetBit)
	return err
}
