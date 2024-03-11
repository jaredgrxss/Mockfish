# Mockfish
Chess engine written in Go !

# Details
I wanted to attempt to make a chess engine that could reliably beat humans
with advanced optimization and best known techniques for search / evaluation.

Also, I wanted it to be realtively fast (in order to search deeper), and so Go
was my chosen language. 

C++ was considered, but in order to take advantage of such optimizations, I would need to manage memory very efficiently, something for a personal project I was not interested in doing. 

*I more so wanted to explore advanced search, pruning, and bit operations in order to make a good engine.*

# Setup / Running Details
1. Download source code from this repository
2. cd into "Mockfish" directory from your download destination
3. run "go build"
4. This command should produce an executible (Mockfish.exe), you can download a UCI compatible interface and insert this file as your engine of choice


# Testing, Performance, and Techniques
*Below I will try my best to update this brief overview of performance the engine develops at every level, seeing which techniques added the most elo boost and when*

- [Perft Test](https://www.chessprogramming.org/Perft_Results):
    - This is the most basic of testing used to test your engine against millions of pre-calculated positions. Without this testing phase passing, your engine cannot be assured to work at even a basic level. I implemented a perft suite that can be used to test various positions to depths of your choosing to ensure any changes made to the move generator still pass a variety of different positions.
        - ***Elo (N/A) / Elo Gain (N/a)***

- [Static Piece Material Score](https://www.dailychess.com/rival/programming/evaluation.php)
    - This is the technique of assigning each piece. In practice this is quite bad, as when you load up the engine to play, it will simply move pieces around quite unintelligently. There's really little to give it 
        - *** Elo (50) / Elo Gain (+50) ***

- [Static Positional Score](https://www.dailychess.com/rival/programming/evaluation.php)
    - This assigns a certain weight to every position for every different piece type (for instance, kings typically do better towards the edge of the board). This improved performance, but not by much...
        - *** Elo (200) / Elo Gain (+150)***

- [Alpha-Beta / Negamax Search](https://web.mit.edu/6.034/wwwbob/handout3-fall11.pdf)
    - This is an enhanced variation of the minimax which is a exhaustive search algorithim in an alternating game. Minimax looks at all possible positions while Alpha-Beta prunes the search trees so it doesn't waste as much time searching non-optimal branches. 
    - Here is about where it would take a game or two off of myself if I made a huge blunder, but would still fall victim to things like the horizon effect and not spotting killer moves, but it was cool to see it actually starting to beat an unskilled human XD. Also, the longer it thought the better it was, but at this point it was very slow past a depth >= 6 / 7

- [Quiescence Search]()
    - This is a method to only explore interesting search paths in a position that lead to devastating results if not accounted for (such as hanging queens and such), this will allow the engine to not make incredible blunders and increased play, but again not by much
        - *** Elo (400) / Elo gain (+200)***

- [Move Ordering]()
    - 

- [Killer Moves]()


- [History Moves]()


- [Principle Variation]()


- [Iterative Deepening]()


- [Principle Variation Serach]()


- [Late Move Reduction]()


- [Null Move Pruning]()


- [Zobrist Hashing / Transposition Tables]()