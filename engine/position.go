package engine

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	Pawn   uint8 = 0
	Knight uint8 = 1
	Bishop uint8 = 2
	Rook   uint8 = 3
	Queen  uint8 = 4
	King   uint8 = 5
	NoType uint8 = 6

	Black   uint8 = 0
	White   uint8 = 1
	NoColor uint8 = 2

	WhiteKingsideRight  uint8 = 0x8
	WhiteQueensideRight uint8 = 0x4
	BlackKingsideRight  uint8 = 0x2
	BlackQueensideRight uint8 = 0x1

	FENStartPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0"
	FENKiwiPete      = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

	A1, B1, C1, D1, E1, F1, G1, H1 = 0, 1, 2, 3, 4, 5, 6, 7
	A2, B2, C2, D2, E2, F2, G2, H2 = 8, 9, 10, 11, 12, 13, 14, 15
	A3, B3, C3, D3, E3, F3, G3, H3 = 16, 17, 18, 19, 20, 21, 22, 23
	A4, B4, C4, D4, E4, F4, G4, H4 = 24, 25, 26, 27, 28, 29, 30, 31
	A5, B5, C5, D5, E5, F5, G5, H5 = 32, 33, 34, 35, 36, 37, 38, 39
	A6, B6, C6, D6, E6, F6, G6, H6 = 40, 41, 42, 43, 44, 45, 46, 47
	A7, B7, C7, D7, E7, F7, G7, H7 = 48, 49, 50, 51, 52, 53, 54, 55
	A8, B8, C8, D8, E8, F8, G8, H8 = 56, 57, 58, 59, 60, 61, 62, 63
	NoSq                           = 64

	NorthDelta = 8
	SouthDelta = -8

	MaxGamePly = 1024
)

var PositionHistories [MaxGamePly]uint64
var HistoryPly uint16

var CastlingRightRemovalMasks = [64]uint8{
	0xb, 0xf, 0xf, 0xf, 0x3, 0xf, 0xf, 0x7,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xe, 0xf, 0xf, 0xf, 0xc, 0xf, 0xf, 0xd,
}

var CastlingRookSq = map[uint8][2]uint8{
	G1: {H1, F1},
	C1: {A1, D1},
	G8: {H8, F8},
	C8: {A8, D8},
}

var CharToPiece = map[byte]Piece{
	'P': {Pawn, White},
	'N': {Knight, White},
	'B': {Bishop, White},
	'R': {Rook, White},
	'Q': {Queen, White},
	'K': {King, White},
	'p': {Pawn, Black},
	'n': {Knight, Black},
	'b': {Bishop, Black},
	'r': {Rook, Black},
	'q': {Queen, Black},
	'k': {King, Black},
}

var PieceTypeToChar = map[uint8]rune{
	Pawn:   'i',
	Knight: 'n',
	Bishop: 'b',
	Rook:   'r',
	Queen:  'q',
	King:   'k',
	NoType: '.',
}

type Piece struct {
	Type  uint8
	Color uint8
}

type IrreversibleState struct {
	CastlingRights uint8
	EPSq           uint8
	Rule50         uint8

	Captured Piece
	Moved    Piece
}

type Position struct {
	PieceBB [2][6]Bitboard
	SideBB  [2]Bitboard
	Squares [64]Piece

	// The castling rights are keep track of using 4-bits:
	// 00001000 = white kingside castling right
	// 00000100 = white queenside castling right
	// 00000010 = black kingside castling right
	// 00000001 = black queenside castling right
	CastlingRights uint8

	Hash uint64

	SideToMove uint8
	EPSq       uint8

	Ply    uint16
	Rule50 uint8

	prevStates [100]IrreversibleState
	StatePly   uint8
}

