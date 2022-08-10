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

	// A constant to scale the evaluation score by if the position is considered
	// drawish (e.g. king and queen vs king and queen).
	ScaleFactor int16 = 16
)

// A constant representing a draw value.
var Draw int16 = 0

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
var OutpostMasks [2][64]Bitboard

var PieceValueMG = [6]int16{84, 333, 346, 441, 921}
var PieceValueEG = [6]int16{106, 244, 268, 478, 886}

var PieceMobilityMG = [5]int16{0, 5, 3, 3, 0}
var PieceMobilityEG = [5]int16{0, 2, 3, 2, 6}

var BishopPairBonusMG int16 = 22
var BishopPairBonusEG int16 = 30

var IsolatedPawnPenatlyMG int16 = 17
var IsolatedPawnPenatlyEG int16 = 6
var DoubledPawnPenatlyMG int16 = 1
var DoubledPawnPenatlyEG int16 = 16

// A middlegame equivelent of this bonus is not missing,
// and one was used at first. But several thousand
// iterations of the tuner indicated that such a "bonus"
// was actually quite bad to give in the middlegame.
var RookOrQueenOnSeventhBonusEG int16 = 23

var KnightOnOutpostBonusMG int16 = 27
var KnightOnOutpostBonusEG int16 = 18

var RookOnOpenFileBonusMG int16 = 23
var TempoBonusMG int16 = 14

var BishopOutPostBonusMG int16 = 10
var BishopOutPostBonusEG int16 = 14

