package engine

const (
	// Constants which map a piece to how much weight it should have on the phase of the game.
	PawnPhase   int16 = 0
	KnightPhase int16 = 1
	BishopPhase int16 = 1
	RookPhase   int16 = 2
	QueenPhase  int16 = 4
	TotalPhase  int16 = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

	// A constant infinite (checkmate) value.
	Inf int16 = 10000
)

// Variables representing values for draws in the middle and
// end-game.
var MiddleGameDraw int16 = 25
var EndGameDraw int16 = 0

type Eval struct {
	MGScores [2]int16
	EGScores [2]int16
	Scores   [2]int16

	KingZones        [2]KingZone
	KingAttackPoints [2]uint8
	KingAttackers    [2]uint8
}

type KingZone struct {
	OuterRing Bitboard
	InnerRing Bitboard
}

var KingZones [64]KingZone
var IsolatedPawnMasks [8]Bitboard
var DoubledPawnMasks [2][64]Bitboard
var KnightOutpustMasks [2][64]Bitboard

var PieceValueMG [6]int16 = [6]int16{94, 313, 344, 462, 957}
var PieceValueEG [6]int16 = [6]int16{138, 274, 294, 510, 952}
var PieceMobilityMG [4]int16 = [4]int16{3, 4, 5, 1}
var PieceMobilityEG [4]int16 = [4]int16{2, 3, 3, 6}

var KnightOutpostBonusMG int16 = 20
var KnightOutpostBonusEG int16 = 10

var IsolatedPawnPenatly int16 = 6
var DoubledPawnPenatly int16 = 17

var MinorAttackOuterRing int16 = 1
var MinorAttackInnerRing int16 = 5
var RookAttackOuterRing int16 = 2
var RookAttackInnerRing int16 = 6
var QueenAttackOuterRing int16 = 2
var QueenAttackInnerRing int16 = 5

var KingAttackTable [64]int16 = [64]int16{
	0, 0, 1, 2, 3, 5, 7, 9, 12, 15,
	18, 22, 26, 30, 35, 39, 44, 50, 56, 62,
	68, 75, 82, 85, 89, 97, 105, 113, 122, 131,
	140, 150, 169, 180, 191, 202, 213, 225, 237, 248,
	260, 272, 283, 295, 307, 319, 330, 342, 354, 366,
	377, 389, 401, 412, 424, 436, 448, 459, 471, 483,
	494, 500, 500, 500,
}

var PhaseValues [6]int16 = [6]int16{
	PawnPhase,
	KnightPhase,
	BishopPhase,
	RookPhase,
	QueenPhase,
}

