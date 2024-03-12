package engine

/*********************************

	EVALUATION DRIVERS AND HELPERS

*********************************/

// static piece scores
var pieceScore = [12]int{
	100,    // white pawn
	300,    // white knight
	350,    // white bishop
	500,    // white rook
	1000,   // white queen
	10000,  // white king
	-100,   // black pawn
	-300,   // black knight
	-350,   // black bishop
	-500,   // black rook
	-1000,  // black queen
	-10000, // black king
}

// static positional scores
var pawnScore = [64]int{
	90, 90, 90, 90, 90, 90, 90, 90,
	30, 30, 30, 40, 40, 30, 30, 30,
	20, 20, 20, 30, 30, 30, 20, 20,
	10, 10, 10, 20, 20, 10, 10, 10,
	5, 5, 10, 20, 20, 5, 5, 5,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 0, -10, -10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var knightScore = [64]int{
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 10, 10, 0, 0, -5,
	-5, 5, 20, 20, 20, 20, 5, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 5, 20, 10, 10, 20, 5, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, -10, 0, 0, 0, 0, -10, -5,
}

var bishopScore = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 10, 0, 0, 0, 0, 10, 0,
	0, 30, 0, 0, 0, 0, 30, 0,
	0, 0, -10, 0, 0, -10, 0, 0,
}

var rookScore = [64]int{
	50, 50, 50, 50, 50, 50, 50, 50,
	50, 50, 50, 50, 50, 50, 50, 50,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 0, 20, 20, 0, 0, 0,
}

var kingScore = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 5, 5, 5, 5, 0, 0,
	0, 5, 5, 10, 10, 5, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 5, 5, -5, -5, 0, 5, 0,
	0, 0, 5, 0, -15, 0, 10, 0,
}

var mirrorScore = [64]Square{
	A1, B1, C1, D1, E1, F1, G1, H1,
	A2, B2, C2, D2, E2, F2, G2, H2,
	A3, B3, C3, D3, E3, F3, G3, H3,
	A4, B4, C4, D4, E4, F4, G4, H4,
	A5, B5, C5, D5, E5, F5, G5, H5,
	A6, B6, C6, D6, E6, F6, G6, H6,
	A7, B7, C7, D7, E7, F7, G7, H7,
	A8, B8, C8, D8, E8, F8, G8, H8,
}

