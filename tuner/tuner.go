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
	Iterations              = 2000
	NumWeights              = 931
	NumSafetyEvalTerms      = 9
	SafetyEvalTermsStartIdx = 922
	ScalingFactor           = 0.01
	Epsilon                 = 0.00000001
	LearningRate            = 0.5

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
	NormalCoefficents []Coefficent
	SafetyCoefficents [][]Coefficent
	MGPhase           float64
	Outcome           float64
}

// An object to store useful data while tracing safety coefficents,
// similar to the Eval object in evaluation.go.
type SafetyTracer struct {
	KingZones     [2]engine.KingZone
	KingAttackers [2]uint8
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

	copy(tempWeights[768:773], engine.PieceValueMG[0:5])
	copy(tempWeights[773:778], engine.PieceValueEG[0:5])

	tempWeights[778] = engine.BishopPairBonusMG
	tempWeights[779] = engine.BishopPairBonusEG

	copy(tempWeights[780:784], engine.PieceMobilityMG[1:5])
	copy(tempWeights[784:788], engine.PieceMobilityEG[1:5])

	copy(tempWeights[788:852], engine.PassedPawnPSQT_MG[:])
	copy(tempWeights[852:916], engine.PassedPawnPSQT_EG[:])

	tempWeights[916] = engine.DoubledPawnPenatlyMG
	tempWeights[917] = engine.DoubledPawnPenatlyEG
	tempWeights[918] = engine.IsolatedPawnPenatlyMG
	tempWeights[919] = engine.IsolatedPawnPenatlyEG

	tempWeights[920] = engine.RookOrQueenOnSeventhBonusMG
	tempWeights[921] = engine.RookOrQueenOnSeventhBonusEG

	copy(tempWeights[922:926], engine.OuterRingAttackPoints[1:5])
	copy(tempWeights[926:930], engine.InnerRingAttackPoints[1:5])
	tempWeights[930] = engine.SemiOpenFileNextToKingPenalty

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
		if result == "[1.0]" {
			outcome = WhiteWin
		} else if result == "[0.0]" {
			outcome = BlackWin
		}

		pos := engine.Position{}
		pos.LoadFEN(fen)

		normalCoefficents, safetyCoefficents := getCoefficents(&pos)
		phase := (pos.Phase*256 + (engine.TotalPhase / 2)) / engine.TotalPhase
		mgPhase := float64(256-phase) / 256

		entries = append(
			entries,
			Entry{
				NormalCoefficents: normalCoefficents,
				SafetyCoefficents: safetyCoefficents,
				Outcome:           outcome,
				MGPhase:           mgPhase,
			},
		)
	}

	fmt.Printf("Done loading %d positions...\n", numPositions)
	return entries
}

