package engine

import (
	"fmt"
)

const (
	MG = 0
	EG = 1

	Material      EvalTerm = 0
	Positional    EvalTerm = 1
	PawnStructure EvalTerm = 2
	Mobility      EvalTerm = 3
	KingSafety    EvalTerm = 4
)

type EvalTerm uint8

type Trace struct {
	Material      [2][2]int16
	Positional    [2][2]int16
	PawnStructure [2][2]int16
	Mobility      [2][2]int16
	KingSafety    [2]int16

	MGScores [2]int16
	EGScores [2]int16

	KingZones        [2]KingZone
	KingAttackPoints [2]uint16
	KingAttackers    [2]uint8
}

func (trace *Trace) AddEvalTerm(evalTerm EvalTerm, mgScore, egScore int16, color uint8) {
	trace.MGScores[color] += mgScore
	trace.EGScores[color] += egScore

	switch evalTerm {
	case Material:
		trace.Material[MG][color] += mgScore
		trace.Material[EG][color] += egScore
	case Positional:
		trace.Positional[MG][color] += mgScore
		trace.Positional[EG][color] += egScore
	case PawnStructure:
		trace.PawnStructure[MG][color] += mgScore
		trace.PawnStructure[EG][color] += egScore
	case Mobility:
		trace.Mobility[MG][color] += mgScore
		trace.Mobility[EG][color] += egScore
	case KingSafety:
		trace.KingSafety[color] += mgScore
	}
}

type Scores struct {
	Material      [2]int16
	Positional    [2]int16
	PawnStructure [2]int16
	Mobility      [2]int16
	KingSafety    [2]int16
	TotalScore    int16
}

func evaluatePosTrace(pos *Position) {
	var trace Trace
	trace.KingZones[White] = KingZones[pos.PieceBB[White][King].Msb()]
	trace.KingZones[Black] = KingZones[pos.PieceBB[Black][King].Msb()]

	phase := TotalPhase
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	for allBB != 0 {
		sq := allBB.PopBit()
		piece := pos.Squares[sq]

		switch piece.Type {
		case Pawn:
			evalPawnTrace(pos, piece.Color, sq, &trace)
		case Knight:
			evalKnightTrace(pos, piece.Color, sq, &trace)
		case Bishop:
			evalBishopTrace(pos, piece.Color, sq, &trace)
		case Rook:
			evalRookTrace(pos, piece.Color, sq, &trace)
		case Queen:
			evalQueenTrace(pos, piece.Color, sq, &trace)
		}

		phase -= PhaseValues[piece.Type]
	}

	evalKingTrace(pos, White, pos.PieceBB[White][King].Msb(), &trace)
	evalKingTrace(pos, Black, pos.PieceBB[Black][King].Msb(), &trace)

	var scores Scores
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase

	scores.Material[White] = interpolate(trace.Material[MG][White], trace.Material[EG][White], phase)
	scores.Material[Black] = interpolate(trace.Material[MG][Black], trace.Material[EG][Black], phase)

	scores.Positional[White] = interpolate(trace.Positional[MG][White], trace.Positional[EG][White], phase)
	scores.Positional[Black] = interpolate(trace.Positional[MG][Black], trace.Positional[EG][Black], phase)

	scores.PawnStructure[White] = interpolate(trace.PawnStructure[MG][White], trace.PawnStructure[EG][White], phase)
	scores.PawnStructure[Black] = interpolate(trace.PawnStructure[MG][Black], trace.PawnStructure[EG][Black], phase)

	scores.Mobility[White] = interpolate(trace.Mobility[MG][White], trace.Mobility[EG][White], phase)
	scores.Mobility[Black] = interpolate(trace.Mobility[MG][Black], trace.Mobility[EG][Black], phase)

	scores.KingSafety[White] = trace.KingSafety[White]
	scores.KingSafety[Black] = trace.KingSafety[Black]

	mgScore := trace.MGScores[pos.SideToMove] - trace.MGScores[pos.SideToMove^1]
	egScore := trace.EGScores[pos.SideToMove] - trace.EGScores[pos.SideToMove^1]
	scores.TotalScore = interpolate(mgScore, egScore, phase)

	printEvalTrace(&scores, pos.SideToMove)

}

func printEvalTraceRow(evalTerm, whiteScore, blackScore, difference string) {
	fmt.Printf(
		"| %s | %s | %s | %s |\n",
		padToCenter(evalTerm, " ", 16),
		padToCenter(whiteScore, " ", 12),
		padToCenter(blackScore, " ", 12),
		padToCenter(difference, " ", 10),
	)
}

func printEvalTraceRowSeparator() {
	fmt.Println("+------------------+--------------+--------------+------------+")
}

