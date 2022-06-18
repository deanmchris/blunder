package engine

const (
	// Constants which map a piece to how much weight it should have on the phase of the game.
	PawnPhase   int16 = 0
	KnightPhase int16 = 1
	BishopPhase int16 = 1
	RookPhase   int16 = 2
	QueenPhase  int16 = 4
	TotalPhase  int16 = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

	// Constants representing a draw or infinite (checkmate) value.
	Inf int16 = 10000
)

type Eval struct {
	MGScores [2]int16
	EGScores [2]int16

	KingZones        [2]KingZone
	KingAttackPoints [2]uint16
	KingAttackers    [2]uint8
}

type KingZone struct {
	OuterRing Bitboard
	InnerRing Bitboard
}

var KingZones [64]KingZone
var DoubledPawnMasks [2][64]Bitboard
var PassedPawnMasks [2][64]Bitboard

var BishopPairBonusMG int16 = 26
var BishopPairBonusEG int16 = 45

var PieceValueMG = [6]int16{96, 351, 361, 476, 975}
var PieceValueEG = [6]int16{121, 255, 279, 500, 930}

var PieceMobilityMG = [5]int16{0, 0, 4, 4, 0}
var PieceMobilityEG = [5]int16{0, 0, 3, 2, 6}

var OuterRingAttackPoints = [5]int16{0, 1, 0, 1, 1}
var InnerRingAttackPoints = [5]int16{0, 3, 4, 5, 3}

var PhaseValues = [6]int16{
	PawnPhase,
	KnightPhase,
	BishopPhase,
	RookPhase,
	QueenPhase,
}

var PSQT_MG = [6][64]int16{
	{
		// MG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		60, 76, 43, 54, 40, 46, 22, 12,
		-33, -13, 5, 10, 44, 43, 6, -39,
		-34, 2, -8, 11, 12, 3, 11, -37,
		-45, -14, -10, 4, 10, 2, 0, -40,
		-37, -16, -13, -13, 1, 0, 22, -24,
		-44, -8, -30, -20, -13, 18, 31, -31,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-70, -16, -10, -12, 3, -32, -7, -35,
		-49, -32, 45, 11, 4, 32, -3, -16,
		-29, 39, 23, 43, 56, 65, 47, 12,
		-12, 13, 5, 40, 18, 54, 6, 11,
		-15, 0, 13, 6, 22, 11, 10, -12,
		-22, -9, 11, 12, 23, 19, 19, -15,
		-28, -33, -11, 6, 7, 18, -13, -18,
		-36, -13, -40, -27, -8, -21, -10, -19,
	},
	{
		// MG Bishop PST
		-13, -4, -25, -12, -7, -12, -1, -1,
		-20, 2, -23, -16, 6, 22, 1, -46,
		-21, 17, 20, 14, 9, 23, 12, -8,
		-8, -3, 1, 27, 15, 14, -2, -10,
		-5, 4, -3, 13, 17, -3, 0, 2,
		-1, 11, 13, 1, 9, 27, 14, 6,
		1, 23, 10, 3, 11, 20, 39, 2,
		-23, 0, 2, -9, -1, 0, -17, -17,
	},
	{
		// MG Rook PST
		11, 17, 4, 22, 19, -1, 2, 5,
		16, 16, 41, 37, 40, 36, 6, 13,
		-14, 6, 10, 11, -2, 14, 25, 0,
		-27, -15, -2, 8, 3, 14, -6, -21,
		-40, -22, -10, -6, -1, -11, 0, -28,
		-36, -20, -8, -10, 0, -1, -7, -28,
		-30, -9, -9, 1, 6, 10, -3, -53,
		-3, 0, 14, 22, 21, 15, -21, -4,
	},
	{
		// MG Queen PST
		-19, -1, 4, 0, 25, 13, 12, 20,
		-18, -42, -2, 2, -4, 22, 13, 24,
		-10, -14, -1, -5, 14, 30, 23, 33,
		-22, -25, -18, -23, -10, 0, -3, -11,
		-4, -27, -7, -12, -8, -5, -3, -3,
		-14, 8, -5, 4, 1, 4, 10, 2,
		-22, -1, 18, 15, 23, 19, -3, -1,
		0, 1, 11, 25, 3, -14, -14, -33,
	},
	{
		// MG King PST
		-3, 0, 1, 0, -2, 0, 0, 0,
		2, 5, 3, 9, 4, 6, 0, -2,
		4, 10, 13, 4, 6, 18, 19, 0,
		0, 4, 6, 0, 0, 0, 4, -11,
		-10, 7, -5, -22, -29, -21, -19, -30,
		0, 2, -8, -31, -33, -28, 6, -16,
		6, 22, 0, -50, -29, -6, 32, 27,
		-13, 49, 28, -54, 9, -18, 41, 28,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		103, 101, 85, 68, 78, 74, 92, 109,
		25, 38, 26, 9, -2, -9, 24, 19,
		-17, -25, -35, -47, -44, -39, -32, -29,
		-23, -26, -41, -45, -45, -46, -38, -41,
		-33, -30, -43, -38, -37, -43, -44, -49,
		-24, -33, -27, -33, -30, -40, -43, -49,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-53, -29, -3, -20, -7, -25, -33, -50,
		-18, 2, -6, 14, 5, -6, -7, -32,
		-16, -2, 23, 21, 11, 14, -2, -20,
		-1, 17, 36, 33, 36, 21, 19, -2,
		-2, 8, 28, 39, 28, 30, 16, -3,
		-8, 11, 9, 25, 20, 5, -6, -9,
		-21, -8, 6, 5, 7, -5, -8, -24,
		-25, -33, -8, 1, -4, -4, -34, -28,
	},
	{
		// EG Bishop PST
		-10, -10, -12, -7, 0, -6, -3, -13,
		0, 0, 6, -7, 0, -1, 0, -5,
		9, -1, -1, -2, -2, 2, 4, 8,
		3, 9, 9, 3, 5, 4, 2, 7,
		1, 3, 9, 10, 0, 5, 0, 0,
		-1, 1, 5, 6, 10, -3, 0, -5,
		-2, -13, -2, 0, 0, -3, -13, -18,
		-9, 2, -7, 5, 3, -4, 0, -5,
	},
	{
		// EG Rook PST
		7, 3, 10, 6, 8, 9, 7, 4,
		6, 7, 2, 4, -5, 0, 7, 4,
		6, 4, 1, 2, 0, -2, -2, -2,
		6, 1, 9, -2, 0, 0, -3, 4,
		10, 5, 6, 1, -3, -3, -7, -4,
		3, 2, -4, -1, -7, -9, -7, -9,
		0, -4, 0, 1, -7, -8, -9, 4,
		-6, -1, -3, -6, -9, -9, 0, -22,
	},
	{
		// EG Queen PST
		-19, 4, 7, 6, 16, 10, 5, 10,
		-17, 3, 6, 15, 17, 12, 7, 1,
		-19, -10, -10, 17, 19, 14, 7, 4,
		1, 5, -1, 12, 21, 11, 26, 21,
		-20, 11, 0, 22, 8, 9, 18, 8,
		-6, -36, 0, -10, -2, 5, 4, 2,
		-11, -18, -34, -23, -23, -16, -19, -9,
		-17, -26, -21, -38, -3, -15, -14, -25,
	},
	{
		// EG King PST
		-24, -16, -10, -12, -9, 5, 2, -5,
		-6, 11, 7, 10, 10, 30, 16, 4,
		1, 13, 16, 10, 13, 38, 37, 5,
		-14, 14, 19, 23, 21, 28, 21, -3,
		-24, -7, 17, 24, 26, 21, 8, -13,
		-22, -5, 10, 21, 24, 18, 2, -9,
		-29, -15, 4, 15, 14, 5, -13, -27,
		-55, -45, -26, -6, -27, -11, -35, -59,
	},
}

