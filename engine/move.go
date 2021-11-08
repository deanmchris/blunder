package engine

import "fmt"

// move.go constaints the implementation of a move datatype.

const (
	// Constants represeting the four possible move types.
	Quiet     uint8 = 0
	Attack    uint8 = 1
	Castle    uint8 = 2
	Promotion uint8 = 3

	// Constants representing move flags indicating what kind of promotion
	// is occuring.
	KnightPromotion uint8 = 0
	BishopPromotion uint8 = 1
	RookPromotion   uint8 = 2
	QueenPromotion  uint8 = 3

	// A constant representing a move flag indicating an attack is an en passant
	// attack.
	AttackEP uint8 = 1

	// A constant representing a null flag
	NoFlag uint8 = 0
)

type Move uint32

// Create a new move. The first 6 bits are the from square, the next 6 bits are the to square,
// the next two represent the move type, the next two are reserved for any speical flags needed
// to give full information concering the move, and the last 16-bits are used for scoring a move
// for move-ordering in the search phase.
func NewMove(from, to, moveType, flag uint8) Move {
	return Move(uint32(from)<<26 | uint32(to)<<20 | uint32(moveType)<<18 | uint32(flag)<<16)
}

// Get the from square of the move.
func (move Move) FromSq() uint8 {
	return uint8((move & 0xfc000000) >> 26)
}

// Get the to square of the move.
func (move Move) ToSq() uint8 {
	return uint8((move & 0x3f00000) >> 20)
}

// Get the type of the move.
func (move Move) MoveType() uint8 {
	return uint8((move & 0xc0000) >> 18)
}

// Get the flag of the move.
func (move Move) Flag() uint8 {
	return uint8((move & 0x30000) >> 16)
}

// Get the score of a move.
func (move Move) Score() uint16 {
	return uint16(move & 0xffff)
}

// Add a score to the move for move ordering.
func (move *Move) AddScore(score uint16) {
	(*move) &= 0xffff0000
	(*move) |= Move(score)
}

// Test if two moves are equal.
func (move Move) Equal(m Move) bool {
	return (move & 0xffff0000) == (m & 0xffff0000)
}

// A helper function to extract the info from a move represented
// as 32-bits, and display it.
func (move Move) String() string {
	from, to, moveType, flag := move.FromSq(), move.ToSq(), move.MoveType(), move.Flag()

	promotionType := ""
	if moveType == Promotion {
		switch flag {
		case KnightPromotion:
			promotionType = "n"
		case BishopPromotion:
			promotionType = "b"
		case RookPromotion:
			promotionType = "r"
		case QueenPromotion:
			promotionType = "q"
		}
	}
	return fmt.Sprintf("%v%v%v", posToCoordinate(from), posToCoordinate(to), promotionType)
}

// Convert a move in UCI format into a Move
func MoveFromCoord(pos *Position, move string) Move {
	from := CoordinateToPos(move[0:2])
	to := CoordinateToPos(move[2:4])
	moved := pos.Squares[from].Type

	var moveType uint8
	flag := NoFlag

	moveLen := len(move)
	if moveLen == 5 {
		moveType = Promotion
		if move[moveLen-1] == 'n' {
			flag = KnightPromotion
		} else if move[moveLen-1] == 'b' {
			flag = BishopPromotion
		} else if move[moveLen-1] == 'r' {
			flag = RookPromotion
		} else if move[moveLen-1] == 'q' {
			flag = QueenPromotion
		}
	} else if move == "e1g1" && moved == King {
		moveType = Castle
	} else if move == "e1c1" && moved == King {
		moveType = Castle
	} else if move == "e8g8" && moved == King {
		moveType = Castle
	} else if move == "e8c8" && moved == King {
		moveType = Castle
	} else if to == pos.EPSq && moved == Pawn {
		moveType = Attack
		flag = AttackEP
	} else {
		captured := pos.Squares[to]
		if captured.Type == NoType {
			moveType = Quiet
		} else {
			moveType = Attack
		}
	}
	return NewMove(from, to, moveType, flag)
}
