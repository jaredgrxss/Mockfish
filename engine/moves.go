package engine

import (
	"bufio"
	"fmt"
	"os"
)

/****************************************
	HELPERS RELATED TO MOVE GENERATION
****************************************/

type Moves struct {
	Move_list  [256]int // big enough to store max legal moves in any pos
	Move_count int      // to keep track of where to insert next move
}

const (
	allMoves = iota
	onlyCaptures
)

// main function to generate all PSUEDO LEGAL moves of a given position
func GeneratePositionMoves(moves *Moves) {
	// loop over every piece
	for i := WhitePawn; i <= BlackKing; i++ {
		// generate based on side moving, and then piece
		if SideToMove == White {
			switch i {
			case WhitePawn:
				genPawnMoves(White, moves)
			case WhiteKing:
				genSlidingPieceMoves(White, i, moves)
				genKingCastleMoves(White, moves)
			case WhiteBishop, WhiteRook, WhiteQueen, WhiteKnight:
				genSlidingPieceMoves(White, i, moves)
			}
		} else {
			switch i {
			case BlackPawn:
				genPawnMoves(Black, moves)
			case BlackKing:
				genSlidingPieceMoves(Black, i, moves)
				genKingCastleMoves(Black, moves)
			case BlackBishop, BlackRook, BlackQueen, BlackKnight:
				genSlidingPieceMoves(Black, i, moves)
			}
		}
	}
}

// helper for generating pawn moves
func genPawnMoves(side int, moves *Moves) {
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
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteQueen), 0, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteRook), 0, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteKnight), 0, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteBishop), 0, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					moves.addMove(encodeMove(source, target, int(WhitePawn), 0, 0, 0, 0, 0))
					// pawn push two squares if on board and not occupied
					if (source >= int(A2) && source <= int(H2)) && GameOccupancy[Both].GetBit(target-8) == 0 {
						moves.addMove(encodeMove(source, target-8, int(WhitePawn), 0, 0, 1, 0, 0))
					}
				}
			}

			// init pawn attacks & find attacks on only black pieces
			attacks = PawnAttacks[side][source] & GameOccupancy[Black]
			// loop over attacks bitboard
			for attacks != 0 {
				target = attacks.LSBIndex()
				attacks.PopBit(target)
				if source >= int(A7) && source <= int(H7) {
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteQueen), 1, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteRook), 1, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteBishop), 1, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(WhitePawn), int(WhiteKnight), 1, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					moves.addMove(encodeMove(source, target, int(WhitePawn), 0, 1, 0, 0, 0))
				}
			}

			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[side][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					moves.addMove(encodeMove(source, target_enpassant, int(WhitePawn), 0, 1, 0, 1, 0))

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
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackQueen), 0, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackRook), 0, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackKnight), 0, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackBishop), 0, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					moves.addMove(encodeMove(source, target, int(BlackPawn), 0, 0, 0, 0, 0))
					// pawn push two squares if on board and not occupied
					if (source >= int(A7) && source <= int(H7)) && GameOccupancy[Both].GetBit(target+8) == 0 {
						moves.addMove(encodeMove(source, target+8, int(BlackPawn), 0, 0, 1, 0, 0))
					}
				}
			}
			// init pawn attacks & find attacks on only white pieces
			attacks = PawnAttacks[side][source] & GameOccupancy[White]
			// loop over attacks bitboard
			for attacks != 0 {
				target = attacks.LSBIndex()
				attacks.PopBit(target)
				if source >= int(A2) && source <= int(H2) {
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackQueen), 1, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackRook), 1, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackBishop), 1, 0, 0, 0))
					moves.addMove(encodeMove(source, target, int(BlackPawn), int(BlackKnight), 1, 0, 0, 0))
				} else {
					// pawn push one square, already checked that is on board
					moves.addMove(encodeMove(source, target, int(BlackPawn), 0, 1, 0, 0, 0))
				}
			}
			// enpssant captures
			if Enpassant != 64 {
				// see if this square is under attack by current pawn placement
				enpassant_attacks := PawnAttacks[side][source] & (1 << Enpassant)
				if enpassant_attacks != 0 {
					target_enpassant := enpassant_attacks.LSBIndex()
					moves.addMove(encodeMove(source, target_enpassant, int(BlackPawn), 0, 1, 0, 1, 0))
				}
			}
		}
	}
}

