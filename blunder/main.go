package main

import (
	"blunder/engine"
)

func main() {
	/*var s engine.Search
	s.Pos.LoadFEN("8/6k1/8/8/4K3/8/8/8 w - - 2 190")
	s.Timer.TimeLeft = 300000
	s.TT.Resize(engine.DefaultTTSize)
	start := time.Now()
	m := s.Search()
	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	engine.UCILoop()
}
