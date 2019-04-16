/*
 * Copyright 2019/04/06- yoya@awm.jp. All rights reserved.
 */
package main

import (
	"bytes"
	"fmt"
	"github.com/yoya/go-iobit/iobit"
	"os"
)

func main() {
	buffio := bytes.NewBufferString("ABCDE")
	var iob iobit.Reader = iobit.NewReader(buffio, iobit.BigEndian)
	for {
		v := iob.GetUIBits_uint32(4)
		if iob.GetLastError() != nil {

			break
		}
		fmt.Printf("%x ", v)
	}
	fmt.Println("")
	var iobw iobit.Writer = iobit.NewWriter(os.Stdout, iobit.BigEndian)
	for i := 1; i < 6; i++ {
		_ = iobw.PutUIBits_uint8(4, 4)
		_ = iobw.PutUIBits_uint8(uint8(i), 4)
	}
	fmt.Println("")
}
