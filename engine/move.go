package engine

import "fmt"

// A file containg utility methods and type definitions
// for moves in Blunder.

const (
	Quiet uint8 = iota
	DoublePawnPush
	Attack
	AttackEP
	KnightPromotion
	BishopPromotion
	RookPromotion
	QueenPromotion
	CastleWKS
	CastleWQS
	CastleBKS
	CastleBQS

	MaxMoves = 50
)

type Move uint16
type Moves = []Move

// A helper function to create moves. Each move generated
// by our move generator is encoded in 16-bits, where the
// first six bits are the from square, the second 6, are the
// to square, and the last four are the move type (see above).
func MakeMove(from, to, moveType uint8) Move {
	return Move(uint16(from)<<10 | uint16(to)<<4 | uint16(moveType))
}

// Get the type of the move.
func (move Move) MoveType() uint8 {
	return uint8(move & 0xf)
}

// Get the from square of the move.
func (move Move) FromSq() uint8 {
	return uint8((move & 0xFC00) >> 10)
}

// Get the to square of the move.
func (move Move) ToSq() uint8 {
	return uint8((move & 0x3F0) >> 4)
}

// A helper function to extract the info from a move represented
// as 32-bits, and display it.
func (move Move) String() string {
	from, to, movType := move.FromSq(), move.ToSq(), move.MoveType()
	promotionType, seperator := "", "-"
	switch movType {
	case Attack:
		fallthrough
	case AttackEP:
		seperator = "x"
	case KnightPromotion:
		promotionType = "n"
	case BishopPromotion:
		promotionType = "b"
	case RookPromotion:
		promotionType = "r"
	case QueenPromotion:
		promotionType = "q"
	}
	return fmt.Sprintf(
		"%v%v%v%v",
		PosToCoordinate(from),
		seperator,
		PosToCoordinate(to),
		promotionType,
	)
}

func MoveFromCoord(board *Board, move string, useChess960Castling bool) Move {
	fromPos := CoordinateToPos(move[0:2])
	toPos := CoordinateToPos(move[2:4])
	movePieceType := board.Squares[fromPos].Type
	var moveType uint8

	moveLen := len(move)
	if moveLen == 5 {
		if move[moveLen-1] == 'n' {
			moveType = KnightPromotion
		} else if move[moveLen-1] == 'b' {
			moveType = BishopPromotion
		} else if move[moveLen-1] == 'r' {
			moveType = RookPromotion
		} else if move[moveLen-1] == 'q' {
			moveType = QueenPromotion
		}
	} else if move == "e1g1" && movePieceType == King && !useChess960Castling {
		moveType = CastleWKS
	} else if move == "e1c1" && movePieceType == King && !useChess960Castling {
		moveType = CastleWQS
	} else if move == "e8g8" && movePieceType == King && !useChess960Castling {
		moveType = CastleBKS
	} else if move == "e8c8" && movePieceType == King && !useChess960Castling {
		moveType = CastleBQS
	} else if move == "e1h1" && movePieceType == King && useChess960Castling {
		moveType = CastleWKS
	} else if move == "e1a1" && movePieceType == King && useChess960Castling {
		moveType = CastleWQS
	} else if move == "e8h8" && movePieceType == King && useChess960Castling {
		moveType = CastleBKS
	} else if move == "e8a8" && movePieceType == King && useChess960Castling {
		moveType = CastleBQS
	} else if toPos == board.EPSquare && movePieceType == Pawn {
		moveType = AttackEP
	} else {
		capturePieceType := board.Squares[toPos].Type
		if capturePieceType == NoType {
			if movePieceType == Pawn && abs(int8(fromPos)-int8(toPos)) == 16 {
				moveType = DoublePawnPush
			} else {
				moveType = Quiet
			}
		} else {
			moveType = Attack
		}
	}
	return MakeMove(fromPos, toPos, moveType)
}
