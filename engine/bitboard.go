package engine

import (
	"fmt"
	"strconv"
)

// Types
type Bitboard uint64
type Piece int
type Square int

// FEN strings
var EMPTY_BOARD = "8/8/8/8/8/8/8/8 w - - "
var START_POSITION = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 "
var TRICKY_POSITION = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1 "
var KILLER_POSITION = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
var CMK_POSITION = "r2q1rk1/ppp2ppp/2n1bn2/2b1p3/3pP3/3P1NPP/PPP1NPB1/R1BQ1RK1 b - - 0 9 "

/****************************************
	BOARD REPRESENTATION HELPERS
****************************************/
/*
	castling representation info... 

	bin		dec		description
	
	0001 	1 		white king can castle king side
	0010	2		white king can castle queen side
	0100	4		black king can castle king side
	1000	8		black king can caslte queen side
*/
var GameBoards [12]Bitboard // 12 boards for each piece + color (6 * 2)
var GameOccupancy [3]Bitboard // white = 0, black = 1, both = 2
var SideToMove = -1
var Enpassant int = 64
var Castle int
const White_king_side, White_queen_side, Black_king_side, Black_queen_side = 1, 2, 4, 8

/******************************
		ENUMERATIONS
******************************/
const ( White = iota; Black; Both)

const ( WhitePawn Piece = iota; WhiteKnight; WhiteBishop; WhiteRook; WhiteQueen; WhiteKing 
		BlackPawn; BlackKnight; BlackBishop; BlackRook; BlackQueen; BlackKing )

// can use either ASCII or Unicode for board debugging
var AsciiPieces string = "PNBRQKpnbrqk" // must be type casted with string()

var UnicodePieces [12]string = [12]string {
	"♙", "♘", "♗", "♖", "♕", "♔", 
	"♟︎", "♞", "♝", "♜", "♛", "♚",
}

var AsciiToType = map[string]Piece {
	"P": WhitePawn, "N": WhiteKnight, "B": WhiteBishop, "R": WhiteRook, "Q": WhiteQueen, "K": WhiteKing,
	"p": BlackPawn, "n": BlackKnight, "b": BlackBishop, "r": BlackRook, "q": BlackQueen, "k": BlackKing,
}

var StringSquareToBit = map[string]Square {
	"a8": A8, "b8": B8, "c8": C8, "d8": D8, "e8": E8, "f8": F8, "g8": G8, "h8": H8,
	"a7": A7, "b7": B7, "c7": C7, "d7": D7, "e7": E7, "f7": F7, "g7": G7, "h7": H7, 
	"a6": A6, "b6": B6, "c6": C6, "d6": D6, "e6": E6, "f6": F6, "g6": G6, "h6": H6,
	"a5": A5, "b5": B5, "c5": C5, "d5": D5, "e5": E5, "f5": F5, "g5": G5, "h5": H5, 
	"a4": A4, "b4": B4, "c4": C4, "d4": D4, "e4": E4, "f4": F4, "g4": G4, "h4": H4, 
	"a3": A3, "b3": B3, "c3": C3, "d3": D3, "e3": E3, "f3": F3, "g3": G3, "h3": H3, 
	"a2": A2, "b2": B2, "c2": C2, "d2": D2, "e2": E2, "f2": F2, "g2": G2, "h2": H2, 
	"a1": A1, "b1": B1, "c1": C1, "d1": D1, "e1": E1, "f1": F1, "g1": G1, "h1": H1, 
}

var IntSquareToString = map[int]string {
	0: "a8", 1: "b8", 2: "c8", 3: "d8", 4: "e8", 5: "f8", 6: "g8", 7: "h8",
	8: "a7", 9: "b7", 10: "c7", 11: "d7", 12: "e7", 13: "f7", 14: "g7", 15: "h7", 
	16: "a6", 17: "b6", 18: "c6", 19: "d6", 20: "e6", 21: "f6", 22: "g6", 23: "h6",
	24: "a5", 25: "b5", 26: "c5", 27: "d5", 28: "e5", 29: "f5", 30: "g5", 31: "h5", 
	32: "a4", 33: "b4", 34: "c4", 35: "d4", 36: "e4", 37: "f4", 38: "g4", 39: "h4", 
	40: "a3", 41: "b3", 42: "c3", 43: "d3", 44: "e3", 45: "f3", 46: "g3", 47: "h3", 
	48: "a2", 49: "b2", 50: "c2", 51: "d2", 52: "e2", 53: "f2", 54: "g2", 55: "h2", 
	56: "a1", 57: "b1", 58: "c1", 59: "d1", 60: "e1", 61: "f1", 62: "g1", 63: "h1", 
}

