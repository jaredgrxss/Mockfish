package engine

/*********************************

	EVALUATION DRIVERS AND HELPERS

*********************************/

// static piece scores
// var pieceScore = [12]int{
// 	100,    // white pawn
// 	300,    // white knight
// 	350,    // white bishop
// 	500,    // white rook
// 	1000,   // white queen
// 	10000,  // white king
// 	-100,   // black pawn
// 	-300,   // black knight
// 	-350,   // black bishop
// 	-500,   // black rook
// 	-1000,  // black queen
// 	-10000, // black king
// }

// game phase scores
var pieceScore = [2][12]int{
	// opening material scores
	{82, 337, 365, 447, 1025, 12000, -82, -337, -365, -477, -1025, -12000},
	// endgame material scores
	{94, 281, 297, 512, 936, 12000, -94, -281, -297, -512, -936, -12000},
}

// game phase scores
const openingPhaseScore = 6192
const endGamePhaseScore = 518

// game phases
const (
	opening int = iota
	endgame
	middlegame
)

// constants for lookups
const (
	evalPawn int = iota
	evalKnight
	evalBishop
	evalRook
	evalQueen
	evalKing
)

var postionalScores = [2][6][64]int{
	// opening game positional scores
	{
		//pawn
		{0, 0, 0, 0, 0, 0, 0, 0,
			98, 134, 61, 95, 68, 126, 34, -11,
			-6, 7, 26, 31, 65, 56, 25, -20,
			-14, 13, 6, 21, 23, 12, 17, -23,
			-27, -2, -5, 12, 17, 6, 10, -25,
			-26, -4, -4, -10, 3, 3, 33, -12,
			-35, -1, -20, -23, -15, 24, 38, -22,
			0, 0, 0, 0, 0, 0, 0, 0,
		},

		// knight
		{-167, -89, -34, -49, 61, -97, -15, -107,
			-73, -41, 72, 36, 23, 62, 7, -17,
			-47, 60, 37, 65, 84, 129, 73, 44,
			-9, 17, 19, 53, 37, 69, 18, 22,
			-13, 4, 16, 13, 28, 19, 21, -8,
			-23, -9, 12, 10, 19, 17, 25, -16,
			-29, -53, -12, -3, -1, 18, -14, -19,
			-105, -21, -58, -33, -17, -28, -19, -23,
		},

		// bishop
		{-29, 4, -82, -37, -25, -42, 7, -8,
			-26, 16, -18, -13, 30, 59, 18, -47,
			-16, 37, 43, 40, 35, 50, 37, -2,
			-4, 5, 19, 50, 37, 37, 7, -2,
			-6, 13, 13, 26, 34, 12, 10, 4,
			0, 15, 15, 15, 14, 27, 18, 10,
			4, 15, 16, 0, 7, 21, 33, 1,
			-33, -3, -14, -21, -13, -12, -39, -21,
		},

		// rook
		{32, 42, 32, 51, 63, 9, 31, 43,
			27, 32, 58, 62, 80, 67, 26, 44,
			-5, 19, 26, 36, 17, 45, 61, 16,
			-24, -11, 7, 26, 24, 35, -8, -20,
			-36, -26, -12, -1, 9, -7, 6, -23,
			-45, -25, -16, -17, 3, 0, -5, -33,
			-44, -16, -20, -9, -1, 11, -6, -71,
			-19, -13, 1, 17, 16, 7, -37, -26,
		},

		// queen
		{
			-28, 0, 29, 12, 59, 44, 43, 45,
			-24, -39, -5, 1, -16, 57, 28, 54,
			-13, -17, 7, 8, 29, 56, 47, 57,
			-27, -27, -16, -16, -1, 17, -2, 1,
			-9, -26, -9, -10, -2, -4, 3, -3,
			-14, 2, -11, -2, -5, 2, 14, 5,
			-35, -8, 11, 2, 8, 15, -3, 1,
			-1, -18, -9, 10, -15, -25, -31, -50,
		},

		// king
		{
			-65, 23, 16, -15, -56, -34, 2, 13,
			29, -1, -20, -7, -8, -4, -38, -29,
			-9, 24, 2, -16, -20, 6, 22, -22,
			-17, -20, -12, -27, -30, -25, -14, -36,
			-49, -1, -27, -39, -46, -44, -33, -51,
			-14, -14, -22, -46, -44, -30, -15, -27,
			1, 7, -8, -64, -43, -16, 9, 8,
			-15, 36, 12, -54, 8, -28, 24, 14,
		},
	},
	// end game positional scores
	{
		// pawn
		{0, 0, 0, 0, 0, 0, 0, 0,
			178, 173, 158, 134, 147, 132, 165, 187,
			94, 100, 85, 67, 56, 53, 82, 84,
			32, 24, 13, 5, -2, 4, 17, 17,
			13, 9, -3, -7, -7, -8, 3, -1,
			4, 7, -6, 1, 0, -5, -1, -8,
			13, 8, 8, 10, 13, 0, 2, -7,
			0, 0, 0, 0, 0, 0, 0, 0,
		},

		// knight
		{-58, -38, -13, -28, -31, -27, -63, -99,
			-25, -8, -25, -2, -9, -25, -24, -52,
			-24, -20, 10, 9, -1, -9, -19, -41,
			-17, 3, 22, 22, 22, 11, 8, -18,
			-18, -6, 16, 25, 16, 17, 4, -18,
			-23, -3, -1, 15, 10, -3, -20, -22,
			-42, -20, -10, -5, -2, -20, -23, -44,
			-29, -51, -23, -15, -22, -18, -50, -64,
		},

		// bishop
		{-14, -21, -11, -8, -7, -9, -17, -24,
			-8, -4, 7, -12, -3, -13, -4, -14,
			2, -8, 0, -1, -2, 6, 0, 4,
			-3, 9, 12, 9, 14, 10, 3, 2,
			-6, 3, 13, 19, 7, 10, -3, -9,
			-12, -3, 8, 10, 13, 3, -7, -15,
			-14, -18, -7, -1, 4, -9, -15, -27,
			-23, -9, -23, -5, -9, -16, -5, -17,
		},

		// rook
		{13, 10, 18, 15, 12, 12, 8, 5,
			11, 13, 13, 11, -3, 3, 8, 3,
			7, 7, 7, 5, 4, -3, -5, -3,
			4, 3, 13, 1, 2, 1, -1, 2,
			3, 5, 8, 4, -5, -6, -8, -11,
			-4, 0, -5, -1, -7, -12, -8, -16,
			-6, -6, 0, 2, -9, -9, -11, -3,
			-9, 2, 3, -1, -5, -13, 4, -20,
		},

		// queen
		{-9, 22, 22, 27, 27, 19, 10, 20,
			-17, 20, 32, 41, 58, 25, 30, 0,
			-20, 6, 9, 49, 47, 35, 19, 9,
			3, 22, 24, 45, 57, 40, 57, 36,
			-18, 28, 19, 47, 31, 34, 39, 23,
			-16, -27, 15, 6, 9, 17, 10, 5,
			-22, -23, -30, -16, -16, -23, -36, -32,
			-33, -28, -22, -43, -5, -32, -20, -41,
		},

		// king
		{-74, -35, -18, -18, -11, 15, 4, -17,
			-12, 17, 14, 17, 17, 38, 23, 11,
			10, 17, 23, 15, 20, 45, 44, 13,
			-8, 22, 24, 27, 26, 33, 26, 3,
			-18, -4, 21, 24, 27, 23, 9, -11,
			-19, -3, 11, 21, 23, 16, 7, -9,
			-27, -11, 4, 13, 14, 4, -5, -17,
			-53, -34, -21, -11, -28, -14, -24, -43,
		},
	},
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
var DoublePawnPenalty int = -10

// penalty for making pawns isolated
var IsolatedPawnPenalty int = -10

// passed pawns are really good
var PassedPawnBonus = [8]int{
	0, 10, 30, 50, 75, 100, 150, 200,
}

// king shield
var KingShieldBonus int = 5

// open files are good
var SemiOpenFileScore = 10

// open files are good
var OpenFileScore = 15

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
	// grab our current game phase score
	gamePhaseScore := getGamePhaseScore()
	gamePhase := -1

	// calculate game phase
	if gamePhaseScore > openingPhaseScore {
		gamePhase = opening
	} else if gamePhaseScore < endGamePhaseScore {
		gamePhase = endgame
	} else {
		gamePhase = middlegame
	}

	// score of the position
	score := 0

	// looping variables
	var current_piece Bitboard
	var square int

	// loop
	for i := WhitePawn; i <= BlackKing; i++ {
		// copy
		current_piece = GameBoards[i]
		// loop over copied piece
		for current_piece != 0 {
			piece := int(i)
			// pop from gameboard copy
			square = current_piece.LSBIndex()
			current_piece.PopBit(square)
			// add piece score based on game phase
			if gamePhase == middlegame {
				score += (pieceScore[opening][piece]*
					gamePhaseScore +
					pieceScore[endgame][piece]*
						(openingPhaseScore-gamePhaseScore)) /
					openingPhaseScore
			} else { // end game or opening here
				score += pieceScore[gamePhase][piece]
			}
			// add position score depending on gamePhase / piece / square
			switch i {
			case WhitePawn:
				if gamePhase == middlegame {
					score += (postionalScores[opening][evalPawn][square]*
						gamePhaseScore +
						postionalScores[endgame][evalPawn][square]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score += postionalScores[gamePhase][evalPawn][square]
				}
				// score += pawnScore[square]
				// get double pawn penalty
				doublePawns := GameBoards[i] & FileMasks[square]
				if doublePawns.CountBits() > 1 {
					score += doublePawns.CountBits() * DoublePawnPenalty
				}
				// get isolated pawn penalty
				if GameBoards[i]&IsolatedMasks[square] == 0 {
					score += IsolatedPawnPenalty
				}

				// // white passed pawn bonus
				if WhitePassedMasks[square]&GameBoards[BlackPawn] == 0 {
					score += PassedPawnBonus[GetRank[square]]
				}
			case WhiteKnight:
				if gamePhase == middlegame {
					score += (postionalScores[opening][evalKnight][square]*
						gamePhaseScore +
						postionalScores[endgame][evalKnight][square]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score += postionalScores[gamePhase][evalKnight][square]
				}
				// score += knightScore[square]
			case WhiteBishop:
				if gamePhase == middlegame {
					score += (postionalScores[opening][evalBishop][square]*
						gamePhaseScore +
						postionalScores[endgame][evalBishop][square]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score += postionalScores[gamePhase][evalBishop][square]
				}
				// score += bishopScore[square]

				// mobility
				score += GetBishopAttack(square, GameOccupancy[Both]).CountBits()

			case WhiteRook:
				if gamePhase == middlegame {
					score += (postionalScores[opening][evalRook][square]*
						gamePhaseScore +
						postionalScores[endgame][evalRook][square]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score += postionalScores[gamePhase][evalRook][square]
				}
				// score += rookScore[square]

				// // semi open files for rook
				if GameBoards[WhitePawn]&FileMasks[square] == 0 {
					score += SemiOpenFileScore
				}

				// // // open files for rook
				// // semi open files for rook
				if (GameBoards[WhitePawn]|GameBoards[BlackPawn])&FileMasks[square] == 0 {
					score += OpenFileScore
				}
			case WhiteQueen:
				if gamePhase == middlegame {
					score += (postionalScores[opening][evalQueen][square]*
						gamePhaseScore +
						postionalScores[endgame][evalQueen][square]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score += postionalScores[gamePhase][evalQueen][square]
				}
				// mobility
				score += GetQueenAttack(square, GameOccupancy[Both]).CountBits()

			case WhiteKing:
				if gamePhase == middlegame {
					score += (postionalScores[opening][evalKing][square]*
						gamePhaseScore +
						postionalScores[endgame][evalKing][square]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score += postionalScores[gamePhase][evalKing][square]
				}
				// score += kingScore[square]

				// semi open files for king are bad
				if GameBoards[WhitePawn]&FileMasks[square] == 0 {
					score -= SemiOpenFileScore
				}

				// // open files for king are bad
				if (GameBoards[WhitePawn]|GameBoards[BlackPawn])&FileMasks[square] == 0 {
					score -= OpenFileScore
				}

				// // king safety is important
				score += (KingAttacks[square] & GameOccupancy[White]).CountBits() * KingShieldBonus

			case BlackPawn:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				if gamePhase == middlegame {
					score -= (postionalScores[opening][evalPawn][index_sq]*
						gamePhaseScore +
						postionalScores[endgame][evalPawn][index_sq]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score -= postionalScores[gamePhase][evalPawn][index_sq]
				}
				// get double pawn penalty
				doublePawns := GameBoards[i] & FileMasks[square]
				if doublePawns.CountBits() > 1 {
					score -= doublePawns.CountBits() * DoublePawnPenalty
				}
				// get isolated pawn penalty
				if GameBoards[i]&IsolatedMasks[square] == 0 {
					score -= IsolatedPawnPenalty
				}
				// give passed pawn bonus
				// white passed pawn bonus
				if BlackPassedMasks[square]&GameBoards[WhitePawn] == 0 {
					score -= PassedPawnBonus[GetRank[square]]
				}
			case BlackKnight:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				if gamePhase == middlegame {
					score -= (postionalScores[opening][evalKnight][index_sq]*
						gamePhaseScore +
						postionalScores[endgame][evalKnight][index_sq]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score -= postionalScores[gamePhase][evalKnight][index_sq]
				}
			case BlackBishop:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				if gamePhase == middlegame {
					score -= (postionalScores[opening][evalBishop][index_sq]*
						gamePhaseScore +
						postionalScores[endgame][evalBishop][index_sq]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score -= postionalScores[gamePhase][evalBishop][index_sq]
				}

				// mobility
				score -= GetBishopAttack(square, GameOccupancy[Both]).CountBits()
			case BlackRook:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				if gamePhase == middlegame {
					score -= (postionalScores[opening][evalRook][index_sq]*
						gamePhaseScore +
						postionalScores[endgame][evalRook][index_sq]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score -= postionalScores[gamePhase][evalRook][index_sq]
				}

				// score -= rookScore[mirrorScore[index_sq]]

				// // semi open files for rook
				if GameBoards[BlackPawn]&FileMasks[square] == 0 {
					score -= SemiOpenFileScore
				}

				// open files for rook
				if (GameBoards[WhitePawn]|GameBoards[BlackPawn])&FileMasks[square] == 0 {
					score -= OpenFileScore
				}
			case BlackQueen:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				if gamePhase == middlegame {
					score -= (postionalScores[opening][evalQueen][index_sq]*
						gamePhaseScore +
						postionalScores[endgame][evalQueen][index_sq]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score -= postionalScores[gamePhase][evalQueen][index_sq]
				}
				// mobility
				score -= GetQueenAttack(square, GameOccupancy[Both]).CountBits()
			case BlackKing:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				if gamePhase == middlegame {
					score -= (postionalScores[opening][evalKing][index_sq]*
						gamePhaseScore +
						postionalScores[endgame][evalKing][index_sq]*
							(openingPhaseScore-gamePhaseScore)) /
						openingPhaseScore
				} else {
					score -= postionalScores[gamePhase][evalKing][index_sq]
				}

				// // semi open files for king are bad
				if GameBoards[BlackPawn]&FileMasks[square] == 0 {
					score += SemiOpenFileScore
				}

				// open files for king are bad
				if (GameBoards[WhitePawn]|GameBoards[BlackPawn])&FileMasks[square] == 0 {
					score += OpenFileScore
				}

				// king safety is important
				score -= (KingAttacks[square] & GameOccupancy[White]).CountBits() * KingShieldBonus
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

// get game phase score
func getGamePhaseScore() int {
	whitePieces := 0
	blackPieces := 0
	/*
		Calculated as follows
		4 * knight material score in the opening +
		4 * bishop material score in the opening +
		4 * rook material score in the opening +
		2 * queen material score in the opening
	*/

	for i := WhiteKnight; i <= WhiteQueen; i++ {
		whitePieces += GameBoards[i].CountBits() * pieceScore[opening][i]
	}

	for i := BlackKnight; i <= BlackQueen; i++ {
		blackPieces += GameBoards[i].CountBits() * -pieceScore[opening][i]
	}
	return whitePieces + blackPieces
}
