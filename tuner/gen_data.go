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
		for moveNum, move := range pgn.Moves {
			search.Pos.MakeMove(move)
			search.Pos.StatePly--

			// Skip the last few moves.
			if (len(pgn.Moves)-1)-moveNum <= 12 {
				continue
			}

			// Ensure the position is quiet by doing a qsearch and using the resulting
			// PV line to construct the final position.
			var pvLine engine.PVLine
			search.Qsearch(-engine.Inf, engine.Inf, 0, 0, &pvLine)

			fields := strings.Fields(getPVPosition(&pvLine, &search.Pos))
			result := OutcomeToResult[pgn.Outcome]

			fens = append(fens, fmt.Sprintf("%s %s %s %s c9 \"%s\"\n", fields[0], fields[1], fields[2], fields[3], result))
			numPositions++
		}
	}

	// randomize positions to avoid overfitting when training
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(fens), func(i, j int) { fens[i], fens[j] = fens[j], fens[i] })

	seen := make(map[string]int)
	duplicates := 0

	for _, fen := range fens {
		seen[fen]++
		if seen[fen] == 1 {
			_, err := file.WriteString(fen)
			if err != nil {
				panic(err)
			}
		} else {
			duplicates++
		}
	}

	fmt.Printf("%d positions succesfully extracted!\n", numPositions)
	fmt.Printf("%d duplicate positions were skipped!\n", duplicates)
	file.Close()
}

func getPVPosition(pvLine *engine.PVLine, pos *engine.Position) (fen string) {
	for _, move := range pvLine.Moves {
		pos.MakeMove(move)
	}

	fen = pos.GenFEN()

	for i := len(pvLine.Moves) - 1; i >= 0; i-- {
		pos.UnmakeMove(pvLine.Moves[i])
	}

	return
}