var OuterRingAttackPoints = [5]int16{0, 1, 0, 1, 1}
var InnerRingAttackPoints = [5]int16{0, 3, 4, 3, 2}
var SemiOpenFileNextToKingPenalty int16 = 4

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
		45, 52, 42, 43, 28, 34, 19, 9,
		-14, -3, 7, 14, 35, 50, 15, -6,
		-27, -6, -8, 13, 16, 4, -3, -25,
		-32, -28, -7, 5, 7, -1, -15, -30,
		-29, -25, -12, -12, -1, -5, 6, -17,
		-34, -23, -27, -18, -14, 10, 13, -22,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-43, -11, -8, -5, 1, -20, -4, -22,
		-31, -22, 19, 7, 5, 13, -8, -11,
		-21, 21, 8, 16, 36, 33, 19, 6,
		-6, 2, 0, 23, 8, 27, 4, 14,
		-3, 10, 12, 8, 16, 10, 19, 1,
		-19, -4, 3, 7, 22, 12, 15, -11,
		-21, -20, -9, 8, 9, 11, -5, 0,
		-19, -13, -20, -14, -2, 3, -11, -8,
	},
	{
		// MG Bishop PST
		-13, 0, -17, -8, -7, -5, -2, -3,
		-21, 0, -16, -10, 4, 1, -6, -41,
		-23, 6, 10, 8, 8, 26, 0, -10,
		-15, -4, 2, 22, 9, 10, -1, -16,
		0, 10, -2, 15, 17, -7, -1, 13,
		-2, 16, 13, 0, 5, 16, 14, 0,
		8, 11, 12, 3, 11, 23, 27, 3,
		-26, 3, -3, -1, 10, -5, -7, -15,
	},
	{
		// MG Rook PST
		3, 1, 0, 7, 7, -1, 0, 0,
		-6, -9, 7, 7, 7, 5, -4, -1,
		-12, 11, 0, 17, -2, 12, 23, -1,
		-17, -9, 4, 0, 3, 15, -1, -2,
		-24, -16, -16, -4, -1, -14, 2, -20,
		-30, -15, -6, -3, 0, 2, 2, -15,
		-25, -6, -6, 5, 8, 6, 8, -46,
		-3, 1, 6, 15, 17, 14, -13, -2,
	},
	{
		// MG Queen PST
		-10, 0, 0, 0, 10, 9, 5, 7,
		-19, -35, -5, 2, -9, 7, 1, 15,
		-10, -7, -4, -9, 15, 29, 24, 22,
		-14, -14, -15, -11, -1, -5, 3, -6,
		-8, -20, -8, -5, -4, -2, 2, -2,
		-13, 5, 2, 1, -1, 8, 4, 2,
		-20, 0, 10, 16, 16, 16, -6, 6,
		-3, -1, 7, 19, 5, -10, -9, -17,
	},
	{
		// MG King PST
		-3, 0, 2, 0, 0, 0, 1, -1,
		1, 4, 0, 7, 4, 2, 3, -2,
		2, 4, 7, 4, 4, 14, 12, 0,
		0, 2, 6, 0, 0, 2, 6, -9,
		-8, 5, 0, -8, -10, -10, -9, -23,
		-3, 5, 1, -8, -12, -12, 8, -24,
		6, 13, 0, -40, -23, -1, 25, 19,
		-28, 29, 17, -53, 2, -25, 34, 15,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		77, 74, 63, 53, 59, 60, 72, 77,
		17, 11, 11, 11, 11, -6, 14, 8,
		-3, -14, -18, -31, -29, -25, -20, -18,
		-12, -14, -24, -31, -29, -28, -27, -28,
		-22, -20, -25, -20, -21, -24, -34, -34,
		-16, -22, -11, -19, -13, -23, -32, -34,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-36, -16, -7, -14, -4, -20, -20, -29,
		-17, 2, -7, 14, 2, -7, -9, -19,
		-13, -7, 14, 12, 4, 6, 0, -13,
		-5, 8, 24, 18, 22, 15, 11, -4,
		-3, 4, 20, 30, 22, 25, 15, -2,
		-7, 1, 3, 19, 10, -2, -4, -4,
		-10, -2, -1, 0, 6, -8, -3, -13,
		-12, -28, -8, 1, -5, -12, -27, -12,
	},
	{
		// EG Bishop PST
		-9, -5, -9, -5, -2, -4, -5, -8,
		0, 2, 8, -7, 1, 0, -2, -8,
		8, 0, 0, 1, 0, 1, 5, 6,
		0, 7, 7, 8, 3, 5, 2, 6,
		-1, 0, 12, 8, 0, 6, 0, -5,
		0, 0, 3, 6, 8, -1, 0, -1,
		-6, -12, -7, 0, 0, -8, -9, -13,
		-11, 0, -6, 0, -3, -4, -5, -9,
	},
	{
		// EG Rook PST
		8, 9, 11, 13, 13, 12, 13, 9,
		3, 5, 1, 0, -1, 0, 6, 2,
		9, 5, 7, 2, 2, 1, 0, 0,
		3, 3, 6, 0, 0, 0, 0, 4,
		5, 4, 9, 0, -3, -2, -6, -2,
		0, 0, -6, -5, -9, -14, -7, -12,
		-2, -5, -1, -7, -9, -11, -13, -1,
		-7, -3, 0, -8, -13, -12, -4, -24,
	},
	{
		// EG Queen PST
		-12, 4, 8, 4, 10, 9, 3, 6,
		-17, -7, -1, 7, 3, 6, 1, 0,
		-5, -1, -4, 12, 14, 20, 12, 14,
		-2, 2, 2, 9, 13, 7, 18, 22,
		-9, 3, 1, 15, 5, 10, 12, 10,
		-6, -20, 0, -15, 0, -1, 10, 7,
		-6, -14, -31, -27, -19, -12, -11, -4,
		-12, -22, -19, -30, -8, -13, -6, -15,
	},
	{
		// EG King PST
		-15, -11, -11, -6, -2, 3, 4, -9,
		-9, 14, 11, 13, 13, 28, 19, 1,
		-1, 18, 19, 15, 16, 35, 34, 4,
		-12, 14, 21, 25, 19, 25, 18, -5,
		-23, -6, 14, 21, 20, 18, 5, -16,
		-21, -6, 5, 13, 15, 9, -2, -12,
		-27, -10, 2, 9, 9, 1, -12, -26,
		-43, -34, -20, -5, -26, -9, -35, -55,
	},
}

var PassedPawnPSQT_MG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	45, 52, 42, 43, 28, 34, 19, 9,
	48, 43, 43, 30, 24, 31, 12, 2,
	28, 17, 13, 10, 10, 19, 6, 1,
	14, 0, -9, -7, -13, -7, 9, 16,
	5, 3, -3, -14, -3, 10, 13, 19,
	8, 9, 2, -8, -3, 8, 16, 9,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPSQT_EG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	77, 74, 63, 53, 59, 60, 72, 77,
	91, 83, 66, 40, 30, 61, 67, 84,
	55, 52, 42, 35, 30, 34, 56, 52,
	29, 26, 21, 18, 17, 19, 34, 30,
	8, 6, 5, 1, 1, -1, 14, 7,
	2, 3, -4, 0, -2, -1, 7, 6,
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

