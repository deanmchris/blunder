package engine

// A file containg various bitboard releated utilities.

import (
	"fmt"
	"math/bits"
)

var SquareBB [65]uint64

func init() {
	for sq := 0; sq < 65; sq++ {
		SquareBB[sq] = 0x8000000000000000 >> sq
	}
}

// Pretty print a bitboard.
func PrintBitboard(bitboard uint64) {
	bitstring := fmt.Sprintf("%064b\n", bitboard)
	fmt.Println()
	for rankStartPos := 56; rankStartPos >= 0; rankStartPos -= 8 {
		fmt.Printf("%v | ", (rankStartPos/8)+1)
		for index := rankStartPos; index < rankStartPos+8; index++ {
			squareChar := bitstring[index]
			if squareChar == '0' {
				squareChar = '.'
			}
			fmt.Printf("%c ", squareChar)
		}
		fmt.Println()
	}
	fmt.Print("   ")
	for fileNo := 0; fileNo < 8; fileNo++ {
		fmt.Print("--")
	}

	fmt.Print("\n    ")
	for _, file := range "abcdefgh" {
		fmt.Printf("%c ", file)
	}
	fmt.Println()
}

// Set the bit of the given bitbord at the given position.
func setBit(bb *uint64, sq int) {
	*bb |= SquareBB[sq]
}

// Clear the bit of the given bitbord at the given position.
func clearBit(bb *uint64, sq int) {
	*bb &= ^SquareBB[sq]
}

// Test whether the bit of the given bitbord at the given
// position is set.
func bitSet(bb uint64, sq int) bool {
	return (bb & (0x8000000000000000 >> sq)) > 0
}

// Get the position of the MSB of the given bitboard.
func msb(bb uint64) int {
	return bits.LeadingZeros64(bb)
}

// Get the position of the LSB of the given bitboard,
// a bitboard with only the LSB set, and clear the LSB.
func PopBit(bb *uint64) int {
	sq := msb(*bb)
	clearBit(bb, sq)
	return sq
}

// Count the bits in a given bitboard using the SWAR-popcount
// algorithm for 64-bit integers.
func CountBits(u uint64) int {
	u = u - ((u >> 1) & 0x5555555555555555)
	u = (u & 0x3333333333333333) + ((u >> 2) & 0x3333333333333333)
	u = (u + (u >> 4)) & 0x0f0f0f0f0f0f0f0f
	u = (u * 0x0101010101010101) >> 56
	return int(u)
}
