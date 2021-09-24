package main

import (
	"blunder/engine"
)

func main() {
	// "8/p1ppk1p1/2n2p2/8/4B3/2P1KPP1/1P5P/8 w - - 0 1"
	// 8/1k1K4/4R3/8/8/8/8/8 b - - 80 221

	/*var s engine.Search
	s.Pos.LoadFEN("1rbq1rk1/p1b1nppp/1p2p3/8/1B1pN3/P2B4/1P3PPP/2RQ1R1K w - - 0 1")
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