var FlipRank = [2][8]uint8{
	{Rank8, Rank7, Rank6, Rank5, Rank4, Rank3, Rank2, Rank1},
	{Rank1, Rank2, Rank3, Rank4, Rank5, Rank6, Rank7, Rank8},
}

// Evaluate a position and give a score, from the perspective of the side to move (
// more positive if it's good for the side to move, otherwise more negative).
func EvaluatePos(pos *Position) int16 {
	if isDrawn(pos) {
		return Draw
	}

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

	eval.MGScores[pos.SideToMove] += TempoBonusMG

	mgScore := eval.MGScores[pos.SideToMove] - eval.MGScores[pos.SideToMove^1]
	egScore := eval.EGScores[pos.SideToMove] - eval.EGScores[pos.SideToMove^1]

	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	score := int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))

	if isDrawish(pos) {
		return score / ScaleFactor
	}

	return score
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
	usPawns := pos.Pieces[color][Pawn]
	enemyPawns := pos.Pieces[color^1][Pawn]

	if OutpostMasks[color][sq]&enemyPawns == 0 &&
		PawnAttacks[color^1][sq]&usPawns != 0 &&
		FlipRank[color][RankOf(sq)] >= Rank5 {

		eval.MGScores[color] += KnightOnOutpostBonusMG
		eval.EGScores[color] += KnightOnOutpostBonusEG
	}

	usBB := pos.Sides[color]
	moves := KnightMoves[sq] & ^usBB
	safeMoves := moves

	for enemyPawns != 0 {
		sq := enemyPawns.PopBit()
		safeMoves &= ^PawnAttacks[color^1][sq]
	}

	mobility := int16(safeMoves.CountBits())
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
	usPawns := pos.Pieces[color][Pawn]
	enemyPawns := pos.Pieces[color^1][Pawn]

	if OutpostMasks[color][sq]&enemyPawns == 0 &&
		PawnAttacks[color^1][sq]&usPawns != 0 &&
		FlipRank[color][RankOf(sq)] >= Rank5 {

		eval.MGScores[color] += BishopOutPostBonusMG
		eval.EGScores[color] += BishopOutPostBonusEG
	}

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
	enemyKingSq := pos.Pieces[color^1][King].Msb()
	if FlipRank[color][RankOf(sq)] == Rank7 && FlipRank[color][RankOf(enemyKingSq)] >= Rank7 {
		eval.EGScores[color] += RookOrQueenOnSeventhBonusEG
	}

	pawns := pos.Pieces[White][Pawn] | pos.Pieces[Black][Pawn]
	if MaskFile[FileOf(sq)]&pawns == 0 {
		eval.MGScores[color] += RookOnOpenFileBonusMG
	}

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
	enemyKingSq := pos.Pieces[color^1][King].Msb()
	if FlipRank[color][RankOf(sq)] == Rank7 && FlipRank[color][RankOf(enemyKingSq)] >= Rank7 {
		eval.EGScores[color] += RookOrQueenOnSeventhBonusEG
	}

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

// Determine if a position is a strict draw (e.g. king versus king,
// or king versus knight and king).
func isDrawn(pos *Position) bool {
	whiteKnights := pos.Pieces[White][Knight].CountBits()
	whiteBishops := pos.Pieces[White][Bishop].CountBits()

	blackKnights := pos.Pieces[Black][Knight].CountBits()
	blackBishops := pos.Pieces[Black][Bishop].CountBits()

	pawns := pos.Pieces[White][Pawn].CountBits() + pos.Pieces[Black][Pawn].CountBits()
	knights := whiteKnights + blackKnights
	bishops := whiteBishops + blackBishops
	rooks := pos.Pieces[White][Rook].CountBits() + pos.Pieces[Black][Rook].CountBits()
	queens := pos.Pieces[White][Queen].CountBits() + pos.Pieces[Black][Queen].CountBits()

	majors := rooks + queens
	miniors := knights + bishops

	if pawns+majors+miniors == 0 {
		// KvK => draw
		return true
	} else if majors+pawns == 0 {
		if miniors == 1 {
			// K & minior v K => draw
			return true
		} else if miniors == 2 && whiteKnights == 1 && blackKnights == 1 {
			// KNvKN => draw
			return true
		} else if miniors == 2 && whiteBishops == 1 && blackBishops == 1 {
			// KBvKB => draw when only when bishops are the same color
			whiteBishopSq := pos.Pieces[White][Bishop].Msb()
			blackBishopSq := pos.Pieces[Black][Bishop].Msb()
			return sqIsDark(whiteBishopSq) == sqIsDark(blackBishopSq)
		}
	}

	return false
}

