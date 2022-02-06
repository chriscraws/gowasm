package wasm

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Instruction interface {
	write(out io.Writer) error
}

type op struct {
	code byte
	args []byte
}

func (o *op) write(out io.Writer) error {
	out.Write(append([]byte{o.code}, o.args...))
	return nil
}

type ops []Instruction

func (o ops) write(out io.Writer) error {
	for _, v := range o {
		if err := v.write(out); err != nil {
			return err
		}
	}
	return nil
}

func AssignF32(dst MutableF32, v F32) Instruction {
	return assignF32{dst: dst, v: v}
}

type assignF32 struct {
	dst MutableF32
	v   F32
}

func (a assignF32) write(out io.Writer) error {
	if err := a.v.write(out); err != nil {
		return err
	}
	if err := a.dst.set(out); err != nil {
		return err
	}
	return nil
}

var opAddF32 = &op{code: 0x92}

func AddF32(a, b F32) F32 {
	return ops{a, b, opAddF32}
}

func ConstF32(v float32) F32 {
	out := new(bytes.Buffer)
	binary.Write(out, binary.LittleEndian, v)
	return &op{
		code: 0x43,
		args: out.Bytes(),
	}
}
