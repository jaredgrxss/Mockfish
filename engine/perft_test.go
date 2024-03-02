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
var startPosition string = TRICKY_POSITION; var perft_depth int = 2

// main perft function
func perftDriver(depth int) {
	// we've reached a leaf
	if depth == 0 {
		nodes++
		return 
	}
	// generate a positions moves
	GeneratePositionMoves()
	// loop over these moves
	for i := 0; i < MoveList.move_count; i++ {
		// copy board position
		COPY()
		// move was illegal, don't make it
		if MakeMove(MoveList.move_list[i], allMoves) == 0 {
			continue
		}
		// step into
		perftDriver(depth - 1)
		// restore prev state
		RESTORE()
	}
}

func TestPerft(t *testing.T) {
	// precompute attack data for pieces
	GeneratePieceAttacks()
	// parse the starting position & setting time
	ParseFen(startPosition)
	GeneratePositionMoves()
	start := getTime()
	fmt.Printf("Starting perft test....\n")
	for i := 0; i < 2; i++ {
		source, target, _, promo, _, _, _, _ := DecodeMove(MoveList.move_list[i])
		var promo_to_print = ""
		if (promo != 0) {
			promo_to_print = string(PromotedPieces[promo])
		}
		fmt.Printf("move: %s%s%s  side: %d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, SideToMove)
		// copy board position
		COPY()
		// move was illegal, don't make it
		if MakeMove(MoveList.move_list[i], allMoves) == 0 {
			nonLegal++
			continue
		}
		cummulative_nodes := nodes;

		// step into
		perftDriver(perft_depth - 1)

		old_nodes := nodes - cummulative_nodes

		// restore prev state
		RESTORE()
		fmt.Printf("move: %s%s%s  side: %d nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, SideToMove, old_nodes)
	}
	// results
	fmt.Printf("total time taken: %d ms\n", getTime() - start)
	fmt.Printf("total nodes reached: %d \n", nodes)
}




/*************
	HELPERS
*************/
func getTime() int64 {
	return time.Now().UnixNano() / 1e6
}