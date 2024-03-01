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
var nodes, captures, ep, castles, promotions int64 = 0, 0, 0, 0, 0
var startPosition string = TRICKY_POSITION

func getTime() int64 {
	return time.Now().UnixNano() / 1e6
}

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
	for move := 0; move < MoveList.move_count; move++ {
		_, _, _, promo, capture, _, enpassant, castling := DecodeMove(MoveList.move_list[move])
		captures += int64(capture)
		castles += int64(castling)
		ep += int64(enpassant)
		if promo > 0 {
			promotions++
		}
		// copy board position
		COPY()
		// move was illegal, don't make it
		if MakeMove(MoveList.move_list[move], allMoves) == 0 {
			continue
		}
		// step into
		perftDriver(depth - 1)
		// restore prev state
		RESTORE()
	}
}

func TestTrickyPosition(t *testing.T) {
	// precompute attack data for pieces
	GeneratePieceAttacks()
	// parse the starting position & setting time
	ParseFen(startPosition)
	start := getTime()
	// perft
	perftDriver(3)
	// results
	fmt.Printf("total time taken: %d ms\n", getTime() - start)
	fmt.Printf("total nodes reached: %d \n", nodes)
	fmt.Printf("other information: \ncaptures: %d\nenpassant moves: %d\ncastles: %d\npromotions: %d\n", captures, ep, castles, promotions)
}

func TestInitialPosition(t *testing.T) {

}