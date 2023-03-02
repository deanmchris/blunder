package tuner

import (
	"blunder/engine"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// gen_data.go generates training data for the texel tuner from the PGNs of games played.

var OutcomeToResult = []string{
	"1.0",
	"0.0",
	"0.5",
}

// Given an infile containg the PGNs, extract quiet positions from the files,
// and write them to the given outfile.
func GenTrainingData(infile, outfile string, samplingSizePerGame int, minElo uint16, maxGames uint32) {
	search := engine.Search{}
	search.Timer.Setup(
		engine.InfiniteTime,
		engine.NoValue,
		engine.NoValue,
		int16(engine.NoValue),
		engine.MaxPly,
		math.MaxUint64,
	)

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	numPositions := 0
	fens := []string{}
	pgnCnt := 1

	for pgn := range parsePGNFile(infile, minElo, maxGames) {
		log.Printf("Extracting positions from game %d\n", pgnCnt)
		pgnCnt++

		search.Pos.LoadFEN(pgn.Fen)
		fullGamePly := len(pgn.Moves)
		possibleFens := []string{}

		if fullGamePly <= 40 {
			continue
		}

		gamePly := 0
		for _, move := range pgn.Moves {
			search.Pos.ComputePinAndCheckInfo()
			search.Pos.DoMove(move)
			search.Pos.StateIdx--
			gamePly++

			if search.Pos.InCheck {
				continue
			}

			if gamePly > 200 {
				continue
			}

			if gamePly < 10 {
				continue
			}

			if fullGamePly-gamePly <= 10 {
				continue
			}

			pvLine := engine.PVLine{}
			search.QuiescenceSearch(-engine.Infinity, engine.Infinity, &pvLine)

			fields := strings.Fields(applyPV(&search.Pos, pvLine))
			result := OutcomeToResult[pgn.Outcome]

			beforeFen := search.Pos.GenFEN()
			fen := fmt.Sprintf("%s %s - - 0 1 [%s]\n", fields[0], fields[1], result)

			search.Pos.LoadFEN(fen)
			search.Pos.ComputePinAndCheckInfo()

			if search.Pos.InCheck {
				search.Pos.LoadFEN(beforeFen)
				continue
			} else {
				search.Pos.LoadFEN(beforeFen)
			}


			possibleFens = append(possibleFens, fen)
		}

		samplingSize := engine.Min(samplingSizePerGame, len(possibleFens))
		for i := 0; i < samplingSize; i++ {
			fens = append(fens, possibleFens[rand.Intn(len(possibleFens))])
			numPositions++
		}
	}

	rand.Shuffle(len(fens), func(i, j int) { fens[i], fens[j] = fens[j], fens[i] })

	file, err := os.OpenFile(outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	seen := make(map[string]bool)
	for _, fen := range fens {
		if seenBefore := seen[fen]; !seenBefore {
			seen[fen] = true
			_, err := file.WriteString(fen)
			if err != nil {
				panic(err)
			}
		} else {
			numPositions--
		}
	}

	log.Printf("%d positions succesfully extracted!\n", numPositions)
	file.Close()
}

func applyPV(pos *engine.Position, pvLine engine.PVLine) (fen string) {
	for _, move := range pvLine.Moves {
		pos.DoMove(move)
	}

	fen = pos.GenFEN()

	for i := len(pvLine.Moves) - 1; i >= 0; i-- {
		pos.UndoMove(pvLine.Moves[i])
	}

	return fen
}
