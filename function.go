package wasm

import (
	"bytes"
	"io"
)

type Function interface {
	Exportable
	AddInstructions(inst ...Instruction)
	LocalF32() MutableF32
}

type function struct {
	idx          uint32
	instructions []Instruction
	localCnt     uint32
}

func (f *function) isExportable() {}

func (f *function) functype() functype {
	return functype{}
}

func (f *function) AddInstructions(inst ...Instruction) {
	f.instructions = append(f.instructions, inst...)
}

func (f *function) LocalF32() MutableF32 {
	l := localF32(f.localCnt)
	f.localCnt++
	return l
}

func (f *function) String() string {
	return f.functype().String()
}

func (f *function) encode(out io.Writer) error {
	buf := new(bytes.Buffer)
	// vec(locals)
	writeu32(1, buf)
	// f32 locals
	writeu32(f.localCnt, buf)
	valuetype{numtype: f32}.encode(buf)
	// for _, l := range f.locals {
	// 	// init expr
	// 	if err := l.write(buf); err != nil {
	// 		return err
	// 	}
	// }
	// expr
	for _, inst := range f.instructions {
		if err := inst.write(buf); err != nil {
			return err
		}
	}
	buf.WriteByte(0x0B) // end

	// codesec
	writeu32(uint32(buf.Len()), out)
	out.Write(buf.Bytes())
	return nil
}
