package engine

const (
	PawnPhase   int16 = 0
	KnightPhase int16 = 1
	BishopPhase int16 = 1
	RookPhase   int16 = 2
	QueenPhase  int16 = 4
	TotalPhase  int16 = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

	Inf int16 = 10000
)

var MiddleGameDraw int16 = 25
var EndGameDraw int16 = 0

type Eval struct {
	MGScores [2]int16
	EGScores [2]int16

	KingSq           [2]uint8
	KingZones        [2]KingZone
	KingAttackPoints [2]uint16
	KingAttackers    [2]uint8
}

type KingZone struct {
	OuterRing Bitboard
	InnerRing Bitboard
}

var KingZones [64]KingZone
var IsolatedPawnMasks [8]Bitboard
var DoubledPawnMasks [2][64]Bitboard
var PassedPawnMasks [2][64]Bitboard
var MiniorOutpostMasks [2][64]Bitboard

var PieceValueMG = [6]int16{83, 328, 365, 473, 968}
var PieceValueEG = [6]int16{98, 273, 303, 522, 976}
var PieceMobilityMG = [4]int16{2, 3, 4, 1}
var PieceMobilityEG = [4]int16{4, 4, 3, 7}

var PassedPawnBonusMG = [8]int16{0, 9, 4, 1, 13, 48, 109, 0}
var PassedPawnBonusEG = [8]int16{0, 1, 5, 25, 50, 103, 149, 0}

var IsolatedPawnPenatlyMG int16 = 19
var IsolatedPawnPenatlyEG int16 = 5
var DoubledPawnPenatlyMG int16 = 1
var DoubledPawnPenatlyEG int16 = 17

var KnightOutpostBonusMG int16 = 41
var KnightOutpostBonusEG int16 = 8
var BishopOutpostBonusMG int16 = 20
var BishopOutpostBonusEG int16 = 10

var RookOnTheSeventhBonusMG int16 = 5
var RookOnTheSeventhBonusEG int16 = 10

var MinorAttackOuterRing int16 = 1
var MinorAttackInnerRing int16 = 3
var RookAttackOuterRing int16 = 1
var RookAttackInnerRing int16 = 4
var QueenAttackOuterRing int16 = 1
var QueenAttackInnerRing int16 = 3
var SemiOpenFileNextToKingPenalty int16 = 2

var KingAttackTable = [100]int16{
	0, 0, 1, 2, 4, 6, 9, 12, 16, 20, 25, 30, 36,
	42, 49, 56, 64, 72, 81, 90, 100, 110, 121, 132,
	144, 156, 169, 182, 196, 210, 225, 240, 256, 272,
	289, 306, 324, 342, 361, 380, 400, 420, 441, 462,
	484, 506, 529, 552, 576, 600, 625, 650, 676, 702,
	729, 756, 784, 812, 841, 870, 900, 930, 961, 992,
	1024, 1056, 1089, 1122, 1156, 1190, 1225, 1260, 1260,
	1260, 1260, 1260, 1260, 1260, 1260, 1260, 1260, 1260,
	1260, 1260, 1260, 1260, 1260, 1260, 1260, 1260, 1260,
	1260, 1260, 1260, 1260, 1260, 1260, 1260, 1260, 1260,
}

var InitKingSafety = [64]uint16{
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 16,
	6, 6, 8, 10, 10, 8, 6, 6,
	2, 2, 4, 8, 8, 4, 2, 2,
	2, 0, 2, 4, 4, 2, 0, 2,
}

var PhaseValues = [6]int16{
	PawnPhase,
	KnightPhase,
	BishopPhase,
	RookPhase,
	QueenPhase,
}

