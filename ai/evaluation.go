package ai

import (
	"blunder/engine"
)

const (
	PawnValue   = 100
	KnightValue = 300
	BishopValue = 330
	RookValue   = 500
	QueenValue  = 900

	PawnPhase   = 0
	KnightPhase = 1
	BishopPhase = 1
	RookPhase   = 2
	QueenPhase  = 4
	TotalPhase  = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

	PosInf = 200000
	NegInf = -PosInf
)

// Below are tables containg piece square tables for each piece, indexed
// by their bitboard index (see the constants in board.go)
//
// Credit to Maksym Korzh for the piece square table values, which can
// be found here:
//
// https://github.com/maksimKorzh/wukongJS/blob/main/tools/eval_tuner/temp_weights/session_weights_2021-01-08-18-52.txt

var PSQT_MG [6][64]int = [6][64]int{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		49, 49, 50, 50, 51, 49, 49, 49,
		11, 11, 19, 31, 29, 21, 10, 9,
		5, 4, 9, 25, 25, 11, 4, 6,
		1, 0, 1, 19, 19, -1, -1, -1,
		4, -6, -10, 0, 0, -10, -4, 4,
		5, 11, 11, -19, -19, 11, 11, 4,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		-51, -40, -30, -29, -29, -30, -40, -50,
		-40, -19, 0, -1, 0, -1, -20, -40,
		-29, 1, 11, 14, 15, 9, 0, -29,
		-30, 6, 16, 19, 19, 15, 6, -31,
		-31, 0, 15, 19, 19, 16, -1, -29,
		-29, 4, 9, 16, 14, 11, 6, -29,
		-39, -19, 1, 4, 6, 0, -19, -39,
		-49, -39, -30, -29, -31, -29, -39, -51,
	},

	// Piece-square table for bishops
	{
		-19, -9, -9, -9, -10, -9, -11, -20,
		-10, 1, 1, 1, -1, -1, 1, -9,
		-11, 1, 4, 11, 9, 5, -1, -10,
		-9, 4, 4, 11, 9, 4, 6, -11,
		-9, 1, 11, 11, 9, 9, 1, -11,
		-9, 9, 11, 9, 10, 11, 11, -10,
		-11, 4, 0, 1, 1, 1, 6, -9,
		-19, -9, -11, -9, -11, -10, -10, -19,
	},

	// Piece square table for rook
	{
		0, 0, 1, -1, 1, 1, -1, -1,
		6, 9, 11, 10, 11, 11, 11, 5,
		-4, 0, 1, 1, -1, 1, -1, -4,
		-6, 0, 1, -1, 1, 1, 1, -4,
		-4, 1, 1, -1, 1, -1, -1, -6,
		-5, 1, 0, -1, 1, 0, -1, -5,
		-6, 0, -1, 0, 1, -1, -1, -6,
		0, 1, 1, 4, 6, 1, -1, 0,
	},

	// Piece square table for queens
	{
		-21, -10, -9, -6, -5, -10, -9, -19,
		-11, 1, 1, -1, 1, -1, 0, -11,
		-11, -1, 4, 5, 6, 5, 1, -9,
		-4, 0, 5, 4, 5, 6, 0, -4,
		-1, 1, 6, 4, 4, 6, -1, -4,
		-11, 4, 4, 6, 6, 6, 1, -11,
		-9, 1, 6, -1, 1, -1, 1, -9,
		-19, -9, -11, -4, -6, -11, -10, -19,
	},

	// Piece square table for kings
	{
		-30, -40, -40, -50, -50, -40, -39, -30,
		-30, -39, -39, -50, -51, -39, -40, -30,
		-30, -39, -39, -50, -49, -39, -41, -31,
		-31, -41, -41, -49, -49, -40, -40, -30,
		-20, -31, -30, -40, -39, -29, -30, -20,
		-10, -21, -20, -19, -20, -20, -20, -10,
		19, 21, 0, 0, 1, 1, 19, 19,
		19, 31, 9, 1, 0, 11, 30, 19,
	},
}

