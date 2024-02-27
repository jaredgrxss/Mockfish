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
		1. source square 		0000 0000 0000 0000 0011 1111 (6 bits / max square 63) HEX = 0x3f
		2. target square 		0000 0000 0000 1111 1100 0000 (6 bits / max square 63) HEX = 0xfc0
		3. piece type 			0000 0000 1111 0000 0000 0000 (4 bits / max value 11) HEX = 0xf000
		4. promoted piece 		0000 1111 0000 0000 0000 0000 (4 bits / max value 11) HEX = 0xf0000
		5. capture flag 		0001 0000 0000 0000 0000 0000 (1 bit) HEX = 0x100000
		6. double push flag 	0010 0000 0000 0000 0000 0000 (1 bit) HEX = 0x200000
		7. enpassant capture 	0100 0000 0000 0000 0000 0000 (1 bit) HEX = 0x400000
		8. castling  flag 		1000 0000 0000 0000 0000 0000 (1 bit) HEX = 0x800000
*/

type Moves struct {
	move_list  [256]int // big enough to store max legal moves in any pos
	move_count int      // to keep track of where to insert next move
}

// singleton for our moves in a game
var MoveList Moves

// main function to generate all PSUEDO LEGAL moves of a given position
func GeneratePositionMoves() {
	// clear any moves from previous position
	MoveList.clearMoveList()
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
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteQueen), 0, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteRook), 0, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteKnight), 0, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteBishop), 0, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), 0, 0, 0, 0, 0))
					// pawn push two squares if on board and not occupied
					if (source >= int(A2) && source <= int(H2)) && GameOccupancy[Both].GetBit(target-8) == 0 {
						MoveList.addMove(encodeMove(source, target-8, int(WhitePawn), 0, 0, 1, 0, 0))
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
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteQueen), 1, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteRook), 1, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteBishop), 1, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteKnight), 1, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					MoveList.addMove(encodeMove(source, target, int(WhitePawn), 0, 1, 0, 0, 0))
				}
			}

			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[White][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					enpassant_attacks.PopBit(target_enpassant)
					MoveList.addMove(encodeMove(source, target_enpassant, int(WhitePawn), 0, 1, 0, 1, 0))

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
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackQueen), 0, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackRook), 0, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackKnight), 0, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackBishop), 0, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), 0, 0, 0, 0, 0))
					// pawn push two squares if on board and not occupied
					if (source >= int(A7) && source <= int(H7)) && GameOccupancy[Both].GetBit(target+8) == 0 {
						MoveList.addMove(encodeMove(source, target+8, int(BlackPawn), 0, 0, 1, 0, 0))
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
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackQueen), 1, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackRook), 1, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackBishop), 1, 0, 0, 0))
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), int(BlackKnight), 1, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					MoveList.addMove(encodeMove(source, target, int(BlackPawn), 0, 1, 0, 0, 0))
				}
			}
			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[Black][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					enpassant_attacks.PopBit(target_enpassant)
					MoveList.addMove(encodeMove(source, target_enpassant, int(BlackPawn), 0, 1, 0, 1, 0))
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
					MoveList.addMove(encodeMove(int(E1), int(G1), int(WhiteKing), 0, 0, 0, 0, 1))
				}
			}
		}

		if Castle&White_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B1)) == 0 && GameOccupancy[Both].GetBit(int(C1)) == 0 &&
				GameOccupancy[Both].GetBit(int(D1)) == 0 {
				// first, make sure that E1 and D1 are not under attack
				if !IsSquareAttacked(int(E1), Black) && !IsSquareAttacked(int(D1), Black) {
					MoveList.addMove(encodeMove(int(E1), int(C1), int(WhiteKing), 0, 0, 0, 0, 1))
				}
			}
		}
	} else {
		if Castle&Black_king_side != 0 {
			// make sure that the squares between king and the kings rook are empty
			if GameOccupancy[Both].GetBit(int(F8)) == 0 && GameOccupancy[Both].GetBit(int(G8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(F8), White) {
					MoveList.addMove(encodeMove(int(E8), int(G8), int(BlackKing), 0, 0, 0, 0, 1))
				}
			}
		}

		if Castle&Black_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B8)) == 0 && GameOccupancy[Both].GetBit(int(C8)) == 0 &&
				GameOccupancy[Both].GetBit(int(D8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(D8), White) {
					MoveList.addMove(encodeMove(int(E8), int(C8), int(BlackKing), 0, 0, 0, 0, 1))
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
				MoveList.addMove(encodeMove(int(source), int(target), int(piece), 0, 0, 0, 0, 0))
			} else if side == Black && GameOccupancy[White].GetBit(target) == 0 {
				MoveList.addMove(encodeMove(int(source), int(target), int(piece), 0, 0, 0, 0, 0))
			} else {
				MoveList.addMove(encodeMove(int(source), int(target), int(piece), 0, 1, 0, 0, 0))
			}
		}
	}
}

// function to encode all possible information about a potential move, used by search function
func encodeMove(source, target, piece, promoted, capture, double, enpassant, castling int) int {
	return source |
		(target << 6) |
		(piece << 12) |
		(promoted << 16) |
		(capture << 20) |
		(double << 21) |
		(enpassant << 22) |
		(castling << 23)
}

// decode a encoded move with the schema mentioned above, only 24 bits of a 32 bit int are used
func decodeMove(encodedMove int) (source, target, piece, promoted, capture, double, enpassant, castling int) {
	return (encodedMove & 0x3f),
		(encodedMove & 0xfc0) >> 6,
		(encodedMove & 0xf000) >> 12,
		(encodedMove & 0xf0000) >> 16,
		(encodedMove & 0x100000) >> 20,
		(encodedMove & 0x200000) >> 21,
		(encodedMove & 0x400000) >> 22,
		(encodedMove & 0x800000) >> 23
}

// function to add move after it has already been encoded
func (moves *Moves) addMove(move int) {
	moves.move_list[moves.move_count] = move
	moves.move_count++
}

// function to print move after decoding
func (moves Moves) printMove(move int) {
	source, target, piece, promo, capture, double, enpassant, castling := decodeMove(move)
	if promo == 0 { 
		fmt.Printf("move	  piece	     capture	 double	 enpassant	castling\n")
		fmt.Println("         ", IntSquareToString[source]+IntSquareToString[target],
			"  ", IntToPieceName[piece], "    ", capture, "       ", double,
			"        ", enpassant, "          ", castling)
	} else {
		fmt.Printf("move	  piece	     capture	 double	 enpassant	castling\n")
		fmt.Println("         ", IntSquareToString[source]+IntSquareToString[target]+string(PromotedPieces[promo]),
			" ", IntToPieceName[piece], "    ", capture, "       ", double,
			"        ", enpassant, "          ", castling)
	}
	fmt.Println()
}

// function to print entire move list information for a given position
func (moves Moves) PrintMoveList() {
	if moves.move_count == 0 {
		fmt.Println("No moves generated for the current position...")
		return
	}
	// loop over move list and print each
	for i := 0; i < moves.move_count; i++ {
		fmt.Print("Move #", i + 1, "   ")
		moves.printMove(moves.move_list[i])
	}
}

func (moves *Moves) clearMoveList() {
	moves.move_list = [256]int{}; moves.move_count = 0
}
