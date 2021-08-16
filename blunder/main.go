package main

import (
	"blunder/engine"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/*var s engine.Search
	s.Pos.LoadFEN(engine.FENStartPosition)
	s.Timer.TimeLeft = 600000
	s.TransTable.Resize(engine.DefaultTTSize)

	start := time.Now()
	m := s.Search()
	fmt.Println("Bestmove:", m)
	elapsed := time.Since(start)
	fmt.Printf("Time: %vms\n", elapsed.Milliseconds())*/

	defer func() {
		if r := recover(); r != nil {
			println(fmt.Sprintf("Internal error: %v", r))
			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, true)
			println(fmt.Sprintf("%s\n", string(buf[0:stackSize])))
		}
	}()
	engine.UCILoop()
}
