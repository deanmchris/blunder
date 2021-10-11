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
	Inf            int16 = 10000
	MiddleGameDraw int16 = 25
	EndGameDraw    int16 = 0
)

var KnightMobility int16 = 1
var BishopMobility int16 = 4
var RookMobilityMG int16 = 6
var RookMobilityEG int16 = 2
var QueenMobilityMG int16 = 1
var QueenMobilityEG int16 = 8

var PhaseValues [6]int16 = [6]int16{
	PawnPhase,
	KnightPhase,
	BishopPhase,
	RookPhase,
	QueenPhase,
}

// Endgame and middlegame piece square tables, with piece values builtin.
//
// https://www.chessprogramming.org/Piece-Square_Tables
//
var PSQT_MG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		124, 108, 101, 107, 99, 101, 100, 98,
		202, 236, 206, 214, 209, 219, 145, 98,
		89, 94, 113, 122, 151, 153, 111, 70,
		72, 102, 90, 110, 112, 102, 106, 64,
		60, 85, 86, 102, 106, 98, 99, 61,
		68, 85, 85, 81, 96, 96, 122, 78,
		63, 90, 66, 71, 78, 115, 128, 71,
		123, 101, 104, 104, 101, 100, 98, 100,
	},

	// Piece-square table for knights
	{
		137, 234, 290, 276, 354, 228, 290, 203,
		225, 264, 384, 345, 343, 383, 317, 301,
		266, 371, 350, 370, 399, 433, 388, 366,
		299, 330, 326, 368, 348, 388, 331, 339,
		302, 325, 332, 325, 342, 332, 335, 303,
		294, 305, 327, 330, 338, 335, 340, 298,
		280, 266, 300, 320, 322, 335, 308, 306,
		211, 303, 262, 290, 314, 293, 303, 279,
	},

	// Piece-square table for bishops
	{
		285, 341, 249, 280, 298, 305, 325, 334,
		312, 351, 317, 325, 371, 394, 362, 294,
		328, 371, 377, 373, 382, 397, 375, 344,
		342, 342, 352, 387, 371, 378, 345, 344,
		343, 361, 351, 366, 374, 345, 353, 355,
		349, 367, 362, 353, 356, 375, 365, 359,
		361, 370, 362, 347, 359, 372, 385, 352,
		323, 354, 346, 339, 346, 345, 312, 328,
	},

	// Piece square table for rook
	{
		502, 521, 501, 532, 528, 484, 506, 505,
		512, 513, 536, 533, 573, 557, 503, 519,
		476, 499, 498, 511, 483, 535, 540, 495,
		460, 470, 479, 501, 497, 513, 482, 461,
		448, 458, 472, 480, 490, 472, 495, 454,
		446, 463, 477, 474, 489, 483, 481, 456,
		455, 478, 473, 482, 498, 506, 483, 420,
		480, 484, 500, 509, 512, 504, 459, 477,
	},

	// Piece square table for queens
	{
		875, 904, 924, 928, 970, 981, 963, 951,
		892, 865, 902, 912, 906, 974, 941, 968,
		907, 890, 913, 904, 938, 963, 954, 971,
		884, 882, 887, 887, 910, 930, 891, 914,
		910, 884, 904, 901, 905, 904, 912, 909,
		904, 920, 906, 914, 916, 920, 928, 921,
		890, 914, 931, 926, 933, 955, 918, 924,
		924, 922, 923, 938, 920, 910, 902, 881,
	},

	// Piece square table for kings
	{
		-16, 122, 100, 59, -5, 33, 57, 73,
		84, 61, 53, 96, 54, 61, 21, -17,
		56, 60, 66, 37, 49, 82, 79, 14,
		-4, 17, 34, 17, 23, 26, 21, -14,
		-19, 59, -2, -30, -22, -14, -17, -31,
		45, 17, -10, -29, -18, -16, 2, -2,
		27, 27, 2, -60, -40, -12, 19, 23,
		-4, 56, 31, -56, 15, -20, 39, 31,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		119, 97, 99, 105, 101, 103, 100, 100,
		277, 270, 251, 222, 226, 227, 251, 294,
		192, 204, 182, 160, 152, 142, 187, 186,
		130, 119, 110, 97, 91, 99, 109, 113,
		110, 104, 90, 86, 84, 84, 94, 94,
		98, 103, 87, 97, 93, 87, 88, 84,
		108, 104, 103, 103, 106, 91, 91, 85,
		123, 107, 101, 98, 100, 102, 99, 99,
	},

	// Piece-square table for knights
	{

		259, 247, 265, 246, 247, 257, 221, 202,
		264, 283, 257, 277, 272, 261, 263, 240,
		260, 263, 294, 285, 273, 273, 256, 234,
		270, 287, 308, 301, 305, 293, 291, 265,
		273, 273, 300, 307, 301, 297, 285, 269,
		260, 280, 280, 291, 291, 281, 263, 261,
		250, 263, 277, 277, 280, 265, 267, 240,
		263, 239, 266, 268, 271, 275, 236, 219,
	},

	// Piece-square table for bishops
	{
		312, 293, 311, 296, 308, 299, 294, 285,
		313, 296, 304, 294, 294, 289, 301, 305,
		311, 293, 297, 293, 292, 293, 295, 317,
		308, 308, 309, 299, 299, 302, 300, 311,
		303, 299, 308, 307, 298, 305, 301, 294,
		299, 301, 306, 310, 313, 297, 302, 296,
		297, 290, 296, 304, 308, 292, 292, 280,
		289, 300, 297, 305, 308, 299, 309, 295,
	},

	// Piece square table for rook
	{
		511, 503, 510, 502, 505, 517, 505, 505,
		510, 509, 504, 502, 486, 496, 509, 499,
		509, 503, 503, 504, 505, 488, 492, 499,
		506, 505, 513, 493, 497, 499, 501, 513,
		509, 512, 506, 502, 497, 492, 492, 500,
		502, 501, 495, 504, 492, 486, 489, 491,
		499, 493, 499, 503, 491, 490, 487, 511,
		491, 500, 496, 491, 489, 487, 503, 480,
	},

	// Piece square table for queens
	{
		929, 974, 956, 951, 956, 968, 948, 971,
		931, 951, 969, 964, 969, 946, 964, 950,
		924, 931, 918, 967, 966, 953, 949, 954,
		953, 953, 944, 956, 975, 964, 1010, 977,
		912, 968, 940, 967, 949, 949, 975, 963,
		938, 895, 935, 924, 938, 949, 952, 952,
		923, 910, 897, 915, 914, 908, 894, 917,
		911, 900, 903, 886, 932, 913, 944, 907,
	},

	// Piece square table for kings
	{
		-93, -29, -9, -28, 2, 34, 21, 11,
		-15, 26, 17, 17, 14, 49, 40, 35,
		27, 34, 32, 24, 35, 48, 59, 35,
		12, 39, 38, 41, 35, 45, 36, 22,
		-4, 6, 40, 47, 46, 39, 25, 8,
		-9, 9, 28, 39, 41, 35, 25, 11,
		-14, 3, 20, 35, 37, 24, 12, -4,
		-35, -28, -9, 11, -12, 10, -12, -35,
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
func EvaluatePos(pos *Position) int16 {
	var mgScores, egScores [2]int16
	phase := TotalPhase

	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]
	for allBB != 0 {
		sq := allBB.PopBit()
		piece := pos.Squares[sq]

		mgScores[piece.Color] += PSQT_MG[piece.Type][FlipSq[piece.Color][sq]]
		egScores[piece.Color] += PSQT_EG[piece.Type][FlipSq[piece.Color][sq]]

		if piece.Type == Knight {
			evalKnight(pos, piece.Color, sq, &mgScores, &egScores)
		} else if piece.Type == Bishop {
			evalBishop(pos, piece.Color, sq, &mgScores, &egScores)
		} else if piece.Type == Rook {
			evalRook(pos, piece.Color, sq, &mgScores, &egScores)
		} else if piece.Type == Queen {
			evalQueen(pos, piece.Color, sq, &mgScores, &egScores)
		}

		phase -= PhaseValues[piece.Type]
	}

	mgScore := mgScores[pos.SideToMove] - mgScores[pos.SideToMove^1]
	egScore := egScores[pos.SideToMove] - egScores[pos.SideToMove^1]

	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	return int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))
}

// Evaluate the score of a knight.
func evalKnight(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	usBB := pos.SideBB[color]
	moves := KnightMoves[sq] & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 4) * KnightMobility
	egScores[color] += (mobility - 4) * KnightMobility
}

// Evaluate the score of a bishop.
func evalBishop(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genBishopMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 7) * BishopMobility
	egScores[color] += (mobility - 7) * BishopMobility
}

// Evaluate the score of a rook.
func evalRook(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genRookMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 7) * RookMobilityMG
	egScores[color] += (mobility - 7) * RookMobilityEG
}

// Evaluate the score of a queen.
func evalQueen(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := (genBishopMoves(sq, allBB) | genRookMoves(sq, allBB)) & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 14) * QueenMobilityMG
	egScores[color] += (mobility - 14) * QueenMobilityEG
}
