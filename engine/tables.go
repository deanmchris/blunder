package engine

import (
	"math/bits"
)

// A file containing various precomputed tables used
// in the engine.

const (
	Rank1 = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

const (
	FileA = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

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

var RookMagicNumbers [64]uint64 = [64]uint64{
	0x2220a09401006042, 0x1022000403018802, 0x2422006150081c02, 0x400200102014081a,
	0x80020410010008a1, 0x2008102003004009, 0x839100214001, 0x6101310024408005,
	0x820018104004600, 0x804021048430400, 0x100020004008080, 0x904802402480080,
	0x1800900100982100, 0x2028200082100080, 0x40028900412500, 0x400080006080,
	0x42041820011, 0x8002000804220011, 0xa000850020024, 0x4202040008008080,
	0x8480120140220008, 0x10008020008052, 0x4020006250004000, 0x4c40004020808002,
	0x4c984082000b04, 0x24814001021, 0x102011052000468, 0x105000801001004,
	0x8201089001002101, 0x1012811000806000, 0x220100040400020, 0x80814002a0800081,
	0x20040420000830c, 0x80a1016400104208, 0x8440020080800400, 0x4008080080140080,
	0x4001300080080080, 0x8120011010040200, 0x1000842300400100, 0x4000400080002085,
	0x20004004081, 0x800040001081290, 0xb001010004001802, 0x80000d0011000800,
	0x44a0020120040, 0x142020040118020, 0x888808020014000, 0x90402c8000804000,
	0x829000042128300, 0xa9cc00150a080410, 0x4001002842440100, 0x2002052000488,
	0x803001000204904, 0x800c802002100880, 0x4014804001842004, 0x408000e0400080,
	0x9080014080042100, 0xa00041122001188, 0x41001281001c0008, 0x2100030010080084,
	0x80080004801000, 0x200181202204080, 0x8040002000409004, 0x380004003b0a080,
}

var BishopMagicNumbers [64]uint64 = [64]uint64{
	0x40484802902044, 0x210401044110050, 0x10006020322084, 0x1a4a0410021200,
	0x800080000c208800, 0x2810084008800, 0x80012608025800, 0x220204a00800,
	0x34840082020002, 0x40184123020024, 0x403120041000800c, 0x6800801810242401,
	0xa005020442088020, 0x12000100a8040020, 0x382004108292000, 0xc002080404040400,
	0x550008100480101, 0x80a28020122c0, 0x504090045040200, 0x828140892000c00,
	0x41144000801, 0x80c0c40c8011000, 0x8004020242201020, 0x8004020242201020,
	0xc010220920198, 0x404032404014400, 0x4801080200802200, 0x8040038020020220,
	0x4040020080080080, 0x44020300081848, 0x8004014c00600400, 0x41041381202000,
	0x4828504005040211, 0x10840101820880, 0x9420121c1101c, 0x900401c004049,
	0x4048080004820002, 0x81004c0018080313, 0x10102858090121, 0x20640050100a00,
	0x8403429080820, 0x1009610822080, 0x32400608200412, 0x9004101202020240,
	0x1108000420401000, 0x114001244008203, 0x8104868204040412, 0x8403429080820,
	0x80148054262004, 0x4040008218200400, 0x10120802080a81, 0x1402440421000210,
	0x80122082080440, 0x8080880801082000, 0x2444681203d200, 0x312208080880,
	0x1008044200440, 0x958808402025, 0x206012462000121, 0xc042000800000,
	0xc014040288010002, 0x2218080110210404, 0x948110c0b2081, 0xa20040301410a,
}

var KnightMoves [64]Bitboard = [64]Bitboard{
	0x20400000000000, 0x10a00000000000, 0x88500000000000, 0x44280000000000, 0x22140000000000, 0x110a0000000000, 0x8050000000000, 0x4020000000000,
	0x2000204000000000, 0x100010a000000000, 0x8800885000000000, 0x4400442800000000, 0x2200221400000000, 0x1100110a00000000, 0x800080500000000, 0x400040200000000,
	0x4020002040000000, 0xa0100010a0000000, 0x5088008850000000, 0x2844004428000000, 0x1422002214000000, 0xa1100110a000000, 0x508000805000000, 0x204000402000000,
	0x40200020400000, 0xa0100010a00000, 0x50880088500000, 0x28440044280000, 0x14220022140000, 0xa1100110a0000, 0x5080008050000, 0x2040004020000,
	0x402000204000, 0xa0100010a000, 0x508800885000, 0x284400442800, 0x142200221400, 0xa1100110a00, 0x50800080500, 0x20400040200,
	0x4020002040, 0xa0100010a0, 0x5088008850, 0x2844004428, 0x1422002214, 0xa1100110a, 0x508000805, 0x204000402,
	0x40200020, 0xa0100010, 0x50880088, 0x28440044, 0x14220022, 0xa110011, 0x5080008, 0x2040004,
	0x402000, 0xa01000, 0x508800, 0x284400, 0x142200, 0xa1100, 0x50800, 0x20400,
}

var KingMoves [64]Bitboard = [64]Bitboard{
	0x40c0000000000000, 0xa0e0000000000000, 0x5070000000000000, 0x2838000000000000, 0x141c000000000000, 0xa0e000000000000, 0x507000000000000, 0x203000000000000,
	0xc040c00000000000, 0xe0a0e00000000000, 0x7050700000000000, 0x3828380000000000, 0x1c141c0000000000, 0xe0a0e0000000000, 0x705070000000000, 0x302030000000000,
	0xc040c000000000, 0xe0a0e000000000, 0x70507000000000, 0x38283800000000, 0x1c141c00000000, 0xe0a0e00000000, 0x7050700000000, 0x3020300000000,
	0xc040c0000000, 0xe0a0e0000000, 0x705070000000, 0x382838000000, 0x1c141c000000, 0xe0a0e000000, 0x70507000000, 0x30203000000,
	0xc040c00000, 0xe0a0e00000, 0x7050700000, 0x3828380000, 0x1c141c0000, 0xe0a0e0000, 0x705070000, 0x302030000,
	0xc040c000, 0xe0a0e000, 0x70507000, 0x38283800, 0x1c141c00, 0xe0a0e00, 0x7050700, 0x3020300,
	0xc040c0, 0xe0a0e0, 0x705070, 0x382838, 0x1c141c, 0xe0a0e, 0x70507, 0x30203,
	0xc040, 0xe0a0, 0x7050, 0x3828, 0x1c14, 0xe0a, 0x705, 0x302,
}

var PawnAttacks [2][64]Bitboard = [2][64]Bitboard{
	{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x4000000000000000, 0xa000000000000000, 0x5000000000000000, 0x2800000000000000, 0x1400000000000000, 0xa00000000000000, 0x500000000000000, 0x200000000000000,
		0x40000000000000, 0xa0000000000000, 0x50000000000000, 0x28000000000000, 0x14000000000000, 0xa000000000000, 0x5000000000000, 0x2000000000000,
		0x400000000000, 0xa00000000000, 0x500000000000, 0x280000000000, 0x140000000000, 0xa0000000000, 0x50000000000, 0x20000000000,
		0x4000000000, 0xa000000000, 0x5000000000, 0x2800000000, 0x1400000000, 0xa00000000, 0x500000000, 0x200000000,
		0x40000000, 0xa0000000, 0x50000000, 0x28000000, 0x14000000, 0xa000000, 0x5000000, 0x2000000,
		0x400000, 0xa00000, 0x500000, 0x280000, 0x140000, 0xa0000, 0x50000, 0x20000,
		0x4000, 0xa000, 0x5000, 0x2800, 0x1400, 0xa00, 0x500, 0x200,
	},

	{
		0x40000000000000, 0xa0000000000000, 0x50000000000000, 0x28000000000000, 0x14000000000000, 0xa000000000000, 0x5000000000000, 0x2000000000000,
		0x400000000000, 0xa00000000000, 0x500000000000, 0x280000000000, 0x140000000000, 0xa0000000000, 0x50000000000, 0x20000000000,
		0x4000000000, 0xa000000000, 0x5000000000, 0x2800000000, 0x1400000000, 0xa00000000, 0x500000000, 0x200000000,
		0x40000000, 0xa0000000, 0x50000000, 0x28000000, 0x14000000, 0xa000000, 0x5000000, 0x2000000,
		0x400000, 0xa00000, 0x500000, 0x280000, 0x140000, 0xa0000, 0x50000, 0x20000,
		0x4000, 0xa000, 0x5000, 0x2800, 0x1400, 0xa00, 0x500, 0x200,
		0x40, 0xa0, 0x50, 0x28, 0x14, 0xa, 0x5, 0x2,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	},
}

var PawnPushes [2][64]Bitboard = [2][64]Bitboard{
	{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x8000000000000000, 0x4000000000000000, 0x2000000000000000, 0x1000000000000000, 0x800000000000000, 0x400000000000000, 0x200000000000000, 0x100000000000000,
		0x80000000000000, 0x40000000000000, 0x20000000000000, 0x10000000000000, 0x8000000000000, 0x4000000000000, 0x2000000000000, 0x1000000000000,
		0x800000000000, 0x400000000000, 0x200000000000, 0x100000000000, 0x80000000000, 0x40000000000, 0x20000000000, 0x10000000000,
		0x8000000000, 0x4000000000, 0x2000000000, 0x1000000000, 0x800000000, 0x400000000, 0x200000000, 0x100000000,
		0x80000000, 0x40000000, 0x20000000, 0x10000000, 0x8000000, 0x4000000, 0x2000000, 0x1000000,
		0x800000, 0x400000, 0x200000, 0x100000, 0x80000, 0x40000, 0x20000, 0x10000,
		0x8000, 0x4000, 0x2000, 0x1000, 0x800, 0x400, 0x200, 0x100,
	},
	{
		0x80000000000000, 0x40000000000000, 0x20000000000000, 0x10000000000000, 0x8000000000000, 0x4000000000000, 0x2000000000000, 0x1000000000000,
		0x800000000000, 0x400000000000, 0x200000000000, 0x100000000000, 0x80000000000, 0x40000000000, 0x20000000000, 0x10000000000,
		0x8000000000, 0x4000000000, 0x2000000000, 0x1000000000, 0x800000000, 0x400000000, 0x200000000, 0x100000000,
		0x80000000, 0x40000000, 0x20000000, 0x10000000, 0x8000000, 0x4000000, 0x2000000, 0x1000000,
		0x800000, 0x400000, 0x200000, 0x100000, 0x80000, 0x40000, 0x20000, 0x10000,
		0x8000, 0x4000, 0x2000, 0x1000, 0x800, 0x400, 0x200, 0x100,
		0x80, 0x40, 0x20, 0x10, 0x8, 0x4, 0x2, 0x1,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	},
}

var MaskDiagonal [15]Bitboard = [15]Bitboard{
	0x80,
	0x8040,
	0x804020,
	0x80402010,
	0x8040201008,
	0x804020100804,
	0x80402010080402,
	0x8040201008040201,
	0x4020100804020100,
	0x2010080402010000,
	0x1008040201000000,
	0x804020100000000,
	0x402010000000000,
	0x201000000000000,
	0x100000000000000,
}
var MaskAntidiagonal [15]Bitboard = [15]Bitboard{
	0x1,
	0x102,
	0x10204,
	0x1020408,
	0x102040810,
	0x10204081020,
	0x1020408102040,
	0x102040810204080,
	0x204081020408000,
	0x408102040800000,
	0x810204080000000,
	0x1020408000000000,
	0x2040800000000000,
	0x4080000000000000,
	0x8000000000000000,
}

var ClearRank [8]Bitboard = [8]Bitboard{
	0xffffffffffffff,
	0xff00ffffffffffff,
	0xffff00ffffffffff,
	0xffffff00ffffffff,
	0xffffffff00ffffff,
	0xffffffffff00ffff,
	0xffffffffffff00ff,
	0xffffffffffffff00,
}

var ClearFile [8]Bitboard = [8]Bitboard{
	0x7f7f7f7f7f7f7f7f,
	0xbfbfbfbfbfbfbfbf,
	0xdfdfdfdfdfdfdfdf,
	0xefefefefefefefef,
	0xf7f7f7f7f7f7f7f7,
	0xfbfbfbfbfbfbfbfb,
	0xfdfdfdfdfdfdfdfd,
	0xfefefefefefefefe,
}

var MaskRank [8]Bitboard = [8]Bitboard{
	0xff00000000000000,
	0xff000000000000,
	0xff0000000000,
	0xff00000000,
	0xff000000,
	0xff0000,
	0xff00,
	0xff,
}

var MaskFile [8]Bitboard = [8]Bitboard{
	0x8080808080808080,
	0x4040404040404040,
	0x2020202020202020,
	0x1010101010101010,
	0x808080808080808,
	0x404040404040404,
	0x202020202020202,
	0x101010101010101,
}

// Generate a blocker mask for rooks.
func GenRookMasks(sq uint8) Bitboard {
	slider := SquareBB[sq]
	sliderPos := slider.Msb()

	sliderBB := uint64(slider)
	var occupiedBB uint64

	fileMask := uint64(MaskFile[FileOf(sliderPos)])
	rankMask := uint64(MaskRank[RankOf(sliderPos)])

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & rankMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & rankMask) - 2*sliderBB
	eastWestMoves := ((rhs ^ lhs) & rankMask) & (uint64(ClearFile[FileA] & ClearFile[FileH]))

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & fileMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & fileMask) - 2*sliderBB
	northSouthMoves := ((rhs ^ lhs) & fileMask) & (uint64(ClearRank[Rank1] & ClearRank[Rank8]))

	return Bitboard(northSouthMoves | eastWestMoves)
}

