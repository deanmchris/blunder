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

// Peform a static exchange evaluation on target square of the move given,
// and return a score of the move from the perspective of the side to move.
func (pos *Position) See(move Move) int16 {
	toSq := move.ToSq()
	frSQ := move.FromSq()
	target := pos.Squares[toSq].Type
	attacker := pos.Squares[frSQ].Type

	gain := [32]int16{}
	depth := uint8(0)
	sideToMove := pos.SideToMove ^ 1

	seenBB := Bitboard(0)
	occupiedBB := pos.Sides[White] | pos.Sides[Black]
	attackerBB := SquareBB[frSQ]

	attadef := pos.allAttackers(toSq, occupiedBB)
	maxXray := occupiedBB & ^(pos.Pieces[White][Knight] | pos.Pieces[White][King] |
		pos.Pieces[Black][Knight] | pos.Pieces[Black][King])

	gain[depth] = PieceValues[target]

	for ok := true; ok; ok = attackerBB != 0 {
		depth++
		gain[depth] = PieceValues[attacker] - gain[depth-1]

		if max16(-gain[depth-1], gain[depth]) < 0 {
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
		gain[depth-1] = -max16(-gain[depth-1], gain[depth])
	}

	return gain[0]
}

func (pos *Position) minAttacker(attadef Bitboard, color uint8, attacker *uint8) Bitboard {
	for *attacker = Pawn; *attacker <= King; *attacker++ {
		subset := attadef & pos.Pieces[color][*attacker]
		if subset != 0 {
			return subset & -subset
		}
	}
	return 0
}

func (pos *Position) considerXrays(sq uint8, occupiedBB Bitboard) (attackers Bitboard) {
	attackingBishops := pos.Pieces[White][Bishop] | pos.Pieces[Black][Bishop]
	attackingRooks := pos.Pieces[White][Rook] | pos.Pieces[Black][Rook]
	attackingQueens := pos.Pieces[White][Queen] | pos.Pieces[Black][Queen]

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
	attackingBishops := pos.Pieces[attackerColor][Bishop]
	attackingRooks := pos.Pieces[attackerColor][Rook]
	attackingQueens := pos.Pieces[attackerColor][Queen]
	attackingKnights := pos.Pieces[attackerColor][Knight]
	attackingKing := pos.Pieces[attackerColor][King]
	attackingPawns := pos.Pieces[attackerColor][Pawn]

	intercardinalRays := genBishopMoves(sq, occupiedBB)
	cardinalRaysRays := genRookMoves(sq, occupiedBB)

	attackers |= intercardinalRays & (attackingBishops | attackingQueens)
	attackers |= cardinalRaysRays & (attackingRooks | attackingQueens)
	attackers |= KnightMoves[sq] & attackingKnights
	attackers |= KingMoves[sq] & attackingKing
	attackers |= PawnAttacks[attackerColor^1][sq] & attackingPawns
	return attackers
}
