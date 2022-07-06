package engine

import "fmt"

// staged_movegen.go implements a simple staged move generator for the search phase.

const (
	HashMoveStage uint8 = iota
	CapturesStage
	QuietsStage
)

// A struct to hold state needed to generate moves in stages.
type StagedMoveGenerator struct {
	search   *Search
	ttMove   Move
	prevMove Move
	ply      uint8

	stage    uint8
	captures MoveList
	quiets   MoveList

	capturesIdx uint8
	quietsIdx   uint8
}

// Generate the given moves in a position in a staged manner, in the hopes that we can get a
// beta cutoff without having to do the work of generate all possible moves, saving time.
func (mg *StagedMoveGenerator) Next() Move {
	nextMove := NullMove

mainLoop:
	for {
	selectionLoop:
		switch mg.stage {
		case HashMoveStage:
			nextMove = mg.ttMove
			mg.stage = CapturesStage

			if nextMove.Equal(NullMove) {
				break selectionLoop
			}

			break mainLoop
		case CapturesStage:
			if mg.captures.Count == 0 {
				mg.capturesIdx = 0
				mg.captures = genCaptures(&mg.search.Pos)
				mg.search.scoreMoves(&mg.captures, mg.ttMove, mg.ply, mg.prevMove)
			}

			if mg.capturesIdx == mg.captures.Count {
				mg.stage = QuietsStage
				break selectionLoop
			}

			orderMoves(mg.capturesIdx, &mg.captures)
			nextMove = mg.captures.Moves[mg.capturesIdx]
			mg.capturesIdx++

			if nextMove.Equal(mg.ttMove) {
				break selectionLoop
			}

			break mainLoop
		case QuietsStage:
			if mg.quiets.Count == 0 {
				mg.quietsIdx = 0
				mg.quiets = genQuiets(&mg.search.Pos)
				mg.search.scoreMoves(&mg.quiets, mg.ttMove, mg.ply, mg.prevMove)
			}

			if mg.quietsIdx == mg.quiets.Count {
				break mainLoop
			}

			orderMoves(mg.quietsIdx, &mg.quiets)
			nextMove = mg.quiets.Moves[mg.quietsIdx]
			mg.quietsIdx++

			if nextMove.Equal(mg.ttMove) {
				break selectionLoop
			}

			break mainLoop
		}
	}

	if nextMove.Equal(NullMove) && mg.stage != QuietsStage {
		fmt.Println(mg.stage)
		panic("Returning null move from staged move generator")
	}

	return nextMove
}
