package engine

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	MaxPly             int8   = 100
	Infinity           int16  = 10000
	CheckmateThreshold int16  = 9000
	Draw               int16  = 0
	MaxPVLength        uint8  = 50
	MaxGamePly         uint16 = 700

	Buffer            uint16 = math.MaxUint16 - 200
	TTMoveScore       uint16 = 60
	FirstKillerScore  uint16 = 6
	SecondKillerScore uint16 = 4

	// Pruning constants

	NMP_Depth_Limit int8 = 2
)

var MVV_LVA = [7][6]uint16{
	{16, 15, 14, 13, 12, 11}, // victim Pawn
	{26, 25, 24, 23, 22, 21}, // victim Knight
	{36, 35, 34, 33, 32, 31}, // victim Bishop
	{46, 45, 44, 43, 42, 41}, // vitcim Rook
	{56, 55, 54, 53, 52, 51}, // victim Queen
	{0, 0, 0, 0, 0, 0},       // victim King
	{0, 0, 0, 0, 0, 0},       // victom no piece
}

type PVLine struct {
	Moves []uint32
}

func (pvLine *PVLine) Clear() {
	pvLine.Moves = nil
}

func (pvLine *PVLine) Update(move uint32, newPVLine PVLine) {
	pvLine.Clear()
	pvLine.Moves = append(pvLine.Moves, move)
	pvLine.Moves = append(pvLine.Moves, newPVLine.Moves...)
}

func (pvLine *PVLine) GetBestMove() uint32 {
	return pvLine.Moves[0]
}

func (pvLine PVLine) String() string {
	pv := strings.Builder{}
	for i := 0; i < len(pvLine.Moves); i++ {
		move := pvLine.Moves[i]
		if move == NullMove {
			break
		}
		pv.WriteString(moveToStr(move))
		pv.WriteString(" ")
	}
	return pv.String()
}

type KillerMovePair struct {
	FirstKiller,
	SecondKiller uint32
}

func (pair *KillerMovePair) AddKillerMove(newKillerMove uint32) {
	if !equals(pair.FirstKiller, newKillerMove) {
		pair.SecondKiller = pair.FirstKiller
		pair.FirstKiller = newKillerMove
	}
}

type Search struct {
	Pos     Position
	Timer   TimeManager
	TT      TransTable[SearchBucket]
	Killers [MaxPly]KillerMovePair

	totalNodes        uint64
	zobristHistoryPly uint16
	zobristHistory    [MaxGamePly]uint64

	ageCounter uint16
	age        uint8
}

func NewSearch(fen string) Search {
	search := Search{}
	search.LoadFEN(fen)
	search.TT.Resize(SearchTTSize)
	return search
}

func (search *Search) ResetInternals(fen string) {
	search.LoadFEN(fen)
	search.TT.Clear()

	for i := range search.Killers {
		search.Killers[i].FirstKiller = NullMove
		search.Killers[i].SecondKiller = NullMove
	}
}

func (search *Search) LoadFEN(fen string) {
	search.Pos.LoadFEN(fen)
	search.zobristHistoryPly = 0
	search.zobristHistory[0] = search.Pos.Hash
	search.ageCounter = 0
	search.age = 0
}

func (search *Search) StopSearch() {
	search.Timer.ForceStop()
}

func (search *Search) AddHistory(hash uint64) {
	search.zobristHistoryPly++
	search.zobristHistory[search.zobristHistoryPly] = hash
}
func (search *Search) RemoveHistory() {
	search.zobristHistoryPly--
}

func (search *Search) RunSearch() uint32 {
	pv := PVLine{}
	bestMove := NullMove
	totalTime := int64(0)
	search.totalNodes = 0

	search.age = uint8(search.ageCounter % 16)
	search.ageCounter += 1

	fmt.Println(search.ageCounter, search.age)

	search.Timer.Start()

	for depth := int8(1); depth <= MaxPly && depth <= search.Timer.MaxDepth; depth++ {
		pv.Clear()

		startTime := time.Now()
		score := search.negamax(depth, 0, -Infinity, Infinity, &pv, false)
		endTime := time.Since(startTime)

		if search.Timer.IsStopped() {
			if bestMove == NullMove && depth == 1 {
				bestMove = pv.GetBestMove()
			}
			break
		}

		totalTime += endTime.Milliseconds() + 1

		bestMove = pv.GetBestMove()
		nps := (search.totalNodes * 1000) / uint64(totalTime)

		fmt.Printf(
			"info depth %d score %s nodes %d nps %d time %d pv %s\n",
			depth, getMateOrCPScore(score),
			search.totalNodes, nps,
			totalTime, pv,
		)
	}

	return bestMove
}

