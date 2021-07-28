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
- interactive: Enter Blunder's interactive debug mode
- help: Quit the program
`

func mainLoop() {
	for {
		reader := bufio.NewReader(os.Stdin)
		programMode, _ := reader.ReadString('\n')

		if programMode == "uci\n" || programMode == "uci" {
			ui.UCILoop()
			break
		} else if programMode == "interactive\n" {
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
