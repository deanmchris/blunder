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

var OutcomeToResult []string = []string{
	"1-0",
	"0-1",
	"1/2-1/2",
}

// Given an infile containg the PGNs, extract quiet positions from the files,
// and write them to the given outfile.
func GenTrainingData(infile, outfile string) {
	file, err := os.OpenFile(outfile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	var search engine.Search
	search.Timer.TimeLeft = engine.InfiniteTime
	search.SpecifiedDepth = engine.MaxPly
	search.SpecifiedNodes = math.MaxUint64
	search.TT.Resize(engine.DefaultTTSize)

	pgns := parsePGNs(infile)
	numPositions := 0
	fens := []string{}

	for i, pgn := range pgns {
		fmt.Printf("Extracting positions from game %d\n", i+1)

		search.Pos.LoadFEN(pgn.Fen)
		for j, move := range pgn.Moves {
			search.Pos.MakeMove(move)
			search.Pos.StatePly--

			eval := engine.EvaluatePos(&search.Pos)
			qeval := search.Qsearch(-engine.Inf, engine.Inf, 0, &engine.PVLine{})

			if search.Pos.InCheck() {
				continue
			}

			if (len(pgn.Moves) - j) <= 10 {
				continue
			}

			if abs16(qeval-eval) > 50 {
				continue
			}

			fields := strings.Fields(search.Pos.GenFEN())
			result := OutcomeToResult[pgn.Outcome]

			fens = append(fens, fmt.Sprintf("%s %s %s %s c9 \"%s\";\n", fields[0], fields[1], fields[2], fields[3], result))
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

func abs16(n int16) int16 {
	if n < 0 {
		return -n
	}
	return n
}
