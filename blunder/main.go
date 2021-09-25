package main

import "blunder/engine"

func main() {
	// "8/p1ppk1p1/2n2p2/8/4B3/2P1KPP1/1P5P/8 w - - 0 1"
	// 8/1k1K4/4R3/8/8/8/8/8 b - - 80 221

	/*var s engine.Search
	s.Pos.LoadFEN(engine.FENStartPosition)
	s.Timer.TimeLeft = 300000
	s.TT.Resize(engine.DefaultTTSize)

	start := time.Now()
	m := s.Search()
	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	/*defer func() {
		if err := recover(); err != nil {
			println(fmt.Sprintf("%s", err))
			println(string(debug.Stack()))
		}
	}()*/

	engine.UCILoop()
}
