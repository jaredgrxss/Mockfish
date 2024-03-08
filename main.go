package main

import (
	"Mockfish/engine"
	// "fmt"
)

var debug int = 0

func main() {
	engine.GeneratePieceAttacks()
	if debug != 1 {
		engine.RunUCI()
	}
}
