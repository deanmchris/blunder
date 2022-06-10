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
	KPrecision         = 10
	Draw       float64 = 0.5
	WhiteWin   float64 = 1.0
	BlackWin   float64 = 0.0
)

// A struct object to hold data concering a position loaded from the training file.
// Each position consists of a position board object and the outcome of the game
// the position was from.
type Entry struct {
	Coefficents []float64
	Outcome     float64
}

// Load the weights for tuning from the current evaluation terms.
func loadWeights(numWeights int) (weights []float64) {
	tempWeights := make([]int16, numWeights)
	weights = make([]float64, numWeights)

	copy(tempWeights[0:64], engine.PSQT_MG[engine.Pawn][:])
	copy(tempWeights[64:128], engine.PSQT_MG[engine.Knight][:])
	copy(tempWeights[128:192], engine.PSQT_MG[engine.Bishop][:])
	copy(tempWeights[192:256], engine.PSQT_MG[engine.Rook][:])
	copy(tempWeights[256:320], engine.PSQT_MG[engine.Queen][:])
	copy(tempWeights[320:384], engine.PSQT_MG[engine.King][:])

	copy(tempWeights[384:448], engine.PSQT_EG[engine.Pawn][:])
	copy(tempWeights[448:512], engine.PSQT_EG[engine.Knight][:])
	copy(tempWeights[512:576], engine.PSQT_EG[engine.Bishop][:])
	copy(tempWeights[576:640], engine.PSQT_EG[engine.Rook][:])
	copy(tempWeights[640:704], engine.PSQT_EG[engine.Queen][:])
	copy(tempWeights[704:768], engine.PSQT_EG[engine.King][:])

	copy(tempWeights[768:773], engine.PieceValueMG[:])
	copy(tempWeights[773:778], engine.PieceValueEG[:])

	for i := range tempWeights {
		weights[i] = float64(tempWeights[i])
	}

	return weights
}

// Load the given number of positions from the training set file.
func loadEntries(infile string, numPositions int, numWeights int) (entries []Entry) {
	file, err := os.Open(infile)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for positionCount := 0; scanner.Scan() && positionCount < numPositions; positionCount++ {
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

		pos := engine.Position{}
		pos.LoadFEN(fen)
		entries = append(entries, Entry{Coefficents: getCoefficents(&pos, numWeights), Outcome: outcome})
	}

	fmt.Printf("Done loading %d positions...\n", numPositions)
	return entries
}

// Get the evaluation coefficents of the position so it can be used to calculate
// the evaluation.
func getCoefficents(pos *engine.Position, numWeights int) (coefficents []float64) {
	phase := (pos.Phase*256 + (engine.TotalPhase / 2)) / engine.TotalPhase
	mgPhase := 256 - phase
	egPhase := phase

	stm := pos.SideToMove
	allBB := pos.Sides[engine.White] | pos.Sides[engine.Black]
	coefficents = make([]float64, numWeights)

	for allBB != 0 {
		sq := allBB.PopBit()
		piece := pos.Squares[sq]

		mgIndex := uint16(piece.Type)*64 + uint16(engine.FlipSq[piece.Color][sq])
		egIndex := 384 + mgIndex
		sign := float64(1)

		if piece.Color != engine.White {
			sign = -1
		}

		coefficents[mgIndex] += sign * float64(mgPhase)
		coefficents[egIndex] += sign * float64(egPhase)
	}

	for piece := 0; piece <= 4; piece++ {
		coefficents[768+piece] = float64(
			(pos.Pieces[stm][piece].CountBits() - pos.Pieces[stm^1][piece].CountBits()),
		) * float64(mgPhase)
		coefficents[768+piece+5] = float64(
			(pos.Pieces[stm][piece].CountBits() - pos.Pieces[stm^1][piece].CountBits()),
		) * float64(egPhase)
	}

	return coefficents
}

// Evaluate the position from the training set file.
func evaluate(weights, coefficents []float64) (score float64) {
	for i := 0; i < len(weights); i++ {
		score += weights[i] * coefficents[i] / 256
	}
	return score
}

