package extra

import (
	"blunder/engine"
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// iq.go is a program to measure Blunder's tatical strength, by trying to have it find
// agreed upon best moves in a variety of positions, under a certian time limit.
// The positions used can be found in testdata/tatical.epd and testdata/win_at_chess.epd.
// All credit goes to the creators.

var TestPositions []TestPosition

type TestPosition struct {
	Fen            string
	FirstBestMove  engine.Move
	SecondBestMove engine.Move
}

func loadTestPositions() {
	wd, _ := os.Getwd()
	parentFolder := filepath.Dir(wd)
	filePath := filepath.Join(parentFolder, "/testdata/win_at_chess.epd")

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
		if best[len(best)-1] == ';' {
			testPos.FirstBestMove = engine.ConvertSANToLAN(&pos, strings.TrimSuffix(best, ";"))
			testPos.SecondBestMove = engine.NullMove
		} else {
			testPos.FirstBestMove = engine.ConvertSANToLAN(&pos, best)
			testPos.SecondBestMove = engine.ConvertSANToLAN(&pos, strings.TrimSuffix(fields[6], ";"))
		}

		TestPositions = append(TestPositions, testPos)
	}
}

func TestIQ(timeAlloted int64) {
	loadTestPositions()

	var search engine.Search
	search.SpecifiedDepth = uint8(engine.MaxPly)
	search.SpecifiedNodes = uint64(math.MaxUint64)
	search.TT.Resize(engine.DefaultTTSize)

	total := len(TestPositions)
	failedPositions := []TestPosition{}

	correct := 0
	failed := 0

	for i, testPos := range TestPositions {
		if i > 0 && i%10 == 0 {
			fmt.Printf(
				"\nPERCENTAGE SCORE: %f (%d of out %d done)\n\n",
				float64(correct)/float64(i), i, len(TestPositions),
			)
		}

		search.Pos.LoadFEN(testPos.Fen)
		search.Timer.SetHardTimeForMove(timeAlloted * 1000)
		bestMove := search.Search()

		if testPos.FirstBestMove.Equal(engine.NullMove) && testPos.SecondBestMove.Equal(engine.NullMove) {
			panic("Invalid best move for position: " + testPos.Fen)
		}

		if bestMove.Equal(testPos.FirstBestMove) {
			fmt.Printf("%s BESTMOVE=%s (CORRECT)\n", testPos.Fen, testPos.FirstBestMove)
			correct++
		} else if bestMove.Equal(testPos.SecondBestMove) {
			fmt.Printf("%s BESTMOVE=%s (CORRECT)\n", testPos.Fen, testPos.SecondBestMove)
			correct++
		} else {
			if testPos.SecondBestMove.Equal(engine.NullMove) {
				fmt.Printf("%s BESTMOVE=%s (FAILED=%s)\n", testPos.Fen, testPos.FirstBestMove, bestMove)
			} else {
				fmt.Printf(
					"%s BESTMOVE=%s OR %s (FAILED=%s)\n",
					testPos.Fen,
					testPos.FirstBestMove, testPos.SecondBestMove,
					bestMove)
			}

			failedPositions = append(failedPositions, testPos)
			failed++
		}
	}

	fmt.Println("TOTAL POSITIONS:", total)
	fmt.Println("TOTAL CORRECT:", correct)
	fmt.Printf("FINAL PERCENTAGE SCORE: %f\n", float64(correct)/float64(total))
	fmt.Printf("FAILED POSITIONS (%d):\n", failed)
	for _, pos := range failedPositions {
		fmt.Println(pos)
	}
}
