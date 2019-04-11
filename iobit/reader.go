package iobit

/*
 * Copyright 2019/04/06- yoya@awm.jp. All rights reserved.
 */

type Reader interface {
	Read(buff []byte) (int, error)
	ReadAll() ([]byte, error)
	GetOffset() (uint64, uint64)
	AlignByte() error
	GetUInt8() (uint8, error)
	GetUInt16() (uint16, error)
	GetUInt24() (uint32, error)
	GetUInt32() (uint32, error)
	GetUIn64() (uint64, error)
	GetUIBit() (uint8, error)
	GetUIBits_uint8(n int) (uint8, error)
	GetUIBits_uint16(n int) (uint16, error)
	GetUIBits_uint32(n int) (uint32, error)
	GetUIBits_uint64(n int) (uint64, error)
}
