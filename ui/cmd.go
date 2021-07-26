package ui

import (
	"blunder/engine"
	"blunder/tests"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DepthLimit  = 7
	HelpMessage = `
Options:
- perft <DEPTH>: Run perft up to <DEPTH>
- fen <FEN>: Load a fen string given by <FEN>
- print: Display the current board state
- options: Display this help message
- quit: Quit the program
- test: Run the perft tests for Blunder (will take a bit!)
`
)

// Run the loop for the command mode of Blunder
func CmdLoop() {
	var board engine.Board
	board.LoadFEN(engine.FENStartPosition)
	fmt.Println(board)

	for {
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')

		if command == "quit\n" {
			break
		} else if command == "print\n" {
			fmt.Println(board)
		} else if command == "options\n" {
			fmt.Println(HelpMessage)
		} else if command == "test\n" {
			fmt.Println()
			tests.RunPerftTests(&board)
			fmt.Println()
		} else if strings.HasPrefix(command, "perft ") {
			perftCommand(&board, command)
		} else if strings.HasPrefix(command, "fen ") {
			fenCommand(&board, command)
		} else {
			fmt.Printf("Unknown command \"%v\"\n", strings.TrimSuffix(command, "\n"))
			fmt.Printf("Enter \"options\" to show available commands\n")
		}
	}
}

// Run the perft command in the command line mode
func perftCommand(board *engine.Board, command string) {
	command = strings.TrimPrefix(command, "perft ")
	command = strings.TrimSuffix(command, "\n")

	depth, err := strconv.Atoi(command)
	if err == nil {
		if depth <= DepthLimit {
			start := time.Now()
			fmt.Println()
			nodes := engine.Perft(board, depth, depth, false)
			fmt.Println("\nNodes:", nodes)
			elapsed := time.Since(start)
			fmt.Printf("Time: %vms\n", elapsed.Milliseconds())
			fmt.Printf("Nps: %d\n", int(float64(nodes)/elapsed.Seconds()))
		} else {
			fmt.Printf("Depth limit for perft is %d", DepthLimit)
		}
	} else {
		fmt.Println("Perft depth should be an integer")
	}
}

// Run the fen command in the command line mode
func fenCommand(board *engine.Board, command string) {
	command = strings.TrimPrefix(command, "fen ")
	command = strings.TrimSuffix(command, "\n")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("fen entered is not valid")
		}
	}()
	board.LoadFEN(command)
}
