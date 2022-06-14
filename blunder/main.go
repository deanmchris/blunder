package main

import (
	"blunder/engine"
	"blunder/tuner"
	"os"
)

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
}

func main() {
	homeDir, _ := os.UserHomeDir()
	tuner.GenTrainingData(homeDir+"\\Desktop\\data\\games.pgn", homeDir+"\\Desktop\\data\\fens.pgn", 2700)
	// tuner.Tune(homeDir+"\\Desktop\\data\\quiet-labeled.epd", 200, 725000, 100000)
	// engine.RunCommLoop()
}
