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
	NumPositions = 600000
	BatchSize    = 10000
	NumWeights   = 778
	KPrecision   = 10
	LearningRate = 1.0

	ScalingFactor float64 = 0.01
	Draw          float64 = 0.5
	WhiteWin      float64 = 1.0
	BlackWin      float64 = 0.0
)

// A struct object to hold data concering a position loaded from the training file.
// Each position consists of a position board object and the outcome of the game
// the position was from.
type Entry struct {
	Coefficents [NumWeights]float64
	Outcome     float64
}

// A global variable to hold the positions loaded from the training file.
var Entries []Entry

// The weights to be adjusted during the tuning process.
var Weights [NumWeights]float64

// Load the weights for tuning from the current evaluation terms.
func loadWeights() (weights [NumWeights]float64) {
	tempWeights := make([]int16, NumWeights)
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
func loadEntries(infile string, numPositions int) (entries []Entry) {
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
		entries = append(entries, Entry{Coefficents: getCoefficents(&pos), Outcome: outcome})
	}

	fmt.Printf("Done loading %d positions...\n", numPositions)
	return entries
}

// Get the evaluation coefficents of the position so it can be used to calculate
// the evaluation.
func getCoefficents(pos *engine.Position) (coefficents [NumWeights]float64) {
	phase := (pos.Phase*256 + (engine.TotalPhase / 2)) / engine.TotalPhase
	mgPhase := 256 - phase
	egPhase := phase

	stm := pos.SideToMove
	allBB := pos.Sides[engine.White] | pos.Sides[engine.Black]

	for allBB != 0 {
		sq := allBB.PopBit()
		piece := pos.Squares[sq]

		mgIndex := uint16(piece.Type)*64 + uint16(engine.FlipSq[piece.Color][sq])
		egIndex := 384 + mgIndex
		sign := float64(1)

		if piece.Color != stm {
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
func evaluate(weights, coefficents [NumWeights]float64) (score float64) {
	for i := 0; i < NumWeights; i++ {
		score += weights[i] * coefficents[i] / 256
	}
	return score
}

func computePartialDerivate(weights [NumWeights]float64, weightIdx int, batch []Entry) (sum float64) {
	for i := 0; i < BatchSize; i++ {
		score := evaluate(weights, batch[i].Coefficents)
		eTerm := math.Exp(-ScalingFactor * score)
		eTerm = eTerm / (math.Pow(1+eTerm, 2))
		eTerm *= batch[i].Coefficents[weightIdx]
		sum += -eTerm
	}
	return sum / NumWeights * 2
}

func prettyPrintPSQT(msg string, psqt []float64) {
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

func printParameters() {
	prettyPrintPSQT("MG Pawn PST:", Weights[0:64])
	prettyPrintPSQT("MG Knight PST:", Weights[64:128])
	prettyPrintPSQT("MG Bishop PST:", Weights[128:192])
	prettyPrintPSQT("MG Rook PST:", Weights[192:256])
	prettyPrintPSQT("MG Queen PST:", Weights[256:320])
	prettyPrintPSQT("MG King PST:", Weights[320:384])

	prettyPrintPSQT("EG Pawn PST:", Weights[384:448])
	prettyPrintPSQT("EG Knight PST:", Weights[448:512])
	prettyPrintPSQT("EG Bishop PST:", Weights[512:576])
	prettyPrintPSQT("EG Rook PST:", Weights[576:640])
	prettyPrintPSQT("EG Queen PST:", Weights[640:704])
	prettyPrintPSQT("EG King PST:", Weights[704:768])

	fmt.Println("\nMG Piece Values:", Weights[768:773])
	fmt.Println("EG Piece Values:", Weights[773:778])
	fmt.Println()
}

func Tune(infile string, epochs int) {
	Weights = loadWeights()
	Entries = loadEntries(infile, NumPositions)

	batches := [NumPositions / BatchSize][]Entry{}
	index := 0

	for batchStart := 0; batchStart < NumPositions; batchStart += BatchSize {
		batches[index] = Entries[batchStart : batchStart+BatchSize]
		index++
	}

	fmt.Printf("%d batches partitioned...\n", NumPositions/BatchSize)

	for i := 0; i < epochs; i++ {
		for k, batch := range batches {
			copyWeights := Weights
			for j := 0; j < NumWeights; j++ {
				Weights[j] += LearningRate * computePartialDerivate(copyWeights, j, batch)
			}

			fmt.Printf("Batch %d/%d completed\n", k+1, NumPositions/BatchSize)
			if k%10 == 0 && k > 0 {
				printParameters()
			}
		}

		fmt.Printf("Epoch number %d completed\n", i+1)
		if i > 0 && i%10 == 0 {
			printParameters()
		}
	}
}
