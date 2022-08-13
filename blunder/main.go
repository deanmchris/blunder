package main

import (
	"blunder/engine"
	"blunder/tuner"
	"flag"
	"fmt"
	"os"
)

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
	engine.InitEvalBitboards()
	engine.InitSearchTables()
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
	tuneLearningRate := tuneCommand.Float64("learning-rate", 0.5, "The learning rate of the gradient descent algorithm.")
	tuneNumCores := tuneCommand.Int("num-cores", 1, "The number of cores to assume can be used while tuning.")
	tuneNumPositions := tuneCommand.Int(
		"num-positions",
		1000000,
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
	genFENsMaxGames := genFENsCommand.Int("max-games", 200000, "The maximum number of games to attempt to extract quiet positions from.")
	genFENsMinElo := genFENsCommand.Int(
		"min-elo",
		0,
		"The minimum Elo that white and black must be rated for a game to be included.\n"+
			"Unrated games are skipped if the value is greater than 0.",
	)

	if len(os.Args) < 2 {
		engine.RunCommLoop()
	} else {
		switch os.Args[1] {
		case "gen-magics":
			fmt.Println("Generate new rook and bishop magic numbers")
			engine.GenBishopMagics()
			engine.GenRookMagics()

			fmt.Print("\nBishop magic numbers:\n\n")
			for i, magic := range engine.BishopMagics {
				if i%4 == 0 && i > 0 {
					fmt.Println()
				}

				fmt.Printf("0x%x,", magic.MagicNo)
			}

			fmt.Printf("\n\nRook magic numbers:\n\n")
			for i, magic := range engine.RookMagics {
				if i%4 == 0 && i > 0 {
					fmt.Println()
				}
				fmt.Printf("0x%x,", magic.MagicNo)
			}

			fmt.Println()
		case "gen-fens":
			genFENsCommand.Parse(os.Args[2:])

			if *genFENsInputFile == "" {
				fmt.Println("\nInput file is needed to generate fens")
				os.Exit(1)
			}

			tuner.GenTrainingData(*genFENsInputFile, *genFENsOutputFile, *genFENsSampleSize, uint16(*genFENsMinElo), uint32(*genFENsMaxGames))
		case "tune":
			tuneCommand.Parse(os.Args[2:])

			if *tuneInputFile == "" {
				fmt.Println("\nInput file is needed for tuning")
				os.Exit(1)
			}

			tuner.Tune(*tuneInputFile, *tuneEpochs, *tuneNumPositions, *tuneLearningRate, *tuneNumCores, false, *tuneUseDefaultWeights)
		case "help":
			fmt.Println("\ngen-magics: Generate a new set of magic numbers for rooks and bishops.")
			fmt.Println("gen-fens: Extract a set of quiet fens from a given PGN file.")
			fmt.Println("tune: Run the engine's tuner.")
			fmt.Println("\nNo command starts the modified UCI protocol")
		default:
			fmt.Printf("unrecognized command: '%s'. Type help for the commands Blunder understands", os.Args[1])
			os.Exit(1)
		}
	}
}
