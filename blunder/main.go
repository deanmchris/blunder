package main

import "blunder/tuner"

func main() {
	/*var s engine.Search
	s.Pos.LoadFEN("rnb1kbnr/pppp1ppp/8/1P2p1q1/8/8/P1PPPPPP/RNBQKBNR w KQkq - 0 2")
	fmt.Println(s.Pos)

	s.Timer.TimeLeft = 300000
	s.TT.Resize(engine.DefaultTTSize)
	start := time.Now()
	m := s.Search()
	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	// engine.UCILoop()
	tuner.RunTuner(true)
	// engine.UCILoop()

	//var pos engine.Position
	//pos.LoadFEN("r1bqk2r/ppp1bp1p/7p/3pP3/3P4/3Q1N1P/PP3PP1/RN2K2R w KQkq - 0 13")
	//fmt.Println(engine.EvaluatePos(&pos))
}
