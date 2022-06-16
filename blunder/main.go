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
	engine.InitEvalBitboards()
}

func main() {
	homeDir, _ := os.UserHomeDir()
	tuner.Tune(homeDir+"\\Desktop\\data\\quiet.epd", 10000, 725000, true)
	// engine.RunCommLoop()
}
