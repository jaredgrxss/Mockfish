package main

import (
	"Mockfish/engine"
)

var debug int = 0

func main() {
	// precompute
	engine.GeneratePieceAttacks()
	engine.ParseFen(engine.TRICKY_POSITION)
	engine.PrintGameboard()
	engine.SearchPosition(7)
	if debug != 1 {
		engine.RunUCI()
	}
}
