package main

import (
	"Mockfish/engine"
	// "fmt"
)

var debug int = 0

func main() {
	/******************
		REQUIRED
	******************/
	// precompute attacks
	engine.GeneratePieceAttacks()
	// clear transposition table
	engine.ClearTranspositionTable()
	// precompute hash keys
	engine.InitZobrist()
	// improving evaulation masks
	engine.InitEvaluationMasks()
	/**********************
		END OF REQUIRED
	**********************/

	// // SEARCH TEST
	// testPos := "1k6/8/8/6Q1/1K6/8/8/8 b - -"
	// engine.ParseFen(engine.START_POSITION)
	// engine.PrintGameboard()
	// fmt.Println(engine.Evaluate())
	// engine.SearchPosition(10)

	if debug != 1 {
		engine.RunUCI()
	}
}

// func TESTING_ZOBRIST(depth int) {
// 	// we've reached a leaf
// 	if depth == 0 {
// 		return
// 	}
// 	var moves engine.Moves
// 	// generate a positions moves
// 	engine.GeneratePositionMoves(&moves)

// 	// loop over these moves
// 	for i := 0; i < moves.Move_count; i++ {
// 		// _, _, _, _, _, _, _, _ := DecodeMove(moves.Move_list[i])

// 		// copy board position
// 		GameBoards_Copy := engine.GameBoards
// 		GameOccupancy_Copy := engine.GameOccupancy
// 		SideToMove_Copy := engine.SideToMove
// 		Enpassant_Copy := engine.Enpassant
// 		Castle_Copy := engine.Castle
// 		HashKeyCopy := engine.HashKey

// 		// move was illegal, don't make it
// 		if engine.MakeMove(moves.Move_list[i], 0) == 0 {
// 			continue
// 		}

// 		// step into
// 		TESTING_ZOBRIST(depth - 1)

// 		// restore board position
// 		engine.GameBoards = GameBoards_Copy
// 		engine.GameOccupancy = GameOccupancy_Copy
// 		engine.SideToMove = SideToMove_Copy
// 		engine.Enpassant = Enpassant_Copy
// 		engine.Castle = Castle_Copy
// 		engine.HashKey = HashKeyCopy
// 	}
// }
