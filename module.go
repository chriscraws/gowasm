package wasm

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

// Module is a repesentation of a WASM module.
// The Compile method will convert the wasm into a binary representation.
type Module struct {
	buf bytes.Buffer

	exportNames map[string]Exportable

	// function state
	functions       []*function
	functionTypeMap map[*function]int
	functionTypes   []functype

	// exports state
	exports [][]byte

	// globals
	globals []global

	// imports
	imports map[[2]string]importable

	doesUseMemory bool
}

// An Exportable type can be exported from the module.
// Exportables must be exported with the Module.Export
// method, and can be created with the following functions:
//
// Module.GlobalF32
// Module.Function
type Exportable interface {
	isExportable()
}

// GlobalF32 creates a global, mutable F32 object.
func (m *Module) GlobalF32(init float32) GlobalF32 {
	g := new(varF32)
	g.idx = uint32(len(m.globals))
	g.init = init
	m.globals = append(m.globals, g)
	return g
}

// GlobalVec4F32 creates a global, mutable Vec4F32 object.
func (m *Module) GlobalVec4F32(init [4]float32) GlobalVec4F32 {
	g := new(vec4F32)
	g.idx = uint32(len(m.globals))
	g.init = init
	m.globals = append(m.globals, g)
	return g
}

// Function instantiates a function.
func (m *Module) Function() Function {
	f := new(function)
	f.idx = uint32(len(m.functions))
	m.functions = append(m.functions, f)
	return f
}

// Export exports v as name. If a previous Exportable has already
// been exported as name, it will be replaced.
func (m *Module) Export(name string, v Exportable) {
	if m.exportNames == nil {
		m.exportNames = make(map[string]Exportable)
	}
	m.exportNames[name] = v
}

func (m *Module) addImport(mod, name string, v importable) {
	if m.imports == nil {
		m.imports = make(map[[2]string]importable)
	}
	key := [2]string{mod, name}
	if _, ok := m.imports[key]; ok {
		panic(fmt.Errorf("duplicate import %q", key))
	}
	for _, g := range m.globals {
		g.incGlobalIndex()
	}
	if g, ok := v.(global); ok {
		g.setGlobalIndex(uint32(len(m.imports)))
	}
	m.imports[key] = v
}

// ImportF32 imports a global F32 value.
// If symbol has already been imported, ImportF32 panics.
// symbol must be in the format "module.name"
func (m *Module) ImportF32(mod, name string) MutableF32 {
	out := new(varF32)
	m.addImport(mod, name, out)
	return out
}

// ImportSliceF32 imports a slice of float32 values located
// in memory. This requires memory to be provided to the wasm
// module.
// SliceF32 is an i64 value that is interpreted as two u32 offsets,
// defining the inclusive lower bound and exclusive upper bound
// (exlusive) of the slice memory offsets.
func (m *Module) ImportSliceF32(name string) SliceF32 {
	out := new(sliceF32)
	m.addImport("_sf32", name, out)
	m.doesUseMemory = true
	return out
}

// Compile compiles the module into binary WASM format.
func (m *Module) Compile() ([]byte, error) {

	m.functionTypeMap = make(map[*function]int)

	// collect exports
	if err := m.collectExports(); err != nil {
		return nil, fmt.Errorf("failed to collect Exports: %s", err)
	}

	// encode
	m.buf.Reset()

	// write magic number
	m.buf.Write([]byte{0x00, 0x61, 0x73, 0x6D})
	// write version
	m.buf.Write([]byte{0x01, 0, 0, 0})

	// (1) type section
	if err := m.writeTypeSection(); err != nil {
		return nil, fmt.Errorf("failed to write type section: %s", err)
	}

	// (2) import section
	if err := m.writeImportSection(); err != nil {
		return nil, fmt.Errorf("failed to write import section: %s", err)
	}

	// (3) function section
	if err := m.writeFunctionSection(); err != nil {
		return nil, fmt.Errorf("failed to write function section: %s", err)
	}

	// (6) global section
	if err := m.writeGlobalSection(); err != nil {
		return nil, fmt.Errorf("failed to write global section: %s", err)
	}

	// (7) export section
	if err := m.writeExportSection(); err != nil {
		return nil, fmt.Errorf("failed to write export section: %s", err)
	}

	// (10) code section
	if err := m.writeCodeSection(); err != nil {
		return nil, fmt.Errorf("failed to write code section: %s", err)
	}

	out := m.buf.Bytes()
	m.buf = bytes.Buffer{}
	return out, nil
}

func (m *Module) addFunction(f *function) {
	if _, ok := m.functionTypeMap[f]; ok {
		return
	}
	t := f.functype()
	for i, ft := range m.functionTypes {
		if ft.equals(t) {
			m.functionTypeMap[f] = i
			break
		}
	}
	m.functionTypeMap[f] = len(m.functionTypes)
	m.functionTypes = append(m.functionTypes, t)
}

