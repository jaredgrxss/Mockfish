package main

import (
	"Mockfish/engine"
	// "fmt"
)

func main() {
	engine.GeneratePieceAttacks()
	engine.ParseUCIPosition("position startpos")
	engine.PrintGameboard()
	engine.RunUCI()
}
