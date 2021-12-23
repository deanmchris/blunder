Overview
--------

[Blunder](http://ccrl.chessdom.com/ccrl/404/cgi/compare_engines.cgi?family=Blunder&print=Rating+list&print=Results+table&print=LOS+table&print=Ponder+hit+table&print=Eval+difference+table&print=Comopp+gamenum+table&print=Overlap+table&print=Score+with+common+opponents) is an open-source UCI compatible chess engine. The philosophy behind Blunder's design is for the code to 
be straightforward and easy to read, so that others can benefit from the project.

| Version     | Estimated Rating (Elo) | CCRL Rating (Elo) | 
| ----------- | -----------------------|-------------------|
| 1.0.0       | 1400                   | N/A               |
| 2.0.0       | 1570                   | N/A               |
| 3.0.0       | 1782                   | N/A               |
| 4.0.0       | 1832                   | 1734              |
| 5.0.0       | 2000                   | 2080              |
| 6.0.0       | 2200                   | N/A               |
| 6.1.0       | 2200                   | 2155              |
| 7.0.0       | 2280                   | N/A               |
| 7.1.0       | 2395                   | N/A               |
| 7.2.0       | 2395                   | 2425              |
| 7.3.0       | 2450                   | N/A               |
| 7.4.0       | 2510                   | 2532              |

Installation
-----

Builds for Windows, Linux, and MacOS are included with each release of Blunder. However, if you
prefer to build Blunder from scratch the steps to do so are outlined below.

Visit [Golang download page](https://golang.org/dl/), and install Golang using the download
package appropriate for your machine. To make using the Golang compiler easier, make sure that if the installer asks,
you let it add the Golang compiler command to your path.

Your installation should be up and running in about 5-7 minutes, and from there, you need to open up a terminal/powershell/
command line, navigate to `blunder/blunder`, and run `go build`. This will create an executable for your computer, which you
should then able to run.

Alternatively, if the `make` build automation tool is installed on your computer (it comes standard on most Linux systems),
simply download this repository's zip file, unzip it, navigate to the primary folder, and run `make` from the command line.
An executable for Windows, Linux, and MacOS will be built and placed inside of the primary directory.

If you're on a windows platform, you'll need to run `make build-windows` instead.

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
    - [Zobrist hashing](https://www.chessprogramming.org/Zobrist_Hashing)
* Search
    - [Negamax search framework](https://www.chessprogramming.org/Negamax)
    - [Alpha-Beta pruning](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning)
    - [MVV-LVA move ordering](https://www.chessprogramming.org/MVV-LVA)
    - [Quiescence search](https://www.chessprogramming.org/Quiescence_Search)
    - [Time-control logic supporting classical, rapid, bullet, and ultra-bullet time formats](https://www.chessprogramming.org/Time_Management).
    - [Repetition detection](https://www.chessprogramming.org/Repetitions)
    - [Killer moves](https://www.chessprogramming.org/Killer_Move)
    - [Transposition table](https://www.chessprogramming.org/Transposition_Table)
    - [Null-move pruning](https://www.chessprogramming.org/Null_Move_Pruning)
    - [Reverse futility pruning](https://www.chessprogramming.org/Reverse_Futility_Pruning)
    - [History Heuristics](https://www.chessprogramming.org/History_Heuristic)
    - [Principal Variation Search](https://www.chessprogramming.org/Principal_Variation_Search)
    - [Fail-Soft](https://www.ics.uci.edu/~eppstein/180a/990202b.html)
    - [Late-move reductions](https://www.chessprogramming.org/Late_Move_Reductions)
    - [Futility pruning](https://www.chessprogramming.org/Futility_Pruning)
    - [Static-exchange evaluation](https://www.chessprogramming.org/Static_Exchange_Evaluation)
    - [Aspiration windows](https://www.chessprogramming.org/Aspiration_Windows)
* Evaluation
    - [Material evaluation](https://www.chessprogramming.org/Material)
    - [Tuned piece-square tables](https://www.chessprogramming.org/Piece-Square_Tables)
    - [Tapered evaluation](https://www.chessprogramming.org/Tapered_Eval)
    - [Mobility](https://www.chessprogramming.org/Mobility)
    - [Basic king safety](https://www.chessprogramming.org/King_Safety)
    - [Basic pawn structure](https://www.chessprogramming.org/Pawn_Structure)
    - [Knight outposts](https://www.chessprogramming.org/Outposts)
    - [Texel Tuner](https://www.chessprogramming.org/Texel%27s_Tuning_Method)
    
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
 
 These credits will be updated from time to time as I remember or encounter more people who have helped me
 in Blunder's development.

 License
 -------
 
 Blunder is licensed under the [MIT license](https://opensource.org/licenses/MIT).
