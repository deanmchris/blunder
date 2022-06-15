package tuner

import (
	"blunder/engine"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// gen_data.go generates training data for the texel tuner from the PGNs of games played.

var OutcomeToResult = [3]float64{
	1.0,
	0.0,
	0.5,
}

// Given an infile containg the PGNs, extract quiet positions from the files,
// and write them to the given outfile.
func GenTrainingData(infile, outfile string, minimumElo, minimumYear int) {
	file, err := os.OpenFile(outfile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	search := engine.Search{}
	search.TT.Resize(engine.DefaultTTSize, engine.SearchEntrySize)
	search.Timer.Setup(
		engine.InfiniteTime,
		engine.NoValue,
		engine.NoValue,
		int16(engine.NoValue),
		engine.MaxDepth,
		math.MaxUint64,
	)

	pgns, skipped := parsePGNs(infile, minimumElo, minimumYear)
	numPositions := 0
	fens := []string{}

	fmt.Printf("Skipped %d unwanted/invalid games", skipped)

	for i, pgn := range pgns {
		fmt.Printf("Extracting positions from game %d\n", i+1)

		search.Pos.LoadFEN(pgn.Fen)
		for moveCnt, move := range pgn.Moves {
			search.Pos.DoMove(move)
			search.Pos.StatePly--

			pvLine := engine.PVLine{}
			search.Qsearch(-engine.Inf, engine.Inf, 0, &pvLine)

			if search.Pos.InCheck() {
				continue
			}

			if (len(pgn.Moves) - moveCnt) <= 10 {
				continue
			}

			fields := strings.Fields(GetFenFromPV(pvLine, &search.Pos))
			result := OutcomeToResult[pgn.Outcome]

			fens = append(fens, fmt.Sprintf("%s %s %s %s 0 1 %f\n", fields[0], fields[1], fields[2], fields[3], result))
			numPositions++
		}
	}

	// randomize positions to avoid overfitting when training
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(fens), func(i, j int) { fens[i], fens[j] = fens[j], fens[i] })

	for _, fen := range fens {
		_, err := file.WriteString(fen)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("%d positions succesfully extracted!\n", numPositions)
	file.Close()
}

func GetFenFromPV(pvLine engine.PVLine, pos *engine.Position) (fen string) {
	for _, move := range pvLine.Moves {
		pos.DoMove(move)
	}

	fen = pos.GenFEN()

	for i := len(pvLine.Moves) - 1; i >= 0; i-- {
		pos.UndoMove(pvLine.Moves[i])
	}

	return fen
}
