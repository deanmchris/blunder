package engine

const (
	// Constants which map a piece to how much weight it should have on the phase of the game.
	PawnPhase   int16 = 0
	KnightPhase int16 = 1
	BishopPhase int16 = 1
	RookPhase   int16 = 2
	QueenPhase  int16 = 4
	TotalPhase  int16 = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2

	// Constants representing a draw or infinite (checkmate) value.
	Inf int16 = 10000
)

type Eval struct {
	MGScores [2]int16
	EGScores [2]int16
}

var IsolatedPawnMasks [8]Bitboard
var DoubledPawnMasks [2][64]Bitboard

var IsolatedPawnPenaltyMG int16 = 15
var IsolatedPawnPenaltyEG int16 = 15
var DoubledPawnPenaltyMG int16 = 15
var DoubledPawnPenaltyEG int16 = 15

var PieceValueMG = [6]int16{146, 523, 541, 694, 1341}
var PieceValueEG = [6]int16{208, 360, 363, 692, 1381}

var PhaseValues = [6]int16{
	PawnPhase,
	KnightPhase,
	BishopPhase,
	RookPhase,
	QueenPhase,
}

var PSQT_MG = [6][64]int16{
	{
		// MG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		149, 121, 74, 108, 96, 70, 0, -26,
		-36, -65, 8, 18, 83, 100, -4, -66,
		-38, 7, -3, 7, 27, 15, 13, -61,
		-72, -16, -25, 2, 10, -3, -3, -65,
		-68, -22, -23, -30, -4, -11, 18, -29,
		-75, -14, -30, -52, -43, 18, 53, -61,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-228, -3, -28, 49, -26, -106, -2, -54,
		-107, -70, 152, -15, 22, 94, -10, -24,
		-69, 68, 78, 90, 90, 169, 73, 26,
		-16, 22, 50, 72, 64, 116, 25, 19,
		-20, 32, 18, 17, 31, 22, 35, -3,
		-34, -9, 30, 15, 25, 33, 37, -24,
		-59, -43, -18, 0, -5, 35, -26, -22,
		-86, -28, -71, -45, -11, -59, -20, -44,
	},
	{
		// MG Bishop PST
		-15, -17, -104, 0, -45, -2, 3, -56,
		-49, 2, -19, -52, 6, 93, 28, -72,
		-14, 41, 26, 47, 67, 58, 6, 1,
		-26, 5, 6, 52, 37, 70, 0, 7,
		-1, 10, -2, 27, 8, 19, -6, 8,
		1, 17, 7, 24, 9, 43, 11, -5,
		33, 29, 17, -4, 14, 21, 41, 30,
		-21, 6, -26, -52, -42, -13, -7, -49,
	},
	{
		// MG Rook PST
		25, 2, 67, 91, 84, 13, 15, 16,
		28, 56, 100, 68, 88, 24, 2, 39,
		-16, 40, -1, 27, 11, 34, 80, 31,
		-51, 0, 2, 63, 41, 53, 6, -71,
		-51, -7, -46, -6, 14, -30, 10, -35,
		-67, -16, -19, -18, -12, -28, -34, -35,
		-77, -8, -22, -19, -10, -8, -12, -123,
		-21, -18, -8, 7, 16, -3, -46, -37,
	},
	{
		// MG Queen PST
		-23, -56, 68, 45, 51, 40, 93, 39,
		-20, -65, 47, 63, 25, 86, 52, 102,
		11, -7, 16, 8, 22, 87, 39, 61,
		-57, -39, -24, -4, -17, -1, 3, -4,
		-7, -65, -33, -10, -5, 11, 14, 1,
		-14, 17, -8, 7, -17, 2, 23, 36,
		-62, -16, 26, 5, 17, 17, -5, -19,
		-13, -14, -15, 19, -22, -58, -16, -50,
	},
	{
		// MG King PST
		34, 28, 48, 59, -29, -14, 35, -9,
		10, 36, 51, 28, 71, 10, 22, -20,
		-3, 55, 23, -6, 72, 93, 20, 23,
		-35, -11, -20, -24, -60, -44, 0, -9,
		6, -34, -11, -108, -85, -73, -72, -82,
		38, 11, -15, -79, -15, -44, -30, -5,
		77, 25, 32, -78, -66, -49, 30, 44,
		-29, 84, 57, -83, 16, -37, 82, 61,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		233, 212, 170, 131, 148, 121, 248, 243,
		78, 112, 83, 40, 19, -2, 64, 61,
		-14, -35, -52, -62, -73, -64, -43, -37,
		-37, -46, -66, -78, -76, -70, -69, -62,
		-51, -48, -74, -65, -61, -74, -87, -55,
		-44, -52, -46, -49, -51, -69, -55, -73,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-89, -38, 30, -26, 9, 36, -76, -120,
		16, 8, -18, 36, 10, -20, -6, -36,
		0, 1, 37, 26, 12, 0, -1, -41,
		0, 34, 49, 71, 64, 28, 28, 7,
		-16, 12, 47, 64, 35, 59, 9, -26,
		6, 9, 13, 46, 34, 2, -8, -16,
		26, -11, 16, 27, 12, -22, -23, -38,
		-22, -63, -7, -12, 7, -22, -35, -78,
	},
	{
		// EG Bishop PST
		3, -18, -21, -19, -11, 16, 0, -10,
		6, -6, 42, 6, 14, -5, -9, -27,
		8, -12, 8, -18, 1, 24, 21, 17,
		-9, 22, 27, 20, 19, -1, 7, 3,
		12, 8, 41, 30, 35, 19, -3, -21,
		-7, -4, 18, 14, 25, 7, 7, -10,
		-12, -17, -2, 11, 16, -11, -24, -40,
		-30, -21, -43, 14, 1, -19, -14, -12,
	},
	{
		// EG Rook PST
		5, 8, 3, 0, -2, 18, 15, 6,
		8, 3, 0, 15, -2, 15, 9, 0,
		5, -2, 4, 8, 2, -6, -22, -6,
		7, 6, 18, -19, -5, -15, 7, 33,
		9, -2, 19, 8, -5, -2, 3, -1,
		9, 1, -7, 4, 0, -6, 6, -8,
		13, -8, 6, 16, -1, 2, -3, 25,
		-1, 10, 6, 5, -5, -4, 2, -22,
	},
	{
		// EG Queen PST
		-10, 43, -7, 64, 31, -17, -20, 7,
		-41, 73, 10, 4, 53, 37, 28, -90,
		-44, 8, 9, 52, 54, 25, 2, -28,
		22, 50, 27, 29, 62, 75, 67, 38,
		-50, 88, 52, 55, 31, 72, 55, 29,
		-14, -45, 16, -3, 43, 0, -33, -1,
		-28, -22, -47, 2, -34, -17, -45, -13,
		-25, -43, -23, -70, 7, -23, -21, -78,
	},
	{
		// EG King PST
		9, -16, 8, -32, 7, 19, 13, 8,
		-26, -5, 21, 19, 5, 35, 19, 20,
		-25, 11, 24, 31, 21, 21, 37, 6,
		-9, 30, 26, 38, 34, 41, 26, -20,
		-52, 6, 27, 40, 40, 30, 16, -16,
		-34, -16, 0, 30, 25, 28, 9, -5,
		-46, -18, -1, 14, 19, 15, -13, -33,
		-87, -87, -47, -11, -55, -19, -62, -102,
	},
}