func (search *Search) negamax(depth int8, ply uint8, alpha, beta int16, pv *PVLine, doNMP bool) int16 {
	search.totalNodes++

	if ply == uint8(MaxPly) {
		return Evaluate(&search.Pos)
	}

	if depth <= 0 {
		search.totalNodes--
		return search.QuiescenceSearch(alpha, beta, pv)
	}

	if search.totalNodes >= search.Timer.MaxNodeCount {
		search.Timer.ForceStop()
	}

	if (search.totalNodes & 2047) == 0 {
		search.Timer.CheckIfTimeIsUp()
	}

	if search.Timer.IsStopped() {
		return 0
	}

	search.Pos.ComputePinAndCheckInfo()

	isRoot := ply == 0
	childPV := PVLine{}

	// =====================================================================//
	// CHECK EXTENSION: Extend the search depth by one if we're in check,   //
	// so that we're less likely to push danger over the search horizon,    //
	// and we won't enter quiescence search while in check.                 //
	// =====================================================================//

	if search.Pos.InCheck {
		depth++
	}

	possibleMateInOne := search.Pos.InCheck && ply == 1
	if !isRoot && ((search.Pos.HalfMoveClock == 100 && !possibleMateInOne) || search.posIsDrawByRepition()) {
		return Draw
	}

	// =====================================================================//
	// TRANSPOSITION TABLE PROBING: Probe the transposition table to see if //
	// we have a useable matching entry for the current position. If we get //
	// a hit, return the score and stop searching.                          //
	// =====================================================================//

	bucket := search.TT.GetBucket(search.Pos.Hash)
	entry := bucket.GetEntryForProbing(search.Pos.Hash, search.age)
	ttScore, ttMove, shouldUse := entry.GetScoreAndBestMove(search.Pos.Hash, ply, depth, alpha, beta)

	if !isRoot && shouldUse {
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

	if doNMP && !search.Pos.InCheck && depth >= NMP_Depth_Limit && !search.Pos.NoMajorsOrMiniors() {
		search.Pos.DoNullMove()
		search.AddHistory(search.Pos.Hash)

		R := int8(2)
		score := -search.negamax(depth-1-R, ply+1, -beta, -beta+1, &childPV, false)

		search.RemoveHistory()
		search.Pos.UndoNullMove()
		childPV.Clear()

		if search.Timer.IsStopped() {
			return 0
		}

		if score >= beta && Abs(score) < CheckmateThreshold {
			return beta
		}
	}

	moves := genAllMoves(&search.Pos)
	bestScore := -Infinity
	numLegalMoves := uint8(0)
	nodeType := FailLowNode
	bestMove := NullMove

	scoreMoves(search, &moves, ttMove, ply)

	for i := uint8(0); i < moves.Count; i++ {
		swapBestMoveToIdx(&moves, i)
		move := moves.Moves[i]

		if !search.Pos.DoMove(move) {
			search.Pos.UndoMove(move)
			continue
		}

		search.AddHistory(search.Pos.Hash)
		numLegalMoves++

		// =====================================================================//
		// PRINCIPAL VARIATION SEARCH: Due to good move ordering, the first     //
		// move we search is likely the best move, and all the remaining moves  //
		// will end up failing-low. To prove this cheaply we search all moves   //
		// after the first move with a null window cenetered around alpha. Thus //
		// beta-cutoffs will occur in children nodes of the same color if a     //
		// value better than alpha is found, since we only care if alpha is     //
		// beaten, not by how much. If beta-cutoffs occur in all of our         //
		// immediate children of the same color, then we've actually found a    //
		// better move than the first and need to research with a full-window   //
		// to get it's exact value. But most of the time we're correct, the     //
		// first move is the best, and we save time verifying this by searching //
		// non-first moves with a null-window.                                  //
		// =====================================================================//

		score := int16(0)
		if numLegalMoves == 1 {
			score = -search.negamax(depth-1, ply+1, -beta, -alpha, &childPV, false)
		} else {
			score = -search.negamax(depth-1, ply+1, -alpha-1, -alpha, &childPV, true)
			if score > alpha && score < beta {
				score = -search.negamax(depth-1, ply+1, -beta, -alpha, &childPV, true)
			}
		}

		search.Pos.UndoMove(move)
		search.RemoveHistory()

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if bestScore >= beta {
			nodeType = FailHighNode
			if search.Pos.GetPieceType(toSq(move)) == NoType {
				search.Killers[ply].AddKillerMove(move)
			}
			break
		}

		if bestScore > alpha {
			alpha = bestScore
			nodeType = PVNode
			pv.Update(move, childPV)
		}

		childPV.Clear()
	}

	if numLegalMoves == 0 {
		if search.Pos.InCheck {
			return -Infinity + int16(ply)
		}
		return Draw
	}

	if !search.Timer.IsStopped() {
		entry := bucket.GetEntryForStoring(search.Pos.Hash, search.age)
		entry.StoreNewInfo(search.Pos.Hash, bestMove, bestScore, depth, nodeType, ply, search.age)
	}

	return bestScore
}

func (search *Search) QuiescenceSearch(alpha, beta int16, pv *PVLine) int16 {
	search.totalNodes++

	if search.totalNodes >= search.Timer.MaxNodeCount {
		search.Timer.ForceStop()
	}

	if (search.totalNodes & 2047) == 0 {
		search.Timer.CheckIfTimeIsUp()
	}

	if search.Timer.IsStopped() {
		return 0
	}

	bestScore := Evaluate(&search.Pos)

	if bestScore >= beta {
		return bestScore
	}

	if bestScore > alpha {
		alpha = bestScore
	}

	search.Pos.ComputePinAndCheckInfo()

	moves := genAttacks(&search.Pos)
	childPV := PVLine{}

	scoreMoves(search, &moves, NullMove, 0)

	for i := uint8(0); i < moves.Count; i++ {
		swapBestMoveToIdx(&moves, i)
		move := moves.Moves[i]

		if !search.Pos.DoMove(move) {
			search.Pos.UndoMove(move)
			continue
		}

		search.AddHistory(search.Pos.Hash)

		score := -search.QuiescenceSearch(-beta, -alpha, &childPV)

		search.Pos.UndoMove(move)
		search.RemoveHistory()

		if score > bestScore {
			bestScore = score
		}

		if bestScore >= beta {
			break
		}

		if bestScore > alpha {
			alpha = bestScore
			pv.Update(move, childPV)
		}

		childPV.Clear()
	}

	return bestScore
}

func (search *Search) posIsDrawByRepition() bool {
	for i := uint16(0); i < search.zobristHistoryPly; i++ {
		if search.zobristHistory[i] == search.Pos.Hash {
			return true
		}
	}
	return false
}

func scoreMoves(search *Search, moves *MoveList, ttMove uint32, ply uint8) {
	for i := uint8(0); i < moves.Count; i++ {
		move := &moves.Moves[i]

		if equals(*move, ttMove) {
			addScore(move, Buffer+TTMoveScore)
		} else if search.Pos.GetPieceType(toSq(*move)) != NoType {
			from, to := fromSq(*move), toSq(*move)
			attackerType, attackedType := search.Pos.GetPieceType(from), search.Pos.GetPieceType(to)
			addScore(move, Buffer+MVV_LVA[attackedType][attackerType])
		} else {
			if equals(search.Killers[ply].FirstKiller, *move) {
				addScore(move, Buffer+FirstKillerScore)
			} else if equals(search.Killers[ply].SecondKiller, *move) {
				addScore(move, Buffer+SecondKillerScore)
			}
		}
	}
}

func swapBestMoveToIdx(moves *MoveList, index uint8) {
	bestMoveScore := score(moves.Moves[index])
	bestMoveIndex := index

	for i := index + 1; i < moves.Count; i++ {
		moveScore := score(moves.Moves[i])
		if moveScore > bestMoveScore {
			bestMoveScore = moveScore
			bestMoveIndex = i
		}
	}

	if bestMoveIndex != index {
		bestMove := moves.Moves[bestMoveIndex]
		moves.Moves[bestMoveIndex] = moves.Moves[index]
		moves.Moves[index] = bestMove
	}
}

func getMateOrCPScore(score int16) string {
	if score > CheckmateThreshold {
		pliesToMate := Infinity - score
		mateInN := (pliesToMate / 2) + (pliesToMate % 2)
		return fmt.Sprintf("mate %d", mateInN)
	}

	if score < -CheckmateThreshold {
		pliesToMate := -Infinity - score
		mateInN := (pliesToMate / 2) + (pliesToMate % 2)
		return fmt.Sprintf("mate %d", mateInN)
	}

	return fmt.Sprintf("cp %d", score)
}
