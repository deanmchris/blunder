Changelog
---------
* Blunder 7.x.x
* Blunder 6.1.0
* Blunder 6.0.0
* Blunder 5.0.0
* Blunder 4.0.0
* Blunder 3.0.0
* Blunder 2.0.0
* Blunder 1.0.0 (Initial release)

Blunder 7.x.x
-------------
* Engine
    - Decreased the memory footprint of Position objects by removing the array holding position history into a global variable.
    - XBoard compatibility
* Search
    - TODO: Fix the node count. You double count nodes before going into QSearch!
    - Re-add HH since fixed history heuristics bug
    - Fixed PV bug
    - Fixed node counting bug
    - Show mate score as mate-in-N
* Evaluation
    - Texel tuner

Blunder 6.1.0
-------------

Blunder 6.1.0 includes several bugfixes throughout the code, and removes history heuristics. Even with history heuristics removed however, the Elo gained from bugfixes puts Blunder 6.1.0 roughly equal to Blunder 6.0.0's strength. After further testing and bug discovery, Blunder's current Elo is now estimated to be around 2200.

* Engine
    - Create rudimentary endgame detection.
* Search
    - Adjust the scores in the transposition table before checking for a hit.
    - Fix the UCI "Hash" command to actually resize the hash table with the given size.
    - Rework the contempt factor used to be more nuanced and accurate.

Blunder 6.0.0
-------------

Blunder 6.0.0 includes an implementation of reverse futility pruning and history heuristics, both which from self-play contributed 80-100 Elo. Additionally, the UCI stop command has been implemented, as well as some general code cleanups. Blunder 6.0.0's estimated rating is ~2100-2150 Elo.

* Engine
    - UCI "stop" command
* Search
    - [Reverse futility pruning](https://www.chessprogramming.org/Reverse_Futility_Pruning)
    - [History heuristics](https://www.chessprogramming.org/History_Heuristic)

Blunder 5.0.0
-------------

Blunder 5.0.0 is a complete rewrite of the engine. Many basic design ideas and principles were kept, and some areas were just ported over, but the majority of the code-base was rewritten, and the layout of the project was completely changed.

The movitations behind this rewrite were twofold: First, I was dissatisfied with Blunder's speed and wanted to take another crack at creating an engine that was simply faster. Second, I didn't like how Blunder was designed in several places, and I quickly realized these "several places" constitued large chunks of the code-base. So Blunder 5.0.0 was the result.

From the testing I've done, Blunder 5.0.0 is 20-30% faster than Blunder 4.0.0, and perft(6) from the starting position was coming in at around 6-8s (14-18 Mnps), whereas perft(6) from the starting position for Blunder 4.0.0 was generally 10-12s (10-12Mnps). And overall, I'm happy with the refractoring I've done and my code feels cleaner in many of the places that bothered me. So both goals, all things considered, were meet.

Although Blunder 5.0.0 is a rewrite, it did build on Blunder 4.0.0, and two new features were added: a transposition table, and null-move pruning. Additional, the tapered evaluation has been refractored and is stronger. And the speed increase should add some Elo to engine, though I didn't test for a specfic amount. Overall, these changes have added about 200-300 Elo to the engine in self-play, and puts Blunder at 2038 Elo on the CCRL.

Since Blunder 5.0.0 is a rewrite, a listing of all of the current features are listed below:

* Engine
    - [Bitboards representation](https://www.chessprogramming.org/Bitboards)
    - [Magic bitboards for slider move generation](https://www.chessprogramming.org/Magic_Bitboards)
    - [Zobrist hashing](https://www.chessprogramming.org/Zobrist_Hashing)
* Search
    - [Negamax search framework](https://www.chessprogramming.org/Negamax)
    - [Alpha-Beta pruning](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning)
    - [MVV-LVA move ordering](https://www.chessprogramming.org/MVV-LVA)
    - [Quiescence search](https://www.chessprogramming.org/Quiescence_Search)
    - [Time-control logic supporting classical, rapid, bullet, and ultra-bullet time formats](https://www.chessprogramming.org/Time_Management)
    - [Repition detection](https://www.chessprogramming.org/Repetitions)
    - [Killer moves](https://www.chessprogramming.org/Killer_Move)
    - [Transposition table](https://www.chessprogramming.org/Transposition_Table)
    - [Null-move pruning](https://www.chessprogramming.org/Null_Move_Pruning)
* Evaluation
    - [Material evaluation](https://www.chessprogramming.org/Material)
    - [Tuned piece-square tables](https://www.chessprogramming.org/Piece-Square_Tables)
    - [Tapered evaluation](https://www.chessprogramming.org/Tapered_Eval)

Blunder 4.0.0
-------------
Blunder 4.0.0 includes "filtered" move generation. In other words, Blunder's move generator can now produce all moves or only captures. The ability to produce only capture moves were added to speed up quiescence search, and this speed-up gave Blunder a ~35 Elo increase (in self-play). Additionally, I did extensive refactoring of the size of types used throughout Blunder's codebase, and shrinking types to only as big as they needed to be speedup blunder and gained roughly another ~15 Elo (in self-play); putting Blunder's total Elo gain between version 3.0.0 and 4.0.0 at ~50 Elo (in self-play).

One more thing of note is that Blunder 4.0.0 is the first version to include releases for Windows and macOS in addition to Linux. Going forward, all future releases of Blunder will include releases targeting all three operating systems. If any of the three releases still don't work for you, see the README on how to build your own fairly easily. With these new releases, a bug had to be fixed with the way Blunder handled IO operations across platforms.

Engine
* Filtered move generation
* Refractor types to more conservative sizes
* Update IO operations to be cross-platform compatible

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
