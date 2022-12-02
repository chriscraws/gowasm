package wasm

import (
	"encoding/binary"
)

// Vec4F32 represents a float32 vector of 4
// elements.
type Vec4F32 interface {
	Instruction
	isVec4F32()
}

type MutableVec4F32 interface {
	Vec4F32
}

type vec4F32 struct {
	idx    uint32
	init   [4]float32
	offset uint32
	align  uint32
}

func (v *vec4F32) write(out instCtx) error {
	out.Write([]byte{0x23})
	writeu32(v.idx, out)
	return nil
}

func (v *vec4F32) isVec4F32() {}

func (v *vec4F32) incGlobalIndex() {
	v.idx++
}

func (v *vec4F32) setGlobalIndex(i uint32) {
	v.idx = i
}

func (v *vec4F32) globalIndex() uint32 {
	return v.idx
}

type ConstVec4F32 [4]float32

func (c ConstVec4F32) isVec4F32() {}

func (c ConstVec4F32) write(out instCtx) error {
	constV128.write(out)
	binary.Write(out, binary.LittleEndian, c[:])
	return nil
}

// GlobalVec4F32 represents a mutable Vec4F32 in the
// global scope of the wasm module. It may be assigned
// or read from any function, and also can be exported
// to be observed by the runtime.
type GlobalVec4F32 interface {
	Vec4F32
}

type extractLaneVec4F32 struct {
	x Vec4F32
	i int
}

func (e extractLaneVec4F32) isF32() {}

func (e extractLaneVec4F32) write(out instCtx) error {
	if err := e.x.write(out); err != nil {
		return err
	}
	extractLanef32x4V128.write(out)
	writeu32(uint32(e.i), out)
	return nil
}

// Return the index i as an F32.
func ExtractLaneVec4F32(x Vec4F32, i int) F32 {
	return extractLaneVec4F32{x: x, i: i}
}

// AbsVec4F32 returns the absolute value of a.
func AbsVec4F32(a Vec4F32) Vec4F32 { return opsVec4F32{a, absf32x4V128} }

// NegVec4F32 returns the reslt of negating a.
func NegVec4F32(a Vec4F32) Vec4F32 { return opsVec4F32{a, negf32x4V128} }

// CeilVec4F32 returns a rounded up.
func CeilVec4F32(a Vec4F32) Vec4F32 { return opsVec4F32{a, ceilf32x4V128} }

// FloorVec4F32 returns a rounded down.
func FloorVec4F32(a Vec4F32) Vec4F32 { return opsVec4F32{a, floorf32x4V128} }

// TruncVec4F32 returns a rounded towards zero.
func TruncVec4F32(a Vec4F32) Vec4F32 { return opsVec4F32{a, truncf32x4V128} }

// NearestVec4F32 returns the neaest integral value to a.
func NearestVec4F32(a Vec4F32) Vec4F32 { return opsVec4F32{a, nearestf32x4V128} }

// SqrtVec4F32 returns the square root of a.
func SqrtVec4F32(a Vec4F32) Vec4F32 { return opsVec4F32{a, sqrtf32x4V128} }

// AddVec4F32 returns the sum of a and b.
func AddVec4F32(a, b Vec4F32) Vec4F32 { return opsVec4F32{a, b, addf32x4V128} }

// SubVec4F32 returns the difference of a and b.
func SubVec4F32(a, b Vec4F32) Vec4F32 { return opsVec4F32{a, b, subf32x4V128} }

// MulVec4F32 returns the product of a and b.
func MulVec4F32(a, b Vec4F32) Vec4F32 { return opsVec4F32{a, b, mulf32x4V128} }

// DivVec4F32 returns the quotient of a and b.
func DivVec4F32(a, b Vec4F32) Vec4F32 { return opsVec4F32{a, b, divf32x4V128} }

// MinVec4F32 returns the minimum of a and b.
func MinVec4F32(a, b Vec4F32) Vec4F32 { return opsVec4F32{a, b, minf32x4V128} }

// MaxVec4F32 returns the maximum of a and b.
func MaxVec4F32(a, b Vec4F32) Vec4F32 { return opsVec4F32{a, b, maxf32x4V128} }
