package engine

// movegen.go implements the move generator for Blunder.

import (
	"fmt"
)

const (
	// These masks help determine whether or not the squares between
	// the king and it's rooks are clear for castling
	F1_G1, B1_C1_D1 = 0x600000000000000, 0x7000000000000000
	F8_G8, B8_C8_D8 = 0x6, 0x70

	FullBB Bitboard = 0xffffffffffffffff
)

func GenAllMoves(pos *Position) (moves MoveList) {
	for piece := uint8(Knight); piece < NoType; piece++ {
		piecesBB := pos.PieceBB[pos.SideToMove][piece]
		for piecesBB != 0 {
			pieceSq := piecesBB.PopBit()
			genPieceMoves(pos, piece, pieceSq, &moves, FullBB)
		}
	}

	genPawnMoves(pos, &moves)
	genCastlingMoves(pos, &moves)
	return moves
}

func genCapturesAndQPromotions(pos *Position) (moves MoveList) {
	targets := pos.SideBB[pos.SideToMove^1]
	for piece := uint8(Knight); piece < NoType; piece++ {
		piecesBB := pos.PieceBB[pos.SideToMove][piece]
		for piecesBB != 0 {
			pieceSq := piecesBB.PopBit()
			genPieceMoves(pos, piece, pieceSq, &moves, targets)
		}
	}

	genPawnAttacksAndQPromotions(pos, &moves)
	return moves
}

func genPieceMoves(pos *Position, piece, sq uint8, moves *MoveList, targets Bitboard) {
	usBB := pos.SideBB[pos.SideToMove]
	enemyBB := pos.SideBB[pos.SideToMove^1]

	switch piece {
	case Knight:
		knightMoves := (KnightMoves[sq] & ^usBB) & targets
		genMovesFromBB(pos, sq, knightMoves, enemyBB, moves)
	case King:
		kingMoves := (KingMoves[sq] & ^usBB) & targets
		genMovesFromBB(pos, sq, kingMoves, enemyBB, moves)
	case Bishop:
		bishopMoves := (genBishopMoves(sq, usBB|enemyBB) & ^usBB) & targets
		genMovesFromBB(pos, sq, bishopMoves, enemyBB, moves)
	case Rook:
		rookMoves := (genRookMoves(sq, usBB|enemyBB) & ^usBB) & targets
		genMovesFromBB(pos, sq, rookMoves, enemyBB, moves)
	case Queen:
		bishopMoves := (genBishopMoves(sq, usBB|enemyBB) & ^usBB) & targets
		rookMoves := (genRookMoves(sq, usBB|enemyBB) & ^usBB) & targets
		genMovesFromBB(pos, sq, bishopMoves|rookMoves, enemyBB, moves)
	}
}

