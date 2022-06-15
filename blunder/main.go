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
	/*tuner.GenTrainingData(
		homeDir+"\\Desktop\\data\\gm_2700.pgn",
		homeDir+"\\Desktop\\data\\gm_2700.epd",
		0,
		2000,
	)*/
	// tuner.Tune(homeDir+"\\Desktop\\data\\quiet.epd", 2000, 725000, 2)
	// tuner.Tune(homeDir+"\\Desktop\\data\\quiet.epd", 25000, 50000, 150)
	engine.RunCommLoop()
}
