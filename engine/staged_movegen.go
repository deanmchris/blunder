package engine

const (
	HashMoveStage uint8 = iota
	CapturesStage
	QuietsStage
)

type StagedMoveGenerator struct {
	search   *Search
	ttMove   uint32
	ply      uint8

	stage    uint8
	captures MoveList
	quiets   MoveList

	capturesIdx uint8
	quietsIdx   uint8
}

func (mg *StagedMoveGenerator) Next() uint32 {
	nextMove := NullMove

mainLoop:
	for {
	selectionLoop:
		switch mg.stage {
		case HashMoveStage:
			nextMove = mg.ttMove
			mg.stage = CapturesStage

			if equals(nextMove, NullMove) {
				break selectionLoop
			}

			break mainLoop
		case CapturesStage:
			if mg.captures.Count == 0 {
				mg.capturesIdx = 0
				mg.captures = genAttacks(&mg.search.Pos)
				scoreMoves(mg.search, &mg.captures, mg.ttMove, mg.ply)
			}

			if mg.capturesIdx == mg.captures.Count {
				mg.stage = QuietsStage
				break selectionLoop
			}

			swapBestMoveToIdx(&mg.captures, mg.capturesIdx)
			nextMove = mg.captures.Moves[mg.capturesIdx]
			mg.capturesIdx++

			if equals(nextMove, mg.ttMove) {
				break selectionLoop
			}

			break mainLoop
		case QuietsStage:
			if mg.quiets.Count == 0 {
				mg.quietsIdx = 0
				mg.quiets = genQuiets(&mg.search.Pos)
				scoreMoves(mg.search, &mg.quiets, mg.ttMove, mg.ply)
			}

			if mg.quietsIdx == mg.quiets.Count {
				break mainLoop
			}

			swapBestMoveToIdx(&mg.quiets, mg.quietsIdx)
			nextMove = mg.quiets.Moves[mg.quietsIdx]
			mg.quietsIdx++

			if equals(nextMove, mg.ttMove) {
				break selectionLoop
			}

			break mainLoop
		}
	}

	return nextMove
}
