package engine

import (
	"fmt"
	"regexp"
	"strings"
)

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

const (
	CastlingMovePattern    = "(?P<CastlingMove>(O-O-O)|(O-O))"
	PromotionMovePattern   = "(?P<PromotionMove>[a-h]?x?[a-h][1-8]=[NBRQ])"
	AttackPawnMovePattern  = "(?P<AttackPawnMove>[a-h]x[a-h][1-8])"
	AttackPieceMovePattern = "(?P<AttackPieceMove>[NBRQK][a-h]?[1-8]?x[a-h][1-8])"
	QuietPawnMovePattern   = "(?P<QuietPawnMove>[a-h][1-8])"
	QuietPieceMovePattern  = "(?P<QuietPieceMove>[NBRQK][a-h]?[1-8]?[a-h][1-8])"
)

var MoveRe = regexp.MustCompile(strings.Join(
	[]string{
		CastlingMovePattern, PromotionMovePattern,
		AttackPawnMovePattern, AttackPieceMovePattern,
		QuietPawnMovePattern, QuietPieceMovePattern,
	},
	"|",
))

var CharToPieceType = map[rune]uint8{
	'N': Knight,
	'B': Bishop,
	'R': Rook,
	'Q': Queen,
	'K': King,
}

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

func ConvertSANToLAN(pos *Position, moveStr string) (uint32, error) {
	matches := MoveRe.FindStringSubmatch(moveStr)
	movePattern := ""

	for i, name := range MoveRe.SubexpNames() {
		if i != 0 && name != "" && matches[i] != "" {
			movePattern = name
		}
	}

	if movePattern == "CastlingMove" {
		if moveStr == "O-O" && pos.SideToMove == White {
			return newMove(E1, G1, Castle, NoFlag), nil
		} else if moveStr == "O-O" && pos.SideToMove == Black {
			return newMove(E8, G8, Castle, NoFlag), nil
		} else if moveStr == "O-O-O" && pos.SideToMove == White {
			return newMove(E1, C1, Castle, NoFlag), nil
		} else if moveStr == "O-O-O" && pos.SideToMove == Black {
			return newMove(E8, C8, Castle, NoFlag), nil
		}
	}

	if movePattern == "QuietPawnMove" {
		delta := -8
		if pos.SideToMove == Black {
			delta = 8
		}

		to := coordToSq(moveStr)

		if pos.GetPieceType(uint8(int(to)+delta)) == Pawn {
			return newMove(uint8(int(to)+delta), to, Quiet, NoFlag), nil
		} else if pos.GetPieceType(uint8(int(to)+delta+delta)) == Pawn {
			return newMove(uint8(int(to)+delta+delta), to, Quiet, NoFlag), nil
		} else {
			return NullMove, fmt.Errorf("could not find pawn moving to sq %s for move %s", sqToCoord(to), moveStr)
		}
	}

	moves := genAllMoves(pos)
	possibleMatchingMoves := []uint32{}

	if movePattern == "PromotionMove" {
		var file, to, promoType uint8
		file = NoSq

		if strings.Contains(moveStr, "x") {
			file, to = moveStr[0]-'a', coordToSq(moveStr[2:4])
			promoType = CharToPieceType[rune(moveStr[5])]
		} else {
			to = coordToSq(moveStr[0:2])
			promoType = CharToPieceType[rune(moveStr[3])]
		}

		for i := 0; i < int(moves.Count); i++ {
			move := moves.Moves[i]
			if moveType(move) == Promotion &&
				toSq(move) == to &&
				flag(move)+1 == promoType &&
				(file == NoSq || fileOf(fromSq(move)) == file) {
				possibleMatchingMoves = append(possibleMatchingMoves, move)
			}
		}
	} else if movePattern == "AttackPawnMove" {
		file, to := moveStr[0]-'a', coordToSq(moveStr[2:4])
		for i := 0; i < int(moves.Count); i++ {
			move := moves.Moves[i]
			if moveType(move) == Attack &&
				pos.GetPieceType(fromSq(move)) == Pawn &&
				toSq(move) == to &&
				fileOf(fromSq(move)) == file {
				possibleMatchingMoves = append(possibleMatchingMoves, move)
			}
		}
	} else if movePattern == "AttackPieceMove" {
		var to, from, file, rank, pieceType uint8
		pieceType = CharToPieceType[rune(moveStr[0])]
		from, file, rank = NoSq, NoSq, NoSq

		if len(moveStr) == 5 {
			to = coordToSq(moveStr[3:5])
			fileOrRank := moveStr[1]
			if fileOrRank >= '1' && fileOrRank <= '8' {
				rank = fileOrRank - '0' - 1
			} else {
				file = fileOrRank - 'a'
			}
		} else if len(moveStr) == 6 {
			from, to = coordToSq(moveStr[1:3]), coordToSq(moveStr[4:6])
		} else {
			to = coordToSq(moveStr[2:4])
		}

		for i := 0; i < int(moves.Count); i++ {
			move := moves.Moves[i]
			moveFromSq := fromSq(move)

			if moveType(move) == Attack &&
				pos.GetPieceType(moveFromSq) == pieceType &&
				toSq(move) == to &&
				(from == NoSq || moveFromSq == from) &&
				(file == NoSq || fileOf(moveFromSq) == file) &&
				(rank == NoSq || rankOf(moveFromSq) == rank) {
				possibleMatchingMoves = append(possibleMatchingMoves, move)
			}
		}
	} else if movePattern == "QuietPieceMove" {
		var to, from, file, rank, pieceType uint8
		pieceType = CharToPieceType[rune(moveStr[0])]
		from, file, rank = NoSq, NoSq, NoSq

		if len(moveStr) == 4 {
			to = coordToSq(moveStr[2:4])
			fileOrRank := moveStr[1]
			if fileOrRank >= '1' && fileOrRank <= '8' {
				rank = fileOrRank - '0' - 1
			} else {
				file = fileOrRank - 'a'
			}
		} else if len(moveStr) == 5 {
			from, to = coordToSq(moveStr[1:3]), coordToSq(moveStr[3:5])
		} else {
			to = coordToSq(moveStr[1:3])
		}

		for i := 0; i < int(moves.Count); i++ {
			move := moves.Moves[i]
			moveFromSq := fromSq(move)

			if moveType(move) == Quiet &&
				pos.GetPieceType(moveFromSq) == pieceType &&
				toSq(move) == to &&
				(from == NoSq || moveFromSq == from) &&
				(file == NoSq || fileOf(moveFromSq) == file) &&
				(rank == NoSq || rankOf(moveFromSq) == rank) {
				possibleMatchingMoves = append(possibleMatchingMoves, move)
			}
		}
	} else {
		panic(fmt.Errorf("move %s not recognized", moveStr))
	}

	if len(possibleMatchingMoves) == 1 {
		return possibleMatchingMoves[0], nil
	} else {
		for _, move := range possibleMatchingMoves {
			if pos.DoMove(move) {
				pos.UndoMove(move)
				return move, nil
			}
			pos.UndoMove(move)
		}
	}

	return NullMove, fmt.Errorf("failed to convert move %s", moveStr)
}
