package engine

import (
	"fmt"
)

// Bitboard typing
type Bitboard uint64

/*

	HELPFUL... MAY BE NEEDED

	var squareToBit = map[string]uint64 {
	"a8": 0, "b8": 1, "c8": 2, "d8": 3, "e8": 4, "f8": 5, "g8": 6, "h8": 7, 
	"a7": 8, "b7": 9, "c7": 10, "d7": 11, "e7": 12, "f7": 13, "g7": 14, "h7": 15, 
	"a6": 16, "b6": 17, "c6": 18, "d6": 19, "e6": 20, "f6": 21, "g6": 22, "h6": 23, 
	"a5": 24, "b5": 25, "c5": 26, "d5": 27, "e5": 28, "f5": 29, "g5": 30, "h5": 31, 
	"a4": 32, "b4": 33, "c4": 34, "d4": 35, "e4": 36, "f4": 37, "g4": 38, "h4": 39, 
	"a3": 40, "b3": 41, "c3": 42, "d3": 43, "e3": 44, "f3": 45, "g3": 46, "h3": 47, 
	"a2": 48, "b2": 49, "c2": 50, "d2": 51, "e2": 52, "f2": 53, "g2": 54, "h2": 55, 
	"a1": 56, "b1": 57, "c1": 58, "d1": 59, "e1": 60, "f1": 61, "g1": 62, "h1": 63,
}

	const (
		A1, B1, C1, D1, E1, F1, G1, H1 = 0, 1, 2, 3, 4, 5, 6, 7
		A2, B2, C2, D2, E2, F2, G2, H2 = 8, 9, 10, 11, 12, 13, 14, 15
		A3, B3, C3, D3, E3, F3, G3, H3 = 16, 17, 18, 19, 20, 21, 22, 23
		A4, B4, C4, D4, E4, F4, G4, H4 = 24, 25, 26, 27, 28, 29, 30, 31
		A5, B5, C5, D5, E5, F5, G5, H5 = 32, 33, 34, 35, 36, 37, 38, 39
		A6, B6, C6, D6, E6, F6, G6, H6 = 40, 41, 42, 43, 44, 45, 46, 47
		A7, B7, C7, D7, E7, F7, G7, H7 = 48, 49, 50, 51, 52, 53, 54, 55
		A8, B8, C8, D8, E8, F8, G8, H8 = 56, 57, 58, 59, 60, 61, 62, 63
	)

	"a8" "b8" "c8" "d8" "e8" "f8" "g8" "h8"
	"a7" "b7" "c7" "d7' "e7" "f7" "g7" "h7"
	"a6" "b6" "c6" "d6" "e6" "f6" "g6" "h6"
	"a5" "b5" "c5" "d5" "e5" "f5" "g5" "h5"
	"a4" "b4" "c4" "d4" "e4" "f4" "g4" "h4"
	"a3" "b3" "c3" "d3" "e3" "f3" "g3" "h3"
	"a2" "b2" "c2" "d2" "e2" "f2" "g2" "h2"
	"a1" "b1" "c1" "d1" "e1" "f1" "g1" "h1"

*/

/**********************************************************************
		- BIT OPERATIONS FOR BITBOARD (LSB == a8, MSB == h1)
		- Board is printed from top to bottom from whites perspective 
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
			if val == true {
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