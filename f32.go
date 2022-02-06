package wasm

import "io"

type F32 interface {
	Instruction
}

type MutableF32 interface {
	F32
	set(out io.Writer) error
}

type GlobalF32 interface {
	MutableF32
	Exportable
}

type varF32 struct {
	init float32
	idx  uint32
}

func (v *varF32) write(out io.Writer) error {
	out.Write([]byte{0x23}) // global.get x
	writeu32(v.idx, out)
	return nil
}

func (v *varF32) set(out io.Writer) error {
	out.Write([]byte{0x24}) // global.set x
	writeu32(v.idx, out)
	return nil
}

func (v *varF32) isExportable() {}

type localF32 uint32

func (l localF32) write(out io.Writer) error {
	out.Write([]byte{0x20}) // local.get x
	writeu32(uint32(l), out)
	return nil
}

func (l localF32) set(out io.Writer) error {
	out.Write([]byte{0x21}) // local.set x
	writeu32(uint32(l), out)
	return nil
}
