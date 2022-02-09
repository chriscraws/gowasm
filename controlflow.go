package wasm

import "fmt"

type controlInst byte

func (ci controlInst) write(c instCtx) error {
	c.Write([]byte{byte(ci), 0x40})
	return nil
}

const (
	unreachableCI op = iota
	nopCI
)

const (
	blockCI controlInst = 0x02 + iota
	loopCI
	ifElseCI
)

const (
	branchCI op = 0x0C + iota
	branchIfCI
	branchTableCI
	returnCI
)

const elseCI op = 0x05
const endCI op = 0x0B

const (
	dropOp op = 0x1A + iota
	selectOp
)

type branch int

func (b branch) write(c instCtx) error {
	branchCI.write(c)
	writeu32(uint32(b), c)
	return nil
}

type branchIf int

func (b branchIf) write(c instCtx) error {
	branchIfCI.write(c)
	writeu32(uint32(b), c)
	return nil
}

type ForRangeF32 struct {
	Begin F32
	End   F32
	Inc   F32
	Do    func(index F32) []Instruction
}

func (fr ForRangeF32) write(c instCtx) error {
	if fr.Inc == nil {
		fr.Inc = ConstF32(1)
	}
	// create locals
	idx := c.fn.LocalF32()
	end := c.fn.LocalF32()
	inc := c.fn.LocalF32()
	body := []Instruction{
		// assign locals
		AssignF32(idx, fr.Begin),
		AssignF32(end, fr.End),
		AssignF32(inc, fr.Inc),

		// ok, time to loop
		blockCI,
		loopCI,

		// check if we're out of bounds
		inc,
		ConstF32(0),
		gef32,

		ifElseCI,
		idx,
		end,
		gef32, // if inc is positive, check if idx >= end
		branchIf(2),
		elseCI,
		idx,
		end,
		lef32, // if inc is negative check if idx <= end
		branchIf(2),
		endCI,
	}

	// ok we're still in the loop, call user code now
	body = append(body, fr.Do(idx)...)

	// user code is done, lets wrap up the loop
	body = append(body,
		AssignF32(idx, AddF32(idx, inc)),
		// return to beginning of loop by default
		branch(0),

		endCI, // end loop
		endCI, // end outer block
	)

	for _, inst := range body {
		if err := inst.write(c); err != nil {
			return fmt.Errorf("failure in for range: %s", err)
		}
	}
	return nil
}
