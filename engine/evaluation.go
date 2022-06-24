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

var PieceValueMG = [6]int16{90, 350, 368, 480, 978}
var PieceValueEG = [6]int16{116, 251, 280, 498, 932}

var PieceMobilityMG = [5]int16{0, 4, 5, 5, 1}
var PieceMobilityEG = [5]int16{0, 0, 2, 2, 6}

var BishopPairBonusMG int16 = 30
var BishopPairBonusEG int16 = 45

var IsolatedPawnPenatlyMG int16 = 15
var IsolatedPawnPenatlyEG int16 = 7
var DoubledPawnPenatlyMG int16 = 4
var DoubledPawnPenatlyEG int16 = 14

// A middlegame equivelent of this bonus is not missing,
// and one was used at first. But several thousand
// iterations of the tuner indicated that such a "bonus"
// was actually quite bad to give in the middlegame.
var RookOrQueenOnSeventhBonusEG int16 = 17

var KnightOnOutpostBonusMG int16 = 18
var KnightOnOutpostBonusEG int16 = 16

var OuterRingAttackPoints = [5]int16{0, 1, 0, 1, 0}
var InnerRingAttackPoints = [5]int16{0, 2, 4, 4, 3}
var SemiOpenFileNextToKingPenalty int16 = 3

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
		54, 72, 39, 53, 35, 46, 13, 2,
		-27, -14, 4, 9, 44, 47, 6, -29,
		-26, -4, -3, 15, 13, 5, 0, -27,
		-37, -25, -8, 5, 10, 1, -12, -33,
		-31, -25, -9, -9, 2, 0, 9, -18,
		-35, -18, -25, -16, -12, 19, 18, -22,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-75, -16, -10, -12, 3, -33, -7, -36,
		-54, -33, 45, 11, 1, 32, -3, -16,
		-29, 37, 19, 40, 55, 66, 45, 12,
		-12, 11, 1, 36, 17, 50, 5, 8,
		-11, 0, 12, 6, 20, 11, 10, -10,
		-20, -6, 10, 13, 24, 17, 20, -12,
		-23, -33, -8, 8, 10, 18, -8, -13,
		-37, -9, -38, -22, -3, -16, -7, -19,
	},
	{
		// MG Bishop PST
		-13, -4, -26, -12, -7, -12, -1, 0,
		-21, 0, -24, -17, 3, 18, 0, -49,
		-21, 16, 18, 10, 4, 23, 7, -9,
		-8, 0, 0, 25, 13, 10, 0, -10,
		-4, 3, 0, 13, 17, -2, 0, 2,
		0, 12, 14, 5, 11, 27, 13, 6,
		1, 22, 11, 5, 14, 20, 38, 2,
		-23, 0, 2, -3, 0, 0, -17, -16,
	},
	{
		// MG Rook PST
		9, 17, 0, 21, 18, -1, 0, 2,
		11, 10, 37, 34, 37, 33, 1, 10,
		-11, 4, 8, 9, -2, 10, 25, -1,
		-25, -14, 0, 6, 0, 12, -6, -20,
		-36, -18, -7, -4, 0, -10, 0, -27,
		-32, -15, -3, -6, 0, 0, -6, -26,
		-26, -6, -5, 3, 7, 9, 0, -49,
		-3, 2, 16, 24, 23, 17, -16, -3,
	},
	{
		// MG Queen PST
		-18, -1, 0, 0, 25, 13, 11, 20,
		-19, -46, -3, 0, -6, 17, 8, 21,
		-10, -13, -1, -7, 10, 29, 18, 28,
		-22, -22, -17, -21, -9, 0, -3, -10,
		-3, -25, -6, -10, -5, -2, -2, -2,
		-13, 8, -1, 4, 1, 4, 10, 2,
		-19, 0, 17, 16, 23, 20, 0, 0,
		1, 2, 12, 25, 4, -10, -10, -33,
	},
	{
		// MG King PST
		-3, 0, 1, 0, -2, 0, 0, 0,
		2, 5, 3, 9, 4, 6, 0, -2,
		4, 10, 13, 4, 6, 18, 19, 0,
		0, 4, 6, 0, 0, 0, 4, -11,
		-10, 7, -5, -22, -30, -21, -19, -31,
		0, 2, -7, -26, -30, -23, 7, -16,
		6, 23, 0, -46, -27, -3, 30, 23,
		-13, 47, 26, -53, 7, -18, 40, 22,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		103, 97, 83, 63, 72, 67, 86, 106,
		13, 25, 22, 9, 0, -11, 16, 12,
		-16, -25, -30, -42, -37, -32, -31, -25,
		-22, -25, -35, -39, -39, -39, -38, -36,
		-32, -31, -37, -30, -31, -36, -45, -44,
		-24, -33, -19, -25, -22, -34, -44, -44,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-54, -29, -3, -20, -7, -25, -34, -52,
		-15, 2, -7, 14, 5, -6, -7, -32,
		-14, -2, 21, 18, 7, 9, -2, -20,
		-1, 14, 33, 30, 33, 19, 17, -2,
		0, 8, 27, 40, 29, 29, 16, 0,
		-5, 11, 9, 24, 20, 5, -5, -6,
		-16, -3, 6, 6, 7, -2, -6, -24,
		-23, -31, -3, 1, -3, -1, -33, -29,
	},
	{
		// EG Bishop PST
		-9, -10, -11, -7, 0, -6, -3, -13,
		0, 0, 5, -7, 0, -1, 0, -2,
		7, -1, -1, -2, -1, 0, 3, 6,
		1, 7, 9, 3, 5, 3, 0, 6,
		0, 2, 9, 10, 0, 4, 0, 0,
		0, 1, 5, 7, 10, -2, 0, -4,
		-1, -12, -1, 0, 0, -3, -13, -16,
		-7, 2, -7, 2, 0, -4, 0, -3,
	},
	{
		// EG Rook PST
		6, 3, 10, 5, 8, 8, 7, 4,
		0, 1, -1, -1, -7, -2, 4, 1,
		6, 4, 1, 2, 0, 0, -1, -1,
		7, 1, 9, 0, 0, 0, -1, 4,
		10, 5, 6, 1, -1, -1, -5, 0,
		3, 2, -4, -1, -4, -7, -3, -7,
		0, -3, 0, 0, -7, -8, -9, 2,
		-7, -1, -3, -7, -9, -9, 0, -23,
	},
	{
		// EG Queen PST
		-18, 2, 6, 4, 13, 10, 4, 9,
		-15, 0, 3, 12, 13, 7, 2, 0,
		-18, -9, -11, 17, 18, 12, 5, 2,
		2, 5, 0, 12, 20, 10, 26, 22,
		-17, 11, 0, 18, 5, 7, 18, 8,
		-2, -35, 0, -9, -1, 4, 4, 2,
		-7, -18, -33, -23, -23, -16, -19, -8,
		-17, -27, -21, -34, -3, -12, -11, -25,
	},
	{
		// EG King PST
		-25, -17, -10, -12, -9, 5, 0, -5,
		-6, 9, 7, 10, 10, 30, 16, 3,
		0, 12, 15, 10, 13, 38, 35, 2,
		-14, 13, 18, 23, 20, 27, 20, -3,
		-23, -6, 17, 23, 26, 21, 7, -12,
		-22, -5, 9, 19, 23, 17, 1, -10,
		-29, -15, 3, 13, 14, 4, -11, -24,
		-55, -42, -23, -4, -23, -8, -34, -56,
	},
}

