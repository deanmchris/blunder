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

	// Constants representing indexes into a table of pawn shield
	// bitboards, kingside and queenside.
	Queenside = 0
	Kingside  = 1
)

// Variables representing values for draws in the middle and
// end-game.
var MiddleGameDraw int16 = 25
var EndGameDraw int16 = 0

type Eval struct {
	MGScores     [2]int16
	EGScores     [2]int16
	CastlingSide [2]int16
}

var IsolatedPawnPenatlyMG int16 = 11
var IsolatedPawnPenatlyEG int16 = 3
var DoubledPawnPenatlyMG int16 = 7
var DoubledPawnPenaltyEG int16 = 21
var PawnShieldBonusZone1 int16 = 9
var PawnShieldBonusZone2 int16 = 8

var PieceValueMG [5]int16 = [5]int16{89, 310, 346, 464, 952}
var PieceValueEG [5]int16 = [5]int16{132, 273, 291, 508, 944}
var PieceMobilityMG [4]int16 = [4]int16{2, 4, 7, 1}
var PieceMobilityEG [4]int16 = [4]int16{3, 4, 2, 8}

var PhaseValues [6]int16 = [6]int16{
	PawnPhase,
	KnightPhase,
	BishopPhase,
	RookPhase,
	QueenPhase,
}

var IsolatedPawnMasks [8]Bitboard
var DoubledPawnMasks [2][64]Bitboard
var Zone1PawnShields [2][2]Bitboard = [2][2]Bitboard{
	{0xe000, 0x700},
	{0xe0000000000000, 0x7000000000000},
}
var Zone2PawnShields [2][2]Bitboard = [2][2]Bitboard{
	{0xe00000, 0x70000},
	{0xe00000000000, 0x70000000000},
}

