package engine

import (
	"fmt"
	"math/bits"
)

// bitboard.go defines the bitboard type and its various methods.

var SquareBB [65]Bitboard

func init() {
	for sq := 0; sq < 65; sq++ {
		SquareBB[sq] = 0x8000000000000000 >> sq
	}
}

type Bitboard uint64

// Pretty print a bitboard.
func (bb Bitboard) String() (boardStr string) {
	bitstring := fmt.Sprintf("%064b\n", bb)
	boardStr += "\n"
	for rankStartPos := 56; rankStartPos >= 0; rankStartPos -= 8 {
		boardStr += fmt.Sprintf("%v | ", (rankStartPos/8)+1)
		for index := rankStartPos; index < rankStartPos+8; index++ {
			squareChar := bitstring[index]
			if squareChar == '0' {
				squareChar = '.'
			}
			boardStr += fmt.Sprintf("%c ", squareChar)
		}
		boardStr += "\n"
	}
	boardStr += "   "
	for fileNo := 0; fileNo < 8; fileNo++ {
		boardStr += "--"
	}

	boardStr += "\n    "
	for _, file := range "abcdefgh" {
		boardStr += fmt.Sprintf("%c ", file)
	}
	boardStr += "\n"
	return boardStr
}

// Set the bit of the given bitbord at the given position.
func (bb *Bitboard) SetBit(sq uint8) {
	*bb |= SquareBB[sq]
}

// Clear the bit of the given bitbord at the given position.
func (bb *Bitboard) ClearBit(sq uint8) {
	*bb &= ^SquareBB[sq]
}

// Test whether the bit of the given bitbord at the given
// position is set.
func (bb Bitboard) BitSet(sq uint8) bool {
	return (bb & Bitboard((0x8000000000000000 >> sq))) > 0
}

// Get the position of the MSB of the given bitboard.
func (bb Bitboard) Msb() uint8 {
	return uint8(bits.LeadingZeros64(uint64(bb)))
}

// Get the position of the LSB of the given bitboard,
// a bitboard with only the LSB set, and clear the LSB.
func (bb *Bitboard) PopBit() uint8 {
	sq := bb.Msb()
	bb.ClearBit(sq)
	return sq
}

// Count the bits in a given bitboard using the SWAR-popcount
// algorithm for 64-bit integers.
func (bb Bitboard) CountBits() int {
	u := uint64(bb)
	u = u - ((u >> 1) & 0x5555555555555555)
	u = (u & 0x3333333333333333) + ((u >> 2) & 0x3333333333333333)
	u = (u + (u >> 4)) & 0x0f0f0f0f0f0f0f0f
	u = (u * 0x0101010101010101) >> 56
	return int(u)
}
