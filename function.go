package wasm

import (
	"bytes"
	"io"
)

// Function represents a callable wasm function.
// Functions can be exported and called externally.
type Function interface {
	Exportable

	// Body sets the body of the Function to inst.
	Body(inst ...Instruction)

	// LocalF32 returns a local MutableF32 that can
	// be used inside the function. Using the returned
	// value in Instructions provided to a different
	// Function will result in
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

func (f *function) Body(inst ...Instruction) {
	f.instructions = inst
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
