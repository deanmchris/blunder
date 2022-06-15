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
	tuner.Tune(homeDir+"\\Desktop\\data\\quiet.epd", 50000, 10000)
	// engine.RunCommLoop()
}