func genRookMoves(sq uint8, blockers Bitboard) Bitboard {
	magic := &RookMagics[sq]
	blockers &= magic.Mask
	return RookAttacks[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}

func genBishopMoves(sq uint8, blockers Bitboard) Bitboard {
	magic := &BishopMagics[sq]
	blockers &= magic.Mask
	return BishopAttacks[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}

func genPawnMoves(pos *Position, moves *MoveList) {
	usBB := pos.SideBB[pos.SideToMove]
	enemyBB := pos.SideBB[pos.SideToMove^1]
	pawnsBB := pos.PieceBB[pos.SideToMove][Pawn]

	for pawnsBB != 0 {
		from := pawnsBB.PopBit()

		pawnOnePush := PawnPushes[pos.SideToMove][from] & ^(usBB | enemyBB)
		pawnTwoPush := ((pawnOnePush & MaskRank[Rank6]) << 8) & ^(usBB | enemyBB)
		if pos.SideToMove == White {
			pawnTwoPush = ((pawnOnePush & MaskRank[Rank3]) >> 8) & ^(usBB | enemyBB)
		}

		pawnPush := pawnOnePush | pawnTwoPush
		pawnAttacks := PawnAttacks[pos.SideToMove][from]

		for pawnPush != 0 {
			to := pawnPush.PopBit()
			if isPromoting(pos.SideToMove, to) {
				makePromotionMoves(pos, from, to, moves)
				continue
			}
			moves.AddMove(NewMove(from, to, Quiet, NoFlag))
		}

		for pawnAttacks != 0 {
			to := pawnAttacks.PopBit()
			toBB := SquareBB[to]

			if to == pos.EPSq {
				moves.AddMove(NewMove(from, to, Attack, AttackEPFlag))
			} else if toBB&enemyBB != 0 {
				if isPromoting(pos.SideToMove, to) {
					makePromotionMoves(pos, from, to, moves)
					continue
				}
				moves.AddMove(NewMove(from, to, Attack, NoFlag))
			}
		}
	}
}

func genPawnAttacksAndQPromotions(pos *Position, moves *MoveList) {
	usBB := pos.SideBB[pos.SideToMove]
	enemyBB := pos.SideBB[pos.SideToMove^1]
	pawnsBB := pos.PieceBB[pos.SideToMove][Pawn]

	for pawnsBB != 0 {
		from := pawnsBB.PopBit()

		pawnOnePush := PawnPushes[pos.SideToMove][from] & ^(usBB | enemyBB)
		pawnAttacks := PawnAttacks[pos.SideToMove][from]

		to := pawnOnePush.PopBit()
		if isPromoting(pos.SideToMove, to) {
			moves.AddMove(NewMove(from, to, Promotion, QueenPromotionFlag))
		}

		for pawnAttacks != 0 {
			to := pawnAttacks.PopBit()
			toBB := SquareBB[to]

			if to == pos.EPSq {
				moves.AddMove(NewMove(from, to, Attack, AttackEPFlag))
			} else if toBB&enemyBB != 0 {
				if isPromoting(pos.SideToMove, to) {
					makePromotionMoves(pos, from, to, moves)
					continue
				}
				moves.AddMove(NewMove(from, to, Attack, NoFlag))
			}
		}
	}
}

func isPromoting(usColor, toSq uint8) bool {
	if usColor == White {
		return toSq >= 56 && toSq <= 63
	}
	return toSq <= 7
}

func makePromotionMoves(pos *Position, from, to uint8, moves *MoveList) {
	moves.AddMove(NewMove(from, to, Promotion, KnightPromotionFlag))
	moves.AddMove(NewMove(from, to, Promotion, BishopPromotionFlag))
	moves.AddMove(NewMove(from, to, Promotion, RookPromotionFlag))
	moves.AddMove(NewMove(from, to, Promotion, QueenPromotionFlag))
}

func genCastlingMoves(pos *Position, moves *MoveList) {
	allBB := pos.SideBB[pos.SideToMove] | pos.SideBB[pos.SideToMove^1]
	if pos.SideToMove == White {
		if pos.CastlingRights&WhiteKingsideRight != 0 && (allBB&F1_G1) == 0 && (!sqIsAttacked(pos, pos.SideToMove, E1) &&
			!sqIsAttacked(pos, pos.SideToMove, F1)) {
			moves.AddMove(NewMove(E1, G1, Castle, NoFlag))
		}
		if pos.CastlingRights&WhiteQueensideRight != 0 && (allBB&B1_C1_D1) == 0 && (!sqIsAttacked(pos, pos.SideToMove, E1) &&
			!sqIsAttacked(pos, pos.SideToMove, D1)) {
			moves.AddMove(NewMove(E1, C1, Castle, NoFlag))
		}
	} else {
		if pos.CastlingRights&BlackKingsideRight != 0 && (allBB&F8_G8) == 0 && (!sqIsAttacked(pos, pos.SideToMove, E8) &&
			!sqIsAttacked(pos, pos.SideToMove, F8)) {
			moves.AddMove(NewMove(E8, G8, Castle, NoFlag))
		}
		if pos.CastlingRights&BlackQueensideRight != 0 && (allBB&B8_C8_D8) == 0 && (!sqIsAttacked(pos, pos.SideToMove, E8) &&
			!sqIsAttacked(pos, pos.SideToMove, D8)) {
			moves.AddMove(NewMove(E8, C8, Castle, NoFlag))
		}
	}
}

func genMovesFromBB(pos *Position, from uint8, movesBB, enemyBB Bitboard, moves *MoveList) {
	for movesBB != 0 {
		to := movesBB.PopBit()
		toBB := SquareBB[to]
		moveType := Quiet
		if toBB&enemyBB != 0 {
			moveType = Attack
		}
		moves.AddMove(NewMove(from, to, moveType, NoFlag))
	}
}

func sqIsAttacked(pos *Position, usColor, sq uint8) bool {
	enemyBB := pos.SideBB[usColor^1]
	usBB := pos.SideBB[usColor]

	enemyBishops := pos.PieceBB[usColor^1][Bishop]
	enemyRooks := pos.PieceBB[usColor^1][Rook]
	enemyQueens := pos.PieceBB[usColor^1][Queen]
	enemyKnights := pos.PieceBB[usColor^1][Knight]
	enemyKing := pos.PieceBB[usColor^1][King]
	enemyPawns := pos.PieceBB[usColor^1][Pawn]

	if KnightMoves[sq]&enemyKnights != 0 {
		return true
	}
	if KingMoves[sq]&enemyKing != 0 {
		return true
	}
	if PawnAttacks[usColor][sq]&enemyPawns != 0 {
		return true
	}

	intercardinalRays := genBishopMoves(sq, enemyBB|usBB)
	if intercardinalRays&(enemyBishops|enemyQueens) != 0 {
		return true
	}

	cardinalRays := genRookMoves(sq, enemyBB|usBB)
	return cardinalRays&(enemyRooks|enemyQueens) != 0
}

func DividePerft(pos *Position, depth, divdeAt uint8) uint64 {
	if depth == 0 {
		return 1
	}

	moves := GenAllMoves(pos)
	var nodes uint64

	var idx uint8
	for idx = 0; idx < moves.Count; idx++ {
		move := moves.Moves[idx]
		if pos.DoMove(move) {
			moveNodes := DividePerft(pos, depth-1, divdeAt)
			if depth == divdeAt {
				fmt.Printf("%v: %v\n", move, moveNodes)
			}

			nodes += moveNodes
		}

		pos.UndoMove(move)
	}

	return nodes
}

func Perft(pos *Position, depth uint8) uint64 {
	if depth == 0 {
		return 1
	}

	moves := GenAllMoves(pos)
	var nodes uint64

	var idx uint8
	for idx = 0; idx < moves.Count; idx++ {
		if pos.DoMove(moves.Moves[idx]) {
			nodes += Perft(pos, depth-1)
		}
		pos.UndoMove(moves.Moves[idx])
	}

	return nodes
}
