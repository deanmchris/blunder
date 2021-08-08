package tests

import (
	"blunder/engine"
	"fmt"
	"os"
	"path/filepath"
)

const (
	StartingPositionHash uint64 = 0x463b96181691fc9c
	PlyLimit                    = engine.MaxHistory
)

// Print a row in the perft test output
func printZobristTestRow(filename, correct string) {
	fmt.Printf(
		"| %s | %s |\n",
		padToCenter(filename, " ", 10),
		padToCenter(correct, " ", 8),
	)
}

// Print a row separator in the perft test output
func printZobristTestRowSeparator() {
	fmt.Println("+" + "------------" + "+" + "----------+")
}

// Read a file at a time from the test files in test_books, and test them to
// verify Zobrist hashing is working correctly (see below).
func RunAllZobristHashingTests(board *engine.Board) {
	printZobristTestRowSeparator()
	printZobristTestRow("filename", "correct")
	printZobristTestRowSeparator()

	for fileNameSuffix := 1; fileNameSuffix < 12; fileNameSuffix++ {
		wd, _ := os.Getwd()
		parentFolder := filepath.Dir(wd)
		filePath := filepath.Join(parentFolder, fmt.Sprintf("/tests/polyglot_test_files/test%d.bin", fileNameSuffix))
		RunZobristHashingTest(board, filePath)
	}
	printZobristTestRowSeparator()
}

// To ensure zobrist hashing is working correctly, Blunder's polyglot
// reader is used to read in a polyglot file from a game played, and
// apply the moves which correspond to the current board Zobrist hash.
// If all the moves are applied successivley, then the hashing is working
// correctly. This is discovered by undoing each move and seeing if the board
// is returned to its correct beginning state, which is always the inital
// position.
func RunZobristHashingTest(board *engine.Board, path string) {
	entries, err := engine.LoadPolyglotFile(path)
	if err != nil {
		panic(err)
	}

	// Implement the 3-repititon rule for draws
	positionRepeats := make(map[uint64]int)
	var movesMade []engine.Move

	for {
		if entry, ok := entries[board.Hash]; ok {
			move := engine.MoveFromCoord(board, entry.Move, true)
			board.DoMove(move, true)
			movesMade = append(movesMade, move)
			positionRepeats[board.Hash]++
			if positionRepeats[board.Hash] == 3 {
				break
			}

			if board.GamePly == PlyLimit {
				break
			}
		} else {
			break
		}
	}

	for len(movesMade) != 0 {
		move := pop(&movesMade)
		board.UndoMove(move)

		_, ok := entries[board.Hash]
		if !ok {
			//invalid hash for position shown, so stop further
			//testing
			break
		}
	}

	correct := "yes"
	if board.Hash != StartingPositionHash || len(movesMade) != 0 {
		correct = "no"
	}
	printZobristTestRow(filepath.Base(path), correct)
}

// A helper function to pop and item from a slice
func pop(s *[]engine.Move) (item engine.Move) {
	item, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return item
}
