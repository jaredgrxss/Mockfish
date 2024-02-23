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
    - Opening Books
- Search 
    - Coming soon...
- Evaluation
    - Coming soon...

# Setup / Running Details
1. Download source code from this repository
2. cd into "Mockfish" directory from your download destination
3. run "go build"
4. This command should produce an executible (Mockfish.exe), you can download a UCI compatible interface and insert this file as your engine of choice


# Testing / Performance
*Below I will try my best to update this brief overview of performance the engine develops at every level, seeing which techniques added the most elo boost and when*