// Get the evaluation coefficents of the position so it can be used to calculate
// the evaluation.
func getCoefficents(pos *engine.Position) (normalCoefficents []Coefficent, safetyCoefficents [][]Coefficent) {
	phase := (pos.Phase*256 + (engine.TotalPhase / 2)) / engine.TotalPhase
	mgPhase := float64(256-phase) / 256
	egPhase := float64(phase) / 256

	allBB := pos.Sides[engine.White] | pos.Sides[engine.Black]

	rawNormCoefficents := make([]float64, NumWeights)
	rawSafetyCoefficents := make([][]float64, 2)
	rawSafetyCoefficents[engine.White] = make([]float64, NumSafetyEvalTerms)
	rawSafetyCoefficents[engine.Black] = make([]float64, NumSafetyEvalTerms)
	safetyCoefficents = make([][]Coefficent, 2)

	safetyTracer := SafetyTracer{
		KingZones: [2]engine.KingZone{
			engine.KingZones[pos.Pieces[engine.Black][engine.King].Msb()],
			engine.KingZones[pos.Pieces[engine.White][engine.King].Msb()],
		},
	}

	for allBB != 0 {
		sq := allBB.PopBit()
		piece := pos.Squares[sq]

		sign := float64(1)
		if piece.Color != engine.White {
			sign = -1
		}

		getPSQT_Coefficents(rawNormCoefficents, piece, sq, sign, mgPhase, egPhase)

		switch piece.Type {
		case engine.Pawn:
			getPawnCoefficents(pos, rawNormCoefficents, sq, mgPhase, egPhase, sign)
		case engine.Knight:
			getKnightCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, sq, mgPhase, egPhase, sign)
		case engine.Bishop:
			getBishopCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, sq, mgPhase, egPhase, sign)
		case engine.Rook:
			getRookCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, sq, mgPhase, egPhase, sign)
		case engine.Queen:
			getQueenCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, sq, mgPhase, egPhase, sign)
		}
	}

	getMaterialCoeffficents(pos, rawNormCoefficents, mgPhase, egPhase)
	getBishopPairCoefficents(pos, rawNormCoefficents, mgPhase, egPhase)

	getPawnShieldCoefficents(pos, pos.Pieces[engine.White][engine.King].Msb(), engine.White, rawSafetyCoefficents)
	getPawnShieldCoefficents(pos, pos.Pieces[engine.Black][engine.King].Msb(), engine.Black, rawSafetyCoefficents)

	for i, coefficent := range rawNormCoefficents {
		if coefficent != 0 {
			normalCoefficents = append(normalCoefficents, Coefficent{Idx: uint16(i), Value: coefficent})
		}
	}

	for i, coefficent := range rawSafetyCoefficents[engine.White] {
		value := float64(0)
		if safetyTracer.KingAttackers[engine.White] >= 2 && pos.Pieces[engine.White][engine.Queen] != 0 {
			value = coefficent
		}

		safetyCoefficents[engine.White] = append(
			safetyCoefficents[engine.White],
			Coefficent{Idx: SafetyEvalTermsStartIdx + uint16(i), Value: value},
		)
	}

	for i, coefficent := range rawSafetyCoefficents[engine.Black] {
		value := float64(0)
		if safetyTracer.KingAttackers[engine.Black] >= 2 && pos.Pieces[engine.Black][engine.Queen] != 0 {
			value = coefficent
		}

		safetyCoefficents[engine.Black] = append(
			safetyCoefficents[engine.Black],
			Coefficent{Idx: SafetyEvalTermsStartIdx + uint16(i), Value: value},
		)
	}

	return normalCoefficents, safetyCoefficents
}

// Get the piece square table coefficents of the position.
func getPSQT_Coefficents(coefficents []float64, piece engine.Piece, sq uint8, sign, mgPhase, egPhase float64) {
	mgIndex := uint16(piece.Type)*64 + uint16(engine.FlipSq[piece.Color][sq])
	egIndex := 384 + mgIndex
	coefficents[mgIndex] += sign * mgPhase
	coefficents[egIndex] += sign * egPhase
}

// Get the material coefficents of the position.
func getMaterialCoeffficents(pos *engine.Position, coefficents []float64, mgPhase, egPhase float64) {
	for piece := 0; piece <= 4; piece++ {
		coefficents[768+piece] = float64(
			(pos.Pieces[engine.White][piece].CountBits() - pos.Pieces[engine.Black][piece].CountBits()),
		) * mgPhase
		coefficents[773+piece] = float64(
			(pos.Pieces[engine.White][piece].CountBits() - pos.Pieces[engine.Black][piece].CountBits()),
		) * egPhase
	}
}

// Get the bishop pair coefficents of the position.
func getBishopPairCoefficents(pos *engine.Position, coefficents []float64, mgPhase, egPhase float64) {
	if pos.Pieces[engine.White][engine.Bishop].CountBits() >= 2 {
		coefficents[778] += mgPhase
		coefficents[779] += egPhase
	}

	if pos.Pieces[engine.Black][engine.Bishop].CountBits() >= 2 {
		coefficents[778] -= mgPhase
		coefficents[779] -= egPhase
	}
}

