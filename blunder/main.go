package main

import "blunder/engine"

func main() {
	// tuner.RunTuner(true)
	// tuner.GenTrainingData("/home/algerbrex/games.pgn", "/home/algerbrex/pos2.epd")
	var inter engine.UCIInterface
	inter.UCILoop()

	/*var pos engine.Position
	pos.LoadFEN("r1bqr1k1/1ppnpp1p/p2p1npQ/6N1/3PP3/2N5/PPP2PPP/R3KB1R w KQ - 3 10")
	fmt.Println(engine.EvaluatePos(&pos))*/
}
