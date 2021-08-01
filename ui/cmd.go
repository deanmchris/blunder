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
- help: Display this help message
- ptest: Run the perft tests for Blunder
- ztest: Run the zobrist hash tests (resets board to starting position)
- quit: Quit the program
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
		} else if command == "help\n" {
			fmt.Println(HelpMessage)
		} else if command == "ptest\n" {
			fmt.Println()
			tests.RunPerftTests(&board)
			fmt.Println()
		} else if command == "ztest\n" {
			fmt.Println()
			board.LoadFEN(engine.FENStartPosition)
			tests.RunAllZobristHashingTests(&board)
			fmt.Println()
		} else if strings.HasPrefix(command, "perft ") {
			perftCommand(&board, command)
		} else if strings.HasPrefix(command, "fen ") {
			fenCommand(&board, command)
		} else {
			fmt.Printf("Unknown command \"%v\"\n", strings.TrimSuffix(command, "\n"))
			fmt.Printf("Enter \"help\" to show available commands\n")
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
			fmt.Printf("Depth limit for perft is %d\n", DepthLimit)
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
