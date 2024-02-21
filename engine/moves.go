package engine

import (
	"fmt"
)

/****************************************
	HELPERS RELATED TO MOVE GENERATION
****************************************/

// storing moves
var Moves []Bitboard

// main function to generate all PSUEDO LEGAL moves of a given position
// MAKE_MOVE function will handle legality
func GeneratePositionMoves() {
	// reset stored moves
	Moves = nil
	// loop over every piece
	for i := WhitePawn; i <= BlackKing; i++ {
		// generate based on side moving, and then piece
		if SideToMove == White {
			switch i {
			case WhitePawn:
				fmt.Println("--- Gen white pawn moves ---")
				genPawnMoves(White)
			case WhiteKing:
				fmt.Println("--- Gen white king moves ---")
				genKingMoves(White)
			case WhiteKnight:
				fmt.Println("--- Gen white knight moves ---")
				genKnightMoves(White)
			case WhiteBishop:
				fmt.Println("--- Gen white bishop moves ---")
				genBishopMoves(White)
			case WhiteRook:
				fmt.Println("--- Gen white rook moves ---")
				genRookMoves(White)
			case WhiteQueen:
				fmt.Println("--- Gen white queen moves ---")
				genQueenMoves(White)
			}
		} else {
			switch i {
			case BlackPawn:
				fmt.Println("--- Gen black pawn moves ---")
				genPawnMoves(Black)
			case BlackKing:
				fmt.Println("--- Gen black king moves ---")
				genKingMoves(Black)
			case BlackKnight:
				fmt.Println("--- Gen black knight moves ---")
				genKnightMoves(Black)
			case BlackBishop:
				fmt.Println("--- Gen black bishop moves ---")
				genBishopMoves(Black)
			case BlackRook:
				fmt.Println("--- Gen black rook moves ---")
				genRookMoves(Black)
			case BlackQueen:
				fmt.Println("--- Gen black queen moves ---")
				genQueenMoves(Black)
			}
		}
	}
}

