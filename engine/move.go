package engine

import "fmt"

// move.go constaints the implementation of a move datatype.

const (
	Quiet     uint8 = 0
	Attack    uint8 = 1
	Castle    uint8 = 2
	Promotion uint8 = 3

	KnightPromotionFlag uint8 = 0
	BishopPromotionFlag uint8 = 1
	RookPromotionFlag   uint8 = 2
	QueenPromotionFlag  uint8 = 3
	AttackEPFlag        uint8 = 1
	NoFlag              uint8 = 0
)

type Move uint32

// Create a new move. The first 6 bits are the from square, the next 6 bits are the to square,
// the next two represent the move type, the next two are reserved for any speical flags needed
// to give full information concering the move, and the last 16-bits are used for scoring a move
// for move-ordering in the search phase.
func NewMove(from, to, moveType, flag uint8) Move {
	return Move(uint32(from)<<26 | uint32(to)<<20 | uint32(moveType)<<18 | uint32(flag)<<16)
}

func (move Move) FromSq() uint8 {
	return uint8((move & 0xfc000000) >> 26)
}

func (move Move) ToSq() uint8 {
	return uint8((move & 0x3f00000) >> 20)
}

func (move Move) MoveType() uint8 {
	return uint8((move & 0xc0000) >> 18)
}

func (move Move) Flag() uint8 {
	return uint8((move & 0x30000) >> 16)
}

func (move Move) Score() uint16 {
	return uint16(move & 0xffff)
}

func (move *Move) AddScore(score uint16) {
	(*move) &= 0xffff0000
	(*move) |= Move(score)
}

func (move Move) Equal(m Move) bool {
	return (move & 0xffff0000) == (m & 0xffff0000)
}

func (move Move) String() string {
	from, to, moveType, flag := move.FromSq(), move.ToSq(), move.MoveType(), move.Flag()

	promotionType := ""
	if moveType == Promotion {
		switch flag {
		case KnightPromotionFlag:
			promotionType = "n"
		case BishopPromotionFlag:
			promotionType = "b"
		case RookPromotionFlag:
			promotionType = "r"
		case QueenPromotionFlag:
			promotionType = "q"
		}
	}
	return fmt.Sprintf("%v%v%v", posToCoordinate(from), posToCoordinate(to), promotionType)
}

func UCIMoveToInternalMove(pos *Position, move string) Move {
	from := CoordinateToPos(move[0:2])
	to := CoordinateToPos(move[2:4])
	moved := pos.Squares[from].Type

	var moveType uint8
	flag := NoFlag

	moveLen := len(move)
	if moveLen == 5 {
		moveType = Promotion
		if move[moveLen-1] == 'n' {
			flag = KnightPromotionFlag
		} else if move[moveLen-1] == 'b' {
			flag = BishopPromotionFlag
		} else if move[moveLen-1] == 'r' {
			flag = RookPromotionFlag
		} else if move[moveLen-1] == 'q' {
			flag = QueenPromotionFlag
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
		flag = AttackEPFlag
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
