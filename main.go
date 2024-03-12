package main

import (
	"Mockfish/engine"
	"fmt"
)

var debug int = 1

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
	fmt.Println()
	/**********************
		END OF REQUIRED
	**********************/

	// // SEARCH TEST
	engine.ParseFen(engine.START_POSITION)
	engine.PrintGameboard()
	engine.SearchPosition(10)
	// // // NODES BEOFRE TRANSPOSITION TABLE
	// // // 3407991
	// // // NODES AFTER TRANSPOSITION TABLE
	// // //
	// fmt.Println("HERE")
	// TESTING_ZOBRIST(4)

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
