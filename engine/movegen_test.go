package engine

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
	PerftTTSize = 10
)

type TestData struct {
	Fen               string
	Depth             uint8
	ExpectedNodeCount uint64
}

func loadTestData() (data []TestData) {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	parentFolder := filepath.Dir(wd)
	filePath := filepath.Join(parentFolder, "/perft_suite/perft_suite.epd")

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(bufio.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ";")

		fen := strings.TrimSpace(fields[0])

		for _, field := range fields[1:] {
			depth, err := strconv.Atoi(string(field[1]))
			if err != nil {
				panic(fmt.Sprintf("Parsing error on line: %s\n", line))
			}

			field = strings.TrimSpace(field[3:])
			nodeCount, err := strconv.Atoi(field)

			if err != nil {
				panic(fmt.Sprintf("Parsing error on line: %s\n", line))
			}

			data = append(data, TestData{fen, uint8(depth), uint64(nodeCount)})
		}
	}

	return data
}

func TestMovegen(t *testing.T) {
	InitBitboard()
	InitTables()
	InitMagics()
	InitZobrist()

	pos := Position{}
	tt := NewTransTable[PerftEntry](PerftTTSize)
	data := loadTestData()

	fmt.Println("Begin perft testing.")
	fmt.Println("Format of output: (<depth>) <fen> => <expected number>/<calculated number> [passed/failed]")
	fmt.Print("------------------------------------------------------------------------------\n\n")

	numberOfTests := 0
	numerOfFailedTests := 0
	totalNodes := 0
	startTime := time.Now()

	for _, testData := range data {
		numberOfTests++
		pos.LoadFEN(testData.Fen)

		nodes := Perft(&pos, testData.Depth, tt)
		totalNodes += int(nodes)

		if nodes == testData.ExpectedNodeCount {
			fmt.Printf("(%d) %s => %d/%d [passed]\n", testData.Depth, testData.Fen, testData.ExpectedNodeCount, nodes)
		} else {
			fmt.Printf("(%d) %s => %d/%d [failed]\n", testData.Depth, testData.Fen, testData.ExpectedNodeCount, nodes)
			numerOfFailedTests++
		}
	}

	if numerOfFailedTests != 0 {
		t.Error("\nTesting failed on some positions.")
	}

	endTime := time.Since(startTime)

	fmt.Printf("Testing completed in %d seconds.\n", endTime.Milliseconds())
	fmt.Printf("Nps: %d\n", int(float64(totalNodes)/endTime.Seconds()))
	fmt.Printf("%d tests were run, and %d were incorrect.\n", numberOfTests, numerOfFailedTests)
}