// helper for generating king castling moves
func genKingCastleMoves(side int, moves *Moves) {
	if side == White {
		if Castle&White_king_side != 0 {
			// make sure that the squares between king and the kings rook are empty
			if GameOccupancy[Both].GetBit(int(F1)) == 0 && GameOccupancy[Both].GetBit(int(G1)) == 0 {
				// first, make sure E1 and F1 are not under attack
				if !IsSquareAttacked(int(E1), Black) && !IsSquareAttacked(int(F1), Black) {
					moves.addMove(encodeMove(int(E1), int(G1), int(WhiteKing), 0, 0, 0, 0, 1))
				}
			}
		}

		if Castle&White_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B1)) == 0 && GameOccupancy[Both].GetBit(int(C1)) == 0 &&
				GameOccupancy[Both].GetBit(int(D1)) == 0 {
				// first, make sure that E1 and D1 are not under attack
				if !IsSquareAttacked(int(E1), Black) && !IsSquareAttacked(int(D1), Black) {
					moves.addMove(encodeMove(int(E1), int(C1), int(WhiteKing), 0, 0, 0, 0, 1))
				}
			}
		}
	} else {
		if Castle&Black_king_side != 0 {
			// make sure that the squares between king and the kings rook are empty
			if GameOccupancy[Both].GetBit(int(F8)) == 0 && GameOccupancy[Both].GetBit(int(G8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(F8), White) {
					moves.addMove(encodeMove(int(E8), int(G8), int(BlackKing), 0, 0, 0, 0, 1))
				}
			}
		}

		if Castle&Black_queen_side != 0 {
			// make sure that the squares between king and queens rook are empty
			if GameOccupancy[Both].GetBit(int(B8)) == 0 && GameOccupancy[Both].GetBit(int(C8)) == 0 &&
				GameOccupancy[Both].GetBit(int(D8)) == 0 {
				// first, make sure E8 and F8 are not under attack
				if !IsSquareAttacked(int(E8), White) && !IsSquareAttacked(int(D8), White) {
					moves.addMove(encodeMove(int(E8), int(C8), int(BlackKing), 0, 0, 0, 0, 1))
				}
			}
		}
	}
}

// helper for generating all sliding piece moves
func genSlidingPieceMoves(side int, piece Piece, moves *Moves) {
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
				moves.addMove(encodeMove(int(source), int(target), int(piece), 0, 0, 0, 0, 0))
			} else if side == Black && GameOccupancy[White].GetBit(target) == 0 {
				moves.addMove(encodeMove(int(source), int(target), int(piece), 0, 0, 0, 0, 0))
			} else {
				moves.addMove(encodeMove(int(source), int(target), int(piece), 0, 1, 0, 0, 0))
			}
		}
	}
}

/***************************************************************************************************************
									----- ENCODED DATA -----
		 Things encoded: 							Offset / Schema:
		1. source square 		0000 0000 0000 0000 0011 1111 (6 bits / max square 63) HEX_OFFSET = 0x3f
		2. target square 		0000 0000 0000 1111 1100 0000 (6 bits / max square 63) HEX_OFFSET = 0xfc0
		3. piece type 			0000 0000 1111 0000 0000 0000 (4 bits / max value 11) HEX_OFFSET = 0xf000
		4. promoted piece 		0000 1111 0000 0000 0000 0000 (4 bits / max value 11) HEX_OFFSET = 0xf0000
		5. capture flag 		0001 0000 0000 0000 0000 0000 (1 bit) HEX_OFFSET = 0x100000
		6. double push flag 	0010 0000 0000 0000 0000 0000 (1 bit) HEX_OFFSET = 0x200000
		7. enpassant capture 	0100 0000 0000 0000 0000 0000 (1 bit) HEX_OFFSET = 0x400000
		8. castling  flag 		1000 0000 0000 0000 0000 0000 (1 bit) HEX_OFFSET = 0x800000
***************************************************************************************************************/

// encode all possible information about a potential move with offsets, used by search function
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
func DecodeMove(encodedMove int) (source, target, piece, promoted, capture, double, enpassant, castling int) {
	return (encodedMove & 0x3f),
		(encodedMove & 0xfc0) >> 6,
		(encodedMove & 0xf000) >> 12,
		(encodedMove & 0xf0000) >> 16,
		(encodedMove & 0x100000) >> 20,
		(encodedMove & 0x200000) >> 21,
		(encodedMove & 0x400000) >> 22,
		(encodedMove & 0x800000) >> 23
}

// add a move to the moves after it has already been encoded
func (moves *Moves) addMove(move int) {
	moves.Move_list[moves.Move_count] = move
	moves.Move_count++
}

