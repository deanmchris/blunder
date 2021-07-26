Overview
--------

Blunder is an open-source UCI compatible chess engine. The philosophy behind Blunder's design is for the code to 
straightforward and easy to read, so that others can benefit from the project. Currently my estimate is that Blunder
is at about ~1400 Elo, so around the level of a decent ameatuer chess player.

Installation
-----

Blunder is current being developed in a linux 64-bit enviorment, so that is the executable that is provided. However, 
compiling Blunder on different machines is fairly simple.

All that is needed is to visit [Golang downlaod page](https://golang.org/dl/), and install Golang using the download
package appropriate for your machine. To make using the Golang compiler easier, make sure that if the installer asks,
you let it add the Golang compiler command to your path.

Your installation should be up and running in about 5-7 minutes, and from there, you need to open up a terminal/powershell/
command line, navigate to `blunder/blunder`, and run `go build`. This will create an executable for your computer, which you
should then able to run.

Usage
-----

Blunder, like many chess engines, does not include its own GUI for chess playing, but supports something
known as the [UCI protocol](http://wbec-ridderkerk.nl/html/UCIProtocol.html). This protocol allows chess engines, like Blunder, 
to communicate with different chess GUI programs.

So to use Blunder, it's reccomend you install one of these programs. Popular free ones include:

* [Arena](http://www.playwitharena.de/)
* [Scid](http://scidvspc.sourceforge.net/)
* [Cute-chess](https://cutechess.com/) 

Once you have a program downloaded, you'll need to follow that specfic programs guide on how to install a chess engine. When prompted 
for a command or executable, direct the GUI to the Golang exectuable you built.

In addition to being UCI compatible, Blunder has a simple command line interface. To use it, start the executable, and type in `debug`. 
But wait 3-5 seconds before entering debug, as Blunder needs a second or two to initalize its internals. After typing in `debug`, you
should be greeted with an ASCII display of Blunder's current board state, and a prompt to enter a command:

```

8 | r n b q k b n r 
7 | p p p p p p p p 
6 | . . . . . . . . 
5 | . . . . . . . . 
4 | . . . . . . . . 
3 | . . . . . . . . 
2 | P P P P P P P P 
1 | R N B Q K B N R 
   ----------------
    a b c d e f g h 

turn: white
castling rights: KQkq
en passant: none
rule 50: 0
game ply: 0

>> 

```

Type "options" into the prompt
to see the available commands.

Features
--------

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

 Future Features
 ---------------
 One of the most fun and exciting parts of chess programming is adding new features to your engine, and watching it
 slowly become stronger and better than previous versions. 
 
 With that said, here are some features that will be added in the upcoming versions of Blunder (in no particular
 order):
 
 * A transposition table using zobrist hashing
 * Repition draw detection
     - I'd like to add this feature especially, since Blunder currently draws many games because
       it sees a position as being equal, or that it has an advantage, no matter what move it makes,
       so it simply repeats a particular move, not realizing that it'll draw the game soon.
 * Tapered evaluation
 * Null-move pruning
 * Killer heuristics
 * History heuristics
 * Pawn structure evaluation
 * Texel tuning
    
 Changelog
 ---------
 
 The changelog of features for Blunder can be found in the `docs/changelog.md`.
 
 Credits
 -------
 
 Although Blunder is an orginal project, there are many people without whom Blunder would not have been finished. 
 The brief listing is included here (in no particular order). For the full listing, with elaborations, 
 see `docs/credits.md`:
 
 ```
 My girlfriend, Marcel Vanthoor, Hart Gert Muller, Sven Schüle, J.V. Merlino, Niels Abildskov, 
 Maksim Korzh, Erik Madsen, Pedro Duran, Nihar Karve, Rhys Rustad Elliott, Lithander, 
 Jonatan Pettersson, Rein Halbersma, Tony Mokonen, SmallChess, Richard Allbert, Spirch, and
 the Stockfish Developers.
 ```
 
 These credits will be updated from time to time as a remember or encounter more people who have helped me
 in Blunder's development.

 License
 -------
 
 Blunder is licensed under the [MIT license](https://opensource.org/licenses/MIT).
