package main

import (
	"Mockfish/engine"
	
)

func main() {

	var board engine.Bitboard = 0
	board.PrintBoard()
	board.SetBit("e4")
	board.SetBit("f2")
	board.SetBit("a8")
	board.PrintBoard()
	board.PopBit("a8")
	board.PopBit("a8")
	board.PrintBoard()
}