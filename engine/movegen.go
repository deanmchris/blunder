package engine

import "fmt"

// A file containg the move generator of Blunder

const (
	// These masks help determine whether or not the squares between
	// the king and it's rooks are clear for castling
	F1_G1, B1_C1_D1 = 0x600000000000000, 0x7000000000000000
	F8_G8, B8_C8_D8 = 0x6, 0x70
	FullBB          = 0xffffffffffffffff
)

// Generate all pseduo-legal moves for a given position.
func GenPseduoLegalMoves(board *Board) (moves Moves) {
	kingPos := board.KingPos[board.ColorToMove][King]
	checkers := attackersOfSquare(board, board.ColorToMove, kingPos)
	var targets Bitboard = FullBB

	if checkers.CountBits() > 1 {
		genPieceMoves(board, King, kingPos, &moves, targets, false)
		return
	} else if checkers.CountBits() == 1 {
		checkerPos := checkers.Msb()
		if board.Squares[checkerPos].Type == Knight {
			targets = SquareBB[checkerPos]
		} else {
			targets = LinesBewteen[kingPos][checkerPos]
		}
	}

	moves = make([]Move, 0, 100)
	var piece uint8

	for piece = Knight; piece < NoType; piece++ {
		piecesBB := board.PieceBB[board.ColorToMove][piece]
		for piecesBB != 0 {
			piecePos := piecesBB.PopBit()
			genPieceMoves(board, piece, piecePos, &moves, targets, false)
		}
	}

	genPawnMoves(board, &moves, targets)
	genCastlingMoves(board, &moves)
	return moves
}

// Generate all pseduo-legal moves for a given position.
func GenPseduoLegalCaptures(board *Board) (moves Moves) {
	kingPos := board.KingPos[board.ColorToMove][King]
	checkers := attackersOfSquare(board, board.ColorToMove, kingPos)
	var targets = board.SideBB[board.ColorToMove^1]

	if checkers.CountBits() > 1 {
		genPieceMoves(board, King, kingPos, &moves, targets, true)
		return
	} else if checkers.CountBits() == 1 {
		checkerPos := checkers.Msb()
		targets = SquareBB[checkerPos]
	}

	moves = make([]Move, 0, 100)
	var piece uint8

	for piece = Knight; piece < NoType; piece++ {
		piecesBB := board.PieceBB[board.ColorToMove][piece]
		for piecesBB != 0 {
			piecePos := piecesBB.PopBit()
			genPieceMoves(board, piece, piecePos, &moves, targets, true)
		}
	}
	genPawnMoves(board, &moves, targets)
	return moves
}

// Generate the moves a single piece, making sure they all align with the squares specified
// by a targets bitboard. Filter king moves using the target bitboard if includeKing is set to
// true.
func genPieceMoves(board *Board, piece, sq uint8, moves *Moves, targets Bitboard, includeKing bool) {
	usBB := board.SideBB[board.ColorToMove]
	enemyBB := board.SideBB[board.ColorToMove^1]

	switch piece {
	case Knight:
		knightMoves := (KnightMoves[sq] & ^usBB) & targets
		genMovesFromBB(board, sq, knightMoves, enemyBB, moves)
	case King:
		kingMoves := KingMoves[sq] & ^usBB
		if includeKing {
			kingMoves &= targets
		}
		genMovesFromBB(board, sq, kingMoves, enemyBB, moves)
	case Bishop:
		bishopMoves := (genBishopMoves(sq, usBB|enemyBB) & ^usBB) & targets
		genMovesFromBB(board, sq, bishopMoves, enemyBB, moves)
	case Rook:
		rookMoves := (genRookMoves(sq, usBB|enemyBB) & ^usBB) & targets
		genMovesFromBB(board, sq, rookMoves, enemyBB, moves)
	case Queen:
		bishopMoves := (genBishopMoves(sq, usBB|enemyBB) & ^usBB) & targets
		rookMoves := (genRookMoves(sq, usBB|enemyBB) & ^usBB) & targets
		genMovesFromBB(board, sq, bishopMoves|rookMoves, enemyBB, moves)
	}
}

