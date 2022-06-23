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
var PieceValueEG = [6]int16{119, 254, 279, 499, 932}

var PieceMobilityMG = [5]int16{0, 1, 4, 5, 0}
var PieceMobilityEG = [5]int16{0, 0, 3, 2, 6}

var OuterRingAttackPoints = [5]int16{0, 1, 0, 1, 0}
var InnerRingAttackPoints = [5]int16{0, 3, 4, 4, 3}
var SemiOpenFileNextToKingPenalty int16 = 3

var IsolatedPawnPenatlyMG int16 = 15
var IsolatedPawnPenatlyEG int16 = 7
var DoubledPawnPenatlyMG int16 = 4
var DoubledPawnPenatlyEG int16 = 14

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
		57, 74, 41, 53, 38, 46, 17, 7,
		-30, -14, 5, 10, 44, 46, 6, -34,
		-28, -2, -3, 16, 15, 6, 4, -30,
		-40, -23, -7, 7, 11, 2, -10, -35,
		-31, -25, -9, -10, 1, 2, 10, -18,
		-38, -18, -28, -18, -13, 18, 18, -25,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-73, -16, -10, -12, 3, -33, -7, -36,
		-52, -33, 45, 11, 2, 32, -3, -16,
		-29, 39, 21, 42, 56, 66, 47, 12,
		-12, 13, 4, 39, 18, 52, 6, 9,
		-13, 0, 13, 6, 21, 11, 10, -12,
		-22, -8, 11, 13, 24, 19, 20, -14,
		-26, -33, -9, 7, 9, 18, -11, -16,
		-37, -12, -39, -25, -6, -19, -9, -19,
	},
	{
		// MG Bishop PST
		-13, -4, -26, -12, -7, -12, -1, -1,
		-21, 0, -24, -17, 4, 21, 0, -48,
		-21, 17, 19, 12, 7, 23, 10, -9,
		-8, -2, 0, 25, 13, 12, -1, -10,
		-4, 4, -1, 13, 17, -2, 0, 2,
		-1, 12, 14, 3, 10, 27, 14, 6,
		1, 23, 11, 5, 13, 20, 39, 2,
		-23, 0, 3, -6, 0, 1, -17, -16,
	},
	{
		// MG Rook PST
		10, 17, 2, 22, 18, -1, 0, 3,
		15, 14, 41, 37, 40, 36, 4, 13,
		-14, 4, 8, 9, -3, 12, 25, -1,
		-27, -15, -2, 6, 0, 12, -7, -21,
		-39, -21, -10, -6, -1, -11, 0, -28,
		-34, -18, -6, -8, 0, 0, -7, -28,
		-28, -7, -7, 2, 7, 10, -1, -51,
		-2, 1, 15, 23, 22, 16, -18, -3,
	},
	{
		// MG Queen PST
		-18, -1, 2, 0, 25, 13, 12, 20,
		-18, -44, -2, 1, -5, 20, 11, 24,
		-10, -14, -1, -7, 12, 29, 21, 31,
		-22, -24, -18, -23, -10, -1, -4, -11,
		-4, -26, -7, -11, -7, -4, -3, -3,
		-13, 8, -3, 5, 1, 4, 10, 2,
		-20, 0, 18, 17, 24, 20, -1, 0,
		1, 2, 12, 25, 4, -12, -12, -33,
	},
	{
		// MG King PST
		-3, 0, 1, 0, -2, 0, 0, 0,
		2, 5, 3, 9, 4, 6, 0, -2,
		4, 10, 13, 4, 6, 18, 19, 0,
		0, 4, 6, 0, 0, 0, 4, -11,
		-10, 7, -5, -22, -30, -21, -19, -31,
		0, 2, -8, -29, -31, -26, 7, -16,
		6, 23, 0, -48, -28, -5, 31, 24,
		-13, 47, 26, -54, 8, -18, 40, 24,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		102, 97, 83, 64, 74, 70, 88, 106,
		19, 31, 25, 10, 0, -10, 20, 16,
		-16, -25, -32, -44, -39, -34, -33, -27,
		-23, -26, -36, -41, -41, -42, -40, -38,
		-33, -32, -39, -33, -33, -38, -48, -46,
		-24, -34, -21, -28, -25, -36, -46, -46,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-54, -29, -3, -20, -7, -25, -34, -52,
		-17, 2, -7, 14, 5, -6, -7, -32,
		-15, -2, 23, 20, 9, 12, -2, -20,
		-1, 17, 36, 33, 36, 21, 19, -2,
		0, 8, 28, 39, 28, 29, 16, -1,
		-7, 11, 9, 25, 20, 5, -6, -8,
		-19, -6, 6, 5, 7, -4, -7, -24,
		-25, -32, -6, 1, -4, -2, -34, -29,
	},
	{
		// EG Bishop PST
		-9, -10, -11, -7, 0, -6, -3, -13,
		0, 0, 6, -7, 0, -1, 0, -3,
		8, -1, -1, -2, -2, 0, 4, 7,
		3, 8, 9, 3, 5, 3, 0, 6,
		1, 3, 9, 10, 0, 5, 0, 0,
		0, 1, 5, 7, 10, -3, 0, -5,
		-1, -12, -2, 0, 0, -3, -13, -16,
		-7, 2, -7, 4, 1, -4, 0, -3,
	},
	{
		// EG Rook PST
		7, 3, 10, 6, 8, 9, 7, 4,
		6, 7, 2, 3, -5, 0, 7, 4,
		6, 4, 1, 2, 0, -2, -2, -2,
		7, 1, 9, -2, 0, 0, -3, 4,
		10, 5, 6, 1, -3, -3, -7, -2,
		3, 2, -4, -1, -7, -9, -5, -8,
		0, -3, 0, 0, -7, -8, -9, 3,
		-7, -2, -4, -7, -10, -10, 0, -23,
	},
	{
		// EG Queen PST
		-18, 3, 6, 5, 14, 10, 5, 10,
		-15, 3, 6, 15, 16, 10, 5, 0,
		-19, -10, -11, 17, 18, 12, 5, 2,
		2, 5, 0, 12, 20, 10, 26, 22,
		-18, 11, 0, 20, 6, 7, 18, 8,
		-4, -36, 0, -10, -2, 4, 4, 2,
		-9, -18, -33, -23, -23, -16, -19, -8,
		-17, -27, -21, -36, -3, -13, -12, -25,
	},
	{
		// EG King PST
		-25, -17, -10, -12, -9, 5, 0, -5,
		-6, 10, 7, 10, 10, 30, 16, 4,
		0, 13, 15, 10, 13, 38, 36, 3,
		-14, 13, 19, 23, 21, 28, 21, -3,
		-23, -6, 17, 24, 27, 22, 8, -12,
		-22, -5, 10, 21, 24, 18, 2, -9,
		-29, -15, 3, 14, 14, 5, -12, -26,
		-55, -45, -26, -6, -26, -10, -36, -59,
	},
}

