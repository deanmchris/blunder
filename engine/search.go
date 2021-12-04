package engine

import (
	"fmt"
	"math"
	"time"
)

const (
	// The maximum depth the engine will attempt to reach.
	MaxPly = 100

	// A constant representing no move.
	NullMove Move = 0

	// A constant representing the score of the principal variation
	// move from the transposition table.
	PVMoveScore uint16 = 60

	// A constant representing the score of a killer move.
	KillerMoveScore uint16 = 10

	// A constant representing the maximum number of killers.
	MaxKillers = 2

	// A constant to offset the score of the pv and MVV-LVA move higher
	// than killers and history heuristic moves.
	MvvLvaOffset uint16 = math.MaxUint16 - 256

	// A constant representing the maximum value a history heuristic score
	// is allowed to reach.
	MaxHistoryScore int32 = int32(MvvLvaOffset) - int32((MaxKillers+1)*KillerMoveScore)

	StaticNullMovePruningBaseMargin int16 = 120
	FirstNullMoveReduction          int8  = 2
	SecondNullMoveReduction         int8  = 3
	LateMoveReduction               int8  = 2
	LMRLegalMovesLimit              int   = 4
	LMRDepthLimit                   int8  = 3
)

var FutilityMargins = [4]int16{0, 200, 300, 500}

type Reduction struct {
	MoveLimit int
	Reduction int8
}

var LateMoveReductions = [10]Reduction{
	{4, 2},
	{8, 3},
	{12, 4},
	{16, 5},
	{20, 6},
	{24, 7},
	{28, 8},
	{32, 9},
	{34, 10},
	{100, 12},
}

// An array that maps move scores to attacker and victim piece types
// for MVV-LVA move ordering: https://www.chessprogramming.org/MVV-LVA.
var MvvLva [7][6]uint16 = [7][6]uint16{
	{15, 14, 13, 12, 11, 10}, // victim Pawn
	{25, 24, 23, 22, 21, 20}, // victim Knight
	{35, 34, 33, 32, 31, 30}, // victim Bishop
	{45, 44, 43, 42, 41, 40}, // victim Rook
	{55, 54, 53, 52, 51, 50}, // victim Queen

	{0, 0, 0, 0, 0, 0}, // victim King
	{0, 0, 0, 0, 0, 0}, // No piece
}

// A struct representing a principal variation line.
type PVLine struct {
	Moves []Move
}

// Clear the principal variation line.
func (pvLine *PVLine) Clear() {
	pvLine.Moves = nil
}

// Update the principal variation line with a new best move,
// and a new line of best play after the best move.
func (pvLine *PVLine) Update(move Move, newPVLine PVLine) {
	pvLine.Clear()
	pvLine.Moves = append(pvLine.Moves, move)
	pvLine.Moves = append(pvLine.Moves, newPVLine.Moves...)
}

// Get the best move from the principal variation line.
func (pvLine *PVLine) GetPVMove() Move {
	return pvLine.Moves[0]
}

// Convert the principal variation line to a string.
func (pvLine PVLine) String() string {
	pv := fmt.Sprintf("%s", pvLine.Moves)
	return pv[1 : len(pv)-1]
}

// A struct that holds state needed during the search phase. The search
// routines are thus implemented as methods of this struct.
type Search struct {
	Pos   Position
	Timer TimeManager
	TT    TransTable

	side       uint8
	nodes      uint64
	totalNodes uint64

	killers [MaxPly + 1][MaxKillers]Move
	history [2][64][64]int32

	SpecifiedDepth uint8
	SpecifiedNodes uint64
}

