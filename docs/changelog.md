Changelog
---------
* Blunder 8.5.5
* Blunder 8.0.0
* Blunder 7.6.0
* Blunder 7.5.0
* Blunder 7.4.0
* Blunder 7.3.0
* Blunder 7.2.0
* Blunder 7.1.0
* Blunder 7.0.0
* Blunder 6.1.0
* Blunder 6.0.0
* Blunder 5.0.0
* Blunder 4.0.0
* Blunder 3.0.0
* Blunder 2.0.0
* Blunder 1.0.0 (Initial release)

Blunder 8.5.5
-------------

Blunder 8.5.5 includes several new features, as well as some strength gaining tweaks to the engine.
Singular extensions have been added, as well as bishop outposts. The time management has been tweaked
to be a bit stronger, as well as late-move reductions. Bugs where the fifty-move rule wasn't reset for
a capture and checkmate didn't trump the fifty-move rule have been corrected, as well as a bug in the
transposition table replacement scheme.

This version of Blunder should be ~40-50 Elo stronger than 8.0.0 at bullet time controls.

Blunder 8.0.0
-------------

This release is very similar in spirit to that of 5.0.0, in the sense that much of the work done has been improvement and optimization of the code base through slow, large-scale refactoring, which have
brought about Elo gains, although several new strength gainging features have been implemented. It'd
be impossible to list every last change, but a summary of some of the major areas will be listed below.

 - Improvements to Blunder's move generator have been made, improving it's speed by 10-15%. Thanks in particular to Koivisto, for giving me the idea to store the hash of a position
as part of it's irreversible state, leading to a nice speed increase.

- Small tweaks have been made to the UCI that's made it a bit more user friendly and hopefully
more straightforward to use.

- Magic bitboard generation code and comments have also been improved and cleaned-up.

- Perhaps the biggest and move effective change, the tuner for Blunder has been updated from using a
naieve local optimization algorithm, to gradient descent, specfically AdaGrad. The speed and efficeny of
the new tuner in allowing me to quickly add, tune, and reject evaluation ideas and improvements allowed
me to very quickly tune a stronger, and more streamlined evaluation, than the one I had in Blunder 7.6.0.
I'll be uploading a paper shortly to this repository mostly outlining the math behind the process of 
getting the tuner working, although parts of the code will be touched upon. The hope is that such a
document will be useful to those who were in the position I was in several months ago.

- The time manager code has had quite a large overhaul, making the code cleaner, and the API easier
to use. Time managament in games with no increment or bonus time (a.k.a "sudden death"), has also
been improved, using the assumption that it's better to spend less time in the opening phase,
and more towards the end of the game.

