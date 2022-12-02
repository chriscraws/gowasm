package main

import (
	"os"

	wasm "github.com/chriscraws/gowasm"
)

func main() {
	m := new(wasm.Module)

	end := m.GlobalF32(0)
	v := m.GlobalF32(0)
	vn := m.GlobalF32(1)
	res := m.GlobalF32(0)

	f := m.Function()
	lv := f.LocalF32()
	f.Body(
		wasm.AssignF32(lv, wasm.ConstF32(0)),
		wasm.AssignF32(v, wasm.ConstF32(0)),
		wasm.AssignF32(vn, wasm.ConstF32(1)),
		wasm.ForRangeF32{
			End: end,
			Do: func(i wasm.F32) []wasm.Instruction {
				return []wasm.Instruction{
					wasm.AssignF32(lv, wasm.AddF32(v, vn)),
					wasm.AssignF32(v, vn),
					wasm.AssignF32(vn, lv),
				}
			},
		},
		wasm.AssignF32(res, v),
	)

	m.Export("main", f)
	m.Export("end", end)
	m.Export("res", res)

	b, err := m.Compile()
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("fib.wasm", b, 0666); err != nil {
		panic(err)
	}
}
