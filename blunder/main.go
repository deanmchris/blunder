package main

import (
	"blunder/engine"
)

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
	engine.InitEvalBitboards()
	engine.InitSearchTables()
}

func main() {
	// homeDir, _ := os.UserHomeDir()
	// tuner.Tune(homeDir+"\\Desktop\\data\\quiet.epd", 500, 725000, true)
	engine.RunCommLoop()
}
