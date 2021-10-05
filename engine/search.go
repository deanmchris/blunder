package engine

import (
	"fmt"
	"time"
)

const (
	// The maximum depth the engine will attempt to reach.
	MaxDepth = 50

	// A constant representing no move.
	NullMove Move = 0

	// Constants representing the score for a killer move and
	// the principal variation move from the transposition table.
	KillerMoveScore int16 = 10
	PVMoveScore     int16 = 60

	// A constant offset added to the pv move, MVV-LVA moves, and killers
	// to give room for scoring below the value with the history heurustic.
	HistoryOffset = 32000
)

// An array that maps move scores to attacker and victim piece types
// for MVV-LVA move ordering: https://www.chessprogramming.org/MVV-LVA.
var MvvLva [7][6]int16 = [7][6]int16{
	{16, 15, 14, 13, 12, 11}, // victim Pawn
	{26, 25, 24, 23, 22, 21}, // victim Knight
	{36, 35, 34, 33, 32, 31}, // victim Bishop
	{46, 45, 44, 43, 42, 41}, // vitcim Rook
	{56, 55, 54, 53, 52, 51}, // victim Queen

	{0, 0, 0, 0, 0, 0}, // victim King
	{0, 0, 0, 0, 0, 0}, // No piece
}

// A struct representing a principal variation line.
type PVLine struct {
	moves []Move
}

// Clear the principal variation line.
func (pvLine *PVLine) Clear() {
	pvLine.moves = nil
}

// Update the principal variation line with a new best move,
// and a new line of best play after the best move.
func (pvLine *PVLine) Update(move Move, newPVLine PVLine) {
	pvLine.Clear()
	pvLine.moves = append(pvLine.moves, move)
	pvLine.moves = append(pvLine.moves, newPVLine.moves...)
}

// Get the best move from the principal variation line.
func (pvLine *PVLine) GetPVMove() Move {
	return pvLine.moves[0]
}

// Convert the principal variation line to a string.
func (pvLine PVLine) String() string {
	pv := fmt.Sprintf("%s", pvLine.moves)
	return pv[1 : len(pv)-1]
}

// A struct that holds state needed during the search phase. The search
// routines are thus implemented as methods of this struct.
type Search struct {
	Pos   Position
	Timer TimeManager
	TT    TransTable

	side  uint8
	nodes uint64

	killers [MaxDepth][2]Move
	history [2][64][64]int16
}

// The main search function for Blunder, implemented as an interative
// deepening loop.
func (search *Search) Search() Move {
	search.side = search.Pos.SideToMove
	var pvLine PVLine
	bestMove := NullMove

	search.ageHistoryTable()
	search.Timer.Start()

	for depth := 1; depth <= MaxDepth; depth++ {
		// Clear the nodes searched and the last iterations pv line.
		search.nodes = 0
		pvLine.Clear()

		// Start a search, and time it for reporting purposes.
		startTime := time.Now()
		score := search.negamax(uint8(depth), 0, -Inf, Inf, &pvLine, false)
		endTime := time.Since(startTime)

		if search.Timer.Stop {
			if bestMove == NullMove && depth == 1 {
				bestMove = pvLine.GetPVMove()
			}
			break
		}

		// Save the best move and report search statistics to the GUI
		bestMove = pvLine.GetPVMove()
		fmt.Printf(
			"info depth %d score cp %d time %d nodes %d\n",
			depth, score,
			endTime.Milliseconds(),
			search.nodes,
			// pvLine,
		)
	}

	// Return the best move found to the GUI.
	return bestMove
}

