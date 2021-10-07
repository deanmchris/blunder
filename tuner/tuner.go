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
	DataFile = "/home/algerbrex/quiet-labeled.epd"
	NumCores = 4

	Draw         float64 = 0.5
	WhiteWin     float64 = 1.0
	BlackWin     float64 = 0.0
	NumPositions float64 = 362500.0
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
var Positions = loadPositions()

// A global variable to hold the parallel computations of the MSE function.
var Answers = make(chan float64)

// The weights to be adjusted during the tuning process.
var Weights []int16 = []int16{
	// PSQT MG pawns
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,

	// PSQT MG knights
	320, 320, 320, 320, 320, 320, 320, 320,
	320, 320, 320, 320, 320, 320, 320, 320,
	320, 320, 320, 320, 320, 320, 320, 320,
	320, 320, 320, 320, 320, 320, 320, 320,
	320, 320, 320, 320, 320, 320, 320, 320,
	320, 320, 320, 320, 320, 320, 320, 320,
	320, 320, 320, 320, 320, 320, 320, 320,
	320, 320, 320, 320, 320, 320, 320, 320,

	// PSQT MG bishops
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,

	// PSQT MG rooks
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,

	// PSQT MG queens
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,

	// PSQT MG kings
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,

	// PSQT EG pawns
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
	100, 100, 100, 100, 100, 100, 100, 100,

	// PSQT EG knights
	300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300,

	// PSQT EG bishops
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,
	330, 330, 330, 330, 330, 330, 330, 330,

	// PSQT EG rooks
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 500, 500, 500, 500,

	// PSQT EG queens
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,
	900, 900, 900, 900, 900, 900, 900, 900,

	// PSQT EG kings
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
	10, 10, 10, 10, 10, 10, 10, 10,
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
	improved := true
	bestK := 0.5
	bestError := meanSquaredError(bestK)

	for iteration := 1; improved; iteration++ {
		improved = false
		fmt.Println("Iteration:", iteration)
		fmt.Println("Best error:", bestError)
		fmt.Println("Best K:", bestK)
		fmt.Println()

		bestK += 0.01
		newError := meanSquaredError(bestK)

		if newError < bestError {
			bestError = newError
			improved = true
		} else {
			bestK -= 0.02
			newError = meanSquaredError(bestK)
			if newError < bestError {
				bestError = newError
				improved = true
			}
		}
	}

	return bestK
}

func tune() {
	numParams := len(Weights)
	bestError := meanSquaredError(K)

	improved := true
	for iteration := 1; improved; iteration++ {
		improved = false
		for weightIdx := 0; weightIdx < numParams; weightIdx++ {
			// fmt.Println("Best error:", bestError)
			// fmt.Printf("Tuning parameter number %d...\n", weightIdx)

			Weights[weightIdx] += 1
			newError := meanSquaredError(K)

			if newError < bestError {
				//fmt.Printf(
				//	"Improved parameter number %d from %d to %d\n",
				//	weight_idx, Weights[weight_idx]-1, Weights[weight_idx],
				//)
				bestError = newError
				improved = true
			} else {
				Weights[weightIdx] -= 2
				newError = meanSquaredError(K)
				if newError < bestError {
					//fmt.Printf(
					//	"Improved parameter number %d from %d to %d\n",
					//	weight_idx, Weights[weight_idx]+1, Weights[weight_idx],
					//)
					bestError = newError
					improved = true
				} else {
					Weights[weightIdx] += 1
				}
			}
		}

		fmt.Printf("Iteration %d complete...\n", iteration)
		fmt.Printf("Best error: %f\n", bestError)

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
}

func RunTuner(verbose bool) {
	// K := findK()
	// fmt.Println("Best K is:", K)

	tune()
	mapWeightsToParameters()

	if verbose {
		printParameters()
	}
}
