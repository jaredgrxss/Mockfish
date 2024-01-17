package engine
import (
	"fmt"
)
/********************************************************
	PRECOMPUTED ATTACK INFO USED FOR MOVE GENERATION
********************************************************/

/* 
	NOT A FILE 
		8 0 1 1 1 1 1 1 1 
		7 0 1 1 1 1 1 1 1 
		6 0 1 1 1 1 1 1 1 
		5 0 1 1 1 1 1 1 1 
		4 0 1 1 1 1 1 1 1 
		3 0 1 1 1 1 1 1 1 
		2 0 1 1 1 1 1 1 1 
		1 0 1 1 1 1 1 1 1 
		a b c d e f g h


	NOT H FILE 

		8 1 1 1 1 1 1 1 0 
		7 1 1 1 1 1 1 1 0 
		6 1 1 1 1 1 1 1 0 
		5 1 1 1 1 1 1 1 0 
		4 1 1 1 1 1 1 1 0 
		3 1 1 1 1 1 1 1 0 
		2 1 1 1 1 1 1 1 0 
		1 1 1 1 1 1 1 1 0 
		a b c d e f g h

*/

/**********************
LEAPING PIECE INFO
**********************/

// constants to help with correctly generating leaper piece moves in case they wrap around
const NOT_A_FILE Bitboard = 18374403900871474942
const NOT_H_FILE Bitboard = 9187201950435737471
const NOT_HG_FILE Bitboard = 4557430888798830399
const NOT_AB_FILE Bitboard = 18229723555195321596

// constants to store attack maps for leaping pieces
const White, Black = 0, 1
var PawnAttacks [2][64]Bitboard
var KnightAttacks [64]Bitboard
var KingAttacks [64]Bitboard


/***************************************
	MAGIC BITBOARD / SLIDING PIECE INFO
***************************************/

// used when generating magic numbers for a certain piece
const Rook, Bishop = 0, 1

// total bits a bishop can move given a sqaure
var BitCountBishop [64]int = [64]int {
	6, 5, 5, 5, 5, 5, 5, 6, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 7, 7, 7, 7, 5, 5, 
	5, 5, 7, 9, 9, 7, 5, 5, 
	5, 5, 7, 9, 9, 7, 5, 5, 
	5, 5, 7, 7, 7, 7, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	6, 5, 5, 5, 5, 5, 5, 6, 
}
// total bits a rook can move given a square
var BitCountRook [64]int = [64]int {
	12, 11, 11, 11, 11, 11, 11, 12, 
	11, 10, 10, 10, 10, 10, 10, 11, 
	11, 10, 10, 10, 10, 10, 10, 11, 
	11, 10, 10, 10, 10, 10, 10, 11, 
	11, 10, 10, 10, 10, 10, 10, 11, 
	11, 10, 10, 10, 10, 10, 10, 11, 
	11, 10, 10, 10, 10, 10, 10, 11, 
	12, 11, 11, 11, 11, 11, 11, 12, 
}

// magic numbers for rook
var RookMagics [64]uint64 = [64]uint64 {
	0x80051082604004,
	0x40200040001000,
	0x880100180486000,
	0x1080040800801000,
	0x280028004000800,
	0x5180018004000200,
	0x80020001000080,
	0x200020348802504,
	0x6008800840008020,
	0x4842804001802004,
	0x2542002200801040,
	0x1101001000090020,
	0x6000600201008,
	0x2000410080200,
	0x12000401020088,
	0x16500021040a100,
	0x4415608002804000,
	0x10004040002001,
	0x1000110041002000,
	0x1010010040821,
	0x94010100080010,
	0x480808004000200,
	0x40040010880201,
	0x20420001004084,
	0x420400080002088,
	0x20200240025001,
	0x8010006060040800,
	0x200090100100020,
	0x7241000500502800,
	0x6040080800200,
	0x1302010400104802,
	0x19000100008042,
	0x6280042004400048,
	0x40401000402000,
	0x2081002001001040,
	0x1001001000821,
	0x400041101000800,
	0x2000802000410,
	0x2000081004000281,
	0x2000140042000081,
	0x400400381228005,
	0x20002040008080,
	0x410040800202002,
	0x4c08008010008008,
	0x2808008004008008,
	0x641000400890002,
	0x20004010100,
	0x861005404820025,
	0x4082050041298200,
	0x40400280200180,
	0x4220100820008080,
	0x890028800918080,
	0x840408021001002,
	0x81000400422900,
	0x69800100020080,
	0x5004400810200,
	0x880800110402b01,
	0x440402100820012,
	0x8085402001000815,
	0x8802040861500101,
	0x422000508201002,
	0xa000c15500822,
	0x8350013000820824,
	0x8482084489002402,   
}


