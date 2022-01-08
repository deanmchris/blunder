package engine

import (
	"fmt"
	"math/bits"
)

// bitboard.go contains the implementation of a bitboard datatype for the engine.

type Bitboard uint64

var SquareBB [64]Bitboard

func (bitboard *Bitboard) SetBit(sq uint8) {
	*bitboard |= SquareBB[sq]
}

func (bitboard *Bitboard) ClearBit(sq uint8) {
	*bitboard &= ^SquareBB[sq]
}

func (bb Bitboard) BitSet(sq uint8) bool {
	return (bb & Bitboard((0x8000000000000000 >> sq))) > 0
}

func (bitboard Bitboard) Msb() uint8 {
	return uint8(bits.LeadingZeros64(uint64(bitboard)))
}

func (bitboard *Bitboard) PopBit() uint8 {
	sq := bitboard.Msb()
	bitboard.ClearBit(sq)
	return sq
}

func (bitboard Bitboard) CountBits() int {
	return bits.OnesCount64(uint64(bitboard))
}

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
	bitboardAsString += "   "
	for fileNo := 0; fileNo < 8; fileNo++ {
		bitboardAsString += "--"
	}

	bitboardAsString += "\n    "
	for _, file := range "abcdefgh" {
		bitboardAsString += fmt.Sprintf("%c ", file)
	}
	bitboardAsString += "\n"
	return bitboardAsString
}

func init() {
	var sq uint8
	for sq = 0; sq < 64; sq++ {
		SquareBB[sq] = 0x8000000000000000 >> sq
	}
}
