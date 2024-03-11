package engine

import (
	"fmt"
)

/*
*******************************

	MAIN DRIVERS FOR SEARCH

*******************************
*/
var Ply int = 0
var NegamaxNodes = 0

const FULL_DEPTH_MOVES int = 4
const REDUCTION_LIMIT int = 3

var nodesSaved = 0

func Negamax(alpha int, beta int, depth int) int {
	// node score
	var score int

	// define hash flag for TT
	hashFlag := HashFlagAlpha

	// read cache
	if Ply > 0 &&
		score == ReadHashData(alpha, beta, depth) &&
		ReadHashData(alpha, beta, depth) != NO_HASH_ENTRY {
		nodesSaved++
		// if we have seen this move, return score without searching
		return score
	}

	// init pv length
	PVLength[Ply] = Ply

	// evaluate position and base case
	if depth == 0 {
		return Quiescence(alpha, beta)
	}

	if Ply > MAX_PLY-1 {
		return Evaluate()
	}

	// increment nodes count
	NegamaxNodes++

	// check for king check
	var inCheck bool
	if SideToMove == White {
		inCheck = IsSquareAttacked(GameBoards[WhiteKing].LSBIndex(), SideToMove^1)
	} else {
		inCheck = IsSquareAttacked(GameBoards[BlackKing].LSBIndex(), SideToMove^1)
	}

	if inCheck {
		depth++
	}

	// legal moves
	legalMoves := 0

	// NULL MOVE PRUNING
	if depth >= 3 && !inCheck && Ply > 0 {
		GameBoardsCopy := GameBoards
		GameOccupancyCopy := GameOccupancy
		SideToMoveCopy := SideToMove
		CastleCopy := Castle
		EnpassantCopy := Enpassant
		Ply++
		HashKeyCopy := HashKey

		// hash enpassant
		if Enpassant != 64 {
			HashKey ^= EnpassantKeys[Enpassant]
		}

		SideToMove ^= 1

		// hash side if it is black
		HashKey ^= SideKey

		Enpassant = 64

		// search moves with reduced depth to find beta cutoffs
		score := -Negamax(-beta, -beta+1, depth-1-2)

		Ply--
		GameBoards = GameBoardsCopy
		GameOccupancy = GameOccupancyCopy
		SideToMove = SideToMoveCopy
		Castle = CastleCopy
		Enpassant = EnpassantCopy
		HashKey = HashKeyCopy

		if score >= beta {
			// fails high
			return beta
		}
	}

	// generate moves
	var moves Moves
	GeneratePositionMoves(&moves)

	if FollowPv == 1 {
		EnablePVMoveScoring(moves)
	}

	// sort moves
	SortMoves(&moves)

	// used for LMR
	movesSearched := 0

	// loop over moves
	for i := 0; i < moves.Move_count; i++ {
		// decode for history
		_, target, piece, promo, capture, _, _, _ := DecodeMove(moves.Move_list[i])

		// copy board
		GameBoardsCopy := GameBoards
		GameOccupancyCopy := GameOccupancy
		SideToMoveCopy := SideToMove
		CastleCopy := Castle
		EnpassantCopy := Enpassant
		HashKeyCopy := HashKey
		Ply++

		// check to see if move was legal
		if MakeMove(moves.Move_list[i], 0) == 0 {
			Ply--
			continue
		}
		// increase legal move counter
		legalMoves++

		// LMR
		if movesSearched == 0 {
			// do a full normal search for non PV nodes
			score = -Negamax(-beta, -alpha, depth-1)
		} else {
			// condition to use LMR
			if movesSearched >= FULL_DEPTH_MOVES && depth >= REDUCTION_LIMIT && !inCheck && capture == 0 && promo == 0 {
				score = -Negamax(-alpha-1, -alpha, depth-2) // search w/ reduced depth
			} else {
				score = alpha + 1 // make sure full detph search is done
			}
			// PV search
			if score > alpha {
				score = -Negamax(-alpha-1, -alpha, depth-1) // make aspiration window smaller

				// if LMK fails, re-search at fll depth with normal aspiration bandwidth
				if (score > alpha) && (score < beta) {
					score = -Negamax(-beta, -alpha, depth-1)
				}
			}
		}

		// restore board
		GameBoards = GameBoardsCopy
		GameOccupancy = GameOccupancyCopy
		SideToMove = SideToMoveCopy
		Castle = CastleCopy
		Enpassant = EnpassantCopy
		HashKey = HashKeyCopy
		Ply--

		// increment number of moves searched
		movesSearched++

		// found a better move
		if score > alpha {
			// switch our hash flag
			hashFlag = HashFlagExact

			// store history moves on non capture moves
			if capture == 0 {
				HistoryMoves[piece][target] += depth
			}

			// PV node
			alpha = score

			// write PV move to PV table
			PVTable[Ply][Ply] = moves.Move_list[i]

			// loop next ply
			for nxtPly := Ply + 1; nxtPly < PVLength[Ply+1]; nxtPly++ {
				// copy moves from deeper ply
				PVTable[Ply][nxtPly] = PVTable[Ply+1][nxtPly]
			}

			// adjust PV length
			PVLength[Ply] = PVLength[Ply+1]

			// fail-hard beta cutoff
			if score >= beta {
				// store hash entry for TT
				WriteHashData(beta, depth, HashFlagBeta)

				// only store killer moves on non capture moves
				if capture == 0 {
					// store killer moves
					KillerMoves[1][Ply] = KillerMoves[0][Ply]
					KillerMoves[0][Ply] = moves.Move_list[i]
				}

				// move fails high
				return beta
			}
		}
	}
	// we have no legal moves
	if legalMoves == 0 {
		if inCheck {
			return -49000 + Ply // returning mate
		} else {
			return 0 // stalemate
		}
	}

	// write hash entry to store alpha
	WriteHashData(alpha, depth, hashFlag)

	// fails low
	return alpha
}

