package engine

// bitboard.go contains the implementation of a bitboard datatype for the engine.

import (
	"fmt"
	"math/bits"
)

// A type representing a bitboard, which is a unsigned 64-bit number. Blunder's
// bitboard representation has the most significant bit being A1 and the least significant
// bit being H8.
type Bitboard uint64

// A constant representing a bitboard with every square set and every square
// empty.
const FullBB Bitboard = 0xffffffffffffffff
const EmptyBB Bitboard = 0x0

// A global constant where each entry represents a square on the chess board,
// and each entry contains a bitboard with the bit set high at that square.
// An extra entry is given so that the invalid square constant NoSq can be
// indexed into the table without the program crashing.
var SquareBB [65]Bitboard

// Set the bit at given square.
func (bitboard *Bitboard) SetBit(sq uint8) {
	*bitboard |= SquareBB[sq]
}

// Clear the bit at given square.
func (bitboard *Bitboard) ClearBit(sq uint8) {
	*bitboard &= ^SquareBB[sq]
}

// Test whether the bit of the given bitboard at the given
// position is set.
func (bb Bitboard) BitSet(sq uint8) bool {
	return (bb & SquareBB[sq]) != 0
}

// Get the position of the MSB of the given bitboard.
func (bitboard Bitboard) Msb() uint8 {
	return uint8(bits.LeadingZeros64(uint64(bitboard)))
}

// Get the position of the MSB of the given bitboard,
// and clear the MSB.
func (bitboard *Bitboard) PopBit() uint8 {
	sq := bitboard.Msb()
	bitboard.ClearBit(sq)
	return sq
}

// Count the bits in a given bitboard using the SWAR-popcount
// algorithm for 64-bit integers.
func (bitboard Bitboard) CountBits() int {
	return bits.OnesCount64(uint64(bitboard))
}

// Return a string representation of the given bitboard
func (bitboard Bitboard) String() (bitboardAsString string) {
	bitstring := fmt.Sprintf("%064b\n", bitboard)
	bitboardAsString += "\n"
	for rankStartPos := 56; rankStartPos >= 0; rankStartPos -= 8 {
		bitboardAsString += fmt.Sprintf("%v | ", (rankStartPos/8)+1)
		for index := rankStartPos; index < rankStartPos+8; index++ {
			squareChar := bitstring[index]
			if squareChar == '0' {
				squareChar = '.'
			}
			bitboardAsString += fmt.Sprintf("%c ", squareChar)
		}
		bitboardAsString += "\n"
	}

	bitboardAsString += "   ----------------"
	bitboardAsString += "\n    a b c d e f g h"
	bitboardAsString += "\n"
	return bitboardAsString
}

// Initialize the bitboard constants.
func InitBitboards() {
	var sq uint8
	for sq = 0; sq < 65; sq++ {
		SquareBB[sq] = 0x8000000000000000 >> sq
	}
}
