package wasm

import "io"

type I32 interface {
	isI32()
}

func isI32() {}

type localI32 uint32

func (l localI32) isI32() {}

func (l localI32) write(out instCtx) error {
	out.Write([]byte{0x20}) // local.get x
	writeu32(uint32(l), out)
	return nil
}

func (l localI32) set(out io.Writer) error {
	out.Write([]byte{0x21}) // local.set x
	writeu32(uint32(l), out)
	return nil
}

type constUI32 uint32

func (c constUI32) isI32() {}

func (c constUI32) write(out instCtx) error {
	out.Write([]byte{0x41}) // i32.cosnt
	writeu32(uint32(c), out)
	return nil
}
