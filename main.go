package main

import (
	"Mockfish/engine"
)
 
func main() {
	engine.GeneratePieceAttacks()
	// engine.ParseFen(engine.TRICKY_POSITION)
	// engine.PrintGameboard()

	var occupancy engine.Bitboard
	occupancy.SetBit(int(engine.E6)); occupancy.SetBit(int(engine.B7))
	occupancy.SetBit(int(engine.B5)); occupancy.SetBit(int(engine.G2))
	var result = engine.GetQueenAttack(int(engine.D5), occupancy)
	result.PrintBitboard()
}