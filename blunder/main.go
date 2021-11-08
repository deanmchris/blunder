package main

import "blunder/engine"

func main() {
	/*var s engine.Search
	s.Pos.LoadFEN(engine.FENKiwiPete)

	s.Timer.TimeLeft = 300000
	s.SpecifiedDepth = engine.MaxPly
	s.SpecifiedNodes = math.MaxUint64
	s.TT.Resize(engine.DefaultTTSize)

	start := time.Now()
	m := s.Search()

	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	// r2k2nr/1pp2n1p/p1qpbpp1/8/3B1P2/1NP3Q1/P1P3PP/1K1R1B1R w - - 3 21

	var inter engine.UCIInterface
	inter.UCILoop()

	// tuner.GenTrainingData("/blunder/games.pgn", "/tuner/positions.txt")
	// tuner.RunTuner(true)
}