- For now [WAC testing](https://www.chessprogramming.org/Win_at_Chess) (a.k.a "iq testing") has been
removed from Blunder for now, as well code in `blunder/tuner` to generate tuning data for Blunder.
Both hadn't been very useful, although I plan on experimenting with generating custom tuning data
back into Blunder very soon, so that code, at least, will be added back into later versions.

- A logo has been created by myself for Blunder. It's nothing too complicated, as I'm not an artist,
but it's an orginaly design I had a while ago, and I think it suits Blunder well. See the README,
or `blunder/logo` for the logo.

- Much of the code has been updated to make benefit of new features added in Go 1.18,
most notably generics, which have allowed me to clean-up much of the code-base, including
the transposition table.

As far as some of the new strength gaining features, they're listed below:

* Engine
    - Update PSQT and material incrementally
* Search
    - [Internal Iterative Deepening](https://www.chessprogramming.org/Internal_Iterative_Deepening)
    - [Razoring](https://www.chessprogramming.org/Razoring)
    - Consider checks for checks in quiescence search.
    - Add buckets to transposition table, and use a depth-replace and always-replace scheme
* Evaluation
    - [Basic rook structure](https://www.chessprogramming.org/Evaluation_of_Pieces#Rook)
    - [Bishop pair](https://www.chessprogramming.org/Bishop_Pair)
    - [Make mobility for knight's safe with regards to pawns](https://www.chessprogramming.org/Mobility#Safe_Mobility)
    - [Drawn and drawish endgame recognition improvement](https://www.chessprogramming.org/Draw_Evaluation)

Lastly, the release of Go 1.18 has also come with the ability to target certain, extended
AMD 64 instruction sets. The makefile has been updated to allow compiling builds which
target some of these instruction sets, although I've not had the opportunity to extensively
test if there are any speed gains. I've also included pre-builds, for windows, macOS, and
linux in the version release.

See the README for more information.

Blunder 7.6.0
-------------

This release contains no new features, but various tweaks of existing ones that have improved
Blunder's strength. History values are decreased when the move fails to raise alpha or cause
a beta-cutoff, queen promotions are considered in quiescence search, pawn pushes to the 6th
rank and above are no longer reduced, and late-move reductions are no longer allowed to drop 
the search directly into quiescence search. These changes are summarized below and resulted 
in a gain of around 80 Elo in self-play:

```
Score of Blunder 7.6.0 vs Blunder 7.5.0: 313 - 147 - 251  [0.617] 711
...      Blunder 7.6.0 playing White: 156 - 80 - 120  [0.607] 356
...      Blunder 7.6.0 playing Black: 157 - 67 - 131  [0.627] 355
...      White vs Black: 223 - 237 - 251  [0.490] 711
Elo difference: 82.6 +/- 20.8, LOS: 100.0 %, DrawRatio: 35.3 %
SPRT: llr 2.95 (100.3%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

The individual gains from each tweak can be seen in the commit history, per usual.

* Search
    - Decrease the history value of moves that don't raise alpha or cause a betacutoff
    - Consider queen promotions in quiescence search
    - Don't reduce pawn pushes to the 6th rank and above
    - Don't allow LMR to drop the search directly into quiescence search

Blunder 7.5.0
-------------
The release of Blunder 7.5.0 consists of mostly strength gaining tweaks to Blunder's previous features, with the exception of the addition of late-move pruning/move-count based pruning. This version is estimated to be roughly 30 Elo stronger than Blunder 7.4.0.

* Search
    - [Late-move pruning/move-count based pruning](https://www.chessprogramming.org/Futility_Pruning#MoveCountBasedPruning) 

Blunder 7.4.0
-------------
The release of this version includes evaluation terms for passed pawns and knight outposts, a new form of dynamic time management, and a better tuned evaluation, particularly with regards to king safety. All of these additions amount to about 50 Elo in gauntlet testing over 7.3.0, and self-play tests show Blunder 7.4.0 is roughly 50 Elo stronger than 7.3.0 (tc=inf/10+0.1):

```
Score of Blunder 7.4.0-1 vs Blunder 7.3.0: 582 - 384 - 350  [0.575] 1316
...      Blunder 7.4.0-1 playing White: 286 - 191 - 181  [0.572] 658
...      Blunder 7.4.0-1 playing Black: 296 - 193 - 169  [0.578] 658
...      White vs Black: 479 - 487 - 350  [0.497] 1316
Elo difference: 52.7 +/- 16.2, LOS: 100.0 %, DrawRatio: 26.6 %
SPRT: llr 2.95 (100.4%), lbound -2.94, ubound 2.94 - H1 was accepted
```

and more than 100 Elo stronger than 7.1.0 (tc=inf/10+0.1):

```
Score of Blunder 7.4.0 vs Blunder 7.1.0: 334 - 132 - 117  [0.673] 583
...      Blunder 7.4.0 playing White: 173 - 63 - 55  [0.689] 291
...      Blunder 7.4.0 playing Black: 161 - 69 - 62  [0.658] 292
...      White vs Black: 242 - 224 - 117  [0.515] 583
Elo difference: 125.6 +/- 26.4, LOS: 100.0 %, DrawRatio: 20.1 %
SPRT: llr 2.95 (100.0%), lbound -2.94, ubound 2.94 - H1 was accepted
```

These testing results would put Blunder 7.4.0 at 2500+ Elo.

Many thanks to those who continue to test and use Blunder. Your support is very much appreicated.

Blunder 7.3.0
-------------

New features included in this version are evaluation terms for king safety and pawn structure. Both are still quite basic and have much lacking, but I believe them to be solid designs to build on in future versions. Addtionally, all evaluation terms have been retuned on a larger dataset, yielding some minior gains. However, most of the strength gain in this new version comes from slight tweaks to the pruning parameters and conditions for late-move reductions, null-move pruning, and futility pruning. Lastly, there are some non-strength gaining tweaks here and there, code clean-ups, and the command line interface is much more user friendly. The major changes as usual are summarized below:

* Engine
    - Improved command-line interface
    - UCI command `movetime` now implemented
* Search
    - Tweaked futility pruning's allowed depth and margins
    - Tweaked the formula for calculating null-move reductions (much more agressive for higher depths)
    - Trasitioned to using a basic table for computing late-move reductions.
* Evaluation
    - [Basic king safety](https://www.chessprogramming.org/King_Safety)
    - [Basic pawn structure](https://www.chessprogramming.org/Pawn_Structure)
    - All terms retuned using the full Zurichess dataset.

Though not as large as past gains, these tweaks and tunings show promise, and indicate Blunder's gained anywhere between 45-60 Elo.

Blunder 7.2.0
-------------

This release is not a notable strength improvement over 7.1.0. However, 7.2.0 does introduce a polyglot opening book loader as the primary new feature. Several UCI options are provided to make usage of the loader. Addtionally, various places in the codebase have been refactored and cleaned-up, and the evaluation for Blunder has been retuned and restructured to allow for more granularity. Lastly, a makefile is now included so Blunder can be more easily compiled across platforms.

* Engine
    - Polyglot opening book loader
* Evaluation
   - Refactored & retuned

Blunder 7.1.0
-------------

The release of this version includes late-move reductions, basic futility pruning, as well as a static-exchange evaluation routine, as well as a little refractoring and bug-fixing here and there. 

* Search
    - [Late-move reductions](https://www.chessprogramming.org/Late_Move_Reductions)
    - [Futility pruning](https://www.chessprogramming.org/Futility_Pruning)
    - [Static-exchange evaluation](https://www.chessprogramming.org/Static_Exchange_Evaluation)

The addition of these features seem promising and show a little over a 100 point Elo gain in self-play (tc=inf/10+0.1):

```
Score of Blunder 7.1.0 vs Blunder 7.0.0: 347 - 146 - 133  [0.661] 626
...      Blunder 7.1.0 playing White: 174 - 66 - 74  [0.672] 314
...      Blunder 7.1.0 playing Black: 173 - 80 - 59  [0.649] 312
...      White vs Black: 254 - 239 - 133  [0.512] 626
Elo difference: 115.6 +/- 25.1, LOS: 100.0 %, DrawRatio: 21.2 %
SPRT: llr 2.96 (100.5%), lbound -2.94, ubound 2.94 - H1 was accepted
Finished match
```

The Elo gains for each feature can be seen in the commit history. Going forward this will be the place where I try to document Elo gains (if any) from new features.

Blunder 7.0.0
-------------

Blunder 7.0.0 includes a variety of new features, some strength gaining, some not. Most notably, I've added basic mobility evaluation to blunder, which was made possible by a Texel Tuner implementation, and I've readded history heuristics, and have gotten Principal Variation Search working. On top of these changes, I've also fixed a bug in the way Blunder counts nodes, and I've added several UCI features such as `go depth`, `go nodes`, `Clear History`, and I'm now collecting and reporting a principal variation. Lastly, Blunder has switched from using a fail-hard to a fail-soft negamax implementation. The updates are summarized below.

* Engine
    - UCI features
    - Fixed node counting bug
* Search
    - [History Heuristics](https://www.chessprogramming.org/History_Heuristic)
    - [Principal Variation Search](https://www.chessprogramming.org/Principal_Variation_Search)
    - [Fail-Soft](https://www.ics.uci.edu/~eppstein/180a/990202b.html)
    - [Principal Variation](https://www.chessprogramming.org/Principal_Variation)
* Evaluation
   - [Mobility](https://www.chessprogramming.org/Mobility)
   - [Texel Tuner](https://www.chessprogramming.org/Texel%27s_Tuning_Method)

These additions show 7.0.0 is a little more than 100 Elo stronger than Blunder 6.0.0 (6.1.0 is not disscussed here as it is a slight, but more stable regression from 6.0.0's strength) in self-play, and combined with gauntlet testing, gives an strength estimate of 2250-2350 Elo.

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

The motivations behind this rewrite were twofold: First, I was dissatisfied with Blunder's speed and wanted to take another crack at creating an engine that was simply faster. Second, I didn't like how Blunder was designed in several places, and I quickly realized these "several places" constitued large chunks of the code-base. So Blunder 5.0.0 was the result.

From the testing I've done, Blunder 5.0.0 is 20-30% faster than Blunder 4.0.0, and perft(6) from the starting position was coming in at around 6-8s (14-18 Mnps), whereas perft(6) from the starting position for Blunder 4.0.0 was generally 10-12s (10-12Mnps). And overall, I'm happy with the refractoring I've done and my code feels cleaner in many of the places that bothered me. So both goals, all things considered, were met.

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
