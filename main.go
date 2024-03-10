package main

import (
	"Mockfish/engine"
)

var debug int = 0

func main() {
	// precompute
	engine.GeneratePieceAttacks()
	// engine.ParseFen(engine.TRICKY_POSITION)
	// // SEARCH TEST
	// engine.PrintGameboard()
	// engine.SearchPosition(9)


	if debug != 1 {
		engine.RunUCI()
	}
}
