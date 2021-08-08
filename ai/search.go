package ai

import (
	"blunder/engine"
	"fmt"
	"time"
)

const (
	SearchDepth             = 50
	NullMove    engine.Move = 0
	DrawValue               = -(PawnValue / 2)
	R                       = 2

	TTMoveScore           = 60
	FirstKillerMoveScore  = 10
	SecondKillerMoveScore = 9
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
	Board engine.Board
	Timer Timer

	killerMoves   [SearchDepth + 1][2]engine.Move
	nodesSearched uint64
	engineColor   uint8
}

// Search for the best move for the side to move in the given position.
// Implemented using iterative deepening.
func (search *Search) Search() engine.Move {
	search.engineColor = search.Board.ColorToMove
	bestMove, bestScore := NullMove, NegInf

	search.Timer.StartSearch()

	for depth := 1; depth <= SearchDepth; depth++ {
		searchTimeStart := time.Now()
		move, score := search.rootNegamax(depth)
		searchTimeEnd := time.Since(searchTimeStart)

		if search.Timer.TimeIsUp() {
			break
		}

		bestMove, bestScore = move, score
		fmt.Printf(
			"info depth %d score cp %d time %d nodes %d\n",
			depth, bestScore,
			searchTimeEnd/time.Millisecond,
			search.nodesSearched,
		)
	}
	return bestMove
}

// The top-level function for negamax, which returns a move and a score.
func (search *Search) rootNegamax(depth int) (engine.Move, int) {
	search.nodesSearched = 0
	moves := engine.GenPseduoLegalMoves(&search.Board)
	scores := search.scoreMoves(&moves, -1, NullMove)

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

		score := -search.negamax(depth-1, 0, -beta, -alpha)
		search.Board.UndoMove(move)

		if score == beta {
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
func (search *Search) negamax(depth, ply, alpha, beta int) int {
	if search.Timer.TimeIsUp() {
		return 0
	}

	search.nodesSearched++
	inCheck := search.Board.KingIsAttacked(search.Board.ColorToMove)

	if inCheck {
		depth++
	}

	if search.isDraw() {
		return search.contempt()
	}

	if depth == 0 {
		return search.quiescence(alpha, beta)
	}

	moves := engine.GenPseduoLegalMoves(&search.Board)
	scores := search.scoreMoves(&moves, depth, NullMove)

	noMoves := true
	for index := range moves {
		orderMoves(index, &moves, &scores)
		move := moves[index]

		search.Board.DoMove(move, true)
		if search.Board.KingIsAttacked(search.Board.ColorToMove ^ 1) {
			search.Board.UndoMove(move)
			continue
		}

		score := -search.negamax(depth-1, ply+1, -beta, -alpha)
		search.Board.UndoMove(move)
		noMoves = false

		if score >= beta {
			search.storeKiller(depth, move)
			return beta
		}

		if score > alpha {
			alpha = score
		}
	}

	if noMoves {
		if inCheck {
			alpha = NegInf + ply
		} else {
			alpha = search.contempt()
		}
	}
	return alpha
}

// An implementation of a quiesence search algorithm.
func (search *Search) quiescence(alpha, beta int) int {
	if search.Timer.TimeIsUp() {
		return 0
	}

	search.nodesSearched++
	standPat := EvaluateBoard(&search.Board)

	if standPat >= beta {
		return beta
	}

	if alpha < standPat {
		alpha = standPat
	}

	moves := engine.GenPseduoLegalCaptures(&search.Board)
	scores := search.scoreMoves(&moves, -1, NullMove)

	for index := range moves {
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
	return alpha
}

func (search *Search) isDraw() bool {
	if search.Board.Rule50 >= 100 {
		return true
	}
	return search.isRepition()
}

// Determine if the current board state is being repeated.
func (search *Search) isRepition() bool {
	var repPly uint16
	for repPly = 0; repPly < search.Board.RepitionPly; repPly++ {
		if search.Board.Repitions[repPly] == search.Board.Hash {
			return true
		}
	}
	return false
}

// Determine the draw score based on whose moving. If the engine is moving,
// return a negative score, and if the opponet is moving, return a positive
// score.
func (search *Search) contempt() int {
	if search.Board.ColorToMove == search.engineColor {
		return DrawValue
	}
	return -DrawValue
}

// Given a "killer move" (a quiet move that caused a beta cut-off), store the
// Move in the slot for the given depth.
func (search *Search) storeKiller(depth int, move engine.Move) {
	if move.MoveType() != engine.Attack && move.MoveType() != engine.AttackEP {
		if move != search.killerMoves[depth][0] {
			search.killerMoves[depth][1] = search.killerMoves[depth][0]
			search.killerMoves[depth][0] = move
		}
	}
}

// Score the moves given
func (search *Search) scoreMoves(moves *engine.Moves, depth int, ttMove engine.Move) (scores []int) {
	scores = make([]int, len(*moves))
	for index, move := range *moves {
		captured := &search.Board.Squares[move.ToSq()]
		moved := &search.Board.Squares[move.FromSq()]
		scores[index] = MvvLva[captured.Type][moved.Type]

		if move == ttMove {
			scores[index] = TTMoveScore
		} else if depth > 0 && move == search.killerMoves[depth][0] {
			scores[index] = FirstKillerMoveScore
		} else if depth > 0 && move == search.killerMoves[depth][1] {
			scores[index] = SecondKillerMoveScore
		}
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
