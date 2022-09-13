package tuner

import (
	"blunder/engine"
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// pgn.go implements a (purposely) partial pgn parser to be used in generating tuning data for Blunder.

const (
	WhiteWon uint8 = 0
	BlackWon uint8 = 1
	Drawn    uint8 = 2

	EventTagPattern    = "\\[Event .*\\]"
	VariantTagPattern  = "\\[Variant .*\\]"
	WhiteEloTagPattern = "\\[WhiteElo .*\\]"
	BlackEloTagPattern = "\\[BlackElo .*\\]"
	FenTagPattern      = "\\[FEN .*\\]"
	ResultTagPattern   = "\\[Result .*\\]"
	TagPattern         = "\\[.*\\]"
	MovePattern        = "([a-h]x)?([a-h][81]=[QRBN])|(([QRBNK])?([a-h])?([1-8])?(x)?[a-h][1-8])|(O-O-O)|(O-O)"

	BufferSize = 1024 * 500
)

type PGN struct {
	Outcome uint8
	Fen     string
	Moves   []engine.Move
}

// Parse a file of PGNs. The file name is assumed to be the name of a
// file in the blunder/tuner directory.
func ParsePGNs(filename string, minElo uint16, maxGames uint32) (pgns []PGN) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var pos engine.Position
	fenTagRegex, _ := regexp.Compile(FenTagPattern)
	variantTagRegex, _ := regexp.Compile(VariantTagPattern)
	whiteEloTagRegex, _ := regexp.Compile(WhiteEloTagPattern)
	blackEloTagRegex, _ := regexp.Compile(BlackEloTagPattern)
	resultTagRegex, _ := regexp.Compile(ResultTagPattern)
	tagRegex, _ := regexp.Compile(TagPattern)
	moveRegex, _ := regexp.Compile(MovePattern)

	reader := bufio.NewReader(file)
	gamesParsed := uint32(0)
	head := ""

mainLoop:
	for {
		buffer := make([]byte, BufferSize)
		n, err := reader.Read(buffer)

		if err != io.EOF && err != nil {
			panic(err)
		}

		buffer = buffer[:n]
		rawPGNs := head + string(buffer)

		eventTagRegex, _ := regexp.Compile(EventTagPattern)
		allIndexes := eventTagRegex.FindAllStringSubmatchIndex(rawPGNs, -1)

		chunks := []string{}
		for i := 0; i < len(allIndexes); i++ {
			chunkStart := allIndexes[i][0]
			chunkEnd := 0

			if i+1 == len(allIndexes) {
				// Since we're reading a fixed buffer size, we're probably going to run into
				// the issue of spliting a PGN between buffer reads. This is remedied by always
				// taking the second to last chunk in chunks, removing it, and adding it to the
				// start of the next buffer we read. This way we're always guaranteed to get a set
				// of full PGNs, but the potentially broken PGN is not discarded and will simply be
				// processed on the next iteration of the loop.
				chunkEnd = len(rawPGNs)
				head = rawPGNs[chunkStart:chunkEnd]
			} else {
				chunkEnd = allIndexes[i+1][0]
				chunks = append(chunks, strings.TrimSpace(rawPGNs[chunkStart:chunkEnd]))
			}
		}

		for _, chunk := range chunks {
			var fen string
			var outcome uint8
			var moves []engine.Move

			// Variants are skipped
			variantTag := variantTagRegex.FindString(chunk)
			if variantTag != "" {
				continue
			}

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
				continue
			} else {
				result := strings.Fields(resultTag)[1]
				result = strings.TrimPrefix(result, "\"")
				result = strings.TrimSuffix(result, "\"]")

				if result == "1-0" {
					outcome = WhiteWon
				} else if result == "0-1" {
					outcome = BlackWon
				} else if result == "1/2-1/2" {
					outcome = Drawn
				} else {
					continue
				}
			}

			whiteEloTag := whiteEloTagRegex.FindString(chunk)
			if whiteEloTag == "" {
				if minElo > 0 {
					continue
				}
			} else {
				elo := strings.Fields(whiteEloTag)[1]
				elo = strings.TrimPrefix(elo, "\"")
				elo = strings.TrimSuffix(elo, "\"]")

				eloNumber, err := strconv.Atoi(elo)
				if err != nil {
					continue
				}

				if uint16(eloNumber) < minElo {
					continue
				}
			}

			blackEloTag := blackEloTagRegex.FindString(chunk)
			if blackEloTag == "" {
				if minElo > 0 {
					continue
				}
			} else {
				elo := strings.Fields(blackEloTag)[1]
				elo = strings.TrimPrefix(elo, "\"")
				elo = strings.TrimSuffix(elo, "\"]")

				eloNumber, err := strconv.Atoi(elo)
				if err != nil {
					continue
				}

				if uint16(eloNumber) < minElo {
					continue
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
			gamesParsed++

			if gamesParsed == maxGames {
				break mainLoop
			}
		}

		if err == io.EOF {
			break
		}
	}

	return pgns
}
