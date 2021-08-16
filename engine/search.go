package engine

import (
	"fmt"
	"time"
)

// search.go implements the search routine for Blunder.

const (
	// A constant representing no move.
	NullMove Move = 0

	// A constant representing the maximum search depth that
	// will be attempted.
	MaxDepth = 50

	// The score the best move from the transposition table will
	// be given.
	TT_BestMoveScore = 200

	// Scores for the two killers from each ply. They're ranked below the hash move,
	// and good captures, but above normal quiet moves.
	FirstKillerMoveScore  = 10
	SecondKillerMoveScore = 9
)

// An array that maps move scores to attacker and victim piece types
// for MVV-LVA move ordering: https://www.chessprogramming.org/MVV-LVA.
var MvvLva [7][6]uint8 = [7][6]uint8{
	{16, 15, 14, 13, 12, 11}, // victim Pawn
	{26, 25, 24, 23, 22, 21}, // victim Knight
	{36, 35, 34, 33, 32, 31}, // victim Bishop
	{46, 45, 44, 43, 42, 41}, // vitcim Rook
	{56, 55, 54, 53, 52, 51}, // victim Queen

	{0, 0, 0, 0, 0, 0}, // victim King
	{0, 0, 0, 0, 0, 0}, // No piece
}

// A struct that holds state needed during the search phase. The search
// routines are thus implemented as methods of this struct.
type Search struct {
	Pos        Position
	TransTable TransTable
	Timer      TimeManager

	killerMoves    [MaxDepth][2]Move
	nodesSearched  uint64
	selectiveDepth uint8
	engineColor    uint8
}

// The main search function for Blunder, implemented as an interative
// deepening loop.
func (search *Search) Search() Move {
	search.engineColor = search.Pos.SideToMove
	bestScore := -Inf
	bestMove := NullMove

	search.Timer.Start()

	for depth := 1; depth <= MaxDepth; depth++ {
		// Start a search, and time it for reporting purposes.
		startTime := time.Now()
		move, score := search.rootNegamax(uint8(depth))
		elapsedTime := time.Since(startTime)

		if search.Timer.Stop {
			break
		}

		// Save the best move and best score
		bestMove, bestScore = move, score

		// Report search statistics to the GUI
		fmt.Printf(
			"info depth %d seldepth %d score cp %d time %d nodes %d\n",
			depth, search.selectiveDepth, bestScore,
			elapsedTime.Milliseconds(),
			search.nodesSearched,
		)
	}
	return bestMove
}

// The top-level function for negamax, which returns a move and a score.
func (search *Search) rootNegamax(depth uint8) (Move, int16) {

	// Reset search statisics
	search.nodesSearched = 0
	search.selectiveDepth = 0

	bestMove := NullMove
	alpha, beta := -Inf, Inf

	// Generate the pseduo-legal moves for the current position.
	moves := genMoves(&search.Pos)

	// Score the moves
	search.scoreMoves(&moves, NullMove, 0)

	for index := 0; index < int(moves.Count); index++ {

		// Order the moves to get the best moves first.
		orderMoves(index, &moves)
		move := moves.Moves[index]

		// Make the move, and if it was illegal, undo it and skip to the next move.
		if !search.Pos.MakeMove(move) {
			search.Pos.UnmakeMove(move)
			continue
		}

		score := -search.negamax(depth-1, 0, -beta, -alpha, true)
		search.Pos.UnmakeMove(move)

		// If we have a beta-cutoff (i.e this move gives us a score better than what
		// our opponet can already guarantee early in the tree), return beta and the move
		// that caused the cutoff as the best move.
		if score == beta {
			alpha = beta
			bestMove = move
			break
		}

		// If the score of this move is better than alpha (i.e better than the score
		// we can currently guarantee), set alpha to be the score and the best move
		// to be the move that raised alpha.
		if score > alpha {
			alpha = score
			bestMove = move
			// do_pvs = true
		}
	}

	// Return the best move, and it's score.
	return bestMove, alpha
}

