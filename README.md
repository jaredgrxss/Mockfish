# Mockfish
A bitboard chess engine written in Go!

# Details
I wanted to attempt to make a chess engine that could reliably beat humans
with advanced optimization and best known techniques for search / evaluation.

Also, I wanted it to be *relatively* fast, and so Go was chosen to implement this engine. 

C++ was considered, but in order to take advantage of such optimizations, I would need to manage memory very efficiently as well as be clever with optimizing functions, something for a personal project I was not interested in doing. 

*I more so wanted to explore advanced search, pruning, bit operations, and algorithms in order to make a good engine.*

# Setup / Running Details
1. Download source code from this repository
2. cd into "Mockfish" directory from your download destination
3. run "go build" (assuming you have Golang installed on your machine)
4. This command should produce an executible (Mockfish.exe), you can download a UCI compatible interface and insert this file as your engine of choice


# Testing, Performance, and Techniques
*Below I will try my best to update this brief overview of performance the engine develops at every level, seeing which techniques added the most elo boost and when*

- [Perft Test](https://www.chessprogramming.org/Perft_Results):
    - This is the most basic of testing used to test your engine against millions of pre-calculated positions. Without this testing phase passing, your engine cannot be assured to work at even a basic level. I implemented a perft suite that can be used to test various positions to depths of your choosing to ensure any changes made to the move generator still pass a variety of different positions.
        - ***Elo (N/A) / Elo Gain (N/a)***

- [Static Piece Material Score](https://www.dailychess.com/rival/programming/evaluation.php)
    - This is the technique of assigning each piece. In practice this is quite bad, as when you load up the engine to play, it will simply move pieces around quite unintelligently. There's really little to give it 
        - ***Elo (100) / Elo Gain (+100) ***

- [Static Positional Score](https://www.dailychess.com/rival/programming/evaluation.php)
    - This assigns a certain weight to every position for every different piece type (for instance, kings typically do better towards the edge of the board). This improved performance, but not by much...
        - ***Elo (~400) / Elo Gain (+300)***

- [Alpha-Beta / Negamax Search](https://web.mit.edu/6.034/wwwbob/handout3-fall11.pdf)
    - This is an enhanced variation of the minimax which is a exhaustive search algorithim in an alternating game. Minimax looks at all possible positions while Alpha-Beta prunes the search trees so it doesn't waste as much time searching non-optimal branches. 
    - Here is about where it would take a game or two off of myself if I made a huge blunder, but would still fall victim to things like the horizon effect and not spotting killer moves, but it was cool to see it actually starting to beat an unskilled human XD. Also, the longer it thought the better it was, but at this point it was very slow past a depth >= 6 / 7
        - ***Elo (~500) / Elo Gain (+100)***

- [Quiescence Search](https://adamberent.com/quiescencesearch-andextensions/)
    - This is a method to only explore interesting search paths in a position that lead to devastating results if not accounted for (such as hanging queens and such), this will allow the engine to not make incredible blunders and increased play, but again not by much
        - ***Elo (~850) / Elo Gain (+350)***

- [Move Ordering](https://www.chessprogramming.org/Move_Ordering)
    - Without proper move ordering, an engine can only hope to reach depths of 6 to 7 in any given position due to the fact that most good moves might be considered at random instead of in order. With this technique I was able to get the engine to search far deeper than before, and faster, allowing it to consider more outcomes of the current position
        - ***Elo (~1300) / Elo Gain (+450)***

- [Principle Variation Search](https://www.chessprogramming.org/Principal_Variation_Search)
    - When implementing principle variation search, you are allowing the engine to search it allows for a more efficient alpha-beta search to occur as the PV-node is the only one search in the full window
        - ***Elo (~1550) / Elo Gain (250) ***

- [Iterative Deepening](https://www.educative.io/answers/what-is-iterative-deepening-search)
    - With iterative deepening we can use the results of our previous search at a depth n - 1 to make enhanced decisions about the n-th layer (with principle variation search for instance), this allows also for time controls, although i didn't implement them...
        - ***Elo (~1600) / Elo Gain (50)*** 

- [Late Move Reduction](https://mediocrechess.blogspot.com/2007/03/other-late-move-reduction-lmr.html)
    - With late move reduction, you can search at a lower depth for moves >= n (n being a certain number of moves searched) as we expect the best moves usually to be in the beginning (thanks to move ordering). This sped the engine up quite a bit, reducing number of nodes search at a depth of 10 from ~13,000,000 to ~1,000,000, making it able to search further ahead, especially when their were less pieces on the board. 
        - ***Elo (~1750) / Elo Gain (150)***

- [Null Move Pruning](https://www.chessprogramming.org/Null_Move_Pruning)
    - This technique allows the chosen side to give the opponent a "free move". In theory, if a position is still very good for side X after 2 or 3 moves by side Y, side X will want to rank these outcomes higher or lower, again further enchancing the pruning dune by our search
        - ***Elo (~1850) / Elo Gain (100)***

- [Zobrist Hashing / Transposition Tables](https://dev.to/larswaechter/zobrist-hashing-72n)
    - Dynamic Programming in real life! This is essentially a cache that the search first looks to see if we have been in a given position at an earlier part of the search, and if so, we can return the outcome of this position before having to search to depth N again. This greatly enhanced performance, and like most caches, it depends on your hit rate. I found that with a cache size of 10MB seemed to be the sweet spot for my engine, but I've seen some engines use as small as a 4MB cache. 
        - ***Elo (~2000) / Elo Gain (150)***

- [Advanced Evaluation]()
    - With a lot of the optimizations in search complete, I wanted the last leg of this journey to be centered around improving the evaluation function that the engine used. A lot of these techniques involved simple bit manipulations to check for open files for rooks, penalties for doubled pawns, and king safety discouraging king moves, with a lot of these simple things in place, the engine saw its final boost and was able to take down consistently bots and players rated ~2200-2300!
    - You could take this even further by adding stockfish's NNUE to your engine, but that felt a bit like cheating, so I left it out :) 
        - ***Elo (2300) / Elo Gain (300)***