package engine

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
	"unicode"
)

const (
	Pawn   uint8 = 0
	Knight uint8 = 1
	Bishop uint8 = 2
	Rook   uint8 = 3
	Queen  uint8 = 4
	King   uint8 = 5
	NoType uint8 = 6

	Black   uint8 = 0
	White   uint8 = 1
	NoColor uint8 = 2

	WhiteKingsideRight  uint8 = 0x8
	WhiteQueensideRight uint8 = 0x4
	BlackKingsideRight  uint8 = 0x2
	BlackQueensideRight uint8 = 0x1

	FENStartPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0"
	FENKiwiPete      = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
)

var PieceTypeToChar = map[uint8]rune{
	Pawn:   'i',
	Knight: 'n',
	Bishop: 'b',
	Rook:   'r',
	Queen:  'q',
	King:   'k',
	NoType: '.',
}

var CharToPiece = map[byte]Piece{
	'P': {Pawn, White},
	'N': {Knight, White},
	'B': {Bishop, White},
	'R': {Rook, White},
	'Q': {Queen, White},
	'K': {King, White},
	'p': {Pawn, Black},
	'n': {Knight, Black},
	'b': {Bishop, Black},
	'r': {Rook, Black},
	'q': {Queen, Black},
	'k': {King, Black},
}

var ToSqToCastlingInfo = map[uint8]CastlingInfo{
	G1: {H1, F1, F1, G1},
	C1: {A1, D1, D1, C1},
	G8: {H8, F8, F8, G8},
	C8: {A8, D8, D8, C8},
}

var Spoilers = [64]uint8{
	0xb, 0xf, 0xf, 0xf, 0x3, 0xf, 0xf, 0x7,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf,
	0xe, 0xf, 0xf, 0xf, 0xc, 0xf, 0xf, 0xd,
}

var PieceTypeLookup = [33]uint8{
	NoType, King, Queen, 0, Rook, 0,
	0, 0, Bishop, 0, 0, 0, 0, 0, 0, 0,
	Knight, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, Pawn,
}

var PieceColorLookup = [3]uint8{
	NoColor, White, Black,
}

type Piece struct {
	Type,
	Color uint8
}

type CastlingInfo struct {
	RookFromSq,
	RookToSq,
	FirstSqCrossed,
	SecondSqCrossed uint8
}

type State struct {
	CastlingRights,
	EPSq,
	HalfMoveClock,
	MovedType,
	MovedColor,
	CapturedType,
	CapturedColor uint8

	Hash,
	Pinned uint64
	InCheck bool

	Phase uint8
	MGScores,
	EGScores [2]int16
}

type Position struct {
	Pieces [6]uint64
	Sides  [2]uint64
	CastlingRights,
	SideToMove,
	HalfMoveClock,
	EPSq uint8

	Hash,
	pinned uint64
	InCheck bool

	prevStates [MaxPly]State
	StateIdx   uint8

	Phase uint8
	MGScores,
	EGScores [2]int16
}

func NewPosition(fen string) Position {
	pos := Position{}
	pos.LoadFEN(fen)
	return pos
}

func (pos *Position) LoadFEN(fen string) {
	pos.Pieces = [6]uint64{}
	pos.Sides = [2]uint64{}
	pos.CastlingRights = 0
	pos.EPSq = NoSq
	pos.pinned = EmptyBB
	pos.InCheck = false

	pos.prevStates = [MaxPly]State{}
	pos.StateIdx = 0

	pos.Phase = 0
	pos.MGScores = [2]int16{}
	pos.EGScores = [2]int16{}

	fields := strings.Fields(fen)
	piecePlacement := fields[0]
	sideToMove := fields[1]
	castlingRights := fields[2]
	epSq := fields[3]
	halfMoveClock := fields[4]

	for index, sq := 0, uint8(56); index < len(piecePlacement); index++ {
		char := piecePlacement[index]
		switch char {
		case 'p', 'n', 'b', 'r', 'q', 'k', 'P', 'N', 'B', 'R', 'Q', 'K':
			piece := CharToPiece[char]
			pos.zobristPutPiece(piece.Type, piece.Color, sq)
			sq++
		case '/':
			sq -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8':
			sq += piecePlacement[index] - '0'
		}
	}

	pos.SideToMove = Black
	if sideToMove == "w" {
		pos.SideToMove = White
	}

	pos.EPSq = NoSq
	if epSq != "-" {
		pos.EPSq = coordToSq(epSq)
	}

	halfMoveCounter, err := strconv.Atoi(halfMoveClock)

	if err != nil {
		panic(err)
	}

	pos.HalfMoveClock = uint8(halfMoveCounter)

	for _, char := range castlingRights {
		switch char {
		case 'K':
			pos.CastlingRights |= WhiteKingsideRight
		case 'Q':
			pos.CastlingRights |= WhiteQueensideRight
		case 'k':
			pos.CastlingRights |= BlackKingsideRight
		case 'q':
			pos.CastlingRights |= BlackQueensideRight
		}
	}

	pos.Hash = Zobrist.GenHash(pos)
}

