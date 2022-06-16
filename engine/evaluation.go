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

var BishopPairBonusMG int16 = 27
var BishopPairBonusEG int16 = 44

var PieceValueMG = [6]int16{90, 344, 347, 461, 959}
var PieceValueEG = [6]int16{139, 247, 270, 498, 924}

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
		70, 75, 49, 55, 40, 42, 31, 27,
		-22, -4, 12, 18, 43, 39, 9, -31,
		-34, -1, -5, 10, 12, 1, 8, -35,
		-47, -15, -16, 1, 7, -1, 2, -37,
		-44, -17, -15, -20, -6, -3, 25, -24,
		-54, -13, -31, -33, -25, 16, 29, -33,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-51, -13, -7, -10, 3, -23, -7, -25,
		-37, -25, 40, 11, 5, 28, -3, -13,
		-23, 36, 25, 46, 51, 53, 45, 9,
		-12, 13, 14, 48, 33, 57, 15, 13,
		-16, 0, 9, 8, 23, 15, 12, -11,
		-25, -12, 9, 5, 14, 15, 20, -18,
		-30, -30, -15, -5, -4, 14, -15, -22,
		-27, -22, -42, -32, -18, -28, -19, -15,
	},
	{
		// MG Bishop PST
		-13, -4, -19, -9, -5, -9, -1, -4,
		-18, 3, -16, -10, 8, 21, 5, -39,
		-19, 16, 24, 21, 13, 25, 17, -4,
		-9, 0, 11, 38, 30, 25, 2, -7,
		-9, 7, 6, 19, 27, 6, 4, -3,
		-5, 9, 8, 10, 7, 22, 11, 2,
		-2, 11, 9, -5, 2, 13, 27, -4,
		-30, -7, -18, -25, -17, -16, -17, -23,
	},
	{
		// MG Rook PST
		13, 18, 12, 25, 22, 2, 5, 8,
		21, 25, 42, 36, 36, 32, 12, 15,
		-6, 12, 15, 19, 6, 16, 22, 6,
		-21, -8, 4, 16, 16, 17, -1, -13,
		-37, -18, -8, -1, 2, -5, 4, -20,
		-45, -20, -12, -16, 0, -1, -3, -29,
		-46, -13, -18, -9, -1, 7, -6, -65,
		-20, -15, -1, 9, 10, -1, -29, -19,
	},
	{
		// MG Queen PST
		-22, 1, 9, 6, 25, 11, 10, 17,
		-23, -35, 2, 8, 4, 28, 19, 25,
		-15, -11, 3, 11, 24, 35, 31, 40,
		-24, -22, -9, -9, 7, 18, 8, 0,
		-8, -21, -7, -8, -1, 1, 5, 0,
		-17, 3, -8, -2, -4, 2, 12, 1,
		-35, -12, 12, 1, 8, 8, -11, -8,
		-13, -19, -10, 12, -17, -33, -19, -30,
	},
	{
		// MG King PST
		-3, 1, 1, 0, -2, 0, 1, 0,
		2, 5, 3, 6, 4, 6, 2, -2,
		4, 9, 10, 4, 5, 15, 16, 1,
		-1, 4, 6, 0, 0, 2, 4, -8,
		-9, 4, -4, -16, -20, -16, -15, -23,
		-1, 1, -7, -30, -31, -30, 0, -16,
		2, 15, -2, -59, -41, -13, 30, 34,
		-15, 55, 27, -54, 17, -24, 48, 43,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		114, 111, 95, 84, 83, 81, 95, 110,
		50, 53, 38, 19, 7, 3, 35, 38,
		-9, -22, -34, -44, -52, -45, -33, -28,
		-27, -37, -50, -56, -57, -58, -47, -47,
		-38, -38, -54, -44, -48, -54, -53, -55,
		-27, -37, -37, -35, -35, -50, -50, -53,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-45, -25, -4, -20, -7, -24, -26, -37,
		-22, -1, -3, 13, 4, -5, -10, -30,
		-19, 0, 24, 23, 16, 20, 0, -20,
		-3, 16, 35, 31, 32, 23, 17, -3,
		-5, 6, 28, 38, 27, 29, 16, -6,
		-10, 10, 9, 26, 21, 7, -8, -10,
		-23, -14, 2, 7, 10, -8, -11, -24,
		-23, -39, -17, -6, -10, -10, -36, -23,
	},
	{
		// EG Bishop PST
		-14, -12, -15, -9, -4, -9, -6, -14,
		-8, 0, 7, -8, 3, 0, 1, -14,
		4, 0, 7, 5, 5, 13, 6, 5,
		0, 12, 16, 14, 16, 14, 5, 4,
		-4, 5, 16, 20, 10, 12, 0, -4,
		-8, 0, 11, 13, 16, 4, -2, -9,
		-9, -15, -3, 3, 6, -5, -14, -22,
		-23, -6, -23, -3, -5, -15, -8, -14,
	},
	{
		// EG Rook PST
		13, 11, 18, 16, 16, 11, 9, 8,
		7, 9, 9, 12, 2, 5, 7, 5,
		5, 5, 5, 6, 3, -1, -1, -4,
		2, 1, 9, 0, 0, 0, -5, 0,
		3, 1, 4, 0, -5, -8, -10, -12,
		-2, -1, -8, -2, -10, -14, -11, -16,
		-3, -8, -1, 0, -10, -11, -12, -1,
		-8, 1, 1, -1, -6, -8, 0, -25,
	},
	{
		// EG Queen PST
		-22, 3, 10, 9, 20, 9, 4, 9,
		-25, 0, 7, 15, 19, 17, 9, 3,
		-23, -10, -3, 22, 25, 20, 10, 7,
		-9, 5, 2, 20, 29, 18, 23, 16,
		-24, 7, 6, 30, 17, 14, 17, 5,
		-11, -34, 2, -1, 0, 7, 1, 0,
		-20, -20, -38, -21, -20, -17, -19, -12,
		-22, -26, -22, -57, -10, -26, -17, -24,
	},
	{
		// EG King PST
		-15, -10, -7, -9, -7, 6, 4, -4,
		-5, 13, 8, 10, 10, 29, 16, 4,
		5, 16, 18, 10, 12, 38, 38, 8,
		-12, 15, 18, 21, 19, 27, 21, -1,
		-24, -7, 16, 21, 24, 19, 7, -14,
		-21, -6, 7, 19, 22, 17, 3, -9,
		-27, -15, 2, 15, 16, 5, -12, -28,
		-48, -46, -29, -8, -34, -12, -37, -60,
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
