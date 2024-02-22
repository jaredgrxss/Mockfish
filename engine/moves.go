package engine

import (
	"fmt"
)

/****************************************
	HELPERS RELATED TO MOVE GENERATION
****************************************/
/*
	----- ENCODED DATA -----
	Things needing to be encoded:
		1. source square 		0000 0000 0000 0000 0011 1111 (6 bits = max square = 63) HEX = 0x3f
		2. target square 		0000 0000 0000 1111 1100 0000 (6 bits = max square = 63) HEX = 0xfc0
		3. piece type 			0000 0000 1111 0000 0000 0000 (4 bits = max value = 11) HEX = 0xf000
		4. promoted piece 		0000 1111 0000 0000 0000 0000 (4 bits = max value = 11) HEX = 0xf0000
		5. capture flag 		0001 0000 0000 0000 0000 0000 (1 bit) HEX = 0x100000
		6. double push flag 	0010 0000 0000 0000 0000 0000 (1 bit) HEX = 0x200000
		7. enpassant capture 	0100 0000 0000 0000 0000 0000 (1 bit) HEX = 0x400000
		8. castling  flag 		1000 0000 0000 0000 0000 0000 (1 bit) HEX = 0x800000
*/

// storing moves
var Moves []Bitboard

// main function to generate all PSUEDO LEGAL moves of a given position
// MAKE_MOVE function will handle legality
func GeneratePositionMoves() {
	// reset moves
	Moves = nil
	// loop over every piece
	for i := WhitePawn; i <= BlackKing; i++ {
		// generate based on side moving, and then piece
		if SideToMove == White {
			switch i {
			case WhitePawn:
				genPawnMoves(White)
				fmt.Println()
			case WhiteKing:
				genSlidingPieceMoves(White, i)
				genKingCastleMoves(White)
				fmt.Println()
			case WhiteBishop, WhiteRook, WhiteQueen, WhiteKnight:
				genSlidingPieceMoves(White, i)
				fmt.Println()
			}
		} else {
			switch i {
			case BlackPawn:
				genPawnMoves(Black)
				fmt.Println()
			case BlackKing:
				genSlidingPieceMoves(Black, i)
				genKingCastleMoves(Black)
				fmt.Println()
			case BlackBishop, BlackRook, BlackQueen, BlackKnight:
				genSlidingPieceMoves(Black, i)
				fmt.Println()
			}
		}
	}
}

// helper for generating pawn moves
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
					fmt.Println("WhitePawn promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("WhitePawn single push", IntSquareToString[source], IntSquareToString[target])
					// pawn push two squares if on board and not occupied
					if (source >= int(A2) && source <= int(H2)) && GameOccupancy[Both].GetBit(target-8) == 0 {
						fmt.Println("WhitePawn double", IntSquareToString[source], IntSquareToString[target-8])
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
					fmt.Println("WhitePawn capture promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("WhitePawn capture", IntSquareToString[source], IntSquareToString[target])
				}
			}

			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[White][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					enpassant_attacks.PopBit(target_enpassant)
					fmt.Println("WhitePawn enpassant capture", IntSquareToString[source], IntSquareToString[target_enpassant])
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
					fmt.Println("BlackPawn promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("BlackPawn single push", IntSquareToString[source], IntSquareToString[target])
					// pawn push two squares if on board and not occupied
					if (source >= int(A7) && source <= int(H7)) && GameOccupancy[Both].GetBit(target+8) == 0 {
						fmt.Println("BlackPawn double push", IntSquareToString[source], IntSquareToString[target+8])
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
					fmt.Println("BlackPawn capture promotion", IntSquareToString[source], IntSquareToString[target])
				} else {
					// pawn push one square, already checked that is on board
					fmt.Println("BlackPawn capture", IntSquareToString[source], IntSquareToString[target])
				}
			}
			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[Black][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					enpassant_attacks.PopBit(target_enpassant)
					fmt.Println("BlackPawn enpassant capture", IntSquareToString[source], IntSquareToString[target_enpassant])
				}
			}
		}
	}
}

// helper for generating king moves
func genKingCastleMoves(side int) {
	if side == White {
		if Castle&White_king_side != 0 {
			// make sure that the squares between king and the kings rook are empty
			if GameOccupancy[Both].GetBit(int(F1)) == 0 && GameOccupancy[Both].GetBit(int(G1)) == 0 {
				// first, make sure E1 and F1 are not under attack
				if !IsSquareAttacked(int(E1), Black) && !IsSquareAttacked(int(F1), Black) {
					fmt.Println("WhiteKing castling move king side... e1 g1")
				}
			}
		}

		if Castle&White_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B1)) == 0 && GameOccupancy[Both].GetBit(int(C1)) == 0 &&
				GameOccupancy[Both].GetBit(int(D1)) == 0 {
				// first, make sure that E1 and D1 are not under attack
				if !IsSquareAttacked(int(E1), Black) && !IsSquareAttacked(int(D1), Black) {
					fmt.Println("WhiteKing castling move queen side... e1 c1")
				}
			}
		}
	} else {
		if Castle&Black_king_side != 0 {
			// make sure that the squares between king and the kings rook are empty
			if GameOccupancy[Both].GetBit(int(F8)) == 0 && GameOccupancy[Both].GetBit(int(G8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(F8), White) {
					fmt.Println("BlackKing castling move king side... e8 g8")
				}
			}
		}

		if Castle&Black_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B8)) == 0 && GameOccupancy[Both].GetBit(int(C8)) == 0 &&
				GameOccupancy[Both].GetBit(int(D8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(D8), White) {
					fmt.Println("BlackKing castling move queen side... e8 c8")
				}
			}
		}
	}
}

// helper for generating all sliding piece moves
func genSlidingPieceMoves(side int, piece Piece) {
	// current source squares for this piece
	bitboard := GameBoards[piece]
	// loop over all postions of this piece
	for bitboard != 0 {
		source := bitboard.LSBIndex()
		bitboard.PopBit(source)
		// board to store target squares
		var attacks Bitboard
		// look up which piece that's needed, make sure it doesn't attack pieces of same color
		switch piece {
		case WhiteBishop, BlackBishop:
			attacks = GetBishopAttack(source, GameOccupancy[Both]) & ^GameOccupancy[side]
		case WhiteRook, BlackRook:
			attacks = GetRookAttack(source, GameOccupancy[Both]) & ^GameOccupancy[side]
		case WhiteQueen, BlackQueen:
			attacks = GetQueenAttack(source, GameOccupancy[Both]) & ^GameOccupancy[side]
		case WhiteKnight, BlackKnight:
			attacks = KnightAttacks[source] & ^GameOccupancy[side]
		case BlackKing, WhiteKing:
			attacks = KingAttacks[source] & ^GameOccupancy[side]
		}
		// loop over all possible target squares
		for attacks != 0 {
			target := attacks.LSBIndex()
			attacks.PopBit(target)
			// quiet moves OR capture moves
			if side == White && GameOccupancy[Black].GetBit(target) == 0 {
				fmt.Println(IntToPieceName[piece], "quiet move", IntSquareToString[source], IntSquareToString[target])
			} else if side == Black && GameOccupancy[White].GetBit(target) == 0 {
				fmt.Println(IntToPieceName[piece], "quiet move", IntSquareToString[source], IntSquareToString[target])
			} else {
				fmt.Println(IntToPieceName[piece], "piece capture", IntSquareToString[source], IntSquareToString[target])
			}
		}
	}
}
