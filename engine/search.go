package engine

import (
	"fmt"
	"math"
	"time"
)

const (
	// The maximum depth the engine will attempt to reach.
	MaxDepth = 100

	// A constant representing no move.
	NullMove Move = 0

	// A constant representing the score of the principal variation
	// move from the transposition table.
	PVMoveScore uint16 = 60

	// A constant representing the score offsets of the killer moves.
	FirstKillerMoveScore  uint16 = 10
	SecondKillerMoveScore uint16 = 20

	// A constant representing the maximum number of killers.
	MaxKillers = 2

	// A constant to offset the score of the pv and MVV-LVA move higher
	// than killers and history heuristic moves.
	MvvLvaOffset uint16 = math.MaxUint16 - 256

	// A constant representing the maximum value a history heuristic score
	// is allowed to reach. This ensures history scores stay well below
	// pv moves, captures, and killer moves.
	MaxHistoryScore int32 = int32(MvvLvaOffset - 30)

	// A constant representing the maximum game ply,
	// used to initalize the array for holding repetition
	// detection history.
	MaxGamePly = 1024

	// A constat representing a draw score value.
	Draw int16 = 0

	// Pruning constants
	NullMoveReduction               int8  = 2
	StaticNullMovePruningBaseMargin int16 = 120
	FutilityPruningBaseMargin       int16 = 200
	LateMoveReduction               int8  = 2
	LMRLegalMovesLimit              int   = 4
	LMRDepthLimit                   int8  = 3
	SingularMoveMargin              int16 = 250
)

// Futility margins
var FutilityMargins = [9]int16{
	0,
	200, // depth 1
	200, // depth 2
	300, // depth 3
	300, // depth 4
	550, // depth 5
	550, // depth 6
	800, // depth 7
	800, // depth 8

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
	TT    TransTable[SearchEntry]
	Timer TimeManager

	side              uint8
	nodes             uint64
	totalNodes        uint64
	killers           [MaxDepth + 1][MaxKillers]Move
	history           [2][64][64]int32
	zobristHistory    [MaxGamePly]uint64
	zobristHistoryPly uint16
}

// Setup the necessary internals of the engine when given a new FEN string.
func (search *Search) Setup(FEN string) {
	search.Pos.LoadFEN(FEN)
	search.zobristHistoryPly = 0
	search.zobristHistory[search.zobristHistoryPly] = search.Pos.Hash
}

// Reset the necessary internals of the engine. Normally used when the
// "ucinewgame" command is sent.
func (search *Search) Reset() {
	search.TT.Clear()
	search.ClearKillers()
	search.ClearHistoryTable()
}

// Add a zobrist hash to the history.
func (search *Search) AddHistory(hash uint64) {
	search.zobristHistoryPly++
	search.zobristHistory[search.zobristHistoryPly] = hash
}

// Remove a zobrist hash from the history.
func (search *Search) RemoveHistory() {
	search.zobristHistoryPly--
}

