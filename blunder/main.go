package main

import (
	"blunder/engine"
)

func main() {
	// tuner.GenTrainingData("C:\\Users\\deanm\\Desktop\\engines\\self_play.pgn", "C:\\Users\\deanm\\Desktop\\data.txt")
	// tuner.Tune("C:\\Users\\deanm\\Desktop\\data.txt", 850000, 100)
	engine.RunCommLoop()

	/*var pos engine.Position
	pos.LoadFEN("rnbqkbnr/pp1p1ppp/8/2pNp3/2P5/8/PPP1PPPP/RNBQKB1R w KQkq - 0 1")
	fmt.Println("CP:", engine.EvaluatePos(&pos))*/
}
