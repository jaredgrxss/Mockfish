package main

import (
	"Mockfish/engine"
)
 
func main() {
	engine.GeneratePieceAttacks()
	engine.ParseFen(engine.TRICKY_POSITION)
	engine.PrintGameboard()
}