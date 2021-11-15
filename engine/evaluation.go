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
}

var PieceValueMG [5]int16 = [5]int16{89, 309, 339, 458, 940}
var PieceValueEG [5]int16 = [5]int16{130, 270, 289, 501, 929}
var PieceMobilityMG [4]int16 = [4]int16{3, 4, 6, 1}
var PieceMobilityEG [4]int16 = [4]int16{3, 4, 2, 8}

var PhaseValues [6]int16 = [6]int16{
	PawnPhase,
	KnightPhase,
	BishopPhase,
	RookPhase,
	QueenPhase,
}

// Endgame and middlegame piece square tables, with piece values builtin.
//
// https://www.chessprogramming.org/Piece-Square_Tables
//
var PSQT_MG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		15, -6, 1, 0, -1, 4, 0, -1,
		111, 131, 99, 97, 104, 97, 90, -7,
		3, 9, 25, 28, 54, 38, -15, -17,
		-26, 4, -4, 15, 14, 2, 10, -37,
		-34, -12, -10, 3, 9, -4, -1, -36,
		-24, -12, -13, -12, -3, -4, 22, -18,
		-30, -5, -28, -19, -13, 17, 29, -24,
		15, 1, 0, 0, -2, 0, 0, -1,
	},

	// Piece-square table for knights
	{
		-132, -66, 9, -22, 31, -90, -31, -93,
		-64, -52, 62, 16, 26, 42, -2, -9,
		-42, 35, 25, 48, 68, 89, 60, 27,
		-3, 15, 12, 42, 26, 61, 19, 32,
		-14, 2, 12, 8, 29, 17, 31, 0,
		-21, -9, 10, 15, 33, 21, 23, -9,
		2, -26, -4, 10, 13, 26, -8, -3,
		-60, -6, -29, -14, -3, 5, -2, -20,
	},

	// Piece-square table for bishops
	{
		-10, 10, -71, -23, -39, -26, 9, 0,
		-33, 14, -8, -1, -12, 50, 10, -36,
		-26, 25, 25, 25, 32, 57, 33, -1,
		-3, 9, 24, 47, 37, 34, 13, 2,
		5, 15, 6, 24, 28, 7, 13, 25,
		8, 13, 22, 9, 15, 30, 19, 10,
		15, 27, 19, 14, 14, 28, 38, 14,
		-14, 14, 11, 5, 10, 3, 6, -5,
	},

	// Piece square table for rook
	{
		-2, 16, -2, 29, 22, 2, 15, -6,
		2, 2, 26, 28, 44, 35, 6, 9,
		-15, -5, -9, 9, -17, 25, 39, -21,
		-33, -26, -5, 17, 0, 24, -17, -31,
		-47, -30, -17, -10, -2, -20, -4, -42,
		-35, -29, -25, -19, 3, -12, -11, -39,
		-37, -20, -14, -1, 2, 5, -6, -55,
		-9, -5, 9, 20, 17, 17, -27, -14,
	},

	// Piece square table for queens
	{
		-18, -16, 25, 13, 37, 35, 40, 29,
		-28, -44, -10, 6, 4, 49, 53, 29,
		-22, -16, 3, -2, 10, 46, 45, 44,
		-36, -34, -22, -18, -8, -4, 15, -11,
		-11, -41, -11, -18, -5, -2, 4, 1,
		-20, 2, -6, -2, -7, -1, 10, 7,
		-17, -5, 13, 7, 17, 15, 5, 20,
		4, -1, 14, 21, 1, -14, -21, -44,
	},

	// Piece square table for kings
	{
		13, 59, 78, 58, -8, 22, 28, 22,
		51, 25, 35, 58, 37, 39, 15, 17,
		25, 38, 41, 26, 17, 57, 49, -14,
		-11, 10, 22, 24, 13, 22, 16, -5,
		-29, 37, -9, -31, -8, -14, -4, -37,
		26, 4, 0, -22, -18, -21, -7, -16,
		7, 17, -17, -57, -42, -19, 0, 7,
		-41, 32, 22, -55, 5, -25, 28, 6,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		14, 7, 1, 0, -1, 1, 2, 1,
		145, 135, 99, 96, 111, 98, 130, 157,
		60, 69, 51, 33, 17, 23, 57, 50,
		4, -8, -18, -29, -36, -31, -18, -11,
		-16, -22, -35, -41, -41, -43, -31, -31,
		-29, -25, -37, -33, -34, -38, -38, -40,
		-22, -25, -23, -26, -26, -35, -34, -40,
		10, -1, 2, 0, 1, 3, 0, 1,
	},

	// Piece-square table for knights
	{

		-65, -19, 2, -14, -6, -18, -52, -94,
		-13, 7, -22, 13, -5, -25, -9, -42,
		-10, -9, 12, 14, -1, -4, -8, -24,
		1, 14, 23, 24, 24, 10, 12, -9,
		-5, 2, 18, 34, 14, 21, 11, -12,
		-3, 7, -5, 16, 4, -7, -15, -6,
		-31, -3, -5, -4, 1, -23, -17, -36,
		-17, -36, -8, -2, -14, -27, -34, -37,
	},

	// Piece-square table for bishops
	{
		-16, -17, 5, 9, 9, 4, 7, -9,
		9, 5, 11, -7, 11, -11, 0, -12,
		18, 6, 5, 2, -5, 2, 6, 14,
		7, 14, 10, 2, 10, 8, 8, 13,
		1, 6, 16, 16, -2, 12, -4, 2,
		1, 8, 10, 14, 19, 3, 1, 1,
		-2, -5, -2, 4, 9, 3, 0, -12,
		-8, 8, -2, 5, 4, 5, -7, 2,
	},

	// Piece square table for rook
	{
		14, 9, 20, 11, 14, 10, 7, 12,
		11, 17, 16, 18, -3, 7, 9, 16,
		9, 13, 14, 5, 12, 0, -5, 10,
		14, 6, 12, 1, 6, 1, -1, 7,
		15, 8, 12, 5, -4, 4, -5, 1,
		0, 7, 3, -3, -9, -1, -1, -4,
		-3, 1, 6, 3, -6, -4, -2, 3,
		-8, -1, -3, -6, -10, -12, -2, -20,
	},

	// Piece square table for queens
	{
		4, 24, 19, 17, 38, 12, 16, 33,
		-3, 13, 17, 30, 35, 19, -11, 29,
		-1, 2, -12, 43, 42, 36, 26, 40,
		25, 19, 11, 18, 23, 20, 34, 53,
		-6, 26, 3, 27, 14, 25, 31, 39,
		11, -21, 0, -1, 0, 17, 23, 25,
		-11, -9, -17, 0, -7, -20, -16, -23,
		-5, -23, -15, -9, 15, -3, 6, -24,
	},

	// Piece square table for kings
	{
		-37, -35, -26, -5, 7, 5, 0, -27,
		0, 19, 21, 15, 20, 37, 24, 10,
		9, 15, 22, 20, 17, 45, 42, 13,
		-9, 24, 26, 28, 31, 32, 30, 1,
		-15, -6, 23, 33, 30, 28, 15, -7,
		-25, -1, 15, 26, 25, 22, 11, -9,
		-30, -13, 13, 20, 21, 13, 4, -13,
		-49, -30, -22, 1, -19, -6, -20, -38,
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

// Evaluate the score of a knight.
func evalPawn(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PieceValueMG[Pawn] + PSQT_MG[Pawn][FlipSq[color][sq]]
	eval.EGScores[color] += PieceValueEG[Pawn] + PSQT_EG[Pawn][FlipSq[color][sq]]
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

// Evaluate the score of a queen.
func evalKing(pos *Position, color, sq uint8, eval *Eval) {
	eval.MGScores[color] += PSQT_MG[King][FlipSq[color][sq]]
	eval.EGScores[color] += PSQT_EG[King][FlipSq[color][sq]]
}
