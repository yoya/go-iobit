# go-iobit - IOBit for Golang

# usage

```
import (
	"fmt"
	"github.com/yoya/go-iobit/iobit"
)
iob := iobit.NewReader(reader, iobit.BigEndian)
v, _ := iob.GetUIBits_uint32(4)
fmt.Printf("%x\n", v)
```

# IOBit API

## IOBitReader

### factory

- func NewReader(reader io.Reader, binary binary.ByteOrder) *Reader

### method

- func (r *Reader) GetUIBits_uint8(n int) (uint8, error)

## IOBitWriter

### factory

- func NewWriter(writer io.Writer, binary binary.ByteOrder) *Writer {

### method

- func (w *Writer) PutUIBits_uint8(v uint8, n int) error
