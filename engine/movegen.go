package engine

import "fmt"

const (
	F1_G1, B1_C1_D1 = 0x600000000000000, 0x7000000000000000
	F8_G8, B8_C8_D8 = 0x6, 0x70
)

func genAllMoves(pos *Position) (moves MoveList) {
	usBB := pos.Sides[pos.SideToMove]
	enemyBB := pos.Sides[pos.SideToMove^1]

	genKingMoves(pos.Pieces[King]&usBB, FullBB, enemyBB, usBB, &moves)
	genKnightMoves(pos.Pieces[Knight]&usBB, FullBB, enemyBB, usBB, &moves)
	genRookMoves(pos.Pieces[Rook]&usBB, FullBB, enemyBB, usBB, &moves)
	genBishopMoves(pos.Pieces[Bishop]&usBB, FullBB, enemyBB, usBB, &moves)
	genQueenMoves(pos.Pieces[Queen]&usBB, FullBB, enemyBB, usBB, &moves)

	genPawnMoves(
		pos.Pieces[Pawn]&usBB,
		enemyBB, usBB, pos.SideToMove, pos.EPSq, &moves,
	)

	genCastlingMoves(pos, &moves)

	return moves
}

func genAttacks(pos *Position) (moves MoveList) {
	usBB := pos.Sides[pos.SideToMove]
	enemyBB := pos.Sides[pos.SideToMove^1]

	genKingMoves(pos.Pieces[King]&usBB, enemyBB, enemyBB, usBB, &moves)
	genKnightMoves(pos.Pieces[Knight]&usBB, enemyBB, enemyBB, usBB, &moves)
	genBishopMoves(pos.Pieces[Bishop]&usBB, enemyBB, enemyBB, usBB, &moves)
	genRookMoves(pos.Pieces[Rook]&usBB, enemyBB, enemyBB, usBB, &moves)
	genQueenMoves(pos.Pieces[Queen]&usBB, enemyBB, enemyBB, usBB, &moves)

	genPawnAttacks(
		pos.Pieces[Pawn]&usBB,
		enemyBB, usBB, pos.SideToMove, pos.EPSq, &moves,
	)

	return moves
}

func genKingMoves(kingBB, filter, enemyBB, usBB uint64, moves *MoveList) {
	from := BitScanAndClear(&kingBB)
	genMovesFromBB(KingMoves[from]&^usBB&filter, enemyBB, from, moves)
}

func genKnightMoves(knightBB, filter, enemyBB, usBB uint64, moves *MoveList) {
	for knightBB != 0 {
		from := BitScanAndClear(&knightBB)
		genMovesFromBB(KnightMoves[from]&^usBB&filter, enemyBB, from, moves)
	}
}

func genBishopMoves(bishopBB, filter, enemyBB, usBB uint64, moves *MoveList) {
	for bishopBB != 0 {
		from := BitScanAndClear(&bishopBB)
		genMovesFromBB(genBishopMovesBB(from, usBB|enemyBB)&^usBB&filter, enemyBB, from, moves)
	}
}

func genRookMoves(rookBB, filter, enemyBB, usBB uint64, moves *MoveList) {
	for rookBB != 0 {
		from := BitScanAndClear(&rookBB)
		genMovesFromBB(genRookMovesBB(from, usBB|enemyBB)&^usBB&filter, enemyBB, from, moves)
	}
}

func genQueenMoves(queenBB, filter, enemyBB, usBB uint64, moves *MoveList) {
	for queenBB != 0 {
		from := BitScanAndClear(&queenBB)
		movesBB := genBishopMovesBB(from, usBB|enemyBB) | genRookMovesBB(from, usBB|enemyBB)
		genMovesFromBB(movesBB&^usBB&filter, enemyBB, from, moves)
	}
}

func genPawnMoves(pawnsBB, enemyBB, usBB uint64, stm, epSq uint8, moves *MoveList) {
	for pawnsBB != 0 {
		from := BitScanAndClear(&pawnsBB)

		pawnOnePush := PawnPushes[stm][from] & ^(usBB | enemyBB)
		pawnTwoPush := EmptyBB

		if stm == White {
			pawnTwoPush = ((pawnOnePush & MaskRank[Rank3]) >> 8) & ^(usBB | enemyBB)
		} else {
			pawnTwoPush = ((pawnOnePush & MaskRank[Rank6]) << 8) & ^(usBB | enemyBB)
		}

		pawnPush := pawnOnePush | pawnTwoPush
		pawnAttacks := PawnAttacks[stm][from] & (enemyBB | SquareBB[epSq])

		for pawnPush != 0 {
			to := BitScanAndClear(&pawnPush)
			if isPromoting(to) {
				makePromotionMoves(from, to, moves)
				continue
			}
			moves.AddMove(newMove(from, to, Quiet, NoFlag))
		}

		for pawnAttacks != 0 {
			to := BitScanAndClear(&pawnAttacks)
			if to == epSq {
				moves.AddMove(newMove(from, to, Attack, AttackEP))
				continue
			}

			if isPromoting(to) {
				makePromotionMoves(from, to, moves)
				continue
			}

			moves.AddMove(newMove(from, to, Attack, NoFlag))
		}
	}
}

