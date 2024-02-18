package engine
import (
	"fmt"
)

/****************************************
	HELPERS RELATED TO MOVE GENERATION
****************************************/

// storing moves
var Moves []Bitboard

// main function to generate all psuedo legal moves
func GenerateMoves() {
	// if we had moves stored, get rid of them
	Moves = nil
	// loop over every piece
	for i := WhitePawn; i <= BlackKing; i++ {
		// special care has to be given to colored pawns, kings, and castling
		if SideToMove == White {
			if i == WhitePawn { genPawnMoves(White) }
			if i == WhiteKing { genKingMoves(White) }
		} else {
			if i == BlackPawn { genPawnMoves(Black) }
			if i == BlackKing { genKingMoves(Black) }
		}

		// knight moves
		genKnightMoves()
		// bishop moves
		genBishopMoves()
		// rook moves
		genRookMoves()
		// queen moves
		genQueenMoves()
	}
}

// helper for pawn moves
func genPawnMoves(side int) {
	// to -> from, aka: the move
	var source, target int
	// copy for current bitboard and all legal attacks given the square
	var bitboard, attacks Bitboard
	if (side == White) {
		bitboard = GameBoards[WhitePawn]
		// loop over all bits in the bitboard
		for bitboard != 0 {
			// get the current position & remove it from the board
			source = bitboard.LSBIndex(); bitboard.PopBit(source)
			// move the piece
			target = source >> 8
			// needs to be on the board and not occupied
			if !(target < int(A8)) && GameOccupancy[Both].GetBit(target) == 0 {
				// pawn promotion case
				if source >= int(A7) && source <= int(H7) {
					fmt.Println()
				} else {
					// pawn push one square if on board and not occupied

					// pawn push two squares if on board and not occupied
					if (source >= int(A2) && source <= int(H2)) && GameOccupancy[Both].GetBit(target >> 8) == 0 {

					}
				}
			}
		}
	} else {
		bitboard = GameBoards[BlackPawn]
		// loop over all bits in the bitboard
		for bitboard != 0 {
			// get the current position & remove it from the board
			source = bitboard.LSBIndex(); bitboard.PopBit(source)
			// move the piece
			target = source << 8
			// needs to be on the board and not occupied
			if !(target < int(A8)) && GameOccupancy[Both].GetBit(target) == 0 {
				// pawn promotion case
				if source >= int(A7) && source <= int(H7) {
					fmt.Println()
				} else {
					// pawn push one square if on board and not occupied

					// pawn push two squares if on board and not occupied
					if (source >= int(A2) && source <= int(H2)) && GameOccupancy[Both].GetBit(target >> 8) == 0 {

					}
				}
			}
		}
	}
}

// helper for king moves
func genKingMoves(side int) {

}

// helper for knight moves
func genKnightMoves() {

}

// helper for bishop moves
func genBishopMoves() {

}
// helper for rook moves
func genRookMoves() {

}

// helper for queen moves
func genQueenMoves() {

}
