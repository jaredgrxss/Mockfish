package engine

import (
	"fmt"
	"strconv"
)

// Types
type Bitboard uint64
type Piece int

// starter FEN strings
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

var Ascii_to_Type = map[string]Piece {
	"P": WhitePawn, "N": WhiteKnight, "B": WhiteBishop, "R": WhiteRook, "Q": WhiteQueen, "K": WhiteKing,
	"p": BlackPawn, "n": BlackKnight, "b": BlackBishop, "r": BlackRook, "q": BlackQueen, "k": BlackKing,
}

const (
	A8 = iota; B8; C8; D8; E8; F8; G8; H8
	A7; B7; C7; D7; E7; F7; G7; H7 
	A6; B6; C6; D6; E6; F6; G6; H6
	A5; B5; C5; D5; E5; F5; G5; H5 
	A4; B4; C4; D4; E4; F4; G4; H4 
	A3; B3; C3; D3; E3; F3; G3; H3 
	A2; B2; C2; D2; E2; F2; G2; H2 
	A1; B1; C1; D1; E1; F1; G1; H1 
)

/**********************************************************************
		BIT OPERATIONS FOR BITBOARD (LSB == a8, MSB == h1)
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
	fmt.Println()
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
 }

 func SetInitBoardState() {
	// white pawns
	GameBoards[WhitePawn].SetBit(A2)
	GameBoards[WhitePawn].SetBit(B2)
	GameBoards[WhitePawn].SetBit(C2)
	GameBoards[WhitePawn].SetBit(D2)
	GameBoards[WhitePawn].SetBit(E2)
	GameBoards[WhitePawn].SetBit(F2)
	GameBoards[WhitePawn].SetBit(G2)
	GameBoards[WhitePawn].SetBit(H2)

	// white knights
	GameBoards[WhiteKnight].SetBit(B1)
	GameBoards[WhiteKnight].SetBit(G1)

	// white bishops
	GameBoards[WhiteBishop].SetBit(C1)
	GameBoards[WhiteBishop].SetBit(F1)

	// white rooks
	GameBoards[WhiteRook].SetBit(A1)
	GameBoards[WhiteRook].SetBit(H1)

	// white queen + king
	GameBoards[WhiteQueen].SetBit(D1)
	GameBoards[WhiteKing].SetBit(E1)

	// black pawns
	GameBoards[BlackPawn].SetBit(A7)
	GameBoards[BlackPawn].SetBit(B7)
	GameBoards[BlackPawn].SetBit(C7)
	GameBoards[BlackPawn].SetBit(D7)
	GameBoards[BlackPawn].SetBit(E7)
	GameBoards[BlackPawn].SetBit(F7)
	GameBoards[BlackPawn].SetBit(G7)
	GameBoards[BlackPawn].SetBit(H7)

	// black knights 
	GameBoards[BlackKnight].SetBit(B8)
	GameBoards[BlackKnight].SetBit(G8)

	// black bishops
	GameBoards[BlackBishop].SetBit(C8)
	GameBoards[BlackBishop].SetBit(F8)

	// black rooks 
	GameBoards[BlackRook].SetBit(A8)
	GameBoards[BlackRook].SetBit(H8)

	// black queen + king
	GameBoards[BlackQueen].SetBit(D8)
	GameBoards[BlackKing].SetBit(E8)

	// side
	SideToMove = 0; Castle = 15;
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
				pieceType := Ascii_to_Type[string(fen[idx])]
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

	idx++
	if string(fen[idx]) == "w" {
		SideToMove = White
	} else {
		SideToMove = Black
	}
	fmt.Println("FEN STRING:", fen)
}