func (pos Position) GenFEN() string {
	positionStr := strings.Builder{}

	for rankStartPos := 56; rankStartPos >= 0; rankStartPos -= 8 {
		emptySquares := 0
		for sq := rankStartPos; sq < rankStartPos+8; sq++ {
			pieceType := pos.GetPieceType(uint8(sq))
			if pieceType == NoType {
				emptySquares++
			} else {
				// If we have some consecutive empty squares, add them to the FEN
				// string board, and reset the empty squares counter.
				if emptySquares > 0 {
					positionStr.WriteString(strconv.Itoa(emptySquares))
					emptySquares = 0
				}

				pieceType := pos.GetPieceType(uint8(sq))
				pieceColor := pos.GetPieceColor(uint8(sq))
				pieceChar := PieceTypeToChar[pieceType]

				if pieceColor == White {
					pieceChar = unicode.ToUpper(pieceChar)
				}

				// In FEN strings pawns are represented with p's not i's
				if pieceChar == 'i' {
					pieceChar = 'p'
				} else if pieceChar == 'I' {
					pieceChar = 'P'
				}

				positionStr.WriteRune(pieceChar)
			}
		}

		if emptySquares > 0 {
			positionStr.WriteString(strconv.Itoa(emptySquares))
			emptySquares = 0
		}

		positionStr.WriteString("/")
	}

	sideToMove := ""
	castlingRights := ""
	epSquare := ""

	if pos.SideToMove == White {
		sideToMove = "w"
	} else {
		sideToMove = "b"
	}

	if pos.CastlingRights&WhiteKingsideRight != 0 {
		castlingRights += "K"
	}
	if pos.CastlingRights&WhiteQueensideRight != 0 {
		castlingRights += "Q"
	}
	if pos.CastlingRights&BlackKingsideRight != 0 {
		castlingRights += "k"
	}
	if pos.CastlingRights&BlackQueensideRight != 0 {
		castlingRights += "q"
	}

	if castlingRights == "" {
		castlingRights = "-"
	}

	if pos.EPSq == NoSq {
		epSquare = "-"
	} else {
		epSquare = sqToCoord(pos.EPSq)
	}

	// Technically this FEN isn't accurate as the full move counter
	// isn't tracked by Blunder's internal position object, so it's
	// always zero. But for our purposes this is good enough.
	return fmt.Sprintf(
		"%s %s %s %s %d %d",
		strings.TrimSuffix(positionStr.String(), "/"),
		sideToMove, castlingRights, epSquare,
		pos.HalfMoveClock, 0,
	)
}

func (pos *Position) StmInCheck() bool {
	kingSq := bitScan(pos.Pieces[King] & pos.Sides[pos.SideToMove])
	return sqIsAttacked(pos, pos.SideToMove, kingSq)
}

func (pos *Position) ComputePinAndCheckInfo() {
	kingSq := bitScan(pos.Pieces[King] & pos.Sides[pos.SideToMove])
	pos.pinned = pos.computePinnedPieces(kingSq, pos.SideToMove)
	pos.InCheck = sqIsAttacked(pos, pos.SideToMove, kingSq)
}

func (pos *Position) NoMajorsOrMiniors() bool {
	knights := bits.OnesCount64(pos.Pieces[Knight])
	bishops := bits.OnesCount64(pos.Pieces[Bishop])
	rook := bits.OnesCount64(pos.Pieces[Rook])
	queen := bits.OnesCount64(pos.Pieces[Queen])
	return knights+bishops+rook+queen == 0
}

