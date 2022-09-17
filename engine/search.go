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
	PVMoveScore uint16 = 65

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

	// Pruning constants
	NMR_Depth_Limit int8 = 2
)

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
	age               uint8
	totalNodes        uint64
	killers           [MaxDepth + 1][MaxKillers]Move
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
	search.age = 0
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
	search.totalNodes = 0
	search.age ^= 1

	pvLine := PVLine{}
	bestMove := NullMove

	totalTime := int64(0)
	depth := uint8(0)
	alpha := -Inf
	beta := Inf

	search.Timer.Start(search.Pos.Ply)

	for depth = 1; depth <= MaxDepth &&
		depth <= search.Timer.MaxDepth &&
		search.Timer.MaxNodeCount > 0; depth++ {

		pvLine.Clear()

		startTime := time.Now()
		score := search.negamax(int8(depth), 0, alpha, beta, &pvLine, true, NullMove, NullMove, false)
		endTime := time.Since(startTime)

		if search.Timer.Stop {
			if bestMove == NullMove && depth == 1 {
				bestMove = pvLine.GetPVMove()
			}
			break
		}

		totalTime += endTime.Milliseconds()

		bestMove = pvLine.GetPVMove()
		nps := uint64(float64(search.totalNodes*1000) / float64(totalTime))

		fmt.Printf(
			"info depth %d score %s nodes %d nps %d time %d pv %s\n",
			depth, getMateOrCPScore(score),
			search.totalNodes, nps,
			totalTime,
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
func (search *Search) negamax(depth int8, ply uint8, alpha, beta int16, pvLine *PVLine, doNull bool, prevMove, skipMove Move, isExtended bool) int16 {
	// Update the number of nodes searched.
	search.totalNodes++

	if ply >= MaxDepth {
		return NNEvaluatePos(&search.Pos)
	}

	// Make sure we haven't gone pass the node count limit.
	if search.totalNodes >= search.Timer.MaxNodeCount {
		search.Timer.Stop = true
	}

	// Every 2048 nodes, check if our time has expired.
	if (search.totalNodes & 2047) == 0 {
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
		search.totalNodes--
		return search.Qsearch(alpha, beta, ply, pvLine, ply)
	}

	// Don't do any extra work if the current position is a draw. We
	// can just return a draw value. We also need to check for the edge
	// case where there's a mate-in-one, but the fifty move rule counter
	// is at 99. Mate should always trumps the counter, so make sure we
	// don't return a draw evaluation for such a situation.
	possibleMateInOne := inCheck && ply == 1
	if !isRoot && ((search.Pos.Rule50 >= 100 && !possibleMateInOne) || search.isDrawByRepition()) {
		return Draw
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

	if shouldUse && !isRoot && !skipMove.Equal(ttMove) {
		return ttScore
	}

	// =====================================================================//
	// NULL MOVE PRUNING: If our opponet is given a free move, can they     //
	// improve their position? If we do a quick search after giving our     //
	// opponet this free move and we still find a move with a score better  //
	// than beta, our opponet can't improve their position and they         //
	// wouldn't take this path, so we have a beta cut-off and can prune     //
	// this branch.                                                         //
	// =====================================================================//

	if doNull && !inCheck && !isPVNode && depth >= NMR_Depth_Limit && !search.Pos.NoMajorsOrMiniors() {
		search.Pos.DoNullMove()
		search.AddHistory(search.Pos.Hash)

		R := int8(2)
		score := -search.negamax(depth-1-R, ply+1, -beta, -beta+1, &childPVLine, false, NullMove, NullMove, isExtended)

		search.RemoveHistory()
		search.Pos.UndoNullMove()
		childPVLine.Clear()

		if search.Timer.Stop {
			return 0
		}

		if score >= beta && abs(score) < Checkmate {
			return beta
		}
	}

	moves := genMoves(&search.Pos)
	search.scoreMoves(&moves, ttMove, ply, prevMove)

	legalMoves := 0
	ttFlag := AlphaFlag

	bestScore := -Inf
	bestMove := NullMove

	for index := uint8(0); index < moves.Count; index++ {
		orderMoves(index, &moves)
		move := moves.Moves[index]

		if !search.Pos.DoMove(move) {
			search.Pos.UndoMove(move)
			continue
		}

		legalMoves++

		search.AddHistory(search.Pos.Hash)

		score := -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine, true, move, NullMove, isExtended)

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
			break
		}

		// If the score of this move is better than alpha (i.e better than the score
		// we can currently guarantee), set alpha to be the score and the best move
		// to be the move that raised alpha.
		if score > alpha {
			alpha = score
			ttFlag = ExactFlag
			pvLine.Update(move, childPVLine)
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
		return Draw
	}

	// If we're not out of time, store the result of the search for this position.
	if !search.Timer.Stop {
		entry := search.TT.Probe(search.Pos.Hash)
		entry.Set(search.Pos.Hash, bestScore, bestMove, ply, uint8(depth), ttFlag)
	}

	// Return the best score.
	return bestScore
}

// Onece we reach a depth of zero in the main negamax search, instead of
// returning a static evaluation right away, continue to search deeper using
// a special form of negamax until the position is quiet (i.e there are no
// winning tatical captures). Doing this is known as quiescence search, and
// it makes the static evaluation much more accurate.
func (search *Search) Qsearch(alpha, beta int16, maxPly uint8, pvLine *PVLine, ply uint8) int16 {
	search.totalNodes++

	if maxPly+ply >= MaxDepth {
		return NNEvaluatePos(&search.Pos)
	}

	if search.totalNodes >= search.Timer.MaxNodeCount {
		search.Timer.Stop = true
	}

	if (search.totalNodes & 2047) == 0 {
		search.Timer.Check()
	}

	if search.Timer.Stop {
		return 0
	}

	bestScore := NNEvaluatePos(&search.Pos)

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

	moves := genCapturesAndQueenPromotions(&search.Pos)
	search.scoreMoves(&moves, NullMove, maxPly, NullMove)
	childPVLine := PVLine{}

	for index := uint8(0); index < moves.Count; index++ {
		orderMoves(index, &moves)
		move := moves.Moves[index]

		if !search.Pos.DoMove(move) {
			search.Pos.UndoMove(move)
			continue
		}

		score := -search.Qsearch(-beta, -alpha, maxPly, &childPVLine, ply+1)
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
func (search *Search) scoreMoves(moves *MoveList, pvMove Move, ply uint8, prevMove Move) {
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
