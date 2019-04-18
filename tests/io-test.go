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
	var iob iobit.Reader = iobit.NewIOReader(buffio, iobit.BigEndian)
	for {
		v := iob.GetUIBits_uint32(4)
		if iob.GetLastError() != nil {

			break
		}
		fmt.Printf("%x ", v)
	}
	fmt.Println("")
	var iobw iobit.Writer = iobit.NewIOWriter(os.Stdout, iobit.BigEndian)
	for i := 1; i < 6; i++ {
		iobw.PutUIBits_uint8(4, 4)
		iobw.PutUIBits_uint8(uint8(i), 4)
	}
	fmt.Println()
	if iobw.GetLastError() != nil {
		fmt.Println(iobw.GetLastError())
	}
}
