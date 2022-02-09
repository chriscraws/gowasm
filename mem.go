package wasm

import "io"

type memImport struct{}

func (imp memImport) writeImportDesc(out io.Writer) error {
	out.Write([]byte{0x02, 0x00, 0x00})
	return nil
}
