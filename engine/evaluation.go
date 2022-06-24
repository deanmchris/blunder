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
var OutpostMasks [2][64]Bitboard

var PieceValueMG = [6]int16{94, 352, 365, 480, 978}
var PieceValueEG = [6]int16{117, 253, 279, 498, 932}

var PieceMobilityMG = [5]int16{0, 1, 4, 5, 0}
var PieceMobilityEG = [5]int16{0, 0, 3, 2, 6}

var BishopPairBonusMG int16 = 28
var BishopPairBonusEG int16 = 45

var IsolatedPawnPenatlyMG int16 = 16
var IsolatedPawnPenatlyEG int16 = 7
var DoubledPawnPenatlyMG int16 = 4
var DoubledPawnPenatlyEG int16 = 14

// A middlegame equivelent of this bonus is not missing,
// and one was used at first. But several thousand
// iterations of the tuner indicated that such a "bonus"
// was actually quite bad to give in the middlegame.
var RookOrQueenOnSeventhBonusEG int16 = 17

var KnightOnOutpostBonusMG int16 = 17
var KnightOnOutpostBonusEG int16 = 16

var OuterRingAttackPoints = [5]int16{0, 1, 0, 1, 0}
var InnerRingAttackPoints = [5]int16{0, 2, 4, 4, 3}
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
		55, 73, 40, 53, 36, 46, 14, 3,
		-28, -14, 5, 10, 44, 47, 6, -30,
		-26, -3, -3, 16, 15, 6, 2, -28,
		-38, -24, -7, 7, 11, 2, -12, -34,
		-30, -26, -8, -9, 2, 1, 9, -17,
		-37, -19, -27, -18, -13, 18, 17, -24,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-75, -16, -10, -12, 3, -33, -7, -36,
		-54, -33, 45, 11, 1, 32, -3, -16,
		-29, 38, 20, 41, 56, 66, 46, 12,
		-12, 12, 2, 37, 17, 51, 6, 8,
		-12, 0, 13, 6, 21, 11, 10, -11,
		-21, -6, 13, 14, 24, 20, 20, -13,
		-24, -33, -8, 8, 10, 18, -9, -14,
		-37, -11, -39, -23, -4, -17, -9, -19,
	},
	{
		// MG Bishop PST
		-13, -4, -26, -12, -7, -12, -1, 0,
		-21, 0, -24, -17, 3, 19, 0, -49,
		-21, 17, 19, 11, 5, 23, 8, -9,
		-8, -1, 0, 25, 13, 10, 0, -10,
		-4, 4, -1, 13, 17, -2, 0, 2,
		0, 12, 14, 3, 10, 27, 14, 6,
		1, 23, 11, 5, 13, 20, 39, 2,
		-23, 0, 3, -4, 0, 1, -17, -16,
	},
	{
		// MG Rook PST
		10, 17, 0, 22, 18, -1, 0, 3,
		12, 11, 38, 35, 38, 34, 2, 11,
		-12, 4, 8, 9, -3, 10, 25, -1,
		-26, -15, -1, 6, 0, 12, -7, -21,
		-37, -19, -8, -5, -1, -11, 0, -28,
		-32, -16, -4, -6, 0, 0, -7, -26,
		-26, -6, -5, 3, 7, 10, 0, -49,
		-1, 2, 16, 24, 23, 17, -17, -2,
	},
	{
		// MG Queen PST
		-18, -1, 0, 0, 25, 13, 12, 20,
		-19, -46, -3, 0, -6, 18, 9, 22,
		-10, -14, -1, -7, 10, 29, 19, 29,
		-22, -23, -18, -22, -10, -1, -4, -10,
		-3, -26, -7, -10, -6, -3, -3, -3,
		-13, 8, -1, 5, 1, 4, 10, 2,
		-19, 0, 18, 17, 24, 20, 0, 0,
		1, 2, 12, 25, 4, -10, -10, -33,
	},
	{
		// MG King PST
		-3, 0, 1, 0, -2, 0, 0, 0,
		2, 5, 3, 9, 4, 6, 0, -2,
		4, 10, 13, 4, 6, 18, 19, 0,
		0, 4, 6, 0, 0, 0, 4, -11,
		-10, 7, -5, -22, -30, -21, -19, -31,
		0, 2, -7, -27, -31, -24, 7, -16,
		6, 23, 0, -46, -27, -4, 31, 23,
		-13, 47, 26, -53, 8, -18, 39, 23,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		103, 97, 83, 64, 72, 68, 86, 106,
		14, 26, 23, 10, 0, -11, 17, 13,
		-16, -25, -30, -42, -37, -32, -32, -26,
		-23, -25, -35, -39, -39, -40, -38, -37,
		-33, -31, -37, -31, -31, -36, -46, -45,
		-25, -33, -20, -26, -23, -35, -44, -45,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-54, -29, -3, -20, -7, -25, -34, -52,
		-16, 2, -7, 14, 5, -6, -7, -32,
		-14, -2, 22, 19, 8, 10, -2, -20,
		-1, 15, 34, 31, 34, 20, 18, -2,
		0, 8, 28, 40, 29, 29, 16, 0,
		-6, 11, 10, 25, 20, 6, -6, -7,
		-17, -4, 6, 6, 7, -2, -7, -24,
		-24, -32, -4, 1, -4, -2, -34, -29,
	},
	{
		// EG Bishop PST
		-9, -10, -11, -7, 0, -6, -3, -13,
		0, 0, 6, -7, 0, -1, 0, -2,
		8, -1, -1, -2, -1, 0, 4, 7,
		2, 8, 9, 3, 5, 3, 0, 6,
		0, 3, 9, 10, 0, 4, 0, 0,
		0, 1, 5, 7, 10, -3, 0, -4,
		-1, -12, -2, 0, 0, -3, -13, -16,
		-7, 2, -7, 3, 1, -4, 0, -3,
	},
	{
		// EG Rook PST
		7, 3, 10, 6, 8, 9, 7, 4,
		1, 2, -1, -1, -7, -2, 5, 2,
		6, 4, 1, 2, 0, -1, -2, -2,
		7, 1, 9, 0, 0, 0, -2, 4,
		10, 5, 6, 1, -2, -2, -6, -1,
		3, 2, -4, -1, -5, -8, -4, -7,
		0, -3, 0, 0, -7, -8, -9, 3,
		-6, -1, -3, -7, -9, -9, 0, -23,
	},
	{
		// EG Queen PST
		-18, 3, 6, 5, 14, 10, 5, 10,
		-15, 1, 4, 13, 14, 8, 3, 0,
		-18, -10, -11, 17, 18, 12, 5, 2,
		2, 5, 0, 12, 20, 10, 26, 22,
		-17, 11, 0, 18, 5, 7, 18, 8,
		-2, -36, 0, -10, -2, 4, 4, 2,
		-7, -18, -33, -23, -23, -16, -19, -8,
		-17, -27, -21, -34, -3, -12, -11, -25,
	},
	{
		// EG King PST
		-25, -17, -10, -12, -9, 5, 0, -5,
		-6, 9, 7, 10, 10, 30, 16, 3,
		0, 12, 15, 10, 13, 38, 35, 2,
		-14, 13, 18, 23, 21, 27, 20, -3,
		-23, -6, 17, 24, 26, 21, 7, -12,
		-22, -5, 9, 20, 23, 17, 1, -10,
		-29, -15, 3, 14, 14, 4, -11, -25,
		-55, -43, -24, -4, -24, -8, -34, -56,
	},
}

var PassedPawnPSQT_MG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	49, 54, 34, 38, 37, 43, 48, 46,
	45, 27, 24, 16, 12, 21, 15, 22,
	20, 10, 11, 2, 7, 12, 4, 3,
	14, -6, -7, -12, -2, 0, 1, 12,
	3, 0, -7, -9, 0, 4, 2, 8,
	1, 5, 2, -3, 0, 4, 5, 2,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPSQT_EG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	59, 57, 54, 51, 55, 54, 63, 63,
	101, 79, 57, 42, 37, 63, 64, 82,
	63, 56, 45, 39, 30, 34, 59, 55,
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
	enemyKingSq := pos.Pieces[color^1][King].Msb()
	if FlipRank[color][RankOf(sq)] == Rank7 && FlipRank[color][RankOf(enemyKingSq)] >= Rank7 {
		eval.EGScores[color] += RookOrQueenOnSeventhBonusEG
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
