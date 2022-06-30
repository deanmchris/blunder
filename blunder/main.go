package main

import (
	"blunder/engine"
)

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
	engine.InitEvalBitboards()
	engine.InitSearchTables()
}

func main() {
	engine.RunCommLoop()
}
