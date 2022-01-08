package engine

// see.go is a simple implementation of a static exchange evaluator.

var PieceValues [7]int16 = [7]int16{
	100,
	300,
	300,
	500,
	900,
	Inf,
	0,
}

func (pos *Position) See(move Move) int16 {
	toSq := move.ToSq()
	frSQ := move.FromSq()
	target := pos.Squares[toSq].Type
	attacker := pos.Squares[frSQ].Type

	gain := [32]int16{}
	depth := uint8(0)
	sideToMove := pos.SideToMove ^ 1

	seenBB := Bitboard(0)
	occupiedBB := pos.SideBB[White] | pos.SideBB[Black]
	attackerBB := SquareBB[frSQ]

	attadef := pos.allAttackers(toSq, occupiedBB)
	maxXray := occupiedBB & ^(pos.PieceBB[White][Knight] | pos.PieceBB[White][King] |
		pos.PieceBB[Black][Knight] | pos.PieceBB[Black][King])

	gain[depth] = PieceValues[target]

	for ok := true; ok; ok = attackerBB != 0 {
		depth++
		gain[depth] = PieceValues[attacker] - gain[depth-1]

		if max(-gain[depth-1], gain[depth]) < 0 {
			break
		}

		attadef &= ^attackerBB
		occupiedBB &= ^attackerBB
		seenBB |= attackerBB

		if (attackerBB & maxXray) != 0 {
			attadef |= pos.considerXrays(toSq, occupiedBB) & ^seenBB
		}

		attackerBB = pos.minAttacker(attadef, sideToMove, &attacker)
		sideToMove ^= 1
	}

	for depth--; depth > 0; depth-- {
		gain[depth-1] = -max(-gain[depth-1], gain[depth])
	}

	return gain[0]
}

func (pos *Position) minAttacker(attadef Bitboard, color uint8, attacker *uint8) Bitboard {
	for *attacker = Pawn; *attacker <= King; *attacker++ {
		subset := attadef & pos.PieceBB[color][*attacker]
		if subset != 0 {
			return subset & -subset
		}
	}
	return 0
}

func (pos *Position) considerXrays(sq uint8, occupiedBB Bitboard) (attackers Bitboard) {
	attackingBishops := pos.PieceBB[White][Bishop] | pos.PieceBB[Black][Bishop]
	attackingRooks := pos.PieceBB[White][Rook] | pos.PieceBB[Black][Rook]
	attackingQueens := pos.PieceBB[White][Queen] | pos.PieceBB[Black][Queen]

	intercardinalRays := genBishopMoves(sq, occupiedBB)
	cardinalRaysRays := genRookMoves(sq, occupiedBB)

	attackers |= intercardinalRays & (attackingBishops | attackingQueens)
	attackers |= cardinalRaysRays & (attackingRooks | attackingQueens)
	return attackers
}

func (pos *Position) allAttackers(sq uint8, occupiedBB Bitboard) (attackers Bitboard) {
	attackers |= pos.attackersForSide(White, sq, occupiedBB)
	attackers |= pos.attackersForSide(Black, sq, occupiedBB)
	return attackers
}

func (pos *Position) attackersForSide(attackerColor, sq uint8, occupiedBB Bitboard) (attackers Bitboard) {
	attackingBishops := pos.PieceBB[attackerColor][Bishop]
	attackingRooks := pos.PieceBB[attackerColor][Rook]
	attackingQueens := pos.PieceBB[attackerColor][Queen]
	attackingKnights := pos.PieceBB[attackerColor][Knight]
	attackingKing := pos.PieceBB[attackerColor][King]
	attackingPawns := pos.PieceBB[attackerColor][Pawn]

	intercardinalRays := genBishopMoves(sq, occupiedBB)
	cardinalRaysRays := genRookMoves(sq, occupiedBB)

	attackers |= intercardinalRays & (attackingBishops | attackingQueens)
	attackers |= cardinalRaysRays & (attackingRooks | attackingQueens)
	attackers |= KnightMoves[sq] & attackingKnights
	attackers |= KingMoves[sq] & attackingKing
	attackers |= PawnAttacks[attackerColor^1][sq] & attackingPawns
	return attackers
}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}