func computeGradient(entries []Entry, weights []float64, scalingFactor float64) (gradients []float64) {
	N := float64(len(entries))
	numWeights := len(weights)
	gradients = make([]float64, numWeights)

	for i := 0; i < len(entries); i++ {
		score := evaluate(weights, entries[i].Coefficents)
		sigmoid := 1 / (1 + math.Exp(-scalingFactor*score))
		err := entries[i].Outcome - sigmoid
		term := -2 * scalingFactor / N * err * (1 - sigmoid) * sigmoid

		for k := 0; k < numWeights; k++ {
			gradients[k] += term * entries[i].Coefficents[k]
		}
	}

	return gradients
}

func meanSquaredError(entries []Entry, weights []float64, scalingFactor float64) (errSum float64) {
	for i := 0; i < len(entries); i++ {
		score := evaluate(weights, entries[i].Coefficents)
		sigmoid := 1 / (1 + math.Exp(-scalingFactor*score))
		errSum += math.Pow(entries[i].Outcome-sigmoid, 2)
	}
	return errSum / float64(len(entries))
}

// Credit to Andrew Grant (author of the chess engine Ethereal), for
// this specfic implementation of computing an appropriate scaling value
// for the logistic function.
func findScalingFactor(entries []Entry, weights []float64) float64 {
	start, end, step := float64(0), float64(10), float64(1)
	err := float64(0)

	curr := start
	best := meanSquaredError(entries, weights, start)

	for i := 0; i < KPrecision; i++ {
		curr = start - step
		for curr < end {
			curr = curr + step
			err = meanSquaredError(entries, weights, curr)
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

func convertFloatSiceToInt(slice []float64) (ints []int16) {
	for _, float := range slice {
		ints = append(ints, int16(float))
	}
	return ints
}

func prettyPrintPSQT(msg string, psqt []int16) {
	fmt.Print("\n")
	fmt.Println(msg)
	for sq := 0; sq < 64; sq++ {
		if sq%8 == 0 {
			fmt.Println()
		}
		fmt.Print(psqt[sq], ", ")
	}
	fmt.Print("\n")
}

func printParameters(weights []float64) {
	prettyPrintPSQT("MG Pawn PST:", convertFloatSiceToInt(weights[0:64]))
	prettyPrintPSQT("MG Knight PST:", convertFloatSiceToInt(weights[64:128]))
	prettyPrintPSQT("MG Bishop PST:", convertFloatSiceToInt(weights[128:192]))
	prettyPrintPSQT("MG Rook PST:", convertFloatSiceToInt(weights[192:256]))
	prettyPrintPSQT("MG Queen PST:", convertFloatSiceToInt(weights[256:320]))
	prettyPrintPSQT("MG King PST:", convertFloatSiceToInt(weights[320:384]))

	prettyPrintPSQT("EG Pawn PST:", convertFloatSiceToInt(weights[384:448]))
	prettyPrintPSQT("EG Knight PST:", convertFloatSiceToInt(weights[448:512]))
	prettyPrintPSQT("EG Bishop PST:", convertFloatSiceToInt(weights[512:576]))
	prettyPrintPSQT("EG Rook PST:", convertFloatSiceToInt(weights[576:640]))
	prettyPrintPSQT("EG Queen PST:", convertFloatSiceToInt(weights[640:704]))
	prettyPrintPSQT("EG King PST:", convertFloatSiceToInt(weights[704:768]))

	fmt.Println("\nMG Piece Values:", convertFloatSiceToInt(weights[768:773]))
	fmt.Println("EG Piece Values:", convertFloatSiceToInt(weights[773:778]))
	fmt.Println()
}

func Tune(infile string, epochs int, numWeights, numPositions int, learningRate float64) {
	weights := loadWeights(numWeights)
	entries := loadEntries(infile, numPositions, numWeights)
	scalingFactor := findScalingFactor(entries, weights)

	fmt.Printf("K-value computed: %f ...\n", scalingFactor)

	for i := 0; i < epochs; i++ {
		gradients := computeGradient(entries, weights, scalingFactor)
		for k, gradient := range gradients {
			weights[k] += learningRate * gradient
		}

		fmt.Printf("Epoch number %d completed\n", i+1)
	}

	printParameters(weights)
}
