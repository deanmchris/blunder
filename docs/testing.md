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

Re-tune Evaluation Using Gradient Descent Tuner
-----------------------------------------------
```
Score of Blunder 8.0.0-evaluation-tuning vs Blunder 8.0.0: 541 - 339 - 284  [0.587] 1164
...      Blunder 8.0.0-evaluation-tuning playing White: 268 - 173 - 141  [0.582] 582
...      Blunder 8.0.0-evaluation-tuning playing Black: 273 - 166 - 143  [0.592] 582
...      White vs Black: 434 - 446 - 284  [0.495] 1164
Elo difference: 60.9 +/- 17.5, LOS: 100.0 %, DrawRatio: 24.4 %
SPRT: llr 2.96 (100.4%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Re-tune Evaluation WIth Bug-Fixed Tuner
-----------------------------------------------
```
Score of Blunder 8.0.0-v1 vs Blunder 8.0.0-v2: 485 - 492 - 442  [0.498] 1419
...      Blunder 8.0.0-v1 playing White: 256 - 237 - 217  [0.513] 710
...      Blunder 8.0.0-v1 playing Black: 229 - 255 - 225  [0.482] 709
...      White vs Black: 511 - 466 - 442  [0.516] 1419
Elo difference: -1.7 +/- 15.0, LOS: 41.1 %, DrawRatio: 31.1 %
SPRT: llr -0.322 (-10.9%), lbound -2.94, ubound 2.94
Finished match
```

Add Evaluation Term for Having The Bishop Pair
----------------------------------------------
```
Score of Blunder 8.0.0-bishop-pair-eval vs Blunder 8.0.0: 774 - 678 - 548  [0.524] 2000
...      Blunder 8.0.0-bishop-pair-eval playing White: 395 - 327 - 278  [0.534] 1000
...      Blunder 8.0.0-bishop-pair-eval playing Black: 379 - 351 - 270  [0.514] 1000
...      White vs Black: 746 - 706 - 548  [0.510] 2000
Elo difference: 16.7 +/- 13.0, LOS: 99.4 %, DrawRatio: 27.4 %
SPRT: llr 1.25 (42.4%), lbound -2.94, ubound 2.94
Finished match
```

Aspiration Windows
------------------
```
Score of Blunder 8.0.0-aspiration-window vs Blunder 8.0.0: 662 - 532 - 806  [0.532] 2000
...      Blunder 8.0.0-aspiration-window playing White: 360 - 263 - 377  [0.548] 1000
...      Blunder 8.0.0-aspiration-window playing Black: 302 - 269 - 429  [0.516] 1000
...      White vs Black: 629 - 565 - 806  [0.516] 2000
Elo difference: 22.6 +/- 11.8, LOS: 100.0 %, DrawRatio: 40.3 %
SPRT: llr 2.18 (74.1%), lbound -2.94, ubound 2.94
Finished match
```

Dynamic Time Management Via Aspiration Windows
----------------------------------------------
```
Score of Blunder 8.0.0-dynamic-tc vs Blunder 8.0.0: 638 - 539 - 823  [0.525] 2000
...      Blunder 8.0.0-dynamic-tc playing White: 352 - 259 - 389  [0.546] 1000
...      Blunder 8.0.0-dynamic-tc playing Black: 286 - 280 - 434  [0.503] 1000
...      White vs Black: 632 - 545 - 823  [0.522] 2000
Elo difference: 17.2 +/- 11.7, LOS: 99.8 %, DrawRatio: 41.1 %
SPRT: llr 1.6 (54.4%), lbound -2.94, ubound 2.94
Finished match
```

Utilize The Fail-Soft Score in Static Null-Move Pruning
-------------------------------------------------------
```
Score of Blunder 8.0.0-rnmp-tuning vs Blunder 8.0.0: 649 - 515 - 836  [0.533] 2000
...      Blunder 8.0.0-rnmp-tuning playing White: 325 - 244 - 431  [0.540] 1000
...      Blunder 8.0.0-rnmp-tuning playing Black: 324 - 271 - 405  [0.526] 1000
...      White vs Black: 596 - 568 - 836  [0.507] 2000
Elo difference: 23.3 +/- 11.6, LOS: 100.0 %, DrawRatio: 41.8 %
SPRT: llr 2.32 (78.7%), lbound -2.94, ubound 2.94
Finished match
```

Add Basic Mobility Evaluation
-----------------------------
```
Post-Convergence rating estimation
done

   # PLAYER                          : RATING    POINTS  PLAYED    (%)
   1 Blunder 8.0.0-mobility          : 2513.8    1078.5    2000   53.9%
   2 Blunder 8.0.0                   : 2486.2     921.5    2000   46.1%

White advantage = 0.00
Draw rate (equal opponents) = 50.00 %
```

Add Internal Iterative Deepening
--------------------------------
```
Score of Blunder 8.0.0-IID vs Blunder 8.0.0: 625 - 562 - 813  [0.516] 2000
...      Blunder 8.0.0-IID playing White: 330 - 270 - 400  [0.530] 1000
...      Blunder 8.0.0-IID playing Black: 295 - 292 - 413  [0.501] 1000
...      White vs Black: 622 - 565 - 813  [0.514] 2000
Elo difference: 10.9 +/- 11.7, LOS: 96.6 %, DrawRatio: 40.6 %
SPRT: llr 0.888 (30.2%), lbound -2.94, ubound 2.94
Finished match
```

Add Basic King Safety
---------------------
```
Score of Blunder 8.0.0-king-safety vs Blunder 8.0.0: 578 - 391 - 435  [0.567] 1404
...      Blunder 8.0.0-king-safety playing White: 312 - 179 - 211  [0.595] 702
...      Blunder 8.0.0-king-safety playing Black: 266 - 212 - 224  [0.538] 702
...      White vs Black: 524 - 445 - 435  [0.528] 1404
Elo difference: 46.6 +/- 15.2, LOS: 100.0 %, DrawRatio: 31.0 %
SPRT: llr 2.94 (100.0%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Add Basic Passed Pawn Evaluation
--------------------------------
```
Score of Blunder 8.0.0-passed-pawns vs Blunder 8.0.0: 608 - 435 - 627  [0.552] 1670
...      Blunder 8.0.0-passed-pawns playing White: 333 - 191 - 310  [0.585] 834
...      Blunder 8.0.0-passed-pawns playing Black: 275 - 244 - 317  [0.519] 836
...      White vs Black: 577 - 466 - 627  [0.533] 1670
Elo difference: 36.1 +/- 13.2, LOS: 100.0 %, DrawRatio: 37.5 %
SPRT: llr 2.95 (100.1%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```