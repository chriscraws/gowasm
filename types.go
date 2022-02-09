// package wasm generates binary wasm from Go types
// it has a focus on mathematical expressions.
package wasm

import (
	"fmt"
	"io"
	"strings"
)

type global interface {
	incGlobalIndex()
	setGlobalIndex(i uint32)
	globalIndex() uint32
}

type importable interface {
	writeImportDesc(out io.Writer) error
}

type symbol interface {
	String() string
	write(out io.Writer) error
}

type numtype byte

const (
	i32 numtype = 0x7F
	i64         = 0x7E
	f32         = 0x7D
	f64         = 0x7C
)

func (nt numtype) String() string {
	switch nt {
	case i32:
		return "i32"
	case i64:
		return "i64"
	case f32:
		return "f32"
	case f64:
		return "f64"
	default:
		panic(fmt.Errorf("%d is not a numtype", nt))
	}
}

func (nt numtype) encode(out io.Writer) error {
	if nt < f64 || nt > i32 {
		return fmt.Errorf("%d is not a numtype", nt)
	}
	out.Write([]byte{byte(nt)})
	return nil
}

type reftype byte

const (
	funcref   reftype = 0x70
	externref         = 0x6F
)

func (rt reftype) encode(out io.Writer) error {
	if rt != funcref && rt != externref {
		return fmt.Errorf("%d is not a reftype", rt)
	}
	out.Write([]byte{byte(rt)})
	return nil
}

func (rt reftype) String() string {
	switch rt {
	case funcref:
		return "funcref"
	case externref:
		return "externref"
	default:
		panic(fmt.Errorf("%d is not a reftype", rt))
	}
}

// if numType is not set, valueType is a refType
type valuetype struct {
	numtype numtype
	reftype reftype
	vectype bool
}

func (vt valuetype) encode(out io.Writer) error {
	if vt.vectype {
		out.Write([]byte{0x7B})
		return nil
	}
	if vt.numtype == 0 && vt.reftype == 0 {
		return fmt.Errorf("%v is invalid valuetype", vt)
	}
	if vt.numtype == 0 {
		return vt.reftype.encode(out)
	}
	return vt.numtype.encode(out)
}

func (vt valuetype) String() string {
	if vt.numtype == 0 && vt.reftype == 0 {
		panic("invalid valuetype")
	}
	if vt.numtype == 0 {
		return vt.reftype.String()
	}
	return vt.numtype.String()
}

type resulttype []valuetype

func (rt resulttype) String() string {
	list := make([]string, len(rt))
	for i, t := range rt {
		list[i] = t.String()
	}
	return "(" + strings.Join(list, ", ") + ")"
}

func (rt resulttype) encode(out io.Writer) error {
	writeu32(uint32(len(rt)), out)
	for _, t := range rt {
		if err := t.encode(out); err != nil {
			return fmt.Errorf("failed to encode resulttype param: %s", err)
		}
	}
	return nil
}

type functype struct {
	params  resulttype
	results resulttype
}

func (ft functype) equals(other functype) bool {
	if len(ft.params) != len(other.params) ||
		len(ft.results) != len(other.results) {
		return false
	}

	for i, p := range ft.params {
		if p != other.params[i] {
			return false
		}
	}

	for i, r := range ft.results {
		if r != other.results[i] {
			return false
		}
	}

	return true
}

func (ft functype) String() string {
	return "func " + ft.params.String() + " -> " + ft.results.String()
}

func (ft functype) encode(out io.Writer) error {
	out.Write([]byte{0x60})
	if err := ft.params.encode(out); err != nil {
		return fmt.Errorf("failed to encode function params: %s", err)
	}
	if err := ft.results.encode(out); err != nil {
		return fmt.Errorf("failed to encode function results: %s", err)
	}
	return nil
}

type globaltype struct {
	mutable   bool
	valuetype valuetype
}

func (gt globaltype) String() string {
	vstr := gt.valuetype.String()
	if gt.mutable {
		return "var " + vstr
	} else {
		return "const " + vstr
	}
}

func (gt globaltype) encode(out io.Writer) error {
	if err := gt.valuetype.encode(out); err != nil {
		return err
	}
	o := []byte{0x0}
	if gt.mutable {
		o[0] = 0x01
	}
	out.Write(o)
	return nil
}
