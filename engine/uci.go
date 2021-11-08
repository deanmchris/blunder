package engine

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
    "path/filepath"
)

const (
	EngineName   = "Blunder 7.2.0"
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

type UCIInterface struct {
	Search        Search
	OpeningBook   map[uint64]PolyglotEntry
	OptionUseBook bool
}

// Respond to the command "uci"
func (inter *UCIInterface) uciCommandResponse() {
	fmt.Printf("\nid name %v\n", EngineName)
	fmt.Printf("id author %v\n", EngineAuthor)
	fmt.Printf("\noption name Hash type spin default 64 min 1 max 32000\n")
	fmt.Print("option name Clear Hash type button\n")
	fmt.Print("option name Clear History type button\n")
	fmt.Print("option name OwnBook type check default false\n")
	fmt.Printf("uciok\n\n")
}

// Respond to the command "position"
func (inter *UCIInterface) positionCommandResponse(command string) {
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
	inter.Search.Pos.LoadFEN(fenString)
	if strings.HasPrefix(args, "moves") {
		args = strings.TrimPrefix(args, "moves ")
		for _, moveAsString := range strings.Fields(args) {
			move := MoveFromCoord(&inter.Search.Pos, moveAsString)
			inter.Search.Pos.MakeMove(move)

			// Decrementing the history counter here makes
			// sure that no state is saved on the position's
			// history stack since this move will never be undone.
			inter.Search.Pos.StatePly--
		}
	}
}

// Respond to the command "setoption"
func (inter *UCIInterface) setOptionCommandResponse(command string) {
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
			inter.Search.TT.Unitialize()
			inter.Search.TT.Resize(uint64(size))
		}
	case "Clear Hash":
		inter.Search.TT.Clear()
	case "Clear History":
		inter.Search.ClearHistoryTable()
	case "OwnBook":
		if value == "true" {
			inter.OptionUseBook = true
		} else if value == "false" {
			inter.OptionUseBook = false
		}
	}
}

// Respond to the command "go"
func (inter *UCIInterface) goCommandResponse(command string) {
    if inter.OptionUseBook {
        if entry, ok := inter.OpeningBook[GenPolyglotHash(&inter.Search.Pos)]; ok {
            move := MoveFromCoord(&inter.Search.Pos, entry.Move)
            if inter.Search.Pos.MoveIsPseduoLegal(move) {
                fmt.Printf("bestmove %v\n", move) 
                return
            }
	    }
    }
        
	command = strings.TrimPrefix(command, "go ")
	fields := strings.Fields(command)

	colorPrefix := "b"
	if inter.Search.Pos.SideToMove == White {
		colorPrefix = "w"
	}

	// Parse the time left, increment, and moves to go from the command parameters.
	timeLeft, increment, movesToGo := -1, 0, 0
	specifiedDepth := uint64(MaxPly)
	specifiedNodes := uint64(math.MaxUint64)

	for index, field := range fields {
		if strings.HasPrefix(field, colorPrefix) {
			if strings.HasSuffix(field, "time") {
				timeLeft, _ = strconv.Atoi(fields[index+1])
			} else if strings.HasSuffix(field, "inc") {
				increment, _ = strconv.Atoi(fields[index+1])
			}
		} else if field == "movestogo" {
			movesToGo, _ = strconv.Atoi(fields[index+1])
		} else if field == "depth" {
			specifiedDepth, _ = strconv.ParseUint(fields[index+1], 10, 8)
		} else if field == "nodes" {
			specifiedNodes, _ = strconv.ParseUint(fields[index+1], 10, 64)
		}

	}

	// Setup the timer with the go command time control information.
	inter.Search.Timer.TimeLeft = int64(timeLeft)
	inter.Search.Timer.Increment = int64(increment)
	inter.Search.Timer.MovesToGo = int64(movesToGo)

	// Setup user defined search options if given.
	inter.Search.SpecifiedDepth = uint8(specifiedDepth)
	inter.Search.SpecifiedNodes = specifiedNodes

	// Report the best move found by the engine to the GUI.
	bestMove := inter.Search.Search()
	fmt.Printf("bestmove %v\n", bestMove)
}

func (inter *UCIInterface) quitCommandResponse() {
	inter.Search.TT.Unitialize()
}

func (inter *UCIInterface) printCommandResponse() {
	// print internal engine info
}

func (inter *UCIInterface) UCILoop() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(Banner)
	fmt.Println("Author:", EngineAuthor)
	fmt.Println("Engine:", EngineName)
	fmt.Println("Email:", EngineEmail)
	fmt.Printf("Hash size: %d MB\n\n", DefaultTTSize)

	inter.Search.TT.Resize(DefaultTTSize)
	inter.Search.Pos.LoadFEN(FENStartPosition) 
    inter.OpeningBook = make(map[uint64]PolyglotEntry)
    
    wd, _ := os.Getwd()
    parentFolder := filepath.Dir(wd)
    inter.OpeningBook, _ = LoadPolyglotFile(filepath.Join(parentFolder, "/book/book.bin"))

	for {
		command, _ := reader.ReadString('\n')
		command = strings.Replace(command, "\r\n", "\n", -1)

		if command == "uci\n" {
			inter.uciCommandResponse()
		} else if command == "isready\n" {
			fmt.Printf("readyok\n")
		} else if strings.HasPrefix(command, "setoption") {
			inter.setOptionCommandResponse(command)
		} else if strings.HasPrefix(command, "ucinewgame") {
			inter.Search.TT.Clear()
			inter.Search.ClearHistoryTable()
		} else if strings.HasPrefix(command, "position") {
			inter.positionCommandResponse(command)
		} else if strings.HasPrefix(command, "go") {
			go inter.goCommandResponse(command)
		} else if strings.HasPrefix(command, "stop") {
			inter.Search.Timer.Stop = true
		} else if command == "quit\n" {
			inter.quitCommandResponse()
			break
		} else if command == "print\n" {
			inter.printCommandResponse()
		}
	}
}
