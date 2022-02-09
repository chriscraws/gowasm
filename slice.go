package wasm

import "io"

// SliceF32 is a contiguous slice of float32 values located in wasm memory.
type SliceF32 interface {
	// LengthF32 returns the number of float32 values.
	LengthF32() F32
	// IndexF32 returns the float32 value at index i.
	IndexF32(i F32) F32
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

func (s *sliceF32) writeImportDesc(m *Module, out io.Writer) error {
	out.Write([]byte{0x03})
	return globaltype{
		mutable: false,
		valuetype: valuetype{
			numtype: i64,
		},
	}.encode(out)
}

func (s *sliceF32) LengthF32() F32 {
	return opsF32{
		// get global
		globalGet,
		u32(s.idx),

		// get higher order bits
		constUI64(32),
		shiftRightUI64,

		// convert the length to an F32
		converti64uF32,
	}
}

func (s *sliceF32) offsetI32() I32 {
	return opsI32{
		// get global
		globalGet,
		u32(s.idx),
		// get lower order bits
		wrapi64I32,
	}
}

func (s *sliceF32) indexI32(i I32) F32 {
	return loadF32(
		addi32(
			s.offsetI32(),
			muli32(i, constUI32(4)),
		),
	)
}

func (s *sliceF32) IndexF32(i F32) F32 {
	return loadF32(
		addi32(
			s.offsetI32(),
			muli32(castF32I32(i), constUI32(4)),
		),
	)
}

// SliceF32RangeF32 is an instruction that runs the instructions
// returned by Do for each value in the slice from index Begin
// to index End-1
type SliceF32RangeF32 struct {
	Slice SliceF32
	Begin F32
	End   F32
	Do    func(v F32) []Instruction
}

func (s SliceF32RangeF32) write(c instCtx) error {
	end := c.fn.localI32()
	idx := c.fn.localI32()
	if s.Begin == nil {
		s.Begin = ConstF32(0)
	}
	if s.End == nil {
		s.End = s.Slice.LengthF32()
	}
	body := ops{
		// begin = uint32(begin)
		assignI32{dst: idx, v: castF32I32(s.Begin)},
		// end = uint32(end)
		assignI32{dst: end, v: castF32I32(s.End)},
		// idx = uint32(0)

		blockCI,
		loopCI,

		// if (idx >= end) break;
		idx,
		end,
		geUI32,
		branchIfCI,
		u32(1),
	}
	slice := s.Slice.(*sliceF32)
	body = append(body, s.Do(slice.indexI32(idx))...)
	body = append(body,
		// idx++
		assignI32{dst: idx, v: addi32(idx, constUI32(1))},

		// continue
		branchCI,
		u32(0),

		endCI, // loopCI
		endCI, // blockCI
	)
	return body.write(c)
}
