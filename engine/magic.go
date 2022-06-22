package engine

import "math/bits"

// magic.go generates the magic numbers used in generating sliders moves.

/*
Magic numbers are essentially randomly generated numbers that happen to allow us
to create a hash function to index a table of pre-computed moves for bishops and
rooks.

The algorithm behind finding and using magic numbers is as follows:
For each square on the board, all the possible squares that could be occupied by a
blocker, whether friendly or not, are generated, as a bitboard. This bitboard is
usually called a blocker mask.

The blocker mask doesn't include moves on the edge squares of the board, since
even if there is a blocker at the edge of the board, the sliders will have to stop
there anyway. So a blocker being there is irrelevant.
So let's say for example that we're considering a rook sitting on D3. Here is what
the blocker mask would look like for that rook.

8 | . . . . . . . .
7 | . . . 1 . . . .
6 | . . . 1 . . . .
5 | . . . 1 . . . .
4 | . . . 1 . . . .
3 | . 1 1 . 1 1 1 .
2 | . . . 1 . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

Next, every permutation of the blocker mask is generated. The number of possible permutations
is 2^n, where n is the number of bits that are within the given blocker mask. It's common
to used a pre-intialized, static array of perumations, so the size of the array is set to be
N, where N is the largest possible number of bits for any given blocker board. For rooks,
N=4096 (2^12), and for bishops N=512 (2^9).

Lastly, the actually moves a slider has needs to be generated, given a square to move from,
and a blocker mask permutation to describe what squares are blocked. Going back to our
example, one possible permutation for our blocker mask would be:

8 | . . . . . . . .
7 | . . . . . . . .
6 | . . . 1 . . . .
5 | . . . . . . . .
4 | . . . 1 . . . .
3 | . . . . . 1 1 .
2 | . . . 1 . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

And the resulting bitboard for all the moves a rook would have would be:

8 | . . . . . . . .
7 | . . . . . . . .
6 | . . . . . . . .
5 | . . . . . . . .
4 | . . . 1 . . . .
3 | 1 1 1 . 1 1 . .
2 | . . . 1 . . . .
1 | . . . . . . . .
   ----------------
    a b c d e f g h

Once we have a blocker mask and all of its permutations, and for each permutation we can
generated the moves a slider has, magic numbers can be found by generating random 64 bit
numbers for each square, and testing if they're magic.

To test if a random number is magic, an index must first be generated using the formula:

	index = (blocker_mask_permutation * magic) >> (64 - n)

Where n is the number of bits that can be set in a blocker mask for a square.
Once we have an index, we'll know the random number is magic because the index will
be unique (i.e. not generated before by using the same magic number in the above formula
with a different blocker mask permutation).

We care about the index being unique, because the index is used to map *each* permutation,
to a moves bitboard. This way, possible moves for sliders can quickly be found by getting a
sliders blocker board permutation, given the current state of the board, and using this
permutation in the same above formula to generate a unique index that maps the permutation to
its specfic moves bitboard, for that slider, on that square.

However, a random number can actually still be magic even if it is creates an index collision
between two, or more, blocker mask permutations. But what must be the case if a collision occurs
is that all of the different blocker mask permutations that the collision occurs between must
all have the same moves bitboard.

So for every square, and random number is generated, and tested using the above two conditions.
If it passes both conditions, it's a magic number. If it doesn't, then it's not, and we need to
generate a new random number and start the process over. This is done repeadtly until we have
magic numbers for every square, for rooks and bishops.

These numbers only need to be generated once, and can then be saved into an array and used in the
program. So magic.go isn't actually used into releases of Blunder, but serves as a didatic example
of how to construct a program to generate magic numbers.
*/

var RookMagics [64]Magic
var BishopMagics [64]Magic

var RookMoves [64][4096]Bitboard
var BishopMoves [64][512]Bitboard

type Magic struct {
	MagicNo     uint64
	BlockerMask Bitboard
	Shift       uint8
}

// Optimized seeding values to find magic numbers based on file/rank, from Stockfish:
var MagicSeeds [8]uint64 = [8]uint64{728, 10316, 55013, 32803, 12281, 15100, 16645, 255}