// Generate rook moves.
func genRookMoves(sq uint8, blockers Bitboard) Bitboard {
	magic := &RookMagics[sq]
	blockers &= magic.Mask
	return RookAttacks[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}

// Generate rook moves.
func genBishopMoves(sq uint8, blockers Bitboard) Bitboard {
	magic := &BishopMagics[sq]
	blockers &= magic.Mask
	return BishopAttacks[sq][(uint64(blockers)*magic.MagicNo)>>magic.Shift]
}

// Generate pawn moves for the current side. Pawns are treated
// seperately from the rest of the pieces as they have more
// complicated and exceptional rules for how they can move.
func genPawnMoves(board *Board, moves *Moves, targets Bitboard) {
	usBB := board.SideBB[board.ColorToMove]
	enemyBB := board.SideBB[board.ColorToMove^1]
	pawnsBB := board.PieceBB[board.ColorToMove][Pawn]

	for pawnsBB != 0 {
		from := pawnsBB.PopBit()
		pawnOnePush := PawnPushes[board.ColorToMove][from] & ^(usBB | enemyBB)

		pawnTwoPush := ((pawnOnePush & MaskRank[Rank6]) << 8) & ^(usBB | enemyBB)
		if board.ColorToMove == White {
			pawnTwoPush = ((pawnOnePush & MaskRank[Rank3]) >> 8) & ^(usBB | enemyBB)
		}

		pawnPush := (pawnOnePush | pawnTwoPush) & targets
		pawnAttacks := PawnAttacks[board.ColorToMove][from] & (targets | SquareBB[board.EPSquare])
		for pawnPush != 0 {
			to := pawnPush.PopBit()
			if isPromoting(board.ColorToMove, to) {
				makePromotionMoves(board, from, to, moves)
				continue
			}
			if abs(int8(from)-int8(to)) == 16 {
				*moves = append(*moves, MakeMove(from, to, DoublePawnPush))
				continue
			}
			*moves = append(*moves, MakeMove(from, to, Quiet))
		}
		for pawnAttacks != 0 {
			to := pawnAttacks.PopBit()
			toBB := SquareBB[to]

			if to == board.EPSquare {
				*moves = append(*moves, MakeMove(from, to, AttackEP))
			} else if toBB&enemyBB != 0 {
				if isPromoting(board.ColorToMove, to) {
					makePromotionMoves(board, from, to, moves)
					continue
				}
				*moves = append(*moves, MakeMove(from, to, Attack))
			}
		}
	}
}

// Get the absolute value of a number n
func abs(n int8) int8 {
	if n < 0 {
		return -n
	}
	return n
}

// A helper function to determine if a pawn has reached the 8th or
// 1st rank and will promote.
func isPromoting(usColor, toSq uint8) bool {
	if usColor == White {
		return toSq >= 56 && toSq <= 63
	}
	return toSq <= 7
}

// Generate promotion moves for pawns
func makePromotionMoves(board *Board, from, to uint8, moves *Moves) {
	*moves = append(*moves, MakeMove(from, to, KnightPromotion))
	*moves = append(*moves, MakeMove(from, to, BishopPromotion))
	*moves = append(*moves, MakeMove(from, to, RookPromotion))
	*moves = append(*moves, MakeMove(from, to, QueenPromotion))
}

// Generate castling moves. Note testing whether or not castling has the king
// crossing attacked squares is not tested for here, as pseduo-legal move
// generation is the focus.
func genCastlingMoves(board *Board, moves *Moves) {
	allPieces := board.SideBB[board.ColorToMove] | board.SideBB[board.ColorToMove^1]
	if board.ColorToMove == White {
		if board.CastlingRights&WhiteKingside == WhiteKingside && (allPieces&F1_G1) == 0 && (!sqIsAttacked(board, board.ColorToMove, E1) &&
			!sqIsAttacked(board, board.ColorToMove, F1) && !sqIsAttacked(board, board.ColorToMove, G1)) {
			*moves = append(*moves, MakeMove(E1, G1, CastleWKS))
		}
		if board.CastlingRights&WhiteQueenside == WhiteQueenside && (allPieces&B1_C1_D1) == 0 && (!sqIsAttacked(board, board.ColorToMove, E1) &&
			!sqIsAttacked(board, board.ColorToMove, D1) && !sqIsAttacked(board, board.ColorToMove, C1)) {
			*moves = append(*moves, MakeMove(E1, C1, CastleWQS))
		}
	} else {
		if board.CastlingRights&BlackKingside == BlackKingside && (allPieces&F8_G8) == 0 && (!sqIsAttacked(board, board.ColorToMove, E8) &&
			!sqIsAttacked(board, board.ColorToMove, F8) && !sqIsAttacked(board, board.ColorToMove, G8)) {
			*moves = append(*moves, MakeMove(E8, G8, CastleBKS))
		}
		if board.CastlingRights&BlackQueenside == BlackQueenside && (allPieces&B8_C8_D8) == 0 && (!sqIsAttacked(board, board.ColorToMove, E8) &&
			!sqIsAttacked(board, board.ColorToMove, D8) && !sqIsAttacked(board, board.ColorToMove, C8)) {
			*moves = append(*moves, MakeMove(E8, C8, CastleBQS))
		}
	}
}

// From a bitboard representing possible squares a piece can move,
// serialize it, and generate a list of moves.
func genMovesFromBB(board *Board, from uint8, movesBB, enemyBB Bitboard, moves *Moves) {
	for movesBB != 0 {
		to := movesBB.PopBit()
		toBB := SquareBB[to]
		moveType := Quiet
		if toBB&enemyBB != 0 {
			moveType = Attack
		}
		*moves = append(*moves, MakeMove(from, to, moveType))
	}
}

func sqIsAttacked(board *Board, usColor, sq uint8) bool {
	enemyBB := board.SideBB[usColor^1]
	usBB := board.SideBB[usColor]

	enemyBishops := board.PieceBB[usColor^1][Bishop]
	enemyRooks := board.PieceBB[usColor^1][Rook]
	enemyQueens := board.PieceBB[usColor^1][Queen]
	enemyKnights := board.PieceBB[usColor^1][Knight]
	enemyKing := board.PieceBB[usColor^1][King]
	enemyPawns := board.PieceBB[usColor^1][Pawn]

	intercardinalRays := genBishopMoves(sq, enemyBB|usBB)
	cardinalRaysRays := genRookMoves(sq, enemyBB|usBB)

	if intercardinalRays&(enemyBishops|enemyQueens) != 0 {
		return true
	}
	if cardinalRaysRays&(enemyRooks|enemyQueens) != 0 {
		return true
	}

	if KnightMoves[sq]&enemyKnights != 0 {
		return true
	}
	if KingMoves[sq]&enemyKing != 0 {
		return true
	}
	if PawnAttacks[usColor][sq]&enemyPawns != 0 {
		return true
	}
	return false
}

// Compute a bitboard representing the enemy attackers of a particular square.
func attackersOfSquare(board *Board, usColor, sq uint8) (attackers Bitboard) {
	enemyBB := board.SideBB[usColor^1]
	usBB := board.SideBB[usColor]

	enemyBishops := board.PieceBB[usColor^1][Bishop]
	enemyRooks := board.PieceBB[usColor^1][Rook]
	enemyQueens := board.PieceBB[usColor^1][Queen]
	enemyKnights := board.PieceBB[usColor^1][Knight]
	enemyKing := board.PieceBB[usColor^1][King]
	enemyPawns := board.PieceBB[usColor^1][Pawn]

	intercardinalRays := genBishopMoves(sq, enemyBB|usBB)
	cardinalRaysRays := genRookMoves(sq, enemyBB|usBB)

	attackers |= intercardinalRays & (enemyBishops | enemyQueens)
	attackers |= cardinalRaysRays & (enemyRooks | enemyQueens)
	attackers |= KnightMoves[sq] & enemyKnights
	attackers |= KingMoves[sq] & enemyKing
	attackers |= PawnAttacks[usColor][sq] & enemyPawns
	return attackers
}

const TTSize = 10000000

type perftTTEntry struct {
	Hash      uint64
	NodeCount uint64
	Depth     int
}

var TT [TTSize]perftTTEntry

// Explore the move tree up to depth, and return the total
// number of nodes explored.  This function is used to
// debug move generation and ensure it is working by comparing
// the results to the known results of other engines
func Perft(board *Board, depth, divdeAt int, silent bool) uint64 {
	if depth == 0 {
		return 1
	}

	if entry := TT[board.Hash%TTSize]; entry.Hash == board.Hash && entry.Depth == depth {
		return entry.NodeCount
	}

	moves := GenPseduoLegalMoves(board)
	var nodes uint64

	for _, move := range moves {
		board.DoMove(move, true)
		if board.KingIsAttacked(board.ColorToMove ^ 1) {
			board.UndoMove(move)
			continue
		}

		moveNodes := Perft(board, depth-1, divdeAt, silent)

		if depth == divdeAt && !silent {
			fmt.Printf("%v: %v\n", move, moveNodes)
		}

		nodes += moveNodes
		board.UndoMove(move)
	}

	TT[board.Hash%TTSize] = perftTTEntry{Hash: board.Hash, NodeCount: nodes, Depth: depth}
	return nodes
}
