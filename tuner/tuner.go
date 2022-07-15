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
	Iterations         = 2000
	NumWeights         = 936
	NumSafetyEvalTerms = 9
	ScalingFactor      = 0.01
	Epsilon            = 0.00000001
	LearningRate       = 0.5

	Draw     float64 = 0.5
	WhiteWin float64 = 1.0
	BlackWin float64 = 0.0
)

type Indexes struct {
	EG_PSQT_StartIndex            uint16
	MG_Material_StartIndex        uint16
	EG_Material_StartIndex        uint16
	MG_BishopPairIndex            uint16
	EG_BishopPairIndex            uint16
	MG_IsoPawnIndex               uint16
	EG_IsoPawnIndex               uint16
	MG_DoubledPawnIndex           uint16
	EG_DoubledPawnIndex           uint16
	MG_PassedPawn_PSQT_StartIndex uint16
	EG_PassedPawn_PSQT_StartIndex uint16
	MG_KnightOutpostIndex         uint16
	EG_KnightOutpostIndex         uint16
	MG_MobilityStartIndex         uint16
	EG_MobilityStartIndex         uint16
	MG_BishopOutpostIndex         uint16
	EG_BishopOutpostIndex         uint16
	EG_RookOrQueenOnSeventhIndex  uint16
	MG_RookOnOpenFileIndex        uint16
	TempoBonusIndex               uint16
	KingSafteyStartIndex          uint16
}

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
func loadWeights() (weights []float64, indexes Indexes) {
	tempWeights := make([]int16, NumWeights)
	weights = make([]float64, NumWeights)

	index := uint16(0)
	indexes.EG_PSQT_StartIndex = 64 * 6

	for piece := engine.Pawn; piece <= engine.King; piece++ {
		copy(tempWeights[index:index+64], engine.PSQT_MG[piece][:])
		copy(tempWeights[384+index:384+index+64], engine.PSQT_EG[piece][:])
		index += 64
	}
	index *= 2

	indexes.MG_Material_StartIndex = index
	indexes.EG_Material_StartIndex = index + 5

	copy(tempWeights[index:index+5], engine.PieceValueMG[0:5])
	copy(tempWeights[index+5:index+10], engine.PieceValueEG[0:5])
	index += 10

	indexes.MG_BishopPairIndex = index
	indexes.EG_BishopPairIndex = index + 1

	tempWeights[index] = engine.BishopPairBonusMG
	tempWeights[index+1] = engine.BishopPairBonusEG
	index += 2

	indexes.MG_MobilityStartIndex = index
	indexes.EG_MobilityStartIndex = index + 4

	copy(tempWeights[index:index+4], engine.PieceMobilityMG[1:5])
	copy(tempWeights[index+4:index+8], engine.PieceMobilityEG[1:5])
	index += 8

	indexes.MG_PassedPawn_PSQT_StartIndex = index
	indexes.EG_PassedPawn_PSQT_StartIndex = index + 64

	copy(tempWeights[index:index+64], engine.PassedPawnPSQT_MG[:])
	copy(tempWeights[index+64:index+128], engine.PassedPawnPSQT_EG[:])
	index += 128

	indexes.MG_DoubledPawnIndex = index
	indexes.EG_DoubledPawnIndex = index + 1
	indexes.MG_IsoPawnIndex = index + 2
	indexes.EG_IsoPawnIndex = index + 3

	tempWeights[index] = engine.DoubledPawnPenatlyMG
	tempWeights[index+1] = engine.DoubledPawnPenatlyEG
	tempWeights[index+2] = engine.IsolatedPawnPenatlyMG
	tempWeights[index+3] = engine.IsolatedPawnPenatlyEG
	index += 4

	indexes.EG_RookOrQueenOnSeventhIndex = index
	indexes.MG_KnightOutpostIndex = index + 1
	indexes.EG_KnightOutpostIndex = index + 2

	tempWeights[index] = engine.RookOrQueenOnSeventhBonusEG
	tempWeights[index+1] = engine.KnightOnOutpostBonusMG
	tempWeights[index+2] = engine.KnightOnOutpostBonusEG
	index += 3

	indexes.MG_BishopOutpostIndex = index
	indexes.EG_BishopOutpostIndex = index + 1

	tempWeights[index] = engine.BishopOutPostBonusMG
	tempWeights[index+1] = engine.BishopOutPostBonusEG
	index += 2

	indexes.MG_RookOnOpenFileIndex = index
	indexes.TempoBonusIndex = index + 1

	tempWeights[index] = engine.RookOnOpenFileBonusMG
	tempWeights[index+1] = engine.TempoBonusMG
	index += 2

	indexes.KingSafteyStartIndex = index
	copy(tempWeights[index:index+4], engine.OuterRingAttackPoints[1:5])
	copy(tempWeights[index+4:index+8], engine.InnerRingAttackPoints[1:5])
	tempWeights[index+8] = engine.SemiOpenFileNextToKingPenalty

	for i := range tempWeights {
		weights[i] = float64(tempWeights[i])
	}

	return weights, indexes
}

