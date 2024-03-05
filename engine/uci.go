package engine

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
				if (move_promoted == int(WhiteQueen) || move_promoted == int(BlackQueen)) || move_string[4] == 'q' {
					return move
				}
				// check for promoted piece being rook
				if (move_promoted == int(WhiteRook) || move_promoted == int(BlackRook)) || move_string[4] == 'r' {
					return move
				}
				// check for promoted piece being bishop
				if (move_promoted == int(WhiteBishop) || move_promoted == int(BlackBishop)) || move_string[4] == 'b' {
					return move
				}
				// check for promoted piece being knight
				if (move_promoted == int(WhiteKnight) || move_promoted == int(BlackKnight)) || move_string[4] == 'n' {
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