func (pos *Position) LoadFEN(fen string) {
	pos.PieceBB = [2][6]Bitboard{}
	pos.SideBB = [2]Bitboard{}
	pos.Squares = [64]Piece{}
	pos.CastlingRights = 0

	for square := range pos.Squares {
		pos.Squares[square] = Piece{Type: NoType, Color: NoColor}
	}

	fields := strings.Fields(fen)
	pieces := fields[0]
	color := fields[1]
	castling := fields[2]
	ep := fields[3]
	halfMove := fields[4]
	fullMove := fields[5]

	for index, sq := 0, 56; index < len(pieces); index++ {
		char := pieces[index]
		switch char {
		case 'p', 'n', 'b', 'r', 'q', 'k', 'P', 'N', 'B', 'R', 'Q', 'K':
			piece := CharToPiece[char]
			pos.putPiece(piece.Type, piece.Color, uint8(sq))
			pos.Squares[sq] = piece
			sq++
		case '/':
			sq -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8':
			sq += int(pieces[index] - '0')
		}
	}

	pos.SideToMove = Black
	if color == "w" {
		pos.SideToMove = White
	}

	pos.EPSq = NoSq
	if ep != "-" {
		pos.EPSq = CoordinateToPos(ep)
		// Don't set the enpassant square of a position if there's no enemy
		// pawn to capture on the next move.
		if (PawnAttacks[pos.SideToMove^1][pos.EPSq] & pos.PieceBB[pos.SideToMove][Pawn]) == 0 {
			pos.EPSq = NoSq
		}
	}

	halfMoveCounter, _ := strconv.Atoi(halfMove)
	pos.Rule50 = uint8(halfMoveCounter)

	gamePly, _ := strconv.Atoi(fullMove)
	gamePly *= 2
	if pos.SideToMove == Black {
		gamePly--
	}
	pos.Ply = uint16(gamePly)

	for _, char := range castling {
		switch char {
		case 'K':
			pos.CastlingRights |= WhiteKingsideRight
		case 'Q':
			pos.CastlingRights |= WhiteQueensideRight
		case 'k':
			pos.CastlingRights |= BlackKingsideRight
		case 'q':
			pos.CastlingRights |= BlackQueensideRight
		}
	}

	pos.Hash = 0
	pos.Hash = Zobrist.GenHash(pos)

	HistoryPly = 0
	PositionHistories[HistoryPly] = pos.Hash
}

func (pos Position) GenFEN() string {
	positionStr := strings.Builder{}

	for rankStartPos := 56; rankStartPos >= 0; rankStartPos -= 8 {
		emptySquares := 0
		for sq := rankStartPos; sq < rankStartPos+8; sq++ {
			piece := pos.Squares[sq]
			if piece.Type == NoType {
				emptySquares++
			} else {
				// If we have some consecutive empty squares, add then to the FEN
				// string board, and reset the empty squares counter.
				if emptySquares > 0 {
					positionStr.WriteString(strconv.Itoa(emptySquares))
					emptySquares = 0
				}

				piece := pos.Squares[sq]
				pieceChar := PieceTypeToChar[piece.Type]
				if piece.Color == White {
					pieceChar = unicode.ToUpper(pieceChar)
				}

				// In FEN strings pawns are represented with p's not i's
				if pieceChar == 'i' {
					pieceChar = 'p'
				} else if pieceChar == 'I' {
					pieceChar = 'P'
				}

				positionStr.WriteRune(pieceChar)
			}
		}

		if emptySquares > 0 {
			positionStr.WriteString(strconv.Itoa(emptySquares))
			emptySquares = 0
		}

		positionStr.WriteString("/")

	}

	sideToMove := ""
	castlingRights := ""
	epSquare := ""

	if pos.SideToMove == White {
		sideToMove = "w"
	} else {
		sideToMove = "b"
	}

	if pos.CastlingRights&WhiteKingsideRight != 0 {
		castlingRights += "K"
	}
	if pos.CastlingRights&WhiteQueensideRight != 0 {
		castlingRights += "Q"
	}
	if pos.CastlingRights&BlackKingsideRight != 0 {
		castlingRights += "k"
	}
	if pos.CastlingRights&BlackQueensideRight != 0 {
		castlingRights += "q"
	}

	if castlingRights == "" {
		castlingRights = "-"
	}

	if pos.EPSq == NoSq {
		epSquare = "-"
	} else {
		epSquare = posToCoordinate(pos.EPSq)
	}

	fullMoveCount := pos.Ply / 2
	if pos.Ply%2 != 0 {
		fullMoveCount = pos.Ply/2 + 1
	}

	return fmt.Sprintf(
		"%s %s %s %s %d %d",
		strings.TrimSuffix(positionStr.String(), "/"),
		sideToMove, castlingRights, epSquare,
		pos.Rule50, fullMoveCount,
	)
}

