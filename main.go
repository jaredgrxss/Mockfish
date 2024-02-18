package main

import (
	"Mockfish/engine"
)
 
func main() {
	engine.GeneratePieceAttacks()
	engine.ParseFen(engine.TRICKY_POSITION)
	
	// testing move generation
	engine.GenerateMoves()
}