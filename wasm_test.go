package wasm_test

import (
	"fmt"
	"testing"

	wasm "github.com/chriscraws/gowasm"
	"github.com/wasmerio/wasmer-go/wasmer"
)

var opf32Tests = []struct {
	what   string
	assign wasm.F32
	expect float32
}{
	{
		what:   "abs",
		assign: wasm.AbsF32(wasm.ConstF32(-10)),
		expect: 10,
	},
	{
		what:   "neg",
		assign: wasm.NegF32(wasm.ConstF32(10)),
		expect: -10,
	},
	{
		what:   "ceil",
		assign: wasm.CeilF32(wasm.ConstF32(-0.2)),
		expect: -0,
	},
	{
		what:   "floor",
		assign: wasm.FloorF32(wasm.ConstF32(-0.2)),
		expect: -1,
	},
	{
		what:   "trunc negative",
		assign: wasm.TruncF32(wasm.ConstF32(-0.2)),
		expect: -0,
	},
	{
		what:   "trunc positive",
		assign: wasm.TruncF32(wasm.ConstF32(0.2)),
		expect: 0,
	},
	{
		what:   "nearest 1",
		assign: wasm.NearestF32(wasm.ConstF32(0.2)),
		expect: 0,
	},
	{
		what:   "nearest 2",
		assign: wasm.NearestF32(wasm.ConstF32(0.6)),
		expect: 1,
	},
	{
		what:   "nearest 3",
		assign: wasm.NearestF32(wasm.ConstF32(-23.2)),
		expect: -23,
	},
	{
		what:   "sqrt",
		assign: wasm.SqrtF32(wasm.ConstF32(4)),
		expect: 2,
	},
	{
		what:   "add",
		assign: wasm.AddF32(wasm.ConstF32(1), wasm.ConstF32(5)),
		expect: 6,
	},
	{
		what:   "sub",
		assign: wasm.SubF32(wasm.ConstF32(1), wasm.ConstF32(5)),
		expect: -4,
	},
	{
		what:   "mul",
		assign: wasm.MulF32(wasm.ConstF32(3), wasm.ConstF32(5)),
		expect: 15,
	},
	{
		what:   "div",
		assign: wasm.DivF32(wasm.ConstF32(30), wasm.ConstF32(5)),
		expect: 6,
	},
	{
		what:   "min",
		assign: wasm.MinF32(wasm.ConstF32(30), wasm.ConstF32(5)),
		expect: 5,
	},
	{
		what:   "max",
		assign: wasm.MaxF32(wasm.ConstF32(30), wasm.ConstF32(5)),
		expect: 30,
	},
	{
		what:   "copysign",
		assign: wasm.CopysignF32(wasm.ConstF32(30), wasm.ConstF32(5)),
		expect: 30,
	},
	{
		what:   "copysign 2",
		assign: wasm.CopysignF32(wasm.ConstF32(30), wasm.ConstF32(-5)),
		expect: -30,
	},
}

type buildContext struct {
	t     *testing.T
	imp   *wasmer.ImportObject
	store *wasmer.Store
	data  *interface{}
}

type testContext struct {
	t    *testing.T
	inst *wasmer.Instance
	imp  *wasmer.ImportObject
	data *interface{}
}

var tcs = []struct {
	what  string
	build func(b buildContext) *wasm.Module
	test  func(ctx testContext)
}{
	{
		what: "an empty module",
		build: func(b buildContext) *wasm.Module {
			return new(wasm.Module)
		},
	},
	{
		what: "an exported f32",
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			m.Export("hello", m.GlobalF32(38.89))
			return m
		},
		test: func(ctx testContext) {
			t := ctx.t
			glob, err := ctx.inst.Exports.GetGlobal("hello")
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
		what: "imported f32",
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			imp := m.ImportF32("root.x")
			f := m.Function()
			f.Body(wasm.AssignF32(imp, wasm.ConstF32(123)))
			m.Export("main", f)
			x := wasmer.NewGlobal(
				b.store,
				wasmer.NewGlobalType(
					wasmer.NewValueType(wasmer.F32), wasmer.MUTABLE),
				wasmer.NewF32(float32(5)),
			)
			*b.data = x
			b.imp.Register("root", map[string]wasmer.IntoExtern{
				"x": x,
			})
			return m
		},
		test: func(ctx testContext) {
			t := ctx.t
			fn, _ := ctx.inst.Exports.GetFunction("main")
			err, _ := fn()
			if err != nil {
				t.Error(err)
			}
			v := (*ctx.data).(*wasmer.Global)
			if vf, err := v.Get(); err != nil || vf.(float32) != 123 {
				t.Errorf("expected %f, got %f", 123.0, vf.(float32))
			}
		},
	},
	{
		what: "an exported function",
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			hello := m.GlobalF32(38.89)
			m.Export("hello", hello)
			fn := m.Function()
			fn.Body(wasm.AssignF32(hello, wasm.ConstF32(10)))
			m.Export("set_ten", fn)
			return m
		},
		test: func(ctx testContext) {
			t := ctx.t
			glob, _ := ctx.inst.Exports.GetGlobal("hello")
			fn, err := ctx.inst.Exports.GetFunction("set_ten")
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
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			hello := m.GlobalF32(38.89)
			m.Export("hello", hello)
			fn := m.Function()
			loc := fn.LocalF32()
			fn.Body(
				wasm.AssignF32(loc, wasm.ConstF32(15)),
				wasm.AssignF32(hello, loc))
			m.Export("set_fifteen", fn)
			return m
		},
		test: func(ctx testContext) {
			t := ctx.t
			glob, _ := ctx.inst.Exports.GetGlobal("hello")
			fn, err := ctx.inst.Exports.GetFunction("set_fifteen")
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
	{
		what: "f32 ops",
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			out := m.GlobalF32(0)
			m.Export("out", out)
			for i, tc := range opf32Tests {
				f := m.Function()
				f.Body(wasm.AssignF32(out, tc.assign))
				m.Export(fmt.Sprintf("f%d", i), f)
			}
			return m
		},
		test: func(ctx testContext) {
			t := ctx.t
			g, _ := ctx.inst.Exports.GetGlobal("out")
			for i, tc := range opf32Tests {
				n := fmt.Sprintf("f%d", i)
				f, _ := ctx.inst.Exports.GetFunction(n)
				_, err := f()
				if err != nil {
					t.Errorf("failed to run f32op %q: %s", tc.what, err)
				}
				v, _ := g.Get()
				vf := v.(float32)
				if vf != tc.expect {
					t.Errorf("%s: expected %f got %f", tc.what, tc.expect, vf)
				}
			}
		},
	},
}

func TestWasm(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.what, func(t *testing.T) {
			var data interface{}
			engine := wasmer.NewEngine()
			store := wasmer.NewStore(engine)
			importObject := wasmer.NewImportObject()
			mod := tc.build(buildContext{
				t:     t,
				imp:   importObject,
				data:  &data,
				store: store,
			})
			buf, err := mod.Compile()
			if err != nil {
				t.Error(err)
			}
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
			inst, err := wasmer.NewInstance(module, importObject)
			if err != nil {
				t.Error(err)
			}
			if tc.test != nil {
				tc.test(testContext{
					t:    t,
					inst: inst,
					imp:  importObject,
					data: &data,
				})
			}
		})
	}
}