// magic numbers for bishop
var BishopMagics [64]uint64 = [64]uint64 {
	0x8020815401104200,
	0x8048108410404020,
	0x42040d08204400,
	0x109082000c2024,
	0x1104044020020,
	0x4c02011008404111,
	0x3011040282c00430,
	0x400802c0b200840,
	0x850100401840400,
	0x80110400948208,
	0x4080204202520,
	0x2084080a00202090,
	0x113062110001000,
	0x9000020124601080,
	0x8241202000,
	0x1120051041042000,
	0x11002022220840,
	0x2400404041400,
	0x401040204010200,
	0x4404000844020900,
	0x84006206110301,
	0xa042000100420215,
	0x20850444442000,
	0xa800840580880,
	0x4204402020828440,
	0x14121204900400,
	0x4110410010040280,
	0x800300c08c004200,
	0x611001025004000,
	0x8004112006200,
	0x220840088842400,
	0xc4403840900b080,
	0x8008380888052000,
	0x2280415210900,
	0x601080200010400,
	0x80080a0080080380,
	0x40004100401100,
	0x101105200210100,
	0x1001080a00008208,
	0x5402009900882c00,
	0x8680804012880,
	0xc24040202410800,
	0x801004022201004,
	0x22080c2011009802,
	0xa080100400400,
	0x2001102100420202,
	0x2020024242000040,
	0x4010820600400028,
	0x402480208a1a004,
	0x981004862080100,
	0x80602128080000,
	0x2801068084040c00,
	0x400004030248000,
	0x2400408408901,
	0x20408080804a500,
	0x8020880080a08030,
	0x2061402010420c0,
	0x206400484410080a,
	0x9002800200840412,
	0x4000020004840421,
	0x18000c06d0112,
	0x800082052040848,
	0x80508450068228,
	0x8600102021410,
}


/************************
	HELPER FUNCTIONS 
************************/

// bitshifts for pawns is = 7,9
func genPawnAttacks(side uint64, sq int) Bitboard {

	// resulting bitboard of attacks
	var attacks Bitboard = 0

	// make new bitboard and set the pos of pawn on it
	var bitboard Bitboard = 0
	bitboard.SetBit(sq)

	// white side (0) black side (1)
	if side == 0 {
		if bitboard >> 7 & NOT_A_FILE != 0 {
			attacks |= (bitboard >> 7)
		}
		if bitboard >> 9 & NOT_H_FILE != 0 {
			attacks |= (bitboard >> 9)
		}
	} else {
		if bitboard << 7 & NOT_H_FILE != 0 {
			attacks |= (bitboard << 7)
		}
		if (bitboard << 9 & NOT_A_FILE != 0) {
			attacks |= (bitboard << 9)
		}
	}
	return attacks
}