func (pos Position) String() (boardAsString string) {
	boardAsString += "\n"
	for rankStartPos := 56; rankStartPos >= 0; rankStartPos -= 8 {
		boardAsString += fmt.Sprintf("%v | ", (rankStartPos/8)+1)
		for index := rankStartPos; index < rankStartPos+8; index++ {
			piece := pos.Squares[index]
			pieceChar := PieceTypeToChar[piece.Type]
			if piece.Color == White {
				pieceChar = unicode.ToUpper(pieceChar)
			}
			boardAsString += fmt.Sprintf("%c ", pieceChar)
		}
		boardAsString += "\n"
	}

	boardAsString += "   ----------------"
	boardAsString += "\n    a b c d e f g h"

	boardAsString += "\n\n"
	if pos.SideToMove == White {
		boardAsString += "turn: white\n"
	} else {
		boardAsString += "turn: black\n"
	}

	boardAsString += "castling rights: "
	if pos.CastlingRights&WhiteKingsideRight != 0 {
		boardAsString += "K"
	}
	if pos.CastlingRights&WhiteQueensideRight != 0 {
		boardAsString += "Q"
	}
	if pos.CastlingRights&BlackKingsideRight != 0 {
		boardAsString += "k"
	}
	if pos.CastlingRights&BlackQueensideRight != 0 {
		boardAsString += "q"
	}

	boardAsString += "\nen passant: "
	if pos.EPSq != NoSq {
		boardAsString += posToCoordinate(pos.EPSq)
	}

	boardAsString += fmt.Sprintf("\nfen: %s", pos.GenFEN())
	boardAsString += fmt.Sprintf("\nzobrist hash: 0x%x", pos.Hash)
	boardAsString += fmt.Sprintf("\nrule 50: %d\n", pos.Rule50)
	boardAsString += fmt.Sprintf("game ply: %d\n", pos.Ply)
	return boardAsString
}

func (pos *Position) DoMove(move Move) bool {
	from := move.FromSq()
	to := move.ToSq()
	moveType := move.MoveType()
	flag := move.Flag()

	state := IrreversibleState{
		CastlingRights: pos.CastlingRights,
		EPSq:           pos.EPSq,
		Rule50:         pos.Rule50,
		Captured:       pos.Squares[to],
		Moved:          pos.Squares[from],
	}

	pos.Ply++
	pos.Rule50++

	// Clear the enpassant square from the zobrist hash.
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.EPSq = NoSq

	pos.removePiece(from)
	if moveType == Quiet {
		pos.putPiece(state.Moved.Type, state.Moved.Color, to)

		if state.Moved.Type == Pawn {
			if abs16(int16(from)-int16(to)) == 16 {
				// Don't set the enpassant square if there's no enemy pawn in position
				// to capture it on the next move.
				pos.EPSq = uint8(int8(to) - getPawnPushDelta(pos.SideToMove))
				if PawnAttacks[pos.SideToMove][pos.EPSq]&pos.PieceBB[pos.SideToMove^1][Pawn] == 0 {
					pos.EPSq = NoSq
				}
			}

			pos.Rule50 = 0
		}
	} else if moveType == Attack {
		if flag == AttackEPFlag {
			capSq := uint8(int8(to) - getPawnPushDelta(pos.SideToMove))
			state.Captured = pos.Squares[capSq]

			pos.removePiece(capSq)
			pos.putPiece(Pawn, pos.SideToMove, to)

		} else {
			pos.removePiece(to)
			pos.putPiece(state.Moved.Type, state.Moved.Color, to)
		}

		pos.Rule50 = 0
	} else if moveType == Castle {
		pos.putPiece(state.Moved.Type, state.Moved.Color, to)
		rookFrom, rookTo := CastlingRookSq[to][0], CastlingRookSq[to][1]
		pos.removePiece(rookFrom)
		pos.putPiece(Rook, pos.SideToMove, rookTo)
	} else {
		if state.Captured.Type != NoType {
			pos.removePiece(to)
		}

		pos.putPiece(uint8(flag+1), pos.SideToMove, to)
		pos.Rule50 = 0
	}

	// Remove the current castling rights.
	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)

	// Update the castling rights and the zobrist hash with the new castling rights.
	pos.CastlingRights = pos.CastlingRights & CastlingRightRemovalMasks[from] & CastlingRightRemovalMasks[to]
	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)

	// Update the zobrist hash if the en passant square was set
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)

	pos.prevStates[pos.StatePly] = state
	pos.StatePly++

	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	HistoryPly++
	PositionHistories[HistoryPly] = pos.Hash

	return !sqIsAttacked(pos, pos.SideToMove^1, pos.PieceBB[pos.SideToMove^1][King].Msb())
}

