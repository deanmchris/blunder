package engine

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	// Constants representing each piece type. The value of the constants
	// are selected so they can be used in Position.PieceBB to index the
	// bitboards representing the given piece.
	Pawn   uint8 = 0
	Knight uint8 = 1
	Bishop uint8 = 2
	Rook   uint8 = 3
	Queen  uint8 = 4
	King   uint8 = 5
	NoType uint8 = 6

	// Constants representing each piece color. The value of the constants
	// are selected so they can be used in Position.PieceBB and Position.SideBB to
	// index the bitboards representing the given piece of the given color, or
	// the given color.
	Black   uint8 = 0
	White   uint8 = 1
	NoColor uint8 = 2

	// Constants representing the four castling rights. Each constant is set to a number
	// with a single high bit, corresponding to each castling right.
	WhiteKingsideRight  uint8 = 0x8
	WhiteQueensideRight uint8 = 0x4
	BlackKingsideRight  uint8 = 0x2
	BlackQueensideRight uint8 = 0x1

	// Common fen strings used in debugging and initalizing the engine.
	FENStartPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0"
	FENKiwiPete      = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

	// Constants mapping each board coordinate to its square
	A1, B1, C1, D1, E1, F1, G1, H1 = 0, 1, 2, 3, 4, 5, 6, 7
	A2, B2, C2, D2, E2, F2, G2, H2 = 8, 9, 10, 11, 12, 13, 14, 15
	A3, B3, C3, D3, E3, F3, G3, H3 = 16, 17, 18, 19, 20, 21, 22, 23
	A4, B4, C4, D4, E4, F4, G4, H4 = 24, 25, 26, 27, 28, 29, 30, 31
	A5, B5, C5, D5, E5, F5, G5, H5 = 32, 33, 34, 35, 36, 37, 38, 39
	A6, B6, C6, D6, E6, F6, G6, H6 = 40, 41, 42, 43, 44, 45, 46, 47
	A7, B7, C7, D7, E7, F7, G7, H7 = 48, 49, 50, 51, 52, 53, 54, 55
	A8, B8, C8, D8, E8, F8, G8, H8 = 56, 57, 58, 59, 60, 61, 62, 63

	// A constant representing no square
	NoSq = 64

	// Constant representing north and south deltas on the board
	NorthDelta = 8
	SouthDelta = -8

	// A constant representing the maximum game ply,
	// used to initalize the array for holding repetition
	// detection history.
	MaxGamePly = 1024
)

// A global array to store position histories. As this array has
// a quite large memory footprint, it's implemented as a single global
// variable, rather than being part of a Position instance.
var PositionHistories [MaxGamePly]uint64

// A global variable to index into the position histories.
var HistoryPly uint16

// A 64 element array where each entry, when bitwise ANDed with the
// castling rights, destorys the correct bit in the castling rights
// if a move to or from that square would take away castling rights.
var Spoilers [64]uint8 = [64]uint8{
	0xb, 0xf, 0xf, 0xf, 0x3, 0xf, 0xf, 0x7,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xe, 0xf, 0xf, 0xf, 0xc, 0xf, 0xf, 0xd,
}

// An array mapping a castling kings destination square,
// to the origin and destination square of the appropriate
// rook to move.
var CastlingRookSq map[uint8][2]uint8 = map[uint8][2]uint8{
	G1: {H1, F1},
	C1: {A1, D1},
	G8: {H8, F8},
	C8: {A8, D8},
}

// A constant mapping piece characters to Piece objects.
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

// A constant mapping piece types to their respective characters.
var PieceTypeToChar map[uint8]rune = map[uint8]rune{
	Pawn:   'i',
	Knight: 'n',
	Bishop: 'b',
	Rook:   'r',
	Queen:  'q',
	King:   'k',
	NoType: '.',
}

// A struct representing a piece
type Piece struct {
	Type  uint8
	Color uint8
}

// A struct representing position state that is irreversible (cannot be undone in
// UnmakeMove). A position state object is used each time a move is made, and then popped
// off of a stack once a move needs to be unmade.
type State struct {
	CastlingRights uint8
	EPSq           uint8
	Rule50         uint8

	Captured Piece
	Moved    Piece
}