// Load the weights for tuning, setting them to reasonable default values.
func loadDefaultWeights() (weights []float64, indexes Indexes) {
	tempWeights := make([]int16, NumWeights)
	weights = make([]float64, NumWeights)

	indexes.EG_PSQT_StartIndex = 64 * 6
	index := uint16(768)

	indexes.MG_Material_StartIndex = index
	indexes.EG_Material_StartIndex = index + 5

	copy(tempWeights[index:index+5], []int16{100, 300, 310, 500, 950})
	copy(tempWeights[index+5:index+10], []int16{100, 300, 310, 500, 950})
	index += 10

	indexes.MG_BishopPairIndex = index
	indexes.EG_BishopPairIndex = index + 1

	tempWeights[index] = 10
	tempWeights[index+1] = 10
	index += 2

	indexes.MG_MobilityStartIndex = index
	indexes.EG_MobilityStartIndex = index + 4

	copy(tempWeights[index:index+4], []int16{1, 1, 1, 1})
	copy(tempWeights[index+4:index+8], []int16{1, 1, 1, 1})
	index += 8
	index += 128

	indexes.MG_DoubledPawnIndex = index
	indexes.EG_DoubledPawnIndex = index + 1
	indexes.MG_IsoPawnIndex = index + 2
	indexes.EG_IsoPawnIndex = index + 3

	tempWeights[index] = 5
	tempWeights[index+1] = 10
	tempWeights[index+2] = 5
	tempWeights[index+3] = 10
	index += 4

	indexes.EG_RookOrQueenOnSeventhIndex = index
	indexes.MG_KnightOutpostIndex = index + 1
	indexes.EG_KnightOutpostIndex = index + 2

	tempWeights[index] = 15
	tempWeights[index+1] = 15
	tempWeights[index+2] = 20
	index += 3

	indexes.MG_BishopOutpostIndex = index
	indexes.EG_BishopOutpostIndex = index + 1

	tempWeights[index] = 10
	tempWeights[index+1] = 15
	index += 2

	indexes.MG_RookOnOpenFileIndex = index
	indexes.TempoBonusIndex = index + 1

	tempWeights[index] = 10
	tempWeights[index+1] = 10
	index += 2

	indexes.KingSafteyStartIndex = index
	copy(tempWeights[index:index+4], []int16{1, 1, 1, 1})
	copy(tempWeights[index+4:index+8], []int16{1, 1, 1, 1})
	tempWeights[index+8] = 1

	for i := range tempWeights {
		weights[i] = float64(tempWeights[i])
	}

	return weights, indexes
}

