package iobit

/*
 * Copyright 2019/04/07- yoya@awm.jp. All rights reserved.
 */

type Writer interface {
	Write(buff []byte) (int, error)
	GetOffset() (uint64, uint64)
	AlignByte()
	PutUInt8(v uint8)
	PutUInt16(v uint16)
	PutUInt24(v uint32)
	PutUInt32(v uint32)
	PutUInt64(v uint64)
	PutUIBit(v uint8)
	PutUIBits_uint8(v uint8, n int)
	PutUIBits_uint16(v uint16, n int)
	PutUIBits_uint32(v uint32, n int)
	PutUIBits_uint64(v uint64, n int)
	PutBytes([]byte)
	PutString(string)
	GetLastError() error
}
