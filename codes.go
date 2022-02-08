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

// vector instructions
const (
	// 0
	loadV128 vecOp = iota
	load8x8sV128
	load8x8uV128
	load16x4sV128
	load16x4uV128
	load32x2sV128
	load32x2uV128
	load8splatV128
	load16splatV128
	load32splatV128
	// 10
	load64splatV128
	storeV128
	constV128
	shuffleV128
	swizzlei8x16V128
	splati8x16V128
	splati16x8V128
	splati32x4V128
	splati64x2V128
	splatf32x4V128
	// 20
	splatf64x2V128
)

const (
	extractLanef32x4V128 vecOp = 31 + iota
)

const (
	load8laneV128 vecOp = 84 + iota
	load16laneV128
	load32laneV128
	load64laneV128
	store8laneV128
	store16laneV128
	store32laneV128
	store64laneV128
	load32zeroV128
	load64zeroV128
)

// vec f32 numeric operations
const (
	ceilf32x4V128 vecOp = 103 + iota
	floorf32x4V128
	truncf32x4V128
	nearestf32x4V128
)

const (
	absf32x4V128 vecOp = 224 + iota
	negf32x4V128
)

const (
	sqrtf32x4V128 vecOp = 227 + iota
	addf32x4V128
	subf32x4V128
	mulf32x4V128
	divf32x4V128
	minf32x4V128
	maxf32x4V128
	pminf32x4V128
	pmaxf32x4V128
)
