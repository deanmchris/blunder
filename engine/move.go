package engine

import "fmt"

// A file containg utility methods and type definitions
// for moves in Blunder.

const (
	Quiet = iota
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

type Move = uint32
type Moves = []Move

// A helper function to create moves. Each move generated
// by our move generator is encoded in 16-bits, where the
// first six bits are the from square, the second 6, are the
// to square, and the last four are the move type (see above).
func MakeMove(from, to, moveType int) Move {
	return Move(from<<10 | to<<4 | moveType)
}

// Get the type of the move.
func MoveType(move Move) int {
	return int(move & 0xf)
}

// Get the from square of the move.
func FromSq(move Move) int {
	return int((move & 0xFC00) >> 10)
}

// Get the to square of the move.
func ToSq(move Move) int {
	return int((move & 0x3F0) >> 4)
}

// A helper function to extract the info from a move represented
// as 32-bits, and display it.
func MoveStr(move Move) string {
	from, to, movType := FromSq(move), ToSq(move), MoveType(move)
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
		PosToCoordinate(int(from)),
		seperator,
		PosToCoordinate(int(to)),
		promotionType,
	)
}
