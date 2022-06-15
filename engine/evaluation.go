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

var BishopPairBonusMG int16 = 23
var BishopPairBonusEG int16 = 45

var PieceValueMG = [6]int16{95, 366, 373, 506, 1116}
var PieceValueEG = [6]int16{133, 276, 292, 490, 864}
var PieceMobilityMG [5]int16 = [5]int16{0, 1, 5, 7, 0}
var PieceMobilityEG [5]int16 = [5]int16{0, 0, 4, 3, 6}

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
		64, 121, 61, 102, 84, 133, 10, -56,
		-25, -11, 9, 14, 53, 50, 5, -36,
		-33, 1, -9, 10, 13, 0, 11, -37,
		-42, -15, -9, 3, 11, 6, 5, -38,
		-33, -14, -12, -9, 3, 4, 25, -18,
		-40, -4, -26, -17, -12, 20, 37, -27,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-193, -113, -63, -60, 51, -105, -22, -122,
		-95, -60, 61, 14, 11, 54, -10, -30,
		-55, 47, 26, 53, 81, 119, 66, 37,
		-16, 15, 8, 47, 30, 65, 17, 17,
		-17, 3, 16, 11, 27, 17, 15, -12,
		-26, -11, 12, 15, 25, 22, 20, -18,
		-28, -49, -13, 6, 9, 20, -19, -16,
		-124, -15, -49, -28, -8, -23, -11, -33,
	},
	{
		// MG Bishop PST
		-31, -9, -121, -77, -53, -50, -29, -3,
		-31, 0, -33, -44, 12, 42, 5, -51,
		-21, 25, 25, 20, 19, 44, 21, -1,
		-7, 3, 3, 37, 25, 20, 4, 0,
		-1, 12, 3, 21, 25, 0, 5, 11,
		5, 17, 18, 8, 13, 31, 17, 13,
		11, 27, 17, 8, 17, 26, 42, 8,
		-19, 9, 9, -1, 4, 3, -29, -12,
	},
	{
		// MG Rook PST
		14, 25, -3, 41, 35, -28, -10, 7,
		10, 8, 50, 54, 75, 68, 2, 23,
		-31, -4, 0, 3, -19, 32, 52, -4,
		-40, -31, -17, -1, -6, 15, -20, -36,
		-49, -36, -19, -16, -4, -19, -1, -36,
		-43, -27, -14, -14, 0, -1, -9, -31,
		-35, -15, -13, -1, 5, 14, 0, -55,
		-9, -5, 10, 19, 20, 13, -20, -6,
	},
	{
		// MG Queen PST
		-31, -38, -19, -19, 84, 87, 54, 41,
		-36, -58, -23, -28, -54, 43, 11, 39,
		-19, -35, -8, -28, 5, 50, 32, 47,
		-43, -36, -34, -33, -20, -6, -18, -11,
		-11, -43, -14, -22, -14, -12, -11, -10,
		-20, 0, -12, -1, -5, -1, 6, -1,
		-26, -3, 11, 11, 20, 21, 5, 12,
		7, 7, 14, 23, 4, -12, -7, -46,
	},
	{
		// MG King PST
		-67, 175, 158, 104, -107, -55, 64, 43,
		196, 63, 52, 114, 45, 41, -13, -126,
		62, 69, 91, 25, 52, 118, 116, 0,
		0, -3, 18, -29, -26, -36, -8, -74,
		-69, 17, -53, -98, -101, -69, -63, -79,
		5, -12, -42, -74, -73, -67, -21, -39,
		8, 8, -26, -75, -57, -35, 6, 15,
		-16, 42, 18, -71, -2, -33, 35, 25,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		143, 128, 114, 86, 97, 81, 131, 162,
		59, 63, 46, 27, 10, 9, 42, 47,
		-4, -17, -27, -38, -43, -39, -25, -23,
		-22, -31, -46, -51, -50, -51, -40, -41,
		-33, -34, -48, -40, -45, -49, -49, -49,
		-25, -32, -30, -34, -34, -47, -44, -51,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-49, -41, -15, -38, -47, -33, -76, -98,
		-30, -14, -36, -10, -18, -38, -33, -59,
		-35, -31, 2, -2, -19, -24, -34, -57,
		-27, -7, 15, 10, 11, -1, -5, -33,
		-31, -16, 5, 16, 5, 6, -7, -32,
		-37, -12, -12, 1, -2, -17, -30, -37,
		-53, -30, -19, -20, -17, -31, -34, -59,
		-29, -62, -32, -26, -35, -32, -65, -76,
	},
	{
		// EG Bishop PST
		-25, -31, -10, -13, -11, -19, -20, -35,
		-16, -22, -11, -23, -23, -31, -22, -21,
		-11, -26, -25, -26, -29, -24, -19, -12,
		-18, -14, -13, -21, -18, -18, -21, -14,
		-22, -21, -14, -14, -25, -16, -24, -24,
		-27, -23, -18, -15, -13, -27, -24, -30,
		-27, -37, -27, -22, -20, -29, -37, -42,
		-31, -21, -32, -20, -21, -28, -15, -26,
	},
	{
		// EG Rook PST
		5, 0, 13, 0, 2, 15, 9, 2,
		6, 8, -1, -2, -18, -11, 6, -1,
		8, 6, 3, 4, 5, -8, -15, -3,
		6, 5, 12, 0, 2, -2, 0, 7,
		8, 6, 5, 3, -4, -4, -9, -4,
		0, 1, -6, -1, -9, -12, -7, -12,
		-3, -6, -2, 0, -12, -11, -12, 1,
		-10, -2, -5, -9, -12, -13, -2, -27,
	},
	{
		// EG Queen PST
		-10, 32, 24, 21, -21, -27, -17, 10,
		-6, 11, 16, 33, 56, 0, 15, 2,
		-21, -3, -25, 34, 19, -4, 5, 6,
		16, 1, -1, 7, 23, 15, 55, 43,
		-28, 19, -15, 13, -2, 8, 32, 31,
		-4, -48, -12, -27, -16, 0, 12, 19,
		-15, -32, -46, -36, -38, -32, -44, -30,
		-36, -52, -43, -45, -16, -22, -22, -26,
	},
	{
		// EG King PST
		-63, -61, -44, -37, 9, 25, -4, -18,
		-48, 5, 3, -4, 7, 30, 26, 35,
		-3, 8, 7, 11, 9, 26, 24, 11,
		-11, 20, 20, 31, 30, 37, 29, 13,
		-10, -7, 28, 42, 45, 35, 20, 1,
		-22, -1, 20, 32, 36, 28, 12, -3,
		-30, -12, 12, 23, 23, 15, -7, -20,
		-55, -42, -23, 0, -25, -6, -34, -54,
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
