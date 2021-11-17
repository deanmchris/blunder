package tuner

// tuner.go is a texel tuning implementation for Blunder.

import (
	"blunder/engine"
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	DataFile   = ""
	NumCores   = 4
	NumWeights = 786
	KPrecision = 10

	Draw         float64 = 0.5
	WhiteWin     float64 = 1.0
	BlackWin     float64 = 0.0
	NumPositions float64 = 400000.0
	K            float64 = 1.62
)

// A struct object to hold data concering a position loaded from the training file.
// Each position consists of a position board object and the outcome of the game
// the position was from.
type Position struct {
	Pos     engine.Position
	Outcome float64
}

// A global variable to hold the positions loaded from the training file.
var Positions []Position

// A global variable to hold the parallel computations of the MSE function.
var Answers = make(chan float64)

// A method to specifiy which weights should be ignored when tuning.
var IgnoreWeights = make([]bool, len(Weights))

func setIgnoredWeights(from, to int) {
	for i := from; i < to; i++ {
		IgnoreWeights[i] = true
	}
}

// The weights to be adjusted during the tuning process.
var Weights []int16 = loadWeights()

// Load the weights for tuning from the current evaluation terms.
func loadWeights() (weights []int16) {
	weights = make([]int16, NumWeights)
	copy(weights[0:64], engine.PSQT_MG[engine.Pawn][:])
	copy(weights[64:128], engine.PSQT_MG[engine.Knight][:])
	copy(weights[128:192], engine.PSQT_MG[engine.Bishop][:])
	copy(weights[192:256], engine.PSQT_MG[engine.Rook][:])
	copy(weights[256:320], engine.PSQT_MG[engine.Queen][:])
	copy(weights[320:384], engine.PSQT_MG[engine.King][:])

	copy(weights[384:448], engine.PSQT_EG[engine.Pawn][:])
	copy(weights[448:512], engine.PSQT_EG[engine.Knight][:])
	copy(weights[512:576], engine.PSQT_EG[engine.Bishop][:])
	copy(weights[576:640], engine.PSQT_EG[engine.Rook][:])
	copy(weights[640:704], engine.PSQT_EG[engine.Queen][:])
	copy(weights[704:768], engine.PSQT_EG[engine.King][:])

	copy(weights[768:773], engine.PieceValueMG[:])
	copy(weights[773:778], engine.PieceValueEG[:])
	copy(weights[778:782], engine.PieceMobilityMG[:])
	copy(weights[782:786], engine.PieceMobilityEG[:])

	return weights
}

// Load the given number of positions from the training set file.
func loadPositions() (positions []Position) {
	file, err := os.Open(DataFile)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for positionCount := 0; scanner.Scan() && positionCount < int(NumPositions); positionCount++ {
		line := scanner.Text()
		fields := strings.Fields(line)

		fen := fields[0] + " " + fields[1] + " - - 0 1"
		result := fields[5]

		outcome := Draw
		if result == "\"1-0\";" {
			outcome = WhiteWin
		} else if result == "\"0-1\";" {
			outcome = BlackWin
		}

		var pos engine.Position
		pos.LoadFEN(fen)
		positions = append(positions, Position{Pos: pos, Outcome: outcome})
	}

	fmt.Printf("Done loading %d positions...\n", int(NumPositions))
	return positions
}

func mapWeightsToParameters() {
	copy(engine.PSQT_MG[engine.Pawn][:], Weights[0:64])
	copy(engine.PSQT_MG[engine.Knight][:], Weights[64:128])
	copy(engine.PSQT_MG[engine.Bishop][:], Weights[128:192])
	copy(engine.PSQT_MG[engine.Rook][:], Weights[192:256])
	copy(engine.PSQT_MG[engine.Queen][:], Weights[256:320])
	copy(engine.PSQT_MG[engine.King][:], Weights[320:384])

	copy(engine.PSQT_EG[engine.Pawn][:], Weights[384:448])
	copy(engine.PSQT_EG[engine.Knight][:], Weights[448:512])
	copy(engine.PSQT_EG[engine.Bishop][:], Weights[512:576])
	copy(engine.PSQT_EG[engine.Rook][:], Weights[576:640])
	copy(engine.PSQT_EG[engine.Queen][:], Weights[640:704])
	copy(engine.PSQT_EG[engine.King][:], Weights[704:768])

	copy(engine.PieceValueMG[:], Weights[768:773])
	copy(engine.PieceValueEG[:], Weights[773:778])
	copy(engine.PieceMobilityMG[:], Weights[778:782])
	copy(engine.PieceMobilityEG[:], Weights[782:786])
}