const (
	A8 Square = iota; B8; C8; D8; E8; F8; G8; H8
	A7; B7; C7; D7; E7; F7; G7; H7 
	A6; B6; C6; D6; E6; F6; G6; H6
	A5; B5; C5; D5; E5; F5; G5; H5 
	A4; B4; C4; D4; E4; F4; G4; H4 
	A3; B3; C3; D3; E3; F3; G3; H3 
	A2; B2; C2; D2; E2; F2; G2; H2 
	A1; B1; C1; D1; E1; F1; G1; H1 
)

/**********************************************************************
		BITBOARD OPERATIONS (LSB == a8, MSB == h1)
		Board is printed from top to bottom from whites perspective 
**********************************************************************/
func (bitboard Bitboard) GetBit(sq int) int {
	isSet := bitboard & (1 << sq)
	if isSet != 0 {
		return 1
	} else {
		return 0
	}
}

func (bitboard *Bitboard) SetBit(sq int) {
	*bitboard |= (1 << sq)
}

func (bitboard *Bitboard) PopBit(sq int) {
	if *bitboard & (1 << sq) == 0 { return }
	*bitboard ^= (1 << sq)
}

 // retreive number of set bits in bitboard
 func (bitboard Bitboard) CountBits() int {
	ans := 0
	for bitboard != 0 { bitboard &= (bitboard - 1); ans++; }
	return ans
 }

 // retreive the LSB set-bit index (0-index) of a bitboard 
 func (bitboard Bitboard) LSBIndex() int {
	if (bitboard == 0) { return 0 }
	// only leaves the LSB on and then flips all bits behind this bit, 
	bitboard &= -bitboard; bitboard--;
	// this allows counting to return the account index of a bit (0 - indexed)
	return bitboard.CountBits()
 }

 /************************************
 	PRINTING GAMEBOARD OPERATIONS 
 ************************************/
 // prints a bitboard from top to bottom with lsb = a8 and msb = h1
func (bitboard Bitboard) PrintBitboard() {
	// print rank (8, 7, 6...)
	for i := 0; i < 8; i++ {
		// print file (a, b, c...)
		for j := 0; j < 8; j++ {
			if j == 0 {
				fmt.Print(8 - i, " ")
			}
			square := i * 8 + j
			val := (bitboard & (1 << square) != 0)
			if val {
				fmt.Print(1, " ")
			} else {
				fmt.Print(0, " ")
			}
		}
		fmt.Println()
	}
	fmt.Println("  a b c d e f g h")
	fmt.Println("Bitboard as decimal:", bitboard)
 }

 // prints game board from top to bottom with lsb = a8 and msb = h1
 func PrintGameboard() {
	// looping ranks (8, 7, 6, ...)
	for i := 0; i < 8; i++ {
		// looping files (a, b, c, ...)
		for j := 0; j < 8; j++ {
			if j == 0 {
				// ranks printed
				fmt.Print(8 - i, " ")
			}
			sq := i * 8 + j
			var piece Piece = -1

			// loop over piece bitboards
			for bb := WhitePawn; bb <= BlackKing; bb++ {
				if GameBoards[bb].GetBit(sq) == 1 {
					piece = bb
				}
			}

			if piece == -1 {
				fmt.Print(" . ")
			} else {
				fmt.Print(" ", UnicodePieces[piece], " ")
			}
		}
		fmt.Println()
	}
	// files printed
	fmt.Println("   a  b  c  d  e  f  g  h")

	// side to move printed
	if SideToMove == 0 || SideToMove == -1 {
		fmt.Println("    Side To Move: White")
	} else {
		fmt.Println("    Side To Move: Black")
	}

	// enpassant printed
	if (Enpassant != 64) {
		fmt.Println("    Enpassant: ", Enpassant)
	} else {
		fmt.Println("    Enpassant: NO")
	}

	// castling 
	var WK, WQ, BK, BQ string
	if (Castle & White_king_side != 0) { WK = "K" } else { WK = "-"}
	if (Castle & White_queen_side != 0) { WQ = "Q" } else { WQ = "-"}
	if (Castle & Black_king_side != 0) { BK = "k" } else { BK = "-"}
	if (Castle & Black_queen_side != 0) { BQ = "q" } else { BQ = "-"}
	fmt.Println("    Castling: ", WK, WQ, BK, BQ)
	fmt.Println()
 }