func (pos *Position) DoMove(move uint32) bool {
	from := fromSq(move)
	to := toSq(move)
	moveType := moveType(move)
	flag := flag(move)

	state := &pos.prevStates[pos.StateIdx]
	state.CastlingRights = pos.CastlingRights
	state.EPSq = pos.EPSq
	state.Hash = pos.Hash
	state.HalfMoveClock = pos.HalfMoveClock

	state.MovedType = pos.GetPieceType(from)
	state.MovedColor = pos.GetPieceColor(from)
	state.CapturedType = pos.GetPieceType(to)
	state.CapturedColor = pos.GetPieceColor(to)

	state.Pinned = pos.pinned
	state.InCheck = pos.InCheck

	state.Phase = pos.Phase
	state.MGScores = pos.MGScores
	state.EGScores = pos.EGScores

	pos.StateIdx++

	pos.HalfMoveClock++

	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.EPSq = NoSq

	switch moveType {
	case Quiet:
		pos.zobristClearPiece(state.MovedType, state.MovedColor, from)
		pos.zobristPutPiece(state.MovedType, state.MovedColor, to)

		if state.MovedType == Pawn {
			pos.HalfMoveClock = 0
			if Abs(int8(from)-int8(to)) == 16 {
				pos.EPSq = to - 8
				if pos.SideToMove == Black {
					pos.EPSq = to + 8
				}
			}
		}
	case Attack:
		if flag == AttackEP {
			captureSq := to - 8
			if pos.SideToMove == Black {
				captureSq = to + 8
			}

			state.CapturedType = Pawn
			state.CapturedColor = pos.GetPieceColor(captureSq)

			pos.zobristClearPiece(Pawn, state.CapturedColor, captureSq)
			pos.zobristClearPiece(Pawn, state.MovedColor, from)
			pos.zobristPutPiece(Pawn, state.MovedColor, to)
		} else {
			pos.zobristClearPiece(state.CapturedType, state.CapturedColor, to)
			pos.zobristClearPiece(state.MovedType, state.MovedColor, from)
			pos.zobristPutPiece(state.MovedType, state.MovedColor, to)
		}

		pos.HalfMoveClock = 0
	case Promotion:
		pos.zobristClearPiece(Pawn, state.MovedColor, from)

		if state.CapturedType != NoType {
			pos.zobristClearPiece(state.CapturedType, state.CapturedColor, to)
		}

		pos.zobristPutPiece(flag+1, state.MovedColor, to)
		pos.HalfMoveClock = 0
	case Castle:
		castlingInfo := ToSqToCastlingInfo[to]

		pos.zobristClearPiece(King, state.MovedColor, from)
		pos.zobristPutPiece(King, state.MovedColor, to)

		pos.zobristClearPiece(Rook, state.MovedColor, castlingInfo.RookFromSq)
		pos.zobristPutPiece(Rook, state.MovedColor, castlingInfo.RookToSq)

		// Make sure to only test if the castling move was illegal
		// and return early AFTER the move has been made. Also
		// make sure to flip the side to move. All of this is done
		// to make sure we're not leaving the board is a
		// corrupted state that pos.UndoMove can't handle
		// (not doing these things led to nasty bugs!)
		if sqIsAttacked(pos, state.MovedColor, from) ||
			sqIsAttacked(pos, state.MovedColor, castlingInfo.FirstSqCrossed) ||
			sqIsAttacked(pos, state.MovedColor, castlingInfo.SecondSqCrossed) {

			pos.SideToMove ^= 1
			return false
		}
	}

	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)

	pos.CastlingRights = pos.CastlingRights & Spoilers[from] & Spoilers[to]
	pos.SideToMove ^= 1

	pos.Hash ^= Zobrist.CastlingNumber(pos.CastlingRights)
	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.Hash ^= Zobrist.SideToMoveNumber()

	kingSq := bitScan(pos.Pieces[King] & pos.Sides[pos.SideToMove^1])

	if pos.InCheck ||
		(state.MovedType == King && moveType != Castle) ||
		(moveType == Attack && flag == AttackEP) ||
		bitIsSet(pos.pinned, from) {
		return !sqIsAttacked(pos, pos.SideToMove^1, kingSq)
	}

	return true
}

