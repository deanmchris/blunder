package engine

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	EngineName   = "Blunder 0.9.0"
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
)

// Respond to the command "uci"
func uciCommandResponse() {
	fmt.Printf("\nid name %v\n", EngineName)
	fmt.Printf("id author %v\n", EngineAuthor)
	fmt.Printf("\noption name Hash type spin default 32 min 1 max 256\n")
	fmt.Print("option name ClearHash type button\n\n")
	fmt.Printf("uciok\n\n")
}

// Respond to the command "position"
func positionCommandResponse(pos *Position, command string) {
	// Load in the fen string describing the position,
	// or load in the starting position.
	args := strings.TrimPrefix(command, "position ")
	var fenString string
	if strings.HasPrefix(args, "startpos") {
		args = strings.TrimPrefix(args, "startpos ")
		fenString = FENStartPosition
	} else if strings.HasPrefix(args, "fen") {
		args = strings.TrimPrefix(args, "fen ")
		remaining_args := strings.Fields(args)
		fenString = strings.Join(remaining_args[0:6], " ")
		args = strings.Join(remaining_args[6:], " ")
	}

	// Set the board to the appropriate position and make
	// the moves that have occured if any to update the position.
	pos.LoadFEN(fenString)
	if strings.HasPrefix(args, "moves") {
		args = strings.TrimPrefix(args, "moves ")
		for _, moveAsString := range strings.Fields(args) {
			move := moveFromCoord(pos, moveAsString)
			pos.MakeMove(move)

			// Decrementing the history counter here makes
			// sure that no state is saved on the position's
			// history stack since this move will never be undone.
			pos.StatePly--
		}
	}
}

// Respond to the command "setoption"
func setOptionCommandResponse(search *Search, command string) {
	fields := strings.Fields(command)
	var option, value string
	parsingWhat := ""

	for _, field := range fields {
		if field == "name" {
			parsingWhat = "name"
		} else if field == "value" {
			parsingWhat = "value"
		} else if parsingWhat == "name" {
			option += field + " "
		} else if parsingWhat == "value" {
			value += field + " "
		}
	}

	option = strings.TrimSuffix(option, " ")
	value = strings.TrimSuffix(value, " ")

	switch option {
	case "Hash":
		size, err := strconv.Atoi(value)
		if err == nil {
			search.TT.Unitialize()
			search.TT.Resize(uint64(size))
		}
	case "Clear Hash":
		search.TT.Clear()
	}
}

// Respond to the command "go"
func goCommandResponse(search *Search, command string) {
	command = strings.TrimPrefix(command, "go ")
	fields := strings.Fields(command)

	colorPrefix := "b"
	if search.Pos.SideToMove == White {
		colorPrefix = "w"
	}

	// Parse the time left, increment, and moves to go from the command parameters.
	timeLeft, increment, movesToGo := -1, 0, 0
	for index, field := range fields {
		if strings.HasPrefix(field, colorPrefix) {
			if strings.HasSuffix(field, "time") {
				timeLeft, _ = strconv.Atoi(fields[index+1])
			} else if strings.HasSuffix(field, "inc") {
				increment, _ = strconv.Atoi(fields[index+1])
			}
		} else if field == "movestogo" {
			movesToGo, _ = strconv.Atoi(fields[index+1])
		}

	}

	// Setup the timer with the go command time control information.
	search.Timer.TimeLeft = int64(timeLeft)
	search.Timer.Increment = int64(increment)
	search.Timer.MovesToGo = int64(movesToGo)

	// Report the best move found by the engine to the GUI.
	bestMove := search.Search()
	fmt.Printf("bestmove %v\n", bestMove)
}

func quitCommandResponse(search *Search) {
	search.TT.Unitialize()
}

func printCommandResponse() {
	// print internal engine info
}

func UCILoop() {
	reader := bufio.NewReader(os.Stdin)
	var search Search

	fmt.Println(Banner)
	fmt.Println("Author:", EngineAuthor)
	fmt.Println("Engine:", EngineName)
	fmt.Println("Email:", EngineEmail)
	fmt.Printf("Hash size: %d MB\n\n", DefaultTTSize)

	search.TT.Resize(DefaultTTSize)

	for {
		command, _ := reader.ReadString('\n')
		command = strings.Replace(command, "\r\n", "\n", -1)

		if command == "uci\n" {
			uciCommandResponse()
		} else if command == "isready\n" {
			fmt.Printf("readyok\n")
		} else if strings.HasPrefix(command, "setoption") {
			setOptionCommandResponse(&search, command)
		} else if strings.HasPrefix(command, "ucinewgame") {
			search.TT.Clear()
		} else if strings.HasPrefix(command, "position") {
			positionCommandResponse(&search.Pos, command)
		} else if strings.HasPrefix(command, "go") {
			go goCommandResponse(&search, command)
		} else if strings.HasPrefix(command, "stop") {
			search.Timer.Stop = true
		} else if command == "quit\n" {
			quitCommandResponse(&search)
			break
		} else if command == "print\n" {
			printCommandResponse()
		}
	}
}