// The primary negamax function.
func (search *Search) negamax(depth, ply uint8, alpha, beta int16, pvLine *PVLine, doNull bool) int16 {
	// Every 2048 nodes, check if our time has expired.
	if (search.nodes&2047) == 0 && search.Timer.Check() {
		return 0
	}

	// Update the number of nodes searched.
	search.nodes++

	isRoot := ply == 0
	inCheck := search.Pos.InCheck()
	var childPVLine PVLine

	// =====================================================================//
	// CHECK EXTENSION: Extend the search depth by one if we're in check,   //
	// so that we're less likely to push danger over the search horizon,    //
	// and we won't enter quiescence search while in check.                 //
	// =====================================================================//
	if inCheck {
		depth++
	}

	// If we've reached a search depth of zero, enter quiescence
	// search to stabilize the position before returning a static
	// score.
	if depth <= 0 {
		return search.qsearch(alpha, beta, ply)
	}

	// Don't do any extra work if the current position is a draw. We
	// can just return a draw value.
	if !isRoot && (search.Pos.Rule50 >= 100 || search.isDrawByRepition()) {
		return search.contempt()
	}

	// Create a variable to store the possible best move we'll get from probing the transposition
	// table. And the best move we'll get from the search if we don't get a hit.
	ttMove := NullMove

	// Probe the transposition table to see if we have a useable matching entry for the current
	// position. If we get a hit, return the score and stop searching.
	score := search.TT.Probe(search.Pos.Hash, ply, depth, alpha, beta, &ttMove)
	if score != Invalid && !isRoot {
		return score
	}

	// =====================================================================//
	// STATIC NULL MOVE PRUNING: If our current material score is so good   //
	// that even if we give ourselves a big hit materially and subtract a   //
	// large amount of our material score (the "score margin") and our      //
	// material score is still greater than beta, we assume this node will  //
	// fail-high and we can prune its branch.                               //
	// =====================================================================//

	if !inCheck && abs16(beta) < Checkmate {
		staticScore := evaluatePos(&search.Pos)
		var scoreMargin int16 = 120 * int16(depth)
		if staticScore-scoreMargin >= beta {
			return beta
		}
	}

	// =====================================================================//
	// NULL MOVE PRUNING: If our opponet is given a free move, can they     //
	// improve their position? If we do a quick search after giving our     //
	// opponet this free move and we still find a move with a score better  //
	// than beta, our opponet can't improve their position and they         //
	// wouldn't take this path, so we have a beta cut-off and can prune     //
	// this branch.                                                         //
	// =====================================================================//

	if doNull && !inCheck && depth >= 3 {
		var R uint8 = 3
		if depth > 6 {
			R = 4
		}

		search.Pos.MakeNullMove()
		score := -search.negamax(depth-R, ply+1, -beta, -beta+1, &childPVLine, false)
		search.Pos.UnmakeNullMove()
		childPVLine.Clear()

		if search.Timer.Check() {
			return 0
		}
		if score >= beta && abs16(score) < Checkmate {
			return beta
		}
	}

	// Generate and score the moves for the side to move
	// in the current position.
	moves := genMoves(&search.Pos)
	search.scoreMoves(&moves, ply, ttMove)

	legalMoves := 0
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

		score := -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine, true)
		search.Pos.UnmakeMove(move)
		legalMoves++

		// If we have a beta-cutoff (i.e this move gives us a score better than what
		// our opponet can already guarantee early in the tree), return beta and the move
		// that caused the cutoff as the best move.
		if score >= beta {
			// If we're not out of time, store beta and the move that caused the beta-cutoff
			// in the transposition table, and update the killer moves, and history heuristic.
			if !search.Timer.Check() {
				search.storeKiller(ply, move)
				search.updateHistoryTable(move, depth)
				search.TT.Store(search.Pos.Hash, ply, depth, beta, BetaFlag, move)
			}
			return beta
		}

		// If the score of this move is better than alpha (i.e better than the score
		// we can currently guarantee), set alpha to be the score and the best move
		// to be the move that raised alpha.
		if score > alpha {
			alpha = score

			// Update the principal variation line.
			pvLine.Update(move, childPVLine)

			// Set the transposition table flag to exact and record the
			// best move.
			ttFlag = ExactFlag
			ttMove = move
		}
	}

	// If we don't have any legal moves, it's either checkmate, or a stalemate.
	if legalMoves == 0 {
		if inCheck {
			// If its checkmate, return a checkmate score of negative infinity,
			// with the current ply added to it. That way, the engine will be
			// rewarded for finding mate quicker, or avoiding mate longer.
			return -Inf + int16(ply)
		}
		// If it's a draw, return the draw value.
		return search.contempt()
	}

	// Store the result of the search for this position only if we haven't run out of time.
	if !search.Timer.Stop {
		bestMove := NullMove
		if ttFlag == ExactFlag {
			bestMove = pvLine.GetPVMove()
		}
		search.TT.Store(search.Pos.Hash, ply, depth, alpha, ttFlag, bestMove)
	}

	// Return the best score, which is alpha.
	return alpha
}

