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

var opvec4f32Tests = []struct {
	what   string
	assign wasm.Vec4F32
	expect [4]float32
}{
	{
		what: "abs",
		assign: wasm.AbsVec4F32(
			wasm.ConstVec4F32{-10, -23, 14, 0}),
		expect: [4]float32{10, 23, 14, 0},
	},
	{
		what: "neg",
		assign: wasm.NegVec4F32(
			wasm.ConstVec4F32{-10, -23, 14, 0}),
		expect: [4]float32{10, 23, -14, -0},
	},
	{
		what: "ceil",
		assign: wasm.CeilVec4F32(
			wasm.ConstVec4F32{-0.2, -2, 0.3, 1}),
		expect: [4]float32{-0, -2, 1, 1},
	},
	{
		what: "floor",
		assign: wasm.FloorVec4F32(
			wasm.ConstVec4F32{-0.2, -2, 0.3, 1}),
		expect: [4]float32{-1, -2, 0, 1},
	},
	{
		what: "trunc",
		assign: wasm.TruncVec4F32(
			wasm.ConstVec4F32{-0.2, -2, 0.3, 1}),
		expect: [4]float32{-0, -2, 0, 1},
	},
	{
		what: "nearest",
		assign: wasm.NearestVec4F32(
			wasm.ConstVec4F32{-0.8, 2.5, 0.3, 1}),
		expect: [4]float32{-1, 2, 0, 1},
	},
	{
		what: "sqrt",
		assign: wasm.SqrtVec4F32(
			wasm.ConstVec4F32{4, 16, 9, 25}),
		expect: [4]float32{2, 4, 3, 5},
	},
	{
		what: "add",
		assign: wasm.AddVec4F32(
			wasm.ConstVec4F32{4, 16, 9, 25},
			wasm.ConstVec4F32{4, 16, 9, 25}),
		expect: [4]float32{8, 32, 18, 50},
	},
	{
		what: "sub",
		assign: wasm.SubVec4F32(
			wasm.ConstVec4F32{4, 16, 9, 25},
			wasm.ConstVec4F32{30, 2, -3, 0}),
		expect: [4]float32{-26, 14, 12, 25},
	},
	{
		what: "mul",
		assign: wasm.MulVec4F32(
			wasm.ConstVec4F32{4, 16, 9, 25},
			wasm.ConstVec4F32{30, 2, -3, 0}),
		expect: [4]float32{120, 32, -27, 0},
	},
	{
		what: "div",
		assign: wasm.DivVec4F32(
			wasm.ConstVec4F32{12, 16, 9, 25},
			wasm.ConstVec4F32{3, 2, -3, 1}),
		expect: [4]float32{4, 8, -3, 25},
	},
	{
		what: "min",
		assign: wasm.MinVec4F32(
			wasm.ConstVec4F32{12, 16, 9, 25},
			wasm.ConstVec4F32{3, 2, -3, 1}),
		expect: [4]float32{3, 2, -3, 1},
	},
	{
		what: "max",
		assign: wasm.MaxVec4F32(
			wasm.ConstVec4F32{12, 16, 9, 25},
			wasm.ConstVec4F32{3, 2, -3, 1}),
		expect: [4]float32{12, 16, 9, 25},
	},
}

