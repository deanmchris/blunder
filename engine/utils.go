package engine

import (
	"strings"
	"unicode"
)

// utils.go contains various utility functions used throughout the engine.

var CharToPieceType map[rune]uint8 = map[rune]uint8{
	'N': Knight,
	'B': Bishop,
	'R': Rook,
	'Q': Queen,
	'K': King,
}

// Convert a string board coordinate to its position
// number.
func CoordinateToPos(coordinate string) uint8 {
	file := coordinate[0] - 'a'
	rank := int(coordinate[1]-'0') - 1
	return uint8(rank*8 + int(file))
}

// Convert a position number to a string board coordinate.
func posToCoordinate(pos uint8) string {
	file := FileOf(pos)
	rank := RankOf(pos)
	return string(rune('a'+file)) + string(rune('0'+rank+1))
}

// Given a board square, return it's file.
func FileOf(sq uint8) uint8 {
	return sq % 8
}

// Given a board square, return it's rank.
func RankOf(sq uint8) uint8 {
	return sq / 8
}

// Get the absolute value of a number.
func abs16(n int16) int16 {
	if n < 0 {
		return -n
	}
	return n
}

// Get the maximum between two signed 8-bit numbers.
func max8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

// Determine if a square is dark.
func sqIsDark(sq uint8) bool {
	fileNo := FileOf(sq)
	rankNo := RankOf(sq)
	return ((fileNo + rankNo) % 2) == 0
}

// Pad a string into the center
func padToCenter(s string, fill string, w int) string {
	spaceLeft := w - len(s)
	extraFill := ""
	if spaceLeft%2 != 0 {
		extraFill = fill
	}
	return strings.Repeat(fill, spaceLeft/2) + extraFill + s + strings.Repeat(fill, spaceLeft/2)
}

// Convert a move in short algebraic notation, to the long algebraic notation used
// by the UCI protocol.
func ConvertSANToLAN(pos *Position, moveStr string) Move {
	if moveStr == "O-O" && pos.SideToMove == White {
		return NewMove(E1, G1, Castle, NoFlag)
	} else if moveStr == "O-O" && pos.SideToMove == Black {
		return NewMove(E8, G8, Castle, NoFlag)
	} else if moveStr == "O-O-O" && pos.SideToMove == White {
		return NewMove(E1, C1, Castle, NoFlag)
	} else if moveStr == "O-O-O" && pos.SideToMove == Black {
		return NewMove(E8, C8, Castle, NoFlag)
	}

	coords := ""
	pieceType := Pawn

	for _, char := range moveStr {
		switch char {
		case 'N', 'B', 'R', 'Q', 'K':
			pieceType = CharToPieceType[char]
		case '1', '2', '3', '4', '5', '6', '7', '8':
			coords += string(char)
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
			coords += string(char)
		}
	}

	moves := GenAllMoves(pos)
	matchingMove := NullMove

	for i := 0; i < int(moves.Count); i++ {
		move := moves.Moves[i]
		moved := pos.Squares[move.FromSq()].Type
		captured := pos.Squares[move.ToSq()].Type

		if len(coords) == 2 {
			if len(moveStr) == 4 && moveStr[2] == '=' {
				promotionType := pieceType - 1
				toSq := CoordinateToPos(coords[0:2])
				if move.ToSq() == toSq && move.MoveType() == Promotion && move.Flag() == promotionType {
					matchingMove = move
				}
			} else {
				toSq := CoordinateToPos(coords)
				if toSq == move.ToSq() && pieceType == moved {
					matchingMove = move
				}
			}
		} else if len(coords) == 3 {
			if len(moveStr) == 6 && moveStr[4] == '=' {
				promotionType := pieceType - 1
				toSq := CoordinateToPos(coords[1:])

				if captured != NoType &&
					move.MoveType() == Promotion &&
					move.Flag() == promotionType &&
					move.ToSq() == toSq {
					matchingMove = move
				}
			} else {

				toSq := CoordinateToPos(coords[1:])
				fileOrRank := coords[0]
				moveCoords := move.String()

				if unicode.IsLetter(rune(fileOrRank)) {
					if toSq == move.ToSq() && fileOrRank == moveCoords[0] && moved == pieceType {
						matchingMove = move
					}
				} else {
					if toSq == move.ToSq() && fileOrRank == moveCoords[1] && moved == pieceType {
						matchingMove = move
					}
				}
			}
		} else if len(coords) == 4 {
			toSq := CoordinateToPos(coords[0:2])
			fromSq := CoordinateToPos(coords[2:4])
			if toSq == move.ToSq() && fromSq == move.FromSq() {
				matchingMove = move
			}
		}

		if matchingMove != NullMove {
			if !pos.DoMove(matchingMove) {
				pos.UndoMove(matchingMove)
				continue
			}
			pos.UndoMove(matchingMove)
			break
		}
	}

	return matchingMove
}
