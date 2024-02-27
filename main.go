package main

import (
	"Mockfish/engine"
)

func main() {
	engine.GeneratePieceAttacks()
	engine.ParseFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1 ")
	engine.PrintGameboard()
	// testing move generation
	engine.GeneratePositionMoves()
	engine.MoveList.PrintMoveList()
}
