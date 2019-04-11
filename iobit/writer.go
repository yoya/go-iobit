package iobit

/*
 * Copyright 2019/04/07- yoya@awm.jp. All rights reserved.
 */

type Writer interface {
	Write(buff []byte) (int, error)
	GetOffset() (uint64, uint64)
	AlignByte() error
	PutUInt8(v uint8) error
	PutUInt16(v uint16) error
	PutUInt24(v uint32) error
	PutUInt32(v uint32) error
	PutUInt64(v uint64) error
	PutUIBit(v uint8) error
	PutUIBits_uint8(v uint8, n int) error
	PutUIBits_uint16(v uint16, n int) error
	PutUIBits_uint32(v uint32, n int) error
	PutUIBits_uint64(v uint64, n int) error
}