var forRangeTests = []struct {
	what     string
	forRange func(o wasm.MutableF32) wasm.ForRangeF32
	expect   float32
}{
	{
		what: "increment to 10",
		forRange: func(o wasm.MutableF32) wasm.ForRangeF32 {
			return wasm.ForRangeF32{
				Begin: wasm.ConstF32(0),
				End:   wasm.ConstF32(10),
				Do: func(i wasm.F32) []wasm.Instruction {
					return []wasm.Instruction{
						wasm.AssignF32(
							o,
							wasm.AddF32(o, wasm.ConstF32(1)),
						),
					}
				},
			}
		},
		expect: 10,
	},
	{
		what: "decrement to -10",
		forRange: func(o wasm.MutableF32) wasm.ForRangeF32 {
			return wasm.ForRangeF32{
				Begin: wasm.ConstF32(10),
				End:   wasm.ConstF32(0),
				Inc:   wasm.ConstF32(-1),
				Do: func(i wasm.F32) []wasm.Instruction {
					return []wasm.Instruction{
						wasm.AssignF32(
							o,
							wasm.AddF32(o, wasm.ConstF32(1)),
						),
					}
				},
			}
		},
		expect: 10,
	},
	{
		what: "do nothing for invalid range",
		forRange: func(o wasm.MutableF32) wasm.ForRangeF32 {
			return wasm.ForRangeF32{
				Begin: wasm.ConstF32(10),
				End:   wasm.ConstF32(11),
				Inc:   wasm.ConstF32(-1),
				Do: func(i wasm.F32) []wasm.Instruction {
					return []wasm.Instruction{
						wasm.AssignF32(
							o,
							wasm.AddF32(o, wasm.ConstF32(1)),
						),
					}
				},
			}
		},
		expect: 0,
	},
	{
		what: "do nothing for invalid range up",
		forRange: func(o wasm.MutableF32) wasm.ForRangeF32 {
			return wasm.ForRangeF32{
				Begin: wasm.ConstF32(10),
				End:   wasm.ConstF32(9),
				Inc:   wasm.ConstF32(1),
				Do: func(i wasm.F32) []wasm.Instruction {
					return []wasm.Instruction{
						wasm.AssignF32(
							o,
							wasm.AddF32(o, wasm.ConstF32(1)),
						),
					}
				},
			}
		},
		expect: 0,
	},
	{
		what: "nonstandard increment",
		forRange: func(o wasm.MutableF32) wasm.ForRangeF32 {
			return wasm.ForRangeF32{
				Begin: wasm.ConstF32(0),
				End:   wasm.ConstF32(15),
				Inc:   wasm.ConstF32(5),
				Do: func(i wasm.F32) []wasm.Instruction {
					return []wasm.Instruction{
						wasm.AssignF32(
							o,
							wasm.AddF32(o, i),
						),
					}
				},
			}
		},
		expect: 15,
	},
}

