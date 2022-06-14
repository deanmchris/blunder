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
	Iterations    = 2000
	NumWeights    = 790
	LearningRate  = 1000000
	ScalingFactor = 0.5

	Draw     float64 = 0.5
	WhiteWin float64 = 1.0
	BlackWin float64 = 0.0
)

// An object to hold the feature coefficents of a positon, as well
// as the index of the weight the feature corresponds to. Structuring
// in the coefficent array in this manner allows for a sparse array,
// which is much more efficent and less memory intensive.
type Coefficent struct {
	Idx   uint16
	Value float64
}

// A struct object to hold data concering a position loaded from the training file.
// Each position consists of a position board object and the outcome of the game
// the position was from.
type Entry struct {
	Coefficents []Coefficent
	Outcome     float64
}

// Load the weights for tuning from the current evaluation terms.
func loadWeights() (weights []float64) {
	tempWeights := make([]int16, NumWeights)
	weights = make([]float64, NumWeights)

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

	tempWeights[778] = engine.BishopPairBonusMG
	tempWeights[779] = engine.BishopPairBonusEG

	copy(tempWeights[780:785], engine.PieceMobilityMG[:])
	copy(tempWeights[785:790], engine.PieceMobilityEG[:])

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
		result := fields[6]

		outcome := Draw
		if result == "1.0" {
			outcome = WhiteWin
		} else if result == "0.0" {
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
func getCoefficents(pos *engine.Position) (coefficents []Coefficent) {
	phase := (pos.Phase*256 + (engine.TotalPhase / 2)) / engine.TotalPhase
	mgPhase := 256 - phase
	egPhase := phase

	staticAllBB := pos.Sides[engine.White] | pos.Sides[engine.Black]
	allBB := staticAllBB
	tempCoefficents := make([]float64, NumWeights)

	for allBB != 0 {
		sq := allBB.PopBit()
		piece := pos.Squares[sq]

		mgIndex := uint16(piece.Type)*64 + uint16(engine.FlipSq[piece.Color][sq])
		egIndex := 384 + mgIndex
		sign := float64(1)

		if piece.Color != engine.White {
			sign = -1
		}

		tempCoefficents[mgIndex] += sign * float64(mgPhase)
		tempCoefficents[egIndex] += sign * float64(egPhase)

		usBB := pos.Sides[piece.Color]
		pieceType := uint16(piece.Type)

		switch piece.Type {
		case engine.Knight:
			moves := engine.KnightMoves[sq] & ^usBB
			mobility := float64(moves.CountBits())
			tempCoefficents[780+pieceType] += (mobility - 4) * sign * float64(mgPhase)
			tempCoefficents[785+pieceType] += (mobility - 4) * sign * float64(egPhase)
		case engine.Bishop:
			moves := engine.GenBishopMoves(sq, staticAllBB) & ^usBB
			mobility := float64(moves.CountBits())
			tempCoefficents[780+pieceType] += (mobility - 7) * sign * float64(mgPhase)
			tempCoefficents[785+pieceType] += (mobility - 7) * sign * float64(egPhase)
		case engine.Rook:
			moves := engine.GenRookMoves(sq, staticAllBB) & ^usBB
			mobility := float64(moves.CountBits())
			tempCoefficents[780+pieceType] += (mobility - 7) * sign * float64(mgPhase)
			tempCoefficents[785+pieceType] += (mobility - 7) * sign * float64(egPhase)
		case engine.Queen:
			moves := (engine.GenBishopMoves(sq, staticAllBB) | engine.GenRookMoves(sq, staticAllBB)) & ^usBB
			mobility := float64(moves.CountBits())
			tempCoefficents[780+pieceType] += (mobility - 14) * sign * float64(mgPhase)
			tempCoefficents[785+pieceType] += (mobility - 14) * sign * float64(egPhase)
		}
	}

	for piece := 0; piece <= 4; piece++ {
		tempCoefficents[768+piece] = float64(
			(pos.Pieces[engine.White][piece].CountBits() - pos.Pieces[engine.Black][piece].CountBits()),
		) * float64(mgPhase)
		tempCoefficents[768+piece+5] = float64(
			(pos.Pieces[engine.White][piece].CountBits() - pos.Pieces[engine.Black][piece].CountBits()),
		) * float64(egPhase)
	}

	if pos.Pieces[engine.White][engine.Bishop].CountBits() >= 2 {
		tempCoefficents[778] += 1
		tempCoefficents[779] += 1
	}

	if pos.Pieces[engine.Black][engine.Bishop].CountBits() >= 2 {
		tempCoefficents[778] -= 1
		tempCoefficents[779] -= 1
	}

	tempCoefficents[778] *= float64(mgPhase)
	tempCoefficents[779] *= float64(egPhase)

	for i, coefficent := range tempCoefficents {
		if coefficent != 0 {
			coefficents = append(coefficents, Coefficent{Idx: uint16(i), Value: coefficent})
		}
	}

	return coefficents
}

// Evaluate the position from the training set file.
func evaluate(weights []float64, coefficents []Coefficent) (score float64) {
	for i := range coefficents {
		coefficent := &coefficents[i]
		score += weights[coefficent.Idx] * coefficent.Value / 256
	}
	return score
}

func computeGradientNumerically(entries []Entry, weights []float64, epsilon float64) (gradients []float64) {
	N := float64(len(entries))
	gradients = make([]float64, len(weights))
	epsilonAddedErrSums := make([]float64, len(entries))
	epsilonSubtractedErrSums := make([]float64, len(entries))

	for i := range entries {
		for k := range weights {
			weights[k] += epsilon

			score := evaluate(weights, entries[i].Coefficents)
			sigmoid := 1 / (1 + math.Exp(-(ScalingFactor * score)))
			err := entries[i].Outcome - sigmoid
			epsilonAddedErrSums[k] += math.Pow(err, 2)

			weights[k] -= epsilon * 2

			score = evaluate(weights, entries[i].Coefficents)
			sigmoid = 1 / (1 + math.Exp(-(ScalingFactor * score)))
			err = entries[i].Outcome - sigmoid
			epsilonSubtractedErrSums[k] += math.Pow(err, 2)

			weights[k] += epsilon
		}
	}

	for i := range gradients {
		errEpsilonAdded := epsilonAddedErrSums[i] / N
		errEpsilonSubtracted := epsilonSubtractedErrSums[i] / N
		gradients[i] = (1 / (2 * epsilon)) * (errEpsilonAdded - errEpsilonSubtracted)
	}

	return gradients
}

func computeGradient(entries []Entry, weights []float64) (gradients []float64) {
	N := float64(len(entries))
	gradients = make([]float64, NumWeights)

	for i := range entries {
		score := evaluate(weights, entries[i].Coefficents)
		sigmoid := 1 / (1 + math.Exp(-(ScalingFactor * score)))
		err := entries[i].Outcome - sigmoid
		term := -2 * ScalingFactor / N * err * (1 - sigmoid) * sigmoid

		for k := range entries[i].Coefficents {
			coefficent := &entries[i].Coefficents[k]
			gradients[coefficent.Idx] += term * coefficent.Value / 256
		}
	}

	return gradients
}

func computeMSE(entries []Entry, weights []float64) (errSum float64) {
	for i := range entries {
		score := evaluate(weights, entries[i].Coefficents)
		sigmoid := 1 / (1 + math.Exp(-(ScalingFactor * score)))
		err := entries[i].Outcome - sigmoid
		errSum += math.Pow(err, 2)
	}
	return errSum / float64(len(entries))
}

func convertFloatSiceToInt(slice []float64) (ints []int16) {
	for _, float := range slice {
		ints = append(ints, int16(float))
	}
	return ints
}

func printSlice(name string, slice []int16) {
	fmt.Print(name + ": {")
	for _, integer := range slice {
		fmt.Printf("%d, ", integer)
	}
	fmt.Print("}\n")
}

func prettyPrintPSQT(name string, psqt []int16) {
	fmt.Print("{\n")
	fmt.Print("    // ", name, "\n    ")
	for sq := 0; sq < 64; sq++ {
		if sq > 0 && sq%8 == 0 {
			fmt.Print("\n    ")
		}
		fmt.Print(psqt[sq], ", ")
	}
	fmt.Print("\n}\n")
}

func printParameters(weights []float64) {
	prettyPrintPSQT("MG Pawn PST", convertFloatSiceToInt(weights[0:64]))
	prettyPrintPSQT("MG Knight PST", convertFloatSiceToInt(weights[64:128]))
	prettyPrintPSQT("MG Bishop PST", convertFloatSiceToInt(weights[128:192]))
	prettyPrintPSQT("MG Rook PST", convertFloatSiceToInt(weights[192:256]))
	prettyPrintPSQT("MG Queen PST", convertFloatSiceToInt(weights[256:320]))
	prettyPrintPSQT("MG King PST", convertFloatSiceToInt(weights[320:384]))

	prettyPrintPSQT("EG Pawn PST", convertFloatSiceToInt(weights[384:448]))
	prettyPrintPSQT("EG Knight PST", convertFloatSiceToInt(weights[448:512]))
	prettyPrintPSQT("EG Bishop PST", convertFloatSiceToInt(weights[512:576]))
	prettyPrintPSQT("EG Rook PST", convertFloatSiceToInt(weights[576:640]))
	prettyPrintPSQT("EG Queen PST", convertFloatSiceToInt(weights[640:704]))
	prettyPrintPSQT("EG King PST", convertFloatSiceToInt(weights[704:768]))

	printSlice("\nMG Piece Values", convertFloatSiceToInt(weights[768:773]))
	printSlice("EG Piece Values", convertFloatSiceToInt(weights[773:778]))

	fmt.Println("\nBishop Pair Bonus MG:", weights[778])
	fmt.Println("\nBishop Pair Bonus EG:", weights[779])

	printSlice("\nMG Piece Mobility Coefficents", convertFloatSiceToInt(weights[780:785]))
	printSlice("EG Piece Mobility Coefficents", convertFloatSiceToInt(weights[785:790]))

	fmt.Println()
}

func Tune(infile string, epochs, numPositions int, learningRate float64) {
	weights := loadWeights()
	entries := loadEntries(infile, numPositions)
	beforeErr := computeMSE(entries, weights)

	for i := 0; i < epochs; i++ {
		gradients := computeGradient(entries, weights)
		for k, gradient := range gradients {
			weights[k] -= learningRate * gradient
		}

		fmt.Printf("Epoch number %d completed\n", i+1)
	}

	printParameters(weights)
	fmt.Println("Best error before tuning:", beforeErr)
	fmt.Println("Best error after tuning:", computeMSE(entries, weights))
}
