package main

import "blunder/engine"

func main() {
	// tuner.SearchTuner("C:\\Users\\deanm\\Desktop\\blunder\\testdata\\win_at_chess.epd", 5)
	// tuner.Tune("C:\\Users\\deanm\\Desktop\\quiet-labeled.epd", 725000)
	var inter engine.UCIInterface
	inter.UCILoop()

	/*var search engine.Search
	search.Pos.LoadFEN("2r2rn1/p3q1bk/n2p2pp/1b1PPp2/1pp2P1N/1P2N1P1/P1QB2BP/4RRK1 w - - 0 1")

	search.Timer.TimeLeft = 300000
	search.SpecifiedDepth = uint8(engine.MaxPly)
	search.SpecifiedNodes = uint64(math.MaxUint64)
	search.TT.Resize(engine.DefaultTTSize)

	start := time.Now()
	bm := search.Search()
	end := time.Since(start)

	fmt.Println("best move:", bm)
	fmt.Println("time (ms):", end.Milliseconds())*/
}
