package engine

import (
	"math/bits"
)

/*

The file containg the code for generating Blunder's magic bitboard numbers.

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
is 2^n, where n is the number of bits that can be set in the blocker mask (in our example, n=12,
so 4096 possibe permutations).

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

To test if a random number is magic, an index must first be generated using the formula

index = (blocker_mask_permutation * magic) >> (64 - n)

where n is the number of bits that can be set in a blocker mask for a square.

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
*/

const EmptyEntry uint64 = 0

var RookMagics [64]Magic
var BishopMagics [64]Magic

var RookAttacks [64][4096]uint64
var BishopAttacks [64][512]uint64

var MagicSeeds [8]uint64 = [8]uint64{728, 10316, 55013, 32803, 12281, 15100, 16645, 255}

// A struct to hold information regarding a magic number
// for a rook or bishop on a particular square.
type Magic struct {
	MagicNo uint64
	Mask    uint64
	Shift   uint64
}

// Generate a blocker mask for rooks.
func genRookMasks(sq int) uint64 {
	var occupiedBB uint64
	sliderBB := SquareBB[sq]
	sliderPos := msb(sliderBB)

	fileMask := MaskFile[FileOf(sliderPos)]
	rankMask := MaskRank[RankOf(sliderPos)]

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & rankMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & rankMask) - 2*sliderBB
	eastWestMoves := ((rhs ^ lhs) & rankMask) & (ClearFile[FileA] & ClearFile[FileH])

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & fileMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & fileMask) - 2*sliderBB
	northSouthMoves := ((rhs ^ lhs) & fileMask) & (ClearRank[Rank1] & ClearRank[Rank8])

	return northSouthMoves | eastWestMoves
}

// Generate a blocker mask for bishops.
func genBishopMasks(sq int) uint64 {
	var occupiedBB uint64
	sliderBB := SquareBB[sq]
	sliderPos := msb(sliderBB)

	diagonalMask := MaskDiagonal[FileOf(sliderPos)-RankOf(sliderPos)+7]
	antidiagonalMask := MaskAntidiagonal[14-(RankOf(sliderPos)+FileOf(sliderPos))]

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & diagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & diagonalMask) - 2*sliderBB
	diagonalMoves := (rhs ^ lhs) & diagonalMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & antidiagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & antidiagonalMask) - 2*sliderBB
	antidiagonalMoves := (rhs ^ lhs) & antidiagonalMask

	edges := ClearFile[FileA] & ClearFile[FileH] & ClearRank[Rank1] & ClearRank[Rank8]
	return (diagonalMoves | antidiagonalMoves) & edges
}

func genRookAttacks(sq int, occupiedBB uint64) uint64 {
	sliderBB := SquareBB[sq]
	sliderPos := msb(sliderBB)

	fileMask := MaskFile[FileOf(sliderPos)]
	rankMask := MaskRank[RankOf(sliderPos)]

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & rankMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & rankMask) - 2*sliderBB
	eastWestMoves := (rhs ^ lhs) & rankMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & fileMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & fileMask) - 2*sliderBB
	northSouthMoves := (rhs ^ lhs) & fileMask

	return northSouthMoves | eastWestMoves
}

// Generate the moves a bishop has given a square and board occupancy.
func genBishopAttacks(sq int, occupiedBB uint64) uint64 {
	sliderBB := SquareBB[sq]
	sliderPos := msb(sliderBB)

	diagonalMask := MaskDiagonal[FileOf(sliderPos)-RankOf(sliderPos)+7]
	antidiagonalMask := MaskAntidiagonal[14-(RankOf(sliderPos)+FileOf(sliderPos))]

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & diagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & diagonalMask) - 2*sliderBB
	diagonalMoves := (rhs ^ lhs) & diagonalMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & antidiagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & antidiagonalMask) - 2*sliderBB
	antidiagonalMoves := (rhs ^ lhs) & antidiagonalMask

	return diagonalMoves | antidiagonalMoves
}

// Find magic numbers for rooks.
func genRookMagics() {
	var prng PseduoRandomGenerator

	for sq := 0; sq < 64; sq++ {
		magic := &RookMagics[sq]

		magic.Mask = genRookMasks(sq)
		no_bits := CountBits(magic.Mask)
		magic.Shift = 64 - uint64(no_bits)

		permutations := make([]uint64, 1<<no_bits)
		var blockers uint64
		var index int

		for ok := true; ok; ok = (blockers != 0) {
			permutations[index] = blockers
			index++
			blockers = (blockers - magic.Mask) & magic.Mask
		}

		searching := true
		tries := 0

		prng.Seed(MagicSeeds[RankOf(sq)])
		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			RookAttacks[sq] = [4096]uint64{}

			for idx := 0; idx < (1 << no_bits); idx++ {
				index := (permutations[idx] * possible_magic) >> magic.Shift
				attacks := genRookAttacks(sq, permutations[idx])

				if RookAttacks[sq][index] != EmptyEntry && RookAttacks[sq][index] != attacks {
					searching = true
					tries++
					break
				}
				RookAttacks[sq][index] = attacks
			}
		}
		// fmt.Printf("Magic 0x%x for square %d, found in %d tries\n", magic.MagicNo, sq, tries)
	}
}

// Find magic numbers for bishops.
func genBishopMagics() {
	var prng PseduoRandomGenerator

	for sq := 0; sq < 64; sq++ {
		magic := &BishopMagics[sq]

		magic.Mask = genBishopMasks(sq)
		no_bits := CountBits(magic.Mask)
		magic.Shift = 64 - uint64(no_bits)

		permutations := make([]uint64, 1<<no_bits)
		var blockers uint64
		var index int

		for ok := true; ok; ok = (blockers != 0) {
			permutations[index] = blockers
			index++
			blockers = (blockers - magic.Mask) & magic.Mask
		}

		searching := true
		tries := 0

		prng.Seed(MagicSeeds[RankOf(sq)])
		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			BishopAttacks[sq] = [512]uint64{}

			for idx := 0; idx < (1 << no_bits); idx++ {
				index := (permutations[idx] * possible_magic) >> magic.Shift
				attacks := genBishopAttacks(sq, permutations[idx])

				if BishopAttacks[sq][index] != EmptyEntry && BishopAttacks[sq][index] != attacks {
					searching = true
					tries++
					break
				}
				BishopAttacks[sq][index] = attacks
			}
		}
		// fmt.Printf("Magic 0x%x for square %d, found in %d tries\n", magic.MagicNo, sq, tries)
	}
}

func init() {
	//fmt.Print("Finding rook magics....\n\n")
	genRookMagics()
	//fmt.Print("\nFinding bishop magics....\n\n")
	genBishopMagics()
}
