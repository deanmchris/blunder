package main

import (
	"blunder/engine"
)

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
	// tuner.RunTuner(true)
	engine.UCILoop()
}
