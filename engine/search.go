package engine

import (
	"fmt"
	"time"
)

const (
	MaxDepth      = 50
	NullMove Move = 0

	KillerMoveScore int16 = 10
	PVMoveScore     int16 = 60
)

var MvvLva [7][6]int16 = [7][6]int16{
	{16, 15, 14, 13, 12, 11}, // victim Pawn
	{26, 25, 24, 23, 22, 21}, // victim Knight
	{36, 35, 34, 33, 32, 31}, // victim Bishop
	{46, 45, 44, 43, 42, 41}, // vitcim Rook
	{56, 55, 54, 53, 52, 51}, // victim Queen

	{0, 0, 0, 0, 0, 0}, // victim King
	{0, 0, 0, 0, 0, 0}, // No piece
}

type PVLine struct {
	moves []Move
}

func (pvLine *PVLine) Clear() {
	pvLine.moves = nil
}

func (pvLine *PVLine) Update(move Move, newPVLine PVLine) {
	pvLine.Clear()
	pvLine.moves = append(pvLine.moves, move)
	pvLine.moves = append(pvLine.moves, newPVLine.moves...)
}

func (pvLine *PVLine) GetPVMove() Move {
	return pvLine.moves[0]
}

func (pvLine PVLine) String() string {
	pv := fmt.Sprintf("%s", pvLine.moves)
	return pv[1 : len(pv)-1]
}

type Search struct {
	Pos   Position
	Timer TimeManager
	TT    TransTable

	side    uint8
	nodes   uint64
	killers [MaxDepth][2]Move
}

func (search *Search) Search() Move {
	search.side = search.Pos.SideToMove
	var pvLine PVLine
	bestMove := NullMove

	search.Timer.Start()
	for depth := 1; depth <= MaxDepth; depth++ {

		pvLine.Clear()
		startTime := time.Now()
		score := search.negamax(uint8(depth), 0, -Inf, Inf, &pvLine)
		endTime := time.Since(startTime)

		if search.Timer.Stop {
			if bestMove == NullMove && depth == 1 {
				bestMove = pvLine.GetPVMove()
			}
			break
		}

		bestMove = pvLine.GetPVMove()
		fmt.Printf(
			"info depth %d score cp %d time %d nodes %d\n",
			depth, score,
			endTime.Milliseconds(),
			search.nodes,
			// pvLine,
		)
	}
	return bestMove
}

func (search *Search) negamax(depth, ply uint8, alpha, beta int16, pvLine *PVLine) int16 {
	if (search.nodes&2047) == 0 && search.Timer.Check() {
		return 0
	}

	search.nodes++

	isRoot := ply == 0
	inCheck := search.Pos.InCheck()

	if inCheck {
		depth++
	}

	if depth == 0 {
		return search.qsearch(alpha, beta, ply)
	}

	if !isRoot && (search.Pos.Rule50 >= 100 || search.isDrawByRepition()) {
		return search.contempt()
	}

	ttBestMove := NullMove
	score := search.TT.Probe(search.Pos.Hash, ply, depth, alpha, beta, &ttBestMove)
	if score != Invalid && !isRoot {
		return score
	}

	moves := genMoves(&search.Pos)
	search.scoreMoves(&moves, ply, ttBestMove)

	var childPVLine PVLine
	doPVS := false
	legalMoves := 0

	ttFlag := AlphaFlag
	ttBestMove = NullMove

	for index := 0; index < int(moves.Count); index++ {
		orderMoves(index, &moves)
		move := moves.Moves[index]

		if !search.Pos.MakeMove(move) {
			search.Pos.UnmakeMove(move)
			continue
		}

		var score int16
		if doPVS {
			score = -search.negamax(depth-1, ply+1, -alpha-1, -alpha, &childPVLine)
			if score > alpha && score < beta {
				score = -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine)
			}
		} else {
			score = -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine)
		}

		/*if doPVS {
		}
		score := -search.negamax(depth-1, ply+1, -beta, -alpha, &childPVLine)*/
		search.Pos.UnmakeMove(move)
		legalMoves++

		if score >= beta {
			alpha = beta
			search.storeKiller(ply, move)
			ttFlag = BetaFlag
			ttBestMove = move
			break
		}

		if score > alpha {
			alpha = score
			pvLine.Update(move, childPVLine)

			doPVS = true
			ttFlag = ExactFlag
			ttBestMove = move
		}
	}

	if legalMoves == 0 {
		if inCheck {
			return -Inf + int16(ply)
		}
		return search.contempt()
	}

	if !search.Timer.Check() {
		search.TT.Store(search.Pos.Hash, ply, depth, alpha, ttFlag, ttBestMove)
	}

	return alpha
}

func (search *Search) qsearch(alpha, beta int16, ply uint8) int16 {
	if (search.nodes&2047) == 0 && search.Timer.Check() {
		return 0
	}

	search.nodes++

	if search.Pos.Rule50 >= 100 || search.isDrawByRepition() {
		return search.contempt()
	}

	staticScore := evaluatePos(&search.Pos)

	if staticScore >= beta {
		return beta
	}

	if staticScore > alpha {
		alpha = staticScore
	}

	moves := genMoves(&search.Pos)
	search.scoreMoves(&moves, ply, NullMove)

	for index := 0; index < int(moves.Count); index++ {
		if moves.Moves[index].MoveType() == Attack {
			orderMoves(index, &moves)
			move := moves.Moves[index]

			if !search.Pos.MakeMove(move) {
				search.Pos.UnmakeMove(move)
				continue
			}

			score := -search.qsearch(-beta, -alpha, ply)
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

func (search *Search) storeKiller(ply uint8, move Move) {
	if search.Pos.Squares[move.ToSq()].Type == NoType {
		if !move.Equal(search.killers[ply][0]) {
			search.killers[ply][1] = search.killers[ply][0]
			search.killers[ply][0] = move
		}
	}
}

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

func (search *Search) isDrawByRepition() bool {
	var repPly uint16
	for repPly = 0; repPly < search.Pos.HistoryPly; repPly++ {
		if search.Pos.History[repPly] == search.Pos.Hash {
			return true
		}
	}
	return false
}

func (search *Search) scoreMoves(moves *MoveList, ply uint8, pvMove Move) {
	for index := 0; index < int(moves.Count); index++ {
		move := &moves.Moves[index]
		captured := &search.Pos.Squares[move.ToSq()]
		if pvMove.Equal(*move) {
			move.AddScore(PVMoveScore)
		} else if captured.Type != NoType {
			moved := &search.Pos.Squares[move.FromSq()]
			move.AddScore(MvvLva[captured.Type][moved.Type])
		} else {
			if search.killers[ply][0].Equal(*move) {
				move.AddScore(KillerMoveScore)
			} else if search.killers[ply][1].Equal(*move) {
				move.AddScore(KillerMoveScore)
			}
		}
	}
}

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
