package engine

import (
	"fmt"
	"testing"
	"time"
)

/***********************************************************************
	TEST SUITE FOR PERFT FUNCTIONALITY FOR SEARCH / MOVE GENERATION
***********************************************************************/

// to test different positions, put fen below
var nodes, captures, ep, castles, promotions, legal, nonLegal int64 = 0, 0, 0, 0, 0, 0, 0
var startPosition string = TEST_PERFTING; var perft_depth int = 5

// main perft function
func perftDriver(depth int) {
	// we've reached a leaf
	if depth == 0 {
		nodes++
		return 
	}
	var moves Moves
	// generate a positions moves
	GeneratePositionMoves(&moves)

	// loop over these moves
	for i := 0; i < moves.Move_count; i++ {
		_, _, _, promo, capture, _, enp, castle := DecodeMove(moves.Move_list[i])
		
		// copy board position
		//COPY()
		var GameBoards_Copy [12]Bitboard
		var GameOccupancy_Copy [3]Bitboard
		var SideToMove_Copy, Enpassant_Copy, Castle_Copy int
		GameBoards_Copy = GameBoards
		GameOccupancy_Copy = GameOccupancy
		SideToMove_Copy = SideToMove
		Enpassant_Copy = Enpassant
		Castle_Copy = Castle
		
		// move was illegal, don't make it
		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}


		// step into
		perftDriver(depth - 1)

		// restore prev state
		// RESTORE()
		ep += int64(enp)
		captures += int64(capture)
		castles += int64(castle)
		if promo > 0 {
			promotions++
		}
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy
	}
}

func TestPerft(t *testing.T) {
	GeneratePieceAttacks()
	var moves Moves 
	ParseFen(startPosition)
	GeneratePositionMoves(&moves)

	start := getTime()
	fmt.Printf("Starting perft test....\n")
	for i := 0; i < moves.Move_count; i++ {
		source, target, _, promo, capture, _, enp, castle := DecodeMove(moves.Move_list[i])
		
		// fmt.Println("GAMEBOARD IN:")
		// PrintGameboard()
		var promo_to_print = ""
		if (promo != 0) {
			promo_to_print = string(PromotedPieces[promo])
		}
		// COPY()
		var GameBoards_Copy [12]Bitboard
		var GameOccupancy_Copy [3]Bitboard
		var SideToMove_Copy, Enpassant_Copy, Castle_Copy int
		GameBoards_Copy = GameBoards
		GameOccupancy_Copy = GameOccupancy
		SideToMove_Copy = SideToMove
		Enpassant_Copy = Enpassant
		Castle_Copy = Castle

		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		
		cummulative_nodes := nodes;
		perftDriver(perft_depth - 1)

		old_nodes := nodes - cummulative_nodes
		ep += int64(enp)
		captures += int64(capture)
		castles += int64(castle)
		if promo > 0 {
			promotions++
		}

		// RESTORE()
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy
		// fmt.Println("GAMEBOARD OUT:")
		// PrintGameboard()
		fmt.Printf("move: %s%s%s  side: %d nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, SideToMove, old_nodes)
	}
	// results
	fmt.Printf("total time taken: %d ms\n", getTime() - start)
	fmt.Printf("total nodes reached: %d \n", nodes)
	fmt.Printf("total captures: %d\ntotal enpassant: %d\ntotal castles: %d\ntotal promotions: %d\n", captures, ep, castles, promotions)
}




/*************
	HELPERS
*************/
func getTime() int64 {
	return time.Now().UnixNano() / 1e6
}