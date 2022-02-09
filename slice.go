package wasm

import "io"

type SliceF32 interface {
}

type sliceF32 struct {
	idx uint32
}

func (s *sliceF32) incGlobalIndex() {
	s.idx++
}

func (s *sliceF32) setGlobalIndex(i uint32) {
	s.idx = i
}

func (s *sliceF32) globalIndex() uint32 {
	return s.idx
}

func (s *sliceF32) writeImportDesc(out io.Writer) error {
	out.Write([]byte{0x03})
	return globaltype{
		mutable: true,
		valuetype: valuetype{
			numtype: i64,
		},
	}.encode(out)
}
