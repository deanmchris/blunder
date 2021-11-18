package main

import "blunder/engine"

func main() {
	// tuner.RunTuner(true)
	var inter engine.UCIInterface
	inter.UCILoop()
	/*var pos engine.Position
	pos.LoadFEN("r1bqkb1r/pppn1ppp/3p1n2/8/3NP3/2N5/PPP1QP1P/R1B1KBR1 b Qkq - 2 8")
	fmt.Println(engine.EvaluatePos(&pos))*/
}
