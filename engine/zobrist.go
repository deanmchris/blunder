package engine

// zobrist.go contains an interface for creating and incrementally updating
// the zobrist hash value of a given position.
//
// https://www.chessprogramming.org/Zobrist_Hashing

const (
	// A constant representing the value to use when seeding the random
	// numbers generated for zobrist hashing.
	ZobristSeedValue = 1

	// A constant which represents when there is no ep square. This value indexes
	// into _Zobrist.epFileRand64 to return a 0, which will not affect the zobrist
	// hash.
	NoEPFile = 8
)

var Zobrist _Zobrist

type _Zobrist struct {
	// Each aspect of the board needs to be a given a unique random 64-bit number
	// that will be xor-ed together with other unique random numbers from the positions
	// other aspects to create a unique zobrist hash.
	//
	// Blunder uses 794 unique random numbers:
	// * 12x64 for each type of piece on each possible square.
	// * 8 for each possible en passant file
	// * 16 for each possible permutation of the castling right bits.
	// * 1 for when it's white to move

	pieceSqRand64        [768]uint64
	epFileRand64         [9]uint64
	castlingRightsRand64 [16]uint64
	sideToMoveRand64     uint64
}

func (zobrist *_Zobrist) init() {
	prng := PseduoRandomGenerator{}
	prng.Seed(ZobristSeedValue)

	for index := 0; index < 768; index++ {
		zobrist.pieceSqRand64[index] = prng.Random64()
	}

	for index := 0; index < 8; index++ {
		zobrist.epFileRand64[index] = prng.Random64()
	}

	for index := 0; index < 16; index++ {
		zobrist.castlingRightsRand64[index] = prng.Random64()
	}

	zobrist.sideToMoveRand64 = prng.Random64()
}

func (zobrist *_Zobrist) PieceNumber(pieceType, pieceColor uint8, sq uint8) uint64 {
	return zobrist.pieceSqRand64[(uint16(pieceType)*2+uint16(pieceColor))*64+uint16(sq)]
}

func (zobrist *_Zobrist) EPNumber(epSq uint8) uint64 {
	return zobrist.epFileRand64[PossibleEPFiles[epSq]]
}

func (zobrist *_Zobrist) CastlingNumber(castlingRights uint8) uint64 {
	return zobrist.castlingRightsRand64[castlingRights]
}

func (zobrist *_Zobrist) SideToMoveNumber() uint64 {
	return zobrist.sideToMoveRand64
}

func (zobrist *_Zobrist) GenHash(pos *Position) (hash uint64) {
	for sq := uint8(0); sq < 64; sq++ {
		pieceType := pos.GetPieceType(sq)
		pieceColor := pos.GetPieceColor(sq)

		if pieceType != NoType {
			hash ^= zobrist.PieceNumber(pieceType, pieceColor, uint8(sq))
		}
	}

	hash ^= zobrist.EPNumber(pos.EPSq)
	hash ^= zobrist.CastlingNumber(pos.CastlingRights)

	if pos.SideToMove == White {
		hash ^= zobrist.SideToMoveNumber()
	}

	return hash
}

// Precomputing all possible en passant file numbers
// is much more efficent for Blunder than calculating
// them on the fly.
var PossibleEPFiles [65]uint8 = [65]uint8{
	8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8,
	0, 1, 2, 3, 4, 5, 6, 7,
	8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8,
	0, 1, 2, 3, 4, 5, 6, 7,
	8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8,
	8,
}

func InitZobrist() {
	Zobrist = _Zobrist{}
	Zobrist.init()
}
