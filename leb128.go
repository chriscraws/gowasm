package wasm

import (
	"io"
)

func writeu32(v uint32, out io.Writer) {
	for {
		b := byte(v & 0b01111111)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		out.Write([]byte{b})
		if b&0x80 == 0 {
			return
		}
	}
}

type u32 uint32

func (u u32) write(c instCtx) error {
	writeu32(uint32(u), c)
	return nil
}