// A struct reprenting Blunder's core internal position representation, which consists
// of 12 bitboards for each piece type, 2 bitboards for each color, 64-square
// mailbox representation of the board for easy accsess to square-centric information,
// and several other state keeping fields (enpassant square, side to move, etc.)
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

	// The zobrist hash of the position
	Hash uint64

	SideToMove uint8
	EPSq       uint8

	Ply    uint16
	Rule50 uint8

	prevStates [100]State
	StatePly   uint8
}

func (pos *Position) MakeMove(move Move) bool {
	// Get the data we need from the given move
	from := move.FromSq()
	to := move.ToSq()
	moveType := move.MoveType()
	flag := move.Flag()

	// Create a new State object to save the state of the irreversible aspects
	// of the position before making the current move.
	state := State{
		CastlingRights: pos.CastlingRights,
		EPSq:           pos.EPSq,
		Rule50:         pos.Rule50,
		Captured:       pos.Squares[to],
		Moved:          pos.Squares[from],
	}

	// Increment the game ply and the fifty-move rule counter
	pos.Ply++
	pos.Rule50++

	// Clear the en passant square and en passant zobrist number
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.EPSq = NoSq

	// Clear the moving piece from its origin square
	pos.clearPiece(from)

	if moveType == Quiet {
		// if the move is quiet, simple put the piece at the destination square.
		pos.putPiece(state.Moved.Type, state.Moved.Color, to)
	} else if moveType == Attack {
		if flag == AttackEP {
			// If it's an attack en passant, get the actually capture square
			// of the pawn being captured, remove it, and put the moving pawn
			// on the destination square...
			capSq := uint8(int8(to) - pawnPush(pos.SideToMove))
			state.Captured = pos.Squares[capSq]

			pos.clearPiece(capSq)
			pos.putPiece(Pawn, pos.SideToMove, to)

		} else {
			// Otherwise if the move is a normal attack, remove the captured piece
			// from the position, and put the moving piece at its destination square...
			pos.clearPiece(to)
			pos.putPiece(state.Moved.Type, state.Moved.Color, to)
		}

		// and reset the fifty-move rule counter.
		pos.Rule50 = 0
	} else if moveType == Castle {
		// If the move is a castle, move the king to the appropriate square...
		pos.putPiece(state.Moved.Type, state.Moved.Color, to)

		// And move the correct rook.
		rookFrom, rookTo := CastlingRookSq[to][0], CastlingRookSq[to][1]
		pos.clearPiece(rookFrom)
		pos.putPiece(Rook, pos.SideToMove, rookTo)

	}

	if state.Moved.Type == Pawn {
		// If a pawn is moving, do some extra work.

		// Reset the fifty-move rule counter.
		pos.Rule50 = 0

		if moveType == Promotion {
			// If a pawn is promoting, check if it's capturing a piece,
			// remove the captured piece if needed, and then put the
			// correct promotion piece type on the to square indicated
			// by the move flag value.
			if state.Captured.Type != NoType {
				pos.clearPiece(to)
			}
			pos.putPiece(uint8(flag+1), pos.SideToMove, to)
		}

		if abs16(int16(from)-int16(to)) == 16 {
			// If the move is a double pawn push, and there is no enemy pawn that's in
			// a position to capture en passant on the next turn, don't set the position's
			// en passant square.

			pos.EPSq = uint8(int8(to) - pawnPush(pos.SideToMove))
			if PawnAttacks[pos.SideToMove][pos.EPSq]&pos.PieceBB[pos.SideToMove^1][Pawn] == 0 {
				pos.EPSq = NoSq
			}
		}

	}

	// Remove the current castling rights.
	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)

	// Update the castling rights and the zobrist hash with the new castling rights.
	pos.CastlingRights = pos.CastlingRights & Spoilers[from] & Spoilers[to]
	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)

	// Update the zobrist hash if the en passant square was set
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)

	// Save the State object and increment the stack counter
	// to point to the next empty slot in the position state history.
	pos.prevStates[pos.StatePly] = state
	pos.StatePly++

	// Flip the side to move and update the zobrist hash
	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	// Save the current zobrist in the position history array.
	HistoryPly++
	PositionHistories[HistoryPly] = pos.Hash

	// Test if the move was legal or not, and let the caller know.
	return !sqIsAttacked(pos, pos.SideToMove^1, pos.PieceBB[pos.SideToMove^1][King].Msb())
}