func (pos *Position) UndoMove(move uint32) {
	pos.StateIdx--
	state := &pos.prevStates[pos.StateIdx]

	pos.CastlingRights = state.CastlingRights
	pos.EPSq = state.EPSq
	pos.HalfMoveClock = state.HalfMoveClock
	pos.Hash = state.Hash
	pos.pinned = state.Pinned
	pos.InCheck = state.InCheck
	pos.Phase = state.Phase
	pos.MGScores = state.MGScores
	pos.EGScores = state.EGScores

	pos.SideToMove ^= 1

	from := fromSq(move)
	to := toSq(move)
	moveType := moveType(move)
	flag := flag(move)

	switch moveType {
	case Quiet:
		pos.putPiece(state.MovedType, state.MovedColor, from)
		pos.clearPiece(state.MovedType, state.MovedColor, to)
	case Attack:
		if flag == AttackEP {
			capSq := to - 8
			if state.MovedColor == Black {
				capSq = to + 8
			}

			pos.putPiece(Pawn, state.MovedColor, from)
			pos.clearPiece(Pawn, state.MovedColor, to)
			pos.putPiece(Pawn, state.CapturedColor, capSq)
		} else {
			pos.putPiece(state.MovedType, state.MovedColor, from)
			pos.clearPiece(state.MovedType, state.MovedColor, to)
			pos.putPiece(state.CapturedType, state.CapturedColor, to)
		}
	case Castle:
		pos.putPiece(King, state.MovedColor, from)
		pos.clearPiece(King, state.MovedColor, to)

		castlingInfo := ToSqToCastlingInfo[to]
		pos.clearPiece(Rook, state.MovedColor, castlingInfo.RookToSq)
		pos.putPiece(Rook, state.MovedColor, castlingInfo.RookFromSq)
	case Promotion:
		pos.putPiece(Pawn, state.MovedColor, from)
		pos.clearPiece(flag+1, state.MovedColor, to)
		if state.CapturedType != NoType {
			pos.putPiece(state.CapturedType, state.CapturedColor, to)
		}
	}
}

func (pos *Position) DoNullMove() {
	state := &pos.prevStates[pos.StateIdx]
	state.Hash = pos.Hash
	state.CastlingRights = pos.CastlingRights
	state.EPSq = pos.EPSq
	state.HalfMoveClock = pos.HalfMoveClock
	state.Pinned = pos.pinned
	state.InCheck = pos.InCheck
	state.Phase = pos.Phase
	state.MGScores = pos.MGScores
	state.EGScores = pos.EGScores

	pos.StateIdx++

	pos.Hash ^= Zobrist.EPNumber(pos.EPSq)
	pos.EPSq = NoSq

	pos.HalfMoveClock = 0

	pos.SideToMove ^= 1
	pos.Hash ^= Zobrist.SideToMoveNumber()
}

func (pos *Position) UndoNullMove() {
	pos.StateIdx--
	state := &pos.prevStates[pos.StateIdx]

	pos.Hash = state.Hash
	pos.CastlingRights = state.CastlingRights
	pos.EPSq = state.EPSq
	pos.HalfMoveClock = state.HalfMoveClock
	pos.pinned = state.Pinned
	pos.InCheck = state.InCheck
	pos.Phase = state.Phase
	pos.MGScores = state.MGScores
	pos.EGScores = state.EGScores

	pos.SideToMove ^= 1
}

func (pos *Position) computePinnedPieces(kingSq, usColor uint8) uint64 {
	usBB := pos.Sides[usColor]
	enemyBB := pos.Sides[usColor^1]

	enemyQueensAndRooks := enemyBB & (pos.Pieces[Queen] | pos.Pieces[Rook])
	enemyQueensAndBishops := enemyBB & (pos.Pieces[Queen] | pos.Pieces[Bishop])

	cardinalRays := genRookMovesBB(kingSq, enemyBB)
	intercardinalRays := genBishopMovesBB(kingSq, enemyBB)

	potentialPinners := (enemyQueensAndRooks & cardinalRays) | (enemyQueensAndBishops & intercardinalRays)
	pinned := EmptyBB

	for potentialPinners != 0 {
		pinnerSq := BitScanAndClear(&potentialPinners)
		pinned |= usBB & (RaysBetween[kingSq][pinnerSq])
	}

	return pinned
}

func (pos *Position) GetPieceType(sq uint8) uint8 {
	sqBB := MostSigBitBB >> uint64(sq)
	shift := uint64(63 - sq)

	kingBB := (sqBB & pos.Pieces[King]) >> shift
	queenBB := (sqBB & pos.Pieces[Queen]) >> shift
	rookBB := (sqBB & pos.Pieces[Rook]) >> shift
	bishopBB := (sqBB & pos.Pieces[Bishop]) >> shift
	knightBB := (sqBB & pos.Pieces[Knight]) >> shift
	pawnBB := (sqBB & pos.Pieces[Pawn]) >> shift

	index := kingBB | queenBB<<1 | rookBB<<2 | bishopBB<<3 | knightBB<<4 | pawnBB<<5
	return PieceTypeLookup[index]
}

