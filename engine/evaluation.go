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
var IsolatedPawnMasks [8]Bitboard
var PassedPawnMasks [2][64]Bitboard

var BishopPairBonusMG int16 = 28
var BishopPairBonusEG int16 = 45

var PieceValueMG = [6]int16{96, 350, 364, 477, 976}
var PieceValueEG = [6]int16{122, 255, 279, 499, 931}

var PieceMobilityMG = [5]int16{0, 1, 4, 5, 0}
var PieceMobilityEG = [5]int16{0, 0, 3, 2, 6}

var OuterRingAttackPoints = [5]int16{0, 1, 0, 1, 1}
var InnerRingAttackPoints = [5]int16{0, 3, 4, 5, 3}

var IsolatedPawnPenatlyMG int16 = 7
var IsolatedPawnPenatlyEG int16 = 16
var DoubledPawnPenatlyMG int16 = 13
var DoubledPawnPenatlyEG int16 = 7

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
		59, 75, 42, 53, 39, 46, 20, 10,
		-32, -13, 5, 10, 44, 44, 6, -37,
		-31, 0, -5, 15, 14, 5, 8, -33,
		-42, -20, -8, 6, 10, 3, -5, -37,
		-33, -23, -11, -11, 1, 2, 14, -19,
		-40, -16, -29, -20, -14, 19, 22, -26,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-71, -16, -10, -12, 3, -32, -7, -35,
		-50, -32, 45, 11, 3, 32, -3, -16,
		-29, 39, 22, 42, 56, 65, 47, 12,
		-12, 13, 4, 40, 17, 53, 6, 10,
		-14, 0, 13, 6, 22, 11, 10, -12,
		-22, -8, 11, 12, 23, 19, 19, -15,
		-27, -33, -10, 6, 8, 18, -12, -17,
		-36, -13, -40, -26, -7, -20, -10, -19,
	},
	{
		// MG Bishop PST
		-13, -4, -25, -12, -7, -12, -1, -1,
		-20, 1, -23, -16, 5, 22, 0, -46,
		-21, 17, 19, 13, 8, 23, 11, -8,
		-8, -2, 0, 26, 14, 13, -1, -10,
		-5, 4, -2, 13, 17, -2, 0, 2,
		-1, 11, 13, 2, 9, 27, 14, 6,
		1, 23, 10, 4, 12, 20, 39, 2,
		-23, 0, 2, -8, 0, 0, -17, -16,
	},
	{
		// MG Rook PST
		10, 17, 3, 22, 18, -1, 1, 4,
		15, 15, 41, 37, 40, 36, 5, 13,
		-14, 5, 9, 10, -2, 13, 25, 0,
		-27, -15, -2, 7, 2, 13, -6, -21,
		-40, -21, -10, -6, -1, -11, 0, -28,
		-35, -19, -7, -9, 0, -1, -7, -28,
		-29, -8, -8, 1, 6, 10, -2, -52,
		-3, 1, 14, 22, 21, 15, -19, -4,
	},
	{
		// MG Queen PST
		-19, -1, 3, 0, 25, 13, 12, 20,
		-18, -43, -2, 1, -4, 21, 12, 24,
		-10, -14, -1, -6, 13, 29, 22, 32,
		-22, -24, -18, -23, -10, 0, -3, -11,
		-4, -26, -7, -11, -7, -5, -3, -3,
		-13, 8, -4, 4, 1, 4, 10, 2,
		-21, 0, 18, 16, 23, 19, -2, 0,
		0, 1, 11, 25, 3, -13, -13, -33,
	},
	{
		// MG King PST
		-3, 0, 1, 0, -2, 0, 0, 0,
		2, 5, 3, 9, 4, 6, 0, -2,
		4, 10, 13, 4, 6, 18, 19, 0,
		0, 4, 6, 0, 0, 0, 4, -11,
		-10, 7, -5, -22, -29, -21, -19, -30,
		0, 2, -8, -30, -32, -27, 6, -16,
		6, 22, 0, -49, -28, -6, 32, 26,
		-13, 48, 27, -54, 8, -18, 40, 27,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		102, 99, 83, 66, 76, 72, 90, 107,
		24, 36, 26, 9, -1, -9, 23, 19,
		-15, -25, -33, -45, -41, -36, -33, -27,
		-22, -27, -38, -43, -43, -43, -41, -38,
		-32, -33, -40, -35, -36, -40, -49, -47,
		-23, -35, -23, -31, -28, -38, -48, -46,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-53, -29, -3, -20, -7, -25, -33, -51,
		-18, 2, -6, 14, 5, -6, -7, -32,
		-15, -2, 23, 20, 10, 13, -2, -20,
		-1, 17, 36, 33, 36, 21, 19, -2,
		-1, 8, 28, 39, 28, 29, 16, -2,
		-8, 11, 9, 25, 20, 5, -6, -9,
		-20, -7, 6, 5, 7, -4, -7, -24,
		-25, -32, -7, 1, -4, -3, -34, -28,
	},
	{
		// EG Bishop PST
		-9, -10, -11, -7, 0, -6, -3, -13,
		0, 0, 6, -7, 0, -1, 0, -4,
		8, -1, -1, -2, -2, 1, 4, 7,
		3, 9, 9, 3, 5, 3, 1, 6,
		1, 3, 9, 10, 0, 5, 0, 0,
		0, 1, 5, 7, 10, -3, 0, -5,
		-1, -13, -2, 0, 0, -3, -13, -17,
		-8, 2, -7, 4, 2, -4, 0, -4,
	},
	{
		// EG Rook PST
		7, 3, 10, 6, 8, 9, 7, 4,
		6, 7, 2, 3, -5, 0, 7, 4,
		6, 4, 1, 2, 0, -2, -2, -2,
		6, 1, 9, -2, 0, 0, -3, 4,
		10, 5, 6, 1, -3, -3, -7, -3,
		3, 2, -4, -1, -7, -9, -6, -8,
		0, -3, 0, 0, -7, -8, -9, 3,
		-7, -1, -3, -7, -9, -9, 0, -22,
	},
	{
		// EG Queen PST
		-19, 4, 6, 5, 15, 10, 5, 10,
		-16, 3, 6, 15, 16, 11, 6, 0,
		-19, -10, -10, 17, 18, 13, 6, 3,
		1, 5, 0, 12, 20, 10, 26, 21,
		-19, 11, 0, 21, 7, 8, 18, 8,
		-5, -36, 0, -10, -2, 4, 4, 2,
		-10, -18, -33, -23, -23, -16, -19, -8,
		-17, -26, -21, -37, -3, -14, -13, -25,
	},
	{
		// EG King PST
		-24, -16, -10, -12, -9, 5, 1, -5,
		-6, 10, 7, 10, 10, 30, 16, 4,
		0, 13, 15, 10, 13, 38, 36, 4,
		-14, 14, 19, 23, 21, 28, 21, -3,
		-23, -6, 17, 24, 26, 21, 8, -13,
		-22, -5, 10, 21, 24, 18, 2, -9,
		-29, -15, 3, 15, 14, 5, -12, -27,
		-55, -45, -26, -6, -26, -11, -36, -59,
	},
}

