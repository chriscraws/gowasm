package wasm

import (
	"encoding/binary"
	"io"
)

type Instruction interface {
	write(out io.Writer) error
}

type op byte

func (o op) write(out io.Writer) error {
	out.Write([]byte{byte(o)})
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

// AssignF32 assigns the value of v to dst.
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

// ConstF32 is a constant F32 value.
type ConstF32 float32

func (c ConstF32) write(out io.Writer) error {
	out.Write([]byte{0x43})
	binary.Write(out, binary.LittleEndian, c)
	return nil
}

// AbsF32 returns the absolute value of a.
func AbsF32(a F32) F32 { return ops{a, absf32} }

// NegF32 returns the reslt of negating a.
func NegF32(a F32) F32 { return ops{a, negf32} }

// CeilF32 returns a rounded up.
func CeilF32(a F32) F32 { return ops{a, ceilf32} }

// FloorF32 returns a rounded down.
func FloorF32(a F32) F32 { return ops{a, floorf32} }

// TruncF32 returns a rounded towards zero.
func TruncF32(a F32) F32 { return ops{a, truncf32} }

// NearestF32 returns the neaest integral value to a.
func NearestF32(a F32) F32 { return ops{a, nearestf32} }

// SqrtF32 returns the square root of a.
func SqrtF32(a F32) F32 { return ops{a, sqrtf32} }

// AddF32 returns the sum of a and b.
func AddF32(a, b F32) F32 { return ops{a, b, addf32} }

// SubF32 returns the difference of a and b.
func SubF32(a, b F32) F32 { return ops{a, b, subf32} }

// MulF32 returns the product of a and b.
func MulF32(a, b F32) F32 { return ops{a, b, mulf32} }

// DivF32 returns the quotient of a and b.
func DivF32(a, b F32) F32 { return ops{a, b, divf32} }

// MinF32 returns the minimum of a and b.
func MinF32(a, b F32) F32 { return ops{a, b, minf32} }

// MaxF32 returns the maximum of a and b.
func MaxF32(a, b F32) F32 { return ops{a, b, maxf32} }

// CopysignF32 returns a with the sign of b.
func CopysignF32(a, b F32) F32 { return ops{a, b, copysignf32} }
