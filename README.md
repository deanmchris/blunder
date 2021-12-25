Overview
--------

Blunder is an open-source UCI compatible chess engine. The philosophy behind Blunder's design is for the code to 
be straightforward and easy to read, so that others can benefit from the project.

History
-------

The inspiration for Blunder started near the beginning of 2021. Me and many of my friends had recently started playing chess more seriously, and having a couple of years of programming knowledge, I imagined it would be fun to create my own chess playing program. I started a very rough first version, written in Python, but soon abandonded it, as I realized writing a chess engine was a much more daunting project then I had first anticpated. 

With my intial failure, I started doing more research and discovered the rich field of computer programming, and the many helpful people that are a part of it. About 5 months and ten attempts later, I released the first version of Blunder! And I've been working to improve Blunder ever since. As for the programming language switch, though Python is an amazing language (I think anyway), and the first language I learned, it's simply not fast enough for the purpose of writing a relatively strong chess engine. So instead of writing another C/C++ chess engine, I decided to give Go a try, and I've enjoyed working with its tools.

I've also [started a blog](https://algerbrex.github.io/) to track and write about the development of Blunder.

Ratings
-------

When discussing an engine's (or human chess player's) strength, it's important to remember that the Elo is always relative to one's testing conditions. One tester may estimate an engine's strength to be 2300 for example, while another may get 2245. Neither tester is "wrong" per se, but they both likely have a different pool of opponets, different hardware, different time controls, etc.

With that said, several people have been kind enough to test various versions of Blunder, and a summary of the rating list and their ratings for the versions are listed below:

| Version     | Estimated Rating (Elo) | CCRL Blitz Rating (Elo) | Bruce's Bullet Rating List (ELo) |
| ----------- | -----------------------|-------------------------|----------------------------------|
| 1.0.0       | 1400                   | -                       | -                                |
| 2.0.0       | 1570                   | -                       | -                                |
| 3.0.0       | 1782                   | -                       | -                                |
| 4.0.0       | 1832                   | 1734                    | -                                |
| 5.0.0       | 2000                   | 2080                    | 2174                             |
| 6.0.0       | 2200                   | -                       | 2248                             |
| 6.1.0       | 2200                   | 2155                    | 2226                             |
| 7.0.0       | 2280                   | -                       | 2374                             |
| 7.1.0       | 2395                   | -                       | 2455                             |
| 7.2.0       | 2395                   | 2425                    | 2472                             |
| 7.3.0       | 2450                   | -                       | 2499                             |
| 7.4.0       | 2510                   | 2532                    | 2554                             |
| 7.5.0       | 2540                   | ?                       | ?                                |

* [CCRL Blitz Rating List](http://ccrl.chessdom.com/ccrl/404/)
* [Bruce's Bullet Rating List](https://e4e6.com/)

A very big thank you to those who have helped and continue to help test Blunder.

Installation
-----

Builds for Windows, Linux, and MacOS are included with each release of Blunder. However, if you
prefer to build Blunder from scratch the steps to do so are outlined below.

Visit the [Golang download page](https://golang.org/dl/), and install Golang using the download
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
    - [Late-move pruning/move-count based pruning](https://www.chessprogramming.org/Futility_Pruning#MoveCountBasedPruning)
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
 
 Resources
 ---------
 
 This list is by no means exhaustive, but here are some of the main resources that I've found and cotinue to find helpful while developing Blunder:
 
* [The Chess Programming Wiki](https://www.chessprogramming.org/Main_Page)
* [The Chess Stack Exchange site](https://chess.stackexchange.com/)
* [Talkchess](http://talkchess.com/forum3/index.php)
* [Programming a chess engine in C](https://www.youtube.com/watch?v=bGAfaepBco4&list=PLZ1QII7yudbc-Ky058TEaOstZHVbT-2hg)
* [Programming a chess engine in Javascript](https://www.youtube.com/watch?v=2eA0bD3wV3Q&list=PLZ1QII7yudbe4gz2gh9BCI6VDA-xafLog)
* [Bitboard engine in C](https://www.youtube.com/watch?v=QUNP-UjujBM&list=PLmN0neTso3Jxh8ZIylk74JpwfiWNI76Cs)
* [Logic Crazy's Chess Engine Tutorial](https://www.youtube.com/watch?v=V_2-LOvr5E8&list=PLQV5mozTHmacMeRzJCW_8K3qw2miYqd0c)

 License
 -------
 
 Blunder is licensed under the [MIT license](https://opensource.org/licenses/MIT).
