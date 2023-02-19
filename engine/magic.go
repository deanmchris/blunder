package engine

import (
	"math/bits"
)

// TODO: Add comment explaning magic number generation algorithm

var RookMagics [64]Magic
var BishopMagics [64]Magic

var RookMoves [64][4096]uint64
var BishopMoves [64][512]uint64

type Magic struct {
	MagicNo     uint64
	BlockerMask uint64
	Shift       uint8
}

// Optimized seeding values to find magic numbers based on file/rank, from Stockfish:
var MagicSeeds [8]uint64 = [8]uint64{728, 10316, 55013, 32803, 12281, 15100, 16645, 255}

func genRookMovesHQ(sq uint8, occupied uint64, genMask bool) uint64 {
	sliderBB := MostSigBitBB >> sq
	fileMask := MaskFile[fileOf(sq)]
	rankMask := MaskRank[rankOf(sq)]

	rhs := bits.Reverse64(bits.Reverse64((occupied & rankMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupied & rankMask) - 2*sliderBB
	eastWestMoves := (rhs ^ lhs) & rankMask

	rhs = bits.Reverse64(bits.Reverse64((occupied & fileMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupied & fileMask) - 2*sliderBB
	northSouthMoves := (rhs ^ lhs) & fileMask

	if genMask {
		northSouthMoves &= ClearRank[Rank1] & ClearRank[Rank8]
		eastWestMoves &= ClearFile[FileA] & ClearFile[FileH]
	}

	return northSouthMoves | eastWestMoves
}

func genBishopMovesHQ(sq uint8, occupied uint64, genMask bool) uint64 {
	sliderBB := MostSigBitBB >> sq
	diagonalMask := MaskDiagonal[fileOf(sq)-rankOf(sq)+7]
	antidiagonalMask := MaskAntidiagonal[14-(rankOf(sq)+fileOf(sq))]

	rhs := bits.Reverse64(bits.Reverse64((occupied & diagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupied & diagonalMask) - 2*sliderBB
	diagonalMoves := (rhs ^ lhs) & diagonalMask

	rhs = bits.Reverse64(bits.Reverse64((occupied & antidiagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupied & antidiagonalMask) - 2*sliderBB
	antidiagonalMoves := (rhs ^ lhs) & antidiagonalMask

	edges := FullBB
	if genMask {
		edges = ClearFile[FileA] & ClearFile[FileH] & ClearRank[Rank1] & ClearRank[Rank8]
	}
	return (diagonalMoves | antidiagonalMoves) & edges
}

func genSubsets(bb uint64, numSubsets uint16) (subsets []uint64) {
	subsets = make([]uint64, numSubsets)
	currSubset := EmptyBB
	index := 0

	for ok := true; ok; ok = currSubset != 0 {
		subsets[index] = currSubset
		currSubset = (currSubset - bb) & bb
		index++
	}

	return subsets
}

func genRookMagics() {
	prng := PseduoRandomGenerator{}

	for sq := uint8(0); sq < 64; sq++ {
		magic := &RookMagics[sq]

		magic.BlockerMask = genRookMovesHQ(sq, EmptyBB, true)
		numBitsSetHigh := bits.OnesCount64(magic.BlockerMask)
		magic.Shift = uint8(64 - numBitsSetHigh)

		subsets := genSubsets(magic.BlockerMask, 1<<numBitsSetHigh)
		subsetMoves := make([]uint64, 1<<numBitsSetHigh)

		for i, subset := range subsets {
			subsetMoves[i] = genRookMovesHQ(sq, subset, false)
		}

		searching := true
		prng.Seed(MagicSeeds[rankOf(sq)])

		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			RookMoves[sq] = [4096]uint64{}

			for i, subset := range subsets {
				index := (subset * possible_magic) >> magic.Shift
				moves := subsetMoves[i]

				if RookMoves[sq][index] != EmptyBB && RookMoves[sq][index] != moves {
					searching = true
					break
				}

				RookMoves[sq][index] = moves
			}
		}
	}
}

func genBishopMagics() {
	prng := PseduoRandomGenerator{}

	for sq := uint8(0); sq < 64; sq++ {
		magic := &BishopMagics[sq]

		magic.BlockerMask = genBishopMovesHQ(sq, EmptyBB, true)
		numBitsSetHigh := bits.OnesCount64(magic.BlockerMask)
		magic.Shift = uint8(64 - numBitsSetHigh)

		subsets := genSubsets(magic.BlockerMask, 1<<numBitsSetHigh)
		subsetMoves := make([]uint64, 1<<numBitsSetHigh)

		for i, subset := range subsets {
			subsetMoves[i] = genBishopMovesHQ(sq, subset, false)
		}

		searching := true
		prng.Seed(MagicSeeds[rankOf(sq)])

		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			BishopMoves[sq] = [512]uint64{}

			for i, subset := range subsets {
				index := (subset * possible_magic) >> magic.Shift
				moves := subsetMoves[i]

				if BishopMoves[sq][index] != EmptyBB && BishopMoves[sq][index] != moves {
					searching = true
					break
				}

				BishopMoves[sq][index] = moves
			}
		}
	}
}

func InitMagics() {
	genRookMagics()
	genBishopMagics()
}