// The main search function for Blunder, implemented as an interative
// deepening loop.
func (search *Search) Search() Move {
	search.side = search.Pos.SideToMove
	var pvLine PVLine
	bestMove := NullMove

	search.ageHistoryTable()
	search.Timer.Start()

	search.totalNodes = 0
	depth := uint8(0)

	for depth = 1; depth <= MaxPly && depth <= search.SpecifiedDepth && search.SpecifiedNodes > 0; depth++ {
		// Clear the nodes searched and the last iterations pv line.
		search.nodes = 0
		pvLine.Clear()

		// Start a search, and time it for reporting purposes.
		startTime := time.Now()
		score := search.negamax(int8(depth), 0, -Inf, Inf, &pvLine, true)
		endTime := time.Since(startTime)

		if search.Timer.Stop {
			if bestMove == NullMove && depth == 1 {
				bestMove = pvLine.GetPVMove()
			}
			break
		}

		// Save the best move and report search statistics to the GUI
		bestMove = pvLine.GetPVMove()

		// Get the nodes per second
		nps := uint64(float64(search.nodes) / float64(endTime.Seconds()))

		// Collect the amount of nodes searched for this iteration.
		search.totalNodes += search.nodes

		// Send search statistics to the GUI.
		fmt.Printf(
			"info depth %d score %s nodes %d nps %d time %d pv %s\n",
			depth, getMateOrCPScore(score),
			search.nodes, nps,
			endTime.Milliseconds(),
			pvLine,
		)
	}

	// Return the best move found to the GUI.
	return bestMove
}

// Display the correct format for the search score if it's a centipawn score
// or a checkmate score.
func getMateOrCPScore(score int16) string {
	if score > Checkmate {
		pliesToMate := Inf - score
		mateInN := (pliesToMate / 2) + (pliesToMate % 2)
		return fmt.Sprintf("mate %d", mateInN)
	}

	if score < -Checkmate {
		pliesToMate := -Inf - score
		mateInN := (pliesToMate / 2) + (pliesToMate % 2)
		return fmt.Sprintf("mate %d", mateInN)
	}

	return fmt.Sprintf("cp %d", score)
}

