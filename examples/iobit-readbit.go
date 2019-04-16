/*
 * Copyright 2019/04/08- yoya@awm.jp. All rights reserved.
 */
package main

import (
	"bytes"
	"fmt"
	"github.com/yoya/go-iobit/iobit"
	"io"
	"os"
)

func main() {
	for _, arg := range os.Args[1:] {
		buffio := bytes.NewBufferString(arg)
		iob := iobit.NewReader(buffio, iobit.BigEndian)
		for {
			v := iob.GetUIBit()
			err := iob.GetLastError()
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			fmt.Printf("%x", v)
			_, bitOff := iob.GetOffset()
			if bitOff == 0 {
				fmt.Printf(" ")
			}
		}
		fmt.Println("")
	}

}
