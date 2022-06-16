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
}

var BishopPairBonusMG int16 = 26
var BishopPairBonusEG int16 = 44

var PieceValueMG = [6]int16{96, 345, 356, 467, 963}
var PieceValueEG = [6]int16{141, 251, 273, 491, 922}
var PieceMobilityMG [5]int16 = [5]int16{0, 0, 4, 5, 0}
var PieceMobilityEG [5]int16 = [5]int16{0, 0, 3, 2, 6}

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
		71, 80, 50, 58, 42, 45, 31, 25,
		-25, -7, 10, 16, 44, 40, 7, -34,
		-37, -2, -11, 7, 8, -1, 6, -38,
		-47, -19, -14, 0, 5, 0, 0, -40,
		-39, -19, -16, -16, -1, -1, 21, -22,
		-45, -11, -33, -25, -17, 16, 30, -29,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-58, -14, -8, -11, 3, -26, -7, -28,
		-41, -28, 43, 11, 5, 30, -3, -14,
		-25, 37, 24, 45, 54, 58, 47, 10,
		-12, 11, 9, 44, 26, 58, 14, 13,
		-18, 0, 11, 6, 23, 14, 12, -13,
		-27, -14, 7, 8, 18, 17, 17, -18,
		-30, -31, -15, 2, 4, 15, -15, -20,
		-30, -16, -42, -29, -14, -26, -13, -17,
	},
	{
		// MG Bishop PST
		-13, -4, -21, -10, -5, -10, -1, -3,
		-18, 2, -19, -12, 8, 22, 4, -41,
		-20, 16, 22, 18, 12, 25, 16, -4,
		-9, -3, 6, 33, 22, 21, -1, -5,
		-7, 6, -3, 15, 19, -5, 1, 0,
		-2, 9, 10, 0, 5, 22, 11, 5,
		0, 19, 9, 0, 6, 16, 35, -1,
		-27, -3, 0, -17, -9, -2, -17, -21,
	},
	{
		// MG Rook PST
		12, 17, 9, 24, 21, 1, 4, 7,
		19, 21, 42, 36, 38, 34, 11, 15,
		-10, 9, 12, 15, 2, 16, 23, 5,
		-24, -13, -1, 11, 10, 16, -2, -16,
		-42, -21, -10, -5, 0, -7, 3, -25,
		-43, -23, -12, -15, 0, 0, -4, -29,
		-39, -13, -16, -5, 2, 9, -4, -59,
		-11, -9, 6, 15, 16, 9, -23, -7,
	},
	{
		// MG Queen PST
		-20, 0, 9, 5, 27, 13, 12, 20,
		-21, -39, 0, 6, 2, 29, 19, 27,
		-13, -14, 0, 5, 21, 35, 31, 42,
		-25, -26, -15, -18, -1, 12, 5, -1,
		-9, -26, -11, -16, -7, -3, 1, 0,
		-17, 1, -10, -1, -4, 1, 11, 2,
		-30, -9, 12, 7, 16, 12, -9, -6,
		-8, -9, 2, 19, -6, -25, -18, -32,
	},
	{
		// MG King PST
		-3, 1, 1, 0, -2, 0, 1, 0,
		2, 5, 3, 7, 4, 6, 2, -2,
		4, 10, 11, 4, 5, 16, 17, 1,
		0, 4, 6, 0, 0, 1, 4, -9,
		-9, 5, -5, -19, -24, -19, -17, -26,
		0, 1, -8, -33, -35, -33, 0, -16,
		3, 17, -3, -58, -40, -15, 28, 33,
		-13, 56, 32, -52, 15, -16, 52, 43,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		125, 121, 101, 87, 88, 86, 103, 121,
		51, 54, 38, 19, 5, 2, 35, 38,
		-10, -23, -33, -44, -52, -45, -34, -29,
		-28, -37, -51, -57, -58, -60, -48, -47,
		-40, -39, -54, -47, -51, -56, -54, -57,
		-31, -40, -38, -38, -38, -52, -52, -57,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-48, -27, -4, -20, -7, -25, -29, -42,
		-21, 0, -4, 13, 5, -5, -9, -31,
		-18, -1, 23, 21, 14, 18, -1, -20,
		-2, 16, 35, 31, 33, 22, 18, -3,
		-4, 7, 28, 39, 28, 29, 16, -5,
		-9, 10, 10, 25, 21, 6, -7, -9,
		-23, -12, 4, 5, 9, -7, -10, -24,
		-24, -35, -13, -1, -7, -7, -34, -25,
	},
	{
		// EG Bishop PST
		-12, -11, -14, -8, -2, -8, -5, -14,
		-4, 0, 7, -8, 2, 0, 0, -11,
		7, 0, 1, 0, 1, 8, 5, 7,
		2, 11, 10, 5, 7, 8, 3, 6,
		-1, 3, 11, 12, 1, 8, 0, -1,
		-4, 0, 8, 9, 13, 0, 0, -6,
		-6, -13, -2, 2, 3, -3, -13, -20,
		-17, -1, -9, 4, 2, -6, -5, -11,
	},
	{
		// EG Rook PST
		8, 5, 11, 8, 10, 10, 7, 5,
		5, 6, 4, 7, -2, 2, 6, 5,
		5, 3, 1, 2, 0, -2, -2, -3,
		4, 0, 8, -3, -2, 0, -4, 2,
		7, 3, 4, 0, -4, -6, -8, -8,
		2, 1, -5, 0, -8, -11, -9, -11,
		1, -5, 1, 3, -6, -8, -10, 4,
		-5, 1, 0, -4, -7, -7, 0, -21,
	},
	{
		// EG Queen PST
		-20, 4, 10, 9, 20, 10, 5, 11,
		-22, 1, 6, 15, 19, 17, 9, 4,
		-22, -11, -7, 20, 23, 19, 10, 8,
		-6, 4, -1, 15, 25, 15, 25, 19,
		-24, 7, 2, 25, 12, 12, 18, 7,
		-10, -36, 0, -6, -2, 6, 2, 1,
		-17, -19, -38, -22, -21, -16, -19, -11,
		-20, -25, -21, -49, -6, -22, -16, -25,
	},
	{
		// EG King PST
		-18, -12, -8, -10, -8, 6, 4, -4,
		-5, 13, 8, 10, 10, 30, 17, 4,
		4, 16, 18, 10, 13, 38, 38, 8,
		-12, 15, 19, 22, 20, 27, 21, -1,
		-24, -7, 16, 22, 25, 19, 7, -14,
		-21, -6, 7, 20, 23, 17, 3, -9,
		-28, -16, 2, 15, 15, 5, -12, -29,
		-50, -47, -29, -7, -30, -12, -39, -62,
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
func evaluatePos(pos *Position) int16 {
	eval := Eval{MGScores: pos.MGScores, EGScores: pos.EGScores}
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
		eval.MGScores[White] += BishopPairBonusEG
	}

	if pos.Pieces[Black][Bishop].CountBits() >= 2 {
		eval.MGScores[Black] += BishopPairBonusMG
		eval.MGScores[Black] += BishopPairBonusEG
	}

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
}

// Evaluate the score of a bishop.
func evalBishop(pos *Position, color, sq uint8, eval *Eval) {
	usBB := pos.Sides[color]
	allBB := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]

	moves := GenBishopMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 7) * PieceMobilityMG[Bishop]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Bishop]
}

// Evaluate the score of a rook.
func evalRook(pos *Position, color, sq uint8, eval *Eval) {
	usBB := pos.Sides[color]
	allBB := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]

	moves := GenRookMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 7) * PieceMobilityMG[Rook]
	eval.EGScores[color] += (mobility - 7) * PieceMobilityEG[Rook]
}

// Evaluate the score of a queen.
func evalQueen(pos *Position, color, sq uint8, eval *Eval) {
	usBB := pos.Sides[color]
	allBB := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]

	moves := (GenBishopMoves(sq, allBB) | GenRookMoves(sq, allBB)) & ^usBB
	mobility := int16(moves.CountBits())

	eval.MGScores[color] += (mobility - 14) * PieceMobilityMG[Queen]
	eval.EGScores[color] += (mobility - 14) * PieceMobilityEG[Queen]
}
