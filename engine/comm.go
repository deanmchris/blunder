package engine

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EngineName   = "Blunder 8.0.0"
	EngineAuthor = "Christian Dean"
	EngineEmail  = "deanmchris@gmail.com"

	Banner = `
██████╗░██╗░░░░░██╗░░░██╗███╗░░██╗██████╗░███████╗██████╗░
██╔══██╗██║░░░░░██║░░░██║████╗░██║██╔══██╗██╔════╝██╔══██╗
██████╦╝██║░░░░░██║░░░██║██╔██╗██║██║░░██║█████╗░░██████╔╝
██╔══██╗██║░░░░░██║░░░██║██║╚████║██║░░██║██╔══╝░░██╔══██╗
██████╦╝███████╗╚██████╔╝██║░╚███║██████╔╝███████╗██║░░██║
╚═════╝░╚══════╝░╚═════╝░╚═╝░░╚══╝╚═════╝░╚══════╝╚═╝░░╚═╝
	`

	PerftDepthLimit = 255
	HelpMessage     = `
Options:
- uci: Start the UCI protocol
- perft <DEPTH>: Run perft up to <DEPTH>
- dperft <DEPTH>: Run divide perft up to <DEPTH>
- fen <FEN>: Load a fen string given by <FEN>
- print: Display the current board state
- eval: Display a break down of the static evaluation of the current position
- help: Display this help message
- quit: Quit the program

`
)

func perftCommand(pos *Position, command string) {
	command = strings.TrimPrefix(command, "perft ")
	command = strings.TrimSuffix(command, "\n")

	depth, err := strconv.Atoi(command)
	if err == nil {
		if depth <= PerftDepthLimit {
			start := time.Now()
			fmt.Println()
			nodes := Perft(pos, uint8(depth))
			fmt.Println("\nNodes:", nodes)
			elapsed := time.Since(start)
			fmt.Printf("Time: %vms\n", elapsed.Milliseconds())
			fmt.Printf("Nps: %d\n", int(float64(nodes)/elapsed.Seconds()))
		} else {
			fmt.Printf("Depth limit for perft is %d", PerftDepthLimit)
		}
	} else {
		fmt.Println("Perft depth should be an integer")
	}
}

func dividePerftCommand(pos *Position, command string) {
	command = strings.TrimPrefix(command, "dperft ")
	command = strings.TrimSuffix(command, "\n")

	depth, err := strconv.Atoi(command)
	if err == nil {
		if depth <= PerftDepthLimit {
			start := time.Now()
			fmt.Println()
			nodes := DividePerft(pos, uint8(depth), uint8(depth))
			fmt.Println("\nNodes:", nodes)
			elapsed := time.Since(start)
			fmt.Printf("Time: %vms\n", elapsed.Milliseconds())
			fmt.Printf("Nps: %d\n\n", int(float64(nodes)/elapsed.Seconds()))
		} else {
			fmt.Printf("Depth limit for perft is %d\n", PerftDepthLimit)
		}
	} else {
		fmt.Println("Perft depth should be an integer")
		fmt.Println()
	}
}

func fenCommand(pos *Position, command string) {
	command = strings.TrimPrefix(command, "fen ")
	command = strings.TrimSuffix(command, "\n")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("fen entered is not valid")
		}
	}()
	pos.LoadFEN(command)
}

func RunCommLoop() {
	fmt.Println(Banner)
	fmt.Println("Author:", EngineAuthor)
	fmt.Println("Engine:", EngineName)
	fmt.Println("Email:", EngineEmail)
	fmt.Printf("Hash size: %d MB\n\n", DefaultTTSize)

	reader := bufio.NewReader(os.Stdin)
	var inter UCIInterface
	inter.Search.Pos.LoadFEN(FENStartPosition)

	for {
		command, _ := reader.ReadString('\n')
		command = strings.Replace(command, "\r\n", "\n", -1)

		if strings.HasPrefix(command, "perft") {
			perftCommand(&inter.Search.Pos, command)
		} else if strings.HasPrefix(command, "dperft ") {
			dividePerftCommand(&inter.Search.Pos, command)
		} else if strings.HasPrefix(command, "fen ") {
			fenCommand(&inter.Search.Pos, command)
		} else if command == "print\n" {
			fmt.Println(inter.Search.Pos)
		} else if command == "uci\n" {
			inter.UCILoop()
			break
		} else if command == "help\n" {
			fmt.Print(HelpMessage)
		} else if command == "eval\n" {
			evaluatePosTrace(&inter.Search.Pos)
		} else if command == "quit\n" {
			break
		} else {
			fmt.Printf("Unknown command \"%v\"\n", strings.TrimSuffix(command, "\n"))
			fmt.Printf("Enter \"help\" to show available commands\n")
		}
	}
}