// Generate a blocker mask for bishops.
func GenBishopMasks(sq uint8) Bitboard {
	slider := SquareBB[sq]
	sliderPos := slider.Msb()

	sliderBB := uint64(slider)
	var occupiedBB uint64

	diagonalMask := uint64(MaskDiagonal[FileOf(sliderPos)-RankOf(sliderPos)+7])
	antidiagonalMask := uint64(MaskAntidiagonal[14-(RankOf(sliderPos)+FileOf(sliderPos))])

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & diagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & diagonalMask) - 2*sliderBB
	diagonalMoves := (rhs ^ lhs) & diagonalMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & antidiagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & antidiagonalMask) - 2*sliderBB
	antidiagonalMoves := (rhs ^ lhs) & antidiagonalMask

	edges := uint64(ClearFile[FileA] & ClearFile[FileH] & ClearRank[Rank1] & ClearRank[Rank8])
	return Bitboard((diagonalMoves | antidiagonalMoves) & edges)
}

func GenRookAttacks(sq uint8, occupied Bitboard) Bitboard {
	slider := SquareBB[sq]
	sliderPos := slider.Msb()

	sliderBB := uint64(slider)
	occupiedBB := uint64(occupied)

	fileMask := uint64(MaskFile[FileOf(sliderPos)])
	rankMask := uint64(MaskRank[RankOf(sliderPos)])

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & rankMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & rankMask) - 2*sliderBB
	eastWestMoves := (rhs ^ lhs) & rankMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & fileMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & fileMask) - 2*sliderBB
	northSouthMoves := (rhs ^ lhs) & fileMask

	return Bitboard(northSouthMoves | eastWestMoves)
}