func (pos *Position) GetPieceColor(sq uint8) uint8 {
	sqBB := MostSigBitBB >> uint64(sq)
	shift := uint64(63 - sq)

	whiteBB := (sqBB & pos.Sides[White]) >> shift
	blackBB := (sqBB & pos.Sides[Black]) >> shift

	index := whiteBB | blackBB<<1
	return PieceColorLookup[index]
}

func (pos *Position) zobristPutPiece(pieceType, pieceColor, sq uint8) {
	setBit(&pos.Pieces[pieceType], sq)
	setBit(&pos.Sides[pieceColor], sq)
	pos.Hash ^= Zobrist.PieceNumber(pieceType, pieceColor, sq)

	colorCorrectedSq := FlipSq[pieceColor][sq]
	pos.Phase += PhaseScores[pieceType]
	pos.MGScores[pieceColor] += MG_PIECE_VALUES[pieceType] + MG_PSQT[pieceType][colorCorrectedSq]
	pos.EGScores[pieceColor] += EG_PIECE_VALUES[pieceType] + EG_PSQT[pieceType][colorCorrectedSq]
}

func (pos *Position) zobristClearPiece(pieceType, pieceColor, sq uint8) {
	clearBit(&pos.Pieces[pieceType], sq)
	clearBit(&pos.Sides[pieceColor], sq)
	pos.Hash ^= Zobrist.PieceNumber(pieceType, pieceColor, sq)

	colorCorrectedSq := FlipSq[pieceColor][sq]
	pos.Phase -= PhaseScores[pieceType]
	pos.MGScores[pieceColor] -= MG_PIECE_VALUES[pieceType] + MG_PSQT[pieceType][colorCorrectedSq]
	pos.EGScores[pieceColor] -= EG_PIECE_VALUES[pieceType] + EG_PSQT[pieceType][colorCorrectedSq]
}

func (pos *Position) putPiece(pieceType, pieceColor, sq uint8) {
	setBit(&pos.Pieces[pieceType], sq)
	setBit(&pos.Sides[pieceColor], sq)
}

func (pos *Position) clearPiece(pieceType, pieceColor, sq uint8) {
	clearBit(&pos.Pieces[pieceType], sq)
	clearBit(&pos.Sides[pieceColor], sq)
}

func (pos Position) String() (boardStr string) {
	boardStr += "\n"

	for i := 56; i >= 0; i -= 8 {
		boardStr += fmt.Sprintf("%v | ", (i/8)+1)
		for j := i; j < i+8; j++ {
			pieceType := pos.GetPieceType(uint8(j))
			pieceColor := pos.GetPieceColor(uint8(j))

			pieceChar := PieceTypeToChar[pieceType]
			if pieceColor == White {
				pieceChar = unicode.ToUpper(pieceChar)
			}

			boardStr += fmt.Sprintf("%c ", pieceChar)
		}
		boardStr += "\n"
	}

	boardStr += "   ----------------"
	boardStr += "\n    a b c d e f g h"

	boardStr += "\n\n"
	if pos.SideToMove == White {
		boardStr += "turn: white\n"
	} else {
		boardStr += "turn: black\n"
	}

	boardStr += "castling rights: "
	if pos.CastlingRights&WhiteKingsideRight != 0 {
		boardStr += "K"
	}
	if pos.CastlingRights&WhiteQueensideRight != 0 {
		boardStr += "Q"
	}
	if pos.CastlingRights&BlackKingsideRight != 0 {
		boardStr += "k"
	}
	if pos.CastlingRights&BlackQueensideRight != 0 {
		boardStr += "q"
	}

	boardStr += "\nen passant: "
	if pos.EPSq == NoSq {
		boardStr += "none"
	} else {
		boardStr += sqToCoord(pos.EPSq)
	}

	boardStr += fmt.Sprintf("\nhalf-move clock: %d\n", pos.HalfMoveClock)
	boardStr += fmt.Sprintf("phase: %d\n", pos.Phase)
	return boardStr
}
