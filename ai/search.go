package ai

import (
	"blunder/engine"
	"fmt"
)

const (
	SearchDepth             = 8
	NullMove    engine.Move = 0
)

// A table used to implement Most-Valuable-Victim,
// Least-Valuabe-Attacker move scoring.
var MvvLva [7][6]int = [7][6]int{
	{16, 15, 14, 13, 12, 11}, // victim Pawn
	{26, 25, 24, 23, 22, 21}, // victim Knight
	{36, 35, 34, 33, 32, 31}, // victim Bishop
	{46, 45, 44, 43, 42, 41}, // vitcim Rook
	{56, 55, 54, 53, 52, 51}, // victim Queen

	{0, 0, 0, 0, 0, 0}, // victim King
	{0, 0, 0, 0, 0, 0}, // No piece
}

// A struct to hold state during a search
type Search struct {
	Board         engine.Board
	Timer         Timer
	nodesSearched uint64
	searchOver    bool
}

// Search for the best move for the side to move in the given position.
// Implemented using iterative deepening.
func (search *Search) Search() engine.Move {
	bestMove, bestScore := NullMove, NegInf
	search.Timer.StartSearch()

	for depth := 1; depth <= SearchDepth; depth++ {
		move, score := search.rootNegamax(depth)
		if search.searchOver || search.Timer.TimeIsUp() {
			break
		}

		bestMove, bestScore = move, score
		fmt.Printf(
			"info depth %d score cp %d time %d nodes %d\n",
			depth, bestScore,
			search.Timer.TimeTaken(),
			search.nodesSearched,
		)
	}
	return bestMove
}

// The top-level function for negamax, which returns a move and a score.
func (search *Search) rootNegamax(depth int) (engine.Move, int) {
	search.nodesSearched = 0
	moves := engine.GenPseduoLegalMoves(&search.Board)
	scores := scoreMoves(&search.Board, &moves)

	bestMove := NullMove
	alpha, beta := NegInf, PosInf

	for index := range moves {
		orderMoves(index, &moves, &scores)
		move := moves[index]

		search.Board.DoMove(move, true)
		if search.Board.KingIsAttacked(search.Board.ColorToMove ^ 1) {
			search.Board.UndoMove(move)
			continue
		}

		score := -search.negamax(depth-1, -beta, -alpha)
		search.Board.UndoMove(move)

		if score == PosInf {
			return move, beta
		}

		if score > alpha {
			alpha = score
			bestMove = move
		}
	}
	return bestMove, alpha
}

// The primary negamax function, which only returns a score and no best move.
func (search *Search) negamax(depth, alpha, beta int) int {
	if search.searchOver || search.Timer.TimeIsUp() {
		return 0
	}

	search.nodesSearched++
	if depth == 0 {
		return search.quiescence(alpha, beta)
	}

	moves := engine.GenPseduoLegalMoves(&search.Board)
	scores := scoreMoves(&search.Board, &moves)

	noMoves := true
	for index := range moves {
		orderMoves(index, &moves, &scores)
		move := moves[index]

		search.Board.DoMove(move, true)
		if search.Board.KingIsAttacked(search.Board.ColorToMove ^ 1) {
			search.Board.UndoMove(move)
			continue
		}

		score := -search.negamax(depth-1, -beta, -alpha)
		search.Board.UndoMove(move)
		noMoves = false

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	if noMoves {
		if search.Board.KingIsAttacked(search.Board.ColorToMove) {
			return NegInf + (SearchDepth - depth - 1)
		}
		return 0
	}

	return alpha
}

// An implementation of a quiesence search algorithm.
func (search *Search) quiescence(alpha, beta int) int {
	if search.searchOver || search.Timer.TimeIsUp() {
		return 0
	}

	search.nodesSearched++
	standPat := evaluateBoard(&search.Board)

	if standPat >= beta {
		return beta
	}

	if alpha < standPat {
		alpha = standPat
	}

	moves := engine.GenPseduoLegalMoves(&search.Board)
	scores := scoreMoves(&search.Board, &moves)

	for index := range moves {
		if engine.MoveType(moves[index]) == engine.Attack || engine.MoveType(moves[index]) == engine.AttackEP {
			orderMoves(index, &moves, &scores)
			move := moves[index]
			search.Board.DoMove(move, true)
			if search.Board.KingIsAttacked(search.Board.ColorToMove ^ 1) {
				search.Board.UndoMove(move)
				continue
			}

			score := -search.quiescence(-beta, -alpha)
			search.Board.UndoMove(move)

			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
	}
	return alpha
}

// Score the moves given
func scoreMoves(board *engine.Board, moves *engine.Moves) (scores []int) {
	scores = make([]int, len(*moves))
	for index, move := range *moves {
		captured := &board.Squares[engine.ToSq(move)]
		moved := &board.Squares[engine.FromSq(move)]
		scores[index] = MvvLva[captured.Type][moved.Type]
	}
	return scores
}

// Order the moves given by finding the best move and putting it
// at the index given.
func orderMoves(index int, moves *engine.Moves, scores *[]int) {
	bestIndex := index
	bestScore := (*scores)[bestIndex]

	for index := bestIndex; index < len(*moves); index++ {
		if (*scores)[index] > bestScore {
			bestIndex = index
			bestScore = (*scores)[index]
		}
	}

	tempMove := (*moves)[index]
	tempScore := (*scores)[index]

	(*moves)[index] = (*moves)[bestIndex]
	(*scores)[index] = (*scores)[bestIndex]

	(*moves)[bestIndex] = tempMove
	(*scores)[bestIndex] = tempScore
}