var PSQT_MG = [6][64]int16{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		13, 36, -13, 9, -2, 33, -64, -65,
		1, -7, 12, -2, 35, 48, 20, -5,
		-13, 3, 4, 25, 23, 17, 6, -19,
		-22, -19, -1, 14, 19, 4, -10, -22,
		-14, -20, -3, -4, 6, 5, 13, -5,
		-20, -16, -22, -12, -8, 20, 18, -11,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		-153, -52, -33, -26, 23, -85, -36, -96,
		-88, -41, 43, 12, 4, 34, -6, -27,
		-44, 30, 17, 30, 40, 59, 40, 23,
		-18, 7, -11, 30, 12, 44, 5, 16,
		-8, 5, 13, 5, 21, 11, 12, -8,
		-18, -3, 13, 17, 27, 23, 26, -9,
		-15, -33, -5, 12, 15, 23, -5, -2,
		-70, -6, -33, -17, 6, -6, -4, -15,
	},

	// Piece-square table for bishops
	{
		-31, -8, -88, -41, -33, -33, -7, -1,
		-34, -1, -27, -34, 10, 26, 7, -58,
		-23, 22, 30, 18, 17, 30, 19, -8,
		-7, -2, 2, 30, 20, 18, 2, -10,
		-4, 11, 2, 21, 24, 2, 4, 10,
		3, 19, 19, 7, 16, 31, 20, 5,
		13, 25, 17, 12, 17, 27, 38, 11,
		-20, 12, 7, 4, 11, 3, -13, -16,
	},

	// Piece square table for rook
	{
		10, 13, -9, 14, 13, -14, 0, 0,
		10, 8, 25, 29, 32, 31, 0, 8,
		-14, 6, 1, 5, -11, 13, 34, -5,
		-26, -19, 0, 2, -2, 10, -23, -26,
		-33, -24, -13, -10, -3, -17, -9, -31,
		-32, -20, -6, -9, -3, -2, -13, -27,
		-30, -9, -6, 3, 9, 6, -5, -54,
		-6, 1, 12, 22, 22, 18, -23, -8,
	},

	// Piece square table for queens
	{
		-22, -16, -6, -8, 33, 23, 10, 25,
		-26, -46, -16, 5, -27, 4, 3, 22,
		-15, -16, -5, -19, 8, 32, 11, 26,
		-28, -31, -26, -26, -16, -9, -12, -18,
		-11, -28, -14, -18, -10, -10, -11, -11,
		-19, 4, -5, -1, -2, 1, 6, 1,
		-21, -3, 13, 14, 19, 22, 2, 14,
		4, 2, 13, 24, 6, -8, -6, -31,
	},

	// Piece square table for kings
	{
		-54, 30, 48, 29, -21, 2, 17, 15,
		36, 33, 20, 54, 34, 30, -1, -15,
		31, 35, 48, 26, 32, 46, 50, 9,
		20, 28, 28, 22, 24, 26, 30, -2,
		-14, 34, 12, 7, 9, 7, 4, -20,
		0, 10, 8, 9, 9, 8, 17, -16,
		0, 16, -1, -35, -14, -5, 12, 3,
		-38, 24, 13, -61, -3, -29, 19, -6,
	},
}

var PSQT_EG = [6][64]int16{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		52, 40, 22, -4, 7, -1, 41, 62,
		34, 34, 11, -18, -30, -7, 14, 20,
		15, 6, -6, -24, -17, -11, 1, 4,
		3, 2, -12, -20, -18, -13, -6, -9,
		-8, -2, -14, -11, -8, -10, -15, -17,
		-2, -4, 0, -10, 1, -8, -12, -19,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{

		-35, -26, 5, -17, -10, -15, -40, -73,
		1, 10, -11, 13, 3, -13, -7, -28,
		-6, -6, 15, 12, 7, 5, -6, -22,
		6, 13, 27, 24, 25, 17, 15, 0,
		5, 7, 24, 32, 23, 24, 21, 3,
		0, 9, 1, 20, 14, -3, -12, -3,
		-12, 1, 6, 4, 8, -8, -2, -23,
		-6, -23, -2, 9, -2, -2, -26, -38,
	},

	// Piece-square table for bishops
	{
		-3, -9, 7, 0, 7, -2, -3, -15,
		7, 4, 14, -1, 2, -2, -2, -3,
		13, -1, -1, 0, -3, 2, 6, 12,
		8, 14, 14, 7, 8, 6, 4, 12,
		7, 5, 14, 13, 0, 7, 0, 2,
		2, 4, 8, 10, 11, -2, 3, 4,
		-1, -8, -1, 1, 6, -1, -3, -15,
		-4, 2, -3, 5, 3, 1, 6, 1,
	},

	// Piece square table for rook
	{
		15, 11, 20, 15, 17, 18, 15, 13,
		14, 17, 16, 13, 5, 10, 15, 13,
		13, 10, 10, 10, 8, 3, -1, 2,
		12, 9, 14, 6, 6, 8, 6, 13,
		12, 12, 13, 7, 2, 4, 2, 4,
		5, 8, -1, 3, 0, -3, 1, -4,
		4, 0, 2, 5, -5, -3, -5, 7,
		-5, 1, 2, -1, -4, -6, 3, -18,
	},

	// Piece square table for queens
	{
		-13, 17, 16, 14, 14, 13, 8, 22,
		-10, 3, 19, 21, 35, 25, 18, 12,
		-9, -12, -12, 28, 29, 21, 27, 20,
		14, 17, 8, 16, 30, 25, 44, 41,
		-6, 17, 8, 27, 14, 21, 34, 33,
		12, -20, 10, -2, 5, 16, 22, 17,
		1, -12, -14, -12, -8, -15, -24, -18,
		-17, -24, -19, -13, 1, -12, -9, -30,
	},

	// Piece square table for kings
	{
		-71, -40, -23, -22, -8, 15, 2, -17,
		-14, 16, 14, 12, 15, 33, 23, 10,
		5, 17, 19, 16, 17, 44, 40, 9,
		-13, 19, 24, 28, 25, 31, 23, 1,
		-19, -3, 23, 27, 27, 24, 10, -10,
		-18, 0, 14, 20, 23, 17, 6, -4,
		-24, -7, 10, 17, 17, 11, 0, -12,
		-45, -29, -17, 1, -15, -4, -21, -39,
	},
}