// Generate the moves a bishop has given a square and board occupancy.
func GenBishopAttacks(sq uint8, occupied Bitboard) Bitboard {
	slider := SquareBB[sq]
	sliderPos := slider.Msb()

	sliderBB := uint64(slider)
	occupiedBB := uint64(occupied)

	diagonalMask := uint64(MaskDiagonal[FileOf(sliderPos)-RankOf(sliderPos)+7])
	antidiagonalMask := uint64(MaskAntidiagonal[14-(RankOf(sliderPos)+FileOf(sliderPos))])

	rhs := bits.Reverse64(bits.Reverse64((occupiedBB & diagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs := (occupiedBB & diagonalMask) - 2*sliderBB
	diagonalMoves := (rhs ^ lhs) & diagonalMask

	rhs = bits.Reverse64(bits.Reverse64((occupiedBB & antidiagonalMask)) - (2 * bits.Reverse64(sliderBB)))
	lhs = (occupiedBB & antidiagonalMask) - 2*sliderBB
	antidiagonalMoves := (rhs ^ lhs) & antidiagonalMask

	return Bitboard(diagonalMoves | antidiagonalMoves)
}

// Initalize the rook tables holding Magic structs and rook attacks
func initRookMagics() {
	var sq uint8
	for sq = 0; sq < 64; sq++ {
		magic := &RookMagics[sq]

		magic.MagicNo = RookMagicNumbers[sq]
		magic.Mask = GenRookMasks(sq)
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

		for idx := 0; idx < (1 << no_bits); idx++ {
			index := (permutations[idx] * Bitboard(magic.MagicNo)) >> magic.Shift
			attacks := GenRookAttacks(sq, permutations[idx])
			RookAttacks[sq][index] = attacks
		}

	}
}

// Initalize the rook tables holding Magic structs and rook attacks
func initBishopMagics() {
	var sq uint8
	for sq = 0; sq < 64; sq++ {
		magic := &BishopMagics[sq]

		magic.MagicNo = BishopMagicNumbers[sq]
		magic.Mask = GenBishopMasks(sq)
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

		for idx := 0; idx < (1 << no_bits); idx++ {
			index := (permutations[idx] * Bitboard(magic.MagicNo)) >> magic.Shift
			attacks := GenBishopAttacks(sq, permutations[idx])
			BishopAttacks[sq][index] = attacks
		}

	}
}

func InitTables() {
	// Initalize the rook and bishop magic bitboard tables.
	initRookMagics()
	initBishopMagics()
}
