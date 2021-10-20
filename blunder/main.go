package main

import "blunder/engine"

func main() {
	/*var s engine.Search
	s.Pos.LoadFEN("1r4k1/5p1p/1p3np1/pR6/P1Pr3P/2R2B2/1P4P1/7K w - - 0 1")

	s.Timer.TimeLeft = 300000
	s.TT.Resize(engine.DefaultTTSize)
	start := time.Now()
	m := s.Search()
	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	engine.UCILoop()
}
