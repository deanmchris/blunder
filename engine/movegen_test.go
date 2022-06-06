package engine

// movegen_test.go implements a file parser to read in test positions to ensure
// Blunder's move generator is working correctly.

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	MaxPerftDepth        = 7
	MaxFenStringLength   = 84
	MaxDepthNumberLength = 2
	MaxMoveNumberLength  = 10
)

type PerftTest struct {
	FEN         string
	DepthValues [MaxPerftDepth]uint64
}

// Load the perft test suite
func loadPerftSuite() (perftTests []PerftTest) {
	wd, _ := os.Getwd()
	parentFolder := filepath.Dir(wd)
	filePath := filepath.Join(parentFolder, "/perft_suite/perft_suite.epd")

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ";")
		perftTest := PerftTest{FEN: strings.TrimSpace(fields[0])}
		perftTest.DepthValues = [MaxPerftDepth]uint64{}

		for _, nodeCountStr := range fields[1:] {
			depth, err := strconv.Atoi(string(nodeCountStr[1]))
			if err != nil {
				panic(fmt.Sprintf("Parsing error on line: %s\n", line))
			}
			nodeCountStr = strings.TrimSpace(nodeCountStr[3:])
			nodeCount, err := strconv.Atoi(nodeCountStr)
			if err != nil {
				panic(fmt.Sprintf("Parsing error on line: %s\n", line))
			}
			perftTest.DepthValues[depth-1] = uint64(nodeCount)
		}
		perftTests = append(perftTests, perftTest)
	}
	return perftTests
}

// Pad a string into the center
func padToCenter(s string, fill string, w int) string {
	spaceLeft := w - len(s)
	extraFill := ""
	if spaceLeft%2 != 0 {
		extraFill = fill
	}
	return strings.Repeat(fill, spaceLeft/2) + extraFill + s + strings.Repeat(fill, spaceLeft/2)
}

// Print a row in the perft test output
func printPerftTestRow(fen, depth, expected, moves, correct string) {
	fmt.Printf(
		"| %s | %s | %s | %s | %s |\n",
		padToCenter(fen, " ", 84),
		padToCenter(depth, " ", 6),
		padToCenter(expected, " ", 10),
		padToCenter(moves, " ", 10),
		padToCenter(correct, " ", 8),
	)
}

// Print a row separator in the perft test output
func printPerftTestRowSeparator() {
	fmt.Println("+" + strings.Repeat("-", 86) + "+" + "--------+------------+------------+----------+")
}

// Test blunder against the perft suite
func TestMovegen(t *testing.T) {
	printPerftTestRowSeparator()
	printPerftTestRow("position", "depth", "expected", "moves", "correct")
	printPerftTestRowSeparator()

	pos := Position{}
	TT := TransTable[PerftEntry]{}
	totalNodes := uint64(0)
	testsPassed := true

	perftTests := loadPerftSuite()
	TT.Resize(DefaultTTSize, PerftEntrySize)
	start := time.Now()

	for _, perftTest := range perftTests {
		pos.LoadFEN(perftTest.FEN)

		for depth, nodeCount := range perftTest.DepthValues {
			if nodeCount == 0 {
				continue
			}

			result := Perft(&pos, uint8(depth)+1, &TT)
			totalNodes += result

			correct := ""
			if nodeCount == result {
				correct = "yes"
			} else {
				testsPassed = false
				correct = "no"
			}

			printPerftTestRow(
				perftTest.FEN,
				strconv.Itoa(depth+1),
				strconv.FormatInt(int64(nodeCount), 10),
				strconv.FormatInt(int64(result), 10),
				correct,
			)
			printPerftTestRowSeparator()
		}
	}

	if !testsPassed {
		t.Error("\nTesting failed on some positions. See table for exact positions.")
	}

	fmt.Println("\nTotal Nodes:", totalNodes)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())
	fmt.Printf("Nps: %d\n", int(float64(totalNodes)/elapsed.Seconds()))
}
