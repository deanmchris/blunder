package main

import (
	"blunder/engine"
	"flag"
	"fmt"
	"os"
	"time"
)

func init() {
	engine.InitBitboard()
	engine.InitTables()
	engine.InitMagics()
	engine.InitZobrist()
}

func main() {
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
			fmt.Println("print: Display a position")
			fmt.Println("help: Show this help message")
			fmt.Println("\nNo command starts the UCI protocol")
		default:
			fmt.Printf("unrecognized command: '%s'. Type help for a list of valid command line arguments.\n", os.Args[1])
			os.Exit(1)
		}
	}
}
