Changelog
---------
* Blunder 4.0.0
* Blunder 3.0.0
* Blunder 2.0.0
* Blunder 1.0.0 (Initial release)

Blunder 4.0.0
-------------
Blunder 4.0.0 includes "filtered" move generation. In other words, Blunder's move generator can now produce all moves, or only captures. The ability to producce only capture moves was added to speed up quiesence search, and this speed-up gave Blunder a ~35 Elo increase (in self play). Addtionally, I did extensive refractoring of the size of types used throughout Blunder's codebase, and shrinking types to only as big as they needed to be speed-up blunder and gained roughly another ~15 Elo (in self play); putting Blunder's total Elo gain between version 3.0.0 and 4.0.0 at ~50 Elo (in self play).

* Engine
    - Filtered move generation

(I don't include the type refractoring in the changelog here since I consider it more of a bug that needed to have been fixed since Blunder first version.)

Blunder 3.0.0
-------------
Blunder 3.0.0, now includes a tapered evaluation (and updated piece-square tables to better suit the update), killer move heuristics, and a transposition table, only for running perft. An upcoming feature will be a transposition table added to the search. These new features added to the search and evaluation phases of Blunder gave a collective increase of about ~207 Elo, putting at roughly ~1782 Elo in self-play.

* Engine
    - Transposition table for perft
* Search
    - Killer moves
* Evaluation
    - Tapered evaluation

Blunder 2.0.0
-------------
Blunder 2.0.0 adds three new features: Zobrist hashing, because of the hashing, three-fold repetition detection, and better piece-square table
values, courtesey of Marcel Vanthoor, author of [Rustic](https://github.com/mvanthoor/rustic). A future goal is to automatically generate piece-square table, and other evaluation
values via [Texel tuning](https://www.chessprogramming.org/Texel%27s_Tuning_Method). These features combined show an increase of ~175 Elo (+/- 38.4) in self-play testing against Blunder 1.0.0, bringing Blunder 2.0.0 to around ~1570 Elo in self play.

* Engine
    - [Zobrist hashing](https://www.chessprogramming.org/Zobrist_Hashing)
* Search & Evaluation
    - Three-fold repetiton detection

Blunder 1.0.0
-------------

* Engine
    - [Bitboards representation](https://www.chessprogramming.org/Bitboards)
    - [Magic bitboards for slider move generation](https://www.chessprogramming.org/Magic_Bitboards)
* Search
    - [Negamax search framework](https://www.chessprogramming.org/Negamax)
    - [Alpha-Beta pruning](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning)
    - [MVV-LVA move ordering](https://www.chessprogramming.org/MVV-LVA)
    - [Quiescence search](https://www.chessprogramming.org/Quiescence_Search)
    - Time-control logic supporting classical, rapid, bullet, and ultra-bullet time formats.
* Evaluation
    - [Material evaluation](https://www.chessprogramming.org/Material)
    - [Hand-written piece-square tables](https://www.chessprogramming.org/Piece-Square_Tables)
