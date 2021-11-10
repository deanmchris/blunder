package engine

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// book.go is an implementation of a polyglot opening book prober for Blunder.

const (
	// The size of a polyglot entry
	EntryByteLength = 16
	// Each masks helps to extract the correct bits
	// from the move part of a polyglot entry.
	ToFileMask         uint16 = 0x7
	ToRankMask         uint16 = 0x38
	FromFileMask       uint16 = 0x1C0
	FromRankMask       uint16 = 0xE00
	PromotionPieceMask uint16 = 0x7000

	// Each shift works together with a move
	// mask above to extract the correct number
	// from the move part of a polyglot entry
	// by shifting the masked bits to the least
	// significant end of a bitstring.
	ToRankShift         = 3
	FromFileShift       = 6
	FromRankShift       = 9
	PromotionPieceShift = 12

	// Characters representing the eight files
	// and ranks
	fileCharacters = "abcdefgh"
	rankCharacters = "12345678"
)

// Each polyglot book is composed of a series of 16-byte entries. Each
// of these entries contains a key, which is the hash of the position
// after the current moves have been made, the moves made, the weight
// those moves are given (i.e. how good they are), and a learn field,
// which, as far as I can tell, is usually ignored and set to zero
// by polyglot book generators, so it's not included here. The key
// element of the entry is the mapping key to a PolyglotEntry.
type PolyglotEntry struct {
	Hash   uint64
	Move   string
	Weight uint16
}

// Parse a polyglot file and create a map of PolyglotEntry's
// from it. Each zobrist hash for an entry maps to the moves
// and weight of the entry.
func LoadPolyglotFile(path string) (map[uint64][]PolyglotEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	entries := make(map[uint64][]PolyglotEntry)

	for {
		var entryBytes [EntryByteLength]byte
		_, err := io.ReadFull(reader, entryBytes[:])
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		var entry PolyglotEntry

		// Load polyglot hash
		bytesBuffer := bytes.NewBuffer(entryBytes[0:8])
		binary.Read(bytesBuffer, binary.BigEndian, &entry.Hash)

		// Load the move
		var move uint16
		bytesBuffer.Reset()
		bytesBuffer.Write(entryBytes[8:10])
		binary.Read(bytesBuffer, binary.BigEndian, &move)

		toFile := fileCharacters[move&ToFileMask]
		toRank := rankCharacters[(move&ToRankMask)>>ToRankShift]
		fromFile := fileCharacters[(move&FromFileMask)>>FromFileShift]
		fromRank := rankCharacters[(move&FromRankMask)>>FromRankShift]
		promotionPiece := (move & PromotionPieceMask) >> PromotionPieceShift

		promotionCharacter := ""
		switch promotionPiece {
		case 1:
			promotionCharacter = "n"
		case 2:
			promotionCharacter = "b"
		case 3:
			promotionCharacter = "r"
		case 4:
			promotionCharacter = "q"
		}

		entry.Move = fmt.Sprintf("%c%c%c%c%v", fromFile, fromRank, toFile, toRank, promotionCharacter)

		// Load the weight
		var weight uint16
		bytesBuffer.Reset()
		bytesBuffer.Write(entryBytes[10:12])
		binary.Read(bytesBuffer, binary.BigEndian, &weight)
		entry.Weight = weight

		// Load the learn data, and discard it
		var learn uint32
		bytesBuffer.Reset()
		bytesBuffer.Write(entryBytes[12:16])
		binary.Read(bytesBuffer, binary.BigEndian, &learn)
		entries[entry.Hash] = append(entries[entry.Hash], entry)
	}
	return entries, nil
}