// Flip white's perspective to black
var FlipSq = [2][64]int{
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

var FlipRank = [2][8]int{
	{Rank8, Rank7, Rank6, Rank5, Rank4, Rank3, Rank2, Rank1},
	{Rank1, Rank2, Rank3, Rank4, Rank5, Rank6, Rank7, Rank8},
}

// Evaluate a position and give a score, from the perspective of the side to move (
// more positive if it's good for the side to move, otherwise more negative).
func EvaluatePos(pos *Position) int16 {
	var eval Eval
	whiteKingSq := pos.PieceBB[White][King].Msb()
	blackKingSq := pos.PieceBB[Black][King].Msb()

	eval.KingZones[White] = KingZones[whiteKingSq]
	eval.KingZones[Black] = KingZones[blackKingSq]

	eval.KingSq[White] = whiteKingSq
	eval.KingSq[Black] = blackKingSq

	phase := TotalPhase
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

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

		phase -= PhaseValues[piece.Type]
	}

	evalKing(pos, White, pos.PieceBB[White][King].Msb(), &eval)
	evalKing(pos, Black, pos.PieceBB[Black][King].Msb(), &eval)

	mgScore := eval.MGScores[pos.SideToMove] - eval.MGScores[pos.SideToMove^1]
	egScore := eval.EGScores[pos.SideToMove] - eval.EGScores[pos.SideToMove^1]

	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	return int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))
}

func evalPawn(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Pawn] + PSQT_MG[Pawn][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Pawn] + PSQT_EG[Pawn][FlipSq[color][sq]]

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	file := FileOf(sq)
	doubled := false

	if IsolatedPawnMasks[file]&usPawns == 0 {
		eval.MGScores[color] -= IsolatedPawnPenatlyMG
		eval.EGScores[color] -= IsolatedPawnPenatlyEG
	}

	if DoubledPawnMasks[color][sq]&usPawns != 0 {
		doubled = true
		eval.MGScores[color] -= DoubledPawnPenatlyMG
		eval.EGScores[color] -= DoubledPawnPenatlyEG
	}

	if PassedPawnMasks[color][sq]&enemyPawns == 0 && !doubled {
		eval.MGScores[color] += PassedPawnBonusMG[FlipRank[color][RankOf(sq)]]
		eval.EGScores[color] += PassedPawnBonusEG[FlipRank[color][RankOf(sq)]]
	}
}

func evalKnight(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Knight] + PSQT_MG[Knight][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Knight] + PSQT_EG[Knight][FlipSq[color][sq]]

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	if MiniorOutpostMasks[color][sq]&enemyPawns == 0 &&
		PawnAttacks[color^1][sq]&usPawns != 0 &&
		FlipRank[color][RankOf(sq)] >= Rank5 {

		eval.MGScores[color] += KnightOutpostBonusMG
		eval.EGScores[color] += KnightOutpostBonusEG
	}

	usBB := pos.SideBB[color]
	moves := KnightMoves[sq] & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 2) * PieceMobilityMG[Knight-1]
	eval.EGScores[color] += (mobility - 4) * PieceMobilityEG[Knight-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(MinorAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(MinorAttackInnerRing)
	}
}

func evalBishop(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Bishop] + PSQT_MG[Bishop][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Bishop] + PSQT_EG[Bishop][FlipSq[color][sq]]

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	if MiniorOutpostMasks[color][sq]&enemyPawns == 0 &&
		PawnAttacks[color^1][sq]&usPawns != 0 &&
		FlipRank[color][RankOf(sq)] >= Rank5 {

		eval.MGScores[color] += BishopOutpostBonusMG
		eval.EGScores[color] += BishopOutpostBonusEG
	}

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genBishopMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 3) * PieceMobilityMG[Bishop-1]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Bishop-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(MinorAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(MinorAttackInnerRing)
	}
}

