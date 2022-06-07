Testing
-------

This file is for precisley documenting each feature that is added to Blunder, after certian base features, and
specfic, test-based evidence that the feature improved the strength of Blunder through
self-play.

Hopefully this file will be useful for anyone looking to implement certain features, but
truthfully, it's mostly meant to keep myself honest, and force me not to cut corners when
implementing and testing new features.

Many of the features that Blunder currently has will either be re-tested and appear in this file, or be removed until further testing can be done. The end goal being creating an engine that has feature parity with v7.6.0 or fewer features, and still be stronger.

Starting Basis
--------------

The following list of features is where the testing in this file begins from:

* Engine
    - Bitboards representation
    - Magic bitboards for slider move generation
    - Zobrist hashing
* Search
    - Alpha-Beta pruning (fail soft)
    - Basic time-control logic
    - MVV-LVA move ordering
    - Quiescence search
    - Repition detection
    - Killer moves
    - Transposition table
    - Principal Variation Search
    - Check extension
* Evaluation
    - Material evaluation
    - Tuned piece-square tables
    - Tapered evaluation

Each subsequent feature re-added to Blunder will be listed below, along with the self-play test results that demonstrate it's approximate gain.

Null-move Pruning
-----------------
```
Score of Blunder 8.0.0-nm vs Blunder 8.0.0: 320 - 129 - 144  [0.661] 593
...      Blunder 8.0.0-nm playing White: 159 - 65 - 73  [0.658] 297
...      Blunder 8.0.0-nm playing Black: 161 - 64 - 71  [0.664] 296
...      White vs Black: 223 - 226 - 144  [0.497] 593
Elo difference: 116.0 +/- 25.2, LOS: 100.0 %, DrawRatio: 24.3 %
SPRT: llr 2.95 (100.1%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Static Null-Move Pruning
------------------------
```
Score of Blunder 8.0.0 vs Blunder 8.0.0: 544 - 347 - 318  [0.581] 1209
...      Blunder 8.0.0 playing White: 272 - 180 - 153  [0.576] 605
...      Blunder 8.0.0 playing Black: 272 - 167 - 165  [0.587] 604
...      White vs Black: 439 - 452 - 318  [0.495] 1209
Elo difference: 57.1 +/- 16.9, LOS: 100.0 %, DrawRatio: 26.3 %
SPRT: llr 2.95 (100.1%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

History Heuristics
------------------
```
Score of Blunder 8.0.0-hh vs Blunder 8.0.0: 767 - 658 - 575  [0.527] 2000
...      Blunder 8.0.0-hh playing White: 398 - 332 - 271  [0.533] 1001
...      Blunder 8.0.0-hh playing Black: 369 - 326 - 304  [0.522] 999
...      White vs Black: 724 - 701 - 575  [0.506] 2000
Elo difference: 19.0 +/- 12.9, LOS: 99.8 %, DrawRatio: 28.7 %
SPRT: llr 1.48 (50.4%), lbound -2.94, ubound 2.94
```

Dynamic Time Management For Fixed Time
--------------------------------------
```
Score of Blunder 8.0.0-dtm vs Blunder 8.0.0: 626 - 497 - 477  [0.540] 1600
...      Blunder 8.0.0-dtm playing White: 307 - 255 - 239  [0.532] 801
...      Blunder 8.0.0-dtm playing Black: 319 - 242 - 238  [0.548] 799
...      White vs Black: 549 - 574 - 477  [0.492] 1600
Elo difference: 28.1 +/- 14.3, LOS: 100.0 %, DrawRatio: 29.8 %
SPRT: llr 1.9 (64.4%), lbound -2.94, ubound 2.94
```

Note that this feature only gains elo for games with a fixed amount of time per-player (plus any 
increments).

Futility Pruning
----------------
```
Score of Blunder 8.0.0-fp vs Blunder 8.0.0: 708 - 517 - 555  [0.554] 1780
...      Blunder 8.0.0-fp playing White: 358 - 262 - 270  [0.554] 890
...      Blunder 8.0.0-fp playing Black: 350 - 255 - 285  [0.553] 890
...      White vs Black: 613 - 612 - 555  [0.500] 1780
Elo difference: 37.4 +/- 13.4, LOS: 100.0 %, DrawRatio: 31.2 %
SPRT: llr 2.96 (100.4%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Skip Moves in QSearch via Static Exchange Evaluation
----------------------------------------------------
```
Score of Blunder 8.0.0-see-pruning vs Blunder 8.0.0: 775 - 626 - 599  [0.537] 2000
...      Blunder 8.0.0-see-pruning playing White: 394 - 294 - 312  [0.550] 1000
...      Blunder 8.0.0-see-pruning playing Black: 381 - 332 - 287  [0.524] 1000
...      White vs Black: 726 - 675 - 599  [0.513] 2000
Elo difference: 25.9 +/- 12.8, LOS: 100.0 %, DrawRatio: 29.9 %
SPRT: llr 2.17 (73.7%), lbound -2.94, ubound 2.94
Finished match
```

Late Move Reductions
--------------------
```
Score of Blunder 8.0.0-lmr vs Blunder 8.0.0: 444 - 263 - 330  [0.587] 1037
...      Blunder 8.0.0-lmr playing White: 234 - 121 - 164  [0.609] 519
...      Blunder 8.0.0-lmr playing Black: 210 - 142 - 166  [0.566] 518
...      White vs Black: 376 - 331 - 330  [0.522] 1037
Elo difference: 61.3 +/- 17.6, LOS: 100.0 %, DrawRatio: 31.8 %
SPRT: llr 2.96 (100.4%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Use R = 3 + depth/6 Formula to Calculate Null Move Reductions
-------------------------------------------------------------
```
Score of Blunder 8.0.0-advanced-nmp vs Blunder 8.0.0: 632 - 552 - 816  [0.520] 2000
...      Blunder 8.0.0-advanced-nmp playing White: 335 - 256 - 409  [0.539] 1000
...      Blunder 8.0.0-advanced-nmp playing Black: 297 - 296 - 407  [0.500] 1000
...      White vs Black: 631 - 553 - 816  [0.519] 2000
Elo difference: 13.9 +/- 11.7, LOS: 99.0 %, DrawRatio: 40.8 %
SPRT: llr 1.22 (41.5%), lbound -2.94, ubound 2.94
Finished match
```

Make Futility Pruning Margins More Agressive
--------------------------------------------
```
Score of Blunder 8.0.0-advanced-fp vs Blunder 8.0.0: 652 - 539 - 809  [0.528] 2000
...      Blunder 8.0.0-advanced-fp playing White: 336 - 254 - 410  [0.541] 1000
...      Blunder 8.0.0-advanced-fp playing Black: 316 - 285 - 399  [0.515] 1000
...      White vs Black: 621 - 570 - 809  [0.513] 2000
Elo difference: 19.7 +/- 11.7, LOS: 99.9 %, DrawRatio: 40.5 %
SPRT: llr 1.86 (63.0%), lbound -2.94, ubound 2.94
Finished match
```

Update Material, Piece-Square Table, and Phase Incrementally
------------------------------------------------------------
```
Score of Blunder 8.0.0 vs Blunder 8.0.0-advanced-fp: 352 - 199 - 401  [0.580] 952
...      Blunder 8.0.0 playing White: 178 - 93 - 206  [0.589] 477
...      Blunder 8.0.0 playing Black: 174 - 106 - 195  [0.572] 475
...      White vs Black: 284 - 267 - 401  [0.509] 952
Elo difference: 56.3 +/- 16.8, LOS: 100.0 %, DrawRatio: 42.1 %
SPRT: llr 2.95 (100.2%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Estimated ELO: ~2450

- Late-move pruning/move-count based pruning
- Aspiration windows
- Contempt
- Incrementally update PSQT

Unexplained crash vs madchess.