// Get the coefficents of the position related to the given pawn.
func getPawnCoefficents(pos *engine.Position, norm []float64, sq uint8, mgPhase, egPhase, sign float64) {
	piece := pos.Squares[sq]
	enemyPawns := pos.Pieces[piece.Color^1][engine.Pawn]
	usPawns := pos.Pieces[piece.Color][engine.Pawn]

	// Evaluate isolated pawns.
	if engine.IsolatedPawnMasks[engine.FileOf(sq)]&usPawns == 0 {
		norm[918] -= sign * mgPhase
		norm[919] -= sign * egPhase
	}

	// Evaluate doubled pawns.
	if engine.DoubledPawnMasks[piece.Color][sq]&usPawns != 0 {
		norm[916] -= sign * mgPhase
		norm[917] -= sign * egPhase
	}

	// Evaluate passed pawns, but make sure they're not behind a friendly pawn.
	if engine.PassedPawnMasks[piece.Color][sq]&enemyPawns == 0 &&
		usPawns&engine.DoubledPawnMasks[piece.Color][sq] == 0 {

		mgIndex := 788 + uint16(engine.FlipSq[piece.Color][sq])
		egIndex := 852 + uint16(engine.FlipSq[piece.Color][sq])

		norm[mgIndex] += sign * mgPhase
		norm[egIndex] += sign * egPhase
	}
}

// Get the coefficents of the position related to the given knight.
func getKnightCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]

	moves := engine.KnightMoves[sq] & ^usBB
	mobility := float64(moves.CountBits())
	norm[780+uint16(piece.Type)-1] += (mobility - 4) * sign * mgPhase
	norm[784+uint16(piece.Type)-1] += (mobility - 4) * sign * egPhase

	outerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].OuterRing
	innerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		safetyTracer.KingAttackers[piece.Color]++
		safety[piece.Color][piece.Type-1] += float64(outerRingAttacks.CountBits())
		safety[piece.Color][4+piece.Type-1] += float64(innerRingAttacks.CountBits())
	}
}

// Get the coefficents of the position related to the given bishop.
func getBishopCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]
	allBB := usBB | pos.Sides[piece.Color^1]

	moves := engine.GenBishopMoves(sq, allBB) & ^usBB
	mobility := float64(moves.CountBits())
	norm[780+uint16(piece.Type)-1] += (mobility - 7) * sign * mgPhase
	norm[784+uint16(piece.Type)-1] += (mobility - 7) * sign * egPhase

	outerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].OuterRing
	innerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		safetyTracer.KingAttackers[piece.Color]++
		safety[piece.Color][piece.Type-1] += float64(outerRingAttacks.CountBits())
		safety[piece.Color][4+piece.Type-1] += float64(innerRingAttacks.CountBits())
	}
}

// Get the coefficents of the position related to the given rook.
func getRookCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]
	allBB := usBB | pos.Sides[piece.Color^1]

	enemyKingSq := pos.Pieces[piece.Color^1][engine.King].Msb()
	if engine.FlipRank[piece.Color][engine.RankOf(sq)] == engine.Rank7 &&
		engine.FlipRank[piece.Color][engine.RankOf(enemyKingSq)] >= engine.Rank7 {

		norm[920] += sign * mgPhase
		norm[921] += sign * egPhase
	}

	moves := engine.GenRookMoves(sq, allBB) & ^usBB
	mobility := float64(moves.CountBits())
	norm[780+uint16(piece.Type)-1] += (mobility - 7) * sign * float64(mgPhase)
	norm[784+uint16(piece.Type)-1] += (mobility - 7) * sign * float64(egPhase)

	outerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].OuterRing
	innerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		safetyTracer.KingAttackers[piece.Color]++
		safety[piece.Color][piece.Type-1] += float64(outerRingAttacks.CountBits())
		safety[piece.Color][4+piece.Type-1] += float64(innerRingAttacks.CountBits())
	}
}

