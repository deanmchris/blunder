package engine

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// A file containg the implementation of Blunder's internal board representation.

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

	MaxHistory   uint16 = 100
	MaxGamePlies uint16 = 1024

	NoEPSquare       uint8 = 64
	FENStartPosition       = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0"
	FENKiwiPete            = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

	WhiteKingside  Bitboard = 0x900000000000000
	WhiteQueenside Bitboard = 0x8800000000000000
	BlackKingside  Bitboard = 0x9
	BlackQueenside Bitboard = 0x88

	CastleWKSRand64  uint16 = 768
	CastleWQSRand64  uint16 = 769
	CastleBKSRand64  uint16 = 770
	CastleBQSRand64  uint16 = 771
	EPRand64         uint16 = 772
	SideToMoveRand64 uint16 = 780
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

var PieceTypeToChar map[uint8]rune = map[uint8]rune{
	Pawn:   'p',
	Knight: 'n',
	Bishop: 'b',
	Rook:   'r',
	Queen:  'q',
	King:   'k',
	NoType: '.',
}

var EPDelta [2]int8 = [2]int8{8, -8}

type Piece struct {
	Type  uint8
	Color uint8
}

type BoardState struct {
	Moved          Piece
	Captured       Piece
	CastlingRights Bitboard
	Rule50         uint8
	EPSquare       uint8
}

type Board struct {
	PieceBB [2][6]Bitboard
	SideBB  [2]Bitboard
	Squares [64]Piece
	KingPos [2][6]uint8

	ColorToMove uint8
	GamePly     uint16
	Rule50      uint8

	EPSquare       uint8
	CastlingRights Bitboard
	Hash           uint64

	history    [MaxHistory]BoardState
	historyPly uint16

	Repitions   [MaxGamePlies]uint64
	RepitionPly uint16
}

