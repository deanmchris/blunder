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

var PieceValueMG [5]int16 = [5]int16{100, 300, 300, 500, 900}
var PieceValueEG [5]int16 = [5]int16{100, 300, 300, 500, 900}
var PieceMobilityMG [4]int16 = [4]int16{1, 1, 1, 1}
var PieceMobilityEG [4]int16 = [4]int16{1, 1, 1, 1}

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
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for bishops
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for rook
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for queens
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for kings
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
}

var PSQT_EG [6][64]int16 = [6][64]int16{

	// Piece-square table for pawns
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for knights
	{

		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece-square table for bishops
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for rook
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for queens
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	},

	// Piece square table for kings
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
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

		switch piece.Type {
		case Pawn:
			evalPawn(pos, piece.Color, sq, &mgScores, &egScores)
		case Knight:
			evalKnight(pos, piece.Color, sq, &mgScores, &egScores)
		case Bishop:
			evalBishop(pos, piece.Color, sq, &mgScores, &egScores)
		case Rook:
			evalRook(pos, piece.Color, sq, &mgScores, &egScores)
		case Queen:
			evalQueen(pos, piece.Color, sq, &mgScores, &egScores)
		case King:
			evalKing(pos, piece.Color, sq, &mgScores, &egScores)
		}

		phase -= PhaseValues[piece.Type]
	}

	mgScore := mgScores[pos.SideToMove] - mgScores[pos.SideToMove^1]
	egScore := egScores[pos.SideToMove] - egScores[pos.SideToMove^1]

	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	return int16(((int32(mgScore) * (int32(256) - int32(phase))) + (int32(egScore) * int32(phase))) / int32(256))
}

// Evaluate the score of a knight.
func evalPawn(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	mgScores[color] += PieceValueMG[Pawn] + PSQT_MG[Pawn][FlipSq[color][sq]]
	egScores[color] += PieceValueEG[Pawn] + PSQT_EG[Pawn][FlipSq[color][sq]]
}

// Evaluate the score of a knight.
func evalKnight(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	mgScores[color] += PieceValueMG[Knight] + PSQT_MG[Knight][FlipSq[color][sq]]
	egScores[color] += PieceValueEG[Knight] + PSQT_EG[Knight][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	moves := KnightMoves[sq] & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 4) * PieceMobilityMG[Knight-1]
	egScores[color] += (mobility - 4) * PieceMobilityEG[Knight-1]
}

// Evaluate the score of a bishop.
func evalBishop(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	mgScores[color] += PieceValueMG[Bishop] + PSQT_MG[Bishop][FlipSq[color][sq]]
	egScores[color] += PieceValueEG[Bishop] + PSQT_EG[Bishop][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genBishopMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 7) * PieceMobilityMG[Bishop-1]
	egScores[color] += (mobility - 7) * PieceMobilityEG[Bishop-1]
}

// Evaluate the score of a rook.
func evalRook(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	mgScores[color] += PieceValueMG[Rook] + PSQT_MG[Rook][FlipSq[color][sq]]
	egScores[color] += PieceValueEG[Rook] + PSQT_EG[Rook][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := genRookMoves(sq, allBB) & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 7) * PieceMobilityMG[Rook-1]
	egScores[color] += (mobility - 7) * PieceMobilityEG[Rook-1]
}

// Evaluate the score of a queen.
func evalQueen(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	mgScores[color] += PieceValueMG[Queen] + PSQT_MG[Queen][FlipSq[color][sq]]
	egScores[color] += PieceValueEG[Queen] + PSQT_EG[Queen][FlipSq[color][sq]]

	usBB := pos.SideBB[color]
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]

	moves := (genBishopMoves(sq, allBB) | genRookMoves(sq, allBB)) & ^usBB
	mobility := int16(moves.CountBits())

	mgScores[color] += (mobility - 14) * PieceMobilityMG[Queen-1]
	egScores[color] += (mobility - 14) * PieceMobilityEG[Queen-1]
}

// Evaluate the score of a queen.
func evalKing(pos *Position, color, sq uint8, mgScores, egScores *[2]int16) {
	mgScores[color] += PSQT_MG[King][FlipSq[color][sq]]
	egScores[color] += PSQT_EG[King][FlipSq[color][sq]]
}