// Get the coefficents of the position related to the given queen.
func getQueenCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]
	allBB := usBB | pos.Sides[piece.Color^1]

	enemyKingSq := pos.Pieces[piece.Color^1][engine.King].Msb()
	if engine.FlipRank[piece.Color][engine.RankOf(sq)] == engine.Rank7 &&
		engine.FlipRank[piece.Color][engine.RankOf(enemyKingSq)] >= engine.Rank7 {

		norm[920] += sign * mgPhase
		norm[921] += sign * egPhase
	}

	moves := (engine.GenBishopMoves(sq, allBB) | engine.GenRookMoves(sq, allBB)) & ^usBB
	mobility := float64(moves.CountBits())
	norm[780+uint16(piece.Type)-1] += (mobility - 14) * sign * float64(mgPhase)
	norm[784+uint16(piece.Type)-1] += (mobility - 14) * sign * float64(egPhase)

	outerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].OuterRing
	innerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		safetyTracer.KingAttackers[piece.Color]++
		safety[piece.Color][piece.Type-1] += float64(outerRingAttacks.CountBits())
		safety[piece.Color][4+piece.Type-1] += float64(innerRingAttacks.CountBits())
	}
}

// Compute the coefficents releated to king safety via pawn shields
func getPawnShieldCoefficents(pos *engine.Position, sq, color uint8, safety [][]float64) {
	kingFile := engine.MaskFile[engine.FileOf(sq)]
	usPawns := pos.Pieces[color][engine.Pawn]

	leftFile := ((kingFile & engine.ClearFile[engine.FileA]) << 1)
	rightFile := ((kingFile & engine.ClearFile[engine.FileH]) >> 1)

	if kingFile&usPawns == 0 {
		safety[color^1][8] += 1
	}

	if leftFile != 0 && leftFile&usPawns == 0 {
		safety[color^1][8] += 1
	}

	if rightFile != 0 && rightFile&usPawns == 0 {
		safety[color^1][8] += 1
	}
}

// Compute the dot product between an array of king safety coefficents and the appropriate
// weight values.
func computeSafetyDotProduct(v1 []float64, v2 []Coefficent) (sum float64) {
	for i, coefficent := range v2 {
		sum += coefficent.Value * v1[i]
	}
	return sum
}

// Evaluate the position from the training set file.
func evaluate(weights []float64, normalCoefficents []Coefficent, safetyCoefficents [][]Coefficent, mgPhase float64) (score float64) {
	for i := range normalCoefficents {
		coefficent := &normalCoefficents[i]
		score += weights[coefficent.Idx] * coefficent.Value
	}

	whiteSafety := computeSafetyDotProduct(weights[SafetyEvalTermsStartIdx:NumWeights], safetyCoefficents[engine.White])
	blackSafety := computeSafetyDotProduct(weights[SafetyEvalTermsStartIdx:NumWeights], safetyCoefficents[engine.Black])

	whiteSafety = ((whiteSafety * whiteSafety) / 4) * mgPhase
	blackSafety = ((blackSafety * blackSafety) / 4) * mgPhase

	return score + whiteSafety - blackSafety
}

func computeGradientNumerically(entries []Entry, weights []float64, epsilon float64) (gradients []float64) {
	N := float64(len(entries))
	gradients = make([]float64, len(weights))
	epsilonAddedErrSums := make([]float64, len(entries))
	epsilonSubtractedErrSums := make([]float64, len(entries))

	for i := range entries {
		for k := range weights {
			weights[k] += epsilon

			score := evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, entries[i].MGPhase)
			sigmoid := 1 / (1 + math.Exp(-(ScalingFactor * score)))
			err := entries[i].Outcome - sigmoid
			epsilonAddedErrSums[k] += math.Pow(err, 2)

			weights[k] -= epsilon * 2

			score = evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, entries[i].MGPhase)
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
	gradients = make([]float64, NumWeights)

	for i := range entries {
		score := evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, entries[i].MGPhase)
		sigmoid := 1 / (1 + math.Exp(-(ScalingFactor * score)))
		err := entries[i].Outcome - sigmoid

		// Note the gradient here is incomplete, and should inclue the -2k/N coefficent. However,
		// algebraically this can be factored out of the equation and done only when we need to use
		// the gradient. This saves time and accuracy. Thanks to Ethereal for this tweak.
		term := err * (1 - sigmoid) * sigmoid

		for k := range entries[i].NormalCoefficents {
			coefficent := &entries[i].NormalCoefficents[k]
			gradients[coefficent.Idx] += term * coefficent.Value
		}

		whiteSafety := computeSafetyDotProduct(weights[SafetyEvalTermsStartIdx:NumWeights], entries[i].SafetyCoefficents[engine.White])
		blackSafety := computeSafetyDotProduct(weights[SafetyEvalTermsStartIdx:NumWeights], entries[i].SafetyCoefficents[engine.Black])

		for k := range entries[i].SafetyCoefficents[engine.White] {
			whiteCoefficent := &entries[i].SafetyCoefficents[engine.White][k]
			blackCoefficent := &entries[i].SafetyCoefficents[engine.Black][k]

			whiteTerm := whiteSafety * whiteCoefficent.Value * entries[i].MGPhase / 2
			blackTerm := blackSafety * blackCoefficent.Value * entries[i].MGPhase / 2

			gradients[whiteCoefficent.Idx] += term * (whiteTerm - blackTerm)
		}
	}

	return gradients
}

