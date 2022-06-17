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

var BishopPairBonusMG int16 = 26
var BishopPairBonusEG int16 = 44

var PieceValueMG = [6]int16{97, 350, 359, 474, 966}
var PieceValueEG = [6]int16{142, 250, 273, 490, 920}

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
		71, 84, 51, 60, 44, 48, 31, 22,
		-26, -8, 10, 16, 46, 42, 7, -35,
		-35, -1, -10, 8, 9, 0, 6, -40,
		-47, -17, -13, 0, 7, 0, -2, -42,
		-39, -18, -16, -16, -1, -2, 19, -26,
		-46, -10, -33, -24, -17, 15, 28, -33,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-64, -15, -9, -12, 3, -29, -7, -32,
		-45, -30, 44, 11, 4, 31, -3, -15,
		-27, 38, 24, 44, 55, 62, 47, 11,
		-12, 13, 7, 41, 19, 56, 7, 12,
		-16, 0, 12, 6, 22, 12, 11, -13,
		-23, -11, 10, 11, 21, 19, 19, -16,
		-29, -32, -13, 5, 6, 17, -14, -19,
		-33, -14, -41, -28, -10, -23, -11, -18,
	},
	{
		// MG Bishop PST
		-13, -4, -23, -11, -6, -11, -1, -2,
		-19, 2, -21, -14, 7, 22, 2, -44,
		-21, 17, 21, 16, 10, 24, 14, -7,
		-8, -3, 3, 29, 17, 17, -2, -8,
		-6, 5, -3, 13, 17, -3, 0, 1,
		-2, 10, 12, 1, 8, 26, 13, 5,
		0, 21, 9, 2, 10, 18, 38, 1,
		-25, -1, 1, -12, -4, 0, -17, -19,
	},
	{
		// MG Rook PST
		11, 17, 6, 23, 20, 0, 3, 6,
		17, 18, 41, 36, 39, 35, 8, 14,
		-12, 7, 11, 13, 0, 15, 24, 2,
		-26, -14, -2, 9, 6, 15, -4, -18,
		-41, -22, -10, -6, -1, -10, 1, -27,
		-39, -21, -10, -12, 0, -1, -6, -29,
		-33, -11, -12, -1, 4, 9, -3, -55,
		-6, -3, 11, 19, 19, 13, -21, -5,
	},
	{
		// MG Queen PST
		-20, -1, 6, 2, 26, 13, 12, 20,
		-20, -41, -1, 4, -1, 25, 16, 25,
		-11, -14, -1, -1, 17, 32, 26, 36,
		-23, -25, -17, -21, -7, 5, 0, -9,
		-5, -27, -9, -13, -8, -5, -2, -3,
		-16, 6, -7, 2, -1, 3, 10, 2,
		-26, -4, 16, 13, 21, 16, -6, -4,
		-3, -3, 7, 23, 0, -19, -16, -33,
	},
	{
		// MG King PST
		-3, 0, 1, 0, -2, 0, 1, 0,
		2, 5, 3, 8, 4, 6, 1, -2,
		4, 10, 12, 4, 6, 17, 18, 0,
		0, 4, 6, 0, 0, 0, 4, -10,
		-10, 6, -5, -20, -26, -20, -18, -28,
		0, 1, -8, -32, -34, -30, 3, -16,
		4, 20, -1, -52, -31, -7, 32, 29,
		-13, 51, 29, -53, 9, -17, 43, 31,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		129, 125, 104, 88, 92, 89, 108, 128,
		51, 55, 38, 19, 5, 1, 35, 38,
		-11, -24, -34, -45, -53, -46, -34, -29,
		-29, -38, -52, -57, -59, -60, -48, -47,
		-41, -40, -55, -48, -52, -57, -54, -57,
		-31, -40, -38, -40, -40, -52, -52, -57,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-51, -28, -3, -20, -7, -25, -31, -46,
		-20, 1, -5, 13, 5, -6, -8, -32,
		-17, -2, 23, 21, 12, 16, -2, -20,
		-1, 17, 36, 32, 35, 21, 18, -3,
		-3, 7, 28, 39, 28, 29, 16, -4,
		-9, 11, 10, 25, 21, 5, -6, -9,
		-22, -10, 5, 5, 8, -6, -9, -24,
		-25, -34, -10, 0, -5, -5, -34, -27,
	},
	{
		// EG Bishop PST
		-11, -11, -13, -7, -1, -7, -4, -14,
		-2, 0, 7, -8, 0, -1, 0, -8,
		8, 0, 0, -1, -1, 4, 4, 7,
		3, 11, 9, 3, 5, 5, 3, 7,
		0, 3, 10, 10, 0, 7, 0, 0,
		-3, 0, 6, 8, 11, -1, 0, -6,
		-4, -13, -2, 1, 2, -3, -12, -19,
		-13, 0, -8, 5, 3, -5, -2, -8,
	},
	{
		// EG Rook PST
		7, 4, 10, 6, 8, 9, 7, 4,
		5, 6, 2, 5, -4, 0, 5, 4,
		6, 3, 1, 2, 0, -2, -2, -2,
		5, 1, 9, -3, -1, 0, -4, 4,
		9, 4, 5, 1, -4, -5, -8, -6,
		3, 2, -4, 0, -7, -10, -8, -10,
		0, -4, 1, 2, -6, -8, -9, 4,
		-6, 0, -2, -5, -8, -9, 0, -22,
	},
	{
		// EG Queen PST
		-20, 4, 8, 7, 18, 10, 5, 10,
		-20, 2, 6, 15, 18, 14, 8, 2,
		-20, -10, -9, 18, 20, 16, 8, 5,
		-2, 5, -1, 13, 22, 12, 25, 19,
		-22, 9, 1, 24, 10, 10, 18, 7,
		-8, -36, 0, -8, -2, 6, 3, 1,
		-14, -18, -35, -22, -22, -16, -19, -10,
		-18, -25, -21, -42, -4, -18, -15, -25,
	},
	{
		// EG King PST
		-21, -14, -9, -11, -8, 6, 4, -4,
		-5, 13, 8, 10, 10, 31, 17, 4,
		3, 15, 18, 10, 13, 38, 38, 8,
		-13, 15, 19, 23, 20, 28, 22, -1,
		-24, -7, 17, 23, 26, 20, 8, -13,
		-22, -6, 8, 20, 24, 18, 2, -9,
		-29, -16, 2, 15, 14, 4, -14, -28,
		-53, -48, -29, -7, -28, -12, -37, -61,
	},
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
	}
}
