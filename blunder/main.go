package main

import (
	"blunder/engine"
	"flag"
	"fmt"
	"os"
	"time"
	"blunder/tuner"
)

func init() {
	engine.InitBitboard()
	engine.InitTables()
	engine.InitMagics()
	engine.InitZobrist()
}

func main() {
	tuneCommand := flag.NewFlagSet("tune", flag.ExitOnError)
	tuneInputFile := tuneCommand.String(
		"input-file",
		"",
		"A file containing a set of fens to use for tuning. Should be in the format\n"+
			"<full fen> [<result float>], where 'result float' is either 1.0 (white won),\n"+
			"0.0 (black won), or 0.5 (draw).",
	)
	tuneEpochs := tuneCommand.Int("epochs", 50000, "The number of epochs to run the tuner for.")
	tuneLearningRate := tuneCommand.Float64("learning-rate", 1e6, "The learning rate of the gradient descent algorithm.")
	tuneNumCores := tuneCommand.Int("num-cores", 1, "The number of cores to assume can be used while tuning.")
	tuneNumPositions := tuneCommand.Int(
		"num-positions",
		1e6,
		"The number of positions to try to load for tuning. If there are fewer\n"+
			"positions, as many will be read as possible.",
	)
	tuneUseDefaultWeights := tuneCommand.Bool(
		"use-default-weights",
		true,
		"Use default weights for a fresh tuning session, or the current ones in evaluation.go",
	)

	genFENsCommand := flag.NewFlagSet("gen-fens", flag.ExitOnError)
	genFENsInputFile := genFENsCommand.String("input-file", "", "The input pgn file to extract quiet fens/positions from.")
	genFENsOutputFile := genFENsCommand.String("output-file", "fens.epd", "The file to output the quiet fens too. If one is not given, a file will be created.")
	genFENsSampleSize := genFENsCommand.Int("sample-size", 10, "The number of random quiet positions to extract from each game.")
	genFENsMaxGames := genFENsCommand.Int("max-games", 1e5, "The maximum number of games to attempt to extract quiet positions from.")
	genFENsMinElo := genFENsCommand.Int(
		"min-elo",
		0,
		"The minimum Elo that white and black must be rated for a game to be included.\n"+
			"Unrated games are skipped if the value is greater than 0.",
	)

	perftCmd := flag.NewFlagSet("perft", flag.ExitOnError)
	perftFen := perftCmd.String("fen", engine.FENStartPosition, "The position to start PERFT from.")
	perftDepth := perftCmd.Int("depth", 1, "The depth to run PERFT up to.")
	perftDivide := perftCmd.Bool("divide", false, "Display the number of nodes each move produces.")

	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	printFen := printCmd.String("fen", engine.FENStartPosition, "The position to display")

	if len(os.Args) < 2 {
		uciInterface := engine.UCIInterface{}
		uciInterface.UCILoop()
	} else {
		switch os.Args[1] {
		case "tune":
			tuneCommand.Parse(os.Args[2:])

			if *tuneInputFile == "" {
				fmt.Println("\nInput file is needed for tuning")
				os.Exit(1)
			}

			tuner.Tune(*tuneInputFile, *tuneEpochs, *tuneNumPositions, *tuneNumCores, *tuneLearningRate, false, *tuneUseDefaultWeights)
		case "gen-fens":
			genFENsCommand.Parse(os.Args[2:])

			if *genFENsInputFile == "" {
				fmt.Println("\nInput file is needed to generate fens")
				os.Exit(1)
			}

			tuner.GenTrainingData(*genFENsInputFile, *genFENsOutputFile, *genFENsSampleSize, uint16(*genFENsMinElo), uint32(*genFENsMaxGames))
		case "perft":
			perftCmd.Parse(os.Args[2:])

			pos := engine.NewPosition(*perftFen)
			nodes := uint64(0)
			var endTime time.Duration

			if *perftDivide {
				startTime := time.Now()
				nodes = engine.DividePerft(&pos, uint8(*perftDepth))
				endTime = time.Since(startTime)
			} else {
				startTime := time.Now()
				nodes = engine.Perft(&pos, uint8(*perftDepth))
				endTime = time.Since(startTime)
			}

			fmt.Printf("nodes: %d\n", nodes)
			fmt.Printf("ms: %d\n", endTime.Milliseconds())
			fmt.Printf("nps: %d\n", int(float64(nodes)/endTime.Seconds()))
		case "print":
			printCmd.Parse(os.Args[2:])
			fmt.Println(engine.NewPosition(*printFen))
		case "help":
			fmt.Println("\nperft: Run PERFT testing ")
			fmt.Println("gen-fens: Extract a set of quiet fens from a given PGN file.")
			fmt.Println("tune: Run the engine's tuner.")
			fmt.Println("print: Display a position")
			fmt.Println("help: Show this help message")
			fmt.Println("\nNo command starts the UCI protocol")
		default:
			fmt.Printf("unrecognized command: '%s'. Type help for a list of valid command line arguments.\n", os.Args[1])
			os.Exit(1)
		}
	}
}
