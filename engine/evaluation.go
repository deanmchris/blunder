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

	// A constant representing a draw value.
	Draw int16 = 0
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
var OutpostMasks [2][64]Bitboard

var PieceValueMG = [6]int16{95, 356, 368, 479, 979}
var PieceValueEG = [6]int16{115, 252, 279, 496, 932}

var PieceMobilityMG = [5]int16{0, 6, 4, 4, 0}
var PieceMobilityEG = [5]int16{0, 0, 3, 1, 6}

var BishopPairBonusMG int16 = 28
var BishopPairBonusEG int16 = 45

var IsolatedPawnPenatlyMG int16 = 16
var IsolatedPawnPenatlyEG int16 = 7
var DoubledPawnPenatlyMG int16 = 1
var DoubledPawnPenatlyEG int16 = 14

// A middlegame equivelent of this bonus is not missing,
// and one was used at first. But several thousand
// iterations of the tuner indicated that such a "bonus"
// was actually quite bad to give in the middlegame.
var RookOrQueenOnSeventhBonusEG int16 = 20

var KnightOnOutpostBonusMG int16 = 21
var KnightOnOutpostBonusEG int16 = 16

var RookOnOpenFileBonusMG int16 = 21
var TempoBonusMG int16 = 24

var OuterRingAttackPoints = [5]int16{0, 1, 0, 1, 0}
var InnerRingAttackPoints = [5]int16{0, 2, 4, 3, 3}
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
		53, 72, 39, 53, 35, 46, 11, 0,
		-25, -14, 4, 9, 44, 47, 6, -27,
		-26, -3, -3, 17, 15, 5, 0, -27,
		-37, -25, -6, 8, 11, 1, -13, -32,
		-30, -24, -8, -10, 2, 0, 10, -17,
		-36, -18, -27, -16, -12, 17, 18, -22,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-75, -16, -10, -12, 3, -33, -7, -36,
		-54, -33, 45, 10, 1, 32, -3, -16,
		-29, 35, 17, 38, 53, 66, 43, 12,
		-12, 9, 0, 34, 15, 48, 5, 8,
		-9, 0, 12, 7, 20, 12, 10, -7,
		-19, -6, 9, 13, 24, 16, 20, -10,
		-21, -33, -7, 10, 12, 18, -6, -11,
		-37, -10, -37, -20, -1, -14, -7, -19,
	},
	{
		// MG Bishop PST
		-13, -4, -26, -12, -7, -12, -1, 0,
		-21, 0, -24, -17, 2, 17, 0, -49,
		-21, 15, 18, 10, 4, 23, 5, -10,
		-7, -1, 0, 25, 13, 8, 0, -10,
		-4, 3, 0, 13, 17, -2, 0, 2,
		0, 12, 15, 2, 10, 26, 13, 5,
		1, 22, 10, 6, 15, 20, 39, 2,
		-22, 0, 0, -1, 0, 0, -17, -16,
	},
	{
		// MG Rook PST
		9, 17, 0, 21, 18, -1, 0, 2,
		9, 8, 35, 32, 35, 31, 0, 9,
		-10, 4, 7, 8, -2, 9, 25, -1,
		-25, -14, 0, 4, 0, 12, -6, -19,
		-35, -18, -7, -4, 0, -10, 0, -26,
		-31, -13, -3, -6, 0, 0, -6, -25,
		-26, -6, -5, 3, 7, 9, 0, -48,
		-1, 3, 15, 23, 23, 19, -15, -3,
	},
	{
		// MG Queen PST
		-18, -1, 0, 0, 25, 13, 11, 20,
		-19, -45, -3, 0, -6, 15, 7, 21,
		-10, -13, 0, -7, 10, 29, 16, 26,
		-22, -22, -17, -20, -8, 0, -3, -10,
		-4, -25, -6, -9, -4, -2, -2, -2,
		-13, 8, 0, 4, 1, 4, 10, 2,
		-19, 0, 17, 17, 24, 20, 0, 0,
		1, 2, 12, 23, 4, -8, -8, -33,
	},
	{
		// MG King PST
		-3, 0, 1, 0, -2, 0, 0, 0,
		2, 5, 3, 9, 4, 6, 0, -2,
		4, 10, 13, 4, 6, 18, 19, 0,
		0, 4, 6, 0, 0, 0, 4, -11,
		-10, 7, -5, -22, -30, -21, -19, -31,
		0, 2, -6, -24, -30, -22, 7, -16,
		6, 23, 0, -45, -26, -3, 30, 22,
		-13, 46, 24, -53, 6, -18, 40, 21,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		103, 97, 83, 63, 72, 67, 86, 106,
		11, 23, 22, 9, 0, -11, 15, 11,
		-16, -24, -29, -40, -35, -31, -30, -25,
		-22, -24, -34, -38, -37, -38, -37, -35,
		-32, -30, -36, -29, -30, -34, -45, -43,
		-24, -33, -19, -25, -21, -33, -43, -44,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-54, -29, -3, -20, -6, -24, -34, -52,
		-13, 2, -7, 14, 5, -6, -7, -32,
		-13, -2, 19, 16, 5, 7, -2, -20,
		0, 12, 31, 28, 31, 17, 15, -1,
		0, 8, 27, 40, 29, 29, 16, 0,
		-4, 11, 9, 24, 20, 5, -5, -4,
		-15, -1, 6, 6, 7, -2, -5, -23,
		-21, -29, -2, 1, -2, -1, -32, -29,
	},
	{
		// EG Bishop PST
		-9, -10, -11, -7, 0, -6, -3, -13,
		0, 0, 5, -7, 0, -1, 0, -2,
		7, -1, -1, -2, -1, 0, 3, 5,
		1, 7, 9, 3, 5, 3, 0, 6,
		0, 2, 9, 10, 0, 4, 0, 0,
		0, 1, 5, 6, 9, -2, 0, -4,
		-1, -12, -1, 0, 0, -3, -13, -16,
		-7, 2, -7, 1, 0, -4, 0, -3,
	},
	{
		// EG Rook PST
		6, 3, 10, 5, 8, 8, 7, 4,
		0, 1, -2, -1, -8, -2, 4, 0,
		6, 4, 1, 2, 0, 0, -1, -1,
		7, 1, 9, 0, 0, 0, -1, 4,
		10, 5, 6, 1, -1, -1, -5, 0,
		3, 2, -4, -1, -4, -7, -3, -7,
		0, -3, 0, 0, -6, -7, -9, 2,
		-5, 0, -3, -7, -9, -9, 0, -22,
	},
	{
		// EG Queen PST
		-18, 2, 6, 4, 13, 10, 4, 9,
		-15, 0, 1, 10, 11, 5, 2, 0,
		-18, -8, -11, 17, 18, 12, 5, 2,
		2, 5, 0, 12, 20, 10, 26, 22,
		-17, 11, 0, 18, 5, 6, 18, 8,
		0, -35, 0, -9, -1, 4, 4, 2,
		-6, -18, -31, -23, -23, -16, -19, -8,
		-17, -27, -21, -33, -3, -12, -11, -25,
	},
	{
		// EG King PST
		-25, -17, -10, -12, -9, 5, 0, -5,
		-6, 9, 7, 10, 10, 30, 16, 2,
		0, 12, 15, 10, 13, 38, 34, 2,
		-14, 13, 18, 23, 20, 27, 20, -3,
		-23, -6, 17, 23, 26, 21, 7, -12,
		-22, -5, 9, 19, 23, 16, 1, -10,
		-29, -15, 3, 13, 14, 4, -10, -24,
		-55, -42, -23, -4, -23, -8, -34, -55,
	},
}

var PassedPawnPSQT_MG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	47, 53, 33, 38, 36, 43, 45, 43,
	45, 27, 24, 16, 11, 21, 15, 20,
	20, 9, 10, 1, 6, 12, 2, 1,
	14, -6, -7, -11, -2, 0, 1, 12,
	3, 0, -7, -9, 0, 4, 2, 8,
	1, 5, 2, -3, 0, 4, 5, 2,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPSQT_EG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	59, 57, 54, 50, 55, 53, 63, 63,
	103, 81, 58, 43, 37, 64, 66, 84,
	63, 56, 45, 38, 29, 34, 59, 55,
	30, 28, 21, 18, 15, 18, 34, 30,
	6, 6, 6, 1, 0, 0, 15, 10,
	1, 4, -4, 0, 4, 0, 7, 6,
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
