package extra

import (
	"blunder/engine"
	"fmt"
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
These numbers only need to be generated once, and can then be saved into an array and used in the
program. So magic.go isn't actually used into releases of Blunder, but serves as a didatic example
of how to construct a program to generate magic numbers.
*/

type Bitboard = engine.Bitboard

const EmptyEntry Bitboard = 0

var MagicSeeds [8]uint64 = [8]uint64{728, 10316, 55013, 32803, 12281, 15100, 16645, 255}

// A struct to hold information regarding a magic number
// for a rook or bishop on a particular square.
type Magic struct {
	MagicNo uint64
	Mask    Bitboard
	Shift   uint8
}

var RookMagics [64]Magic
var BishopMagics [64]Magic

var RookAttacks [64][4096]Bitboard
var BishopAttacks [64][512]Bitboard

// An implementation of a xorshift pseudo-random number
// generator for 64 bit numbers, based on the implementation
// by Stockfish.
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

// Find magic numbers for rooks.
func genRookMagics() {
	var prng PseduoRandomGenerator

	var sq uint8
	for sq = 0; sq < 64; sq++ {
		magic := &RookMagics[sq]

		magic.Mask = engine.GenRookMasks(sq)
		no_bits := magic.Mask.CountBits()
		magic.Shift = uint8(64 - no_bits)

		permutations := make([]Bitboard, 1<<no_bits)
		var blockers engine.Bitboard
		var index int

		for ok := true; ok; ok = (blockers != 0) {
			permutations[index] = blockers
			index++
			blockers = (blockers - magic.Mask) & magic.Mask
		}

		searching := true
		tries := 0

		prng.Seed(MagicSeeds[engine.RankOf(sq)])
		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			RookAttacks[sq] = [4096]Bitboard{}

			for idx := 0; idx < (1 << no_bits); idx++ {
				index := (uint64(permutations[idx]) * possible_magic) >> magic.Shift
				attacks := engine.GenRookAttacks(sq, permutations[idx])

				if RookAttacks[sq][index] != EmptyEntry && RookAttacks[sq][index] != attacks {
					searching = true
					tries++
					break
				}
				RookAttacks[sq][index] = attacks
			}
		}
		fmt.Printf("Magic 0x%x for square %d, found in %d tries\n", magic.MagicNo, sq, tries)
	}
}

// Find magic numbers for bishops.
func genBishopMagics() {
	var prng PseduoRandomGenerator

	var sq uint8
	for sq = 0; sq < 64; sq++ {
		magic := &BishopMagics[sq]

		magic.Mask = engine.GenBishopMasks(sq)
		no_bits := magic.Mask.CountBits()
		magic.Shift = uint8(64 - no_bits)

		permutations := make([]Bitboard, 1<<no_bits)
		var blockers Bitboard
		var index int

		for ok := true; ok; ok = (blockers != 0) {
			permutations[index] = blockers
			index++
			blockers = (blockers - magic.Mask) & magic.Mask
		}

		searching := true
		tries := 0

		prng.Seed(MagicSeeds[engine.RankOf(sq)])
		for searching {
			possible_magic := prng.SparseRandom64()
			magic.MagicNo = possible_magic
			searching = false

			BishopAttacks[sq] = [512]Bitboard{}

			for idx := 0; idx < (1 << no_bits); idx++ {
				index := (uint64(permutations[idx]) * possible_magic) >> magic.Shift
				attacks := engine.GenBishopAttacks(sq, permutations[idx])

				if BishopAttacks[sq][index] != EmptyEntry && BishopAttacks[sq][index] != attacks {
					searching = true
					tries++
					break
				}
				BishopAttacks[sq][index] = attacks
			}
		}
		fmt.Printf("Magic 0x%x for square %d, found in %d tries\n", magic.MagicNo, sq, tries)
	}
}

func GenMagics() {
	fmt.Println("Generating rook magic numbers...")
	genRookMagics()
	fmt.Print("\nGenerating bishop magic numbers...\n")
	genBishopMagics()

	fmt.Println("\nvar RookMagicNumbers [64]uint64 = [64]uint64{")
	fmt.Print("    ")

	for idx, magic := range RookMagics {
		if idx%4 == 0 && idx != 0 {
			fmt.Println()
			fmt.Print("    ")
		}
		fmt.Printf("0x%x, ", magic.MagicNo)
	}
	fmt.Print("\n}\n\n")

	fmt.Println("var BishopMagicNumbers [64]uint64 = [64]uint64{")
	fmt.Print("    ")

	for idx, magic := range BishopMagics {
		if idx%4 == 0 && idx != 0 {
			fmt.Println()
			fmt.Print("    ")
		}
		fmt.Printf("0x%x, ", magic.MagicNo)
	}
	fmt.Print("\n}\n\n")
}