// The primary negamax function, which only returns a score and no best move.
func (search *Search) negamax(depth, ply uint8, alpha, beta int16, do_null bool) int16 {
	// Every 2048 nodes, check if our time has expired.
	if (search.nodesSearched&2047) == 0 && search.Timer.Check() {
		return 0
	}

	// Update the number of nodes searched.
	search.nodesSearched++

	// Check extension extends the search depth by one if we're in check,
	// so that we're less likely to push danger over the search horizon.
	inCheck := sqIsAttacked(
		&search.Pos,
		search.Pos.SideToMove,
		search.Pos.PieceBB[search.Pos.SideToMove][King].Msb())

	if inCheck {
		depth++
	}

	// If we've reached a search depth of zero, enter quiescence
	// search.
	if depth <= 0 {
		return search.quiescence(alpha, beta, ply, 0)
	}

	// Don't do any extra work if the current position is a draw. We
	// can just return a draw value.
	if search.isDraw() {
		return search.contempt()
	}

	// Create a variable to store the possible best move we'll get from probing the transposition
	// table. And the best move we'll get from the search if we don't get a hit.
	ttBestMove := NullMove

	// Probe the transposition table to see if we have a useable matching entry for the current
	// position.
	score := search.TransTable.Probe(search.Pos.Hash, ply, depth, alpha, beta, &ttBestMove)
	if score != Invalid && ply != 0 {
		// If we get a hit, return the score and stop searching.
		return score
	}

	// Do null-move pruning:
	//
	// https://www.chessprogramming.org/Null_Move_Pruning
	//
	// If our opponet is given a free move, can they improve their position? If we do a quick
	// search after giving our opponet this free move and we still find a move with a score better
	// than beta, our opponet can't improve their position and they wouldn't take this path, so we
	// have a beta cut-off and can prune this branch.
	//
	if do_null && !inCheck && depth >= 3 {
		// Do the null move.
		search.Pos.MakeNullMove()
		score := -search.negamax(depth-2-1, ply+1, -beta, -beta+1, false)
		search.Pos.UnmakeNullMove()

		// If we've run out of time, abort the search.
		if search.Timer.Check() {
			return 0
		}

		// If we get a beta cut-off, and it's not a checkmate score,
		// we can use the beta cut-off to send the search and avoid
		// wasting anymore time.
		if score >= beta && abs(score) < Checkmate {
			return beta
		}
	}

	// Generate the moves for the current position.
	moves := genMoves(&search.Pos)
	noMoves := true

	// Score the moves
	search.scoreMoves(&moves, ttBestMove, ply)

	// Set the transposition table entry flag for this node to alpha by default,
	// assuming that we won't raise alpha, and create a variable to store the best
	// move.
	ttFlag := AlphaFlag

	for index := 0; index < int(moves.Count); index++ {

		// Order the moves to get the best moves first.
		orderMoves(index, &moves)
		move := moves.Moves[index]

		// Make the move, and if it was illegal, undo it and skip to the next move.
		if !search.Pos.MakeMove(move) {
			search.Pos.UnmakeMove(move)
			continue
		}

		score := -search.negamax(depth-1, ply+1, -beta, -alpha, true)
		search.Pos.UnmakeMove(move)
		noMoves = false

		// If we have a beta-cutoff (i.e this move gives us a score better than what
		// our opponet can already guarantee early in the tree), return beta and the move
		// that caused the cutoff as the best move.
		if score >= beta {
			alpha = beta

			// Store the killer move for this ply
			search.storeKiller(ply, move)

			// Set the transposition table flag to beta and record the
			// best move.
			ttFlag = BetaFlag
			ttBestMove = move
			break
		}

		// If the score of this move is better than alpha (i.e better than the score
		// we can currently guarantee), set alpha to be the score and the best move
		// to be the move that raised alpha.
		if score > alpha {
			alpha = score

			// Set the transposition table flag to exact and record the
			// best move.
			ttFlag = ExactFlag
			ttBestMove = move

			// do_pvs = true
		}
	}

	// If we don't have any legal moves, it's either checkmate, or a stalemate.
	if noMoves {
		if inCheck {
			// If its checkmate, return a checkmate score of negative infinity,
			// with the current ply added to it. That way, the engine will be
			// rewarded for finding mate quicker, or avoiding mate longer.
			return -Inf + int16(ply)
		} else {
			// If it's a draw, return the draw value.
			return search.contempt()
		}
	}

	// Store the result of the search for this position only if we haven't run out of time.
	if !search.Timer.Check() {
		search.TransTable.Store(search.Pos.Hash, ply, depth, alpha, ttFlag, ttBestMove)
	}

	// Return the best score, which is alpha.
	return alpha
}

