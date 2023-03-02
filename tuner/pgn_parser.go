package tuner

import (
	"blunder/engine"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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
	TagPattern         = "\\[[a-zA-Z]+ .*\\]"
	MovePattern        = "([a-h]x)?([a-h][81]=[QRBN])|(([QRBNK])?([a-h])?([1-8])?(x)?[a-h][1-8])|(O-O-O)|(O-O)"
)

type PGN struct {
	Outcome uint8
	Fen     string
	Moves   []uint32
}

func GenerateChunks(filename string) chan string {
	chunks := make(chan string)

	go func() {
		defer close(chunks)

		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(file)
		scanner := bufio.NewScanner(reader)

		chunk := strings.Builder{}
		eventTagsSeen := 0

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "Event") {
				eventTagsSeen++

				if eventTagsSeen > 1 {
					chunks <- strings.TrimSpace(chunk.String())
					chunk.Reset()
					chunk.WriteString(line)
					chunk.WriteString("\n")
				} else {
					chunk.WriteString(line)
					chunk.WriteString("\n")
				}
			} else {
				chunk.WriteString(line)
				chunk.WriteString("\n")
			}
		}

		chunks <- strings.TrimSpace(chunk.String())
	}()

	return chunks
}

// Parse a file of PGNs. The file name is assumed to be the name of a
// file in the blunder/tuner directory.
func parsePGNFile(filename string, minElo uint16, maxGames uint32) chan *PGN {
	pgns := make(chan *PGN)

	go func() {
		defer close(pgns)

		fenTagRegex, _ := regexp.Compile(FenTagPattern)
		variantTagRegex, _ := regexp.Compile(VariantTagPattern)
		whiteEloTagRegex, _ := regexp.Compile(WhiteEloTagPattern)
		blackEloTagRegex, _ := regexp.Compile(BlackEloTagPattern)
		resultTagRegex, _ := regexp.Compile(ResultTagPattern)
		tagRegex, _ := regexp.Compile(TagPattern)
		moveRegex, _ := regexp.Compile(MovePattern)

		var pos engine.Position
		var gamesParsed uint32

		for chunk := range GenerateChunks(filename) {
			var fen string
			var outcome uint8
			var moves []uint32

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

			fullChunk := chunk
			chunk = chunk[lastTagIndex:]
			chunk = strings.TrimSpace(chunk)

			moveStrings := moveRegex.FindAllString(chunk, -1)

			pos.LoadFEN(fen)

			for i, moveStr := range moveStrings {
				pos.ComputePinAndCheckInfo()
				move, err := engine.ConvertSANToLAN(&pos, moveStr)

				if err != nil {
					panic(fmt.Errorf("%s. For move %d. %s for chunk: \n%s", err.Error(), (i+1)/2, moveStr, fullChunk))
				}

				if move == engine.NullMove {
					panic(fmt.Errorf("could not convert move %d. %s for chunk: \n%s", (i+1)/2, moveStr, fullChunk))
				}

				pos.DoMove(move)
				pos.StateIdx--
				moves = append(moves, move)
			}

			pgns <- &PGN{Fen: fen, Outcome: outcome, Moves: moves}
			gamesParsed++

			if gamesParsed == maxGames {
				break
			}
		}
	}()

	return pgns
}
