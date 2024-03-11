package engine

import (
	"fmt"
	"testing"
	"time"
)

/*
**********************************************************************

	TEST SUITE FOR PERFT FUNCTIONALITY FOR SEARCH / MOVE GENERATION

**********************************************************************
*/

var PERFT_TEST1 = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var PERFT_TEST2 = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
var PERFT_TEST3 = "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - -"
var PERFT_TEST4 = "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
var PERFT_TEST5 = "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"
var PERFT_TEST6 = "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"

// to test different positions, put fen below
var nodes, captures, ep, castles, promotions int64 = 0, 0, 0, 0, 0

// main perft function
func PerftDriver(depth int, move int) {
	// we've reached a leaf
	if depth == 0 {
		_, _, _, promo, capture, _, enp, castling := DecodeMove(move)
		captures += int64(capture) 
		ep += int64(enp) 
		castles += int64(castling)
		if promo != 0 {
			promotions += 1
		}
		nodes++
		return
	}
	var moves Moves
	// generate a positions moves
	GeneratePositionMoves(&moves)

	// loop over these moves
	for i := 0; i < moves.Move_count; i++ {
		// _, _, _, _, _, _, _, _ := DecodeMove(moves.Move_list[i])

		// copy board position
		GameBoards_Copy := GameBoards
		GameOccupancy_Copy := GameOccupancy
		SideToMove_Copy := SideToMove
		Enpassant_Copy := Enpassant
		Castle_Copy := Castle

		// move was illegal, don't make it
		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		// step into
		PerftDriver(depth - 1, moves.Move_list[i])

		// restore board position
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy
	}
}