func genPawnAttacks(pawnsBB, enemyBB, usBB uint64, stm, epSq uint8, moves *MoveList) {
	for pawnsBB != 0 {
		from := BitScanAndClear(&pawnsBB)
		pawnAttacks := PawnAttacks[stm][from] & (enemyBB | SquareBB[epSq])

		for pawnAttacks != 0 {
			to := BitScanAndClear(&pawnAttacks)
			if to == epSq {
				moves.AddMove(newMove(from, to, Attack, AttackEP))
				continue
			}

			if isPromoting(to) {
				moves.AddMove(newMove(from, to, Promotion, QueenPromotion))
				continue
			}

			moves.AddMove(newMove(from, to, Attack, NoFlag))
		}
	}
}

func genCastlingMoves(pos *Position, moves *MoveList) {
	allPieces := pos.Sides[pos.SideToMove] | pos.Sides[pos.SideToMove^1]
	if pos.SideToMove == White {
		if (pos.CastlingRights&WhiteKingsideRight) != 0 && (allPieces&F1_G1) == 0 {
			moves.AddMove(newMove(E1, G1, Castle, NoFlag))
		}

		if (pos.CastlingRights&WhiteQueensideRight) != 0 && (allPieces&B1_C1_D1) == 0 {
			moves.AddMove(newMove(E1, C1, Castle, NoFlag))
		}
	} else {
		if (pos.CastlingRights&BlackKingsideRight) != 0 && (allPieces&F8_G8) == 0 {
			moves.AddMove(newMove(E8, G8, Castle, NoFlag))
		}

		if (pos.CastlingRights&BlackQueensideRight) != 0 && (allPieces&B8_C8_D8) == 0 {
			moves.AddMove(newMove(E8, C8, Castle, NoFlag))
		}
	}
}

func sqIsAttacked(pos *Position, usColor, sq uint8) bool {
	enemyBB := pos.Sides[usColor^1]
	usBB := pos.Sides[usColor]

	enemyKnights := pos.Pieces[Knight] & pos.Sides[usColor^1]
	enemyKing := pos.Pieces[King] & pos.Sides[usColor^1]
	enemyPawns := pos.Pieces[Pawn] & pos.Sides[usColor^1]

	if KnightMoves[sq]&enemyKnights != 0 {
		return true
	}
	if KingMoves[sq]&enemyKing != 0 {
		return true
	}
	if PawnAttacks[usColor][sq]&enemyPawns != 0 {
		return true
	}

	enemyBishops := pos.Pieces[Bishop] & pos.Sides[usColor^1]
	enemyRooks := pos.Pieces[Rook] & pos.Sides[usColor^1]
	enemyQueens := pos.Pieces[Queen] & pos.Sides[usColor^1]

	intercardinalRays := genBishopMovesBB(sq, enemyBB|usBB)
	cardinalRaysRays := genRookMovesBB(sq, enemyBB|usBB)

	if intercardinalRays&(enemyBishops|enemyQueens) != 0 {
		return true
	}
	if cardinalRaysRays&(enemyRooks|enemyQueens) != 0 {
		return true
	}

	return false
}

func genMovesFromBB(movesBB, enemyBB uint64, from uint8, moves *MoveList) {
	for movesBB != 0 {
		to := BitScanAndClear(&movesBB)
		moveType := Quiet
		if bitIsSet(enemyBB, to) {
			moveType = Attack
		}
		moves.AddMove(newMove(from, to, moveType, NoFlag))
	}
}

func genRookMovesBB(sq uint8, blockers uint64) uint64 {
	magic := &RookMagics[sq]
	blockers &= magic.BlockerMask
	return RookMoves[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}

func genBishopMovesBB(sq uint8, blockers uint64) uint64 {
	magic := &BishopMagics[sq]
	blockers &= magic.BlockerMask
	return BishopMoves[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}

func isPromoting(toSq uint8) bool {
	// Since a pawn can never move to the backrank
	// of it's own color, we don't care about
	// checking the side to move.
	return (toSq >= 56 && toSq <= 63) || toSq <= 7
}

func makePromotionMoves(from, to uint8, moves *MoveList) {
	moves.AddMove(newMove(from, to, Promotion, KnightPromotion))
	moves.AddMove(newMove(from, to, Promotion, BishopPromotion))
	moves.AddMove(newMove(from, to, Promotion, RookPromotion))
	moves.AddMove(newMove(from, to, Promotion, QueenPromotion))
}

func Perft(pos *Position, depth uint8) uint64 {
	if depth == 0 {
		return 1
	}

	moves := genAllMoves(pos)
	nodes := uint64(0)

	pos.ComputePinAndCheckInfo()

	for i := uint8(0); i < moves.Count; i++ {
		move := moves.Moves[i]
		if pos.DoMove(move) {
			nodes += Perft(pos, depth-1)
		}
		pos.UndoMove(move)
	}

	return nodes
}

func DividePerft(pos *Position, depth uint8) uint64 {
	moves := genAllMoves(pos)
	nodes := uint64(0)

	pos.ComputePinAndCheckInfo()

	for i := uint8(0); i < moves.Count; i++ {
		move := moves.Moves[i]
		if pos.DoMove(move) {
			moveNodeCnt := Perft(pos, depth-1)
			fmt.Printf("%s: %d\n", moveToStr(move), moveNodeCnt)
			nodes += moveNodeCnt
		}
		pos.UndoMove(move)
	}

	return nodes
}