var PassedPawnPSQT_MG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	51, 55, 35, 38, 39, 43, 51, 50,
	43, 25, 22, 16, 12, 21, 15, 24,
	19, 10, 12, 3, 7, 12, 6, 5,
	12, -6, -7, -12, -2, 0, 1, 11,
	2, 0, -7, -9, 0, 4, 2, 8,
	0, 3, 2, -3, 0, 4, 5, 2,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPSQT_EG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	58, 57, 54, 51, 57, 56, 65, 63,
	93, 71, 53, 40, 35, 59, 58, 76,
	60, 54, 45, 38, 30, 34, 56, 53,
	28, 26, 20, 18, 15, 18, 32, 29,
	4, 5, 5, 1, 0, 0, 13, 9,
	-1, 3, -4, 0, 4, 0, 6, 5,
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
	enemyPoints := eval.KingAttackPoints[color^1]

	// Evaluate semi-open files adjacent to the enemy king
	kingFile := MaskFile[FileOf(sq)]
	usPawns := pos.Pieces[color][Pawn]

	leftFile := ((kingFile & ClearFile[FileA]) << 1)
	rightFile := ((kingFile & ClearFile[FileH]) >> 1)

	if kingFile&usPawns == 0 {
		enemyPoints += uint16(SemiOpenFileNextToKingPenalty)
	}

	if leftFile != 0 && leftFile&usPawns == 0 {
		enemyPoints += uint16(SemiOpenFileNextToKingPenalty)
	}

	if rightFile != 0 && rightFile&usPawns == 0 {
		enemyPoints += uint16(SemiOpenFileNextToKingPenalty)
	}

	// Take all the king saftey points collected for the enemy,
	// and see what kind of penatly we should get.
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