func toCentiPawn(score int16) string {
	return fmt.Sprintf("%.2f", float64(score)/float64(100))
}

func printEvalTrace(scores *Scores, stm uint8) {
	fmt.Println()
	printEvalTraceRowSeparator()

	printEvalTraceRow("Evaluation Term", "White Score", "Black Score", "Difference")
	printEvalTraceRowSeparator()

	printEvalTraceRow(
		"Material",
		toCentiPawn(scores.Material[White]),
		toCentiPawn(scores.Material[Black]),
		toCentiPawn(scores.Material[stm]-scores.Material[stm^1]),
	)
	printEvalTraceRowSeparator()

	printEvalTraceRow(
		"Positional",
		toCentiPawn(scores.Positional[White]),
		toCentiPawn(scores.Positional[Black]),
		toCentiPawn(scores.Positional[stm]-scores.Positional[stm^1]),
	)
	printEvalTraceRowSeparator()

	printEvalTraceRow(
		"Pawn Structure",
		toCentiPawn(scores.PawnStructure[White]),
		toCentiPawn(scores.PawnStructure[Black]),
		toCentiPawn(scores.PawnStructure[stm]-scores.PawnStructure[stm^1]),
	)
	printEvalTraceRowSeparator()

	printEvalTraceRow(
		"Mobility",
		toCentiPawn(scores.Mobility[White]),
		toCentiPawn(scores.Mobility[Black]),
		toCentiPawn(scores.Mobility[stm]-scores.Mobility[stm^1]),
	)
	printEvalTraceRowSeparator()

	printEvalTraceRow(
		"King Safety",
		toCentiPawn(scores.KingSafety[White]),
		toCentiPawn(scores.KingSafety[Black]),
		toCentiPawn(scores.KingSafety[stm]-scores.KingSafety[stm^1]),
	)
	printEvalTraceRowSeparator()

	printEvalTraceRow("-", "-", "-", toCentiPawn(scores.TotalScore))
	printEvalTraceRowSeparator()
	fmt.Println()

}

func interpolate(mgScore, egScore, phase int16) int16 {
	return int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))
}

func evalPawnTrace(pos *Position, color, sq uint8, trace *Trace) {
	trace.AddEvalTerm(Material, PieceValueMG[Pawn], PieceValueEG[Pawn], color)
	trace.AddEvalTerm(Positional, PSQT_MG[Pawn][FlipSq[color][sq]], PSQT_EG[Pawn][FlipSq[color][sq]], color)

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	file := FileOf(sq)
	doubled := false

	if IsolatedPawnMasks[file]&usPawns == 0 {
		trace.AddEvalTerm(PawnStructure, IsolatedPawnPenatlyMG, IsolatedPawnPenatlyEG, color)
	}

	if DoubledPawnMasks[color][sq]&usPawns != 0 {
		doubled = true
		trace.AddEvalTerm(PawnStructure, DoubledPawnPenatlyMG, DoubledPawnPenatlyEG, color)
	}

	if PassedPawnMasks[color][sq]&enemyPawns == 0 && !doubled {
		trace.AddEvalTerm(
			PawnStructure,
			PassedPawnBonusMG[FlipRank[color][RankOf(sq)]],
			PassedPawnBonusEG[FlipRank[color][RankOf(sq)]],
			color,
		)
	}
}

func evalKnightTrace(pos *Position, color, sq uint8, trace *Trace) {
	trace.AddEvalTerm(Material, PieceValueMG[Knight], PieceValueEG[Knight], color)
	trace.AddEvalTerm(Positional, PSQT_MG[Knight][FlipSq[color][sq]], PSQT_EG[Knight][FlipSq[color][sq]], color)

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	if MiniorOutpostMasks[color][sq]&enemyPawns == 0 &&
		PawnAttacks[color^1][sq]&usPawns != 0 &&
		FlipRank[color][RankOf(sq)] >= Rank5 {
		trace.AddEvalTerm(Positional, KnightOutpostBonusMG, KnightOutpostBonusEG, color)
	}

	usBB := pos.SideBB[color]
	moves := KnightMoves[sq] & ^usBB
	mobility := int16(moves.CountBits())

	trace.AddEvalTerm(
		Mobility,
		(mobility-4)*PieceMobilityMG[Knight-1],
		(mobility-4)*PieceMobilityEG[Knight-1],
		color,
	)

	outerRingAttacks := moves & trace.KingZones[color^1].OuterRing
	innerRingAttacks := moves & trace.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		trace.KingAttackers[color]++
		trace.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(MinorAttackOuterRing)
		trace.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(MinorAttackInnerRing)
	}
}