func computeMSE(entries []Entry, weights []float64) (errSum float64) {
	for i := range entries {
		score := evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, entries[i].MGPhase)
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
	fmt.Print("\n},\n")
}

func printParameters(weights []float64) {
	printSlice("\nMG Piece Values", convertFloatSiceToInt(weights[768:773]))
	printSlice("EG Piece Values", convertFloatSiceToInt(weights[773:778]))

	printSlice("\nMG Piece Mobility Coefficents", convertFloatSiceToInt(weights[780:784]))
	printSlice("EG Piece Mobility Coefficents", convertFloatSiceToInt(weights[784:788]))

	fmt.Println("\nBishop Pair Bonus MG:", weights[778])
	fmt.Println("Bishop Pair Bonus EG:", weights[779])

	fmt.Println("\nIsolated Pawn Penalty MG:", weights[918])
	fmt.Println("Isolated Pawn Penalty EG:", weights[919])

	fmt.Println("\nDoubled Pawn Penalty MG:", weights[916])
	fmt.Println("Doubled Pawn Penalty EG:", weights[917])

	fmt.Println("\nRook Or Queen On Seventh Bonus MG:", weights[920])
	fmt.Println("Rook Or Queen On Seventh Bonus EG:", weights[921])

	printSlice("\nOuter Ring Attack Coefficents", convertFloatSiceToInt(weights[922:926]))
	printSlice("Inner Ring Attack Coefficents", convertFloatSiceToInt(weights[926:930]))
	fmt.Println("Semi-Open File Next To King Penalty:", weights[930])

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

	prettyPrintPSQT("MG Passed Pawn PST", convertFloatSiceToInt(weights[788:852]))
	prettyPrintPSQT("EG Passed Pawn PST", convertFloatSiceToInt(weights[852:916]))

	fmt.Println()
}

func Tune(infile string, epochs, numPositions int, recordErrorRate bool) {
	weights := loadWeights()
	entries := loadEntries(infile, numPositions)

	gradientsSumsSquared := make([]float64, len(weights))
	beforeErr := computeMSE(entries, weights)

	N := float64(numPositions)
	learningRate := LearningRate

	errors := []float64{beforeErr}
	errorRecordingRate := epochs / 100

	for epoch := 0; epoch < epochs; epoch++ {
		gradients := computeGradient(entries, weights)
		for k, gradient := range gradients {
			leadingCoefficent := (-2 * ScalingFactor) / N
			gradientsSumsSquared[k] += (leadingCoefficent * gradient) * (leadingCoefficent * gradient)
			weights[k] += (leadingCoefficent * gradient) * (-learningRate / math.Sqrt(gradientsSumsSquared[k]+Epsilon))
		}

		fmt.Printf("Epoch number %d completed\n", epoch+1)

		if recordErrorRate && epoch > 0 && epoch%errorRecordingRate == 0 {
			errors = append(errors, computeMSE(entries, weights))
		}
	}

	if recordErrorRate {
		errors = append(errors, computeMSE(entries, weights))
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

	printParameters(weights)
	fmt.Println("Best error before tuning:", beforeErr)
	fmt.Println("Best error after tuning:", computeMSE(entries, weights))
}
