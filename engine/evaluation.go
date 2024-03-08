package engine

/*********************************

	EVALUATION DRIVERS AND HELPERS

*********************************/

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

func Evaluate() int {
	// score of the position
	score := 0
	// looping variables
	var current_piece Bitboard
	var piece, square int
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
			case WhiteKnight:
				score += knightScore[square]
			case WhiteBishop:
				score += bishopScore[square]
			case WhiteRook:
				score += rookScore[square]
			case WhiteKing:
				score += kingScore[square]
			// for black pieces mirror indices with square conversion
			case BlackPawn:
				string_sq := IntSquareToString[square]
				index_sq := StringSquareToBit[string_sq]
				score -= pawnScore[mirrorScore[index_sq]]
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
