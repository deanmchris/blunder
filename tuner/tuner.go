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
	NumCores   = 8
	NumWeights = 815
	KPrecision = 10
	Report     = 10

	Draw     float64 = 0.5
	WhiteWin float64 = 1.0
	BlackWin float64 = 0.0
)

// A struct object to hold data concering a position loaded from the training file.
// Each position consists of a position board object and the outcome of the game
// the position was from.
type Entry struct {
	Pos     engine.Position
	Outcome float64
}

// A global variable to hold the positions loaded from the training file.
var Entries []Entry

// A global variable to hold the parallel computations of the MSE function.
var Answers = make(chan float64)

// The weights to be adjusted during the tuning process.
var Weights []int16

// The steps to adjust each weight.

var Steps []int16

// Boolean map of futile weights that shouldn't be tuned by the tuner (i.e 1st and 8th rank
// weights for pawns)
var FutileIndexes []bool

// Intialize futile indexes.
func initFutileIndexes() (indexes []bool) {
	futileIndexes := []int{
		0, 1, 2, 3, 4, 5, 6, 7,
		56, 57, 58, 59, 60, 61, 62, 63,
		384, 385, 386, 387, 388, 389, 390, 391,
		440, 441, 442, 443, 444, 445, 446, 447,
		786, 793, 794, 801,
	}

	indexes = make([]bool, NumWeights)
	for _, idx := range futileIndexes {
		indexes[idx] = true
	}

	return indexes
}

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

	copy(weights[786:794], engine.PassedPawnBonusMG[:])
	copy(weights[794:802], engine.PassedPawnBonusEG[:])

	weights[802] = engine.IsolatedPawnPenatlyMG
	weights[803] = engine.IsolatedPawnPenatlyEG
	weights[804] = engine.DoubledPawnPenatlyMG
	weights[805] = engine.DoubledPawnPenatlyEG

	weights[806] = engine.KnightOutpostBonusMG
	weights[807] = engine.KnightOutpostBonusEG

	weights[808] = engine.MinorAttackOuterRing
	weights[809] = engine.MinorAttackInnerRing
	weights[810] = engine.RookAttackOuterRing
	weights[811] = engine.RookAttackInnerRing
	weights[812] = engine.QueenAttackOuterRing
	weights[813] = engine.QueenAttackInnerRing
	weights[814] = engine.SemiOpenFileNextToKingPenalty

	return weights
}

// Load the steps for each weight.
func loadSteps() (steps []int16) {
	steps = make([]int16, NumWeights)
	copy(steps[0:NumWeights], make_int16_slice(NumWeights, 1))
	return steps
}

// Create an int16 slice of a certian size and filled
// with a specified default value.
func make_int16_slice(size int, defaultValue int16) (slice []int16) {
	slice = make([]int16, size)
	for i := range slice {
		slice[i] = defaultValue
	}
	return slice
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

		var pos engine.Position
		pos.LoadFEN(fen)
		entries = append(entries, Entry{Pos: pos, Outcome: outcome})
	}

	fmt.Printf("Done loading %d positions...\n", numPositions)
	return entries
}

func mapWeights(weights []int16) {
	copy(engine.PSQT_MG[engine.Pawn][:], weights[0:64])
	copy(engine.PSQT_MG[engine.Knight][:], weights[64:128])
	copy(engine.PSQT_MG[engine.Bishop][:], weights[128:192])
	copy(engine.PSQT_MG[engine.Rook][:], weights[192:256])
	copy(engine.PSQT_MG[engine.Queen][:], weights[256:320])
	copy(engine.PSQT_MG[engine.King][:], weights[320:384])

	copy(engine.PSQT_EG[engine.Pawn][:], weights[384:448])
	copy(engine.PSQT_EG[engine.Knight][:], weights[448:512])
	copy(engine.PSQT_EG[engine.Bishop][:], weights[512:576])
	copy(engine.PSQT_EG[engine.Rook][:], weights[576:640])
	copy(engine.PSQT_EG[engine.Queen][:], weights[640:704])
	copy(engine.PSQT_EG[engine.King][:], weights[704:768])

	copy(engine.PieceValueMG[:], weights[768:773])
	copy(engine.PieceValueEG[:], weights[773:778])
	copy(engine.PieceMobilityMG[:], weights[778:782])
	copy(engine.PieceMobilityEG[:], weights[782:786])

	copy(engine.PassedPawnBonusMG[:], weights[786:794])
	copy(engine.PassedPawnBonusEG[:], weights[794:802])

	engine.IsolatedPawnPenatlyMG = weights[802]
	engine.IsolatedPawnPenatlyEG = weights[803]
	engine.DoubledPawnPenatlyMG = weights[804]
	engine.DoubledPawnPenatlyEG = weights[805]

	engine.KnightOutpostBonusMG = weights[806]
	engine.KnightOutpostBonusEG = weights[807]

	engine.MinorAttackOuterRing = weights[808]
	engine.MinorAttackInnerRing = weights[809]
	engine.RookAttackOuterRing = weights[810]
	engine.RookAttackInnerRing = weights[811]
	engine.QueenAttackOuterRing = weights[812]
	engine.QueenAttackInnerRing = weights[813]
	engine.SemiOpenFileNextToKingPenalty = weights[814]
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
		score := float64(evaluate(Entries[i].Pos))
		sigmoid := 1 / (1 + math.Pow(10, -K*score/400))
		errorSum += math.Pow(Entries[i].Outcome-sigmoid, 2)
	}
	Answers <- errorSum
}

