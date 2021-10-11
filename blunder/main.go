package main

import (
	"blunder/engine"
)

func main() {
	/*var s engine.Search
	s.Pos.LoadFEN("rn1q1rk1/1b2bppp/1pn1p3/p2pP3/3P4/P2BBN1P/1P1N1PP1/R2Q1RK1 b - - 0 1")

	s.Timer.TimeLeft = 300000
	s.TT.Resize(engine.DefaultTTSize)
	start := time.Now()
	m := s.Search()
	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	engine.UCILoop()
	// tuner.RunTuner(true)
	// engine.UCILoop()

	//var pos engine.Position
	//pos.LoadFEN("r1bqk2r/ppp1bp1p/7p/3pP3/3P4/3Q1N1P/PP3PP1/RN2K2R w KQkq - 0 13")
	//fmt.Println(engine.EvaluatePos(&pos))
}