// The primary negamax function.
func (search *Search) negamax(depth int8, ply uint8, alpha, beta int16, pvLine *PVLine, doNull bool) int16 {
	// Update the number of nodes searched.
	search.nodes++

	if ply >= MaxPly {
		return EvaluatePos(&search.Pos)
	}

	// If a given node amount to search was given, make sure we haven't passed it
	// and if so stop the search.
	if search.totalNodes+search.nodes >= search.SpecifiedNodes {
		search.Timer.Stop = true
		return 0
	}

	// Every 2048 nodes, check if our time has expired.
	if (search.nodes & 2047) == 0 {
		search.Timer.Check()
	}

	// If we're told to stop, abort the current search and return 0. This won't
	// affect anything, as the previous search's best move will be used, and
	// everything from the current search will be discarded.
	if search.Timer.Stop {
		return 0
	}

	isRoot := ply == 0
	isPVNode := beta-alpha != 1
	inCheck := search.Pos.InCheck()
	canFutilityPrune := false
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
		search.nodes--
		return search.Qsearch(alpha, beta, ply, pvLine)
	}

	// Don't do any extra work if the current position is a draw. We
	// can just return a draw value.
	if !isRoot && (search.Pos.Rule50 >= 100 || search.isDrawByRepition() || search.Pos.EndgameIsDrawn()) {
		return search.contempt()
	}

	// Create a variable to store the possible best move we'll get from probing the transposition
	// table. And the best move we'll get from the search if we don't get a hit.
	ttMove := NullMove

	// Probe the transposition table to see if we have a useable matching entry for the current
	// position. If we get a hit, return the score and stop searching.
	score := search.TT.Probe(search.Pos.Hash, ply, uint8(depth), alpha, beta, &ttMove)
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

	if !inCheck && !isPVNode && abs16(beta) < Checkmate {
		staticScore := EvaluatePos(&search.Pos)
		scoreMargin := StaticNullMovePruningBaseMargin * int16(depth)
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

	if doNull && !inCheck && !isPVNode && depth >= 3 && !search.Pos.NoMajorsOrMiniors() {
		R := 3 + depth/6
		search.Pos.MakeNullMove()
		score := -search.negamax(depth-R-1, ply+1, -beta, -beta+1, &childPVLine, false)
		search.Pos.UnmakeNullMove()
		childPVLine.Clear()

		if search.Timer.Stop {
			return 0
		}

		if score >= beta && abs16(score) < Checkmate {
			return beta
		}
	}

	// =====================================================================//
	// FUTILITY PRUNING: If we're close to the horizon, and even with a     //
	// large margin the static evaluation can't be raised above alpha,      //
	// we're probably in a fail-low node, and many moves can be probably    //
	// be pruned. So set a flag so we don't waste time searching moves that //
	// suck and probably don't even have a chance of raising alpha.         //
	// =====================================================================//

	if depth <= 3 && !isPVNode && !inCheck && alpha < Checkmate {
		staticScore := EvaluatePos(&search.Pos)
		if staticScore+FutilityMargins[depth] <= alpha {
			canFutilityPrune = true
		}
	}

	// Generate and score the moves for the side to move
	// in the current position.
	moves := GenMoves(&search.Pos)
	search.scoreMoves(&moves, ttMove, ply)

	// Set up variables to record the number of legal moves and
	// the transposition table entry flag.
	legalMoves := 0
	ttFlag := AlphaFlag

	// Set up variables to record the best move and best score of the
	// position.
	bestMove := NullMove
	bestScore := -Inf

	for index := 0; index < int(moves.Count); index++ {
		// Order the moves to get the best moves first.
		orderMoves(index, &moves)
		move := moves.Moves[index]

		// Make the move, and if it was illegal, undo it and skip to the next move.
		if !search.Pos.MakeMove(move) {
			search.Pos.UnmakeMove(move)
			continue
		}

		legalMoves++

		// =====================================================================//
		// FUTILITY PRUNING: If we're close to the horizon, and even with a     //
		// large margin the static evaluation can't be raised above alpha,      //
		// we're probably in a fail-low node, and many moves can be probably    //
		// be pruned. So set a flag so we don't waste time searching moves that //
		// suck and probably don't even have a chance of raising alpha.         //
		// =====================================================================//

		if canFutilityPrune && legalMoves > 1 {
			tactical := search.Pos.InCheck() || move.MoveType() == Attack || move.MoveType() == Promotion
			if !tactical {
				search.Pos.UnmakeMove(move)
				continue
			}
		}

		// =====================================================================//
		// LATE MOVE REDUCTION: Since our move ordering is good, the            //
		// first move is likely to be the best move in the position, which      //
		// means it's part of the principal variation. So instead of searching  //
		// every move equally, search the first move with full-depth and full-  //
		// window, and search every move after with a reduced-depth and null-   //
		// window to prove it'll fail low cheaply. If it raises alpha however,  //
		// we have to use a full-window, a full-depth, or both to get an        //
		// accurate score for the move.                                         //
		// =====================================================================//

		score := int16(0)
		if legalMoves == 1 {
			score = -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine, true)
		} else {
			tactical := inCheck || move.MoveType() == Attack
			reduction := int8(0)

			if !isPVNode && legalMoves >= LMRLegalMovesLimit && depth >= LMRDepthLimit && !tactical {
				for _, red := range LateMoveReductions {
					if red.MoveLimit >= legalMoves {
						reduction = red.Reduction
						break
					}
				}
			}

			score = -search.negamax(depth-1-reduction, ply+1, -alpha-1, -alpha, &childPVLine, true)

			if score > alpha && reduction > 0 {
				score = -search.negamax(depth-1-reduction, ply+1, -beta, -alpha, &childPVLine, true)
				if score > alpha {
					score = -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine, true)
				}
			} else if score > alpha && score < beta {
				score = -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine, true)
			}
		}

		search.Pos.UnmakeMove(move)

		// If the current score is better than the best score so far,
		// update the best score and the best move.
		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		// If we have a beta-cutoff (i.e this move gives us a score better than what
		// our opponet can already guarantee early in the tree), return beta and the move
		// that caused the cutoff as the best move.
		if score >= beta {
			// Set the transposition table entry flag to beta
			ttFlag = BetaFlag

			// Store the possible killer.
			search.storeKiller(ply, move)

			// Update the history table.
			search.updateHistoryTable(move, depth)

			// Break he move loop, since we have a beta cutoff.
			break
		}

		// If the score of this move is better than alpha (i.e better than the score
		// we can currently guarantee), set alpha to be the score and the best move
		// to be the move that raised alpha.
		if score > alpha {
			// Update alpha and the best move.
			alpha = score

			// Update the history table.
			search.updateHistoryTable(move, depth)

			// Update the principal variation line.
			pvLine.Update(move, childPVLine)

			// Set the transposition table flag to exact.
			ttFlag = ExactFlag
		}

		// Clear this child node's principal variation line for the
		// next child node.
		childPVLine.Clear()
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

	// If we're not out of time, store the result of the search for this position.
	if !search.Timer.Stop {
		search.TT.Store(search.Pos.Hash, ply, uint8(depth), bestScore, ttFlag, bestMove)
	}

	// Return the best score, which is alpha.
	return bestScore
}