// bitshifts for knights = 6, 10, 15, 17
func genKnightAttacks(sq int) Bitboard {
	// returning attack board
	var attacks Bitboard = 0

	// set bit for knight
	var bitboard Bitboard = 0
	bitboard.SetBit(sq)
	
	// calculating attacks
	if (bitboard >> 17) & NOT_H_FILE != 0 {
		attacks |= (bitboard >> 17)
	}
	if (bitboard >> 15) & NOT_A_FILE != 0 {
		attacks |= (bitboard >> 15)
	}
	if (bitboard >> 10) & NOT_HG_FILE != 0 {
		attacks |= (bitboard >> 10)
	}
	if (bitboard >> 6) & NOT_AB_FILE != 0 {
		attacks |= (bitboard >> 6)
	}

	// oposite side bit shifts
	if (bitboard << 17) & NOT_A_FILE != 0 {
		attacks |= (bitboard << 17)
	}
	if (bitboard << 15) & NOT_H_FILE != 0 {
		attacks |= (bitboard << 15)
	}
	if (bitboard << 10) & NOT_AB_FILE != 0 {
		attacks |= (bitboard << 10)
	}
	if (bitboard << 6) & NOT_HG_FILE != 0 {
		attacks |= (bitboard << 6)
	}

	return attacks
}

// bitshifts for king is = 1, 7, 8, 9
func genKingAttacks(sq int) Bitboard {
	// return attack board
	var attacks Bitboard = 0
	
	// set king on temp board
	var bitboard Bitboard = 0
	bitboard.SetBit(sq)

	// calculating attacks
	// up / down don't need checks
	attacks |= (bitboard << 8)
	attacks |= (bitboard >> 8)

	if (bitboard >> 1) & NOT_H_FILE != 0 {
		attacks |= (bitboard >> 1)
	}
	if (bitboard >> 7) & NOT_A_FILE != 0 {
		attacks |= (bitboard >> 7)
	}
	
	if (bitboard >> 9) & NOT_H_FILE != 0 {
		attacks |= (bitboard >> 9)
	}

	// oposite side bit shifts
	if (bitboard << 1) & NOT_A_FILE != 0 {
		attacks |= (bitboard << 1)
	}
	if (bitboard << 7) & NOT_H_FILE != 0 {
		attacks |= (bitboard << 7)
	}

	if (bitboard << 9) & NOT_A_FILE != 0 {
		attacks |= (bitboard << 9)
	}

	return attacks
}

// function to generate bishop attack rays (assuming NO blockers)
func genBishopAttacks(sq int) Bitboard {
	// return board for all valid attacks
	var attacks Bitboard = 0

	var rank, file int

	var target_rank = sq / 8; var target_file = sq % 8

	// mask relevant occupancy bishop bits
	for rank, file = target_rank + 1, target_file + 1; rank < 7 && file < 7; rank, file = rank + 1, file + 1 {
		attacks |= (1 << (rank * 8 + file))
	}
	for rank, file = target_rank - 1, target_file + 1; rank > 0 && file < 7; rank, file = rank - 1, file + 1 {
		attacks |= (1 << (rank * 8 + file))
	}
	for rank, file = target_rank + 1, target_file - 1; rank < 7 && file > 0; rank, file = rank + 1, file - 1 {
		attacks |= (1 << (rank * 8 + file))
	}
	for rank, file = target_rank - 1, target_file - 1; rank > 0 && file > 0; rank, file = rank - 1, file - 1 {
		attacks |= (1 << (rank * 8 + file))
	}
	return attacks
}

// function to generate rook attack rays (assuming NO blockers)
func genRookAttacks(sq int) Bitboard {
	var attacks Bitboard = 0

	var rank, file int


	var target_rank = sq / 8; var target_file = sq % 8

	// mask relevant occupancy rook bits
	for rank = target_rank + 1; rank < 7; rank++ {
		attacks |= (1 << (rank * 8 + target_file))
	}
	for rank = target_rank - 1; rank > 0; rank-- {
		attacks |= (1 << (rank * 8 + target_file))
	}
	for file = target_file + 1; file < 7; file++ {
		attacks |= (1 << (target_rank * 8 + file))
	}
	for file = target_file - 1; file > 0; file-- {
		attacks |= (1 << (target_rank * 8 + file))
	}

	return attacks
}