// Setup a new Board with the internal fields set using a
// Forsyth–Edwards Notation (FEN) line.
func (board *Board) LoadFEN(fen string) {
	board.PieceBB = [2][6]Bitboard{}
	board.SideBB = [2]Bitboard{}
	board.Squares = [64]Piece{}
	board.CastlingRights = 0
	board.Hash = 0

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
			board.putPiece(piece.Type, piece.Color, uint8(sq))

			if char == 'K' {
				board.KingPos[White][King] = uint8(sq)
			} else if char == 'k' {
				board.KingPos[Black][King] = uint8(sq)
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

	if board.ColorToMove == White {
		board.Hash ^= Random64[SideToMoveRand64]
	}

	board.EPSquare = NoEPSquare
	if ep != "-" {
		board.EPSquare = CoordinateToPos(ep)
		if (PawnAttacks[board.ColorToMove^1][board.EPSquare] & board.PieceBB[board.ColorToMove][Pawn]) == 0 {
			board.EPSquare = NoEPSquare
		}

		if board.EPSquare != NoEPSquare {
			board.Hash ^= Random64[EPRand64+uint16(FileOf(board.EPSquare))]
		}
	}

	halfMoveCounter, _ := strconv.Atoi(halfMove)
	board.Rule50 = uint8(halfMoveCounter)

	gamePly, _ := strconv.Atoi(fullMove)
	gamePly *= 2
	if board.ColorToMove == Black {
		gamePly--
	}
	board.GamePly = uint16(gamePly)

	for _, char := range castling {
		switch char {
		case 'K':
			board.CastlingRights.SetBit(E1)
			board.CastlingRights.SetBit(H1)
			board.Hash ^= Random64[CastleWKSRand64]
		case 'Q':
			board.CastlingRights.SetBit(E1)
			board.CastlingRights.SetBit(A1)
			board.Hash ^= Random64[CastleWQSRand64]
		case 'k':
			board.CastlingRights.SetBit(E8)
			board.CastlingRights.SetBit(H8)
			board.Hash ^= Random64[CastleBKSRand64]
		case 'q':
			board.CastlingRights.SetBit(E8)
			board.CastlingRights.SetBit(A8)
			board.Hash ^= Random64[CastleBQSRand64]
		}
	}

	board.RepitionPly = 0
	board.Repitions[board.RepitionPly] = board.Hash
}

// Get the random 64-bit number corresponding to the given piece
// of a certian color and type, on a certian square.
func getPieceHash(pieceType, pieceColor uint8, sq uint8) uint64 {
	if pieceColor == White {
		return Random64[(uint16(pieceType)*2+1)*64+uint16(sq)]
	}
	return Random64[(uint16(pieceType)*2)*64+uint16(sq)]

}

// hash the castling rights into, or out of, the current board Zobrist
// hash. A BoardState object is needed to figure how the castling rights
// changed, and how they need to be updated.
func (board *Board) hashCastlingRights(state *BoardState) {
	if board.CastlingRights != state.CastlingRights {
		if state.CastlingRights&WhiteKingside == WhiteKingside &&
			board.CastlingRights&WhiteKingside != WhiteKingside {

			board.Hash ^= Random64[CastleWKSRand64]
		}
		if state.CastlingRights&WhiteQueenside == WhiteQueenside &&
			board.CastlingRights&WhiteQueenside != WhiteQueenside {
			board.Hash ^= Random64[CastleWQSRand64]
		}
		if state.CastlingRights&BlackKingside == BlackKingside &&
			board.CastlingRights&BlackKingside != BlackKingside {
			board.Hash ^= Random64[CastleBKSRand64]
		}
		if state.CastlingRights&BlackQueenside == BlackQueenside &&
			board.CastlingRights&BlackQueenside != BlackQueenside {
			board.Hash ^= Random64[CastleBQSRand64]
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

	str += fmt.Sprintf("\nzobrist hash: 0x%x\n", board.Hash)
	str += fmt.Sprintf("\nrule 50: %d\n", board.Rule50)
	str += fmt.Sprintf("game ply: %d\n", board.GamePly)
	return str
}

func (board *Board) DoMove(move Move, saveState bool) {
	from, to, movType := move.FromSq(), move.ToSq(), move.MoveType()
	state := &board.history[board.historyPly]

	if saveState {
		board.historyPly++
	}

	state.Moved = board.Squares[from]
	state.Captured = board.Squares[to]
	state.CastlingRights = board.CastlingRights
	state.Rule50 = board.Rule50
	state.EPSquare = board.EPSquare

	if board.EPSquare != NoEPSquare {
		board.Hash ^= Random64[EPRand64+uint16(FileOf(board.EPSquare))]
	}

	board.EPSquare = NoEPSquare

	board.Rule50++
	board.GamePly++

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
		capturePos := uint8(int8(to) + EPDelta[board.ColorToMove])
		state.Captured = board.Squares[capturePos]
		board.removePiece(capturePos)
		board.movePiece(from, to)
		board.Rule50 = 0
	case Attack:
		board.removePiece(to)
		board.movePiece(from, to)
		board.Rule50 = 0
	case DoublePawnPush:
		board.EPSquare = uint8(int8(to) + EPDelta[board.ColorToMove])
		if (PawnAttacks[board.ColorToMove][board.EPSquare] & board.PieceBB[board.ColorToMove^1][Pawn]) == 0 {
			board.EPSquare = NoEPSquare
		}
		fallthrough
	case Quiet:
		board.movePiece(from, to)

		if state.Moved.Type == Pawn {
			board.Rule50 = 0
		}
	}

	board.CastlingRights.ClearBit(from)
	board.CastlingRights.ClearBit(to)
	board.hashCastlingRights(state)

	if board.EPSquare != NoEPSquare {
		board.Hash ^= Random64[EPRand64+uint16(FileOf(board.EPSquare))]
	}

	board.KingPos[board.ColorToMove][state.Moved.Type] = to
	board.ColorToMove ^= 1
	board.Hash ^= Random64[SideToMoveRand64]

	board.RepitionPly++
	board.Repitions[board.RepitionPly] = board.Hash
}

func (board *Board) UndoMove(move Move) {
	board.historyPly--
	state := &board.history[board.historyPly]

	board.hashCastlingRights(state)
	if board.EPSquare != NoEPSquare {
		board.Hash ^= Random64[EPRand64+uint16(FileOf(board.EPSquare))]
	}

	board.CastlingRights = state.CastlingRights
	board.Rule50 = state.Rule50
	board.EPSquare = state.EPSquare

	board.ColorToMove ^= 1
	board.Hash ^= Random64[SideToMoveRand64]

	from, to, movType := move.FromSq(), move.ToSq(), move.MoveType()

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
		capturePos := uint8(int8(to) + EPDelta[board.ColorToMove])
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

	if board.EPSquare != NoEPSquare {
		board.Hash ^= Random64[EPRand64+uint16(FileOf(board.EPSquare))]
	}

	board.KingPos[board.ColorToMove][state.Moved.Type] = from
	board.RepitionPly--
}

// Move a piece from the given square to the given square.
// For this function, the move is guaranteed to be quiet.
func (board *Board) movePiece(from, to uint8) {
	piece := &board.Squares[from]
	board.PieceBB[piece.Color][piece.Type].ClearBit(from)
	board.SideBB[piece.Color].ClearBit(from)
	board.Hash ^= getPieceHash(piece.Type, piece.Color, from)

	board.PieceBB[piece.Color][piece.Type].SetBit(to)
	board.SideBB[piece.Color].SetBit(to)
	board.Hash ^= getPieceHash(piece.Type, piece.Color, to)

	board.Squares[to].Type = piece.Type
	board.Squares[to].Color = piece.Color
	piece.Type = NoType
	piece.Color = NoColor
}

// Put the piece given on the given square
func (board *Board) putPiece(pieceType, pieceColor, to uint8) {
	board.PieceBB[pieceColor][pieceType].SetBit(to)
	board.SideBB[pieceColor].SetBit(to)
	board.Squares[to].Type = pieceType
	board.Squares[to].Color = pieceColor
	board.Hash ^= getPieceHash(pieceType, pieceColor, to)
}

// Remove the piece given on the given square.
func (board *Board) removePiece(from uint8) {
	piece := &board.Squares[from]
	board.PieceBB[piece.Color][piece.Type].ClearBit(from)
	board.SideBB[piece.Color].ClearBit(from)
	board.Hash ^= getPieceHash(piece.Type, piece.Color, from)

	piece.Type = NoType
	piece.Color = NoColor
}

// Test whether or not the king is attacked for the side
// who moved. Called after board.DoMove.
func (board *Board) KingIsAttacked(kingColor uint8) bool {
	return sqIsAttacked(board, kingColor, board.KingPos[kingColor][King])
}

// Given a board square, return it's file.
func FileOf(sq uint8) uint8 {
	return sq % 8
}

// Given a board square, return it's rank.
func RankOf(sq uint8) uint8 {
	return sq / 8
}

// Convert a string board coordinate to its position
// number.
func CoordinateToPos(coordinate string) uint8 {
	file := coordinate[0] - 'a'
	rank := int(coordinate[1]-'0') - 1
	return uint8(rank*8 + int(file))
}

// Convert a position number to a string board coordinate.
func PosToCoordinate(pos uint8) string {
	file := FileOf(pos)
	rank := RankOf(pos)
	return string(rune('a'+file)) + string(rune('0'+rank+1))
}