func evalRook(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Rook] + PSQT_MG[Rook][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Rook] + PSQT_EG[Rook][FlipSq[color][sq]]

	enemyPawns := pos.PieceBB[color^1][Pawn]
	if FlipRank[color][RankOf(sq)] == Rank7 &&
		(MaskRank[Rank7]&enemyPawns != 0 || FlipRank[color][RankOf(eval.KingSq[color^1])] >= Rank7) {
		eval.MGScores[color] += RookOnTheSeventhBonusMG
		eval.EGScores[color] += RookOnTheSeventhBonusEG
	}

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genRookMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 3) * PieceMobilityMG[Rook-1]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Rook-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(RookAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(RookAttackInnerRing)
	}
}

func evalQueen(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Queen] + PSQT_MG[Queen][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Queen] + PSQT_EG[Queen][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := (genBishopMoves(sq, allBB) | genRookMoves(sq, allBB)) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 7) * PieceMobilityMG[Queen-1]
	eval.EGScores[color] += (mobility - 14) * PieceMobilityEG[Queen-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(QueenAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(QueenAttackInnerRing)
	}
}

func evalKing(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PSQT_MG[King][FlipSq[color][sq]]
	eval.EGScores[color] += PSQT_EG[King][FlipSq[color][sq]]

	enemyPoints := InitKingSafety[FlipSq[color][sq]] + eval.KingAttackPoints[color^1]
	kingFile := MaskFile[FileOf(sq)]
	usPawns := pos.PieceBB[color][Pawn]

	// Evaluate semi-open files adjacent to the enemy king
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
	// and see what kind of penatly we should get by indexing the
	// non-linear king-saftey table.
	enemyPoints = min_u16(enemyPoints, uint16(len(KingAttackTable)-1))
	if eval.KingAttackers[color^1] >= 2 && pos.PieceBB[color^1][Queen] != 0 {
		eval.MGScores[color] -= KingAttackTable[enemyPoints]
	}
}

func min_u16(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func init() {
	for sq := 0; sq < 64; sq++ {
		// Create king zones.
		sqBB := SquareBB[sq]
		zone := ((sqBB & ClearFile[FileH]) >> 1) | ((sqBB & (ClearFile[FileG] & ClearFile[FileH])) >> 2)
		zone |= ((sqBB & ClearFile[FileA]) << 1) | ((sqBB & (ClearFile[FileB] & ClearFile[FileA])) << 2)
		zone |= sqBB

		zone |= ((zone >> 8) | (zone >> 16))
		zone |= ((zone << 8) | (zone << 16))
		KingZones[sq] = KingZone{OuterRing: zone &^ (KingMoves[sq] | sqBB), InnerRing: KingMoves[sq] | sqBB}

		// Create isolated pawn masks.
		file := FileOf(uint8(sq))
		fileBB := MaskFile[file]

		mask := (fileBB & ClearFile[FileA]) << 1
		mask |= (fileBB & ClearFile[FileH]) >> 1
		IsolatedPawnMasks[file] = mask

		// Create doubled pawns masks.
		rank := int(RankOf(uint8(sq)))

		mask = fileBB
		for r := 0; r <= rank; r++ {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[White][sq] = mask

		mask = fileBB
		for r := 7; r >= rank; r-- {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[Black][sq] = mask

		// Create minior outpost masks & passed pawn masks.
		knightMask := (fileBB & ClearFile[FileA]) << 1
		knightMask |= (fileBB & ClearFile[FileH]) >> 1

		frontSpanMask := fileBB
		frontSpanMask |= (fileBB & ClearFile[FileA]) << 1
		frontSpanMask |= (fileBB & ClearFile[FileH]) >> 1

		whiteKnightMask := knightMask
		whiteFrontSpan := frontSpanMask

		for r := 0; r <= rank; r++ {
			whiteKnightMask &= ClearRank[r]
			whiteFrontSpan &= ClearRank[r]
		}

		MiniorOutpostMasks[White][sq] = whiteKnightMask
		PassedPawnMasks[White][sq] = whiteFrontSpan

		blackKnightMask := knightMask
		blackFrontSpan := frontSpanMask

		for r := 7; r >= rank; r-- {
			blackKnightMask &= ClearRank[r]
			blackFrontSpan &= ClearRank[r]
		}

		MiniorOutpostMasks[Black][sq] = blackKnightMask
		PassedPawnMasks[Black][sq] = blackFrontSpan
	}
}