func evalBishopTrace(pos *Position, color, sq uint8, trace *Trace) {
	trace.AddEvalTerm(Material, PieceValueMG[Bishop], PieceValueEG[Bishop], color)
	trace.AddEvalTerm(Positional, PSQT_MG[Bishop][FlipSq[color][sq]], PSQT_EG[Bishop][FlipSq[color][sq]], color)

	usPawns := pos.PieceBB[color][Pawn]
	enemyPawns := pos.PieceBB[color^1][Pawn]

	if MiniorOutpostMasks[color][sq]&enemyPawns == 0 &&
		PawnAttacks[color^1][sq]&usPawns != 0 &&
		FlipRank[color][RankOf(sq)] >= Rank5 {
		trace.AddEvalTerm(Positional, BishopOutpostBonusMG, BishopOutpostBonusEG, color)
	}

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genBishopMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	trace.AddEvalTerm(
		Mobility,
		(mobility-7)*PieceMobilityMG[Bishop-1],
		(mobility-7)*PieceMobilityEG[Bishop-1],
		color,
	)

	outerRingAttacks := moves & trace.KingZones[color^1].OuterRing
	innerRingAttacks := moves & trace.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		trace.KingAttackers[color]++
		trace.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(MinorAttackOuterRing)
		trace.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(MinorAttackInnerRing)
	}
}

func evalRookTrace(pos *Position, color, sq uint8, trace *Trace) {
	trace.AddEvalTerm(Material, PieceValueMG[Rook], PieceValueEG[Rook], color)
	trace.AddEvalTerm(Positional, PSQT_MG[Rook][FlipSq[color][sq]], PSQT_EG[Rook][FlipSq[color][sq]], color)

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genRookMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	trace.AddEvalTerm(
		Mobility,
		(mobility-7)*PieceMobilityMG[Rook-1],
		(mobility-7)*PieceMobilityEG[Rook-1],
		color,
	)
	outerRingAttacks := moves & trace.KingZones[color^1].OuterRing
	innerRingAttacks := moves & trace.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		trace.KingAttackers[color]++
		trace.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(RookAttackOuterRing)
		trace.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(RookAttackInnerRing)
	}
}

func evalQueenTrace(pos *Position, color, sq uint8, trace *Trace) {
	trace.AddEvalTerm(Material, PieceValueMG[Queen], PieceValueEG[Queen], color)
	trace.AddEvalTerm(Positional, PSQT_MG[Queen][FlipSq[color][sq]], PSQT_EG[Queen][FlipSq[color][sq]], color)

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := (genBishopMoves(sq, allBB) | genRookMoves(sq, allBB)) & ^usBB
	mobility := int16(moves.CountBits())

	trace.AddEvalTerm(
		Mobility,
		(mobility-14)*PieceMobilityMG[Queen-1],
		(mobility-14)*PieceMobilityEG[Queen-1],
		color,
	)

	outerRingAttacks := moves & trace.KingZones[color^1].OuterRing
	innerRingAttacks := moves & trace.KingZones[color^1].InnerRing

	if outerRingAttacks != 0 || innerRingAttacks != 0 {
		trace.KingAttackers[color]++
		trace.KingAttackPoints[color] += uint16(outerRingAttacks.CountBits()) * uint16(QueenAttackOuterRing)
		trace.KingAttackPoints[color] += uint16(innerRingAttacks.CountBits()) * uint16(QueenAttackInnerRing)
	}
}

func evalKingTrace(pos *Position, color, sq uint8, trace *Trace) {
	trace.AddEvalTerm(Positional, PSQT_MG[King][FlipSq[color][sq]], PSQT_EG[King][FlipSq[color][sq]], color)

	enemyPoints := InitKingSafety[FlipSq[color][sq]] + trace.KingAttackPoints[color^1]
	kingFile := MaskFile[FileOf(sq)]
	usPawns := pos.PieceBB[color][Pawn]

	leftFile := ((kingFile & ClearFile[FileA]) << 1)
	rightFile := ((kingFile & ClearFile[FileH]) >> 1)

	if kingFile&usPawns == 0 {
		enemyPoints += uint16(SemiOpenFileNextToKingPenalty)
	}

	if leftFile != 0 && leftFile&usPawns == 0 {
		enemyPoints += uint16(SemiOpenFileNextToKingPenalty)
	}

	if rightFile != 0 && rightFile&usPawns == 0 {
		enemyPoints += uint16(SemiOpenFileNextToKingPenalty)
	}

	enemyPoints = min_u16(enemyPoints, uint16(len(KingAttackTable)-1))
	if trace.KingAttackers[color^1] >= 2 && pos.PieceBB[color^1][Queen] != 0 {
		trace.AddEvalTerm(KingSafety, -KingAttackTable[enemyPoints], 0, color)
	}
}
