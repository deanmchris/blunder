package neural_network

/*import (
	"blunder/engine"
	"blunder/tuner"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// gen_data.go generates training data for training the neural network.

var OutcomeToResult = []string{
	"1.0",
	"0.0",
	"0.5",
}

// Given an infile containg the PGNs, extract positions from the files,
// and write them to the given outfile.
func GenTrainingData(infile, outfile string, samplingSizePerGame int, minElo uint16, maxGames uint32) {
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

	rand.Seed(time.Now().UnixNano())

	pgns := tuner.ParsePGNs(infile, minElo, maxGames)
	numPositions := 0
	fens := []string{}

	for i, pgn := range pgns {
		log.Printf("Extracting positions from game %d\n", i+1)

		search.Pos.LoadFEN(pgn.Fen)
		gamePly := len(pgn.Moves)
		possibleFens := []string{}

		if gamePly <= 40 {
			continue
		}

		for _, move := range pgn.Moves {
			search.Pos.DoMove(move)
			search.Pos.StatePly--

			if search.Pos.InCheck() {
				continue
			}

			if search.Pos.Ply > 200 {
				continue
			}

			if search.Pos.Ply < 10 {
				continue
			}

			if gamePly-int(search.Pos.Ply) <= 10 {
				continue
			}

			// pvLine := engine.PVLine{}
			//search.Qsearch(-engine.Inf, engine.Inf, 0, &pvLine, 0)

			fields := strings.Fields(search.Pos.GenFEN())

			// result := OutcomeToResult[pgn.Outcome]

			eval := float64(engine.EvaluatePos(&search.Pos)) / 100
			if search.Pos.SideToMove == engine.Black {
				eval = -eval
			}

			possibleFens = append(possibleFens, fmt.Sprintf("%s %s - - 0 1|%f\n", fields[0], fields[1], eval))
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
}*/
