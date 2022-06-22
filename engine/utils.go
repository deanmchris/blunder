package engine

// utils.go contains various utility functions used throughout the engine.

// Convert a string board coordinate to its position
// number.
func coordinateToPos(coordinate string) uint8 {
	file := coordinate[0] - 'a'
	rank := int(coordinate[1]-'0') - 1
	return uint8(rank*8 + int(file))
}

// Convert a position number to a string board coordinate.
func posToCoordinate(pos uint8) string {
	file := fileOf(pos)
	rank := rankOf(pos)
	return string(rune('a'+file)) + string(rune('0'+rank+1))
}

// Given a board square, return it's file.
func fileOf(sq uint8) uint8 {
	return sq % 8
}

// Given a board square, return it's rank.
func rankOf(sq uint8) uint8 {
	return sq / 8
}

// Get the absolute value of a signed 16-bit number.
func abs16(n int16) int16 {
	if n < 0 {
		return -n
	}
	return n
}

// Get the maximum between two signed 8-bit numbers.
func max8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

// Get the maximum between two signed 8-bit numbers.
func max16(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

// Determine if a square is dark.
func sqIsDark(sq uint8) bool {
	fileNo := fileOf(sq)
	rankNo := rankOf(sq)
	return ((fileNo + rankNo) % 2) == 0
}

// An implementation of a xorshift pseudo-random number
// generator for 64 bit numbers, based on the implementation
// by Stockfish:
//
// https://github.com/official-stockfish/Stockfish/blob/master/src/misc.h#L146
//
type PseduoRandomGenerator struct {
	state uint64
}

// Seed the generator.
func (prng *PseduoRandomGenerator) Seed(seed uint64) {
	prng.state = seed
}

// Generator a random 64 bit number.
func (prng *PseduoRandomGenerator) Random64() uint64 {
	prng.state ^= prng.state >> 12
	prng.state ^= prng.state << 25
	prng.state ^= prng.state >> 27
	return prng.state * 2685821657736338717
}

// Generate a random 64 bit number with few bits. This method is
// useful in finding magic numbers faster for generating slider
// attacks.
func (prng *PseduoRandomGenerator) SparseRandom64() uint64 {
	return prng.Random64() & prng.Random64() & prng.Random64()
}
