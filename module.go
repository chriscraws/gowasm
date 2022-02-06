package wasm

import (
	"bytes"
	"fmt"
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
	globals []*varF32
}

type Exportable interface {
	isExportable()
}

func (m *Module) GlobalF32(init float32) GlobalF32 {
	g := new(varF32)
	g.idx = uint32(len(m.globals))
	g.init = init
	m.globals = append(m.globals, g)
	return g
}

func (m *Module) Function() Function {
	f := new(function)
	f.idx = uint32(len(m.functions))
	m.functions = append(m.functions, f)
	return f
}

func (m *Module) Export(name string, v Exportable) {
	if m.exportNames == nil {
		m.exportNames = make(map[string]Exportable)
	}
	m.exportNames[name] = v
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

	return m.buf.Bytes(), nil
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
		err := globaltype{
			valuetype: valuetype{
				numtype: f32,
			},
			mutable: true,
		}.encode(buf)
		if err != nil {
			return err
		}
		if err := ConstF32(v.init).write(buf); err != nil {
			return err
		}
		buf.WriteByte(0x0B) // end expression
	}

	m.buf.WriteByte(0x06)
	writeu32(uint32(buf.Len()), &m.buf)
	m.buf.Write(buf.Bytes())
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