func (pos *Position) UndoMove(move Move) {
	pos.StatePly--
	state := pos.prevStates[pos.StatePly]

	// remove the en passant zobrist number if there was one in the position
	// we're undoing, and remove the castling rights zobrist number.
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)

	// Remove the current positions from the position history
	HistoryPly--

	// Restore the irreversible aspects of the position using the State object.
	pos.CastlingRights = state.CastlingRights
	pos.EPSq = state.EPSq
	pos.Rule50 = state.Rule50

	// Update the zobrist hash with the restored values castling rights
	// and en passant square.
	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)

	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	pos.Ply--

	from := move.FromSq()
	to := move.ToSq()
	moveType := move.MoveType()
	flag := move.Flag()

	pos.putPiece(state.Moved.Type, state.Moved.Color, from)
	if moveType == Quiet {
		pos.removePiece(to)
	} else if moveType == Attack {
		if flag == AttackEPFlag {
			capSq := uint8(int8(to) - getPawnPushDelta(pos.SideToMove))
			pos.removePiece(to)
			pos.putPiece(Pawn, state.Captured.Color, capSq)
		} else {
			pos.removePiece(to)
			pos.putPiece(state.Captured.Type, state.Captured.Color, to)
		}
	} else if moveType == Castle {
		pos.removePiece(to)
		rookFrom, rookTo := CastlingRookSq[to][0], CastlingRookSq[to][1]
		pos.removePiece(rookTo)
		pos.putPiece(Rook, pos.SideToMove, rookFrom)
	} else {
		pos.removePiece(to)
		if state.Captured.Type != NoType {
			pos.putPiece(state.Captured.Type, state.Captured.Color, to)
		}
	}
}

func (pos *Position) DoNullMove() {
	state := IrreversibleState{
		CastlingRights: pos.CastlingRights,
		EPSq:           pos.EPSq,
		Rule50:         pos.Rule50,
	}

	pos.prevStates[pos.StatePly] = state
	pos.StatePly++

	// Clear the en passant square and en passant zobrist number
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.EPSq = NoSq

	// Set the fifty move rule counter to 0, since we're
	// making a null-move.
	pos.Rule50 = 0

	pos.Ply++

	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	HistoryPly++
	PositionHistories[HistoryPly] = pos.Hash

}

func (pos *Position) UndoNullMove() {
	pos.StatePly--
	state := pos.prevStates[pos.StatePly]

	pos.CastlingRights = state.CastlingRights
	pos.EPSq = state.EPSq
	pos.Rule50 = state.Rule50

	pos.Ply--

	// Update the zobrist hash with the restored en passant square.
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)

	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	HistoryPly--
}

func (pos *Position) putPiece(pieceType, pieceColor, to uint8) {
	pos.PieceBB[pieceColor][pieceType].SetBit(to)
	pos.SideBB[pieceColor].SetBit(to)

	pos.Squares[to].Type = pieceType
	pos.Squares[to].Color = pieceColor
	pos.Hash ^= Zobrist.PieceNumber(pieceType, pieceColor, to)
}

func (pos *Position) removePiece(from uint8) {
	piece := &pos.Squares[from]
	pos.PieceBB[piece.Color][piece.Type].ClearBit(from)
	pos.SideBB[piece.Color].ClearBit(from)

	pos.Hash ^= Zobrist.PieceNumber(piece.Type, piece.Color, from)
	piece.Type = NoType
	piece.Color = NoColor
}

func (pos *Position) InCheck() bool {
	return sqIsAttacked(
		pos,
		pos.SideToMove,
		pos.PieceBB[pos.SideToMove][King].Msb())
}

func (pos *Position) IsEndgame() bool {
	pawnMaterial := int16(pos.PieceBB[White][Pawn].CountBits()+pos.PieceBB[Black][Pawn].CountBits()) * 100
	knightMaterial := int16(pos.PieceBB[White][Knight].CountBits()+pos.PieceBB[Black][Knight].CountBits()) * 320
	bishopMaterial := int16(pos.PieceBB[White][Bishop].CountBits()+pos.PieceBB[Black][Bishop].CountBits()) * 330
	rookMaterial := int16(pos.PieceBB[White][Rook].CountBits()+pos.PieceBB[Black][Rook].CountBits()) * 500
	queenMaterial := int16(pos.PieceBB[White][Queen].CountBits()+pos.PieceBB[Black][Queen].CountBits()) * 950
	return (pawnMaterial + knightMaterial + bishopMaterial + rookMaterial + queenMaterial) < 2600
}

