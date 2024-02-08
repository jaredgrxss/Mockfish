package main

import (
	"Mockfish/engine"
)
 
func main() {
	engine.GeneratePieceAttacks()
	// engine.ParseFen("8/8/4R3/3b4/8/8/8/8 w - - ")
	engine.ParseFen(engine.TRICKY_POSITION)
	engine.PrintGameboard()
	engine.PrintAttackedSquares(engine.Black)
}