// Evaluate the position from the training set file.
func evaluate(pos engine.Position) int16 {
	score := engine.EvaluatePos(&pos)

	// For texel tuning, we always score a position from white's perspective
	if pos.SideToMove == engine.Black {
		return -score
	}

	return score
}

func processor(start, end int, K float64) {
	var errorSum float64
	for i := start; i < end; i++ {
		score := float64(evaluate(Positions[i].Pos))
		sigmoid := 1 / (1 + math.Pow(10, -K*score/400))
		errorSum += math.Pow(Positions[i].Outcome-sigmoid, 2)
	}
	Answers <- errorSum
}

// Calculate the mean square error given the current weights. Credit to
// Amanj Sherwany (author of Zahak) for this parallelized implementation.
func meanSquaredError(K float64) float64 {
	mapWeightsToParameters()
	var errorSum float64

	batchSize := len(Positions) / NumCores
	for i := 0; i < NumCores; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == NumCores-1 {
			end = len(Positions)
		}
		go processor(start, end, K)
	}

	for i := 0; i < NumCores; i++ {
		ans := <-Answers
		errorSum += ans
	}

	return errorSum / float64(len(Positions))
}

func findK() float64 {
	start, end, step := float64(0), float64(10), float64(1)
	err := float64(0)

	curr := start
	best := meanSquaredError(start)

	for i := 0; i < KPrecision; i++ {
		curr = start - step
		for curr < end {
			curr = curr + step
			err = meanSquaredError(curr)
			if err <= best {
				best = err
				start = curr
			}
		}

		fmt.Printf("Best K of %f on iteration %d\n", start, i)

		end = start + step
		start = start - step
		step = step / 10.0
	}

	return start
}

func tune() {
	numParams := len(Weights)
	bestError := meanSquaredError(K)

	improved := true
	for iteration := 1; improved; iteration++ {
		improved = false
		for weightIdx := 0; weightIdx < numParams; weightIdx++ {
			if IgnoreWeights[weightIdx] {
				continue
			}

			Weights[weightIdx] += 1
			newError := meanSquaredError(K)

			if newError < bestError {
				bestError = newError
				improved = true
			} else {
				Weights[weightIdx] -= 2

				if weightIdx >= 768 && Weights[weightIdx] <= 0 {
					Weights[weightIdx] += 1
					continue
				}

				newError = meanSquaredError(K)
				if newError < bestError {
					bestError = newError
					improved = true
				} else {
					Weights[weightIdx] += 1
				}
			}
		}

		fmt.Printf("Iteration %d complete...\n", iteration)
		fmt.Printf("Best error: %.15f\n", bestError)

		if iteration%10 == 0 {
			printParameters()
		}
	}

	fmt.Println("Done tuning!")
}

func prettyPrintPSQT(psqt [64]int16) {
	fmt.Print("\n")
	for sq := 0; sq < 64; sq++ {
		if sq%8 == 0 {
			fmt.Println()
		}
		fmt.Print(psqt[sq], ", ")
	}
	fmt.Print("\n")
}

func printParameters() {
	prettyPrintPSQT(engine.PSQT_MG[engine.Pawn])
	prettyPrintPSQT(engine.PSQT_MG[engine.Knight])
	prettyPrintPSQT(engine.PSQT_MG[engine.Bishop])
	prettyPrintPSQT(engine.PSQT_MG[engine.Rook])
	prettyPrintPSQT(engine.PSQT_MG[engine.Queen])
	prettyPrintPSQT(engine.PSQT_MG[engine.King])

	prettyPrintPSQT(engine.PSQT_EG[engine.Pawn])
	prettyPrintPSQT(engine.PSQT_EG[engine.Knight])
	prettyPrintPSQT(engine.PSQT_EG[engine.Bishop])
	prettyPrintPSQT(engine.PSQT_EG[engine.Rook])
	prettyPrintPSQT(engine.PSQT_EG[engine.Queen])
	prettyPrintPSQT(engine.PSQT_EG[engine.King])

	fmt.Println(engine.PieceValueMG)
	fmt.Println(engine.PieceValueEG)
	fmt.Println(engine.PieceMobilityMG)
	fmt.Println(engine.PieceMobilityEG)
}

func RunTuner(verbose bool) {
	// K := findK()
	// fmt.Println("Best K is:", K)

	Positions = loadPositions()
	tune()
	mapWeightsToParameters()

	if verbose {
		printParameters()
	}
}