// Onece we reach a depth of zero in the main negamax search, instead of
// returning a static evaluation right away, continue to search deeper
// until the position is quiet (i.e there are no winning tatical captures).
// Doing this is known as quiescence search, and it makes the static evaluation
// much more accurate.
func (search *Search) qsearch(alpha, beta int16, negamaxPly uint8) int16 {
	if (search.nodes&2047) == 0 && search.Timer.Check() {
		return 0
	}

	search.nodes++

	if search.Pos.Rule50 >= 100 || search.isDrawByRepition() {
		return search.contempt()
	}

	staticScore := evaluatePos(&search.Pos)

	// If the score is greater than beta, what our opponet can
	// already guarantee early in the search tree, then we
	// have a beta-cutoff.
	if staticScore >= beta {
		return beta
	}

	// If the score is greater than alpha, what score we can guarantee
	// to get, raise alpha.
	if staticScore > alpha {
		alpha = staticScore
	}

	moves := genMoves(&search.Pos)
	search.scoreMoves(&moves, negamaxPly, NullMove)

	for index := 0; index < int(moves.Count); index++ {
		if moves.Moves[index].MoveType() == Attack {
			orderMoves(index, &moves)
			move := moves.Moves[index]

			if !search.Pos.MakeMove(move) {
				search.Pos.UnmakeMove(move)
				continue
			}

			score := -search.qsearch(-beta, -alpha, negamaxPly)
			search.Pos.UnmakeMove(move)

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

// Update the history heuristics table if the move that caused a beta-cutoff is quiet.
func (search *Search) updateHistoryTable(move Move, depth uint8) {
	if search.Pos.Squares[move.ToSq()].Type == NoType {
		search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] += int16(depth) * int16(depth)
	}

	if search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] > HistoryOffset+KillerMoveScore-1 {
		search.ageHistoryTable()
	}
}

// Age the values in the history table by halving them.
func (search *Search) ageHistoryTable() {
	for sq1 := 0; sq1 < 64; sq1++ {
		for sq2 := 0; sq2 < 64; sq2++ {
			search.history[search.Pos.SideToMove][sq1][sq2] /= 2
		}
	}
}

// Clear the values in the history table..
func (search *Search) ClearHistoryTable() {
	for sq1 := 0; sq1 < 64; sq1++ {
		for sq2 := 0; sq2 < 64; sq2++ {
			search.history[search.Pos.SideToMove][sq1][sq2] = 0
		}
	}
}

// Given a "killer move" (a quiet move that caused a beta cut-off), store the
// Move in the slot for the given depth.
func (search *Search) storeKiller(ply uint8, move Move) {
	if search.Pos.Squares[move.ToSq()].Type == NoType {
		if !move.Equal(search.killers[ply][0]) {
			search.killers[ply][1] = search.killers[ply][0]
			search.killers[ply][0] = move
		}
	}
}

// Determine the draw score based on the phase of the game and whose moving,
// to encourge the engine to strive for a win in the middle-game, but be
// satisified with a draw in the endgame.
func (search *Search) contempt() int16 {
	drawValue := MiddleGameDraw
	if search.Pos.IsEndgameForSide() {
		drawValue = EndGameDraw
	}

	if search.Pos.SideToMove == search.side {
		return -drawValue
	}
	return drawValue
}

// Determine if the current board state is being repeated.
func (search *Search) isDrawByRepition() bool {
	var repPly uint16
	for repPly = 0; repPly < search.Pos.HistoryPly; repPly++ {
		if search.Pos.History[repPly] == search.Pos.Hash {
			return true
		}
	}
	return false
}

// Score the moves generated.
func (search *Search) scoreMoves(moves *MoveList, ply uint8, pvMove Move) {
	for index := 0; index < int(moves.Count); index++ {
		move := &moves.Moves[index]
		captured := &search.Pos.Squares[move.ToSq()]

		if pvMove.Equal(*move) {
			move.AddScore(HistoryOffset + PVMoveScore)
		} else if captured.Type != NoType {
			moved := &search.Pos.Squares[move.FromSq()]
			move.AddScore(HistoryOffset + MvvLva[captured.Type][moved.Type])
		} else {
			if search.killers[ply][0].Equal(*move) {
				move.AddScore(HistoryOffset + KillerMoveScore)
			} else if search.killers[ply][1].Equal(*move) {
				move.AddScore(HistoryOffset + KillerMoveScore)
			} else {
				move.AddScore(search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()])
			}
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
