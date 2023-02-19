package engine

const MaxMoves = 255

type MoveList struct {
	Moves [MaxMoves]uint32
	Count uint8
}

func (moveList *MoveList) AddMove(move uint32) {
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}