var PSQT_EG [6][64]int = [6][64]int{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		-1, 50, 49, 49, 50, 50, 49, 1,
		-1, 40, 41, 40, 40, 39, 39, 1,
		-1, -1, -1, -1, 1, -1, 0, 1,
		0, -1, 1, -1, 0, -1, -1, 0,
		-1, 1, -1, 1, 0, 1, 0, 0,
		0, 0, 1, 1, -1, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		-51, -39, -29, -29, -30, -31, -41, -50,
		-39, -19, 1, -1, 1, -1, -19, -41,
		-30, 0, 11, 14, 15, 10, -1, -30,
		-30, 4, 15, 21, 20, 14, 6, -29,
		-31, -1, 16, 19, 19, 16, 0, -29,
		-29, 4, 9, 16, 14, 11, 4, -29,
		-40, -20, 0, 4, 5, -1, -21, -40,
		-50, -40, -31, -31, -31, -30, -40, -50,
	},

	// Piece-square table for bishops
	{
		-19, -9, -9, -10, -11, -9, -11, -19,
		-11, 0, -1, -1, 0, -1, 1, -9,
		-10, 0, 6, 10, 9, 6, 1, -9,
		-10, 4, 5, 9, 9, 6, 6, -11,
		-10, 1, 10, 11, 9, 10, 1, -9,
		-9, 9, 11, 10, 9, 9, 10, -11,
		-10, 5, 0, 1, 0, 1, 4, -10,
		-20, -9, -11, -9, -10, -9, -11, -21,
	},

	// Piece square table for rook
	{
		-1, 1, 1, -1, -1, 0, 1, 0,
		-9, 0, 1, 0, 1, -1, 1, -11,
		-10, 0, 0, 0, -1, 1, 1, -9,
		-10, 1, 0, -1, 0, -1, 1, -9,
		-11, 0, 0, -1, 1, 0, -1, -9,
		-9, 0, 0, -1, 1, 0, -1, -10,
		-11, 0, 1, 0, 0, -1, -1, -10,
		1, 1, 0, 0, 0, 0, -1, 0,
	},

	// Piece square table for queens
	{
		-21, -11, -9, -4, -5, -11, -9, -19,
		-11, 1, 1, 1, 1, 1, 1, -10,
		-11, -1, 6, 4, 6, 5, 1, -10,
		-4, 1, 4, 4, 5, 5, 0, -4,
		-1, 0, 5, 5, 5, 5, 0, -4,
		-11, 4, 4, 5, 5, 6, 1, -11,
		-10, 1, 6, 0, 0, -1, 1, -10,
		-19, -9, -11, -5, -5, -9, -10, -19,
	},

	// Piece square table for kings
	{
		-50, -41, -30, -19, -19, -30, -39, -51,
		-30, -20, -11, 0, 1, -9, -19, -29,
		-31, -11, 21, 30, 30, 19, -11, -30,
		-31, -10, 30, 41, 39, 31, -10, -30,
		-29, -11, 30, 39, 39, 29, -9, -29,
		-30, -10, 21, 31, 31, 20, -11, -31,
		-30, -29, 1, 0, 1, 0, -30, -30,
		-51, -29, -29, -30, -30, -31, -31, -51,
	},
}

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

// Evaluate a board and give a score, from the perspective of the side to move (
// more positive if it's good for the side to move, otherwise more negative).
func EvaluateBoard(board *engine.Board) int {
	whiteScore := evaluateMaterial(board, engine.White)
	whiteScore += evaluatePosition(board, engine.White)

	blackScore := evaluateMaterial(board, engine.Black)
	blackScore += evaluatePosition(board, engine.Black)

	if board.ColorToMove == engine.White {
		return whiteScore - blackScore
	}
	return blackScore - whiteScore
}

// Evaluate the material a side has.
func evaluateMaterial(board *engine.Board, usColor int) (score int) {
	score += engine.CountBits(board.PieceBB[usColor][engine.Pawn]) * PawnValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Knight]) * KnightValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Bishop]) * BishopValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Rook]) * RookValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Queen]) * QueenValue
	return score
}

// Evaluate a board position using piece-square tables.
func evaluatePosition(board *engine.Board, usColor int) int {
	usBB := board.SideBB[usColor]
	var opening, endgame int

	for usBB != 0 {
		sq := engine.PopBit(&usBB)
		pieceType := board.Squares[sq].Type

		opening += PSQT_MG[pieceType][FlipSq[usColor][sq]]
		endgame += PSQT_EG[pieceType][FlipSq[usColor][sq]]
	}

	phase := calcGamePhase(board)
	return ((opening * (256 - phase)) + (endgame * phase)) / 256
}

// Calculate the current phase of the game
func calcGamePhase(board *engine.Board) int {
	phase := TotalPhase

	phase -= PawnPhase * engine.CountBits(board.PieceBB[engine.White][engine.Pawn])
	phase -= KnightPhase * engine.CountBits(board.PieceBB[engine.White][engine.Knight])
	phase -= BishopPhase * engine.CountBits(board.PieceBB[engine.White][engine.Bishop])
	phase -= RookPhase * engine.CountBits(board.PieceBB[engine.White][engine.Rook])
	phase -= QueenPhase * engine.CountBits(board.PieceBB[engine.White][engine.Queen])

	phase -= PawnPhase * engine.CountBits(board.PieceBB[engine.Black][engine.Pawn])
	phase -= KnightPhase * engine.CountBits(board.PieceBB[engine.Black][engine.Knight])
	phase -= BishopPhase * engine.CountBits(board.PieceBB[engine.Black][engine.Bishop])
	phase -= RookPhase * engine.CountBits(board.PieceBB[engine.Black][engine.Rook])
	phase -= QueenPhase * engine.CountBits(board.PieceBB[engine.Black][engine.Queen])

	return (phase*256 + (TotalPhase / 2)) / TotalPhase
}
