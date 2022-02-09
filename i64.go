package wasm

import "io"

type I64 interface {
	isI64()
}

func isI64() {}

type localI64 uint64

func (l localI64) isI64() {}

func (l localI64) write(out instCtx) error {
	out.Write([]byte{0x20}) // local.get x
	writeu64(uint64(l), out)
	return nil
}

func (l localI64) set(out io.Writer) error {
	out.Write([]byte{0x21}) // local.set x
	writeu64(uint64(l), out)
	return nil
}

type constUI64 uint64

func (c constUI64) isI64() {}

func (c constUI64) write(out instCtx) error {
	constI64.write(out)
	writeu64(uint64(c), out)
	return nil
}
