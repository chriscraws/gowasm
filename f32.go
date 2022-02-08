package wasm

import "io"

// F32 represents a float32 node
type F32 interface {
	Instruction
	isF32()
}

// F32 represents a mutable float32 node
type MutableF32 interface {
	F32
	set(out io.Writer) error
}

// GlobalF32 represents a mutable float32 defined
// in the global scope of the WASM module.
type GlobalF32 interface {
	MutableF32
	Exportable
}

type varF32 struct {
	init float32
	idx  uint32
}

func (v *varF32) isF32() {}

func (v *varF32) incGlobalIndex() {
	v.idx++
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

func (l localF32) isF32() {}

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