// Load the given number of positions from the training set file.
func loadEntries(infile string, numPositions int, indexes Indexes) (entries []Entry) {
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

		normalCoefficents, safetyCoefficents := getCoefficents(&pos, indexes)
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
func getCoefficents(pos *engine.Position, indexes Indexes) (normalCoefficents []Coefficent, safetyCoefficents [][]Coefficent) {
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

		getPSQT_Coefficents(rawNormCoefficents, indexes, piece, sq, sign, mgPhase, egPhase)

		switch piece.Type {
		case engine.Pawn:
			getPawnCoefficents(pos, rawNormCoefficents, indexes, sq, mgPhase, egPhase, sign)
		case engine.Knight:
			getKnightCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, indexes, sq, mgPhase, egPhase, sign)
		case engine.Bishop:
			getBishopCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, indexes, sq, mgPhase, egPhase, sign)
		case engine.Rook:
			getRookCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, indexes, sq, mgPhase, egPhase, sign)
		case engine.Queen:
			getQueenCoefficents(pos, rawNormCoefficents, rawSafetyCoefficents, &safetyTracer, indexes, sq, mgPhase, egPhase, sign)
		}
	}

	getMaterialCoeffficents(pos, rawNormCoefficents, indexes, mgPhase, egPhase)
	getBishopPairCoefficents(pos, rawNormCoefficents, indexes, mgPhase, egPhase)

	getPawnShieldCoefficents(pos, pos.Pieces[engine.White][engine.King].Msb(), engine.White, rawSafetyCoefficents)
	getPawnShieldCoefficents(pos, pos.Pieces[engine.Black][engine.King].Msb(), engine.Black, rawSafetyCoefficents)

	getTempoBonusCoefficent(pos, rawNormCoefficents, indexes, mgPhase)

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
			Coefficent{Idx: indexes.KingSafteyStartIndex + uint16(i), Value: value},
		)
	}

	for i, coefficent := range rawSafetyCoefficents[engine.Black] {
		value := float64(0)
		if safetyTracer.KingAttackers[engine.Black] >= 2 && pos.Pieces[engine.Black][engine.Queen] != 0 {
			value = coefficent
		}

		safetyCoefficents[engine.Black] = append(
			safetyCoefficents[engine.Black],
			Coefficent{Idx: indexes.KingSafteyStartIndex + uint16(i), Value: value},
		)
	}

	return normalCoefficents, safetyCoefficents
}

// Get the piece square table coefficents of the position.
func getPSQT_Coefficents(coefficents []float64, indexes Indexes, piece engine.Piece, sq uint8, sign, mgPhase, egPhase float64) {
	mgIndex := uint16(piece.Type)*64 + uint16(engine.FlipSq[piece.Color][sq])
	egIndex := indexes.EG_PSQT_StartIndex + mgIndex
	coefficents[mgIndex] += sign * mgPhase
	coefficents[egIndex] += sign * egPhase
}

// Get the material coefficents of the position.
func getMaterialCoeffficents(pos *engine.Position, coefficents []float64, indexes Indexes, mgPhase, egPhase float64) {
	for piece := uint16(0); piece <= 4; piece++ {
		coefficents[indexes.MG_Material_StartIndex+piece] = float64(
			(pos.Pieces[engine.White][piece].CountBits() - pos.Pieces[engine.Black][piece].CountBits()),
		) * mgPhase
		coefficents[indexes.EG_Material_StartIndex+piece] = float64(
			(pos.Pieces[engine.White][piece].CountBits() - pos.Pieces[engine.Black][piece].CountBits()),
		) * egPhase
	}
}

// Get the bishop pair coefficents of the position.
func getBishopPairCoefficents(pos *engine.Position, coefficents []float64, indexes Indexes, mgPhase, egPhase float64) {
	if pos.Pieces[engine.White][engine.Bishop].CountBits() >= 2 {
		coefficents[indexes.MG_BishopPairIndex] += mgPhase
		coefficents[indexes.EG_BishopPairIndex] += egPhase
	}

	if pos.Pieces[engine.Black][engine.Bishop].CountBits() >= 2 {
		coefficents[indexes.MG_BishopPairIndex] -= mgPhase
		coefficents[indexes.EG_BishopPairIndex] -= egPhase
	}
}

