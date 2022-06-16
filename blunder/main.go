package main

import (
	"blunder/engine"
)

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
}

func main() {
	// homeDir, _ := os.UserHomeDir()
	// tuner.Tune(homeDir+"\\Desktop\\data\\quiet.epd", 10000, 725000, true)
	engine.RunCommLoop()
}
