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

	PosInf = 200000
	NegInf = -PosInf
)

// Table containg piece square tables for each piece, indexed
// by their bitboard index (see the constants in board.go)
var PieceSquareTables [6][64]int = [6][64]int{

	// Piece-square table for pawns
	{
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		-5, -5, 0, 25, 25, 0, -5, -5,
		-5, -5, 5, 10, 10, 5, -5, -5,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 15, 15, 15, 15, 0, 0,
		0, 0, 15, 0, 0, 15, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for bishops
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 15, 10, 10, 15, 0, 0,
		0, 0, 15, 10, 10, 15, 0, 0,
		0, 15, 0, 0, 0, 0, 15, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for rook
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 50, 50, 0, 0, 0,
		0, 0, 0, 50, 50, 0, 0, 0,
	},

	// Piece square table for queens
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 20, 10, 10, 20, 0, 0,
		0, 0, 20, 15, 15, 20, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for kings
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, -25, -25, 0, 0, 0,
		0, 0, 0, -25, -25, 0, 0, 0,
		50, 50, 0, 0, 0, 0, 50, 50,
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

func evaluateBoard(board *engine.Board) int {
	whiteScore := evaluateMaterial(board, engine.White)
	whiteScore += evaluatePosition(board, engine.White)

	blackScore := evaluateMaterial(board, engine.Black)
	blackScore += evaluatePosition(board, engine.Black)

	if board.ColorToMove == engine.White {
		return whiteScore - blackScore
	}
	return blackScore - whiteScore
}

func evaluateMaterial(board *engine.Board, usColor int) (score int) {
	score += engine.CountBits(board.PieceBB[usColor][engine.Pawn]) * PawnValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Knight]) * KnightValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Bishop]) * BishopValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Rook]) * RookValue
	score += engine.CountBits(board.PieceBB[usColor][engine.Queen]) * QueenValue
	return score
}

// Evaluate a board position using piece-square tables.
func evaluatePosition(board *engine.Board, usColor int) (score int) {
	usBB := board.SideBB[usColor]
	for usBB != 0 {
		sq := engine.PopBit(&usBB)
		pieceType := board.Squares[sq].Type
		score += PieceSquareTables[pieceType][FlipSq[usColor][sq]]
	}
	return score
}
