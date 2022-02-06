package wasm_test

import (
	"testing"

	wasm "github.com/chriscraws/gowasm"
	"github.com/wasmerio/wasmer-go/wasmer"
)

var tcs = []struct {
	what  string
	build func() *wasm.Module
	test  func(t *testing.T, inst *wasmer.Instance)
}{
	{
		what: "an empty module",
		build: func() *wasm.Module {
			return new(wasm.Module)
		},
	},
	{
		what: "an exported f32",
		build: func() *wasm.Module {
			m := new(wasm.Module)
			m.Export("hello", m.GlobalF32(38.89))
			return m
		},
		test: func(t *testing.T, inst *wasmer.Instance) {
			glob, err := inst.Exports.GetGlobal("hello")
			if err != nil {
				t.Error(err)
			}
			if glob.Type().ValueType().Kind() != wasmer.F32 {
				t.Errorf("invalid type for global hello")
			}
			vf, err := glob.Get()
			if err != nil || vf.(float32) != 38.89 {
				t.Errorf("invalid global init value")
			}
		},
	},
	{
		what: "an exported function",
		build: func() *wasm.Module {
			m := new(wasm.Module)
			hello := m.GlobalF32(38.89)
			m.Export("hello", hello)
			fn := m.Function()
			fn.AddInstructions(
				wasm.AssignF32(hello, wasm.ConstF32(10)))
			m.Export("set_ten", fn)
			return m
		},
		test: func(t *testing.T, inst *wasmer.Instance) {
			glob, _ := inst.Exports.GetGlobal("hello")
			fn, err := inst.Exports.GetFunction("set_ten")
			if err != nil {
				t.Error(err)
			}
			vf, err := glob.Get()
			if err != nil || vf.(float32) != 38.89 {
				t.Errorf("invalid global init value")
			}
			_, err = fn()
			vf, _ = glob.Get()
			if err != nil || vf.(float32) != 10 {
				t.Errorf("function did not set to ten")
			}
		},
	},
	{
		what: "a function with locals",
		build: func() *wasm.Module {
			m := new(wasm.Module)
			hello := m.GlobalF32(38.89)
			m.Export("hello", hello)
			fn := m.Function()
			loc := fn.LocalF32()
			fn.AddInstructions(
				wasm.AssignF32(loc, wasm.ConstF32(15)),
				wasm.AssignF32(hello, loc))
			m.Export("set_fifteen", fn)
			return m
		},
		test: func(t *testing.T, inst *wasmer.Instance) {
			glob, _ := inst.Exports.GetGlobal("hello")
			fn, err := inst.Exports.GetFunction("set_fifteen")
			if err != nil {
				t.Error(err)
			}
			vf, err := glob.Get()
			if err != nil || vf.(float32) != 38.89 {
				t.Errorf("invalid global init value")
			}
			_, err = fn()
			vf, _ = glob.Get()
			if err != nil || vf.(float32) != 15 {
				t.Errorf("function did not set to fifteen")
			}
		},
	},
}

func TestWasm(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.what, func(t *testing.T) {
			mod := tc.build()
			buf, err := mod.Compile()
			if err != nil {
				t.Error(err)
			}
			engine := wasmer.NewEngine()
			store := wasmer.NewStore(engine)
			err = wasmer.ValidateModule(store, buf)
			if err != nil {
				t.Error(err)
			}
			module, err := wasmer.NewModule(store, buf)
			if err != nil {
				t.Error(err)
			}
			if module == nil {
				return
			}
			importObject := wasmer.NewImportObject()
			inst, err := wasmer.NewInstance(module, importObject)
			if err != nil {
				t.Error(err)
			}
			if tc.test != nil {
				tc.test(t, inst)
			}
		})
	}
}