var PassedPawnPSQT_MG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	53, 56, 36, 38, 40, 43, 54, 53,
	39, 22, 20, 15, 12, 20, 15, 24,
	17, 9, 12, 4, 6, 10, 6, 6,
	9, -6, -5, -11, -1, 0, 1, 9,
	0, -2, -6, -7, 0, 3, 1, 7,
	-1, 1, 2, -2, 0, 4, 4, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPSQT_EG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	58, 59, 54, 53, 59, 58, 67, 64,
	82, 61, 47, 36, 33, 53, 51, 68,
	56, 48, 42, 35, 29, 32, 50, 49,
	25, 22, 19, 17, 14, 16, 28, 27,
	1, 1, 4, 0, 0, 0, 9, 7,
	-4, 0, -3, 0, 3, 0, 3, 3,
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

	// Evaluate isolated pawns.
	if IsolatedPawnMasks[FileOf(sq)]&usPawns == 0 {
		eval.MGScores[color] -= IsolatedPawnPenatlyMG
		eval.EGScores[color] -= IsolatedPawnPenatlyEG
	}

	// Evaluate doubled pawns.
	if DoubledPawnMasks[color][sq]&usPawns != 0 {
		eval.MGScores[color] -= DoubledPawnPenatlyMG
		eval.EGScores[color] -= DoubledPawnPenatlyEG
	}

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
	for file := FileA; file <= FileH; file++ {
		fileBB := MaskFile[file]
		mask := (fileBB & ClearFile[FileA]) << 1
		mask |= (fileBB & ClearFile[FileH]) >> 1
		IsolatedPawnMasks[file] = mask
	}

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
