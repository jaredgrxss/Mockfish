package engine

import (
	"fmt"
)

/********************************
	MAIN DRIVERS FOR SEARCH
********************************/
var ply int
var BestMove int
var NegamaxNodes = 0

func Negamax(alpha int, beta int, depth int) int {
	// evaluate position
	if depth == 0 {
		return Evaluate()
	}
	// store our best move found searching so far
	var BestSoFar int

	// old value of alpha
	var oldAlpha = alpha

	// incriment nodes & generate moves
	NegamaxNodes++
	var moves Moves
	GeneratePositionMoves(&moves)

	// loop over moves
	for i := 0; i < moves.Move_count; i++ {
		// copy board
		GameBoardsCopy := GameBoards
		GameOccupancyCopy := GameOccupancy
		SideToMoveCopy := SideToMove
		CastleCopy := Castle 
		EnpassantCopy := Enpassant
		ply++
		// check to see if move was legal
		if MakeMove(moves.Move_list[i], allMoves) == 0 {
			ply--
			continue
		}
		score := -Negamax(-beta, -alpha, depth - 1)
		// restore board
		GameBoards = GameBoardsCopy
		GameOccupancy = GameOccupancyCopy
		SideToMove = SideToMoveCopy
		Castle = CastleCopy
		Enpassant = EnpassantCopy
		ply--

		if score >= beta {
			// move fails high
			return beta
		}
		
		if score > alpha {
			// PV node
			alpha = score
			if ply == 0 {
				BestSoFar = moves.Move_list[i]
			}
		}
	}
	// found better move
	if oldAlpha != alpha {
		BestMove = BestSoFar
	}

	// fails low
	return alpha
}

func SearchPosition(depth int) {
	// find best move
	Negamax(-50000, 50000, depth)
	// best move printed
	fmt.Print("bestmove ")
	PrintUCICompatibleMove(BestMove)
	fmt.Println()
}