func TestPerft(t *testing.T) {
	// reset variables
	RESET_TEST()
	depth := 5

	// precomputation
	GeneratePieceAttacks()
	var moves Moves
	ParseFen(PERFT_TEST1)
	GeneratePositionMoves(&moves)

	start := getTime()
	fmt.Printf("Starting perft test....\n")
	for i := 0; i < moves.Move_count; i++ {
		source, target, _, promo, _, _, _, _ := DecodeMove(moves.Move_list[i])
		var promo_to_print = ""
		if promo != 0 {
			promo_to_print = string(PromotedPieces[promo])
		}

		// copy board position
		GameBoards_Copy := GameBoards
		GameOccupancy_Copy := GameOccupancy
		SideToMove_Copy := SideToMove
		Enpassant_Copy := Enpassant
		Castle_Copy := Castle

		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		cummulative_nodes := nodes
		PerftDriver(depth - 1,  moves.Move_list[i])

		old_nodes := nodes - cummulative_nodes

		// restore
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy
		fmt.Printf("move: %s%s%s nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, old_nodes)
	}
	if nodes != 4865609 {
		t.Fatalf("Expected: 4,865,609		Found: %d\n", nodes)
	}
	// results
	fmt.Printf("\n	Depth: %d\n", depth)
	fmt.Printf("	Total Time: %d ms\n", getTime()-start)
	fmt.Printf("	Total Nodes: %d \n", nodes)
	fmt.Printf("	Total Captures: %d, Total Ep: %d, Total Castles: %d, Total Promotions: %d \n", captures, ep, castles, promotions)

}

func TestPerft2(t *testing.T) {
	// constants for position
	depth := 5
	RESET_TEST()

	// precomputation
	GeneratePieceAttacks()
	var moves Moves
	ParseFen(PERFT_TEST2)
	GeneratePositionMoves(&moves)

	// start test
	start := getTime()
	fmt.Printf("Starting PERFT 2 TEST....\n")
	for i := 0; i < moves.Move_count; i++ {
		source, target, _, promo, _, _, _, _ := DecodeMove(moves.Move_list[i])
		var promo_to_print = ""
		if promo != 0 {
			promo_to_print = string(PromotedPieces[promo])
		}

		// copy board position
		GameBoards_Copy := GameBoards
		GameOccupancy_Copy := GameOccupancy
		SideToMove_Copy := SideToMove
		Enpassant_Copy := Enpassant
		Castle_Copy := Castle

		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		cummulative_nodes := nodes
		PerftDriver(depth - 1,  moves.Move_list[i])

		old_nodes := nodes - cummulative_nodes

		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy

		fmt.Printf("move: %s%s%s nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, old_nodes)
	}
	if nodes != 193690690 {
		t.Fatalf("Error nodes expected: 193,690,690  but found: %d", nodes)
	}
	// results
	fmt.Printf("\n	Depth: %d\n", depth)
	fmt.Printf("	Total Time: %d ms\n", getTime()-start)
	fmt.Printf("	Total Nodes: %d \n", nodes)
	fmt.Printf("	Total Captures: %d, Total Ep: %d, Total Castles: %d, Total Promotions: %d \n", captures, ep, castles, promotions)

}

func TestPerft3(t *testing.T) {
	// constants for position
	depth := 6
	RESET_TEST()

	// precomputation
	GeneratePieceAttacks()
	var moves Moves
	ParseFen(PERFT_TEST3)
	GeneratePositionMoves(&moves)

	// start perft test
	start := getTime()
	fmt.Printf("Starting PERFT 3 TEST....\n")
	for i := 0; i < moves.Move_count; i++ {
		source, target, _, promo, _, _, _, _ := DecodeMove(moves.Move_list[i])
		var promo_to_print = ""
		if promo != 0 {
			promo_to_print = string(PromotedPieces[promo])
		}

		// copy board position
		GameBoards_Copy := GameBoards
		GameOccupancy_Copy := GameOccupancy
		SideToMove_Copy := SideToMove
		Enpassant_Copy := Enpassant
		Castle_Copy := Castle

		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		cummulative_nodes := nodes
		PerftDriver(depth - 1, moves.Move_list[i])

		old_nodes := nodes - cummulative_nodes

		// restore board state
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy

		fmt.Printf("move: %s%s%s nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, old_nodes)
	}
	if nodes != 11030083 {
		t.Fatalf("Error nodes expected: 11,030,083  but found: %d", nodes)
	}
	// results
	fmt.Printf("\n	Depth: %d\n", depth)
	fmt.Printf("	Total Time: %d ms\n", getTime()-start)
	fmt.Printf("	Total Nodes: %d \n", nodes)
	fmt.Printf("	Total Captures: %d, Total Ep: %d, Total Castles: %d, Total Promotions: %d \n", captures, ep, castles, promotions)

}

func TestPerft4(t *testing.T) {
	// constants for position
	depth := 5
	RESET_TEST()

	// precomputation
	GeneratePieceAttacks()
	var moves Moves
	ParseFen(PERFT_TEST4)
	GeneratePositionMoves(&moves)

	// start perft test
	start := getTime()
	fmt.Printf("Starting PERFT 4 TEST....\n")
	for i := 0; i < moves.Move_count; i++ {
		source, target, _, promo, _, _, _, _ := DecodeMove(moves.Move_list[i])
		var promo_to_print = ""
		if promo != 0 {
			promo_to_print = string(PromotedPieces[promo])
		}

		// copy board position
		GameBoards_Copy := GameBoards
		GameOccupancy_Copy := GameOccupancy
		SideToMove_Copy := SideToMove
		Enpassant_Copy := Enpassant
		Castle_Copy := Castle

		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		cummulative_nodes := nodes
		PerftDriver(depth - 1,  moves.Move_list[i])

		old_nodes := nodes - cummulative_nodes

		// restore board state
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy

		fmt.Printf("move: %s%s%s nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, old_nodes)
	}
	if nodes != 15833292 {
		t.Fatalf("Error nodes expected: 15,833,292  but found: %d", nodes)
	}
	// results
	fmt.Printf("\n	Depth: %d\n", depth)
	fmt.Printf("	Total Time: %d ms\n", getTime()-start)
	fmt.Printf("	Total Nodes: %d \n", nodes)
	fmt.Printf("	Total Captures: %d, Total Ep: %d, Total Castles: %d, Total Promotions: %d \n", captures, ep, castles, promotions)

}

func TestPerft5(t *testing.T) {
	// constants for position
	depth := 5
	RESET_TEST()

	// precomputation
	GeneratePieceAttacks()
	var moves Moves
	ParseFen(PERFT_TEST5)
	GeneratePositionMoves(&moves)

	// start perft test
	start := getTime()
	fmt.Printf("Starting PERFT 5 TEST....\n")
	for i := 0; i < moves.Move_count; i++ {
		source, target, _, promo, _, _, _, _ := DecodeMove(moves.Move_list[i])
		var promo_to_print = ""
		if promo != 0 {
			promo_to_print = string(PromotedPieces[promo])
		}

		// copy board position
		GameBoards_Copy := GameBoards
		GameOccupancy_Copy := GameOccupancy
		SideToMove_Copy := SideToMove
		Enpassant_Copy := Enpassant
		Castle_Copy := Castle

		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		cummulative_nodes := nodes
		PerftDriver(depth - 1,  moves.Move_list[i])

		old_nodes := nodes - cummulative_nodes

		// restore board state
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy

		fmt.Printf("move: %s%s%s nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, old_nodes)
	}
	if nodes != 89941194 {
		t.Fatalf("Error nodes expected: 89,941,194  but found: %d", nodes)
	}
	// results
	fmt.Printf("\n	Depth: %d\n", depth)
	fmt.Printf("	Total Time: %d ms\n", getTime()-start)
	fmt.Printf("	Total Nodes: %d \n", nodes)
	fmt.Printf("	Total Captures: %d, Total Ep: %d, Total Castles: %d, Total Promotions: %d \n", captures, ep, castles, promotions)

}

func TestPerft6(t *testing.T) {
	// constants for position
	depth := 5
	RESET_TEST()

	// precomputation
	GeneratePieceAttacks()
	var moves Moves
	ParseFen(PERFT_TEST6)
	GeneratePositionMoves(&moves)

	// start perft test
	start := getTime()
	fmt.Printf("Starting PERFT 6 TEST....\n")
	for i := 0; i < moves.Move_count; i++ {
		source, target, _, promo, _, _, _, _ := DecodeMove(moves.Move_list[i])
		var promo_to_print = ""
		if promo != 0 {
			promo_to_print = string(PromotedPieces[promo])
		}

		// copy board position
		GameBoards_Copy := GameBoards
		GameOccupancy_Copy := GameOccupancy
		SideToMove_Copy := SideToMove
		Enpassant_Copy := Enpassant
		Castle_Copy := Castle

		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			continue
		}

		cummulative_nodes := nodes
		PerftDriver(depth - 1,  moves.Move_list[i])

		old_nodes := nodes - cummulative_nodes

		// restore board state
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy

		fmt.Printf("move: %s%s%s nodes: %1d\n", IntSquareToString[source], IntSquareToString[target], promo_to_print, old_nodes)
	}
	if nodes != 164075551 {
		t.Fatalf("Error nodes expected: 164,075,551  but found: %d", nodes)
	}
	// results
	fmt.Printf("\n	Depth: %d\n", depth)
	fmt.Printf("	Total Time: %d ms\n", getTime()-start)
	fmt.Printf("	Total Nodes: %d \n", nodes)
	fmt.Printf("	Total Captures: %d, Total Ep: %d, Total Castles: %d, Total Promotions: %d \n", captures, ep, castles, promotions)
}

/*
************

	HELPERS

************
*/
func getTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func RESET_TEST() {
	nodes = 0
	captures = 0
	ep = 0
	castles = 0
	promotions = 0
}
