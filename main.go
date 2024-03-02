package main

import (
	"Mockfish/engine"
)

func main() {
	engine.GeneratePieceAttacks()
	engine.ParseFen("r1b1k2r/pppp2pp/3b1p1n/2n1pQq1/P3P3/1PNB4/2PP1PPP/R1B1K1NR b KQkq - 0 1")
	engine.GeneratePositionMoves()
	engine.TestMakeMove()
}
