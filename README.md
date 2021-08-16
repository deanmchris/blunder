Overview
--------

Blunder is an open-source UCI compatible chess engine. The philosophy behind Blunder's design is for the code to 
straightforward and easy to read, so that others can benefit from the project. Currently my estimate is that Blunder
is at about 1900-2000 Elo (in self-play), and Blunder 4.0.0 is rated [1700 on the CCRL Blitz list](http://ccrl.chessdom.com/ccrl/404/cgi/engine_details.cgi?print=Details&each_game=1&eng=Blunder%204.0.0%2064-bit#Blunder_4_0_0_64-bit). Blunder 5.0.0 has yet to
be tested.

Installation
-----

Compiling Blunder is fairly simple.

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

Features
--------

* Engine
    - [Bitboards representation](https://www.chessprogramming.org/Bitboards)
    - [Magic bitboards for slider move generation](https://www.chessprogramming.org/Magic_Bitboards)
    - [Zobrist hashing]().
    - [Transposition table for perft]().
* Search
    - [Negamax search framework](https://www.chessprogramming.org/Negamax)
    - [Alpha-Beta pruning](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning)
    - [MVV-LVA move ordering](https://www.chessprogramming.org/MVV-LVA)
    - [Quiescence search](https://www.chessprogramming.org/Quiescence_Search)
    - [Time-control logic supporting classical, rapid, bullet, and ultra-bullet time formats](https://www.chessprogramming.org/Time_Management).
    - [Repition detection](https://www.chessprogramming.org/Repetitions).
    - [Killer moves](https://www.chessprogramming.org/Killer_Move).
    - [Transposition table](https://www.chessprogramming.org/Transposition_Table).
    - [Null-move pruning](https://www.chessprogramming.org/Null_Move_Pruning).
* Evaluation
    - [Material evaluation](https://www.chessprogramming.org/Material)
    - [Tuned piece-square tables](https://www.chessprogramming.org/Piece-Square_Tables)
    - [Tapered evaluation](https://www.chessprogramming.org/Tapered_Eval).
    
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
