package wasm

// f32 boolean operations
const (
	eqf32 op = 0x5B + iota
	nef32
	ltf32
	gtf32
	lef32
	gef32
)

// f32 numeric operations
const (
	absf32 op = 0x8B + iota
	negf32
	ceilf32
	floorf32
	truncf32
	nearestf32
	sqrtf32
	addf32
	subf32
	mulf32
	divf32
	minf32
	maxf32
	copysignf32
)
