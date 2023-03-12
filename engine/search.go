package engine

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	MaxPly             uint8  = 100
	Infinity           int16  = 10000
	CheckmateThreshold int16  = 9000
	Draw               int16  = 0
	MaxPVLength        uint8  = 50
	MaxGamePly         uint16 = 700

	Buffer            uint16 = math.MaxUint16 - 200
	TTMoveScore       uint16 = 60
	FirstKillerScore  uint16 = 6
	SecondKillerScore uint16 = 4
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
	Pos   Position
	Timer TimeManager
	TT    TransTable[SearchEntry]

	Killers [MaxPly]KillerMovePair

	totalNodes        uint64
	zobristHistoryPly uint16
	zobristHistory    [MaxGamePly]uint64
}

func NewSearch(fen string) Search {
	search := Search{}
	search.Setup(fen)
	return search
}

func (search *Search) Setup(fen string) {
	search.Pos.LoadFEN(fen)
	search.TT.Resize(SearchTTSize)
	search.zobristHistoryPly = 0
	search.zobristHistory[0] = search.Pos.Hash
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

	search.Timer.Start()

	for depth := uint8(1); depth <= MaxPly && depth <= search.Timer.MaxDepth; depth++ {
		pv.Clear()

		startTime := time.Now()
		score := search.negamax(depth, 0, -Infinity, Infinity, &pv)
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

func (search *Search) negamax(depth, ply uint8, alpha, beta int16, pv *PVLine) int16 {
	search.totalNodes++

	if ply == MaxPly {
		return Evaluate(&search.Pos)
	}

	if depth == 0 {
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

	if search.Pos.InCheck {
		depth++
	}

	possibleMateInOne := search.Pos.InCheck && ply == 1
	isRoot := ply == 0
	if !isRoot && ((search.Pos.HalfMoveClock == 100 && !possibleMateInOne) || search.posIsDrawByRepition()) {
		return Draw
	}

	entry := search.TT.GetEntry(search.Pos.Hash)
	ttScore, ttMove, shouldUse := entry.GetScoreAndBestMove(search.Pos.Hash, ply, depth, alpha, beta)

	if !isRoot && shouldUse {
		return ttScore
	}

	moves := genAllMoves(&search.Pos)
	childPV := PVLine{}
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

		score := -search.negamax(depth-1, ply+1, -beta, -alpha, &childPV)

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
		entry := search.TT.GetEntry(search.Pos.Hash)
		entry.StoreNewInfo(search.Pos.Hash, bestMove, bestScore, depth, nodeType, ply)
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
