package main

import (
	"Mockfish/engine"
)

func main() {
	engine.GeneratePieceAttacks()
	engine.ParseFen(engine.TRICKY_POSITION)
	// engine.GameBoards[engine.WhitePawn].SetBit(int(engine.B7))
	engine.TestMakeMove()
}