var ifElseTests = []struct {
	what   string
	ifElse func(o wasm.MutableF32) wasm.IfF32
	expect float32
}{
	{
		what: "truthy condition",
		ifElse: func(o wasm.MutableF32) wasm.IfF32 {
			return wasm.IfF32{
				Condition: wasm.ConstF32(1),
				Then: []wasm.Instruction{
					wasm.AssignF32(o, wasm.ConstF32(1)),
				},
			}
		},
		expect: 1,
	},
	{
		what: "falsey condition",
		ifElse: func(o wasm.MutableF32) wasm.IfF32 {
			return wasm.IfF32{
				Condition: wasm.ConstF32(0),
				Then: []wasm.Instruction{
					wasm.AssignF32(o, wasm.ConstF32(1)),
				},
				Else: []wasm.Instruction{
					wasm.AssignF32(o, wasm.ConstF32(-1)),
				},
			}
		},
		expect: -1,
	},
	{
		what: "only falsey condition",
		ifElse: func(o wasm.MutableF32) wasm.IfF32 {
			return wasm.IfF32{
				Condition: wasm.ConstF32(0),
				Else: []wasm.Instruction{
					wasm.AssignF32(o, wasm.ConstF32(-1)),
				},
			}
		},
		expect: -1,
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
		what: "a global vec4f32",
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			vec := m.GlobalVec4F32([4]float32{12, -14, 2, 1000})
			// can't export vec4f32 from wasmer :(
			dst := [4]wasm.MutableF32{}
			body := make([]wasm.Instruction, 4)
			for i := range dst {
				g := m.GlobalF32(0)
				m.Export(fmt.Sprintf("f%d", i), g)
				body[i] = wasm.AssignF32(
					g,
					wasm.ExtractLaneVec4F32(vec, i),
				)
				dst[i] = g
			}
			f := m.Function()
			f.Body(body...)
			m.Export("main", f)
			return m
		},
		test: func(ctx testContext) {
			t := ctx.t
			f, err := ctx.inst.Exports.GetFunction("main")
			if err != nil {
				t.Error(err)
			}
			if _, err := f(); err != nil {
				t.Error(err)
			}
			exp := [4]float32{12, -14, 2, 1000}
			for i := 0; i < 4; i++ {
				g, err := ctx.inst.Exports.GetGlobal(fmt.Sprintf("f%d", i))
				if err != nil {
					t.Error(err)
				}
				v, err := g.Get()
				if err != nil {
					t.Error(err)
				}
				vf, ok := v.(float32)
				if !ok {
					t.Error("exported member was not float")
				}
				if vf != exp[i] {
					t.Errorf("[%d] expected %f got %f", i, exp[i], vf)
				}
			}
		},
	},
	{
		what: "imported f32",
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			imp := m.ImportF32("root", "x")
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
	{
		what: "vecf32 ops",
		build: func(b buildContext) *wasm.Module {
			m := new(wasm.Module)
			dst := [4]wasm.MutableF32{}
			for i := range dst {
				g := m.GlobalF32(0)
				m.Export(fmt.Sprintf("o%d", i), g)
				dst[i] = g
			}
			assign := func(v wasm.Vec4F32) []wasm.Instruction {
				// can't export vec4f32 from wasmer :(
				body := make([]wasm.Instruction, 4)
				for i := range dst {
					body[i] = wasm.AssignF32(
						dst[i],
						wasm.ExtractLaneVec4F32(v, i),
					)
				}
				return body
			}
			for i, tc := range opvec4f32Tests {
				f := m.Function()
				f.Body(assign(tc.assign)...)
				m.Export(fmt.Sprintf("f%d", i), f)
			}
			return m
		},
		test: func(ctx testContext) {
			t := ctx.t
			for i, tc := range opvec4f32Tests {
				fail := func(err error) {
					t.Errorf("%s: %s", tc.what, err)
				}
				f, err := ctx.inst.Exports.GetFunction(fmt.Sprintf("f%d", i))
				if err != nil {
					fail(err)
				}
				if _, err := f(); err != nil {
					fail(err)
				}
				for i := 0; i < 4; i++ {
					g, err := ctx.inst.Exports.GetGlobal(fmt.Sprintf("o%d", i))
					if err != nil {
						fail(err)
					}
					v, err := g.Get()
					if err != nil {
						fail(err)
					}
					vf, ok := v.(float32)
					if !ok {
						fail(fmt.Errorf("exported member was not float"))
					}
					if vf != tc.expect[i] {
						fail(fmt.Errorf("[%d] expected %f got %f", i, tc.expect[i], vf))
					}
				}
			}
		},
	},
	{
		what: "for-range tests",
		build: func(ctx buildContext) *wasm.Module {
			m := new(wasm.Module)
			o := m.GlobalF32(0)
			m.Export("o", o)
			reset := m.Function()
			reset.Body(wasm.AssignF32(o, wasm.ConstF32(0)))
			m.Export("reset", reset)
			for i, tc := range forRangeTests {
				f := m.Function()
				f.Body(tc.forRange(o))
				m.Export(fmt.Sprintf("f%d", i), f)
			}
			return m
		},
		test: func(ctx testContext) {
			reset, _ := ctx.inst.Exports.GetFunction("reset")
			res, _ := ctx.inst.Exports.GetGlobal("o")
			for i, tc := range forRangeTests {
				ctx.t.Run(tc.what, func(t *testing.T) {
					tc := tc
					reset()
					f, _ := ctx.inst.Exports.GetFunction(fmt.Sprintf("f%d", i))
					_, err := f()
					if err != nil {
						t.Error(err)
					}
					v, _ := res.Get()
					vf := v.(float32)
					if vf != tc.expect {
						t.Errorf("expected %f, got %f", tc.expect, vf)
					}
				})
			}
		},
	},
	{
		what: "if-else tests",
		build: func(ctx buildContext) *wasm.Module {
			m := new(wasm.Module)
			o := m.GlobalF32(0)
			m.Export("o", o)
			reset := m.Function()
			reset.Body(wasm.AssignF32(o, wasm.ConstF32(0)))
			m.Export("reset", reset)
			for i, tc := range ifElseTests {
				f := m.Function()
				f.Body(tc.ifElse(o))
				m.Export(fmt.Sprintf("f%d", i), f)
			}
			return m
		},
		test: func(ctx testContext) {
			reset, _ := ctx.inst.Exports.GetFunction("reset")
			res, _ := ctx.inst.Exports.GetGlobal("o")
			for i, tc := range ifElseTests {
				ctx.t.Run(tc.what, func(t *testing.T) {
					tc := tc
					reset()
					f, _ := ctx.inst.Exports.GetFunction(fmt.Sprintf("f%d", i))
					_, err := f()
					if err != nil {
						t.Error(err)
					}
					v, _ := res.Get()
					vf := v.(float32)
					if vf != tc.expect {
						t.Errorf("expected %f, got %f", tc.expect, vf)
					}
				})
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
