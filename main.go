package main

import (
	"Mockfish/engine"
	// "fmt"
)

func main() {
	engine.GeneratePieceAttacks()
	engine.ParseUCIPosition("position startpos moves e2e4 e7e5 b1c3 g8f6 d1g4 g7g236 g4g6 h7g6")
	engine.PrintGameboard()
}
