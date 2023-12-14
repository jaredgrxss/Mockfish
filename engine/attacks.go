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

// constants to store piece info
const White, Black = 0, 1
var PawnAttacks [2][64] Bitboard
var KnightAttacks [2][64] Bitboard


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

func genKnightAttacks(side uint64, sq int) Bitboard {
	return 0
}

func GeneratePieceAttacks() {
	for i := 0; i < 64; i++ {
		PawnAttacks[0][i] = genPawnAttacks(0, i)
		PawnAttacks[1][i] = genPawnAttacks(1, i)
	}
}