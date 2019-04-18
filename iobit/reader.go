package iobit

/*
 * Copyright 2019/04/06- yoya@awm.jp. All rights reserved.
 */

// imitate io package
import (
	"errors"
)

var EOF = errors.New("EOF")
var ErrUnexpectedEOF = errors.New("unexpected EOF")

type Reader interface {
	Read(buff []byte) (int, error)
	ReadAll() ([]byte, error)
	ReadUntil(byte, bool) ([]byte, error)
	GetOffset() (uint64, uint64)
	AlignByte()
	GetUInt8() uint8
	GetUInt16() uint16
	GetUInt24() uint32
	GetUInt32() uint32
	GetUInt64() uint64
	GetUIBit() uint8
	GetUIBits_uint8(n int) uint8
	GetUIBits_uint16(n int) uint16
	GetUIBits_uint32(n int) uint32
	GetUIBits_uint64(n int) uint64
	GetBytes(n int) []byte
	GetString(n int) string
	GetLastError() error
}