// lookup table for MVV_LVA
var MVV_LVA = [12][12]int{
	{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
	{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
	{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
	{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
	{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
	{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},
	{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
	{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
	{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
	{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
	{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
	{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},
}

// get the rank from the square
var GetRank = [64]int{
	7, 7, 7, 7, 7, 7, 7, 7,
	6, 6, 6, 6, 6, 6, 6, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// penalty for doubling pawns
var DoublePawnPenalty = -10

// penalty for making pawns isolated
var IsolatedPawnPenalty = -10

// passed pawns are really good
var PassedPawnBonus = [8]int{
	0, 10, 30, 50, 75, 100, 150, 200,
}

// file masking
var FileMasks [64]Bitboard

// rank masking
var RankMasks [64]Bitboard

// isolated mask
var IsolatedMasks [64]Bitboard

// passed pawn mask
var WhitePassedMasks [64]Bitboard

var BlackPassedMasks [64]Bitboard

// constants
const MAX_PLY = 64

// lookup table for killer moves [index][ply]
var KillerMoves [2][MAX_PLY]int

// lookup table for history moves [piece][square]
var HistoryMoves [12][64]int

// PV nodes
var PVLength [MAX_PLY]int

// PV table
var PVTable [MAX_PLY][MAX_PLY]int

// if we are following PV line
var FollowPv int

// scoring the PV line
var ScorePv int

// define highest possible value
var Infinity int = 50000

// defining mate values
var MateValue int = 49000

// define mate score
var MateScore int = 48000

func EnablePVMoveScoring(moves Moves) {
	// disable PV following
	FollowPv = 0

	for i := 0; i < moves.Move_count; i++ {
		if PVTable[0][Ply] == moves.Move_list[i] {
			// enable move scoring
			ScorePv = 1

			// enable following pv
			FollowPv = 1
		}
	}
}

// using MVV_LVA
func ScoreMove(move int) int {
	_, target, piece, _, capture, _, _, _ := DecodeMove(move)

	// PV move scoring
	if ScorePv == 1 {
		if PVTable[0][Ply] == move {
			ScorePv = 0
			// give PV move the highest score
			return 20000
		}
	}

	if capture == 1 {
		var targetPiece = WhitePawn

		var startPiece, endPiece Piece

		if SideToMove == White {
			startPiece = BlackPawn
			endPiece = BlackKing
		} else {
			startPiece = WhitePawn
			endPiece = WhiteKing
		}
		for i := startPiece; i <= endPiece; i++ {
			if GameBoards[i].GetBit(target) != 0 {
				targetPiece = i
				break
			}
		}
		return MVV_LVA[piece][targetPiece] + 10000
	} else {
		if KillerMoves[0][Ply] == move { // 1st killer move
			return 9000
		} else if KillerMoves[1][Ply] == move { // 2nd killer move
			return 8000
		} else { // history move
			return HistoryMoves[piece][target]
		}
	}
}

func SortMoves(moves *Moves) {
	var moveScores []int
	// loop
	for i := 0; i < moves.Move_count; i++ {
		// score
		score := ScoreMove(moves.Move_list[i])
		// append
		moveScores = append(moveScores, score)
	}
	// sort
	for current := 0; current < moves.Move_count; current++ {
		for next := current + 1; next < moves.Move_count; next++ {
			if moveScores[current] < moveScores[next] {
				// swap the scores
				tempScore := moveScores[current]
				moveScores[current] = moveScores[next]
				moveScores[next] = tempScore
				// swap the moves
				tempMove := moves.Move_list[current]
				moves.Move_list[current] = moves.Move_list[next]
				moves.Move_list[next] = tempMove
			}
		}
	}
}

func Evaluate() int {
	// score of the position
	score := 0

	// looping variables
	var current_piece Bitboard
	var piece, square int

	// loop
	for i := WhitePawn; i <= BlackKing; i++ {
		// copy
		current_piece = GameBoards[i]
		// loop over copied piece
		for current_piece != 0 {
			piece = int(i)
			// pop from gameboard copy
			square = current_piece.LSBIndex()
			current_piece.PopBit(square)
			// add piece score
			score += pieceScore[piece]
			// add position score depending on piece
			switch i {
			case WhitePawn:
				score += pawnScore[square]
				// get double pawn penalty
				doublePawns := GameBoards[i] & FileMasks[square]
				if doublePawns.CountBits() > 1 {
					score += doublePawns.CountBits() * DoublePawnPenalty
				}
				// get isolated pawn penalty
				if GameBoards[i] & IsolatedMasks[square] == 0 {
					score += IsolatedPawnPenalty
				}

				// white passed pawn bonus 
				if WhitePassedMasks[square] & GameBoards[BlackPawn] == 0 {
					score += PassedPawnBonus[GetRank[square]]
				}
			case WhiteKnight:
				score += knightScore[square]
			case WhiteBishop:
				score += bishopScore[square]
			case WhiteRook:
				score += rookScore[square]
			case WhiteKing:
				score += kingScore[square]
			case BlackPawn:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				score -= pawnScore[mirrorScore[index_sq]]
				// get double pawn penalty
				doublePawns := GameBoards[i] & FileMasks[square]
				if doublePawns.CountBits() > 1 {
					score -= doublePawns.CountBits() * DoublePawnPenalty
				}
				// get isolated pawn penalty
				if GameBoards[i] & IsolatedMasks[square] == 0 {
					score -= IsolatedPawnPenalty
				}
				// give passed pawn bonus 
				// white passed pawn bonus 
				if BlackPassedMasks[square] & GameBoards[WhitePawn] == 0 {
					score -= PassedPawnBonus[GetRank[square]]
				}
			case BlackKnight:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				score -= knightScore[mirrorScore[index_sq]]
			case BlackBishop:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				score -= bishopScore[mirrorScore[index_sq]]
			case BlackRook:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				score -= rookScore[mirrorScore[index_sq]]
			case BlackKing:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				score -= kingScore[mirrorScore[index_sq]]
			}
		}
	}
	if SideToMove == White {
		return score
	} else {
		return -score
	}
}

// check for repetition on board
func IsRepetition() bool {
	for i := 0; i < RepetitionIndex; i++ {
		// if we found the hash, it is a repetition
		if RepetitionTable[i] == HashKey {
			return true
		}
	}
	// by default it is not a repetition
	return false
}

// set file or rank mask
func SetFileRankMask(fileNumber int, rankNumber int) Bitboard {
	// file or rank
	var mask Bitboard

	// loop over ranks and files
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			if fileNumber != -1 && file == fileNumber {
				mask.SetBit(sq)
			} else if rankNumber != -1 && rank == rankNumber {
				mask.SetBit(sq)
			}
		}
	}

	return mask
}

// set eval masks
func InitEvaluationMasks() {

	// INIT FILE MASKS
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			FileMasks[sq] |= SetFileRankMask(file, -1)

		}
	}
	// INIT RANK MASKS
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			RankMasks[sq] |= SetFileRankMask(-1, rank)
		}
	}

	// INIT ISOLATED MASKS
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			IsolatedMasks[sq] |= SetFileRankMask(file-1, -1)
			IsolatedMasks[sq] |= SetFileRankMask(file+1, -1)
		}
	}

	// White Passed Masks
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			WhitePassedMasks[sq] |= SetFileRankMask(file-1, -1)
			WhitePassedMasks[sq] |= SetFileRankMask(file, -1)
			WhitePassedMasks[sq] |= SetFileRankMask(file+1, -1)

			for i := 0; i < (7 - rank); i++ {
				WhitePassedMasks[sq] &= ^RankMasks[(7-i)*8+file]
			}
		}
	}

	// White Black Masks
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			BlackPassedMasks[sq] |= SetFileRankMask(file-1, -1)
			BlackPassedMasks[sq] |= SetFileRankMask(file, -1)
			BlackPassedMasks[sq] |= SetFileRankMask(file+1, -1)

			for i := 0; i < rank+1; i++ {
				BlackPassedMasks[sq] &= ^RankMasks[i*8+file]
			}
		}
	}
}
