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
var KnightOutpustMasks [2][64]Bitboard

var PieceValueMG [6]int16 = [6]int16{183, 677, 729, 1069, 2242}
var PieceValueEG [6]int16 = [6]int16{267, 638, 661, 1223, 2240}
var PieceMobilityMG [4]int16 = [4]int16{5, 8, 13, 1}
var PieceMobilityEG [4]int16 = [4]int16{3, 9, 6, 15}

var PassedPawnBonusMG [8]int16 = [8]int16{0, 6, 8, 4, 30, 91, 135, 0}
var PassedPawnBonusEG [8]int16 = [8]int16{0, 1, 1, 44, 100, 199, 205, 0}

var IsolatedPawnPenatlyMG int16 = 41
var IsolatedPawnPenatlyEG int16 = 10
var DoubledPawnPenatlyMG int16 = 1
var DoubledPawnPenatlyEG int16 = 48

var KnightOutpostBonusMG int16 = 87
var KnightOutpostBonusEG int16 = 24

var MinorAttackOuterRing int16 = 31
var MinorAttackInnerRing int16 = 63
var RookAttackOuterRing int16 = 21
var RookAttackInnerRing int16 = 64
var QueenAttackOuterRing int16 = 29
var QueenAttackInnerRing int16 = 55
var SemiOpenFileNextToKingPenalty int16 = 60

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
		139, 165, 85, 120, 108, 155, -14, -52,
		-26, -30, 26, 9, 96, 142, 34, -36,
		-38, 0, 6, 48, 47, 31, 12, -54,
		-59, -47, -11, 24, 34, 10, -26, -58,
		-40, -52, -15, -19, 8, 6, 26, -20,
		-53, -40, -57, -38, -24, 45, 35, -33,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		-312, -131, -66, -72, 70, -175, -66, -210,
		-154, -96, 132, 51, 31, 103, 7, -56,
		-80, 87, 64, 92, 142, 170, 114, 90,
		-9, 45, 5, 93, 53, 129, 35, 52,
		4, 45, 56, 36, 69, 48, 51, 5,
		-15, 18, 54, 65, 89, 74, 83, 5,
		-6, -53, 13, 47, 54, 78, 19, 19,
		-159, 12, -45, -17, 40, 12, 17, -8,
	},

	// Piece-square table for bishops
	{
		-47, 25, -151, -87, -48, -61, 1, 44,
		-27, 59, -4, -23, 72, 90, 61, -92,
		9, 104, 116, 97, 91, 124, 80, 38,
		48, 53, 67, 113, 101, 79, 58, 27,
		46, 78, 56, 99, 100, 58, 55, 72,
		60, 91, 94, 68, 88, 119, 97, 62,
		76, 109, 89, 76, 89, 116, 133, 85,
		16, 72, 70, 66, 81, 60, 14, 13,
	},

	// Piece square table for rook
	{
		-16, 17, -48, 26, 31, -40, -20, -17,
		-4, -4, 55, 59, 84, 74, -23, 23,
		-51, -13, -11, -16, -53, 26, 73, -35,
		-86, -76, -30, -7, -39, 6, -74, -89,
		-105, -82, -52, -49, -34, -63, -31, -97,
		-98, -65, -40, -47, -27, -30, -64, -89,
		-85, -45, -42, -19, -8, -1, -43, -138,
		-27, -18, 13, 33, 31, 22, -70, -33,
	},

	// Piece square table for queens
	{
		-89, -42, -31, -35, 66, 61, 30, 40,
		-70, -118, -45, -18, -104, 6, -3, 39,
		-47, -44, -21, -64, -3, 56, 4, 46,
		-77, -81, -66, -64, -52, -41, -54, -59,
		-31, -84, -42, -45, -40, -36, -39, -36,
		-56, 0, -25, -11, -15, -12, -2, -16,
		-65, -14, 18, 18, 28, 40, -11, 16,
		9, -1, 17, 38, 7, -24, -21, -95,
	},

	// Piece square table for kings
	{
		-81, 117, 134, 87, -21, 8, 48, 33,
		138, 94, 64, 118, 65, 79, 24, -52,
		77, 108, 121, 50, 80, 138, 142, 13,
		45, 57, 67, 31, 42, 38, 60, -49,
		-57, 80, 2, -35, -44, -10, -43, -96,
		15, 20, -7, -32, -37, -27, 14, -50,
		19, 49, -8, -111, -74, -9, 32, 18,
		-72, 73, 41, -126, 0, -57, 57, 3,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		194, 177, 127, 82, 103, 88, 173, 225,
		71, 70, 11, -52, -83, -49, 19, 38,
		2, -19, -50, -89, -79, -64, -36, -28,
		-22, -27, -60, -78, -77, -70, -50, -55,
		-50, -37, -63, -56, -54, -60, -69, -75,
		-31, -39, -26, -47, -34, -57, -61, -77,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		-108, -57, 15, -33, -23, -34, -96, -181,
		-19, 29, -20, 37, 16, -23, -17, -65,
		-19, -4, 51, 49, 25, 28, -6, -59,
		4, 31, 78, 73, 76, 49, 43, -2,
		5, 18, 68, 94, 70, 72, 50, -1,
		-12, 29, 19, 58, 48, 11, -17, -10,
		-44, -2, 17, 21, 23, -10, -13, -62,
		-14, -64, -13, 10, -12, -7, -60, -102,
	},

	// Piece-square table for bishops
	{
		51, 24, 54, 53, 62, 54, 46, 14,
		60, 46, 70, 41, 49, 46, 39, 51,
		72, 41, 40, 40, 38, 48, 62, 70,
		54, 72, 69, 59, 61, 61, 53, 72,
		55, 55, 75, 74, 48, 60, 46, 48,
		49, 54, 65, 68, 72, 47, 51, 49,
		43, 28, 45, 53, 59, 42, 46, 5,
		27, 54, 37, 56, 49, 51, 62, 48,
	},

	// Piece square table for rook
	{
		9, -5, 21, 1, 3, 12, 4, -1,
		5, 9, -3, -2, -30, -18, 6, -8,
		-1, -6, -11, -6, -9, -29, -37, -21,
		0, -7, 6, -21, -10, -15, -15, 5,
		2, 1, 0, -9, -20, -19, -31, -18,
		-12, -10, -28, -20, -30, -35, -17, -36,
		-17, -24, -16, -15, -35, -36, -34, -14,
		-45, -27, -30, -42, -45, -50, -21, -73,
	},

	// Piece square table for queens
	{
		-54, -1, 7, 3, 0, -10, -21, 13,
		-58, -15, 3, 22, 69, 10, 5, -12,
		-66, -58, -64, 42, 38, 6, 24, -4,
		-7, 2, -15, 1, 36, 22, 81, 63,
		-63, 24, -21, 31, -3, 5, 47, 21,
		-10, -95, -19, -46, -31, -5, 14, 1,
		-42, -75, -83, -66, -65, -76, -98, -81,
		-88, -109, -85, -76, -47, -79, -67, -102,
	},

	// Piece square table for kings
	{
		-166, -95, -55, -54, -22, 35, 6, -34,
		-47, 31, 27, 28, 37, 79, 49, 30,
		9, 33, 41, 35, 38, 93, 84, 23,
		-31, 41, 52, 67, 60, 75, 57, 13,
		-36, -11, 55, 70, 75, 60, 36, -11,
		-41, -1, 38, 58, 66, 52, 23, -4,
		-57, -21, 24, 48, 51, 26, 0, -32,
		-102, -74, -40, 0, -37, -7, -50, -91,
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
	return int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))
}

