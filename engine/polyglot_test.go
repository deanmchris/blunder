package engine

import (
	"testing"
)

// Make sure to initalize the engine internals, since these
// tests are run independently of main.go.
func init() {
	InitBitboards()
	InitTables()
	InitZobrist()
}

// polyglot_test.go provides tests to ensure polyglot hashing it working correctly.

type PolyglotPosition struct {
	Fen  string
	Hash uint64
}

var PolyglotTestPositions []PolyglotPosition = []PolyglotPosition{
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 0x463b96181691fc9c},
	{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", 0x823c9b50fd114196},
	{"rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 2", 0x0756b94461c50fb0},
	{"rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2", 0x662fafb965db29d4},
	{"rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3", 0x22a48b5a8e47ff78},
	{"rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR b kq - 0 3", 0x652a607ca3f242c1},
	{"rnbq1bnr/ppp1pkpp/8/3pPp2/8/8/PPPPKPPP/RNBQ1BNR w - - 0 4", 0x00fdd303c946bdd9},
	{"rnbqkbnr/p1pppppp/8/8/PpP4P/8/1P1PPPP1/RNBQKBNR b KQkq c3 0 3", 0x3c8123ea7b067637},
	{"rnbqkbnr/p1pppppp/8/8/P6P/R1p5/1P1PPPP1/1NBQKBNR b Kkq - 0 4", 0x5c3f9b829b279560},
}

func TestPolyglotHashing(t *testing.T) {
	var pos Position

	for _, polyglotPos := range PolyglotTestPositions {
		pos.LoadFEN(polyglotPos.Fen)
		hash := GenPolyglotHash(&pos)
		if hash != polyglotPos.Hash {
			t.Errorf(
				"Polyglot hash generation failed for position %s. Got 0x%x instead of 0x%x",
				polyglotPos.Fen, hash, polyglotPos.Hash,
			)
		}
	}
}