// Onece we reach a depth of zero in the main negamax search, instead of
// returning a static evaluation right away, continue to search deeper
// until the position is quiet (i.e there are no winning tatical captures).
// Doing this is known as quiescence search, and it makes the static evaluation
// much more accurate.
func (search *Search) quiescence(alpha, beta int16, negamax_ply uint8, ply uint8) int16 {
	// Every 2048 nodes, check if our time has expired.
	if (search.nodesSearched&2047) == 0 && search.Timer.Check() {
		return 0
	}

	// Update the number of nodes searched.
	search.nodesSearched++

	// Get a static evaluation score for the position.
	score := evaluatePos(&search.Pos)

	// If the score is greater than beta, what our opponet can
	// already guarantee early in the search tree, then we
	// have a beta-cutoff.
	if score >= beta {
		// Update the seldepth to report to the UCI before
		// we return.
		if ply > search.selectiveDepth {
			search.selectiveDepth = ply
		}
		return beta
	}

	// If the score is greater than alpha, what score we can guarantee
	// to get, raise alpha.
	if score > alpha {
		alpha = score
	}

	// Generate all the captures for the current position.
	captures := genMoves(&search.Pos)

	// Score the moves
	search.scoreMoves(&captures, NullMove, negamax_ply)

	for index := 0; index < int(captures.Count); index++ {
		if captures.Moves[index].MoveType() == Attack {

			// Order the moves to get the best moves first.
			orderMoves(index, &captures)
			move := captures.Moves[index]

			// Make the move, and if it was illegal, undo it and skip to the next move.
			if !search.Pos.MakeMove(move) {
				search.Pos.UnmakeMove(move)
				continue
			}

			// While the position is not quiet, continue
			// to search new captures and tatical sequences.
			score := -search.quiescence(-beta, -alpha, negamax_ply, ply+1)
			search.Pos.UnmakeMove(move)

			// If we have a beta-cutoff (i.e this move gives us a score better than what
			// our opponet can already guarantee early in the tree), return beta and the move
			// that caused the cutoff as the best move.
			if score >= beta {
				alpha = beta
				break
			}

			// If the score of this move is better than alpha (i.e better than the score
			// we can currently guarantee), set alpha to be the score and the best move
			// to be the move that raised alpha.
			if score > alpha {
				alpha = score
			}
		}
	}

	// Update the seldepth to report to the UCI
	// before we return.
	if ply > search.selectiveDepth {
		search.selectiveDepth = ply
	}

	// Return the best score, which is alpha.
	return alpha
}

// Given a "killer move" (a quiet move that caused a beta cut-off), store the
// Move in the slot for the given depth.
func (search *Search) storeKiller(ply uint8, move Move) {
	if move.MoveType() != Attack {
		if !move.Equal(search.killerMoves[ply][0]) {
			search.killerMoves[ply][1] = search.killerMoves[ply][0]
			search.killerMoves[ply][0] = move
		}
	}
}

// Determine if the current position is a draw.
func (search *Search) isDraw() bool {
	if search.Pos.Rule50 >= 100 {
		return true
	}
	return search.isRepition()
}

// Determine if the current board state is being repeated.
func (search *Search) isRepition() bool {
	var repPly uint16
	for repPly = 0; repPly < search.Pos.HistoryPly; repPly++ {
		if search.Pos.History[repPly] == search.Pos.Hash {
			return true
		}
	}
	return false
}

// Determine the draw score based on whose moving. If the engine is moving,
// return a negative score, and if the opponet is moving, return a positive
// score.
func (search *Search) contempt() int16 {
	if search.Pos.SideToMove == search.engineColor {
		return Draw
	}
	return -Draw
}

// Score the moves given
func (search *Search) scoreMoves(moves *MoveList, ttBestMove Move, ply uint8) {
	for index := 0; index < int(moves.Count); index++ {
		move := &moves.Moves[index]

		if ttBestMove.Equal(*move) {
			move.AddScore(TT_BestMoveScore)
		} else if search.killerMoves[ply][0].Equal(*move) {
			move.AddScore(FirstKillerMoveScore)
		} else if search.killerMoves[ply][1].Equal(*move) {
			move.AddScore(SecondKillerMoveScore)
		} else {
			captured := &search.Pos.Squares[move.ToSq()]
			moved := &search.Pos.Squares[move.FromSq()]
			move.AddScore(MvvLva[captured.Type][moved.Type])
		}
	}
}

// Order the moves given by finding the best move and putting it
// at the index given.
func orderMoves(index int, moves *MoveList) {
	bestIndex := index
	bestScore := moves.Moves[bestIndex].Score()

	for index := bestIndex; index < int(moves.Count); index++ {
		if moves.Moves[index].Score() > bestScore {
			bestIndex = index
			bestScore = (*moves).Moves[index].Score()
		}
	}

	tempMove := moves.Moves[index]
	moves.Moves[index] = moves.Moves[bestIndex]
	moves.Moves[bestIndex] = tempMove
}