// The main search function for Blunder, implemented as an interative
// deepening loop.
func (search *Search) Search() Move {
	search.side = search.Pos.SideToMove
	pvLine := PVLine{}
	bestMove := NullMove

	search.totalNodes = 0
	depth := uint8(0)
	alpha := -Inf
	beta := Inf

	search.ageHistoryTable()
	search.Timer.Start(search.Pos.Ply)

	for depth = 1; depth <= MaxDepth &&
		depth <= search.Timer.MaxDepth &&
		search.Timer.MaxNodeCount > 0; depth++ {

		search.nodes = 0
		pvLine.Clear()

		startTime := time.Now()
		score := search.negamax(int8(depth), 0, alpha, beta, &pvLine, true, NullMove)
		endTime := time.Since(startTime)

		if search.Timer.Stop {
			if bestMove == NullMove && depth == 1 {
				bestMove = pvLine.GetPVMove()
			}
			break
		}

		bestMove = pvLine.GetPVMove()
		nps := uint64(float64(search.nodes) / float64(endTime.Seconds()))
		search.totalNodes += search.nodes

		fmt.Printf(
			"info depth %d score %s nodes %d nps %d time %d pv %s\n",
			depth, getMateOrCPScore(score),
			search.nodes, nps,
			endTime.Milliseconds(),
			pvLine,
		)
	}

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
func (search *Search) negamax(depth int8, ply uint8, alpha, beta int16, pvLine *PVLine, doNull bool, ignoreMove Move) int16 {
	// Update the number of nodes searched.
	search.nodes++

	if ply >= MaxDepth {
		return evaluatePos(&search.Pos)
	}

	// Make sure we haven't gone pass the node count limit.
	if search.totalNodes+search.nodes >= search.Timer.MaxNodeCount {
		search.Timer.Stop = true
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

	// Variables used throughout the rest of the search routine.
	isRoot := ply == 0
	inCheck := search.Pos.InCheck()
	isPVNode := beta-alpha != 1
	childPVLine := PVLine{}
	canFutilityPrune := false

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
		return search.qsearch(alpha, beta, ply, pvLine)
	}

	// Don't do any extra work if the current position is a draw. We
	// can just return a draw value.
	if !isRoot && (search.Pos.Rule50 >= 100 || search.isDrawByRepition()) {
		return search.contempt()
	}

	// Create a variable to store the possible best move we'll get from probing the transposition
	// table. And the best move we'll get from the search if we don't get a hit.
	ttMove := NullMove

	// =====================================================================//
	// TRANSPOSITION TABLE PROBING: Probe the transposition table to see if //
	// we have a useable matching entry for the current position. If we get //
	// a hit, return the score and stop searching.                          //
	// =====================================================================//

	entry := search.TT.Probe(search.Pos.Hash)
	ttScore, shouldUse := entry.Get(search.Pos.Hash, ply, uint8(depth), alpha, beta, &ttMove)

	if shouldUse && !isRoot && ignoreMove.Equal(NullMove) {
		return ttScore
	}

	// =====================================================================//
	// STATIC NULL MOVE PRUNING: If our current material score is so good   //
	// that even if we give ourselves a big hit materially and subtract a   //
	// large amount of our material score (the "score margin") and our      //
	// material score is still greater than beta, we assume this node will  //
	// fail-high and we can prune its branch.                               //
	// =====================================================================//

	if !inCheck && !isPVNode && abs(beta) < Checkmate {
		staticScore := evaluatePos(&search.Pos)
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

	if doNull && !inCheck && !isPVNode && depth >= NullMoveReduction+1 {
		search.Pos.DoNullMove()
		search.AddHistory(search.Pos.Hash)

		R := 3 + depth/6
		score := -search.negamax(depth-1-R, ply+1, -beta, -beta+1, &PVLine{}, false, NullMove)

		search.RemoveHistory()
		search.Pos.UndoNullMove()

		if search.Timer.Stop {
			return 0
		}

		if score >= beta && abs(score) < Checkmate {
			return beta
		}
	}

	// =====================================================================//
	// FUTILITY PRUNING: If we're close to the horizon, and even with a     //
	// large margin the static evaluation can't be raised above alpha,      //
	// we're probably in a fail-low node, and many moves can be probably    //
	// be pruned. So set a flag so we don't waste time searching moves that //
	// suck and probably don't even have a chance of raise alpha.           //
	// =====================================================================//

	if depth <= 8 && !isPVNode && !inCheck && alpha < Checkmate && beta < Checkmate {
		staticScore := evaluatePos(&search.Pos)
		margin := FutilityMargins[depth]
		canFutilityPrune = staticScore+margin <= alpha
	}

	moves := genMoves(&search.Pos)
	search.scoreMoves(&moves, ttMove, ply)

	legalMoves := 0
	ttFlag := AlphaFlag

	bestScore := -Inf
	bestMove := NullMove

	for index := uint8(0); index < moves.Count; index++ {
		orderMoves(index, &moves)
		move := moves.Moves[index]

		// Don't get caught in an infinite loop via extensions.
		if ignoreMove.Equal(move) {
			continue
		}

		if !search.Pos.DoMove(move) {
			search.Pos.UndoMove(move)
			continue
		}

		legalMoves++

		// Futility prune if possible

		if canFutilityPrune &&
			legalMoves > 1 &&
			!search.Pos.InCheck() &&
			move.MoveType() != Attack &&
			move.MoveType() != Promotion {

			search.Pos.UndoMove(move)
			continue
		}

		// Singular move extension
		extension := int8(0)

		if depth >= 8 &&
			!inCheck &&
			!isRoot &&
			ttMove.Equal(move) &&
			legalMoves == 1 &&
			entry.Depth >= uint8(depth-2) &&
			abs(entry.Score) < Checkmate &&
			entry.Flag != AlphaFlag {

			scoreToBeat := entry.Score - SingularMoveMargin - int16(depth*2)

			search.Pos.UndoMove(move)
			score := search.negamax(depth/2, ply, scoreToBeat-1, scoreToBeat, &PVLine{}, true, move)
			search.Pos.DoMove(move)

			if score < scoreToBeat {
				extension++
			} else if score >= beta {
				// Multi-cut pruning
				search.Pos.UndoMove(move)
				return score
			}
		}

		search.AddHistory(search.Pos.Hash)

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
			score = -search.negamax(depth-1+extension, ply+1, -beta, -alpha, &childPVLine, true, NullMove)
		} else {
			tactical := inCheck || move.MoveType() == Attack
			reduction := int8(0)

			if !isPVNode && legalMoves >= LMRLegalMovesLimit && depth >= LMRDepthLimit && !tactical {
				reduction = 2
			}

			score = -search.negamax(depth-1-reduction+extension, ply+1, -(alpha + 1), -alpha, &childPVLine, true, NullMove)

			if score > alpha && reduction > 0 {
				score = -search.negamax(depth-1-reduction+extension, ply+1, -beta, -alpha, &childPVLine, true, NullMove)
				if score > alpha {
					score = -search.negamax(depth-1+extension, ply+1, -beta, -alpha, &childPVLine, true, NullMove)
				}
			} else if score > alpha && score < beta {
				score = -search.negamax(depth-1+extension, ply+1, -beta, -alpha, &childPVLine, true, NullMove)
			}
		}

		search.Pos.UndoMove(move)
		search.RemoveHistory()

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		// If we have a beta-cutoff (i.e this move gives us a score better than what
		// our opponet can already guarantee early in the tree), return beta and the move
		// that caused the cutoff as the best move.
		if score >= beta {
			ttFlag = BetaFlag
			search.storeKiller(ply, move)
			search.incrementHistoryScore(move, depth)
			break
		} else {
			search.decrementHistoryScore(move)
		}

		// If the score of this move is better than alpha (i.e better than the score
		// we can currently guarantee), set alpha to be the score and the best move
		// to be the move that raised alpha.
		if score > alpha {
			alpha = score
			ttFlag = ExactFlag
			pvLine.Update(move, childPVLine)
			search.incrementHistoryScore(move, depth)
		} else {
			search.decrementHistoryScore(move)
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

	// If we're not out of time, and we haven't skipped moves
	// because of extensions, then store the result of the
	// search for this position.
	if !search.Timer.Stop && ignoreMove.Equal(NullMove) {
		search.TT.Probe(search.Pos.Hash).Set(
			search.Pos.Hash, ply, uint8(depth), bestScore, ttFlag, bestMove,
		)
	}

	// Return the best score.
	return bestScore
}

// Onece we reach a depth of zero in the main negamax search, instead of
// returning a static evaluation right away, continue to search deeper using
// a special form of negamax until the position is quiet (i.e there are no
// winning tatical captures). Doing this is known as quiescence search, and
// it makes the static evaluation much more accurate.
func (search *Search) qsearch(alpha, beta int16, maxPly uint8, pvLine *PVLine) int16 {
	search.nodes++

	if search.totalNodes+search.nodes >= search.Timer.MaxNodeCount {
		search.Timer.Stop = true
	}

	if (search.nodes & 2047) == 0 {
		search.Timer.Check()
	}

	if search.Timer.Stop {
		return 0
	}

	bestScore := evaluatePos(&search.Pos)

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
	search.scoreMoves(&moves, NullMove, maxPly)
	childPVLine := PVLine{}

	for index := uint8(0); index < moves.Count; index++ {
		orderMoves(index, &moves)
		move := moves.Moves[index]
		see := search.Pos.See(move)

		if see < 0 {
			continue
		}

		if !search.Pos.DoMove(move) {
			search.Pos.UndoMove(move)
			continue
		}

		score := -search.qsearch(-beta, -alpha, maxPly, &childPVLine)
		search.Pos.UndoMove(move)

		if score > bestScore {
			bestScore = score
		}

		if score >= beta {
			break
		}

		if score > alpha {
			alpha = score
			pvLine.Update(move, childPVLine)
		}

		childPVLine.Clear()
	}

	return bestScore
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

// Clear the killer moves table.
func (search *Search) ClearKillers() {
	for ply := 0; ply < MaxDepth+1; ply++ {
		search.killers[ply][0] = NullMove
		search.killers[ply][1] = NullMove
	}
}

// Increment the history score for the given move if it caused a beta-cutoff and is quiet.
func (search *Search) incrementHistoryScore(move Move, depth int8) {
	if search.Pos.Squares[move.ToSq()].Type == NoType {
		search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] += int32(depth) * int32(depth)
	}

	if search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] >= MaxHistoryScore {
		search.ageHistoryTable()
	}
}

// Decrement the history score for the given move if it didn't cause a beta-cutoff and is quiet.
func (search *Search) decrementHistoryScore(move Move) {
	if search.Pos.Squares[move.ToSq()].Type == NoType {
		if search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] > 0 {
			search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()] -= 1
		}
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

// Determine the draw score based on the phase of the game and whose moving,
// to encourge the engine to strive for a win in the middle-game, but be
// satisified with a draw in the endgame.
func (search *Search) contempt() int16 {
	return Draw
}

// Determine if the current board state is being repeated.
func (search *Search) isDrawByRepition() bool {
	for repPly := uint16(0); repPly < search.zobristHistoryPly; repPly++ {
		if search.zobristHistory[repPly] == search.Pos.Hash {
			return true
		}
	}
	return false
}

// Score the moves generated.
func (search *Search) scoreMoves(moves *MoveList, pvMove Move, ply uint8) {
	for index := uint8(0); index < moves.Count; index++ {
		move := &moves.Moves[index]
		capturedType := search.Pos.Squares[move.ToSq()].Type

		if move.Equal(pvMove) {
			move.AddScore(MvvLvaOffset + PVMoveScore)
		} else if capturedType != NoType {
			movedType := search.Pos.Squares[move.FromSq()].Type
			move.AddScore(MvvLvaOffset + MvvLva[capturedType][movedType])
		} else if move.Equal(search.killers[ply][0]) {
			move.AddScore(MvvLvaOffset - FirstKillerMoveScore)
		} else if move.Equal(search.killers[ply][1]) {
			move.AddScore(MvvLvaOffset - SecondKillerMoveScore)
		} else {
			move.AddScore(uint16(search.history[search.Pos.SideToMove][move.FromSq()][move.ToSq()]))
		}
	}
}

// Order the moves given by finding the best move and putting it
// at the index given.
func orderMoves(currIndex uint8, moves *MoveList) {
	bestIndex := currIndex
	bestScore := moves.Moves[bestIndex].Score()

	for index := bestIndex; index < moves.Count; index++ {
		if moves.Moves[index].Score() > bestScore {
			bestIndex = index
			bestScore = moves.Moves[index].Score()
		}
	}

	tempMove := moves.Moves[currIndex]
	moves.Moves[currIndex] = moves.Moves[bestIndex]
	moves.Moves[bestIndex] = tempMove
}
