package tuner

import (
	"blunder/engine"
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strings"
)

const (
	WhiteWin float64 = 1.0
	BlackWin float64 = 0.0
	Draw     float64 = 0.5
	K        float64 = 0.008

	// Make PSQT tags correspond to piece types to use
	// piece type iterator variable to compute right tag value.
)

const (
	MG_PAWN_PSQT_TAG uint8 = iota
	MG_KNIGHT_PSQT_TAG
	MG_BISHOP_PSQT_TAG
	MG_ROOK_PSQT_TAG
	MG_QUEEN_PSQT_TAG
	MG_KING_PSQT_TAG

	EG_PAWN_PSQT_TAG
	EG_KNIGHT_PSQT_TAG
	EG_BISHOP_PSQT_TAG
	EG_ROOK_PSQT_TAG
	EG_QUEEN_PSQT_TAG
	EG_KING_PSQT_TAG

	MG_PIECE_VALUES_TAG
	EG_PIECE_VALUES_TAG

	// This constant needs to be last to ensure it correctly counts the number
	// of tag constants that appear before it.
	NUM_TAGS
)

type Weight struct {
	Value float64
	Tag   uint8
}

type Feature struct {
	Value float64
	Index uint16
}

type TuningPosition struct {
	Features []Feature
	Outcome  float64
}

func printPSQTWeights(weights []int16) {
	fmt.Println("    {")
	for idx, weight := range weights {
		if idx%8 == 0 {
			if idx == 0 {
				fmt.Print("        ")
			} else {
				fmt.Print("\n        ")
			}
		}
		fmt.Printf("%3d, ", weight)
	}
	fmt.Println("\n    },")
}

func printArray(weights []int16) {
	arrayAsString := strings.Builder{}
	for _, weight := range weights {
		arrayAsString.WriteString(fmt.Sprintf("%d, ", weight))
	}

	fmt.Printf("{%s}\n", strings.TrimSuffix(arrayAsString.String(), ", "))
}

func printWeights(weights []Weight) {
	weightsGroupedByType := [NUM_TAGS][]int16{}
	for _, weight := range weights {
		weightsGroupedByType[weight.Tag] = append(weightsGroupedByType[weight.Tag], int16(weight.Value))
	}

	fmt.Print("var MG_PIECE_VALUES = [6]int16")
	printArray(weightsGroupedByType[MG_PIECE_VALUES_TAG])

	fmt.Print("var EG_PIECE_VALUES = [6]int16")
	printArray(weightsGroupedByType[EG_PIECE_VALUES_TAG])

	fmt.Println("\nvar MG_PSQT = [6][64]int16{")
	for tag := uint8(0); tag <= MG_KING_PSQT_TAG; tag++ {
		printPSQTWeights(weightsGroupedByType[tag])
	}
	fmt.Println("}")

	fmt.Println("\nvar EG_PSQT = [6][64]int16{")
	for tag := EG_PAWN_PSQT_TAG; tag <= EG_KING_PSQT_TAG; tag++ {
		printPSQTWeights(weightsGroupedByType[tag])
	}
	fmt.Println("}")
}

func loadPositions(infile string, numPositions int) (positions []TuningPosition) {
	file, err := os.Open(infile)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	pos := engine.Position{}

	for positionCount := 0; scanner.Scan() && positionCount < numPositions; positionCount++ {
		line := scanner.Text()
		fields := strings.Fields(line)

		fen := fields[0] + " " + fields[1] + " - - 0 1"
		result := fields[6]

		outcome := Draw
		if result == "[1.0]" {
			outcome = WhiteWin
		} else if result == "[0.0]" {
			outcome = BlackWin
		}

		pos.LoadFEN(fen)
		positions = append(
			positions,
			TuningPosition{Features: computeFeatureVector(&pos), Outcome: outcome},
		)
	}

	fmt.Printf("Done loading %d positions...\n", numPositions)
	return positions
}