var PSQT_MG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		18, -3, 3, 0, -1, 5, 0, -2,
		104, 143, 115, 102, 102, 108, 69, -4,
		-5, 2, 23, 26, 57, 49, -10, -22,
		-22, 3, -1, 20, 18, 5, 7, -32,
		-29, -13, -6, 7, 13, 3, 0, -31,
		-24, -17, -11, -11, 3, -1, 17, -18,
		-30, -12, -31, -18, -12, 17, 23, -21,
		19, 0, 0, 0, -1, 0, 0, -1,
	},

	// Piece-square table for knights
	{
		-145, -67, -22, -22, 36, -80, -17, -96,
		-72, -48, 64, 20, 28, 53, 3, -11,
		-44, 48, 30, 49, 76, 99, 58, 29,
		-11, 17, 16, 56, 30, 67, 21, 37,
		-11, 8, 18, 11, 31, 19, 24, 0,
		-20, -6, 15, 21, 29, 24, 27, -10,
		-3, -31, -2, 11, 16, 28, -5, -3,
		-82, -5, -39, -11, 1, -2, -2, -28,
	},

	// Piece-square table for bishops
	{
		-15, 14, -76, -21, -38, -24, 4, 12,
		-27, 17, -20, -10, -13, 59, 18, -44,
		-18, 22, 28, 30, 27, 64, 31, 0,
		-4, 5, 14, 53, 31, 35, 9, 3,
		4, 25, 13, 30, 37, 11, 16, 23,
		11, 24, 26, 15, 19, 38, 22, 19,
		21, 31, 22, 14, 18, 36, 43, 18,
		-19, 22, 12, 2, 13, 7, -1, -7,
	},

	// Piece square table for rook
	{
		4, 16, -4, 29, 19, 1, 15, -3,
		2, 6, 28, 28, 50, 38, 3, 7,
		-16, -1, -8, 9, -15, 31, 46, -15,
		-40, -22, -6, 8, 3, 20, -14, -34,
		-52, -28, -18, -14, -8, -29, -7, -45,
		-39, -35, -25, -16, -4, -10, -8, -40,
		-40, -21, -17, -1, 3, 8, -8, -57,
		-8, -3, 11, 17, 19, 17, -22, -11,
	},

	// Piece square table for queens
	{
		-22, -10, 25, 10, 44, 44, 47, 46,
		-32, -51, -7, 10, -12, 50, 50, 36,
		-23, -27, 2, -8, 4, 47, 41, 51,
		-41, -36, -27, -23, -2, -1, 8, -1,
		-8, -35, -10, -15, -10, -2, 2, -5,
		-13, 2, -6, 2, -1, 3, 14, -4,
		-24, -12, 12, 8, 17, 20, 7, 20,
		-1, 4, 15, 22, 1, -8, -13, -37,
	},

	// Piece square table for kings
	{
		-2, 64, 88, 43, -10, 26, 32, 30,
		59, 37, 39, 71, 40, 33, 5, -5,
		29, 42, 45, 33, 24, 64, 58, -14,
		-2, 17, 24, 21, 10, 23, 18, -8,
		-24, 36, -13, -39, -30, -18, -10, -38,
		32, 5, -3, -32, -23, -23, -8, -14,
		14, 23, -15, -62, -45, -19, 2, 8,
		-25, 35, 23, -61, 6, -24, 28, 10,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		18, 7, 0, 0, -1, 2, 1, 1,
		158, 158, 121, 104, 111, 107, 137, 174,
		70, 75, 56, 40, 23, 29, 64, 59,
		7, -7, -16, -29, -35, -24, -14, -9,
		-15, -24, -35, -40, -43, -42, -32, -30,
		-28, -26, -38, -28, -35, -38, -39, -40,
		-19, -26, -20, -28, -22, -35, -37, -42,
		12, 0, 2, 1, 2, 3, -1, 1,
	},

	// Piece-square table for knights
	{

		-53, -14, -1, -22, -13, -15, -52, -78,
		-1, 7, -11, 11, -9, -18, -12, -40,
		-9, -6, 14, 13, -2, -7, -13, -30,
		2, 11, 18, 20, 26, 13, 9, -11,
		-2, 5, 17, 32, 22, 17, 9, -6,
		-3, 6, -3, 16, 9, -2, -10, -4,
		-19, -12, -1, 1, 2, -16, -19, -24,
		-11, -33, -14, 0, -10, -18, -37, -49,
	},

	// Piece-square table for bishops
	{
		-10, -17, 2, 4, 12, 10, -4, -2,
		8, 5, 15, -4, 13, -3, 2, 11,
		16, 4, 8, 8, -2, 2, 5, 13,
		12, 16, 11, 3, 10, 7, 8, 11,
		2, 8, 15, 15, -1, 10, 0, 1,
		0, 10, 12, 12, 20, 3, 2, 2,
		1, -5, 5, 5, 9, 1, -2, -10,
		-6, 6, 1, 3, 7, 5, 9, -4,
	},

	// Piece square table for rook
	{
		20, 11, 20, 8, 14, 12, 8, 13,
		13, 18, 13, 13, -4, 5, 8, 13,
		14, 13, 14, 9, 15, -2, -1, 9,
		12, 11, 18, 4, 4, 5, 2, 11,
		14, 14, 13, 7, -5, 1, -2, 5,
		3, 11, 3, 0, -5, -7, -2, -2,
		-1, -1, 4, 4, -2, -4, -5, 2,
		-8, -1, -3, -7, -9, -13, -2, -18,
	},

	// Piece square table for queens
	{
		3, 27, 20, 22, 39, 23, 13, 42,
		-4, 19, 25, 35, 31, 22, -5, 34,
		-1, 6, -15, 41, 44, 30, 20, 35,
		32, 20, 18, 23, 33, 24, 41, 52,
		-12, 30, 5, 27, 15, 26, 32, 38,
		11, -19, 5, 0, 5, 19, 33, 25,
		-10, -10, -16, -4, -6, -17, -19, -13,
		-7, -28, -15, -12, 12, -6, 14, -21,
	},

	// Piece square table for kings
	{
		-66, -47, -20, -33, 7, 24, 7, -21,
		-6, 17, 11, 18, 22, 34, 25, 10,
		9, 23, 27, 17, 25, 39, 49, 20,
		-1, 22, 26, 33, 30, 33, 25, 9,
		-13, -8, 22, 36, 39, 26, 10, -2,
		-24, -2, 21, 27, 29, 23, 14, -4,
		-26, -11, 15, 25, 24, 14, 4, -13,
		-43, -32, -22, 0, -20, -2, -22, -41,
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
	var eval Eval
	phase := TotalPhase

	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	eval.CastlingSide[White] = Queenside
	eval.CastlingSide[Black] = Queenside

	if FileOf(pos.PieceBB[White][King].Msb()) > FileE {
		eval.CastlingSide[White] = Kingside
	}

	if FileOf(pos.PieceBB[Black][King].Msb()) > FileE {
		eval.CastlingSide[Black] = Kingside
	}

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
	file := FileOf(sq)

	// Evaluate isolated pawns.
	if IsolatedPawnMasks[file]&usPawns == 0 {
		eval.MGScores[color] -= IsolatedPawnPenatlyMG
		eval.EGScores[color] -= IsolatedPawnPenatlyEG
	}

	// Evaluate doubled pawns.
	if DoubledPawnMasks[color][sq]&usPawns != 0 {
		eval.MGScores[color] -= DoubledPawnPenatlyMG
		eval.EGScores[color] -= DoubledPawnPenaltyEG
	}

	// Evalue pawns part of the the pawn shield.
	castlingSide := eval.CastlingSide[color]
	if SquareBB[sq]&Zone1PawnShields[color][castlingSide] != 0 {
		eval.MGScores[color] += PawnShieldBonusZone1
	} else if SquareBB[sq]&Zone2PawnShields[color][castlingSide] != 0 {
		eval.MGScores[color] += PawnShieldBonusZone2
	}
}

// Evaluate the score of a knight.
func evalKnight(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Knight] + PSQT_MG[Knight][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Knight] + PSQT_EG[Knight][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	moves := KnightMoves[sq] & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 4) * PieceMobilityMG[Knight-1]
	eval.EGScores[color] += (mobility - 4) * PieceMobilityEG[Knight-1]
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
}

// Evaluate the score of a king.
func evalKing(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PSQT_MG[King][FlipSq[color][sq]]
	eval.EGScores[color] += PSQT_EG[King][FlipSq[color][sq]]
}

func init() {
	for sq := 0; sq < 64; sq++ {
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
	}
}