// Flip white's perspective to black
var FlipSq [2][64]int = [2][64]int{
	{
		0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
	},

	{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	},
}

// Evaluate a position and give a score, from the perspective of the side to move (
// more positive if it's good for the side to move, otherwise more negative).
func evaluatePos(pos *Position) int16 {
	eval := Eval{MGScores: pos.MGScores, EGScores: pos.EGScores}
	phase := pos.Phase

	pawnBB := pos.Pieces[White][Pawn] | pos.Pieces[Black][Pawn]
	for pawnBB != 0 {
		sq := pawnBB.PopBit()
		piece := pos.Squares[sq]
		evalPawn(pos, piece.Color, sq, &eval)
	}

	mgScore := eval.MGScores[pos.SideToMove] - eval.MGScores[pos.SideToMove^1]
	egScore := eval.EGScores[pos.SideToMove] - eval.EGScores[pos.SideToMove^1]

	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	return int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))
}

// Evaluate the score of a pawn.
func evalPawn(pos *Position, color, sq uint8, eval *Eval) {
	usPawns := pos.Pieces[color][Pawn]

	// Evaluate isolated pawns.
	if IsolatedPawnMasks[FileOf(sq)]&usPawns == 0 {
		eval.MGScores[color] -= IsolatedPawnPenaltyMG
		eval.EGScores[color] -= IsolatedPawnPenaltyEG
	}

	// Evaluate doubled pawns.
	if DoubledPawnMasks[color][sq]&usPawns != 0 {
		eval.MGScores[color] -= DoubledPawnPenaltyMG
		eval.EGScores[color] -= DoubledPawnPenaltyEG
	}
}

func InitEvalBitboards() {
	for file := FileA; file <= FileH; file++ {
		// Create isolated pawn masks.
		fileBB := MaskFile[file]
		mask := (fileBB & ClearFile[FileA]) << 1
		mask |= (fileBB & ClearFile[FileH]) >> 1
		IsolatedPawnMasks[file] = mask
	}

	for sq := uint8(0); sq < 64; sq++ {
		// Create doubled pawns masks.
		fileBB := MaskFile[FileOf(sq)]
		rank := int(RankOf(sq))

		mask := fileBB
		for r := 0; r <= rank; r++ {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[White][sq] = mask

		mask = fileBB
		for r := 7; r >= rank; r-- {
			mask &= ClearRank[r]
		}
		DoubledPawnMasks[Black][sq] = mask
	}
}