// Onece we reach a depth of zero in the main negamax search, instead of
// returning a static evaluation right away, continue to search deeper using
// a special form of negamax until the position is quiet (i.e there are no
// winning tatical captures). Doing this is known as quiescence search, and
// it makes the static evaluation much more accurate.
func (search *Search) Qsearch(alpha, beta int16, negamaxPly uint8, pvLine *PVLine) int16 {
	search.nodes++

	if (search.nodes & 2047) == 0 {
		search.Timer.Check()
	}

	if search.totalNodes+search.nodes >= search.SpecifiedNodes {
		search.Timer.Stop = true
	}

	if search.Timer.Stop {
		return 0
	}

	bestScore := EvaluatePos(&search.Pos)

	// If the score is greater than beta, what our opponet can
	// already guarantee early in the search tree, then we
	// have a beta-cutoff.
	if bestScore >= beta {
		return bestScore
	}

	// If the score is greater than alpha, what score we can guarantee
	// to get, raise alpha.
	if bestScore > alpha {
		alpha = bestScore
	}

	moves := genCaptures(&search.Pos)
	search.scoreMoves(&moves, NullMove, negamaxPly)
	var childPVLine PVLine

	for index := 0; index < int(moves.Count); index++ {
		orderMoves(index, &moves)
		move := moves.Moves[index]

		// Prune moves based on the static exchange evaluation being negative.
		if search.Pos.See(move) < 0 {
			continue
		}

		if !search.Pos.MakeMove(move) {
			search.Pos.UnmakeMove(move)
			continue
		}

		score := -search.Qsearch(-beta, -alpha, negamaxPly, &childPVLine)
		search.Pos.UnmakeMove(move)

		if score > bestScore {
			bestScore = score
		}

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
			pvLine.Update(move, childPVLine)
		}

		childPVLine.Clear()
	}

	return bestScore
}

// Update the history heuristics table if the move that caused a beta-cutoff is quiet.
func (search *Search) updateHistoryTable(move Move, depth int8) {
	if search.Pos.Squares[move.ToSq()].Type == NoType {
		search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] += int32(depth) * int32(depth)
	}

	if search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] >= MaxHistoryScore {
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

// Clear the values in the history table.
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
	for repPly = 0; repPly < HistoryPly; repPly++ {
		if PositionHistories[repPly] == search.Pos.Hash {
			return true
		}
	}
	return false
}

// Score the moves generated.
func (search *Search) scoreMoves(moves *MoveList, pvMove Move, ply uint8) {
	for index := 0; index < int(moves.Count); index++ {
		move := &moves.Moves[index]
		captured := search.Pos.Squares[move.ToSq()]

		if move.Equal(pvMove) {
			move.AddScore(MvvLvaOffset + PVMoveScore)
		} else if captured.Type != NoType {
			moved := search.Pos.Squares[move.FromSq()]
			move.AddScore(MvvLvaOffset + MvvLva[captured.Type][moved.Type])
		} else {
			moveScore := uint16(0)
			for n, killer := range search.killers[ply] {
				if move.Equal(killer) {
					moveScore = MvvLvaOffset - (uint16(n+1) * KillerMoveScore)
					break
				}
			}

			if moveScore == 0 {
				moveScore = uint16(search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()])
			}

			move.AddScore(moveScore)
		}
	}
}

// Order the moves given by finding the best move and putting it
// at the index given.
func orderMoves(currIndex int, moves *MoveList) {
	bestIndex := currIndex
	bestScore := moves.Moves[bestIndex].Score()

	for index := bestIndex; index < int(moves.Count); index++ {
		if moves.Moves[index].Score() > bestScore {
			bestIndex = index
			bestScore = moves.Moves[index].Score()
		}
	}

	tempMove := moves.Moves[currIndex]
	moves.Moves[currIndex] = moves.Moves[bestIndex]
	moves.Moves[bestIndex] = tempMove
}