// Evaluate if an endgame is drawn.
func (pos *Position) EndgameIsDrawn() bool {
	whiteKnights := pos.PieceBB[White][Knight].CountBits()
	whiteBishops := pos.PieceBB[White][Bishop].CountBits()

	blackKnights := pos.PieceBB[Black][Knight].CountBits()
	blackBishops := pos.PieceBB[Black][Bishop].CountBits()

	pawns := pos.PieceBB[White][Pawn].CountBits() + pos.PieceBB[Black][Pawn].CountBits()
	knights := whiteKnights + blackKnights
	bishops := whiteBishops + blackBishops
	rooks := pos.PieceBB[White][Rook].CountBits() + pos.PieceBB[Black][Rook].CountBits()
	queens := pos.PieceBB[White][Queen].CountBits() + pos.PieceBB[Black][Queen].CountBits()

	majors := rooks + queens
	miniors := knights + bishops

	if pawns+majors+miniors == 0 {
		// KvK => draw
		return true
	} else if majors+pawns == 0 {
		if miniors == 1 {
			// K & minior v K & minior => draw
			return true
		} else if miniors == 2 && whiteKnights == 1 && blackKnights == 1 {
			// KNvKN => draw
			return true
		} else if miniors == 2 && whiteBishops == 1 && blackBishops == 1 {
			// KBvKB => draw when only when bishops are the same color
			whiteBishopSq := pos.PieceBB[White][Bishop].Msb()
			blackBishopSq := pos.PieceBB[Black][Bishop].Msb()
			return sqIsDark(whiteBishopSq) == sqIsDark(blackBishopSq)
		}
	}
	return false
}

func (pos *Position) MoveIsPseduoLegal(move Move) bool {
	fromSq, toSq := move.FromSq(), move.ToSq()
	moved := pos.Squares[fromSq]
	captured := pos.Squares[toSq]

	toBB := SquareBB[toSq]
	allBB := pos.SideBB[White] | pos.SideBB[Black]
	sideToMove := pos.SideToMove

	if moved.Color != sideToMove ||
		captured.Type == King ||
		captured.Color == sideToMove {
		return false
	}

	if moved.Type == Pawn {
		if fromSq > 55 || fromSq < 8 {
			return false
		}

		// Credit to the Stockfish team for the idea behind this section of code to
		// verify pseduo-legal pawn moves.
		if ((PawnAttacks[sideToMove][fromSq] & toBB & allBB) == 0) &&
			!((fromSq+uint8(getPawnPushDelta(sideToMove)) == toSq) && (captured.Type == NoType)) &&
			!((fromSq+uint8(getPawnPushDelta(sideToMove)*2) == toSq) &&
				captured.Type == NoType &&
				pos.Squares[toSq-uint8(getPawnPushDelta(sideToMove))].Type == NoType &&
				canDoublePush(fromSq, sideToMove)) {
			return false
		}
	} else {
		if (moved.Type == Knight && ((KnightMoves[fromSq] & toBB) == 0)) ||
			(moved.Type == Bishop && ((genBishopMoves(fromSq, allBB) & toBB) == 0)) ||
			(moved.Type == Rook && ((genRookMoves(fromSq, allBB) & toBB) == 0)) ||
			(moved.Type == Queen && (((genBishopMoves(fromSq, allBB) | genRookMoves(fromSq, allBB)) & toBB) == 0)) ||
			(moved.Type == King && ((KingMoves[fromSq] & toBB) == 0)) {
			return false
		}
	}

	return true
}

func (pos *Position) NoMajorsOrMiniors() bool {
	knights := pos.PieceBB[White][Knight].CountBits() + pos.PieceBB[Black][Knight].CountBits()
	bishops := pos.PieceBB[White][Bishop].CountBits() + pos.PieceBB[Black][Bishop].CountBits()
	rook := pos.PieceBB[White][Rook].CountBits() + pos.PieceBB[Black][Rook].CountBits()
	queen := pos.PieceBB[White][Queen].CountBits() + pos.PieceBB[Black][Queen].CountBits()
	return knights+bishops+rook+queen == 0
}

func getPawnPushDelta(color uint8) int8 {
	if color == White {
		return NorthDelta
	}
	return SouthDelta
}

func canDoublePush(fromSq uint8, color uint8) bool {
	if color == White {
		return RankOf(fromSq) == Rank2
	}
	return RankOf(fromSq) == Rank6
}