// print a single move after decoding
func (moves Moves) PrintMove(move int) {
	source, target, piece, promo, capture, double, enpassant, castling := DecodeMove(move)
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

// print entire move list information for a given position
func (moves Moves) Printmoves() {
	if moves.Move_count == 0 {
		fmt.Println("No moves generated for the current position...")
		return
	}
	// loop over move list and print each
	for i := 0; i < moves.Move_count; i++ {
		fmt.Print("Move #", i+1, "   ")
		moves.PrintMove(moves.Move_list[i])
	}
}

// main make move function
func MakeMove(move int, move_flag int) int {
	// parse the move information
	source, target, piece, promo, capture, double, enpassant, castling := DecodeMove(move)
	// distinguish between quiet / capture moves
	if move_flag == allMoves {
		// preserve the board state
		GameBoardsCopy := GameBoards
		GameOccupancyCopy := GameOccupancy
		SideToMoveCopy := SideToMove
		EnpassantCopy := Enpassant
		CastleCopy := Castle
		HashKeyCopy := HashKey

		// perform the move
		GameBoards[piece].PopBit(source)
		GameBoards[piece].SetBit(target)

		// hash piece info
		HashKey ^= PieceKeys[piece][source] // remove source piece info
		HashKey ^= PieceKeys[piece][target] // addd target piece info

		// handle if move was a capture
		if capture == 1 {
			handleCaptureMove(target)
		}

		// handle if move was a promotion
		if promo != 0 {
			handlePawnPromotions(target, promo)
		}

		// handle if move was an enpassant move
		if enpassant == 1 {
			handleEnpassantMove(target)
		}

		// hash enpassant (remove enpassant)
		if Enpassant != 64 {
			HashKey ^= EnpassantKeys[Enpassant]
		}

		// reset enpassant if it is not chosen as a move
		Enpassant = 64

		// handle if move was a double pawn push
		if double == 1 {
			handleDoublePawnPush(target)
		}

		// handle if move was a castling move
		if castling == 1 {
			handleCastlingMove(target)
		}

		// update castling rights every move
		updateCastlingRights(source, target)

		// update occupancy boards with every move
		updateOccupancyBoard()

		// change side
		SideToMove ^= 1

		// hash the side
		HashKey ^= SideKey

		/***********************
			ZOBRIST HASHING DEBUGGING
		***********************/
		// build hash key after move made
		// NewHash := GenerateHashKey()
		// if NewHash != HashKey {
		// 	fmt.Print("move: ")
		// 	PrintUCICompatibleMove(move)
		// 	fmt.Printf(" hash key should be: %d but is %d\n", NewHash, HashKey)
		// }

		// check to see if check was put in check from move
		if SideToMove == White {
			if IsSquareAttacked(GameBoards[BlackKing].LSBIndex(), SideToMove) {
				GameBoards = GameBoardsCopy
				GameOccupancy = GameOccupancyCopy
				SideToMove = SideToMoveCopy
				Enpassant = EnpassantCopy
				Castle = CastleCopy
				HashKey = HashKeyCopy
				return 0
			}
		} else {
			if IsSquareAttacked(GameBoards[WhiteKing].LSBIndex(), SideToMove) {
				GameBoards = GameBoardsCopy
				GameOccupancy = GameOccupancyCopy
				SideToMove = SideToMoveCopy
				Enpassant = EnpassantCopy
				Castle = CastleCopy
				HashKey = HashKeyCopy
				return 0
			}
		}
		// this was a legal move
		return 1
	} else if capture == 1 {
		MakeMove(move, allMoves)
	}
	// bad move / bad input
	return 0
}

// captures made on the board
func handleCaptureMove(target int) {
	var startPiece, endPiece Piece
	// have to loop over opposite sides boards for captures
	if SideToMove == White {
		startPiece = BlackPawn
		endPiece = BlackKing
	} else {
		startPiece = WhitePawn
		endPiece = WhiteKing
	}
	// find which bitboard has the piece being capture
	for i := startPiece; i <= endPiece; i++ {
		if GameBoards[i].GetBit(target) == 1 {
			GameBoards[i].PopBit(target)

			// remove piece from hashkey
			HashKey ^= PieceKeys[i][target]
			break
		}
	}
}

// promotions made on the board
func handlePawnPromotions(target int, promo int) {
	// erase pawn from target square
	if SideToMove == White {
		GameBoards[WhitePawn].PopBit(target)

		// remove pawn hash key
		HashKey ^= PieceKeys[WhitePawn][target]
	} else {
		GameBoards[BlackPawn].PopBit(target)
		// remove pawn hash key
		HashKey ^= PieceKeys[BlackPawn][target]
	}
	// set promoted piece
	GameBoards[promo].SetBit(target)
	// add promoted
	HashKey ^= PieceKeys[promo][target]
}

// enpassent captures made on the board
func handleEnpassantMove(target int) {
	// need to remove captured piece from the board
	if SideToMove == White {
		GameBoards[BlackPawn].PopBit(target + 8)
	} else {
		GameBoards[WhitePawn].PopBit(target - 8)
	}

	// hashing info
	if SideToMove == White {
		GameBoards[BlackPawn].PopBit(target + 8)

		// remove hash
		HashKey ^= PieceKeys[BlackPawn][target+8]
	} else {
		GameBoards[WhitePawn].PopBit(target - 8)
		// remove hash
		HashKey ^= PieceKeys[WhitePawn][target-8]
	}
}

// a double pawn push was made on the board
func handleDoublePawnPush(target int) {
	if SideToMove == White {
		Enpassant = target + 8
		// hash enpassant
		HashKey ^= EnpassantKeys[target+8]
	} else {
		Enpassant = target - 8
		// has nepassant
		HashKey ^= EnpassantKeys[target-8]
	}
}

// a castling move was made on the board, need to update rook
func handleCastlingMove(target int) {
	switch target {
	// white castle king side, move H rook
	case int(G1):
		GameBoards[WhiteRook].PopBit(int(H1))
		GameBoards[WhiteRook].SetBit(int(F1))
		// hash rook
		HashKey ^= PieceKeys[WhiteRook][int(H1)]
		HashKey ^= PieceKeys[WhiteRook][int(F1)]

	// white castle queen side, move A rook
	case int(C1):
		GameBoards[WhiteRook].PopBit(int(A1))
		GameBoards[WhiteRook].SetBit(int(D1))
		// hash rook
		HashKey ^= PieceKeys[WhiteRook][int(A1)]
		HashKey ^= PieceKeys[WhiteRook][int(D1)]

	// black castle king side, move H rook
	case int(G8):
		GameBoards[BlackRook].PopBit(int(H8))
		GameBoards[BlackRook].SetBit(int(F8))
		// hash rook
		HashKey ^= PieceKeys[BlackRook][int(H8)]
		HashKey ^= PieceKeys[BlackRook][int(F8)]

	// black castle queen side, move A rook
	case int(C8):
		GameBoards[BlackRook].PopBit(int(A8))
		GameBoards[BlackRook].SetBit(int(D8))
		// hash rook
		HashKey ^= PieceKeys[BlackRook][int(A8)]
		HashKey ^= PieceKeys[BlackRook][int(D8)]
	}
}

// update castlign rights of new board position
func updateCastlingRights(source int, target int) {
	// unset before castling
	HashKey ^= CastlingKeys[Castle]

	Castle &= CastlingRightsHelper[source]
	Castle &= CastlingRightsHelper[target]

	// set after castling
	HashKey ^= CastlingKeys[Castle]

}

// update occupancies to reflect new board position
func updateOccupancyBoard() {
	// reset the boards
	for i := 0; i < 3; i++ {
		GameOccupancy[i] = 0
	}

	// update white
	for i := WhitePawn; i <= WhiteKing; i++ {
		GameOccupancy[White] |= GameBoards[i]
	}

	// update black
	for i := BlackPawn; i <= BlackKing; i++ {
		GameOccupancy[Black] |= GameBoards[i]
	}

	// update both
	GameOccupancy[Both] |= GameOccupancy[White]
	GameOccupancy[Both] |= GameOccupancy[Black]
}

func PrintMoveScores(moves Moves) {
	for i := 0; i < moves.Move_count; i++ {
		PrintUCICompatibleMove(moves.Move_list[i])
		fmt.Printf(" score: %d\n", ScoreMove(moves.Move_list[i]))
	}
}

// quick test function for debugging
func TestMakeMove() {
	var moves Moves
	moves.Move_count = 0
	GeneratePositionMoves(&moves)
	for i := 0; i < moves.Move_count; i++ {
		move := moves.Move_list[i]
		// copy position
		var GameBoards_Copy [12]Bitboard
		var GameOccupancy_Copy [3]Bitboard
		var SideToMove_Copy, Enpassant_Copy, Castle_Copy int
		GameBoards_Copy = GameBoards
		GameOccupancy_Copy = GameOccupancy
		SideToMove_Copy = SideToMove
		Enpassant_Copy = Enpassant
		Castle_Copy = Castle

		if MakeMove(move, allMoves) == 0 {
			continue
		}
		PrintGameboard()
		fmt.Print("Move #", i+1, "   ")
		moves.PrintMove(move)
		// restore position
		GameBoards = GameBoards_Copy
		GameOccupancy = GameOccupancy_Copy
		SideToMove = SideToMove_Copy
		Enpassant = Enpassant_Copy
		Castle = Castle_Copy
		PrintGameboard()
		buf := bufio.NewReader(os.Stdin)
		buf.ReadBytes('\n')
	}
}
