package engine

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// A file containg the implementation of Blunder's internal board representation.

const (
	Pawn   = 0
	Knight = 1
	Bishop = 2
	Rook   = 3
	Queen  = 4
	King   = 5
	NoType = 6

	Black   = 0
	White   = 1
	NoColor = 2

	MaxHistory = 50

	NoEPSquare       = 64
	FENStartPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0"
	FENKiwiPete      = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

	WhiteKingside  uint64 = 0x900000000000000
	WhiteQueenside uint64 = 0x8800000000000000
	BlackKingside  uint64 = 0x9
	BlackQueenside uint64 = 0x88
)
const (
	A1, B1, C1, D1, E1, F1, G1, H1 = 0, 1, 2, 3, 4, 5, 6, 7
	A2, B2, C2, D2, E2, F2, G2, H2 = 8, 9, 10, 11, 12, 13, 14, 15
	A3, B3, C3, D3, E3, F3, G3, H3 = 16, 17, 18, 19, 20, 21, 22, 23
	A4, B4, C4, D4, E4, F4, G4, H4 = 24, 25, 26, 27, 28, 29, 30, 31
	A5, B5, C5, D5, E5, F5, G5, H5 = 32, 33, 34, 35, 36, 37, 38, 39
	A6, B6, C6, D6, E6, F6, G6, H6 = 40, 41, 42, 43, 44, 45, 46, 47
	A7, B7, C7, D7, E7, F7, G7, H7 = 48, 49, 50, 51, 52, 53, 54, 55
	A8, B8, C8, D8, E8, F8, G8, H8 = 56, 57, 58, 59, 60, 61, 62, 63
)