// quiescence search for horizon effect
func Quiescence(alpha int, beta int) int {
	// score of node
	var score int

	NegamaxNodes++

	if Ply > MAX_PLY-1 {
		return Evaluate()
	}

	eval := Evaluate()

	if eval >= beta {
		// move fails high
		return beta
	}

	if eval > alpha {
		// PV node
		alpha = eval
	}

	// generate moves
	var q_moves Moves
	GeneratePositionMoves(&q_moves)

	// sort moves with MVV_LVA heuristic
	SortMoves(&q_moves)

	// loop over moves
	for i := 0; i < q_moves.Move_count; i++ {
		// copy board
		GameBoardsCopy := GameBoards
		GameOccupancyCopy := GameOccupancy
		SideToMoveCopy := SideToMove
		CastleCopy := Castle
		EnpassantCopy := Enpassant
		HashKeyCopy := HashKey
		Ply++

		// check to see if move was legal
		if MakeMove(q_moves.Move_list[i], onlyCaptures) == 0 {
			Ply--
			continue
		}

		// current score
		score = -Quiescence(-beta, -alpha)

		// restore board
		GameBoards = GameBoardsCopy
		GameOccupancy = GameOccupancyCopy
		SideToMove = SideToMoveCopy
		Castle = CastleCopy
		Enpassant = EnpassantCopy
		HashKey = HashKeyCopy
		Ply--

		if score > alpha {
			// PV node
			alpha = score

			if score >= beta {
				// move fails high
				return beta
			}
		}
	}
	// fails low
	return alpha
}

func SearchPosition(depth int) {
	NegamaxNodes = 0

	// reset our tables and vars
	FollowPv = 0
	ScorePv = 0
	KillerMoves = [2][MAX_PLY]int{}
	HistoryMoves = [12][64]int{}
	PVTable = [MAX_PLY][MAX_PLY]int{}
	PVLength = [MAX_PLY]int{}

	// search variables
	alpha := -50000
	beta := 50000
	search_depth := 1

	// clear transposition table
	ClearTranspositionTable()

	// perform iterative deepening
	for {
		FollowPv = 1

		// find best move
		score := Negamax(alpha, beta, search_depth)

		/***************
			LOGGING
		***************/
		fmt.Printf("ITERATIVE DEEPENING: eval: %d depth: %d nodes: %d nodes saved: %d pv ", score, search_depth, NegamaxNodes, nodesSaved)
		for cnt := 0; cnt < PVLength[0]; cnt++ {
			// print move
			PrintUCICompatibleMove(PVTable[0][cnt])
			fmt.Printf(" ")
		}
		fmt.Printf("\n")

		search_depth++
		if search_depth > depth {
			break
		}
	}

	// best move printed
	fmt.Print("bestmove ")
	PrintUCICompatibleMove(PVTable[0][0])
	fmt.Println()
}
