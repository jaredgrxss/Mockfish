# Mockfish
Chess engine written in Go !

# Details
I wanted to attempt to make a chess engine that could reliably beat humans
with advanced optimization and best known techniques for search / evaluation.

Also, I wanted it to be realtively fast (in order to search deeper), and so Go
was my chosen language. 

C++ was considered, but in order to take advantage of such optimizations, I would need to manage memory very efficiently, something for a personal project I was not interested in doing. 

*I more so wanted to explore advanced search, pruning, and bit operations in order to make a good engine.*

# All Added Techniques
- Board & Misc
    - [Bitboards](https://www.chessprogramming.org/Bitboards)
    - [Magic bitboards](https://www.chessprogramming.org/Magic_Bitboards#:~:text=Magic%20bitboards%20applies%20perfect%20hashing,different%2C%20but%20redundant%20outer%20squares.)
    - [Transposition Tables](https://www.chessprogramming.org/Transposition_Table)
    - [Attack Tables](https://www.chessprogramming.org/Attack_and_Defend_Maps)
    - [Opening Books]()
- Search 
    - Alpha-Beta Pruing
    - Iterative Deepening
    - Move Ordering 
        -
- Evaluation
    - Coming soon...

# Setup / Running Details
1. Download source code from this repository
2. cd into "Mockfish" directory from your download destination
3. run "go build"
4. This command should produce an executible (Mockfish.exe), you can download a UCI compatible interface and insert this file as your engine of choice


# Testing / Performance
*Below I will try my best to update this brief overview of performance the engine develops at every level, seeing which techniques added the most elo boost and when*

- [Perft Test](https://www.chessprogramming.org/Perft_Results):
    - This is the most basic of testing used to test your engine against millions of pre-calculated positions. Without this testing phase passing, your engine cannot be assured to work at even a basic level. I implemented a perft suite that can be used to test various positions to depths of your choosing to ensure any changes made to the move generator still pass a variety of different positions.
        - ***Elo (N/A) / Elo Gain (N/a)***

- Static Piece Material Score 
    - This is the technique of assigning each piece. In practice this is quite bad, as when you load up the engine to play, it will simply move pieces around quite unintelligently. There's really little to give it 
        - *** Elo (50) / Elo Gain (+50) ***

- Static Positional Score
    - This assigns a certain weight to every position for every different piece type (for instance, kings typically do better towards the edge of the board). This improved performance, but not by much...
        - *** Elo (100) / Elo Gain (+50)***