// Calculate the mean square error given the current weights. Credit to
// Amanj Sherwany (author of Zahak) for this parallelized implementation.
func meanSquaredError(weights []int16, K float64) (errorSum float64) {
	mapWeights(weights)
	batchSize := len(Entries) / NumCores

	for i := 0; i < NumCores; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == NumCores-1 {
			end = len(Entries)
		}
		go processor(start, end, K)
	}

	for i := 0; i < NumCores; i++ {
		ans := <-Answers
		errorSum += ans
	}

	return errorSum / float64(len(Entries))
}

// Credit to Andrew Grant (author of the chess engine Ethereal), for
// this specfic implementation of computing an appropriate K value.
func findK() float64 {
	start, end, step := float64(0), float64(10), float64(1)
	err := float64(0)

	curr := start
	best := meanSquaredError(Weights, start)

	for i := 0; i < KPrecision; i++ {
		curr = start - step
		for curr < end {
			curr = curr + step
			err = meanSquaredError(Weights, curr)
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
	fmt.Println("MG Piece Mobility Bonuses:", Weights[778:782])
	fmt.Println("EG Piece Mobility Bonuses:", Weights[782:786])

	fmt.Println("\nMG Passed Pawn Bonus:", engine.PassedPawnBonusMG)
	fmt.Println("EG Passed Pawn Bonus:", engine.PassedPawnBonusEG)

	fmt.Println("\nMG Isolated Pawn Penalty:", engine.IsolatedPawnPenatlyMG)
	fmt.Println("EG Isolated Pawn Penalty:", engine.IsolatedPawnPenatlyEG)
	fmt.Println("MG Doubled Pawn Penalty:", engine.DoubledPawnPenatlyMG)
	fmt.Println("EG Doubled Pawn Penalty:", engine.DoubledPawnPenatlyEG)

	fmt.Println("\nMG Knight Outpost Bonus:", engine.KnightOutpostBonusMG)
	fmt.Println("EG Knight Outpost Bonus:", engine.KnightOutpostBonusEG)

	fmt.Println("\nMinor Attacking Outer Ring:", engine.MinorAttackOuterRing)
	fmt.Println("Minor Attacking Inner Ring:", engine.MinorAttackInnerRing)
	fmt.Println("Rook Attacking Outer Ring:", engine.RookAttackOuterRing)
	fmt.Println("Rook Attacking Inner Ring:", engine.RookAttackInnerRing)
	fmt.Println("Queen Attacking Outer Ring:", engine.QueenAttackOuterRing)
	fmt.Println("Queen Attacking Inner Ring:", engine.QueenAttackInnerRing)
	fmt.Println("Semi-Open File Next To King Penalty:", engine.SemiOpenFileNextToKingPenalty)
}

func Tune(infile string, numPositions int) {
	Entries = loadEntries(infile, numPositions)
	Weights = loadWeights()
	Steps = loadSteps()
	FutileIndexes = initFutileIndexes()

	K := findK()
	bestError := meanSquaredError(Weights, K)
	improved := true

	for iteration := 1; improved; iteration++ {
		improved = false
		for weightIdx := 0; weightIdx < NumWeights; weightIdx++ {
			if FutileIndexes[weightIdx] {
				continue
			}

			Weights[weightIdx] += Steps[weightIdx]
			newError := meanSquaredError(Weights, K)

			if newError < bestError {
				bestError = newError
				improved = true
			} else {
				Weights[weightIdx] -= Steps[weightIdx] * 2

				// All weights but those in the piece-square tables should be
				// positive.
				if weightIdx >= 768 && Weights[weightIdx] <= 0 {
					Weights[weightIdx] += Steps[weightIdx]
					continue
				}

				newError = meanSquaredError(Weights, K)
				if newError < bestError {
					bestError = newError
					improved = true
				} else {
					Weights[weightIdx] += Steps[weightIdx]
				}
			}
		}

		if iteration%Report == 0 {
			fmt.Printf("Best evaluation terms (iteration %d):", iteration)
			printParameters()
		} else {
			fmt.Printf("iteration %d, error = %.15f\n", iteration, bestError)
		}
	}

	fmt.Println("\nDone tuning! Best weights:")
	mapWeights(Weights)
	printParameters()
}
