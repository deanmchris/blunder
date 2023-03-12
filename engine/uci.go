package engine

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	EngineName        = "Blunder 9.0.0"
	EngineAuthor      = "Christian Dean"
	EngineAuthorEmail = "deanmchris@gmail.com"
)

type UCIInterface struct {
	search Search
}

func (inter *UCIInterface) reset() {
	inter.search = NewSearch(FENStartPosition)
}

func (inter *UCIInterface) uciCommandResponse() {
	fmt.Printf("\nid name %s\n", EngineName)
	fmt.Printf("id author %s\n", EngineAuthor)
	fmt.Printf("\noption name Hash type spin default 64 min 1 max 32000\n")
	fmt.Print("option name Clear Hash type button\n")

	fmt.Print("\nAvailable UCI commands:\n")
	fmt.Print("    * uci\n    * isready\n    * ucinewgame")
	fmt.Print("\n    * setoption name <NAME> value <VALUE>")

	fmt.Print("\n    * position")
	fmt.Print("\n\t* fen <FEN>")
	fmt.Print("\n\t* startpos")

	fmt.Print("\n    * go")
	fmt.Print("\n\t* wtime <MILLISECONDS>\n\t* btime <MILLISECONDS>")
	fmt.Print("\n\t* winc <MILLISECONDS>\n\t* binc <MILLISECONDS>")
	fmt.Print("\n\t* movestogo <INTEGER>\n\t* depth <INTEGER>\n\t* nodes <INTEGER>\n\t* movetime <MILLISECONDS>")
	fmt.Print("\n\t* infinite")

	fmt.Print("\n    * stop\n    * quit\n\n")
	fmt.Printf("uciok\n\n")
}

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
			inter.search.TT.Resize(uint64(size))
		}
	case "Clear Hash":
		inter.search.TT.Clear()
	}
}

func (inter *UCIInterface) positionCommandResponse(command string) {
	args := strings.TrimPrefix(command, "position ")
	fenString := ""

	if strings.HasPrefix(args, "startpos") {
		args = strings.TrimPrefix(args, "startpos ")
		fenString = FENStartPosition
	} else if strings.HasPrefix(args, "fen") {
		args = strings.TrimPrefix(args, "fen ")
		remaining_args := strings.Fields(args)
		fenString = strings.Join(remaining_args[0:6], " ")
		args = strings.Join(remaining_args[6:], " ")
	}

	inter.search.Setup(fenString)
	if strings.HasPrefix(args, "moves") {
		args = strings.TrimSuffix(strings.TrimPrefix(args, "moves"), " ")
		if args != "" {
			for _, moveAsString := range strings.Fields(args) {
				move := moveFromCoord(&inter.search.Pos, moveAsString)
				valid := inter.search.Pos.DoMove(move)
				inter.search.AddHistory(inter.search.Pos.Hash)

				// Every move played from the UCI input should be valid,
				// if we detect it's not it's some sort of bug, so raise
				// an error.
				if !valid {
					panic(fmt.Sprintf("Illegal moved detected from UCI input: %s", moveToStr(move)))
				}

				// Decrementing the history counter here makes
				// sure that no state is saved on the position's
				// history stack since this move will never be undone.
				inter.search.Pos.StateIdx--
			}
		}
	}
}

func (inter *UCIInterface) goCommandResponse(command string) {
	command = strings.TrimPrefix(command, "go")
	command = strings.TrimPrefix(command, " ")
	fields := strings.Fields(command)

	colorPrefix := "b"
	if inter.search.Pos.SideToMove == White {
		colorPrefix = "w"
	}

	// Parse the go command arguments.
	timeLeft, increment, movesToGo := int(InfiniteTime), int(NoValue), int(NoValue)
	maxDepth, maxNodeCount, moveTime := uint64(MaxPly), uint64(math.MaxUint64), uint64(NoValue)

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
			maxDepth, _ = strconv.ParseUint(fields[index+1], 10, 8)
		} else if field == "nodes" {
			maxNodeCount, _ = strconv.ParseUint(fields[index+1], 10, 64)
		} else if field == "movetime" {
			moveTime, _ = strconv.ParseUint(fields[index+1], 10, 64)
		}
	}

	inter.search.Timer.Setup(
		int64(timeLeft),
		int64(increment),
		int64(moveTime),
		int16(movesToGo),
		uint8(maxDepth),
		maxNodeCount,
	)

	bestMove := inter.search.RunSearch()
	fmt.Printf("bestmove %v\n", moveToStr(bestMove))
}

func (inter *UCIInterface) UCILoop() {
	fmt.Println("Author:", EngineAuthor)
	fmt.Println("Engine:", EngineName)
	fmt.Println("Email:", EngineAuthorEmail)
	fmt.Printf("Default hash size: %d\n", SearchTTSize)
	fmt.Printf("Default PERFT hash size: %d\n", PerftTTSize)

	reader := bufio.NewReader(os.Stdin)

	inter.uciCommandResponse()
	inter.reset()

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
			inter.reset()
		} else if strings.HasPrefix(command, "position") {
			inter.positionCommandResponse(command)
		} else if strings.HasPrefix(command, "go") {
			go inter.goCommandResponse(command)
		} else if strings.HasPrefix(command, "stop") {
			inter.search.StopSearch()
		} else if command == "quit\n" {
			break
		}
	}
}
