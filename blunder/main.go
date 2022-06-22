package main

import "blunder/engine"

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
}

func main() {
	engine.RunCommLoop()
}