// Evaluate the score of a pawn.
func evalPawn(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Pawn] + PSQT_MG[Pawn][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Pawn] + PSQT_EG[Pawn][FlipSq[color][sq]]

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]
	file := FileOf(sq)

	// Evaluate isolated pawns.
	if IsolatedPawnMasks[file]&usPawns == 0 {
		eval.MGScores[color] -= IsolatedPawnPenatlyMG
		eval.EGScores[color] -= IsolatedPawnPenatlyEG
	}

	// Evaluate doubled pawns.
	if DoubledPawnMasks[color][sq]&usPawns != 0 {
		eval.MGScores[color] -= DoubledPawnPenatlyMG
		eval.EGScores[color] -= DoubledPawnPenatlyEG
	}

	// Evaluate passed pawns.
	if PassedPawnMasks[color][sq]&enemyPawns == 0 {
		eval.MGScores[color] += PassedPawnBonusMG[FlipRank[color][RankOf(sq)]]
		eval.EGScores[color] += PassedPawnBonusEG[FlipRank[color][RankOf(sq)]]
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
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(MinorAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(MinorAttackInnerRing)
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
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(MinorAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(MinorAttackInnerRing)
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
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(RookAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(RookAttackInnerRing)
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
		eval.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(QueenAttackOuterRing)
		eval.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(QueenAttackInnerRing)
	}
}

func evalKing(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PSQT_MG[King][FlipSq[color][sq]]
	eval.EGScores[color] += PSQT_EG[King][FlipSq[color][sq]]
}

func evalKingAttack(pos *Position, color uint8, eval *Eval) {
	points := eval.KingAttackPoints[color]
	enemyKingSq := pos.PieceBB[color^1][King].Msb()
	enemyKingFile := MaskFile[FileOf(enemyKingSq)]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	// Evaluate semi-open files adjacent to the enemy king
	leftFile := ((enemyKingFile & ClearFile[FileA]) << 1)
	rightFile := ((enemyKingFile & ClearFile[FileH]) >> 1)

	if enemyKingFile&enemyPawns == 0 {
		points += uint16(SemiOpenFileNextToKingPenalty)
	}

	if leftFile != 0 && leftFile&enemyPawns == 0 {
		points += uint16(SemiOpenFileNextToKingPenalty)
	}

	if rightFile != 0 && rightFile&enemyPawns == 0 {
		points += uint16(SemiOpenFileNextToKingPenalty)
	}

	// Take all the king saftey points collected for the enemy,
	// and see what kind of penatly we should get by indexing the
	// non-linear king-saftey table.
	points = min_u16(points/10, uint16(len(KingAttackTable)-1))
	if eval.KingAttackers[color] >= 2 && pos.PieceBB[color][Queen] != 0 {
		eval.MGScores[color] += KingAttackTable[points]
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

		// Create knight outpost masks & passed pawn masks.
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

		KnightOutpustMasks[White][sq] = whiteKnightMask
		PassedPawnMasks[White][sq] = whiteFrontSpan

		blackKnightMask := knightMask
		blackFrontSpan := frontSpanMask

		for r := 7; r >= rank; r-- {
			blackKnightMask &= ClearRank[r]
			blackFrontSpan &= ClearRank[r]
		}

		KnightOutpustMasks[Black][sq] = blackKnightMask
		PassedPawnMasks[Black][sq] = blackFrontSpan
	}
}
