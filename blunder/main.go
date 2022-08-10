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
	tuneEpochs := tuneCommand.Int("epochs", 10000, "The number of epochs to run the tuner for.")
	tuneNumPositions := tuneCommand.Int(
		"num-positions",
		500000,
		"The number of positions to try to load for tuning. If there are fewer\n"+
			"positions, as many will be read as possible.",
	)
	tuneUseDefaultWeights := tuneCommand.Bool(
		"use-default-weights",
		false,
		"Use default weights for a fresh tuning session, or the current ones in evaluation.go",
	)
	tuneInputFile := tuneCommand.String(
		"input-file",
		"",
		"A file containing a set of fens to use for tuning. Should be in the format\n"+
			"<full fen> [<result float>], where 'result float' is either 1.0 (white won),\n"+
			"0.0 (black won), or 0.5 (draw).",
	)

	genFENsCommand := flag.NewFlagSet("gen-fens", flag.ExitOnError)
	genFENsInputFile := genFENsCommand.String("input-file", "", "The input pgn file to extract quiet fens/positions from.")
	genFENsOutputFile := genFENsCommand.String("output-file", "fens.epd", "The file to output the quiet fens too. If one is not given, a file will be created.")
	genFENsSampleSize := genFENsCommand.Int("sample-size", 10, "The number of random quiet positions to extract from each game.")

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

			if *genFENsOutputFile == "fens.epd" {
				fmt.Println("No output file given...attempting to create 'fens.epd'")
				_, err := os.Create("fens.epd")

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			tuner.GenTrainingData(*genFENsInputFile, *genFENsOutputFile, *genFENsSampleSize)
		case "tune":
			tuneCommand.Parse(os.Args[2:])

			if *tuneInputFile == "" {
				fmt.Println("\nInput file is needed for tuning")
				os.Exit(1)
			}

			tuner.Tune(*tuneInputFile, *tuneEpochs, *tuneNumPositions, false, *tuneUseDefaultWeights)
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