// Get the coefficents of the position related to the given pawn.
func getPawnCoefficents(pos *engine.Position, norm []float64, indexes Indexes, sq uint8, mgPhase, egPhase, sign float64) {
	piece := pos.Squares[sq]
	enemyPawns := pos.Pieces[piece.Color^1][engine.Pawn]
	usPawns := pos.Pieces[piece.Color][engine.Pawn]

	// Evaluate isolated pawns.
	if engine.IsolatedPawnMasks[engine.FileOf(sq)]&usPawns == 0 {
		norm[indexes.MG_IsoPawnIndex] -= sign * mgPhase
		norm[indexes.EG_IsoPawnIndex] -= sign * egPhase
	}

	// Evaluate doubled pawns.
	if engine.DoubledPawnMasks[piece.Color][sq]&usPawns != 0 {
		norm[indexes.MG_DoubledPawnIndex] -= sign * mgPhase
		norm[indexes.EG_DoubledPawnIndex] -= sign * egPhase
	}

	// Evaluate passed pawns, but make sure they're not behind a friendly pawn.
	if engine.PassedPawnMasks[piece.Color][sq]&enemyPawns == 0 &&
		usPawns&engine.DoubledPawnMasks[piece.Color][sq] == 0 {

		mgIndex := indexes.MG_PassedPawn_PSQT_StartIndex + uint16(engine.FlipSq[piece.Color][sq])
		egIndex := indexes.EG_PassedPawn_PSQT_StartIndex + uint16(engine.FlipSq[piece.Color][sq])

		norm[mgIndex] += sign * mgPhase
		norm[egIndex] += sign * egPhase
	}
}

// Get the coefficents of the position related to the given knight.
func getKnightCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer, indexes Indexes,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]

	usPawns := pos.Pieces[piece.Color][engine.Pawn]
	enemyPawns := pos.Pieces[piece.Color^1][engine.Pawn]

	if engine.OutpostMasks[piece.Color][sq]&enemyPawns == 0 &&
		engine.PawnAttacks[piece.Color^1][sq]&usPawns != 0 &&
		engine.FlipRank[piece.Color][engine.RankOf(sq)] >= engine.Rank5 {

		norm[indexes.MG_KnightOutpostIndex] += sign * mgPhase
		norm[indexes.EG_KnightOutpostIndex] += sign * egPhase
	}

	moves := engine.KnightMoves[sq] & ^usBB
	safeMoves := moves

	for enemyPawns != 0 {
		sq := enemyPawns.PopBit()
		safeMoves &= ^engine.PawnAttacks[piece.Color^1][sq]
	}

	mobility := float64(safeMoves.CountBits())
	norm[indexes.MG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 4) * sign * mgPhase
	norm[indexes.EG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 4) * sign * egPhase

	outerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].OuterRing
	innerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		safetyTracer.KingAttackers[piece.Color]++
		safety[piece.Color][piece.Type-1] += float64(outerRingAttacks.CountBits())
		safety[piece.Color][4+piece.Type-1] += float64(innerRingAttacks.CountBits())
	}
}

// Get the coefficents of the position related to the given bishop.
func getBishopCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer, indexes Indexes,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]
	allBB := usBB | pos.Sides[piece.Color^1]

	usPawns := pos.Pieces[piece.Color][engine.Pawn]
	enemyPawns := pos.Pieces[piece.Color^1][engine.Pawn]

	if engine.OutpostMasks[piece.Color][sq]&enemyPawns == 0 &&
		engine.PawnAttacks[piece.Color^1][sq]&usPawns != 0 &&
		engine.FlipRank[piece.Color][engine.RankOf(sq)] >= engine.Rank5 {

		norm[indexes.MG_BishopOutpostIndex] += sign * mgPhase
		norm[indexes.EG_BishopOutpostIndex] += sign * egPhase
	}

	moves := engine.GenBishopMoves(sq, allBB) & ^usBB
	mobility := float64(moves.CountBits())

	norm[indexes.MG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 7) * sign * mgPhase
	norm[indexes.EG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 7) * sign * egPhase

	outerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].OuterRing
	innerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		safetyTracer.KingAttackers[piece.Color]++
		safety[piece.Color][piece.Type-1] += float64(outerRingAttacks.CountBits())
		safety[piece.Color][4+piece.Type-1] += float64(innerRingAttacks.CountBits())
	}
}