var CharToPiece map[byte]Piece = map[byte]Piece{
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

var PieceTypeToChar map[int]rune = map[int]rune{
	Pawn:   'p',
	Knight: 'n',
	Bishop: 'b',
	Rook:   'r',
	Queen:  'q',
	King:   'k',
	NoType: '.',
}

var EPDelta [2]int = [2]int{8, -8}

type Piece struct {
	Type  int
	Color int
}

type BoardState struct {
	Moved          Piece
	Captured       Piece
	CastlingRights uint64
	Rule50         int
	EPSquare       int
}

type Board struct {
	PieceBB [2][6]uint64
	SideBB  [2]uint64
	Squares [64]Piece
	KingPos [2][6]int

	ColorToMove int
	GamePly     int
	Rule50      int

	EPSquare       int
	CastlingRights uint64

	history    [MaxHistory]BoardState
	historyPly int
}

// Setup a new Board with the internal fields set using a
// Forsyth–Edwards Notation (FEN) line.
func (board *Board) LoadFEN(fen string) {
	board.PieceBB = [2][6]uint64{}
	board.SideBB = [2]uint64{}
	board.Squares = [64]Piece{}
	board.CastlingRights = 0

	fields := strings.Fields(fen)
	pieces := fields[0]
	color := fields[1]
	castling := fields[2]
	ep := fields[3]
	halfMove := fields[4]
	fullMove := fields[5]

	for square := range board.Squares {
		board.Squares[square] = Piece{Type: NoType, Color: NoColor}
	}

	for index, sq := 0, 56; index < len(pieces); index++ {
		char := pieces[index]
		switch char {
		case 'p', 'n', 'b', 'r', 'q', 'k', 'P', 'N', 'B', 'R', 'Q', 'K':
			piece := CharToPiece[char]
			board.putPiece(piece.Type, piece.Color, sq)

			if char == 'K' {
				board.KingPos[White][King] = sq
			} else if char == 'k' {
				board.KingPos[Black][King] = sq
			}

			board.Squares[sq] = piece
			sq++
		case '/':
			sq -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8':
			sq += int(pieces[index] - '0')
		}
	}

	board.ColorToMove = Black
	if color == "w" {
		board.ColorToMove = White
	}

	board.EPSquare = NoEPSquare
	if ep != "-" {
		board.EPSquare = CoordinateToPos(ep)
	}

	halfMoveCounter, _ := strconv.Atoi(halfMove)
	board.Rule50 = halfMoveCounter

	gamePly, _ := strconv.Atoi(fullMove)
	gamePly *= 2
	if board.ColorToMove == Black {
		gamePly--
	}
	board.GamePly = gamePly

	for _, char := range castling {
		switch char {
		case 'K':
			setBit(&board.CastlingRights, E1)
			setBit(&board.CastlingRights, H1)
		case 'Q':
			setBit(&board.CastlingRights, E1)
			setBit(&board.CastlingRights, A1)
		case 'k':
			setBit(&board.CastlingRights, E8)
			setBit(&board.CastlingRights, H8)
		case 'q':
			setBit(&board.CastlingRights, E8)
			setBit(&board.CastlingRights, A8)
		}
	}
}

// Return a pretty string representation of the board. Useful for debugging
// and command-line interaction purposes.
func (board Board) String() (str string) {
	str += "\n"
	for rankStartPos := 56; rankStartPos >= 0; rankStartPos -= 8 {
		str += fmt.Sprintf("%v | ", (rankStartPos/8)+1)
		for index := rankStartPos; index < rankStartPos+8; index++ {
			piece := board.Squares[index]
			pieceChar := PieceTypeToChar[piece.Type]
			if piece.Color == White {
				pieceChar = unicode.ToUpper(pieceChar)
			}
			str += fmt.Sprintf("%c ", pieceChar)
		}
		str += "\n"
	}

	str += "   "
	for fileNo := 0; fileNo < 8; fileNo++ {
		str += "--"
	}

	str += "\n    "
	for _, file := range "abcdefgh" {
		str += fmt.Sprintf("%c ", file)
	}

	str += "\n\n"
	if board.ColorToMove == White {
		str += "turn: white\n"
	} else {
		str += "turn: black\n"
	}

	str += "castling rights: "
	if board.CastlingRights&WhiteKingside == WhiteKingside {
		str += "K"
	}
	if board.CastlingRights&WhiteQueenside == WhiteQueenside {
		str += "Q"
	}
	if board.CastlingRights&BlackKingside == BlackKingside {
		str += "k"
	}
	if board.CastlingRights&BlackQueenside == BlackQueenside {
		str += "q"
	}

	str += "\nen passant: "
	if board.EPSquare == NoEPSquare {
		str += "none"
	} else {
		str += PosToCoordinate(board.EPSquare)
	}

	str += fmt.Sprintf("\nrule 50: %d\n", board.Rule50)
	str += fmt.Sprintf("game ply: %d\n", board.GamePly)
	return str
}

func (board *Board) DoMove(move Move, saveState bool) {
	from, to, movType := FromSq(move), ToSq(move), MoveType(move)
	state := &board.history[board.historyPly]

	if saveState {
		board.historyPly++
	}

	state.Moved = board.Squares[from]
	state.Captured = board.Squares[to]
	state.CastlingRights = board.CastlingRights
	state.Rule50 = board.Rule50
	state.EPSquare = board.EPSquare

	board.EPSquare = NoEPSquare

	switch movType {
	case CastleWKS:
		board.movePiece(E1, G1)
		board.movePiece(H1, F1)
	case CastleWQS:
		board.movePiece(E1, C1)
		board.movePiece(A1, D1)
	case CastleBKS:
		board.movePiece(E8, G8)
		board.movePiece(H8, F8)
	case CastleBQS:
		board.movePiece(E8, C8)
		board.movePiece(A8, D8)
	case KnightPromotion:
		board.removePiece(from)
		if state.Captured.Type != NoType {
			board.removePiece(to)
		}
		board.putPiece(Knight, board.ColorToMove, to)
	case BishopPromotion:
		board.removePiece(from)
		if state.Captured.Type != NoType {
			board.removePiece(to)
		}
		board.putPiece(Bishop, board.ColorToMove, to)
	case RookPromotion:
		board.removePiece(from)
		if state.Captured.Type != NoType {
			board.removePiece(to)
		}
		board.putPiece(Rook, board.ColorToMove, to)
	case QueenPromotion:
		board.removePiece(from)
		if state.Captured.Type != NoType {
			board.removePiece(to)
		}
		board.putPiece(Queen, board.ColorToMove, to)
	case AttackEP:
		capturePos := to + EPDelta[board.ColorToMove]
		state.Captured = board.Squares[capturePos]
		board.removePiece(capturePos)
		board.movePiece(from, to)
		board.Rule50 = 0
	case Attack:
		if state.Captured.Type == King {
			fmt.Println(board)
			panic("illegal king capture")
		}
		board.removePiece(to)
		board.movePiece(from, to)
		board.Rule50 = 0
	case DoublePawnPush:
		board.EPSquare = to + EPDelta[board.ColorToMove]
		fallthrough
	case Quiet:
		board.movePiece(from, to)
	}

	clearBit(&board.CastlingRights, from)
	clearBit(&board.CastlingRights, to)

	board.Rule50++
	board.GamePly++

	board.KingPos[board.ColorToMove][state.Moved.Type] = to
	board.ColorToMove ^= 1
}

func (board *Board) UndoMove(move Move) {
	board.historyPly--
	state := &board.history[board.historyPly]

	board.CastlingRights = state.CastlingRights
	board.Rule50 = state.Rule50
	board.EPSquare = state.EPSquare
	board.ColorToMove ^= 1

	from, to, movType := FromSq(move), ToSq(move), MoveType(move)
	board.GamePly--

	switch movType {
	case CastleWKS:
		board.movePiece(G1, E1)
		board.movePiece(F1, H1)
	case CastleWQS:
		board.movePiece(C1, E1)
		board.movePiece(D1, A1)
	case CastleBKS:
		board.movePiece(G8, E8)
		board.movePiece(F8, H8)
	case CastleBQS:
		board.movePiece(C8, E8)
		board.movePiece(D8, A8)
	case KnightPromotion:
		fallthrough
	case BishopPromotion:
		fallthrough
	case RookPromotion:
		fallthrough
	case QueenPromotion:
		board.removePiece(to)
		if state.Captured.Type != NoType {
			board.putPiece(state.Captured.Type, state.Captured.Color, to)
		}
		board.putPiece(Pawn, board.ColorToMove, from)
	case AttackEP:
		capturePos := to + EPDelta[board.ColorToMove]
		board.movePiece(to, from)
		board.putPiece(Pawn, state.Captured.Color, capturePos)
	case Attack:
		board.removePiece(to)
		board.putPiece(state.Captured.Type, state.Captured.Color, to)
		board.putPiece(state.Moved.Type, board.ColorToMove, from)
	case DoublePawnPush:
		fallthrough
	case Quiet:
		board.movePiece(to, from)
	}

	board.KingPos[board.ColorToMove][state.Moved.Type] = from
}

// Move a piece from the given square to the given square.
// For this function, the move is guaranteed to be quiet.
func (board *Board) movePiece(from, to int) {
	piece := &board.Squares[from]
	clearBit(&board.PieceBB[piece.Color][piece.Type], from)
	clearBit(&board.SideBB[piece.Color], from)
	setBit(&board.PieceBB[piece.Color][piece.Type], to)
	setBit(&board.SideBB[piece.Color], to)

	board.Squares[to].Type = piece.Type
	board.Squares[to].Color = piece.Color
	piece.Type = NoType
	piece.Color = NoColor
}

// Put the piece given on the given square
func (board *Board) putPiece(pieceType, pieceColor, to int) {
	setBit(&board.PieceBB[pieceColor][pieceType], to)
	setBit(&board.SideBB[pieceColor], to)
	board.Squares[to].Type = pieceType
	board.Squares[to].Color = pieceColor
}

// Remove the piece given on the given square.
func (board *Board) removePiece(from int) {
	piece := &board.Squares[from]
	clearBit(&board.PieceBB[piece.Color][piece.Type], from)
	clearBit(&board.SideBB[piece.Color], from)
	piece.Type = NoType
	piece.Color = NoColor
}

// Test whether or not the king is attacked for the side
// who moved. Called after board.DoMove.
func (board *Board) KingIsAttacked(kingColor int) bool {
	return sqIsAttacked(board, kingColor, board.KingPos[kingColor][King])
}

// Given a board square, return it's file.
func FileOf(sq int) int {
	return sq % 8
}

// Given a board square, return it's rank.
func RankOf(sq int) int {
	return sq / 8
}

// Convert a string board coordinate to its position
// number.
func CoordinateToPos(coordinate string) int {
	file := coordinate[0] - 'a'
	rank := int(coordinate[1]-'0') - 1
	return int(rank*8 + int(file))
}

// Convert a position number to a string board coordinate.
func PosToCoordinate(pos int) string {
	file := FileOf(pos)
	rank := RankOf(pos)
	return string(rune('a'+file)) + string(rune('0'+rank+1))
}
