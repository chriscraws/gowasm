package wasm

import "io"

type I32 interface {
	Instruction
	isI32()
}

func isI32() {}

type opsI32 ops

func (o opsI32) isI32() {}

func (o opsI32) write(out instCtx) error {
	return ops(o).write(out)
}

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

type assignI32 struct {
	dst localI32
	v   I32
}

func (a assignI32) write(out instCtx) error {
	if err := a.v.write(out); err != nil {
		return err
	}
	if err := a.dst.set(out); err != nil {
		return err
	}
	return nil
}

func addi32(a, b I32) I32 {
	return opsI32{
		a,
		b,
		op(0x6A),
	}
}

func muli32(a, b I32) I32 {
	return opsI32{
		a,
		b,
		op(0x6C),
	}
}

func castF32I32(a F32) I32 {
	return opsI32{
		a,
		truncf32ui32,
	}
}
