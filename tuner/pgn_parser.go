package tuner

import (
	"blunder/engine"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// pgn.go implements a (purposely) partial pgn parser to be used in generating tuning data for Blunder.

const (
	WhiteWon uint8 = 0
	BlackWon uint8 = 1
	Drawn    uint8 = 2

	EventTagPattern  = "\\[Event .*\\]"
	FenTagPattern    = "\\[FEN .*\\]"
	ResultTagPattern = "\\[Result .*\\]"
	TagPattern       = "\\[.*\\]"
	MovePattern      = "([a-h]x)?([a-h][81]=[QRBN])|(([QRBNK])?([a-h])?([1-8])?(x)?[a-h][1-8])|(O-O-O)|(O-O)"
)

type PGN struct {
	Outcome uint8
	Fen     string
	Moves   []engine.Move
}

func parsePGNs(filename string) (pgns []PGN) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	file := string(buf)

	eventTagRegex, _ := regexp.Compile(EventTagPattern)
	allIndexes := eventTagRegex.FindAllStringSubmatchIndex(file, -1)

	chunks := []string{}
	for i := 0; i < len(allIndexes); i++ {
		chunkStart := allIndexes[i][0]
		chunkEnd := 0

		if i+1 == len(allIndexes) {
			chunkEnd = len(file)
		} else {
			chunkEnd = allIndexes[i+1][0]
		}

		chunks = append(chunks, strings.TrimSpace(file[chunkStart:chunkEnd]))
	}

	var pos engine.Position
	fenTagRegex, _ := regexp.Compile(FenTagPattern)
	resultTagRegex, _ := regexp.Compile(ResultTagPattern)
	tagRegex, _ := regexp.Compile(TagPattern)
	moveRegex, _ := regexp.Compile(MovePattern)

	for _, chunk := range chunks {
		var fen string
		var outcome uint8
		var moves []engine.Move

		fenTag := fenTagRegex.FindString(chunk)
		if fenTag == "" {
			fen = engine.FENStartPosition
		} else {
			fields := strings.Fields(fenTag)
			fen = fields[1] + " " + fields[2] + " " + fields[3] + " "
			fen += fields[4] + " " + fields[5] + " " + fields[6]
			fen = strings.TrimPrefix(fen, "\"")
			fen = strings.TrimSuffix(fen, "\"]")
		}

		resultTag := resultTagRegex.FindString(chunk)
		if resultTag == "" {
			log.Println("invalid pgn skipped")
			continue
		} else {
			result := strings.Fields(resultTag)[1]
			result = strings.TrimPrefix(result, "\"")
			result = strings.TrimSuffix(result, "\"]")

			switch result {
			case "1-0":
				outcome = WhiteWon
			case "0-1":
				outcome = BlackWon
			default:
				outcome = Drawn
			}
		}

		// Once we've read the tags we care about, all tags should be removed
		// from this "chunk". Do this by getting the end index of the last tag
		// in the chunk, and cutting everything off in the string before that
		// index.
		tagIndexes := tagRegex.FindAllStringSubmatchIndex(chunk, -1)
		lastTagIndex := tagIndexes[len(tagIndexes)-1][1]
		chunk = chunk[lastTagIndex:]
		chunk = strings.TrimSpace(chunk)

		moveStrings := moveRegex.FindAllString(chunk, -1)
		pos.LoadFEN(fen)

		for _, moveStr := range moveStrings {
			move := engine.ConvertSANToLAN(&pos, moveStr)
			pos.DoMove(move)
			pos.StatePly--
			moves = append(moves, move)
		}

		pgns = append(pgns, PGN{Fen: fen, Outcome: outcome, Moves: moves})
	}

	return pgns
}
