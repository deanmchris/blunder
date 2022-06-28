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

Add Basic Late-Move Reduction Table
-----------------------------------
```
Rank Name                          Elo     +/-   Games   Score    Draw
   0 Blunder 8.0.0                  -7      15    1166   49.0%   42.5%
   1 Blunder 8.0.0-d6-m12           23      27     388   53.4%   39.2%
   2 Blunder 8.0.0-d3-m14            6      26     389   50.9%   44.2%
   3 Blunder 8.0.0-d5-m14           -8      26     389   48.8%   44.2%
```

Here the rather cryptic ```Blunder 8.0.0-d6-m12``` is the variation of LMR tweaking I
tried that worked the best.

Tweak the Futility Pruning Margins
----------------------------------
```
Score of Blunder 8.0.0-fp-tuning vs Blunder 8.0.0: 580 - 518 - 902  [0.515] 2000
...      Blunder 8.0.0-fp-tuning playing White: 299 - 254 - 447  [0.522] 1000
...      Blunder 8.0.0-fp-tuning playing Black: 281 - 264 - 455  [0.508] 1000
...      White vs Black: 563 - 535 - 902  [0.507] 2000
Elo difference: 10.8 +/- 11.3, LOS: 96.9 %, DrawRatio: 45.1 %
SPRT: llr 0.939 (31.9%), lbound -2.94, ubound 2.94
Finished match
```

Add Doubled and Isolated Pawn Structure Evaluation Terms
--------------------------------------------------------
```
Score of Blunder 8.0.0-pawn-structure vs Blunder 8.0.0: 730 - 545 - 655  [0.548] 1930
...      Blunder 8.0.0-pawn-structure playing White: 368 - 265 - 332  [0.553] 965
...      Blunder 8.0.0-pawn-structure playing Black: 362 - 280 - 323  [0.542] 965
...      White vs Black: 648 - 627 - 655  [0.505] 1930
Elo difference: 33.4 +/- 12.6, LOS: 100.0 %, DrawRatio: 33.9 %
SPRT: llr 2.95 (100.2%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Add Late-Move Pruning To Search
-------------------------------
```
Score of Blunder 8.0.0-lmp vs Blunder 8.0.0: 629 - 503 - 868  [0.531] 2000
...      Blunder 8.0.0-lmp playing White: 322 - 239 - 439  [0.541] 1000
...      Blunder 8.0.0-lmp playing Black: 307 - 264 - 429  [0.521] 1000
...      White vs Black: 586 - 546 - 868  [0.510] 2000
Elo difference: 21.9 +/- 11.4, LOS: 100.0 %, DrawRatio: 43.4 %
SPRT: llr 2.22 (75.4%), lbound -2.94, ubound 2.94
Finished match
```

Tweak Static-Null Move Pruning Margins
--------------------------------------
```
Rank Name                          Elo     +/-   Games   Score    Draw
   0 Blunder 8.0.0-snmp-120        -11      10    2519   48.4%   44.1%
   1 Blunder 8.0.0-snmp-85          28      20     630   54.0%   45.6%
   2 Blunder 8.0.0-snmp-110          9      21     629   51.4%   42.3%
   3 Blunder 8.0.0-snmp-90           6      20     630   50.8%   45.1%
   4 Blunder 8.0.0-snmp-95           2      20     630   50.3%   43.5%