// Determine if a position is drawish.
func isDrawish(pos *Position) bool {
	whiteKnights := pos.Pieces[White][Knight].CountBits()
	whiteBishops := pos.Pieces[White][Bishop].CountBits()
	whiteRooks := pos.Pieces[White][Rook].CountBits()
	whiteQueens := pos.Pieces[White][Queen].CountBits()

	blackKnights := pos.Pieces[Black][Knight].CountBits()
	blackBishops := pos.Pieces[Black][Bishop].CountBits()
	blackRooks := pos.Pieces[Black][Rook].CountBits()
	blackQueens := pos.Pieces[Black][Queen].CountBits()

	pawns := pos.Pieces[White][Pawn].CountBits() + pos.Pieces[Black][Pawn].CountBits()
	knights := whiteKnights + blackKnights
	bishops := whiteBishops + blackBishops
	rooks := whiteRooks + blackRooks
	queens := whiteQueens + blackQueens

	whiteMinors := whiteBishops + whiteKnights
	blackMinors := blackBishops + blackKnights

	majors := rooks + queens
	miniors := knights + bishops
	all := majors + miniors

	if pawns == 0 {
		if all == 2 && blackQueens == 1 && whiteQueens == 1 {
			// KQ v KQ => drawish
			return true
		} else if all == 2 && blackRooks == 1 && whiteRooks == 1 {
			// KR v KR => drawish
			return true
		} else if all == 2 && whiteMinors == 1 && blackMinors == 1 {
			// KN v KB => drawish
			// KB v KB => drawish
			return true
		} else if all == 3 && ((whiteQueens == 1 && blackRooks == 2) || (blackQueens == 1 && whiteRooks == 2)) {
			// KQ v KRR => drawish
			return true
		} else if all == 3 && ((whiteQueens == 1 && blackBishops == 2) || (blackQueens == 1 && whiteBishops == 2)) {
			// KQ vs KBB => drawish
			return true
		} else if all == 3 && ((whiteQueens == 1 && blackKnights == 2) || (blackQueens == 1 && whiteKnights == 2)) {
			// KQ vs KNN => drawish
			return true
		} else if all <= 3 && ((whiteKnights == 2 && blackMinors <= 1) || (blackKnights == 2 && whiteMinors <= 1)) {
			// KNN v KN => drawish
			// KNN v KB => drawish
			// KNN v K => drawish
			return true
		} else if all == 3 &&
			((whiteQueens == 1 && blackRooks == 1 && blackMinors == 1) ||
				(blackQueens == 1 && whiteRooks == 1 && whiteMinors == 1)) {
			// KQ vs KRN => drawish
			// KQ vs KRB => drawish
			return true
		} else if all == 3 &&
			((whiteRooks == 1 && blackRooks == 1 && blackMinors == 1) ||
				(blackRooks == 1 && whiteRooks == 1 && whiteMinors == 1)) {
			// KR vs KRB => drawish
			// KR vs KRN => drawish
		} else if all == 4 &&
			((whiteRooks == 2 && blackRooks == 1 && blackMinors == 1) ||
				(blackRooks == 2 && whiteRooks == 1 && whiteMinors == 1)) {
			// KRR v KRB => drawish
			// KRR v KRN => drawish
			return true
		}
	}

	return false
}

// Determine if a square is dark.
func sqIsDark(sq uint8) bool {
	fileNo := FileOf(sq)
	rankNo := RankOf(sq)
	return ((fileNo + rankNo) % 2) == 0
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

		// Passed pawn masks and outpost masks.
		frontSpanMask := fileBB
		frontSpanMask |= (fileBB & ClearFile[FileA]) << 1
		frontSpanMask |= (fileBB & ClearFile[FileH]) >> 1

		whiteFrontSpan := frontSpanMask
		for r := 0; r <= rank; r++ {
			whiteFrontSpan &= ClearRank[r]
		}

		PassedPawnMasks[White][sq] = whiteFrontSpan
		OutpostMasks[White][sq] = whiteFrontSpan & ^fileBB

		blackFrontSpan := frontSpanMask
		for r := 7; r >= rank; r-- {
			blackFrontSpan &= ClearRank[r]
		}

		PassedPawnMasks[Black][sq] = blackFrontSpan
		OutpostMasks[Black][sq] = blackFrontSpan & ^fileBB
	}
}