// function to generate bishop attack rays  given a blocker configuration
func onTheFlyBishopAttacks(sq int, blockers Bitboard) Bitboard {
	// return board for all valid attacks
	var attacks Bitboard = 0

	var rank, file int

	var target_rank = sq / 8; var target_file = sq % 8

	// gen bishop actual attacks given a blocker config, if we hit a blocker, break
	for rank, file = target_rank + 1, target_file + 1; rank < 8 && file < 8; rank, file = rank + 1, file + 1 {
		attacks |= (1 << (rank * 8 + file))
		if (1 << (rank * 8 + file)) & blockers != 0 { break }
	}
	for rank, file = target_rank - 1, target_file + 1; rank >= 0 && file < 8; rank, file = rank - 1, file + 1 {
		attacks |= (1 << (rank * 8 + file))
		if (1 << (rank * 8 + file)) & blockers != 0 { break }
	}
	for rank, file = target_rank + 1, target_file - 1; rank < 8 && file >= 0; rank, file = rank + 1, file - 1 {
		attacks |= (1 << (rank * 8 + file))
		if (1 << (rank * 8 + file)) & blockers != 0 { break }
	}
	for rank, file = target_rank - 1, target_file - 1; rank >= 0 && file >= 0; rank, file = rank - 1, file - 1 {
		attacks |= (1 << (rank * 8 + file))
		if (1 << (rank * 8 + file)) & blockers != 0 { break }
	}
	return attacks
}

// function to generate rook attack rays given a blocker configuration
func onTheFlyRookAttacks(sq int, blockers Bitboard) Bitboard {
	// return board for all valid attacks given blockers
	var attacks Bitboard = 0

	var rank, file int

	var target_rank = sq / 8; var target_file = sq % 8

	// gen rook actual attacks given a blocker config, if we hit a blocker, break
	for rank = target_rank + 1; rank < 8; rank++ {
		attacks |= (1 << (rank * 8 + target_file))
		if (1 << (rank * 8 + target_file)) & blockers != 0 { break }
	}
	for rank = target_rank - 1; rank >= 0; rank-- {
		attacks |= (1 << (rank * 8 + target_file))
		if (1 << (rank * 8 + target_file)) & blockers != 0 { break }
	}
	for file = target_file + 1; file < 8; file++ {
		attacks |= (1 << (target_rank * 8 + file))
		if (1 << (target_rank * 8 + file)) & blockers != 0 { break }
	}
	for file = target_file - 1; file >= 0; file-- {
		attacks |= (1 << (target_rank * 8 + file))
		if (1 << (target_rank * 8 + file)) & blockers != 0 { break }
	}
	
	return attacks
}

// find occupancy board/bits on a given attack board for sliding pieces
// (4096 or 2^12 possible configurations have to consider, which is not a lot)
func setOccupancy(idx int, bitCnt int, attackMask Bitboard) Bitboard {
	// our final board that will given occupant bits
	var occupancy Bitboard = 0

	// loop for number of bits we have in a current board
	for i := 0; i < bitCnt; i++ {
		// get LSB idx of attackMask then pop it off
		sq := attackMask.LSBIndex(); attackMask.PopBit(sq)
		// set occupant sq if it is present 
		if idx & (1 << i) != 0 {
			occupancy |= (1 << sq)
		}
	}

	// return back this occupancy configuration
	return occupancy
}


/***********************************************
	HELPERS FOR MAGIC BITBOARD HASHKEY
***********************************************/

// state for PRNG func
var random_state uint32 = 1010243240

// psuedo 32 bit random number generator for PRNG 
func Get_u32_rand() uint32 {
	// retreive current state
	x := random_state

	// xor shift 32 algorithm
	x ^= (x << 13)
	x ^= (x >> 17)
	x ^= (x << 5)
	random_state = x 

	// return PRN
	return x
}