var PassedPawnPSQT_MG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	48, 53, 33, 38, 36, 43, 47, 45,
	45, 27, 24, 16, 11, 21, 15, 21,
	20, 9, 10, 1, 6, 12, 3, 3,
	14, -6, -7, -12, -2, 0, 1, 12,
	3, 0, -7, -9, 0, 4, 2, 8,
	1, 5, 2, -3, 0, 4, 5, 2,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPSQT_EG = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	59, 57, 54, 50, 55, 53, 63, 63,
	101, 79, 57, 42, 37, 63, 64, 82,
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
	enemyPawns := pos.Pieces[color^1][Pawn]
	safeMoves := moves

	for enemyPawns != 0 {
		sq := enemyPawns.PopBit()
		safeMoves &= ^PawnAttacks[color^1][sq]
	}

	mobility := int16(safeMoves.CountBits())
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
	enemyPawns := pos.Pieces[color^1][Pawn]
	safeMoves := moves

	for enemyPawns != 0 {
		sq := enemyPawns.PopBit()
		safeMoves &= ^PawnAttacks[color^1][sq]
	}

	mobility := int16(safeMoves.CountBits())
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
	enemyPawns := pos.Pieces[color^1][Pawn]
	safeMoves := moves

	for enemyPawns != 0 {
		sq := enemyPawns.PopBit()
		safeMoves &= ^PawnAttacks[color^1][sq]
	}

	mobility := int16(safeMoves.CountBits())
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
