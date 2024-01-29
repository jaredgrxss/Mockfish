package main

import (
	"Mockfish/engine"
	// "fmt"
)
 
func main() {
	engine.GeneratePieceAttacks()
	engine.SetInitBoardState()
	engine.PrintGameboard()
	engine.ParseFen(engine.START_POSITION)
	engine.PrintGameboard()
}