func (pos *Position) UnmakeMove(move Move) {
	// Get the State object for this move
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

	// Flip the side to move and update the zobrist hash
	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	// Decrement the game ply
	pos.Ply--

	// Get the data we need from the given move
	from := move.FromSq()
	to := move.ToSq()
	moveType := move.MoveType()
	flag := move.Flag()

	// Put the moving piece back on it's orgin square
	pos.putPiece(state.Moved.Type, state.Moved.Color, from)

	if moveType == Quiet {
		// if the move is quiet, remove the piece from its destination square.
		pos.clearPiece(to)
	} else if moveType == Attack {
		if flag == AttackEP {
			// If it was an attack en passant, put the pawn back that
			// was captured, and clear the moving pawn from the destination
			// square.
			capSq := uint8(int8(to) - pawnPush(pos.SideToMove))
			pos.clearPiece(to)
			pos.putPiece(Pawn, state.Captured.Color, capSq)
		} else {
			// Otherwise If the move was a normal attack, put the captured piece
			// back on the destination square, and remove the attacking piece.
			pos.clearPiece(to)
			pos.putPiece(state.Captured.Type, state.Captured.Color, to)
		}
	} else if moveType == Castle {
		// If the move was a castle, clear the king from the destination square...
		pos.clearPiece(to)

		// and move the castled rook back to the right square.
		rookFrom, rookTo := CastlingRookSq[to][0], CastlingRookSq[to][1]
		pos.clearPiece(rookTo)
		pos.putPiece(Rook, pos.SideToMove, rookFrom)
	}

	if state.Moved.Type == Pawn {
		// If a pawn was moving, do some extra work.

		if moveType == Promotion {
			// If the pawn was promoted, remove the promoted piece, and if
			// the promotion was a capture, put the captured piece back on
			// the destination square.
			pos.clearPiece(to)
			if state.Captured.Type != NoType {
				pos.putPiece(state.Captured.Type, state.Captured.Color, to)
			}
		}
	}
}

// Make a "null"-move for null-move pruning:
// https://www.chessprogramming.org/Null_Move_Pruning
//
func (pos *Position) MakeNullMove() {
	state := State{
		CastlingRights: pos.CastlingRights,
		EPSq:           pos.EPSq,
		Rule50:         pos.Rule50,
	}

	// Save the State object and increment the stack counter
	// to point to the next empty slot in the position state history.
	pos.prevStates[pos.StatePly] = state
	pos.StatePly++

	// Clear the en passant square and en passant zobrist number
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.EPSq = NoSq

	// Set the fifty move rule counter to 0, since we're
	// making a null-move.
	pos.Rule50 = 0

	// Increment the game ply.
	pos.Ply++

	// Flip the side to move and update the zobrist hash
	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	// Save the current zobrist in the position history array.
	HistoryPly++
	PositionHistories[HistoryPly] = pos.Hash

}

func (pos *Position) UnmakeNullMove() {
	// Get the State object for the null move
	pos.StatePly--
	state := pos.prevStates[pos.StatePly]

	// Restore the irreversible aspects of the position using the State object.
	pos.CastlingRights = state.CastlingRights
	pos.EPSq = state.EPSq
	pos.Rule50 = state.Rule50

	// Decrement the game ply.
	pos.Ply--

	// Update the zobrist hash with the restored en passant square.
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)

	// Flip the side to move and update the zobrist hash
	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber(pos.SideToMove)

	// Remove the current positions from the position history
	HistoryPly--
}

// Put the piece given on the given square
func (pos *Position) putPiece(pieceType, pieceColor, to uint8) {
	pos.PieceBB[pieceColor][pieceType].SetBit(to)
	pos.SideBB[pieceColor].SetBit(to)
	pos.Squares[to].Type = pieceType
	pos.Squares[to].Color = pieceColor
	pos.Hash ^= Zobrist.PieceNumber(pieceType, pieceColor, to)
}

