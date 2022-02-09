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

	localI32() localI32
}

type function struct {
	idx          uint32
	instructions []Instruction
	localF32Cnt  uint32
	localI32Cnt  uint32
}

func (f *function) isExportable() {}

func (f *function) functype() functype {
	return functype{}
}

func (f *function) Body(inst ...Instruction) {
	f.instructions = inst
}

func (f *function) LocalF32() MutableF32 {
	l := localF32(f.localF32Cnt)
	f.localF32Cnt++
	return l
}

func (f *function) localI32() localI32 {
	l := localI32(f.localI32Cnt)
	f.localI32Cnt++
	return l
}

func (f *function) String() string {
	return f.functype().String()
}

func (f *function) encode(out io.Writer) error {
	// write body first to collect additional locals
	body := new(bytes.Buffer)
	for _, inst := range f.instructions {
		if err := inst.write(instCtx{Writer: body, fn: f}); err != nil {
			return err
		}
	}
	body.WriteByte(0x0B) // end

	// write complete function definition
	buf := new(bytes.Buffer)
	// vec(locals)
	writeu32(2, buf)
	// f32 locals
	writeu32(f.localF32Cnt, buf)
	valuetype{numtype: f32}.encode(buf)
	// i32 locals
	writeu32(f.localI32Cnt, buf)
	valuetype{numtype: i32}.encode(buf)
	// expr
	buf.Write(body.Bytes())

	// codesec
	writeu32(uint32(buf.Len()), out)
	out.Write(buf.Bytes())

	return nil
}