```

Generate Queen Promotions in Quiescence Search
----------------------------------------------

After ~2K games a very slight to no Elo gain was noticed, although
at several instances it looked promising. The engine seems to play
bit more tactically sound however, so the change will be kept.

Check For Semi-Open Files Next To King
--------------------------------------
```
Score of Blunder 8.0.0-pawn-shields vs Blunder 8.0.0: 586 - 578 - 886  [0.502] 2050
...      Blunder 8.0.0-pawn-shields playing White: 313 - 271 - 440  [0.521] 1024
...      Blunder 8.0.0-pawn-shields playing Black: 273 - 307 - 446  [0.483] 1026
...      White vs Black: 620 - 544 - 886  [0.519] 2050
Elo difference: 1.4 +/- 11.3, LOS: 59.3 %, DrawRatio: 43.2 %
SPRT: llr -0.197 (-6.7%), lbound -2.94, ubound 2.94
```

Keep it for now.

Reward Rooks And Queens on Seventh Rank
---------------------------------------
```
Score of Blunder 8.0.0-rook-or-queen-on-seventh vs Blunder 8.0.0: 1152 - 1046 - 1802  [0.513] 4000
...      Blunder 8.0.0-rook-or-queen-on-seventh playing White: 604 - 487 - 909  [0.529] 2000
...      Blunder 8.0.0-rook-or-queen-on-seventh playing Black: 548 - 559 - 893  [0.497] 2000
...      White vs Black: 1163 - 1035 - 1802  [0.516] 4000
Elo difference: 9.2 +/- 8.0, LOS: 98.8 %, DrawRatio: 45.1 %
SPRT: llr 1.5 (50.9%), lbound -2.94, ubound 2.94
Finished match
```

Add A Small Counter-Move Bonus
------------------------------
```
Score of Blunder 8.0.0-nmp-tweaking vs Blunder 8.0.0: 722 - 685 - 1093  [0.507] 2500
...      Blunder 8.0.0-nmp-tweaking playing White: 398 - 316 - 536  [0.533] 1250
...      Blunder 8.0.0-nmp-tweaking playing Black: 324 - 369 - 557  [0.482] 1250
...      White vs Black: 767 - 640 - 1093  [0.525] 2500
Elo difference: 5.1 +/- 10.2, LOS: 83.8 %, DrawRatio: 43.7 %
SPRT: llr 0.315 (10.7%), lbound -2.94, ubound 2.94
```

Add Knight Outpost
------------------
```
Score of Blunder 8.0.0-kight-outposts vs Blunder 8.0.0: 1200 - 1141 - 1659  [0.507] 4000
...      Blunder 8.0.0-kight-outposts playing White: 646 - 534 - 820  [0.528] 2000
...      Blunder 8.0.0-kight-outposts playing Black: 554 - 607 - 839  [0.487] 2000
...      White vs Black: 1253 - 1088 - 1659  [0.521] 4000
Elo difference: 5.1 +/- 8.2, LOS: 88.9 %, DrawRatio: 41.5 %
SPRT: llr 0.481 (16.3%), lbound -2.94, ubound 2.94
Finished match
```

Factor Enemy Pawns Into Knight Mobility Evaluation
--------------------------------------------------
```
Score of Blunder 8.0.0-safe-mobility vs Blunder 8.0.0: 929 - 838 - 1233  [0.515] 3000
...      Blunder 8.0.0-safe-mobility playing White: 490 - 381 - 629  [0.536] 1500
...      Blunder 8.0.0-safe-mobility playing Black: 439 - 457 - 604  [0.494] 1500
...      White vs Black: 947 - 820 - 1233  [0.521] 3000
Elo difference: 10.5 +/- 9.5, LOS: 98.5 %, DrawRatio: 41.1 %
SPRT: llr 1.27 (43.3%), lbound -2.94, ubound 2.94
```

Increase Scale Factor To Divide Eval By In Drawish Endgames
-----------------------------------------------------------
```
Score of Blunder 8.0.0-draw-knowledge vs Blunder 8.0.0: 940 - 838 - 1372  [0.516] 3150
...      Blunder 8.0.0-draw-knowledge playing White: 497 - 389 - 690  [0.534] 1576
...      Blunder 8.0.0-draw-knowledge playing Black: 443 - 449 - 682  [0.498] 1574
...      White vs Black: 946 - 832 - 1372  [0.518] 3150
Elo difference: 11.3 +/- 9.1, LOS: 99.2 %, DrawRatio: 43.6 %
SPRT: llr 1.53 (51.9%), lbound -2.94, ubound 2.94
```

Progress Check of Dev Vs Latest Release
---------------------------------------
```
Score of Blunder 8.0.0 vs Blunder 7.6.0: 734 - 550 - 676  [0.547] 1960
...      Blunder 8.0.0 playing White: 373 - 265 - 343  [0.555] 981
...      Blunder 8.0.0 playing Black: 361 - 285 - 333  [0.539] 979
...      White vs Black: 658 - 626 - 676  [0.508] 1960
Elo difference: 32.7 +/- 12.5, LOS: 100.0 %, DrawRatio: 34.5 %
SPRT: llr 2.95 (100.3%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

Add A Multi-Bucket Scheme To The Transposition Table
----------------------------------------------------

The Elo gain looked promising after several hundred games:

```
Score of Blunder 8.0.0-tt-tweaking vs Blunder 8.0.0: 186 - 142 - 322  [0.534] 650
...      Blunder 8.0.0-tt-tweaking playing White: 92 - 69 - 164  [0.535] 325
...      Blunder 8.0.0-tt-tweaking playing Black: 94 - 73 - 158  [0.532] 325
...      White vs Black: 165 - 163 - 322  [0.502] 650
Elo difference: 23.6 +/- 19.0, LOS: 99.2 %, DrawRatio: 49.5 %
SPRT: llr 0.881 (29.9%), lbound -2.94, ubound 2.94
```

But evaporated after a couple hundred more:

```
Score of Blunder 8.0.0-tt-tweaking vs Blunder 8.0.0: 305 - 302 - 597  [0.501] 1204
...      Blunder 8.0.0-tt-tweaking playing White: 150 - 146 - 307  [0.503] 603
...      Blunder 8.0.0-tt-tweaking playing Black: 155 - 156 - 290  [0.499] 601
...      White vs Black: 306 - 301 - 597  [0.502] 1204
Elo difference: 0.9 +/- 13.9, LOS: 54.8 %, DrawRatio: 49.6 %
SPRT: llr -0.169 (-5.7%), lbound -2.94, ubound 2.94
Finished match
```

From some synthetic testing as well as the testing above, the feature looks promising,
but the time control of 10+0.1 may be too short to see any significant difference in
strength, and longer time control testing isn't particularly feasible at the moment,
so the change will be accepted for now, assuming it might contribute something like
~10-15 Elo at bullet-time controls.

Generate All Moves If In Check And <= 3 Plies Under The Horizon
---------------------------------------------------------------
```
Score of Blunder 8.0.0-qsearch-tweaking vs Blunder 8.0.0: 1077 - 1038 - 1885  [0.505] 4000
...      Blunder 8.0.0-qsearch-tweaking playing White: 596 - 466 - 938  [0.532] 2000
...      Blunder 8.0.0-qsearch-tweaking playing Black: 481 - 572 - 947  [0.477] 2000
...      White vs Black: 1168 - 947 - 1885  [0.528] 4000
Elo difference: 3.4 +/- 7.8, LOS: 80.2 %, DrawRatio: 47.1 %
SPRT: llr 0.0972 (3.3%), lbound -2.94, ubound 2.94
Finished match
```

Tweak LMR Reductions To Be Slightly More Agressive
--------------------------------------------------
```
Score of Blunder 8.0.0-lmr-tweaking vs Blunder 8.0.0: 1115 - 1081 - 1804  [0.504] 4000
...      Blunder 8.0.0-lmr-tweaking playing White: 586 - 524 - 890  [0.515] 2000
...      Blunder 8.0.0-lmr-tweaking playing Black: 529 - 557 - 914  [0.493] 2000
...      White vs Black: 1143 - 1053 - 1804  [0.511] 4000
Elo difference: 3.0 +/- 8.0, LOS: 76.6 %, DrawRatio: 45.1 %
SPRT: llr -0.0113 (-0.4%), lbound -2.94, ubound 2.94
Finished match
```

Keep for now.

Add Basic Razoring Scheme
-------------------------
```
Score of Blunder 8.0.0-razoring vs Blunder 8.0.0: 1291 - 1187 - 2072  [0.511] 4550
...      Blunder 8.0.0-razoring playing White: 700 - 544 - 1031  [0.534] 2275
...      Blunder 8.0.0-razoring playing Black: 591 - 643 - 1041  [0.489] 2275
...      White vs Black: 1343 - 1135 - 2072  [0.523] 4550
Elo difference: 7.9 +/- 7.4, LOS: 98.2 %, DrawRatio: 45.5 %
SPRT: llr 1.37 (46.5%), lbound -2.94, ubound 2.94
```

Progress Check Of Dev Vs Latest Release
---------------------------------------
```
Score of Blunder 8.0.0 vs Blunder 7.6.0: 527 - 350 - 469  [0.566] 1346
...      Blunder 8.0.0 playing White: 276 - 163 - 235  [0.584] 674
...      Blunder 8.0.0 playing Black: 251 - 187 - 234  [0.548] 672
...      White vs Black: 463 - 414 - 469  [0.518] 1346
Elo difference: 46.0 +/- 15.0, LOS: 100.0 %, DrawRatio: 34.8 %
SPRT: llr 2.95 (100.4%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```
