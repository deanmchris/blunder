package engine

import (
	"fmt"
	"testing"
)

// Make sure to initalize the engine internals, since these
// tests are run independently of main.go.
func init() {
	InitBitboards()
	InitTables()
	InitZobrist()
}

type SeePosition struct {
	Fen   string
	Move  Move
	Score int16
}

var SeeTestPositions []SeePosition = []SeePosition{
	{
		"1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - - 0 1",
		NewMove(E1, E5, Attack, NoFlag), PieceValues[Pawn],
	},

	{
		"1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1",
		NewMove(D3, E5, Attack, NoFlag), PieceValues[Pawn] - PieceValues[Knight],
	},

	{
		"4q3/1p1pr1kb/1B2rp2/6p1/p3PP2/P3R1P1/1P2R1K1/4Q3 b - - 0 1",
		NewMove(H7, E4, Attack, NoFlag), PieceValues[Pawn],
	},
}

func TestSee(t *testing.T) {
	var pos Position

	for _, seePos := range SeeTestPositions {
		pos.LoadFEN(seePos.Fen)
		result := pos.See(seePos.Move)
		if result != seePos.Score {
			t.Error(
				fmt.Sprintf(
					"SEE test failed for position %s, move %s. Got %d instead of %d",
					seePos.Fen, seePos.Move, result, seePos.Score,
				),
			)
		}
	}
}