// Clear the piece given from the given square.
func (pos *Position) clearPiece(from uint8) {
	piece := &pos.Squares[from]
	pos.PieceBB[piece.Color][piece.Type].ClearBit(from)
	pos.SideBB[piece.Color].ClearBit(from)

	pos.Hash ^= Zobrist.PieceNumber(piece.Type, piece.Color, from)
	piece.Type = NoType
	piece.Color = NoColor
}

// Load in a FEN string and use it to setup the position.
func (pos *Position) LoadFEN(fen string) {
	// Reset the internal fields of the position
	pos.PieceBB = [2][6]Bitboard{}
	pos.SideBB = [2]Bitboard{}
	pos.Squares = [64]Piece{}
	pos.CastlingRights = 0

	for square := range pos.Squares {
		pos.Squares[square] = Piece{Type: NoType, Color: NoColor}
	}

	// Load in each field of the FEN string.
	fields := strings.Fields(fen)
	pieces := fields[0]
	color := fields[1]
	castling := fields[2]
	ep := fields[3]
	halfMove := fields[4]
	fullMove := fields[5]

	// Loop over each square of the board, rank by rank, from left to right,
	// loading in pieces at squares described by the FEN string.
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

	// Set the side to move for the position.
	pos.SideToMove = Black
	if color == "w" {
		pos.SideToMove = White
	}

	// Set the en passant square for the position.
	pos.EPSq = NoSq
	if ep != "-" {
		pos.EPSq = coordinateToPos(ep)
		if (PawnAttacks[pos.SideToMove^1][pos.EPSq] & pos.PieceBB[pos.SideToMove][Pawn]) == 0 {
			pos.EPSq = NoSq
		}
	}

	// Set the half move counter and game ply for the position.
	halfMoveCounter, _ := strconv.Atoi(halfMove)
	pos.Rule50 = uint8(halfMoveCounter)

	gamePly, _ := strconv.Atoi(fullMove)
	gamePly *= 2
	if pos.SideToMove == Black {
		gamePly--
	}
	pos.Ply = uint16(gamePly)

	// Set the castling rights, for the position.
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

	// Generate the zobrist hash for the position...
	pos.Hash = 0
	pos.Hash = Zobrist.GenHash(pos)

	// and add the hash as the first entry in the position history.
	HistoryPly = 0
	PositionHistories[HistoryPly] = pos.Hash
}

// Return a string representation of the board.
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

	boardAsString += "   "
	for fileNo := 0; fileNo < 8; fileNo++ {
		boardAsString += "--"
	}

	boardAsString += "\n    "
	for _, file := range "abcdefgh" {
		boardAsString += fmt.Sprintf("%c ", file)
	}

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
	if pos.EPSq == NoSq {
		boardAsString += "none"
	} else {
		boardAsString += posToCoordinate(pos.EPSq)
	}

	boardAsString += fmt.Sprintf("\nzobrist hash: 0x%x", pos.Hash)
	boardAsString += fmt.Sprintf("\nrule 50: %d\n", pos.Rule50)
	boardAsString += fmt.Sprintf("game ply: %d\n", pos.Ply)
	return boardAsString
}

// Determine if the side to move is in check.
func (pos *Position) InCheck() bool {
	return sqIsAttacked(
		pos,
		pos.SideToMove,
		pos.PieceBB[pos.SideToMove][King].Msb())
}

// Determine if the current position should be considered an endgame
// position for the current side to move.
func (pos *Position) IsEndgameForSide() bool {
	pawnMaterial := int16(pos.PieceBB[pos.SideToMove][Pawn].CountBits()) * 100
	knightMaterial := int16(pos.PieceBB[pos.SideToMove][Knight].CountBits()) * 320
	bishopMaterial := int16(pos.PieceBB[pos.SideToMove][Bishop].CountBits()) * 330
	rookMaterial := int16(pos.PieceBB[pos.SideToMove][Rook].CountBits()) * 500
	queenMaterial := int16(pos.PieceBB[pos.SideToMove][Queen].CountBits()) * 950

	return (pawnMaterial + knightMaterial + bishopMaterial + rookMaterial + queenMaterial) < 1300
}

// Given a color, return the delta for a single pawn push for that
// color.
func pawnPush(color uint8) int8 {
	if color == White {
		return NorthDelta
	}
	return SouthDelta
}