// Get the coefficents of the position related to the given rook.
func getRookCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer, indexes Indexes,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]
	allBB := usBB | pos.Sides[piece.Color^1]

	enemyKingSq := pos.Pieces[piece.Color^1][engine.King].Msb()
	if engine.FlipRank[piece.Color][engine.RankOf(sq)] == engine.Rank7 &&
		engine.FlipRank[piece.Color][engine.RankOf(enemyKingSq)] >= engine.Rank7 {

		norm[indexes.EG_RookOrQueenOnSeventhIndex] += sign * egPhase
	}

	pawns := pos.Pieces[engine.White][engine.Pawn] | pos.Pieces[engine.Black][engine.Pawn]
	if engine.MaskFile[engine.FileOf(sq)]&pawns == 0 {
		norm[indexes.MG_RookOnOpenFileIndex] += sign * mgPhase
	}

	moves := engine.GenRookMoves(sq, allBB) & ^usBB
	mobility := float64(moves.CountBits())

	norm[indexes.MG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 7) * sign * float64(mgPhase)
	norm[indexes.EG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 7) * sign * float64(egPhase)

	outerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].OuterRing
	innerRingAttacks := moves & safetyTracer.KingZones[piece.Color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		safetyTracer.KingAttackers[piece.Color]++
		safety[piece.Color][piece.Type-1] += float64(outerRingAttacks.CountBits())
		safety[piece.Color][4+piece.Type-1] += float64(innerRingAttacks.CountBits())
	}
}

