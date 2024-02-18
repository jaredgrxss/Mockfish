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
	// reset stored moves
	Moves = nil
	// loop over every piece
	for i := WhitePawn; i <= BlackKing; i++ {
		// special care has to be given to colored pawns, kings, and castling
		if SideToMove == White {
			if i == WhitePawn {
				fmt.Println("Gen white pawn moves...")
				genPawnMoves(White)
			}
			if i == WhiteKing {
				fmt.Println("Get white king moves...")
				genKingMoves(White)
			}
		} else {
			if i == BlackPawn {
				fmt.Println("Gen black pawn moves...")
				genPawnMoves(Black)
			}
			if i == BlackKing {
				fmt.Println("Gen black king moves...")
				genKingMoves(Black)
			}
		}
		if i == WhiteKnight || i == BlackKnight {
			// knight moves
			fmt.Println("Gen knight moves...")
			genKnightMoves()
		}
		if i == WhiteBishop || i == BlackBishop {
			// bishop moves
			fmt.Println("Gen bishop moves...")
			genBishopMoves()
		}
		if i == WhiteRook || i == BlackRook {
			// rook moves
			fmt.Println("Gen rook moves...")
			genRookMoves()
		}
		if i == WhiteQueen || i == BlackQueen {
			// queen moves
			fmt.Println("Gen queen moves...")
			genQueenMoves()
		}
	}
}

// helper for pawn moves
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
			// move the piece
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
					target_enpassant := enpassant_attacks.LSBIndex(); enpassant_attacks.PopBit(target_enpassant)
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
			// move the piece
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
					target_enpassant := enpassant_attacks.LSBIndex(); enpassant_attacks.PopBit(target_enpassant)
					fmt.Println("pawn enpassant capture", IntSquareToString[source], IntSquareToString[target_enpassant])
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