func (m *Module) collectExports() error {
	exportNames := make(sort.StringSlice, len(m.exportNames))
	m.exports = make([][]byte, len(m.exportNames))
	{
		var i int
		for k := range m.exportNames {
			exportNames[i] = k
			i++
		}
		exportNames.Sort()
	}
	for i, name := range exportNames {
		e := m.exportNames[name]
		var ei uint32
		var eid byte
		switch v := e.(type) {
		case *function:
			ei = v.idx
			eid = 0x0
			m.addFunction(v)
		case *varF32:
			ei = v.idx
			eid = 0x03
		default:
			return fmt.Errorf("%v is unsupported export type", v)
		}
		buf := new(bytes.Buffer)
		writeu32(uint32(len(name)), buf)
		buf.WriteString(name)
		buf.WriteByte(eid)
		writeu32(uint32(ei), buf)
		m.exports[i] = buf.Bytes()
	}
	return nil
}

func (m *Module) writeTypeSection() error {
	if len(m.functionTypes) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	writeu32(uint32(len(m.functionTypes)), buf)
	for _, ft := range m.functionTypes {
		if err := ft.encode(buf); err != nil {
			return fmt.Errorf("failed to write functype: %s", err)
		}
	}

	m.buf.WriteByte(0x01)
	writeu32(uint32(buf.Len()), &m.buf)
	m.buf.Write(buf.Bytes())
	return nil
}

func (m *Module) writeImportSection() error {
	if len(m.imports) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	// vec(import)
	writeu32(uint32(len(m.imports)), buf)
	imports := make([][2]string, len(m.imports))
	for k, m := range m.imports {
		imports[m.(global).globalIndex()] = k
	}
	if m.doesUseMemory {
		key := [2]string{"wasm", "memory"}
		m.imports[key] = memImport{}
		imports = append(imports, key)
	}
	for _, imp := range imports {
		// module
		writeu32(uint32(len(imp[0])), buf)
		buf.WriteString(imp[0])
		// name
		writeu32(uint32(len(imp[1])), buf)
		buf.WriteString(imp[1])
		g := m.imports[imp]
		if err := g.writeImportDesc(buf); err != nil {
			return fmt.Errorf("failed to write import %s.%s: %s", imp[0], imp[1], err)
		}
	}
	m.buf.WriteByte(2)
	writeu32(uint32(buf.Len()), &m.buf)
	m.buf.Write(buf.Bytes())
	return nil
}

func (m *Module) writeFunctionSection() error {
	if len(m.functions) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	writeu32(uint32(len(m.functions)), buf)
	for _, f := range m.functions {
		if i, ok := m.functionTypeMap[f]; ok {
			writeu32(uint32(i), buf)
		} else {
			return fmt.Errorf("failed to write function: %s", f)
		}
	}

	m.buf.WriteByte(0x03)
	writeu32(uint32(buf.Len()), &m.buf)
	m.buf.Write(buf.Bytes())
	return nil
}

func (m *Module) writeGlobalSection() error {
	if len(m.globals) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	writeu32(uint32(len(m.globals)), buf)
	for _, v := range m.globals {
		var err error
		switch v := v.(type) {
		case *varF32:
			err = m.writeF32Global(v, buf)
		case *vec4F32:
			err = m.writeVec4F32Global(v, buf)
		default:
			err = fmt.Errorf("%v is not a global-compatible type", v)
		}
		if err != nil {
			return err
		}
	}

	m.buf.WriteByte(0x06)
	writeu32(uint32(buf.Len()), &m.buf)
	m.buf.Write(buf.Bytes())
	return nil
}

func (m *Module) writeF32Global(v *varF32, out io.Writer) error {
	err := globaltype{
		valuetype: valuetype{
			numtype: f32,
		},
		mutable: true,
	}.encode(out)
	if err != nil {
		return err
	}
	if err := ConstF32(v.init).write(instCtx{Writer: out}); err != nil {
		return err
	}
	out.Write([]byte{0x0B}) // end expression
	return nil
}

func (m *Module) writeVec4F32Global(v *vec4F32, out io.Writer) error {
	err := globaltype{
		valuetype: valuetype{
			vectype: true,
		},
	}.encode(out)
	if err != nil {
		return err
	}
	if err := ConstVec4F32(v.init).write(instCtx{Writer: out}); err != nil {
		return err
	}
	out.Write([]byte{0x0B}) // end expression
	return nil
}

func (m *Module) writeExportSection() error {
	if len(m.exports) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	writeu32(uint32(len(m.exports)), buf)
	for _, b := range m.exports {
		buf.Write(b)
	}

	m.buf.WriteByte(0x07)
	writeu32(uint32(buf.Len()), &m.buf)
	m.buf.Write(buf.Bytes())
	return nil
}

func (m *Module) writeCodeSection() error {
	if len(m.functions) == 0 {
		return nil
	}
	buf := new(bytes.Buffer)
	writeu32(uint32(len(m.functions)), buf)
	for _, f := range m.functions {
		if err := f.encode(buf); err != nil {
			return err
		}
	}

	m.buf.WriteByte(10)
	writeu32(uint32(buf.Len()), &m.buf)
	m.buf.Write(buf.Bytes())
	return nil
}
