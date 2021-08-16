package engine

// movelist.go implements a very basic stack for holding moves

const MaxMoves = 255

type MoveList struct {
	Moves [MaxMoves]Move
	Count uint8
}

func (moveList *MoveList) AddMove(move Move) {
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}
