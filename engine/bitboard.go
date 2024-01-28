package engine

import (
	"fmt"
)

// Types
type Bitboard uint64
type Piece int


/****************************************
	BOARD REPRESENTATION HELPERS
****************************************/
var GameBoards [12]Bitboard // 12 boards for each piece + color
var GameOccupancy [3]Bitboard // white = 0, black = 1, both = 2
var SideToMove = -1
var Enpassant int = 65

/*
	castling representation info... 

	bin		dec		description
	
	0001 	1 		white king can castle king side
	0010	2		white king can castle queen side
	0100	4		black king can castle king side
	1000	8		black king can caslte queen side
*/

var Castle int
const White_king_side, White_queen_side, Wlack_king_side, Black_queen_side = 1, 2, 4, 8

/******************************
		ENUMERATIONS
******************************/
const ( White = iota; Black; Both)

const ( WhitePawn Piece = iota; WhiteKnight; WhiteBishop; WhiteRook; WhiteQueen; WhiteKing 
		BlackPawn; BlackKnight; BlackBishiop; BlackRook; BlackQueen; BlackKing )

// can use either ASCII or Unicode for board debugging
var AsciiPieces string = "PNBRQKpnbrqk" // must be type casted with string()

var UnicodePieces [12]string = [12]string {
	"♙", "♘", "♗", "♖", "♕", "♔", 
	"♟︎", "♞", "♝", "♜", "♛", "♚",
}

var Ascii_to_Type = map[string]Piece {
	"P": WhitePawn, "N": WhiteKnight, "B": WhiteBishop, "R": WhiteRook, "Q": WhiteQueen, "K": WhiteKing,
	"p": BlackPawn, "n": BlackKnight, "b": BlackBishiop, "r": BlackRook, "q": BlackQueen, "k": BlackKing,
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

 // prints board from top to bottom with lsb = a8 and msb = h1
func (bitboard Bitboard) PrintBoard() {
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