var PassedPawnPSQT_MG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	54, 57, 37, 39, 41, 43, 56, 55,
	36, 21, 19, 14, 12, 19, 15, 24,
	15, 9, 12, 4, 6, 9, 6, 6,
	7, -6, -4, -11, -1, 0, 1, 8,
	-2, -3, -5, -6, 0, 3, 1, 6,
	-3, 0, 2, -2, 0, 4, 3, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPSQT_EG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	59, 61, 56, 55, 61, 60, 69, 66,
	74, 55, 43, 34, 32, 49, 47, 62,
	52, 45, 39, 33, 29, 31, 47, 46,
	22, 19, 18, 16, 14, 15, 26, 25,
	0, 0, 4, 0, 0, 0, 7, 7,
	-6, -2, -3, 0, 3, 0, 2, 2,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// Flip white's perspective to black
var FlipSq [2][64]int = [2][64]int{
	{
		0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
	},

	{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	},
}

// Evaluate a position and give a score, from the perspective of the side to move (
// more positive if it's good for the side to move, otherwise more negative).
func EvaluatePos(pos *Position) int16 {
	eval := Eval{
		MGScores: pos.MGScores,
		EGScores: pos.EGScores,
		KingZones: [2]KingZone{
			KingZones[pos.Pieces[Black][King].Msb()],
			KingZones[pos.Pieces[White][King].Msb()],
		},
	}

	phase := pos.Phase
	allBB := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]

	for allBB != 0 {
		sq := allBB.PopBit()
		piece := pos.Squares[sq]

		switch piece.Type {
		case Pawn:
			evalPawn(pos, piece.Color, sq, &eval)
		case Knight:
			evalKnight(pos, piece.Color, sq, &eval)
		case Bishop:
			evalBishop(pos, piece.Color, sq, &eval)
		case Rook:
			evalRook(pos, piece.Color, sq, &eval)
		case Queen:
			evalQueen(pos, piece.Color, sq, &eval)
		}
	}

	if pos.Pieces[White][Bishop].CountBits() >= 2 {
		eval.MGScores[White] += BishopPairBonusMG
		eval.EGScores[White] += BishopPairBonusEG
	}

	if pos.Pieces[Black][Bishop].CountBits() >= 2 {
		eval.MGScores[Black] += BishopPairBonusMG
		eval.EGScores[Black] += BishopPairBonusEG
	}

	evalKing(pos, White, pos.Pieces[White][King].Msb(), &eval)
	evalKing(pos, Black, pos.Pieces[Black][King].Msb(), &eval)

	mgScore := eval.MGScores[pos.SideToMove] - eval.MGScores[pos.SideToMove^1]
	egScore := eval.EGScores[pos.SideToMove] - eval.EGScores[pos.SideToMove^1]

	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	return int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))
}

