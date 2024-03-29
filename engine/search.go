package engine

import (
	"fmt"
	"time"
)

/*
*******************************

	MAIN DRIVERS FOR SEARCH

*******************************
*/
var Ply int = 0
var NegamaxNodes = 0
var movesInQuiesence = 0
var transposition int = 0

const FULL_DEPTH_MOVES int = 4
const REDUCTION_LIMIT int = 3

func Negamax(alpha int, beta int, depth int) int {
	// node score
	score := ReadHashData(alpha, beta, depth)

	// define hash flag for TT
	hashFlag := HashFlagAlpha

	// check for repetition
	if Ply > 0 && IsRepetition() {
		// return draw
		return 0
	}

	// figure out if this is a pv node
	pvNode := beta-alpha > 1

	// return cached score
	if Ply > 0 &&
		score != NO_HASH_ENTRY &&
		!pvNode {
		transposition++
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
		// copy board
		state := COPY()
		// increment our ply
		Ply++

		// increment repetition index & store hash key
		RepetitionIndex++
		RepetitionTable[RepetitionIndex] = HashKey

		// hash enpassant
		if Enpassant != 64 {
			HashKey ^= EnpassantKeys[Enpassant]
		}

		SideToMove ^= 1

		// hash side if it is black
		HashKey ^= SideKey

		Enpassant = 64

		// search moves with reduced depth to find beta cutoffs
		score := -Negamax(-beta, -beta+1, depth-3)

		Ply--

		// decrement repetition index
		RepetitionIndex--

		RESTORE(state)

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
		state := COPY()
		Ply++

		// increment repetition index & store hash key
		RepetitionIndex++
		RepetitionTable[RepetitionIndex] = HashKey

		// check to see if move was legal
		if MakeMove(moves.Move_list[i], 0) == 0 {
			Ply--
			// decrement repetition index
			RepetitionIndex--
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
			if movesSearched >= FULL_DEPTH_MOVES &&
				depth >= REDUCTION_LIMIT &&
				!inCheck &&
				capture == 0 && promo == 0 {
				score = -Negamax(-alpha-1, -alpha, depth-2) // search w/ reduced depth
			} else {
				score = alpha + 1 // make sure full detph search is done
			}
			// PV search
			if score > alpha {
				score = -Negamax(-alpha-1, -alpha, depth-1) // make aspiration window smaller

				// if LMK fails, re-search at full depth with normal aspiration bandwidth
				if (score > alpha) && (score < beta) {
					score = -Negamax(-beta, -alpha, depth-1)
				}
			}
		}

		// restore board
		RESTORE(state)
		Ply--
		// decrement repetition index
		RepetitionIndex--

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
			return -MateValue + Ply // returning mate
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
		state := COPY()
		Ply++

		// increment repetition index & store hash key
		RepetitionIndex++
		RepetitionTable[RepetitionIndex] = HashKey

		// check to see if move was legal
		if MakeMove(q_moves.Move_list[i], onlyCaptures) == 0 {
			Ply--
			// decrement repetition index
			RepetitionIndex--
			continue
		}
		movesInQuiesence++
		// current score
		score := -Quiescence(-beta, -alpha)

		// restore board
		RESTORE(state)
		Ply--

		// decrement repetition index
		RepetitionIndex--

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
	alpha := -Infinity
	beta := Infinity
	search_depth := 1

	// clear transposition table
	ClearTranspositionTable()

	// perform iterative deepening
	for {
		FollowPv = 1
		startTime := GetTime()
		// find best move
		score := Negamax(alpha, beta, search_depth)

		if score <= alpha || score >= beta {
			alpha = -Infinity
			beta = Infinity
			continue
		}
		alpha = score - 100
		beta = score + 100

		/***************
			LOGGING
		***************/
		if score > -MateValue && score < -MateScore {
			fmt.Printf("info score mate %d depth %d nodes %d time %d pv ", -(score+MateValue)/2-1, search_depth, NegamaxNodes, GetTime()-startTime)
		} else if score > MateScore && score < MateValue {
			fmt.Printf("info score mate %d depth %d nodes %d time %d pv ", (MateValue-score)/2+1, search_depth, NegamaxNodes, GetTime()-startTime)
		} else {
			fmt.Printf("cp: %d depth: %d nodes: %d moves in qsearch: %d cache hits: %d time: %dms pv ", score, search_depth, NegamaxNodes, movesInQuiesence, transposition, GetTime()-startTime)
		}
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

func GetTime() int64 {
	return time.Now().UnixNano() / 1e6
}
