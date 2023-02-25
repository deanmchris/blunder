package engine

import (
	"fmt"
	"math/bits"
)

const (
	MostSigBitBB uint64 = 0x8000000000000000
	FullBB       uint64 = 0xffffffffffffffff
	EmptyBB      uint64 = 0x0
)

var SquareBB [65]uint64

func setBit(bb *uint64, sq uint8) {
	*bb |= MostSigBitBB >> sq
}

func clearBit(bb *uint64, sq uint8) {
	*bb ^= MostSigBitBB >> sq
}

func bitIsSet(bb uint64, sq uint8) bool {
	return (bb & (MostSigBitBB >> sq)) != 0
}

func bitScan(bb uint64) uint8 {
	return uint8(bits.LeadingZeros64(bb))
}

func BitScanAndClear(bb *uint64) uint8 {
	sq := uint8(bits.LeadingZeros64(*bb))
	*bb ^= MostSigBitBB >> sq
	return sq
}

func PrintBitboard(bb uint64) {
	bbStr := fmt.Sprintf("%064b", bb)
	for i := 56; i >= 0; i -= 8 {
		fmt.Printf("%d| ", i/8+1)
		for j := i; j < i+8; j++ {
			bit := bbStr[j]
			if bit == '0' {
				bit = '.'
			}
			fmt.Printf("%c ", bit)
		}
		fmt.Println()
	}
	fmt.Println("  ----------------\n   a b c d e f g h")
}

func InitBitboard() {
	for i := 0; i < 64; i++ {
		SquareBB[i] = MostSigBitBB >> i
	}
}