// Evaluate the score of a pawn.
func evalPawn(pos *Position, color, sq uint8, eval *Eval) {
	enemyPawns := pos.Pieces[color^1][Pawn]
	usPawns := pos.Pieces[color][Pawn]

	// Evaluate passed pawns, but make sure they're not behind a friendly pawn.
	if PassedPawnMasks[color][sq]&enemyPawns == 0 && usPawns&DoubledPawnMasks[color][sq] == 0 {
		eval.MGScores[color] += PassedPawnPSQT_MG[FlipSq[color][sq]]
		eval.EGScores[color] += PassedPawnPSQT_EG[FlipSq[color][sq]]
	}
}

// Evaluate the score of a knight.
func evalKnight(pos *Position, color, sq uint8, eval *Eval) {
	usBB := pos.Sides[color]
	moves := KnightMoves[sq] & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 4) * PieceMobilityMG[Knight]
	eval.EGScores[color] += (mobility - 4) * PieceMobilityEG[Knight]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(OuterRingAttackPoints[Knight])
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(InnerRingAttackPoints[Knight])
	}
}

// Evaluate the score of a bishop.
func evalBishop(pos *Position, color, sq uint8, eval *Eval) {
	usBB := pos.Sides[color]
	allBB := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]

	moves := GenBishopMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 7) * PieceMobilityMG[Bishop]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Bishop]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(OuterRingAttackPoints[Bishop])
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(InnerRingAttackPoints[Bishop])
	}
}

// Evaluate the score of a rook.
func evalRook(pos *Position, color, sq uint8, eval *Eval) {
	usBB := pos.Sides[color]
	allBB := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]

	moves := GenRookMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 7) * PieceMobilityMG[Rook]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Rook]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(OuterRingAttackPoints[Rook])
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(InnerRingAttackPoints[Rook])
	}
}

// Evaluate the score of a queen.
func evalQueen(pos *Position, color, sq uint8, eval *Eval) {
	usBB := pos.Sides[color]
	allBB := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]

	moves := (GenBishopMoves(sq, allBB) | GenRookMoves(sq, allBB)) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 14) * PieceMobilityMG[Queen]
	eval.EGScores[color] += (mobility - 14) * PieceMobilityEG[Queen]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(OuterRingAttackPoints[Queen])
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(InnerRingAttackPoints[Queen])
	}
}

// Evaluate the score of a king.
func evalKing(pos *Position, color, sq uint8, eval *Eval) {
	// Take all the king saftey points collected for the enemy,
	// and see what kind of penatly we should get.
	enemyPoints := eval.KingAttackPoints[color^1]
	penatly := int16((enemyPoints * enemyPoints) / 4)
	if eval.KingAttackers[color^1] >= 2 && pos.Pieces[color^1][Queen] != 0 {
		eval.MGScores[color] -= penatly
	}
}

func InitEvalBitboards() {
	for sq := 0; sq < 64; sq++ {
		// Create king zones.
		sqBB := SquareBB[sq]
		zone := ((sqBB & ClearFile[FileH]) >> 1) | ((sqBB & (ClearFile[FileG] & ClearFile[FileH])) >> 2)
		zone |= ((sqBB & ClearFile[FileA]) << 1) | ((sqBB & (ClearFile[FileB] & ClearFile[FileA])) << 2)
		zone |= sqBB

		zone |= ((zone >> 8) | (zone >> 16))
		zone |= ((zone << 8) | (zone << 16))

		KingZones[sq] = KingZone{OuterRing: zone & ^(KingMoves[sq] | sqBB), InnerRing: KingMoves[sq] | sqBB}

		file := FileOf(uint8(sq))
		fileBB := MaskFile[file]
		rank := int(RankOf(uint8(sq)))

		// Create doubled pawns masks.
		mask := fileBB
		for r := 0; r <= rank; r++ {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[White][sq] = mask

		mask = fileBB
		for r := 7; r >= rank; r-- {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[Black][sq] = mask

		// Passed pawn masks.
		frontSpanMask := fileBB
		frontSpanMask |= (fileBB & ClearFile[FileA]) << 1
		frontSpanMask |= (fileBB & ClearFile[FileH]) >> 1

		whiteFrontSpan := frontSpanMask
		for r := 0; r <= rank; r++ {
			whiteFrontSpan &= ClearRank[r]
		}
		PassedPawnMasks[White][sq] = whiteFrontSpan

		blackFrontSpan := frontSpanMask
		for r := 7; r >= rank; r-- {
			blackFrontSpan &= ClearRank[r]
		}
		PassedPawnMasks[Black][sq] = blackFrontSpan
	}
}
