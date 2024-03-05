package main

import (
	"Mockfish/engine"
	"fmt"
)

func main() {
	engine.GeneratePieceAttacks()
	engine.ParseFen(engine.START_POSITION)
	move := engine.ParseUCIMove("b2b6")
	if move != 0 {
		engine.MakeMove(move, 0)
		engine.PrintGameboard()
	} else {
		fmt.Println("illegal move!")
	}
}