// Get the coefficents of the position related to the given queen.
func getQueenCoefficents(pos *engine.Position, norm []float64, safety [][]float64, safetyTracer *SafetyTracer, indexes Indexes,
	sq uint8, mgPhase, egPhase, sign float64) {

	piece := pos.Squares[sq]
	usBB := pos.Sides[piece.Color]
	allBB := usBB | pos.Sides[piece.Color^1]

	enemyKingSq := pos.Pieces[piece.Color^1][engine.King].Msb()
	if engine.FlipRank[piece.Color][engine.RankOf(sq)] == engine.Rank7 &&
		engine.FlipRank[piece.Color][engine.RankOf(enemyKingSq)] >= engine.Rank7 {

		norm[indexes.EG_RookOrQueenOnSeventhIndex] += sign * egPhase
	}

	moves := (engine.GenBishopMoves(sq, allBB) | engine.GenRookMoves(sq, allBB)) & ^usBB
	mobility := float64(moves.CountBits())

	norm[indexes.MG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 14) * sign * float64(mgPhase)
	norm[indexes.EG_MobilityStartIndex+uint16(piece.Type)-1] += (mobility - 14) * sign * float64(egPhase)

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

// Compute the tempo bonus coefficent
func getTempoBonusCoefficent(pos *engine.Position, coefficents []float64, indexes Indexes, mgPhase float64) {
	sign := float64(1)
	if pos.SideToMove != engine.White {
		sign = -1
	}
	coefficents[indexes.TempoBonusIndex] = sign * mgPhase
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
func evaluate(weights []float64, normalCoefficents []Coefficent, safetyCoefficents [][]Coefficent, indexes Indexes, mgPhase float64) (score float64) {
	for i := range normalCoefficents {
		coefficent := &normalCoefficents[i]
		score += weights[coefficent.Idx] * coefficent.Value
	}

	whiteSafety := computeSafetyDotProduct(weights[indexes.KingSafteyStartIndex:NumWeights], safetyCoefficents[engine.White])
	blackSafety := computeSafetyDotProduct(weights[indexes.KingSafteyStartIndex:NumWeights], safetyCoefficents[engine.Black])

	whiteSafety = ((whiteSafety * whiteSafety) / 4) * mgPhase
	blackSafety = ((blackSafety * blackSafety) / 4) * mgPhase

	return score + whiteSafety - blackSafety
}

func computeGradientNumerically(entries []Entry, weights []float64, indexes Indexes, epsilon float64) (gradients []float64) {
	N := float64(len(entries))
	gradients = make([]float64, len(weights))
	epsilonAddedErrSums := make([]float64, len(entries))
	epsilonSubtractedErrSums := make([]float64, len(entries))

	for i := range entries {
		for k := range weights {
			weights[k] += epsilon

			score := evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, indexes, entries[i].MGPhase)
			sigmoid := 1 / (1 + math.Exp(-(ScalingFactor * score)))
			err := entries[i].Outcome - sigmoid
			epsilonAddedErrSums[k] += math.Pow(err, 2)

			weights[k] -= epsilon * 2

			score = evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, indexes, entries[i].MGPhase)
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

func computeGradient(entries []Entry, weights []float64, indexes Indexes) (gradients []float64) {
	gradients = make([]float64, NumWeights)

	for i := range entries {
		score := evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, indexes, entries[i].MGPhase)
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

		whiteSafety := computeSafetyDotProduct(weights[indexes.KingSafteyStartIndex:NumWeights], entries[i].SafetyCoefficents[engine.White])
		blackSafety := computeSafetyDotProduct(weights[indexes.KingSafteyStartIndex:NumWeights], entries[i].SafetyCoefficents[engine.Black])

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

func computeMSE(entries []Entry, weights []float64, indexes Indexes) (errSum float64) {
	for i := range entries {
		score := evaluate(weights, entries[i].NormalCoefficents, entries[i].SafetyCoefficents, indexes, entries[i].MGPhase)
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

func printParameters(weights []float64, indexes Indexes) {
	printSlice("\nMG Piece Values", convertFloatSiceToInt(weights[indexes.MG_Material_StartIndex:indexes.MG_Material_StartIndex+5]))
	printSlice("EG Piece Values", convertFloatSiceToInt(weights[indexes.EG_Material_StartIndex:indexes.EG_Material_StartIndex+5]))

	printSlice("\nMG Piece Mobility Coefficents", convertFloatSiceToInt(weights[indexes.MG_MobilityStartIndex:indexes.MG_MobilityStartIndex+4]))
	printSlice("EG Piece Mobility Coefficents", convertFloatSiceToInt(weights[indexes.EG_MobilityStartIndex:indexes.EG_MobilityStartIndex+4]))

	fmt.Println("\nBishop Pair Bonus MG:", weights[indexes.MG_BishopPairIndex])
	fmt.Println("Bishop Pair Bonus EG:", weights[indexes.EG_BishopPairIndex])

	fmt.Println("\nIsolated Pawn Penalty MG:", weights[indexes.MG_IsoPawnIndex])
	fmt.Println("Isolated Pawn Penalty EG:", weights[indexes.EG_IsoPawnIndex])

	fmt.Println("\nDoubled Pawn Penalty MG:", weights[indexes.MG_DoubledPawnIndex])
	fmt.Println("Doubled Pawn Penalty EG:", weights[indexes.EG_DoubledPawnIndex])

	fmt.Println("\nRook Or Queen On Seventh Bonus EG:", weights[indexes.EG_RookOrQueenOnSeventhIndex])

	fmt.Println("\nKnight On Outpost Bonus MG:", weights[indexes.MG_KnightOutpostIndex])
	fmt.Println("Knight On Outpost Bonus EG:", weights[indexes.EG_KnightOutpostIndex])

	fmt.Println("\nRook On Open File Bonus MG:", weights[indexes.MG_RookOnOpenFileIndex])
	fmt.Println("Tempo Bonus MG:", weights[indexes.TempoBonusIndex])

	fmt.Println("\nBishop On Outpost Bonus MG:", weights[indexes.MG_BishopOutpostIndex])
	fmt.Println("Bishop On Outpost Bonus EG:", weights[indexes.EG_BishopOutpostIndex])

	printSlice("\nOuter Ring Attack Coefficents", convertFloatSiceToInt(weights[indexes.KingSafteyStartIndex:indexes.KingSafteyStartIndex+4]))
	printSlice("Inner Ring Attack Coefficents", convertFloatSiceToInt(weights[indexes.KingSafteyStartIndex+4:indexes.KingSafteyStartIndex+8]))
	fmt.Println("Semi-Open File Next To King Penalty:", weights[indexes.KingSafteyStartIndex+8])

	pieceNames := []string{"Pawn", "Knight", "Bishop", "Rook", "Queen", "King"}

	for piece, index := uint8(0), uint16(0); piece <= engine.King; piece, index = piece+1, index+64 {
		tableName := fmt.Sprintf("MG %s PST", pieceNames[piece])
		prettyPrintPSQT(tableName, convertFloatSiceToInt(weights[index:index+64]))
	}

	for piece, index := uint8(0), uint16(0); piece <= engine.King; piece, index = piece+1, index+64 {
		tableName := fmt.Sprintf("MG %s PST", pieceNames[piece])
		prettyPrintPSQT(tableName, convertFloatSiceToInt(weights[indexes.EG_PSQT_StartIndex+index:indexes.EG_PSQT_StartIndex+index+64]))
	}

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

	prettyPrintPSQT("MG Passed Pawn PST", convertFloatSiceToInt(weights[indexes.MG_PassedPawn_PSQT_StartIndex:indexes.MG_PassedPawn_PSQT_StartIndex+64]))
	prettyPrintPSQT("EG Passed Pawn PST", convertFloatSiceToInt(weights[indexes.EG_PassedPawn_PSQT_StartIndex:indexes.EG_PassedPawn_PSQT_StartIndex+64]))

	fmt.Println()
}

func Tune(infile string, epochs, numPositions int, recordErrorRate bool, useDefaultWeights bool) {
	var weights []float64
	var indexes Indexes

	if useDefaultWeights {
		weights, indexes = loadDefaultWeights()
	} else {
		weights, indexes = loadWeights()
	}

	printParameters(weights, indexes)
	return

	entries := loadEntries(infile, numPositions, indexes)

	gradientsSumsSquared := make([]float64, len(weights))
	beforeErr := computeMSE(entries, weights, indexes)

	N := float64(numPositions)
	learningRate := LearningRate

	errors := []float64{beforeErr}
	errorRecordingRate := epochs / 100

	for epoch := 0; epoch < epochs; epoch++ {
		gradients := computeGradient(entries, weights, indexes)
		for k, gradient := range gradients {
			leadingCoefficent := (-2 * ScalingFactor) / N
			gradientsSumsSquared[k] += (leadingCoefficent * gradient) * (leadingCoefficent * gradient)
			weights[k] += (leadingCoefficent * gradient) * (-learningRate / math.Sqrt(gradientsSumsSquared[k]+Epsilon))
		}

		fmt.Printf("Epoch number %d completed\n", epoch+1)

		if recordErrorRate && epoch > 0 && epoch%errorRecordingRate == 0 {
			errors = append(errors, computeMSE(entries, weights, indexes))
		}
	}

	if recordErrorRate {
		errors = append(errors, computeMSE(entries, weights, indexes))
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

	printParameters(weights, indexes)
	fmt.Println("Best error before tuning:", beforeErr)
	fmt.Println("Best error after tuning:", computeMSE(entries, weights, indexes))
}