var PSQT_MG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		126, 136, 104, 97, 104, 109, 64, 70,
		-5, 6, 26, 26, 50, 56, 20, -19,
		-25, 0, -2, 17, 16, 7, 5, -35,
		-35, -20, -10, 4, 9, -2, -9, -36,
		-28, -20, -11, -14, -2, -1, 13, -20,
		-34, -14, -30, -21, -16, 16, 19, -25,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		-146, -59, -32, -38, 33, -87, -26, -93,
		-75, -47, 52, 10, 4, 45, 1, -27,
		-33, 49, 22, 43, 52, 72, 53, 31,
		-5, 18, 7, 42, 16, 51, 8, 18,
		-4, 9, 15, 9, 25, 14, 15, -3,
		-14, -1, 16, 21, 31, 24, 29, -7,
		-11, -31, -2, 15, 18, 26, 2, 0,
		-75, -2, -32, -12, 9, 0, -1, -13,
	},

	// Piece-square table for bishops
	{
		-27, 1, -76, -38, -38, -46, -12, 14,
		-25, 15, -23, -23, 14, 32, 17, -45,
		-14, 33, 36, 32, 22, 45, 26, 0,
		2, 8, 9, 37, 33, 18, 11, -3,
		8, 20, 12, 29, 31, 11, 12, 19,
		12, 27, 28, 16, 23, 40, 29, 15,
		19, 35, 25, 19, 25, 38, 47, 22,
		-7, 19, 16, 10, 18, 11, -10, -4,
	},

	// Piece square table for rook
	{
		1, 14, -9, 12, 17, -22, -14, -14,
		-1, 7, 28, 26, 34, 32, -5, 13,
		-19, -1, -2, -1, -19, 10, 36, -14,
		-31, -24, -3, 4, -9, 10, -23, -29,
		-39, -28, -9, -8, -6, -23, -14, -33,
		-36, -24, -9, -13, -2, -4, -18, -32,
		-31, -12, -6, -1, 6, 9, -13, -57,
		-6, -2, 13, 22, 21, 19, -26, -10,
	},

	// Piece square table for queens
	{
		-27, -19, -14, -28, 26, 24, 18, 27,
		-30, -50, -21, -9, -55, -7, -3, 23,
		-16, -20, -6, -32, -5, 25, -2, 20,
		-32, -31, -27, -30, -22, -17, -21, -21,
		-8, -33, -12, -17, -12, -12, -11, -12,
		-18, 5, -5, 0, -2, 0, 5, -3,
		-24, -3, 14, 14, 20, 24, 3, 12,
		10, 4, 13, 24, 9, -5, -3, -33,
	},

	// Piece square table for kings
	{
		-22, 49, 70, 44, -9, 8, 31, 10,
		61, 43, 35, 64, 36, 38, 13, -21,
		38, 34, 50, 34, 39, 63, 72, 12,
		20, 30, 30, 24, 23, 15, 23, -22,
		-27, 38, 12, -6, -20, -10, -15, -45,
		5, 10, 4, -6, -9, -6, 10, -17,
		6, 20, 1, -48, -31, -6, 15, 5,
		-35, 33, 18, -53, 1, -27, 23, 0,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		139, 135, 111, 96, 102, 102, 132, 142,
		63, 67, 51, 31, 19, 15, 48, 51,
		2, -11, -22, -33, -40, -33, -23, -16,
		-19, -26, -39, -47, -47, -46, -37, -36,
		-31, -30, -43, -34, -39, -43, -45, -45,
		-21, -29, -24, -29, -27, -38, -41, -45,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{

		-39, -29, 1, -18, -17, -18, -46, -83,
		-12, 7, -16, 12, 1, -18, -12, -35,
		-16, -10, 18, 15, 5, 5, -10, -32,
		-2, 13, 31, 28, 32, 18, 17, -5,
		-6, 2, 23, 32, 22, 23, 15, -6,
		-10, 5, 0, 16, 11, -4, -15, -10,
		-26, -8, 0, 0, 2, -11, -10, -33,
		-18, -35, -10, -2, -11, -10, -31, -49,
	},

	// Piece-square table for bishops
	{
		0, -8, 8, 4, 10, 5, 2, -14,
		7, 3, 17, 0, 6, 1, 0, 0,
		16, 1, 4, 2, 2, 4, 7, 14,
		10, 18, 19, 10, 11, 13, 9, 15,
		5, 6, 17, 16, 4, 13, 4, 2,
		1, 4, 10, 14, 17, 2, 4, 2,
		-2, -9, 3, 7, 10, -3, -1, -17,
		-8, 1, -3, 7, 5, 5, 7, -1,
	},

	// Piece square table for rook
	{
		17, 9, 17, 13, 11, 18, 15, 14,
		15, 14, 12, 12, -1, 4, 12, 9,
		11, 10, 8, 9, 8, 0, -5, 5,
		10, 7, 13, 1, 6, 5, 3, 10,
		11, 10, 8, 3, 0, 2, -1, -1,
		4, 6, -2, 1, -6, -6, 1, -6,
		1, 0, 2, 5, -5, -5, -3, 5,
		-9, 0, -3, -7, -8, -13, 1, -22,
	},

	// Piece square table for queens
	{
		-19, 14, 20, 21, 12, 12, 4, 13,
		-12, 7, 23, 24, 51, 25, 21, 3,
		-18, -7, -16, 39, 31, 11, 24, 16,
		14, 14, 6, 20, 31, 19, 45, 39,
		-15, 21, 2, 27, 12, 17, 32, 24,
		3, -31, 4, -6, -2, 12, 18, 15,
		-6, -19, -25, -18, -12, -23, -27, -23,
		-27, -32, -26, -21, -12, -19, -19, -33,
	},

	// Piece square table for kings
	{
		-72, -34, -22, -21, -5, 16, 10, -7,
		-13, 15, 13, 12, 14, 36, 25, 18,
		8, 21, 23, 13, 17, 40, 38, 13,
		-12, 19, 24, 27, 25, 33, 27, 7,
		-16, -6, 22, 28, 33, 26, 14, -4,
		-19, -3, 12, 21, 25, 19, 8, -5,
		-28, -11, 7, 18, 20, 10, -2, -14,
		-45, -37, -22, -4, -19, -4, -23, -42,
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

var FlipRank [2][8]int = [2][8]int{
	{Rank8, Rank7, Rank6, Rank5, Rank4, Rank3, Rank2, Rank1},
	{Rank1, Rank2, Rank3, Rank4, Rank5, Rank6, Rank7, Rank8},
}

// Evaluate a position and give a score, from the perspective of the side to move (
// more positive if it's good for the side to move, otherwise more negative).
func EvaluatePos(pos *Position) int16 {
	var eval Eval
	eval.KingZones[White] = KingZones[pos.PieceBB[White][King].Msb()]
	eval.KingZones[Black] = KingZones[pos.PieceBB[Black][King].Msb()]

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
		case King:
			evalKing(pos, piece.Color, sq, &eval)
		}

		phase -= PhaseValues[piece.Type]
	}

	evalKingAttack(pos, White, &eval)
	evalKingAttack(pos, Black, &eval)

	mgScore := eval.MGScores[pos.SideToMove] - eval.MGScores[pos.SideToMove^1]
	egScore := eval.EGScores[pos.SideToMove] - eval.EGScores[pos.SideToMove^1]

	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	score := int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))

	score += eval.Scores[pos.SideToMove]
	score -= eval.Scores[pos.SideToMove^1]
	return score
}

