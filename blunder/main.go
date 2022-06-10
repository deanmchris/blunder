package main

import (
	"blunder/engine"
	"blunder/tuner"
)

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
}

func main() {
	tuner.Tune("C:\\Users\\deanm\\Desktop\\data\\quiet-labeled.epd", 100, 778, 100000, 0.1)

	// engine.RunCommLoop()

	/*var search engine.Search
	search.Setup(engine.FENKiwiPete)

	search.TT.Resize(engine.DefaultTTSize, engine.SearchEntrySize)
	search.Timer.Setup(
		engine.NoValue,
		engine.NoValue,
		7500,
		int16(engine.NoValue),
		engine.MaxDepth,
		math.MaxUint64,
	)

	fmt.Println()
	start := time.Now()
	bestMove := search.Search()
	elapsed := time.Since(start)

	fmt.Println("\nBest move:", bestMove)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	/*pos := engine.Position{}
	pos.LoadFEN(engine.FENStartPosition)

	fmt.Println()

	start := time.Now()
	nodes := engine.Perft(&pos, 7)
	elapsed := time.Since(start)

	fmt.Println("\nNodes:", nodes)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())
	fmt.Printf("Nps: %d\n", int(float64(nodes)/elapsed.Seconds()))*/

	/*
		cutechess-cli.exe ^
		-srand %RANDOM% ^
		-engine cmd=blunder.exe stderr=log.txt ^
		-engine cmd=blunder-old.exe ^
		-openings file=C:\%HOMEPATH%\Desktop\misc\2moves_v2a.pgn format=pgn order=random ^
		-each option.Hash=32 tc=inf/10+0.1 proto=uci ^
		-games 2 -rounds 1000 -repeat 2 ^
		-sprt elo0=1 elo1=5 alpha=0.05 beta=0.05 ^
		-concurrency 8 ^
		-ratinginterval 50 ^
		-recover

		cutechess-cli.exe ^
		-srand %RANDOM% ^
		-engine cmd=blunder.exe stderr=log.txt ^
		-engine cmd=blunder-5.0.0.exe ^
		-engine cmd=blunder-4.0.0.exe ^
		-openings file=C:\Users\deanm\Desktop\misc\2moves_v2a.pgn format=pgn order=random ^
		-each option.Hash=32 timemargin=60000 tc=inf/8+0.08 proto=uci ^
		-games 2 -rounds 100 -repeat 2 -maxmoves 200 ^
		-concurrency 8 ^
		-ratinginterval 50 ^
		-recover ^
		-tournament gauntlet
	*/
}
