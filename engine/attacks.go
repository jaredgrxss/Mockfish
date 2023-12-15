package engine

/*******************************************************
	INFORMATION ABOUT PRE-CALCULATED PIECE ATTACKS
	
	Structure :
	2D Array -> [sideToMove][sqaure]
********************************************************/


/* 
	NOT A FILE 
		8 0 1 1 1 1 1 1 1 
		7 0 1 1 1 1 1 1 1 
		6 0 1 1 1 1 1 1 1 
		5 0 1 1 1 1 1 1 1 
		4 0 1 1 1 1 1 1 1 
		3 0 1 1 1 1 1 1 1 
		2 0 1 1 1 1 1 1 1 
		1 0 1 1 1 1 1 1 1 
		a b c d e f g h


	NOT H FILE 

		8 1 1 1 1 1 1 1 0 
		7 1 1 1 1 1 1 1 0 
		6 1 1 1 1 1 1 1 0 
		5 1 1 1 1 1 1 1 0 
		4 1 1 1 1 1 1 1 0 
		3 1 1 1 1 1 1 1 0 
		2 1 1 1 1 1 1 1 0 
		1 1 1 1 1 1 1 1 0 
		a b c d e f g h

*/

// constants to help with correctly generating piece moves
// in case of overflow
const NOT_A_FILE Bitboard = 18374403900871474942
const NOT_H_FILE Bitboard = 9187201950435737471
const NOT_HG_FILE Bitboard = 4557430888798830399
const NOT_AB_FILE Bitboard = 18229723555195321596

// constants to store piece info for leaping pieces
const White, Black = 0, 1
var PawnAttacks [2][64] Bitboard
var KnightAttacks [64] Bitboard
var KingAttacks[64] Bitboard

// bitshifts for pawns is = 7,9
func genPawnAttacks(side uint64, sq int) Bitboard {

	// resulting bitboard of attacks
	var attacks Bitboard = 0

	// make new bitboard and set the pos of pawn on it
	var bitboard Bitboard = 0
	bitboard.SetBit(sq)

	// white side (0) black side (1)
	if side == 0 {
		if bitboard >> 7 & NOT_A_FILE != 0 {
			attacks |= (bitboard >> 7)
		}
		if bitboard >> 9 & NOT_H_FILE != 0 {
			attacks |= (bitboard >> 9)
		}
	} else {
		if bitboard << 7 & NOT_H_FILE != 0 {
			attacks |= (bitboard << 7)
		}
		if (bitboard << 9 & NOT_A_FILE != 0) {
			attacks |= (bitboard << 9)
		}
	}
	return attacks
}

// bitshifts for knights = 6, 10, 15, 17
func genKnightAttacks(sq int) Bitboard {
	// returning attack board
	var attacks Bitboard = 0

	// set bit for knight
	var bitboard Bitboard = 0
	bitboard.SetBit(sq)
	
	// calculating attacks
	if (bitboard >> 17) & NOT_H_FILE != 0 {
		attacks |= (bitboard >> 17)
	}
	if (bitboard >> 15) & NOT_A_FILE != 0 {
		attacks |= (bitboard >> 15)
	}
	if (bitboard >> 10) & NOT_HG_FILE != 0 {
		attacks |= (bitboard >> 10)
	}
	if (bitboard >> 6) & NOT_AB_FILE != 0 {
		attacks |= (bitboard >> 6)
	}

	// oposite side bit shifts
	if (bitboard << 17) & NOT_A_FILE != 0 {
		attacks |= (bitboard << 17)
	}
	if (bitboard << 15) & NOT_H_FILE != 0 {
		attacks |= (bitboard << 15)
	}
	if (bitboard << 10) & NOT_AB_FILE != 0 {
		attacks |= (bitboard << 10)
	}
	if (bitboard << 6) & NOT_HG_FILE != 0 {
		attacks |= (bitboard << 6)
	}

	return attacks
}

// bitshifts for king is = 1, 7, 8, 9
func genKingAttacks(sq int) Bitboard {
	// return attack board
	var attacks Bitboard = 0
	
	// set king on temp board
	var bitboard Bitboard = 0
	bitboard.SetBit(sq)

	// calculating attacks
	// up / down don't need checks
	attacks |= (bitboard << 8)
	attacks |= (bitboard >> 8)

	if (bitboard >> 1) & NOT_H_FILE != 0 {
		attacks |= (bitboard >> 1)
	}
	if (bitboard >> 7) & NOT_A_FILE != 0 {
		attacks |= (bitboard >> 7)
	}
	
	if (bitboard >> 9) & NOT_H_FILE != 0 {
		attacks |= (bitboard >> 9)
	}

	// oposite side bit shifts
	if (bitboard << 1) & NOT_A_FILE != 0 {
		attacks |= (bitboard << 1)
	}
	if (bitboard << 7) & NOT_H_FILE != 0 {
		attacks |= (bitboard << 7)
	}

	if (bitboard << 9) & NOT_A_FILE != 0 {
		attacks |= (bitboard << 9)
	}

	return attacks
}

func GeneratePieceAttacks() {
	for i := 0; i < 64; i++ {
		PawnAttacks[0][i] = genPawnAttacks(0, i)
		PawnAttacks[1][i] = genPawnAttacks(1, i)
		KnightAttacks[i] = genKnightAttacks(i)
		KingAttacks[i] = genKingAttacks(i)
	}
}