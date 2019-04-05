/*
 * Copyright 2019/04/06- yoya@awm.jp. All rights reserved.
 */
package main

import (
	"bytes"
	"fmt"
	"github.com/yoya/go-iobit/iobit"
)

func main() {
	buffio := bytes.NewBufferString("ABCDE")
	iob := iobit.Reader(buffio, iobit.BigEndian)
	for {
		v, err := iob.GetUIBits(4)
		if err != nil {
			break
		}
		fmt.Printf("%x ", v)
	}
	fmt.Println("")
}