// psudo 64 bit random number generator 
func Get_u64_rand() uint64 {
	// get 4 random numbers
	var n1, n2, n3, n4 uint64
	n1 = uint64(Get_u32_rand()) & 0xFFFF; n2 = uint64(Get_u32_rand()) & 0xFFFF
	n3 = uint64(Get_u32_rand()) & 0xFFFF; n4 = uint64(Get_u32_rand()) & 0xFFFF

	// return 64 bit PRNG
	return n1 | (n2 << 16) | (n3 << 32) | (n4 << 48)
}

// count bits of a regular unsigned 64 bit number
func count_bits(number uint64) int {
	ans := 0
	for number != 0 { number &= (number - 1); ans++ }
	return ans
}

// using xor32algo and xor64algo, generate magic numbers
func GenerateMagicNumber() uint64 {
	return Get_u64_rand() & Get_u64_rand() & Get_u64_rand()
}

// generating magic number candidates
func find_magic_number(sq int, bitCnt int, flag int) uint64 {
	// get all occupancies for all configurations
	var occupancies [4096]Bitboard

	// attack tables
	var attacks [4096]Bitboard

	// used attacks
	var usedAttacks [4096]Bitboard

	// get attack rays for current piece
	var pieceAttacks Bitboard
	if (flag == 0) {
		pieceAttacks = genRookAttacks(sq)
	} else {
		pieceAttacks = genBishopAttacks(sq)
	}

	var occupant_indices = 1 << bitCnt
	for i := 0; i < occupant_indices; i++ {
		// get occupancies for this configuration
		occupancies[i] = setOccupancy(i, bitCnt, pieceAttacks)

		if (flag == 0) {
			attacks[i] = onTheFlyRookAttacks(sq, occupancies[i])
		} else {
			attacks[i] = onTheFlyBishopAttacks(sq, occupancies[i])
		}
	}

	for i := 0; i < 100000000; i++ {
		// lets generate magic number candidate
		var magic uint64 = GenerateMagicNumber()

		// pass magic numbers that won't work
		if (count_bits((uint64(pieceAttacks) * magic) & 0xFF00000000000000) < 6) { continue }

		// init our used attacks
		for j := 0; j < len(usedAttacks); j++ { usedAttacks[j] = 0 }

		var index, fail int
		for index, fail = 0, 0; fail == 0 && index < occupant_indices; index++ {

			
			magic_index := int((occupancies[index] * Bitboard(magic)) >> (64 - bitCnt))
			
			// if magic index works
			if (usedAttacks[magic_index] == 0) {
				// set usedAttacks to the attacks of this current square
				usedAttacks[magic_index] = attacks[index]
			} else if usedAttacks[magic_index] != attacks[index] {
				// magic index doesn't work
				fail = 1
			}
		}
		if (fail == 0) {
			return magic
		}
	}
	fmt.Println("MAGIC NUMBER FAILS")
	return 0
}

// func Init_magics() {
// 	// lets loop over 64 board squares 
// 	for i := 0; i < 64; i++ {
// 		// init rook magic numbers
// 		fmt.Printf(" 0x%x,\n", find_magic_number(i, BitCountRook[i], 0))
// 	}
// 	fmt.Println()
// 	for i := 0; i < 64; i++ {
// 		fmt.Printf(" 0x%x,\n",find_magic_number(i, BitCountBishop[i], 1))
// 	}
// }



/****************************
	PRECOMPUTE ATTACKS 	
****************************/
func GeneratePieceAttacks() {
	for i := 0; i < 64; i++ {
		// pawns white / black
		PawnAttacks[0][i] = genPawnAttacks(0, i)
		PawnAttacks[1][i] = genPawnAttacks(1, i)

		// knight attacks
		KnightAttacks[i] = genKnightAttacks(i)

		// king attacks
		KingAttacks[i] = genKingAttacks(i)
	}
}