// helper for pawn moves (REWRITE THIS TO MAKE IT SIMPLER)
func genPawnMoves(side int) {
	// to -> from, aka: the move
	var source, target int
	// copy for current bitboard and all legal attacks given the square
	var bitboard, attacks Bitboard
	if side == White {
		bitboard = GameBoards[WhitePawn]
		// loop over all bits in the white pawn game bitboard
		for bitboard != 0 {
			// get the current position & remove it from the board
			source = bitboard.LSBIndex()
			bitboard.PopBit(source)
			// move the piece (this is based on INDEX of bit, so it is a little different)
			target = source - 8
			// needs to be on the board and not occupied
			if !(target < int(A8)) && GameOccupancy[Both].GetBit(target) == 0 {
				// pawn promotion case
				if source >= int(A7) && source <= int(H7) {
					fmt.Println("promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("single", IntSquareToString[source], IntSquareToString[target])
					// pawn push two squares if on board and not occupied
					if (source >= int(A2) && source <= int(H2)) && GameOccupancy[Both].GetBit(target-8) == 0 {
						fmt.Println("double", IntSquareToString[source], IntSquareToString[target-8])
					}
				}
			}

			// init pawn attacks & find attacks on only black pieces
			attacks = PawnAttacks[White][source] & GameOccupancy[Black]
			// loop over attacks bitboard
			for attacks != 0 {
				target = attacks.LSBIndex()
				attacks.PopBit(target)
				if source >= int(A7) && source <= int(H7) {
					fmt.Println("pawn capture promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("pawn capture", IntSquareToString[source], IntSquareToString[target])
				}
			}

			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[White][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					enpassant_attacks.PopBit(target_enpassant)
					fmt.Println("pawn enpassant capture", IntSquareToString[source], IntSquareToString[target_enpassant])
				}
			}
		}
	} else {
		bitboard = GameBoards[BlackPawn]
		// loop over all bits in the game black pawn bitboard
		for bitboard != 0 {
			// get the current position & remove it from the board
			source = bitboard.LSBIndex()
			bitboard.PopBit(source)
			// move the piece (this is based on INDEX of piece, so it is a little different)
			target = source + 8
			// needs to be on the board and not occupied
			if !(target > int(H1)) && GameOccupancy[Both].GetBit(target) == 0 {
				// pawn promotion case
				if source >= int(A2) && source <= int(H2) {
					fmt.Println("promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("single", IntSquareToString[source], IntSquareToString[target])
					// pawn push two squares if on board and not occupied
					if (source >= int(A7) && source <= int(H7)) && GameOccupancy[Both].GetBit(target+8) == 0 {
						fmt.Println("double", IntSquareToString[source], IntSquareToString[target+8])
					}
				}
			}
			// init pawn attacks & find attacks on only white pieces
			attacks = PawnAttacks[Black][source] & GameOccupancy[White]
			// loop over attacks bitboard
			for attacks != 0 {
				target = attacks.LSBIndex()
				attacks.PopBit(target)
				if source >= int(A2) && source <= int(H2) {
					fmt.Println("pawn capture promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("pawn capture", IntSquareToString[source], IntSquareToString[target])
				}
			}
			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[Black][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					enpassant_attacks.PopBit(target_enpassant)
					fmt.Println("pawn enpassant capture", IntSquareToString[source], IntSquareToString[target_enpassant])
				}
			}
		}
	}
}

// helper for king moves
func genKingMoves(side int) {
	if side == White {
		if Castle & White_king_side != 0 {
			// make sure that the squares between king and the kings rook are empty
			if GameOccupancy[Both].GetBit(int(F1)) == 0 && GameOccupancy[Both].GetBit(int(G1)) == 0 {
				// first, make sure E1 and F1 are not under attack
				if !IsSquareAttacked(int(E1), Black) && !IsSquareAttacked(int(F1), Black) {
					fmt.Println("castling move king side... e1 g1")
				}
			}
		}

		if Castle & White_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B1)) == 0 && GameOccupancy[Both].GetBit(int(C1)) == 0 &&
				GameOccupancy[Both].GetBit(int(D1)) == 0 {
				// first, make sure that E1 and D1 are not under attack
				if !IsSquareAttacked(int(E1), Black) && !IsSquareAttacked(int(D1), Black) {
					fmt.Println("castling move queen side... e1 c1")
				}
			}
		}
	} else {
		if Castle & Black_king_side != 0 {
			// make sure that the squares between king and the kings rook are empty
			if GameOccupancy[Both].GetBit(int(F8)) == 0 && GameOccupancy[Both].GetBit(int(G8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(F8), White) {
					fmt.Println("castling move king side... e8 g8")
				}
			}
		}

		if Castle & Black_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B8)) == 0 && GameOccupancy[Both].GetBit(int(C8)) == 0 &&
				GameOccupancy[Both].GetBit(int(D8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(D8), White) {
					fmt.Println("castling move queen side... e8 c8")
				}
			}
		}
	}
}

/*
	MAY END UP REWRITING THESE INTO ONE FUNCTION 
	THAT CAN BE CUSTOMIZED FOR SLIDING PIECE
*/
// helper for knight moves
func genKnightMoves(side int) {
	// to -> from movement
	var source, target int
	// copy our current piece on the actual game board
	var bitboard Bitboard
	// set our copied bitboard
	if side == White { bitboard = GameBoards[WhiteKnight] } else { bitboard = GameBoards[BlackKnight] }

	for bitboard != 0 {
		source = bitboard.LSBIndex(); bitboard.PopBit(source)
		
	}
}

// helper for bishop moves 
func genBishopMoves(side int) {
	// to -> from movement
	var source, target int
	// copy our current piece on the actual game board
	var bitboard Bitboard
	// set our copied bitoard
	if side == White { bitboard = GameBoards[WhiteBishop] } else { bitboard = GameBoards[BlackBishop] }

	for bitboard != 0 {
		source = bitboard.LSBIndex(); bitboard.PopBit(source)
	}
}

// helper for rook moves
func genRookMoves(side int) {

}

// helper for queen moves
func genQueenMoves(side int) {

}