// Generate the moves a rook has given a square and board occupancy.
func genRookMovesHQ(sq uint8, occupied Bitboard, genMask bool) Bitboard {
	slider := SquareBB[sq]
	sliderPos := slider.Msb()

	sliderBB := uint64(slider)
	occupiedBB := uint64(occupied)

	fileMask := uint64(MaskFile[fileOf(sliderPos)])
	rankMask := uint64(MaskRank[rankOf(sliderPos)])

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & rankMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & rankMask) - 2*sliderBB
	eastWestMoves := (rhs ^ lhs) & rankMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & fileMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & fileMask) - 2*sliderBB
	northSouthMoves := (rhs ^ lhs) & fileMask

	if genMask {
		// If we're told to generate a blocker board mask, the egdes of each move
		// bitboard need to cleared.
		northSouthMoves &= uint64(ClearRank[Rank1] & ClearRank[Rank8])
		eastWestMoves &= uint64(ClearFile[FileA] & ClearFile[FileH])
	}

	return Bitboard(northSouthMoves | eastWestMoves)
}

// Generate the moves a bishop has given a square and board occupancy.
func genBishopMovesHQ(sq uint8, occupied Bitboard, genMask bool) Bitboard {
	slider := SquareBB[sq]
	sliderPos := slider.Msb()

	sliderBB := uint64(slider)
	occupiedBB := uint64(occupied)

	diagonalMask := uint64(MaskDiagonal[fileOf(sliderPos)-rankOf(sliderPos)+7])
	antidiagonalMask := uint64(MaskAntidiagonal[14-(rankOf(sliderPos)+fileOf(sliderPos))])

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & diagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & diagonalMask) - 2*sliderBB
	diagonalMoves := (rhs ^ lhs) & diagonalMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & antidiagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & antidiagonalMask) - 2*sliderBB
	antidiagonalMoves := (rhs ^ lhs) & antidiagonalMask

	edges := uint64(FullBB)
	if genMask {
		// If we're told to generate a blocker board mask, the egdes of each move
		// bitboard need to cleared.
		edges = uint64(ClearFile[FileA] & ClearFile[FileH] & ClearRank[Rank1] & ClearRank[Rank8])
	}
	return Bitboard((diagonalMoves | antidiagonalMoves) & edges)
}

// Find magic numbers for rooks.
func genRookMagics() {
	var prng PseduoRandomGenerator

	for sq := uint8(0); sq < 64; sq++ {
		magic := &RookMagics[sq]

		magic.BlockerMask = genRookMovesHQ(sq, EmptyBB, true)
		no_bits := magic.BlockerMask.CountBits()
		magic.Shift = uint8(64 - no_bits)

		permutations := make([]Bitboard, 1<<no_bits)
		blockers := EmptyBB
		index := 0

		for ok := true; ok; ok = (blockers != 0) {
			permutations[index] = blockers
			index++
			blockers = (blockers - magic.BlockerMask) & magic.BlockerMask
		}

		searching := true
		prng.Seed(MagicSeeds[rankOf(sq)])

		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			RookMoves[sq] = [4096]Bitboard{}

			for idx := 0; idx < (1 << no_bits); idx++ {
				index := (uint64(permutations[idx]) * possible_magic) >> magic.Shift
				attacks := genRookMovesHQ(sq, permutations[idx], false)

				if RookMoves[sq][index] != EmptyBB && RookMoves[sq][index] != attacks {
					searching = true
					break
				}
				RookMoves[sq][index] = attacks
			}
		}
	}
}

// Find magic numbers for bishops.
func genBishopMagics() {
	var prng PseduoRandomGenerator
	for sq := uint8(0); sq < 64; sq++ {
		magic := &BishopMagics[sq]

		magic.BlockerMask = genBishopMovesHQ(sq, 0, true)
		no_bits := magic.BlockerMask.CountBits()
		magic.Shift = uint8(64 - no_bits)

		permutations := make([]Bitboard, 1<<no_bits)
		blockers := EmptyBB
		index := 0

		for ok := true; ok; ok = (blockers != 0) {
			permutations[index] = blockers
			index++
			blockers = (blockers - magic.BlockerMask) & magic.BlockerMask
		}

		searching := true
		prng.Seed(MagicSeeds[rankOf(sq)])

		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			BishopMoves[sq] = [512]Bitboard{}

			for idx := 0; idx < (1 << no_bits); idx++ {
				index := (uint64(permutations[idx]) * possible_magic) >> magic.Shift
				attacks := genBishopMovesHQ(sq, permutations[idx], false)

				if BishopMoves[sq][index] != EmptyBB && BishopMoves[sq][index] != attacks {
					searching = true
					break
				}
				BishopMoves[sq][index] = attacks
			}
		}
	}
}
