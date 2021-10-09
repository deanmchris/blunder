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
var Positions = loadPositions(0)

// A global variable to hold the parallel computations of the MSE function.
var Answers = make(chan float64)

// The weights to be adjusted during the tuning process.
var Weights []int16 = []int16{
	120, 102, 99, 105, 99, 100, 101, 97,
	248, 251, 214, 213, 222, 218, 237, 151,
	102, 113, 126, 136, 171, 144, 91, 82,
	78, 104, 99, 112, 115, 104, 108, 69,
	67, 88, 86, 102, 108, 97, 99, 66,
	71, 87, 87, 81, 91, 92, 123, 81,
	59, 89, 70, 65, 73, 115, 125, 71,
	117, 101, 103, 105, 101, 99, 100, 99,

	117, 242, 322, 296, 356, 203, 277, 196,
	231, 254, 386, 348, 353, 375, 311, 308,
	270, 368, 355, 381, 390, 420, 388, 345,
	307, 332, 338, 365, 350, 385, 336, 341,
	296, 316, 331, 327, 346, 337, 346, 311,
	288, 304, 327, 326, 342, 332, 338, 298,
	304, 263, 308, 313, 312, 341, 306, 303,
	240, 293, 269, 280, 295, 307, 297, 299,

	312, 362, 258, 299, 291, 301, 336, 325,
	301, 352, 331, 327, 337, 389, 347, 292,
	313, 368, 369, 373, 374, 401, 378, 335,
	334, 350, 364, 394, 384, 377, 352, 333,
	342, 351, 350, 359, 373, 353, 353, 345,
	341, 347, 356, 353, 351, 366, 353, 345,
	336, 353, 356, 340, 345, 359, 370, 350,
	312, 345, 326, 321, 321, 325, 318, 323,

	503, 528, 513, 539, 536, 489, 524, 508,
	520, 522, 542, 544, 574, 547, 508, 529,
	492, 513, 507, 524, 493, 541, 539, 494,
	475, 490, 518, 524, 516, 533, 492, 468,
	458, 472, 484, 496, 506, 491, 510, 476,
	457, 476, 466, 474, 502, 488, 489, 462,
	452, 477, 480, 487, 492, 503, 493, 430,
	478, 485, 499, 516, 513, 504, 463, 470,

	885, 891, 933, 938, 954, 939, 945, 931,
	881, 867, 901, 931, 940, 974, 970, 950,
	889, 911, 928, 927, 939, 976, 960, 962,
	878, 878, 895, 901, 914, 924, 940, 907,
	902, 876, 904, 900, 914, 919, 921, 920,
	894, 918, 906, 911, 902, 910, 923, 922,
	887, 904, 926, 911, 919, 927, 919, 923,
	915, 897, 909, 924, 896, 888, 867, 860,

	30, 93, 79, 63, -1, 37, 42, 36,
	61, 46, 57, 73, 63, 49, 28, -3,
	38, 49, 49, 39, 28, 63, 58, 1,
	-11, 19, 30, 27, 23, 34, 28, -5,
	-12, 44, -3, -18, -7, -5, 4, -24,
	41, 13, 10, -13, -3, -4, 0, 0,
	24, 25, 9, -51, -36, -3, 17, 24,
	-15, 41, 28, -50, 17, -23, 36, 24,

	109, 95, 99, 102, 100, 102, 99, 100,
	256, 252, 217, 211, 225, 219, 238, 270,
	183, 191, 176, 154, 136, 146, 178, 176,
	122, 114, 102, 93, 82, 91, 105, 107,
	102, 99, 86, 79, 80, 78, 92, 89,
	90, 96, 82, 88, 88, 85, 84, 79,
	101, 98, 100, 100, 102, 88, 91, 81,
	114, 107, 103, 100, 99, 100, 99, 101,

	211, 244, 269, 260, 263, 262, 226, 183,
	267, 287, 257, 284, 266, 250, 264, 224,
	264, 268, 291, 298, 286, 278, 265, 252,
	273, 291, 307, 307, 306, 292, 292, 269,
	271, 284, 301, 315, 298, 306, 284, 262,
	270, 287, 287, 303, 290, 289, 266, 270,
	233, 268, 268, 278, 284, 253, 251, 244,
	235, 234, 264, 273, 257, 243, 241, 222,

	276, 270, 290, 300, 298, 289, 290, 277,
	293, 299, 309, 287, 305, 284, 297, 279,
	307, 299, 304, 300, 293, 302, 306, 304,
	295, 311, 310, 307, 318, 306, 306, 305,
	288, 306, 316, 325, 304, 314, 291, 296,
	290, 299, 311, 314, 319, 306, 294, 285,
	287, 289, 293, 305, 309, 294, 291, 273,
	277, 292, 284, 295, 291, 292, 284, 284,

	517, 509, 520, 514, 514, 514, 505, 508,
	507, 513, 516, 514, 488, 502, 509, 510,
	508, 510, 513, 503, 510, 494, 495, 503,
	508, 502, 507, 502, 504, 497, 494, 505,
	508, 504, 509, 502, 492, 498, 486, 486,
	495, 502, 500, 497, 486, 494, 494, 487,
	492, 497, 502, 501, 492, 490, 487, 497,
	489, 501, 501, 496, 489, 482, 498, 474,

	925, 955, 954, 947, 959, 932, 943, 951,
	908, 949, 959, 962, 956, 945, 923, 933,
	909, 922, 926, 980, 981, 966, 953, 939,
	934, 957, 951, 971, 973, 954, 946, 955,
	901, 959, 941, 962, 954, 951, 953, 941,
	914, 881, 928, 923, 926, 941, 937, 926,
	895, 899, 886, 913, 908, 883, 891, 886,
	904, 892, 896, 874, 919, 889, 921, 887,

	-30, -24, -15, 7, 9, 19, 14, -20,
	18, 28, 26, 23, 32, 41, 36, 26,
	15, 26, 33, 29, 24, 53, 49, 24,
	6, 34, 35, 38, 38, 43, 44, 16,
	-4, 2, 33, 40, 43, 38, 26, 4,
	-17, 10, 24, 35, 34, 31, 22, -1,
	-20, -2, 20, 32, 32, 20, 12, -3,
	-40, -22, -12, 8, -16, 3, -9, -32,
}

// Load the given number of positions from the training set file.
func loadPositions(start int) (positions []Position) {
	file, err := os.Open(DataFile)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for positionCount := 0; scanner.Scan() && positionCount < start+int(NumPositions); positionCount++ {
		if positionCount < start {
			continue
		}

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