func ResetGameStateVariables() {
	// reset game boards
	for i := WhitePawn; i <= BlackKing; i++ {
		GameBoards[i] = 0
	}

	// reset occupancy boards
	for i := 0; i < 3; i++ {
		GameOccupancy[i] = 0
	}
	
	// reset game state
	SideToMove = 0; Enpassant = 64; Castle = 0
}

func ParseFen(fen string) {
	ResetGameStateVariables()
	idx := 0
	// looping ranks (8, 7, 6, ...)
	for i := 0; i < 8; i++ {
		// looping files (a, b, c, ...)
		for j := 0; j < 8; j++ {
			// square 
			sq := i * 8 + j

			// piece is present
			if (string(fen[idx]) >= "A" && string(fen[idx]) <= "Z") || (string(fen[idx]) >= "a" && string(fen[idx]) <= "z") {
				pieceType := AsciiToType[string(fen[idx])]
				GameBoards[pieceType].SetBit(sq)
				idx++
			}

			// file separator
			if string(fen[idx]) >= "0" && string(fen[idx]) <= "9" {
				offset, err := strconv.Atoi(string(fen[idx]))
				if err != nil {
					fmt.Println("ERROR PARSING FEN STRING", err)
					return
				}

				// piece := -1
				// for i := WhitePawn; i <= BlackKing; i++ {
				// 	if GameBoards[i].GetBit(sq) == 1 {
				// 		piece = int(i)
				// 	}
				// }
				// if piece == -1 {
				// 	j--;
				// }

				var present int
				for i := WhitePawn; i <= BlackKing; i++ {
					present |= int(GameBoards[i])
				}
				if present & (1 << sq) == 0 {
					j--
				}

				j += offset
				idx++
			}

			// rank separator
			if string(fen[idx]) == "/" {
				idx++
			}			
		}
	}

	// get side to move
	idx++
	if string(fen[idx]) == "w" {
		SideToMove = White
	} else {
		SideToMove = Black
	}

	// get castling rights 
	idx += 2;
	for string(fen[idx]) != " " {
		if string(fen[idx]) == "K" {
			Castle |= White_king_side
		} else if string(fen[idx]) == "Q" {
			Castle |= White_queen_side
		} else if string(fen[idx]) == "k" {
			Castle |= Black_king_side
		} else if string(fen[idx]) == "q" {
			Castle |= Black_queen_side
		}
		idx++
	}

	// get enpessant 
	idx++
	if string(fen[idx]) != "-" {
		Enpassant = int(StringSquareToBit[string(fen[idx]) + string(fen[idx + 1])])
		idx++
	}

	// loop over white bitboards
	for i := WhitePawn; i <= WhiteKing; i++ {
		GameOccupancy[White] |= GameBoards[i]
	}

	// loop over black bitboards
	for i := BlackPawn; i <= BlackKing; i++ {
		GameOccupancy[Black] |= GameBoards[i]
		GameOccupancy[Both] |= GameBoards[i]
	}

	// set occupancy for both
	GameOccupancy[Both] |= GameOccupancy[White]
	GameOccupancy[Both] |= GameOccupancy[Black]

	fmt.Println("FEN STRING:", fen)
}