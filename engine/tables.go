package engine

// tables.go contains various precomputed tables used in the engine.

const (
	Rank1 uint8 = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

const (
	FileA uint8 = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

const (
	North uint8 = 8
	South uint8 = 8
	East  uint8 = 1
	West  uint8 = 1
)

var ClearRank = [8]uint64{}
var ClearFile = [8]uint64{}
var MaskRank = [8]uint64{}
var MaskFile = [8]uint64{}

var KingMoves = [64]uint64{}
var KnightMoves = [64]uint64{}
var PawnAttacks = [2][64]uint64{}
var PawnPushes = [2][64]uint64{}

var RaysBetween = [64][64]uint64{}

var MaskDiagonal = [15]uint64{
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

var MaskAntidiagonal = [15]uint64{
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

func InitTables() {
	// Generate useful masking lookup tables.

	for i := uint8(0); i < 8; i++ {
		emptyBB := EmptyBB
		fullBB := FullBB

		for j := i; j <= 63; j += 8 {
			setBit(&emptyBB, j)
			clearBit(&fullBB, j)
		}

		MaskFile[i] = emptyBB
		ClearFile[i] = fullBB
	}

	for i := uint8(0); i <= 56; i += 8 {
		emptyBB := EmptyBB
		fullBB := FullBB

		for j := i; j < i+8; j++ {
			setBit(&emptyBB, j)
			clearBit(&fullBB, j)
		}

		MaskRank[i/8] = emptyBB
		ClearRank[i/8] = fullBB
	}

	for sq := 0; sq < 64; sq++ {
		sqBB := MostSigBitBB >> sq
		sqBBClippedHFile := sqBB & ClearFile[FileH]
		sqBBClippedAFile := sqBB & ClearFile[FileA]
		sqBBClippedHGFile := sqBB & ClearFile[FileH] & ClearFile[FileG]
		sqBBClippedABFile := sqBB & ClearFile[FileA] & ClearFile[FileB]

		// Generate king moves lookup table.

		top := sqBB >> North
		topRight := sqBBClippedHFile >> North >> East
		topLeft := sqBBClippedAFile >> North << West

		right := sqBBClippedHFile >> East
		left := sqBBClippedAFile << West

		bottom := sqBB << South
		bottomRight := sqBBClippedHFile << South >> East
		bottomLeft := sqBBClippedAFile << South << West

		kingMoves := top | topRight | topLeft | right | left | bottom | bottomRight | bottomLeft
		KingMoves[sq] = kingMoves

		// Generate knight moves lookup table.

		northNorthEast := sqBBClippedHFile >> North >> North >> East
		northEastEast := sqBBClippedHGFile >> North >> East >> East

		southEastEast := sqBBClippedHGFile << South >> East >> East
		southSouthEast := sqBBClippedHFile << South << South >> East

		southSouthWest := sqBBClippedAFile << South << South << West
		southWestWest := sqBBClippedABFile << South << West << West

		northNorthWest := sqBBClippedAFile >> North >> North << West
		northWestWest := sqBBClippedABFile >> North << West << West

		knightMoves := northNorthEast | northEastEast | southEastEast | southSouthEast |
			southSouthWest | southWestWest | northNorthWest | northWestWest
		KnightMoves[sq] = knightMoves

		// Generate pawn pushes lookup table.

		whitePawnPush := sqBB >> North
		blackPawnPush := sqBB << South

		PawnPushes[White][sq] = whitePawnPush
		PawnPushes[Black][sq] = blackPawnPush

		// Generate pawn attacks lookup table.

		whitePawnRightAttack := sqBBClippedHFile >> North >> East
		whitePawnLeftAttack := sqBBClippedAFile >> North << West

		blackPawnRightAttack := sqBBClippedHFile << South >> East
		blackPawnLeftAttack := sqBBClippedAFile << South << West

		PawnAttacks[White][sq] = whitePawnRightAttack | whitePawnLeftAttack
		PawnAttacks[Black][sq] = blackPawnRightAttack | blackPawnLeftAttack
	}

	for i := uint8(0); i < 64; i++ {
		for j := uint8(0); j < 64; j++ {

			smallerSq := i
			biggerSq := j

			if i > j {
				smallerSq = j
				biggerSq = i
			}

			for k := 0; k < 15; k++ {
				diagonal := MaskDiagonal[k]
				antidiagonal := MaskAntidiagonal[k]

				if bitIsSet(diagonal, smallerSq) && bitIsSet(diagonal, biggerSq) {
					RaysBetween[i][j] = diagonal & (FullBB >> smallerSq) & ^(FullBB >> biggerSq)
					RaysBetween[i][j] &= ^(MostSigBitBB >> smallerSq)
				}

				if bitIsSet(antidiagonal, i) && bitIsSet(antidiagonal, j) {
					RaysBetween[i][j] = antidiagonal & (FullBB >> smallerSq) & ^(FullBB >> biggerSq)
					RaysBetween[i][j] &= ^(MostSigBitBB >> smallerSq)
				}
			}

			for k := 0; k < 8; k++ {
				rank := MaskRank[k]
				file := MaskFile[k]

				if bitIsSet(rank, i) && bitIsSet(rank, j) {
					RaysBetween[i][j] = rank & (FullBB >> smallerSq) & ^(FullBB >> biggerSq)
					RaysBetween[i][j] &= ^(MostSigBitBB >> smallerSq)
				}

				if bitIsSet(file, i) && bitIsSet(file, j) {
					RaysBetween[i][j] = file & (FullBB >> smallerSq) & ^(FullBB >> biggerSq)
					RaysBetween[i][j] &= ^(MostSigBitBB >> smallerSq)
				}
			}
		}
	}
}
