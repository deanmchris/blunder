package engine

// Credit to Stockfish for the implementation details
// https://github.com/official-stockfish/Stockfish/blob/master/src/misc.h#L145
//
type PseduoRandomGenerator struct {
	state uint64
}

func (prng *PseduoRandomGenerator) Seed(seed uint64) {
	prng.state = seed
}

func (prng *PseduoRandomGenerator) Random64() uint64 {
	prng.state ^= prng.state >> 12
	prng.state ^= prng.state << 25
	prng.state ^= prng.state >> 27
	return prng.state * 2685821657736338717
}

func (prng *PseduoRandomGenerator) SparseRandom64() uint64 {
	return prng.Random64() & prng.Random64() & prng.Random64()
}
