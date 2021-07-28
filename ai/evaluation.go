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
//
// Credit to Marcel Vanthoor for the piece square table values, which can
// be found here:
//
// https://github.com/mvanthoor/rustic/blob/3f5436b01bdd40244b55f939e755569f9fc4070a/src/evaluation/psqt.rs
//
var PieceSquareTables [6][64]int = [6][64]int{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		60, 60, 60, 60, 70, 60, 60, 60,
		40, 40, 40, 50, 60, 40, 40, 40,
		20, 20, 20, 40, 50, 20, 20, 20,
		5, 5, 15, 30, 40, 10, 5, 5,
		5, 5, 10, 20, 30, 5, 5, 5,
		5, 5, 5, -30, -30, 5, 5, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, -5, -5, -5, -5, -5, -5, -10,
		-10, -5, 15, 15, 15, 15, -5, -10,
		-10, -5, 15, 15, 15, 15, -5, -10,
		-10, -5, 15, 15, 15, 15, -5, -10,
		-10, -5, 10, 15, 15, 15, -5, -10,
		-10, -5, -5, -5, -5, -5, -5, -10,
		-20, 0, -10, -10, -10, -10, 0, -20,
	},

	// Piece-square table for bishops
	{
		-20, 0, 0, 0, 0, 0, 0, -20,
		-15, 0, 0, 0, 0, 0, 0, -15,
		-10, 0, 0, 5, 5, 0, 0, -10,
		-10, 10, 10, 30, 30, 10, 10, -10,
		5, 5, 10, 25, 25, 10, 5, 5,
		5, 5, 5, 10, 10, 5, 5, 5,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	},

	// Piece square table for rook
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		15, 15, 15, 20, 20, 15, 15, 15,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 10, 10, 10, 0, 0,
	},

	// Piece square table for queens
	{
		-30, -20, -10, -10, -10, -10, -20, -30,
		-20, -10, -5, -5, -5, -5, -10, -20,
		-10, -5, 10, 10, 10, 10, -5, -10,
		-10, -5, 10, 20, 20, 10, -5, -10,
		-10, -5, 10, 20, 20, 10, -5, -10,
		-10, -5, -5, -5, -5, -5, -5, -10,
		-20, -10, -5, -5, -5, -5, -10, -20,
		-30, -20, -10, -10, -10, -10, -20, -30,
	},

	// Piece square table for kings
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 20, 20, 0, 0, 0,
		0, 0, 0, 20, 20, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, -10, -10, 0, 0, 0,
		0, 0, 20, -10, -10, 0, 20, 0,
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
