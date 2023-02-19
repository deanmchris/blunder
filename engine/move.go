package engine

import "fmt"

const (
	Quiet     uint8 = 0
	Attack    uint8 = 1
	Castle    uint8 = 2
	Promotion uint8 = 3

	KnightPromotion uint8 = 0
	BishopPromotion uint8 = 1
	RookPromotion   uint8 = 2
	QueenPromotion  uint8 = 3

	AttackEP uint8 = 1
	NoFlag   uint8 = 0

	NullMove uint32 = 0
)

func newMove(from, to, moveType, flag uint8) uint32 {
	return uint32(from)<<26 | uint32(to)<<20 | uint32(moveType)<<18 | uint32(flag)<<16
}

func fromSq(move uint32) uint8 {
	return uint8((move & 0xfc000000) >> 26)
}

func toSq(move uint32) uint8 {
	return uint8((move & 0x3f00000) >> 20)
}

func moveType(move uint32) uint8 {
	return uint8((move & 0xc0000) >> 18)
}

func flag(move uint32) uint8 {
	return uint8((move & 0x30000) >> 16)
}

func score(move uint32) uint16 {
	return uint16(move & 0xffff)
}

func addScore(move *uint32, score uint16) {
	*move &= 0xffff0000
	*move |= uint32(score)
}

func equals(m1, m2 uint32) bool {
	return (m1 & 0xffff0000) == (m2 & 0xffff0000)
}

func moveToStr(move uint32) string {
	from := fromSq(move)
	to := toSq(move)
	moveType := moveType(move)
	flag := flag(move)

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
	return fmt.Sprintf("%v%v%v", sqToCoord(from), sqToCoord(to), promotionType)
}

func moveFromCoord(pos *Position, move string) uint32 {
	from := coordToSq(move[0:2])
	to := coordToSq(move[2:4])
	movedType := pos.GetPieceType(from)

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
	} else if move == "e1g1" && movedType == King {
		moveType = Castle
	} else if move == "e1c1" && movedType == King {
		moveType = Castle
	} else if move == "e8g8" && movedType == King {
		moveType = Castle
	} else if move == "e8c8" && movedType == King {
		moveType = Castle
	} else if to == pos.EPSq && movedType == Pawn {
		moveType = Attack
		flag = AttackEP
	} else {
		capturedType := pos.GetPieceType(to)
		if capturedType == NoType {
			moveType = Quiet
		} else {
			moveType = Attack
		}
	}
	return newMove(from, to, moveType, flag)
}
