package engine

import (
	"fmt"
)

/*
*******************************

	MAIN DRIVERS FOR SEARCH

*******************************
*/
var ply int = 0
var BestMove int
var NegamaxNodes = 0

func Negamax(alpha int, beta int, depth int) int {
	// evaluate position
	if depth == 0 {
		return Evaluate()
	}
	// increment nodes count
	NegamaxNodes++

	// legal moves
	legalMoves := 0

	// check for king check
	var inCheck bool
	if SideToMove == White {
		inCheck = IsSquareAttacked(GameBoards[WhiteKing].LSBIndex(), SideToMove^1)
	} else {
		inCheck = IsSquareAttacked(GameBoards[BlackKing].LSBIndex(), SideToMove^1)
	}

	// store our best move found searching so far
	var BestSoFar int

	// old value of alpha
	var oldAlpha = alpha

	// generate moves
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
		legalMoves++

		score := -Negamax(-beta, -alpha, depth-1)
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
	// we have no legal moves
	if legalMoves == 0 {
		if inCheck {
			return -49000 + ply
		} else {
			return 0
		}
	}

	// found better move than previous
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
