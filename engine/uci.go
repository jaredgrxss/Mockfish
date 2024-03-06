package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/************************************************
	ALL COMMANDS TO INTERACT WITH UCI PROTOCOL
************************************************/

// parse the GUI move string given (i.e "e7e8q")
func ParseUCIMove(move_string string) int {

	// create moves list & generate current board state moves
	var moves Moves
	GeneratePositionMoves(&moves)

	// get the source square from the uci given string
	move_source := int((move_string[0] - 'a') + (8-(move_string[1]-'0'))*8)
	// get the target square from the uci given string
	move_target := int((move_string[2] - 'a') + (8-(move_string[3]-'0'))*8)

	// loop over given moes for this position
	for i := 0; i < moves.Move_count; i++ {
		// get the current move
		move := moves.Move_list[i]
		// decode moves
		source, target, _, promo, _, _, _, _ := DecodeMove(move)
		if source == move_source && target == move_target {
			move_promoted := promo
			if move_promoted != 0 {
				// check for promoted piece being queen
				if (move_promoted == int(WhiteQueen) || move_promoted == int(BlackQueen)) && move_string[4] == 'q' {
					return move
				}
				// check for promoted piece being rook
				if (move_promoted == int(WhiteRook) || move_promoted == int(BlackRook)) && move_string[4] == 'r' {
					return move
				}
				// check for promoted piece being bishop
				if (move_promoted == int(WhiteBishop) || move_promoted == int(BlackBishop)) && move_string[4] == 'b' {
					return move
				}
				// check for promoted piece being knight
				if (move_promoted == int(WhiteKnight) || move_promoted == int(BlackKnight)) && move_string[4] == 'n' {
					return move
				}
				// try to find wrong promotion before returning out
				continue
			}
			// return legal move
			return move
		}
	}
	// this was an illegal move passed
	return 0
}

/*
------------------ UCI POSITION COMMANDS ------------------

	make the start position of the board
	---- example command ----
	position startpos

	intialize board from the start position and then make moves
	---- example command ----
	position startpos moves a1a2 e7e6

	parse any uci fen that is given
	---- example command ----
	position fen r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1

	parse any uci fen that is given and then make moves
	---- example command ----
	position fen r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1 moves e2a6 e8g8
*/
func ParseUCIPosition(command string) {
	// make sure this is a position command
	position, _ := regexp.MatchString("position", command)
	if !position {
		fmt.Println("Not a UCI position command.")
		return
	}
	// see if the start position was asked for
	startpos, _ := regexp.MatchString("startpos", command)
	if startpos {
		// this case put chess board to initial position
		ParseFen(START_POSITION)
	}
	// see if a FEN string was passed
	fen, _ := regexp.MatchString("fen", command)
	if fen {
		// gather the fen string and parse it
		re := regexp.MustCompile("fen")
		idx := re.FindStringIndex(command)[1] + 1
		fmt.Println("An UCI FEN string was passed.")
		ParseFen(command[idx:])

	}
	// default behavior, don't look any further and just load starting position
	if !startpos && !fen {
		ParseFen(START_POSITION)
		return
	}
	// check to see if we have any moves
	uci_moves, _ := regexp.MatchString("moves", command)
	if uci_moves {
		// find the index of where the moves command ends
		re := regexp.MustCompile("moves")
		idx := re.FindStringIndex(command)[1] + 1
		// gather the moves at the end of the string
		moves := strings.Split(command[idx:], " ")
		for i := 0; i < len(moves); i++ {
			// parse the move given by UCI string w/ helper
			move := ParseUCIMove(moves[i])
			// illegal moves was passed, break out
			if move == 0 {
				break
			}
			MakeMove(move, allMoves)

		}
	}

}

/*
command to make engine search for best move
FIXED DEPTH = 6 at the moment
*/
func ParseUCIGo(command string) {
	// find depth command and parse value
	var depth int
	find_depth, _ := regexp.MatchString("depth", command)
	if find_depth {
		re := regexp.MustCompile("depth")
		idx := re.FindStringIndex(command)[1] + 1
		// parse out the depth integer from our command
		found_depth, err := strconv.Atoi(string(command[idx:]))
		// some error occured
		if err != nil {
			fmt.Println("Error reading in depth")
			return
		}
		depth = found_depth
	} else {
		depth = 6
	}

	// search best move in given position
	// SearchPosition(depth)
	fmt.Println(depth)
}
