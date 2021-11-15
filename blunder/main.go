package main

import "blunder/engine"

func main() {
	var inter engine.UCIInterface
	inter.UCILoop()
	// remove KSE, add PV lines and checks to qsearch, cleaned-up code,
	// Remove king safety code; add PV & checks to qsearch; add UCI options; tweak & clean up code;
}
