package engine

// CONSTANTS for transposition table
const (
	HashFlagExact int = iota
	HashFlagAlpha
	HashFlagBeta
)

// No hash entry found
const NO_HASH_ENTRY = 100000

// hash table size (10 MB)
const HASH_SIZE = 0x1000000

// data structure for transposition table
type TT struct {
	hashKey uint64 // position identifier
	depth   int    // search depth
	flag    int    // exact / alpha / beta
	score   int    // value of a position
}

// instance of a transposition table
var TranspositionTable [HASH_SIZE]TT

// clear transposition table
func ClearTranspositionTable() {
	for i := 0; i < HASH_SIZE; i++ {
		TranspositionTable[i].hashKey = 0
		TranspositionTable[i].depth = 0
		TranspositionTable[i].flag = 0
		TranspositionTable[i].score = 0
	}
}

// read our cache for the position
func ReadHashData(alpha int, beta int, depth int) int {
	// get instance of hash entry in transposition table
	transpositionTableEntry := &TranspositionTable[HashKey%HASH_SIZE]
	// make sure unique hash matches the entry
	if transpositionTableEntry.hashKey == HashKey {
		// make sure we are at the same depth
		if transpositionTableEntry.depth >= depth {
			// matches exact score
			if transpositionTableEntry.flag == HashFlagExact {
				return transpositionTableEntry.score
			}
			// matches alpha score
			if transpositionTableEntry.flag == HashFlagAlpha &&
				transpositionTableEntry.score <= alpha {
				// fail low node
				return alpha
			}
			// matches beta score
			if transpositionTableEntry.flag == HashFlagBeta &&
				transpositionTableEntry.score >= beta {
				// fail high node
				return beta
			}
		}
	}
	// did not match on cache
	return NO_HASH_ENTRY
}

// write our cache for the position
func WriteHashData(score int, depth int, flag int) {
	// get position data should be in cache
	transpositionTableEntry := &TranspositionTable[HashKey%HASH_SIZE]
	// write to position in cache
	transpositionTableEntry.hashKey = HashKey
	transpositionTableEntry.score = score
	transpositionTableEntry.flag = flag
	transpositionTableEntry.depth = depth
}
