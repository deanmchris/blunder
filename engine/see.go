package engine

import "fmt"

// see.go is a simple implementation of a static exchange evaluator.

var PieceValues [7]int16 = [7]int16{
	100,
	300,
	300,
	500,
	900,
	0,
	0,
}

// Peform a static exchange evaluation on target square of the move given,
// and return a score of the move from the perspective of the side to move.
// Credit to the Stockfish 7 team for the structure of the algorithm below:
//
// https://github.com/official-stockfish/Stockfish/blob/dd9cf305816c84c2acfa11cae09a31c4d77cc5a5/src/position.cpp
//
func (pos *Position) See(move Move) int16 {
	fromSq := move.FromSq()
	toSq := move.ToSq()
	scores := [32]int16{}
	sideToMove := pos.Squares[fromSq].Color

	scores[0] = PieceValues[pos.Squares[toSq].Type]
	occupiedBB := pos.SideBB[White] | pos.SideBB[Black]
	occupiedBB &= ^SquareBB[fromSq]

	if move.MoveType() == Castle {
		return 0
	}

	if move.Flag() == AttackEP {
		capSq := uint8(int8(toSq) - pawnPush(sideToMove))
		occupiedBB &= ^SquareBB[capSq]
		scores[0] = PieceValues[pos.Squares[capSq].Type]
	}

	attackers := allAttackers(pos, toSq, occupiedBB) & occupiedBB
	sideToMove ^= 1
	sideToMoveAttackers := attackers & pos.SideBB[sideToMove]

	if sideToMoveAttackers == 0 {
		return scores[0]
	}

	captured := pos.Squares[fromSq].Type
	scoresIndex := 1

	for {
		scores[scoresIndex] = -scores[scoresIndex-1] + PieceValues[captured]
		captured = getLeastValuableAttacker(pos, &attackers, sideToMoveAttackers, toSq)

		sideToMove ^= 1
		sideToMoveAttackers = attackers & pos.SideBB[sideToMove]

		if sideToMoveAttackers == 0 || captured == King {
			break
		}

		scoresIndex++
	}

	fmt.Println(scores, scoresIndex)

	for ; scoresIndex > 0; scoresIndex-- {
		scores[scoresIndex-1] = min16(-scores[scoresIndex], scores[scoresIndex-1])
	}

	return scores[0]
}

// Calculate a bitboard with all of the attackers for each side of a certian square,
// including x-ray attacks.
func allAttackers(pos *Position, sq uint8, occupiedBB Bitboard) (attackers Bitboard) {
	attackers |= attackersForSide(pos, White, sq, occupiedBB)
	attackers |= attackersForSide(pos, Black, sq, occupiedBB)

	for {
		occupiedBB &= ^attackers
		newAttackers := attackersForSide(pos, White, sq, occupiedBB) & ^attackers
		newAttackers |= attackersForSide(pos, Black, sq, occupiedBB) & ^attackers

		if newAttackers == 0 {
			break
		}

		attackers |= newAttackers
	}

	return attackers
}

func attackersForSide(pos *Position, attackerColor, sq uint8, occupiedBB Bitboard) (attackers Bitboard) {
	enemyBishops := pos.PieceBB[attackerColor][Bishop]
	enemyRooks := pos.PieceBB[attackerColor][Rook]
	enemyQueens := pos.PieceBB[attackerColor][Queen]
	enemyKnights := pos.PieceBB[attackerColor][Knight]
	enemyKing := pos.PieceBB[attackerColor][King]
	enemyPawns := pos.PieceBB[attackerColor][Pawn]

	intercardinalRays := genBishopMoves(sq, occupiedBB)
	cardinalRaysRays := genRookMoves(sq, occupiedBB)

	attackers |= intercardinalRays & (enemyBishops | enemyQueens)
	attackers |= cardinalRaysRays & (enemyRooks | enemyQueens)
	attackers |= KnightMoves[sq] & enemyKnights
	attackers |= KingMoves[sq] & enemyKing
	attackers |= PawnAttacks[attackerColor^1][sq] & enemyPawns
	return attackers
}

// Get the last valuable attacker from a bitboard of attackers.
func getLeastValuableAttacker(pos *Position, allAttackers *Bitboard, sideToMoveAttackers Bitboard, toSq uint8) uint8 {
	lowestValue := Inf
	pieceType := NoType
	pieceSq := uint8(NoSq)

	for sideToMoveAttackers != 0 {
		sq := sideToMoveAttackers.PopBit()
		currPieceType := pos.Squares[sq].Type
		value := PieceValues[currPieceType]

		if value < lowestValue {
			if currPieceType == Bishop || currPieceType == Rook || currPieceType == Queen {
				lineBetween := Between[sq][toSq] & ^SquareBB[sq]
				if (lineBetween & *allAttackers) != 0 {
					// Even if we've found the least valuable attacker, we need to make sure it isn't
					// blocked by another piece it's xraying through. If it is, we have to choose
					// another piece instead.
					continue
				}
			}

			lowestValue = value
			pieceType = currPieceType
			pieceSq = sq
		}
	}

	// Remove the least valuable attacker from the attacker bitboard.
	*allAttackers &= ^SquareBB[pieceSq]

	return pieceType
}

// Get the minimum between two numbers.
func min16(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}
