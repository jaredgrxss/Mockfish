package engine


/*******************************

	ZOBRIST HASHING HELPERS

*******************************/

// piece keys [piece][square]
var PieceKeys [12][64] uint64 

// enpassant keys [square]
var EnpassantKeys [64] uint64 

// castling keys [castling state]
var CastlingKeys [16] uint64

// side key [only has if side is black]
var SideKey uint64 

// ZOBRIST KEY for global game state
var HashKey uint64

// init all random hash keys
func InitZobrist() {
	// re-init seed
	RandomState = 1010243240
	// loop pieces
	for piece := WhitePawn; piece <= BlackKing; piece++ {
		// loop squares 
		for sq := 0; sq < 64; sq++ {
			// init piece keys 
			PieceKeys[piece][sq] = Get_u64_rand()
		}
	}

	// init enpassant keys 
	for sq := 0; sq < 64; sq++ {
		EnpassantKeys[sq] = Get_u64_rand()
	}

	// init castling keys 
	for i := 0; i < 16; i++ {
		CastlingKeys[i] = Get_u64_rand()
	}

	// init side key
	SideKey = Get_u64_rand()
}

// generate game state hash key
func GenerateHashKey() uint64 {
	// final hash 
	var FinalKey uint64 = 0

	// individual piece temporary bitboard
	var pieceBoard Bitboard

	// loop over all pieces
	for piece := WhitePawn; piece <= BlackKing; piece++ {
		pieceBoard = GameBoards[piece]
		// loop piece bitboard
		for pieceBoard != 0 {
			// get square for piece and remove it from copy
			sq := pieceBoard.LSBIndex()
			pieceBoard.PopBit(sq)
			// add to hash
			FinalKey ^= PieceKeys[piece][sq]
		}
	}

	// check to see if enpassant needs to be hashed
	if Enpassant != 64 {
		FinalKey ^= EnpassantKeys[Enpassant]
	}

	// hash castling rights 
	FinalKey ^= CastlingKeys[Castle]

	// hash side to move ONLY if black
	if SideToMove == Black {
		FinalKey ^= SideKey
	}

	return FinalKey
}