func computeWeightVector() (weights []Weight) {
	for pieceType := uint8(0); pieceType < engine.King; pieceType++ {
		weights = append(weights, Weight{float64(engine.MG_PIECE_VALUES[pieceType]), MG_PIECE_VALUES_TAG})
	}

	for pieceType := uint8(0); pieceType < engine.King; pieceType++ {
		weights = append(weights, Weight{float64(engine.EG_PIECE_VALUES[pieceType]), EG_PIECE_VALUES_TAG})
	}

	for pieceType := uint8(0); pieceType < engine.NoType; pieceType++ {
		for j := 0; j < 64; j++ {
			weights = append(weights, Weight{float64(engine.MG_PSQT[pieceType][j]), pieceType})
		}
	}

	for pieceType := uint8(0); pieceType < engine.NoType; pieceType++ {
		for sq := 0; sq < 64; sq++ {
			weights = append(weights, Weight{float64(engine.EG_PSQT[pieceType][sq]), 6 + pieceType})
		}
	}

	return weights
}

func computeDefaultWeightVector() (weights []Weight) {
	weights = append(weights, Weight{100, MG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{300, MG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{300, MG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{500, MG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{900, MG_PIECE_VALUES_TAG})

	weights = append(weights, Weight{100, EG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{300, EG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{300, EG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{500, EG_PIECE_VALUES_TAG})
	weights = append(weights, Weight{900, EG_PIECE_VALUES_TAG})


	for pieceType := uint8(0); pieceType < engine.NoType; pieceType++ {
		for j := 0; j < 64; j++ {
			weights = append(weights, Weight{0, pieceType})
		}
	}

	for pieceType := uint8(0); pieceType < engine.NoType; pieceType++ {
		for sq := 0; sq < 64; sq++ {
			weights = append(weights, Weight{0, 6 + pieceType})
		}
	}

	return weights
}

func computeFeatureVector(pos *engine.Position) (features []Feature) {
	mgPhasePercentage := float64(pos.Phase) / float64(engine.TOTAL_PHASE)
	egPhasePercentage := 1.0 - mgPhasePercentage

	denseFeatureVector := []float64{}
	denseFeatureVector = append(denseFeatureVector, computeMaterialFeatures(pos, mgPhasePercentage)...)
	denseFeatureVector = append(denseFeatureVector, computeMaterialFeatures(pos, egPhasePercentage)...)

	denseFeatureVector = append(denseFeatureVector, computePSQTFeatures(pos, mgPhasePercentage)...)
	denseFeatureVector = append(denseFeatureVector, computePSQTFeatures(pos, egPhasePercentage)...)

	for idx, value := range denseFeatureVector {
		if value != 0 {
			features = append(features, Feature{Index: uint16(idx), Value: value})
		}
	}

	return features
}

func computeMaterialFeatures(pos *engine.Position, phasePercentage float64) (features []float64) {
	for pieceType := uint8(0); pieceType < engine.King; pieceType++ {
		whitePieceBB := pos.Pieces[pieceType] & pos.Sides[engine.White]
		blackPieceBB := pos.Pieces[pieceType] & pos.Sides[engine.Black]
		pieceTypeCountDifference := bits.OnesCount64(whitePieceBB) - bits.OnesCount64(blackPieceBB)
		features = append(features, phasePercentage*float64(pieceTypeCountDifference))
	}
	return features
}

func computePSQTFeatures(pos *engine.Position, phasePercentage float64) (features []float64) {
	features = append(features, computePSQTFeaturesPerPieceType(pos, phasePercentage, engine.Pawn)...)
	features = append(features, computePSQTFeaturesPerPieceType(pos, phasePercentage, engine.Knight)...)
	features = append(features, computePSQTFeaturesPerPieceType(pos, phasePercentage, engine.Bishop)...)
	features = append(features, computePSQTFeaturesPerPieceType(pos, phasePercentage, engine.Rook)...)
	features = append(features, computePSQTFeaturesPerPieceType(pos, phasePercentage, engine.Queen)...)
	features = append(features, computePSQTFeaturesPerPieceType(pos, phasePercentage, engine.King)...)
	return features
}

func computePSQTFeaturesPerPieceType(pos *engine.Position, phasePercentage float64, pieceTypeTarget uint8) (features []float64) {
	features = make([]float64, 64)
	allBB := pos.Sides[engine.White] | pos.Sides[engine.Black]

	for allBB != 0 {
		sq := engine.BitScanAndClear(&allBB)
		if pos.GetPieceType(uint8(sq)) == pieceTypeTarget {
			if pos.GetPieceColor(uint8(sq)) == engine.White {
				features[engine.FlipSq[engine.White][sq]] += phasePercentage
			} else {
				features[sq] += -phasePercentage
			}
		}
	}

	return features
}

func evaluate(weights []Weight, features []Feature) (sum float64) {
	for _, feature := range features {
		sum += weights[feature.Index].Value * feature.Value
	}
	return sum
}

func computeMSE(weights []Weight, positions []TuningPosition) (errSum float64) {
	for i := range positions {
		pos := &positions[i]
		eval := evaluate(weights, pos.Features)
		sigmoid := 1 / (1 + math.Exp(-K * eval))
		errorTerm := pos.Outcome - sigmoid
		errSum += math.Pow(errorTerm, 2)
	}
	return errSum / float64(len(positions))
}

func computePartialGradient(partialGradients chan []float64, weights []Weight, positions []TuningPosition) {
	numWeights := len(weights)
	gradients := make([]float64, numWeights)

	for i := range positions {
		pos := &positions[i]

		eval := evaluate(weights, pos.Features)
		sigmoid := 1 / (1 + math.Exp(-K * eval))

		errorTerm := pos.Outcome - sigmoid
		evalTerm := sigmoid * (1 - sigmoid)

		// Note the gradient here is incomplete, and should inclue the -2k/N coefficent. However,
		// algebraically this can be factored out of the equation and done only when we need to use
		// the gradient. This saves time and floating point accuracy. Thanks to Ethereal for this tweak.
		for j := range pos.Features {
			feature := &pos.Features[j]
			gradients[feature.Index] += errorTerm * evalTerm * feature.Value
		}
	}

	partialGradients <- gradients
}

func computeGradient(weights []Weight, positions []TuningPosition, numCores int) (gradients []float64) {
	gradients = make([]float64, len(weights))
	partialGradients := make(chan []float64, numCores)

	posPerProcess := len(positions) / numCores
	for i := 0; i < len(positions); i += posPerProcess {
		go computePartialGradient(partialGradients, weights, positions[i:i+posPerProcess])
	}

	for i := 0; i < numCores; i++ {
		partialGradient := <-partialGradients
		for i := range partialGradient {
			gradients[i] += partialGradient[i]
		}
	}

	close(partialGradients)
	return gradients
}

func Tune(infile string, epochs, numPositions, numCores int, learningRate float64, recordErrorRate bool, useDefaultWeights bool) {
	var weights []Weight

	if useDefaultWeights {
		weights = computeDefaultWeightVector()
	} else {
		weights = computeWeightVector()
	}

	positions := loadPositions(infile, numPositions, )
	beforeErr := computeMSE(weights, positions)

	N := float64(numPositions)

	errors := []float64{beforeErr}
	errorRecordingRate := epochs / 100

	for epoch := 0; epoch < epochs; epoch++ {
		gradients := computeGradient(weights, positions, numCores)
		for k, gradient := range gradients {
			weights[k].Value += -(learningRate * -2 * K / N * gradient)
		}

		fmt.Printf("Epoch number %d completed\n", epoch+1)

		if recordErrorRate && epoch > 0 && epoch%errorRecordingRate == 0 {
			errors = append(errors, computeMSE(weights, positions))
		}
	}

	if recordErrorRate {
		errors = append(errors, computeMSE(weights, positions))
		file, err := os.Create("errors.txt")
		if err != nil {
			fmt.Println("Couldn't create \"errors.txt\" to store recored error rates")
		} else {
			fmt.Println("Storing error rates in errors.txt")
		}

		for _, err := range errors {
			_, e := file.WriteString(fmt.Sprintf("%f\n", err))
			if e != nil {
				panic(e)
			}
		}
	}

	printWeights(weights)
	fmt.Println("Best error before tuning:", beforeErr)
	fmt.Println("Best error after tuning:", computeMSE(weights, positions))
}
