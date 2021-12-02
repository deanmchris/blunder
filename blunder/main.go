package main

import (
	"blunder/engine"
)

func main() {
	// tuner.Tune("C:\\Users\\deanm\\Desktop\\quiet-labeled.epd", 725000)
	var inter engine.UCIInterface
	inter.UCILoop()
}
