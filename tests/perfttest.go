package tests

import (
	"blunder/engine"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	filePath := filepath.Join(parentFolder, "/tests/perftsuite.txt")

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
func printRow(fen, depth, expected, moves, correct string) {
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
func printRowSeparator() {
	fmt.Println("+" + strings.Repeat("-", 86) + "+" + "--------+------------+------------+----------+")
}

// Test blunder against the perft suite
func RunPerftTests(board *engine.Board) {
	printRowSeparator()
	printRow("position", "depth", "expected", "moves", "correct")
	printRowSeparator()

	perftTests := loadPerftSuite()
	var totalNodes uint64
	testsPassed := true
	start := time.Now()

	for _, perftTest := range perftTests {
		board.LoadFEN(perftTest.FEN)

		for depth, nodeCount := range perftTest.DepthValues {
			if nodeCount == 0 {
				continue
			}

			result := engine.Perft(board, depth+1, depth+1, true)
			totalNodes += result

			var correct string
			if nodeCount == result {
				correct = "yes"
			} else {
				testsPassed = false
				correct = "no"
			}

			printRow(
				perftTest.FEN,
				strconv.Itoa(depth+1),
				strconv.FormatInt(int64(nodeCount), 10),
				strconv.FormatInt(int64(result), 10),
				correct,
			)
			printRowSeparator()
		}
	}

	if testsPassed {
		fmt.Println("\nAll tests passed")
	} else {
		fmt.Println("\nTesting failed on some positions")
	}

	fmt.Println("\nTotal Nodes:", totalNodes)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())
	fmt.Printf("Nps: %d\n", int(float64(totalNodes)/elapsed.Seconds()))
}