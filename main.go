package main

import (
	"Mockfish/engine"
	"fmt"
)
 
func main() {
	engine.GeneratePieceAttacks()
	var test engine.Bitboard
	test.SetBit(engine.A8); test.SetBit(engine.E6)
	test.PrintBoard()
	engine.GameBoards[engine.WhitePawn].SetBit(engine.E2)
	engine.GameBoards[engine.WhitePawn].PrintBoard()
	fmt.Println(string(engine.AsciiPieces[engine.BlackBishiop]))
	fmt.Println(engine.UnicodePieces[engine.BlackBishiop])
	fmt.Println(engine.UnicodePieces[engine.Ascii_to_Type["r"]])

}