package wasm

import "io"

type memImport struct{}

func (imp memImport) writeImportDesc(m *Module, out io.Writer) error {
	out.Write([]byte{0x02, 0x00})
	writeu32(1, out)
	return nil
}

func loadF32(offset I32) F32 {
	return opsF32{
		offset,
		op(0x2A),
		u32(0), // static align
		u32(0), // static offset
	}
}
