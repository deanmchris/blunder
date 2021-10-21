package extra

import (
	"blunder/engine"
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// iq.go is a program to measure Blunder's tatical strength, by trying to have it find
// agreed upon best moves in a variety of positions, under a certian time limit.
// The positions used can be found in testdata/tatical.epd, along with credit to their
// creators.

var TestPositions []TestPosition

var CharToPieceType map[rune]uint8 = map[rune]uint8{
	'N': engine.Knight,
	'B': engine.Bishop,
	'R': engine.Rook,
	'Q': engine.Queen,
	'K': engine.King,
}

// An object representing a test position, and the best move
// in the position.
type TestPosition struct {
	Fen      string
	BestMove engine.Move
}

// Convert a move in short algebraic notation, to the long algebraic notation used
// by the UCI protocol.
func convertSANToLAN(pos *engine.Position, move string) engine.Move {
	coords := ""
	pieceType := engine.Pawn

	for _, char := range move {
		switch char {
		case 'N', 'B', 'R', 'Q', 'K':
			pieceType = CharToPieceType[char]
		case '1', '2', '3', '4', '5', '6', '7', '8':
			coords += string(char)
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
			coords += string(char)
		}
	}

	moves := engine.GenMoves(pos)
	matchingMove := engine.NullMove

	for _, move := range moves.Moves {
		if len(coords) == 2 {
			moved := pos.Squares[move.FromSq()].Type
			toSq := engine.CoordinateToPos(coords)
			if toSq == move.ToSq() && pieceType == moved {
				matchingMove = move
				break
			}
		} else if len(coords) == 3 {
			toSq := engine.CoordinateToPos(coords[1:])
			fileOrRank := coords[0]
			moveCoords := move.String()

			if unicode.IsLetter(rune(fileOrRank)) {
				if toSq == move.ToSq() && fileOrRank == moveCoords[0] {
					matchingMove = move
					break
				}
			} else {
				if toSq == move.ToSq() && fileOrRank == moveCoords[1] {
					matchingMove = move
					break
				}
			}
		} else if len(coords) == 4 {
			toSq := engine.CoordinateToPos(coords[0:2])
			fromSq := engine.CoordinateToPos(coords[2:4])
			if toSq == move.ToSq() && fromSq == move.FromSq() {
				matchingMove = move
				break
			}
		}
	}

	return matchingMove
}

// Load the test positions from the file in testdata
func loadTestPositions() {
	wd, _ := os.Getwd()
	parentFolder := filepath.Dir(wd)
	filePath := filepath.Join(parentFolder, "/testdata/tactical.epd")

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	var pos engine.Position

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		var testPos TestPosition
		testPos.Fen = fields[0] + " " + fields[1] + " " + fields[2] + " " + fields[3] + " 0 1"

		best := fields[5]
		pos.LoadFEN(testPos.Fen)
		testPos.BestMove = convertSANToLAN(&pos, strings.TrimSuffix(best, ";"))

		TestPositions = append(TestPositions, testPos)
	}
}

// Test the "iq" of Blunder by testing it's ability to find the best move in a given position
// within a certian amount of time.
func TestIQ(timeAlloted int64) {
	loadTestPositions()

	var search engine.Search
	search.Timer.TimeLeft = timeAlloted * 1000 * 40
	search.SpecifiedDepth = uint8(engine.MaxPly)
	search.SpecifiedNodes = uint64(math.MaxUint64)
	search.TT.Resize(engine.DefaultTTSize)

	total := len(TestPositions)
	correct := 0

	for i, testPos := range TestPositions {
		search.Pos.LoadFEN(testPos.Fen)
		bestMove := search.Search()

		if testPos.BestMove.Equal(engine.NullMove) {
			panic("Invalid best move for position: " + testPos.Fen)
		}

		if bestMove.Equal(testPos.BestMove) {
			fmt.Printf("%s BESTMOVE=%s (CORRECT)\n", testPos.Fen, testPos.BestMove)
			correct++
		} else {
			fmt.Printf("%s BESTMOVE=%s (FAILED=%s)\n", testPos.Fen, testPos.BestMove, bestMove)
		}

		if i > 0 && i%10 == 0 {
			fmt.Printf("\nPERCENTAGE SCORE: %f\n\n", float64(correct)/float64(i))
		}
	}

	fmt.Println("TOTAL POSITIONS:", total)
	fmt.Println("TOTAL CORRECT:", correct)
	fmt.Printf("FINAL PERCENTAGE SCORE: %f\n", float64(correct)/float64(total))
}
