package main

import (
	"blunder/extra"
)

func main() {
	/*var s engine.Search
	s.Pos.LoadFEN("r1bqkb1r/pppp1ppp/2n2n2/4p3/3PP3/2P2N2/PP3PPP/RNBQKB1R b KQkq - 0 4")

	s.Timer.TimeLeft = 300000
	s.SpecifiedDepth = uint8(engine.MaxPly)
	s.SpecifiedNodes = uint64(math.MaxUint64)

	s.TT.Resize(engine.DefaultTTSize)
	start := time.Now()
	m := s.Search()
	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	// engine.UCILoop()
	// Horrible move choice by Blunder here. Why? 1q4r1/3Q1Npk/p6p/1p5N/8/7P/Pn3PP1/6K1 w - - 0 1
	extra.TestIQ(10)
}
