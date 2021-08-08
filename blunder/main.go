package main

import (
	"blunder/ui"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const HelpMessage = `
Options:
- uci: Begin Blunder's UCI protocol
- cli: Enter Blunder's command-line interface
- help: Quit the program
`

func mainLoop() {
	for {
		reader := bufio.NewReader(os.Stdin)
		programMode, _ := reader.ReadString('\n')

		if programMode == "uci\n" || programMode == "uci" {
			ui.UCILoop()
			break
		} else if programMode == "cli\n" {
			ui.CmdLoop()
			break
		} else if programMode == "quit\n" {
			break
		} else if programMode == "help\n" {
			fmt.Println(HelpMessage)
		} else {
			fmt.Printf("\nUnknown command \"%v\"\n", strings.TrimSuffix(programMode, "\n"))
			fmt.Printf("Enter \"help\" to show available commands\n\n")
		}
	}

}

func main() {
	mainLoop()
}