// Evaluate the score of a pawn.
func evalPawn(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Pawn] + PSQT_MG[Pawn][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Pawn] + PSQT_EG[Pawn][FlipSq[color][sq]]

	usPawns := pos.PieceBB[color][Pawn]
	file := FileOf(sq)

	// Evaluate isolated pawns.
	if IsolatedPawnMasks[file]&usPawns == 0 {
		eval.Scores[color] -= IsolatedPawnPenatly
	}

	// Evaluate doubled pawns.
	if DoubledPawnMasks[color][sq]&usPawns != 0 {
		eval.Scores[color] -= DoubledPawnPenatly
	}
}

// Evaluate the score of a knight.
func evalKnight(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Knight] + PSQT_MG[Knight][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Knight] + PSQT_EG[Knight][FlipSq[color][sq]]

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	if KnightOutpustMasks[color][sq]&enemyPawns == 0 &&
		PawnAttacks[color^1][sq]&usPawns != 0 &&
		FlipRank[color][RankOf(sq)] >= Rank5 {

		eval.MGScores[color] += KnightOutpostBonusMG
		eval.EGScores[color] += KnightOutpostBonusEG
	}

	usBB := pos.SideBB[color]
	moves := KnightMoves[sq] & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 4) * PieceMobilityMG[Knight-1]
	eval.EGScores[color] += (mobility - 4) * PieceMobilityEG[Knight-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint8(outerRingAttacks.CountBits()) * uint8(MinorAttackOuterRing)
		eval.KingAttackPoints[color] += uint8(innerRingAttacks.CountBits()) * uint8(MinorAttackInnerRing)
	}
}

// Evaluate the score of a bishop.
func evalBishop(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Bishop] + PSQT_MG[Bishop][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Bishop] + PSQT_EG[Bishop][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genBishopMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 7) * PieceMobilityMG[Bishop-1]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Bishop-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint8(outerRingAttacks.CountBits()) * uint8(MinorAttackOuterRing)
		eval.KingAttackPoints[color] += uint8(innerRingAttacks.CountBits()) * uint8(MinorAttackInnerRing)
	}
}

// Evaluate the score of a rook.
func evalRook(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Rook] + PSQT_MG[Rook][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Rook] + PSQT_EG[Rook][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genRookMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 7) * PieceMobilityMG[Rook-1]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Rook-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint8(outerRingAttacks.CountBits()) * uint8(RookAttackOuterRing)
		eval.KingAttackPoints[color] += uint8(innerRingAttacks.CountBits()) * uint8(RookAttackInnerRing)
	}
}

// Evaluate the score of a queen.
func evalQueen(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Queen] + PSQT_MG[Queen][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Queen] + PSQT_EG[Queen][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := (genBishopMoves(sq, allBB) | genRookMoves(sq, allBB)) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 14) * PieceMobilityMG[Queen-1]
	eval.EGScores[color] += (mobility - 14) * PieceMobilityEG[Queen-1]

	outerRingAttacks := moves & eval.KingZones[color^1].OuterRing
	innerRingAttacks := moves & eval.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		eval.KingAttackers[color]++
		eval.KingAttackPoints[color] += uint8(outerRingAttacks.CountBits()) * uint8(QueenAttackOuterRing)
		eval.KingAttackPoints[color] += uint8(innerRingAttacks.CountBits()) * uint8(QueenAttackInnerRing)
	}
}

func evalKing(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PSQT_MG[King][FlipSq[color][sq]]
	eval.EGScores[color] += PSQT_EG[King][FlipSq[color][sq]]
}

func evalKingAttack(pos *Position, color uint8, eval *Eval) {
	points := min_u8(eval.KingAttackPoints[color], uint8(len(KingAttackTable)-1))
	if eval.KingAttackers[color] >= 2 && pos.PieceBB[color][Queen] != 0 {
		eval.MGScores[color] += KingAttackTable[points]
	}
}

func min_u8(a, b uint8) uint8 {
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

		// Create knight outpost masks.
		mask = (fileBB & ClearFile[FileA]) << 1
		mask |= (fileBB & ClearFile[FileH]) >> 1

		whiteKnightMask := mask
		for r := 0; r <= rank; r++ {
			whiteKnightMask &= ClearRank[r]
		}
		KnightOutpustMasks[White][sq] = whiteKnightMask

		blackKnightMask := mask
		for r := 7; r >= rank; r-- {
			blackKnightMask &= ClearRank[r]
		}
		KnightOutpustMasks[Black][sq] = blackKnightMask
	}
}
