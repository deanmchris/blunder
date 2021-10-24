package engine

// see.go is a simple implementation of a static exchange evaluator.

var PieceValues [7]int16 = [7]int16{
	100,
	300,
	300,
	500,
	900,
	5000,
	0,
}

// Peform a static exchange evaluation on target square of the move given,
// and return a score of the move from the perspective of the side to move.
func (pos *Position) See(move Move) int16 {
	if move.MoveType() == Promotion {
		// Let promotions be evaluated dynamically.
		return 0
	}

	fromSq := move.FromSq()
	toSq := move.ToSq()
	attackerSq := fromSq

	sideToMove := pos.SideToMove
	occupiedBB := pos.SideBB[sideToMove] | pos.SideBB[sideToMove^1]

	scores := [32]int16{}
	attackers := Bitboard(0)
	seen := Bitboard(0)

	scores[0] = PieceValues[pos.Squares[toSq].Type]
	attackerType := pos.Squares[attackerSq].Type
	attackerValue := PieceValues[pos.Squares[attackerSq].Type]

	if move.Flag() == AttackEP {
		capSq := uint8(int8(toSq) - pawnPush(sideToMove))
		occupiedBB &= ^SquareBB[capSq]
		scores[0] = PieceValues[Pawn]
	}

	seen |= SquareBB[attackerSq]
	sideToMove ^= 1
	attackers = attackersOfSquare(pos, sideToMove, toSq, occupiedBB, seen)

	if attackers == 0 {
		return scores[0]
	}

	scoresIndex := 1
	for {
		occupiedBB &= ^SquareBB[attackerSq]
		scores[scoresIndex] = attackerValue - scores[scoresIndex-1]

		attackerSq = getLeastValuableAttackerSq(pos, attackers)
		attackerType = pos.Squares[attackerSq].Type
		attackerValue = PieceValues[pos.Squares[attackerSq].Type]

		if attackerValue == PieceValues[Pawn] && isPromoting(sideToMove, toSq) {
			// Pawn promotions aren't handled in the capture sequence either,
			// so let them be handled dynamically.
			return 0
		}

		seen |= SquareBB[attackerSq]
		sideToMove ^= 1
		attackers = attackersOfSquare(pos, sideToMove, toSq, occupiedBB, seen)

		if attackers == 0 {
			break
		}

		if attackerType == King {
			scores[scoresIndex] = 0
			scoresIndex--
			break
		}

		scoresIndex++
	}

	for ; scoresIndex > 0; scoresIndex-- {
		scores[scoresIndex-1] = min16(-scores[scoresIndex], scores[scoresIndex-1])
	}

	return scores[0]
}

func attackersOfSquare(pos *Position, enemyColor, sq uint8, occupiedBB, seen Bitboard) (attackers Bitboard) {
	enemyBishops := pos.PieceBB[enemyColor][Bishop]
	enemyRooks := pos.PieceBB[enemyColor][Rook]
	enemyQueens := pos.PieceBB[enemyColor][Queen]
	enemyKnights := pos.PieceBB[enemyColor][Knight]
	enemyKing := pos.PieceBB[enemyColor][King]
	enemyPawns := pos.PieceBB[enemyColor][Pawn]

	intercardinalRays := genBishopMoves(sq, occupiedBB)
	cardinalRaysRays := genRookMoves(sq, occupiedBB)

	attackers |= intercardinalRays & (enemyBishops | enemyQueens)
	attackers |= cardinalRaysRays & (enemyRooks | enemyQueens)
	attackers |= KnightMoves[sq] & enemyKnights
	attackers |= KingMoves[sq] & enemyKing
	attackers |= PawnAttacks[enemyColor^1][sq] & enemyPawns
	return attackers & ^seen
}

// Get the last valuable attacker from a bitboard of attackers.
func getLeastValuableAttackerSq(pos *Position, attackers Bitboard) uint8 {
	lowestValue := Inf
	attackerSq := uint8(0)

	for attackers != 0 {
		sq := attackers.PopBit()
		value := PieceValues[pos.Squares[sq].Type]

		if value < lowestValue {
			lowestValue = value
			attackerSq = sq
		}
	}

	return attackerSq
}

// Get the minimum between two numbers.
func min16(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}
