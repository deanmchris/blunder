package tuner

import (
	"blunder/engine"
	"testing"
)

func init() {
	engine.InitBitboards()
	engine.InitTables()
	engine.InitZobrist()
	engine.InitEvalBitboards()
}

var TestFENs = []string{
	"8/3k1p2/R7/P1Pb2P1/3P4/4KP2/6P1/8 b - - 0 1",
	"8/8/8/4k3/6K1/7N/7P/3n4 w - - 0 1",
	"3R4/5kpp/1p3n2/p1p2p2/8/P3PP2/1r4PP/R4K2 b - - 0 1",
	"1R6/8/8/4pP1k/2r5/5K2/8/8 w - - 0 1",
	"6rk/1pqb3r/1p6/4Pp2/2p2P2/P1N1P1P1/1P2QKB1/8 b - - 0 1",
	"8/5p2/Q3p1k1/3pP1qp/1P6/4P1PP/4K3/8 w - - 0 1",
	"2Q5/5q1p/6k1/pp4p1/3B4/P7/5P2/6K1 w - - 0 1",
	"8/4R3/3K4/P6p/r5p1/4PkP1/8/8 w - - 0 1",
	"5b2/1p3k2/p7/5bpp/P3p3/1P2P2P/2R2PP1/6K1 w - - 0 1",
	"3K4/8/6k1/1p6/p2PpP2/8/8/8 b - - 0 1",
	"8/8/6k1/K5pp/6Pn/4R2P/8/8 w - - 0 1",
	"8/7k/4pp1p/1p6/4K2P/2P2P2/6P1/8 b - - 0 1",
	"4k3/8/1p3p2/8/PP6/1KPr4/8/8 w - - 0 1",
	"B5k1/5p2/p3p1p1/7p/bPP2P1b/4N2P/6r1/1RR2K2 b - - 0 1",
	"8/3K4/7r/p7/1P1p1P2/8/3R3p/5k2 w - - 0 1",
	"8/2P5/6pp/1K5k/8/8/7P/8 b - - 0 1",
	"6rr/5kp1/2p2p2/1p1p4/1P1P3P/1KP3P1/5R2/5R2 b - - 0 1",
	"8/8/4p3/1p2p2P/2p1P2K/3k4/8/8 w - - 0 1",
	"8/6pk/p1N4p/6p1/1P4K1/7P/3R4/4r3 b - - 0 1",
	"rn1qk2r/pppbbppp/5n2/3p4/3P4/4PN1P/PP3PP1/RNBQKB1R w KQkq - 0 1",
	"8/5kp1/3n4/1P1p4/6PP/3K1P2/5B2/8 w - - 0 1",
	"8/6R1/7k/8/8/5K2/7p/4b3 w - - 0 1",
	"8/6p1/2Q1pnkp/5p2/8/1q3BP1/5PK1/8 b - - 0 1",
	"8/3RP3/1kp5/p2b3P/6B1/8/1K6/8 b - - 0 1",
	"2N5/1p1R3p/p1n2ppk/3B4/P5P1/8/5P2/6K1 b - - 0 1",
	"rn1qk1nr/2p2ppp/b7/3p4/8/3P1N2/PP3PPP/R2QKBNR b KQkq - 0 1",
	"5k2/6pb/7p/3P1p1P/3P2Pq/1R3P1B/8/7K w - - 0 1",
	"r1b1kbnr/1q3ppp/p3p3/2ppP3/5P2/P1NP1N2/1P1QB1PP/1R2K2R w Kkq - 0 1",
	"r2qkb1r/pp1b1pp1/3p1n1p/4p3/2Pp4/2NPB2P/PPQ1PPP1/2R1KB1R w Kkq - 0 1",
	"r1bqk2r/1p2bp2/2n1p3/3p4/p2PN2p/2PBPNP1/P4P1P/1R1Q1RK1 w q - 0 1",
	"r1bqkbnr/pp1n1p1p/2pP4/7p/4PQ2/7P/PPP2PP1/RNB1KB1R b KQkq - 0 1",
	"2k5/2p5/pp1p4/3Pn3/1P2P3/5p2/P7/6K1 w - - 0 1",
	"1r6/5kp1/5p2/3Pp3/8/P1N1PP2/8/K6n w - - 0 1",
	"1r6/2r1nk2/3ppppp/3p4/3P2P1/1RP1P1KP/1P1N1P2/2R5 b - - 0 1",
	"6k1/1q3rpp/p3p3/2pn2N1/P1Nb4/8/1PQ3PP/R6K b - - 0 1",
	"8/3K4/5k2/8/8/5p2/1p6/8 w - - 0 1",
	"8/6b1/8/5p1p/6kP/8/3KN3/8 w - - 0 1",
	"rnbqkbnr/3p1ppp/4p3/p2p4/Pp5P/2P1B3/1P2PPP1/RN1QKBNR b KQkq - 0 1",
	"2b2kr1/rpp1bp2/2P1pn2/4N2p/1P4P1/PB3Q1P/1B1N4/2K1R2R w - - 0 1",
	"8/6k1/2b1Q1p1/3p1p2/3P1P2/2PB4/7K/8 b - - 0 1",
	"rnbqk2r/pppn1ppp/4p3/4P3/3pNP1b/P5P1/1PPPB2P/R1BQK1NR b KQkq - 0 1",
	"4q2k/1p2r1p1/p4pbp/P4P2/2P5/3Q3P/6PK/1R6 b - - 0 1",
	"r7/3nkpp1/1p2p2p/2Pp3P/b5P1/2KNBP2/2P1B3/4R3 b - - 0 1",
	"2r5/q7/6k1/8/8/2N4P/1P4P1/7K w - - 0 1",
	"3r4/3b1p2/8/R3P3/4P1k1/8/4K3/8 w - - 0 1",
	"8/3bk3/p3p2r/2P5/2P2R2/P4P2/1K2B3/8 w - - 0 1",
	"rn3rk1/pp1bq2p/1npb2p1/3p3Q/3P1N1P/1PPBP3/P4P2/RNB1K2R w KQ - 0 1",
	"r7/p3Qbk1/1p3r2/2p5/4P2P/2P1BP2/PP2B3/2K5 b - - 0 1",
	"8/Q4nk1/6pp/6q1/2P5/6P1/PP5P/6K1 w - - 0 1",
	"6k1/4b1p1/7p/2pPp3/2P1P3/1q5P/6P1/6K1 w - - 0 1",
	"3r4/1B3kp1/1b2pn1p/1p3p1P/1P3P2/6PK/1P2R3/4B3 w - - 0 1",
	"6R1/3r1p1k/8/5nBp/1P5p/3p2P1/7K/8 w - - 0 1",
	"3r4/5kp1/8/p7/2R1B3/6P1/P4b1P/7K w - - 0 1",
	"8/8/p7/4bp2/1k5P/p4r2/P1P5/3K4 b - - 0 1",
	"6k1/1R4p1/r3p3/1p1pPn1p/p6P/P5P1/1P3PK1/6N1 b - - 0 1",
	"8/8/1n6/4p3/1B2k3/8/5K2/8 w - - 0 1",
	"1r6/1p6/2kp1p1b/4pP1p/1B2K2P/P2RP3/8/8 b - - 0 1",
	"4r3/2Q3pk/p6p/p4q2/4r3/6PP/P6K/1RB5 b - - 0 1",
	"4r1k1/5rpp/3p4/8/7P/1Q2Pq2/PPP3R1/1K6 w - - 0 1",
	"8/8/4B3/6K1/1p6/3k2b1/5p2/8 w - - 0 1",
	"5K2/2N1B2p/1p4k1/1b4P1/8/r7/8/8 b - - 0 1",
	"6k1/p6p/6p1/4b3/2RN2p1/3nP2P/3r1P2/6K1 w - - 0 1",
	"8/8/k5N1/1p2P3/5PK1/r7/6R1/8 b - - 0 1",
	"8/1R6/K1k2p1p/p2r3P/P4P2/bP6/3p4/3R4 b - - 0 1",
	"5rk1/p3R1pp/1p3r2/2p5/2Pp4/P2K2P1/1P5P/4R3 b - - 0 1",
	"r1b1k1nr/p3ppbp/p1p3p1/5qN1/3P4/N1P5/PP3PPP/R1BQK2R w KQkq - 0 1",
	"1R6/6p1/5pk1/P6R/3P2KP/5P2/r7/r7 w - - 0 1",
	"6k1/5p2/6p1/4P2p/1BP4N/4P3/5KPP/8 b - - 0 1",
	"rnb2rk1/ppqpb1pp/4pn2/5p2/2PP4/P1N2N2/1P2PPPP/R2QKB1R w KQ - 0 1",
	"6k1/3q2p1/8/5pN1/1p1p1P2/1P2r3/P1Q2KR1/8 w - - 0 1",
	"3r2k1/1p3ppb/pP5p/P4P2/4p1P1/4P2P/4B3/R5K1 b - - 0 1",
	"8/8/KP4k1/8/8/1r3p2/8/8 w - - 0 1",
	"5R2/8/1kp5/p1p1PpP1/P7/3K4/7r/8 b - - 0 1",
	"1k6/1p6/4Kpr1/3Nb3/7p/4P2P/8/1R6 w - - 0 1",
	"r2q1rk1/1p3p2/p1np3p/P1p1p1pP/4P1b1/2PP4/1P1NBPP1/R2Q1RK1 b - - 0 1",
	"8/5pp1/4k3/p2r3p/R5P1/1K5P/8/8 w - - 0 1",
	"r3kbnr/ppp2pp1/B7/3p2Bp/3n4/P6P/1PP2P1P/RN1QK1R1 b Qkq - 0 1",
	"4rrk1/2Rbb1p1/4q3/1p1pPp2/4p3/P3P1PB/1P5P/3Q1RK1 w - - 0 1",
	"8/4k3/6K1/7R/8/8/8/5r2 w - - 0 1",
	"7k/p3bp2/1p6/2p3pp/4P3/2PrBRP1/PP5P/1K6 b - - 0 1",
	"2n4k/8/4r3/p7/8/5P1P/1RP4K/8 b - - 0 1",
	"3r3k/1p4p1/5p1p/P3p3/8/P3P2P/2Q2PPK/2R5 b - - 0 1",
	"r1b3k1/8/2n5/1p1p4/1P1Pp2p/2P4P/2Q1N1q1/1K3R2 w - - 0 1",
	"4R3/5p1k/8/p4N1p/2P5/5KP1/8/1r6 w - - 0 1",
	"Q7/5pk1/6pp/8/3P4/3BPK1P/1q3PP1/8 b - - 0 1",
	"8/5pp1/4n2p/1R6/8/5k1P/3N4/6K1 b - - 0 1",
	"6k1/q3bbpp/1p1r1p2/r1p1p3/2P4P/3P1Q2/5PP1/RRB1N1K1 w - - 0 1",
	"5rk1/pp4p1/2B3q1/7p/2n5/7P/PP1Q1P1K/R7 w - - 0 1",
	"r7/3b4/1p2p1k1/6pp/2PB4/P4PK1/1P5P/4R3 b - - 0 1",
	"7k/7q/p1p1p2p/P2p2pN/1P1P4/7P/7K/5R2 w - - 0 1",
	"8/8/8/3p1k2/6rP/P6R/5K2/8 b - - 0 1",
	"5r2/1p2k3/6p1/1pp4p/7P/1P1RK1P1/8/8 b - - 0 1",
	"2r3k1/2P2p2/p1P1b2p/B7/4R2P/1P4P1/P7/6K1 b - - 0 1",
	"8/1p2r3/6r1/p1P3k1/P3P2p/5RPB/5P2/5K2 w - - 0 1",
	"3B4/1B3p2/r3p1p1/p1P4p/P2P1k2/4n2P/4KP2/8 b - - 0 1",
	"8/8/8/7p/7r/6R1/4pk2/1K6 w - - 0 1",
	"4k1r1/3b1p2/7p/PNq1p1bP/2PpPr2/3B1P2/3Q2PR/R6K b - - 0 1",
	"3r1k2/5p1p/2B3b1/1P2R3/1r5p/6P1/P4P1P/R5K1 w - - 0 1",
	"8/8/8/2k1p1p1/3n2rp/8/R7/7K w - - 0 1",
	"4q3/8/5k2/8/2K2P2/7P/8/8 w - - 0 1",
	"rr4k1/5npp/4Rp2/1BpP4/5P1P/1q4P1/8/5RK1 w - - 0 1",
	"8/3k1p2/1p2p2p/4P2P/P1BnK3/8/8/8 b - - 0 1",
	"r2r4/1p5p/p1k5/2pNn1p1/P3P1b1/1P6/2P3PP/2R1R1K1 w - - 0 1",
	"7r/1p1k2R1/p4n1p/3p4/PP1R1P1P/2rB1K2/7P/8 b - - 0 1",
	"8/8/P7/6r1/3K4/7R/5k1P/8 b - - 0 1",
	"rr4k1/2pn1pb1/p3pnpp/8/P1BPpBP1/2P1P2P/1P1N1P2/R3K2R w KQ - 0 1",
	"r1bqk3/ppp1bppr/5n2/7p/1Pp5/2N1P3/1B1P1PPP/R2QKBNR b KQq - 0 1",
	"3Q4/k7/8/1P1B1b2/8/3p4/7q/2K5 w - - 0 1",
	"8/3K4/8/2P3p1/8/7k/q7/8 w - - 0 1",
	"4r3/2pq1pk1/bp6/p2p2P1/3P3p/2Q1P1n1/PPNB4/1K4R1 w - - 0 1",
	"8/8/4N3/3K4/6pk/1p6/7n/8 b - - 0 1",
	"4Q1k1/pq3ppp/2n5/2P3b1/8/PN1R1P2/1P2R1KP/8 b - - 0 1",
	"4r1k1/2Q2qpp/4p3/3pP3/p1p2P1P/P7/2P3KP/5R2 w - - 0 1",
	"r5k1/pQ3ppp/3b4/1b6/3n3N/PP6/3B2PP/R3R1K1 b - - 0 1",
	"8/5pp1/6kp/q7/3Pp3/4P3/1r3P2/2R3K1 w - - 0 1",
	"4k3/2R2p2/8/3K2pp/8/8/8/8 w - - 0 1",
	"r7/1bb2pk1/p6p/Pp6/4r3/5N1P/QR3PP1/2R3K1 b - - 0 1",
	"5rk1/4ppbp/bp2n1p1/p2rP3/3P1P2/P3B1P1/RP3N1P/4R1K1 w - - 0 1",
	"8/5p1p/5p2/2P3k1/3RB1P1/1r5P/r1PK1P2/8 w - - 0 1",
	"8/p7/6K1/1P6/4k3/Q7/8/8 b - - 0 1",
	"8/8/1p1K4/3P4/B2k4/8/5b2/8 b - - 0 1",
	"r1b1kb1r/pp1p1ppp/5p2/8/1n2N3/1P2P3/2P2PPP/1K1R1BNR b kq - 0 1",
	"qr4k1/5pp1/8/1p1p3P/8/Q3P1P1/p4PK1/2R5 w - - 0 1",
	"r4rk1/1p1bbppp/3q4/p7/2QN1P2/4P3/PP4PP/R1B2RK1 b - - 0 1",
	"8/1R3p2/8/3r4/1P5k/b7/7K/5B2 w - - 0 1",
	"8/1p4pk/7p/5P2/3r1P2/4K2P/3P4/8 b - - 0 1",
	"3r2k1/1b3pp1/pqp1p2p/1P1n4/1b1P4/5N1P/1PBQ1PP1/2R3KR w - - 0 1",
	"8/B4kp1/n3p3/6p1/rbR1P1P1/8/2K5/6N1 w - - 0 1",
	"8/8/8/Bp4k1/p1p5/5K2/8/7r w - - 0 1",
	"8/3r3k/1p2N1pp/pP2n3/P3p3/4P1P1/5RK1/1R6 b - - 0 1",
	"3r4/1kp2p2/1n2p1p1/7p/p2nR3/B4PP1/1R4KP/8 b - - 0 1",
	"r2qkb1r/pb1ppp1p/1pn2n2/2p3B1/P3P3/2NP2P1/1PP2PBP/R2QK1NR b KQkq - 0 1",
	"6k1/3b1p2/6pp/p2B4/1p3KP1/1P3P2/P5P1/8 w - - 0 1",
	"3b4/8/1p3N2/pP6/k3pPPK/4P3/8/8 w - - 0 1",
	"r1bk3r/1ppp1pp1/4q3/B3P2P/3Q2B1/P3P3/1PP2PP1/R4RK1 b - - 0 1",
	"8/8/3k4/pP6/2P1Pp1B/3Kb3/8/8 w - - 0 1",
	"rnbqk1nr/ppp4p/3p1pp1/4p3/8/2P3PN/P1P1PPBP/1RBQK2R b Kkq - 0 1",
	"7k/2r3pp/3N4/8/P3Pp2/7P/4KPP1/8 b - - 0 1",
	"r2b1rk1/1ppq2p1/3pp1bp/5p2/PnNPn3/3BPN1P/1B2QPP1/R3R1K1 w - - 0 1",
	"8/8/BP1k4/4p3/4Pn1p/8/3K4/8 b - - 0 1",
	"8/8/3p1P1k/1b1Pp1p1/p2p2B1/B7/5K2/8 b - - 0 1",
	"4r1k1/1p3b1p/1qnb1pp1/1p1p4/3P4/P4N1P/1PQ1NPP1/3R2K1 w - - 0 1",
	"8/8/5R2/8/8/4K1kP/1r6/8 w - - 0 1",
	"8/4nk1p/p4ppb/1p1b4/3P2KP/5B2/1B3PP1/4N3 b - - 0 1",
	"6k1/R1nr1q2/1p2b1p1/1Pp1Q3/5P2/2BP1P1p/7P/5RK1 w - - 0 1",
	"8/6P1/3R4/2r5/6K1/8/2kp4/8 b - - 0 1",
	"8/6kp/5n2/p2p4/P2PR3/1p4P1/1P3KP1/4R3 w - - 0 1",
	"Q3k3/1Q6/2K1p1r1/8/8/8/8/3r4 b - - 0 1",
	"8/1b6/3b4/1P1Pp3/3pB1pk/8/5KP1/8 b - - 0 1",
	"r2qkb1r/pppb1p1p/5n2/3pn3/PP2p3/2PP4/3BPPPP/RN1QKB1R b KQkq - 0 1",
	"6k1/1R2b1p1/1pn4p/1r2p3/4r3/P1B2NP1/1P3P2/2R2K2 w - - 0 1",
	"5R2/8/5P2/3k3p/7P/8/6K1/r7 w - - 0 1",
	"3k4/1p6/2p5/5pp1/8/2P1K2P/P1P1B2n/8 b - - 0 1",
	"8/8/2N1p3/5p1k/7b/1R4K1/8/3r4 w - - 0 1",
	"r2q1rk1/pb2bppp/3p1n2/1p4N1/5P2/P3BB2/1PP3PP/R2Q1K1R b - - 0 1",
	"2r5/1p1k1p2/1n6/p2p4/Pb1qp1b1/2N4p/1P1NP2P/3QK2B w - - 0 1",
	"2r2Bk1/1b1q1pp1/1p2p2p/pN1pP3/P1PP3b/1P1Q3P/4B1P1/5R1K w - - 0 1",
	"8/4kp2/PK6/6rp/8/1R6/8/8 w - - 0 1",
	"r1bq1k1r/ppppb1p1/2n1p3/2N2p1n/3P4/4P3/PPP2PPP/1RBQK1NR w K - 0 1",
	"4r3/1pb2p1k/r6p/q1pNp1P1/P1P1P3/3P2P1/5R2/R1Q3K1 b - - 0 1",
	"1q1n2k1/1p2npp1/p7/6Pp/7P/P1Qb4/4PP2/4NK1R b - - 0 1",
	"5b2/1r6/1BNp1kp1/2p1p3/4Pp2/7b/1PP1KP2/6R1 w - - 0 1",
	"6k1/5p2/2P3pp/p1b5/P7/7P/5PP1/3B2K1 b - - 0 1",
	"2r1k2r/2B1b3/4p2p/2p5/2P1p3/1P1nQ2P/P4PP1/qN3KR1 w k - 0 1",
	"5k2/R7/2r4p/p3K1p1/8/7P/8/8 b - - 0 1",
	"4b3/8/p6p/1p2P1k1/1P1K4/8/1P6/7R b - - 0 1",
	"8/6p1/4kp2/7P/P5R1/6P1/7K/8 b - - 0 1",
	"r1bqk2r/1pp2pb1/4p3/p5p1/3P2p1/1P4B1/P1P2P1P/R2QKBR1 w Qkq - 0 1",
	"2br4/p1r3kp/4p1p1/1RN1R3/8/2p3P1/P4P1P/6K1 b - - 0 1",
	"r1bqk1r1/pp3p2/3p1n1b/3Pp3/P1Pn3P/2NB4/RP3PP1/1N1QR1K1 b q - 0 1",
	"8/6k1/p5Pp/Pp2r3/1P1R1KP1/5P2/8/8 b - - 0 1",
	"2r5/5p1k/7p/1K6/7b/2P4P/R5P1/8 w - - 0 1",
	"4r1k1/pp2nppp/8/1PPp4/P4P1P/2P5/5Q2/4R1K1 b - - 0 1",
	"8/4b3/7p/3KN3/8/2R5/8/1k6 b - - 0 1",
	"8/8/2k5/R5N1/7K/5P2/8/6q1 w - - 0 1",
	"rnbqk2r/p4ppp/1ppbpn2/1B1p4/3P4/2N1PN2/PPP2PPP/R1BQK2R w KQkq - 0 1",
	"r1bqkb1r/ppp1np2/4p2p/4P1p1/4p3/2P3P1/PP3PBP/RNB1QRK1 b kq - 0 1",
	"r7/p2n2pk/2pNB2p/1p6/1P5q/P3P3/5PPP/R2R2K1 b - - 0 1",
	"8/2p4k/p1n2R2/8/3p2P1/1p1P2P1/2PKn3/8 b - - 0 1",
	"8/5p1k/4pp2/2R5/6n1/r6p/3K4/8 b - - 0 1",
	"6k1/6p1/p1p1r1b1/1p6/3P1P2/P7/1P1N1KP1/6R1 b - - 0 1",
	"6k1/5p2/6p1/R3P2p/5P1P/P4BPb/5K2/8 b - - 0 1",
	"r6r/1n2kp2/2p2p2/1p3P1p/3PP1p1/pPP5/P4NPP/1R2R1K1 w - - 0 1",
	"5k2/1R6/8/8/8/3K4/4r3/8 b - - 0 1",
	"5r2/3b3k/1p1p3p/p1pP4/P1P1P1q1/1PB4p/5p1n/1KQ2R2 w - - 0 1",
	"1rr3k1/5ppp/1pBp1nq1/p2Pp3/2P1n3/4RN1P/4QPP1/1R4K1 w - - 0 1",
	"6k1/8/3n4/8/P4p2/2P3p1/K7/5B2 b - - 0 1",
	"r1b4r/ppp1kppp/2n5/3Pp3/8/7P/PP2PPP1/3RKBNR b K - 0 1",
	"6k1/R6p/3N4/3K4/P7/4r3/6p1/8 w - - 0 1",
	"q3k3/3n1p2/2R5/1N1Ppp1r/1R6/P3P1P1/5PKP/8 w - - 0 1",
	"8/8/2p5/p3k2r/P1R5/3K4/2P2P1p/8 w - - 0 1",
	"5N2/8/8/7p/5P2/2k1pK2/8/8 b - - 0 1",
	"8/8/5kpp/p7/P1K2P1P/8/8/8 w - - 0 1",
	"4k3/6p1/4r2p/5p2/2R1P3/1p2K1PP/1P6/8 w - - 0 1",
	"8/8/3k1b2/K4R2/8/3BP3/1r6/8 b - - 0 1",
	"8/1PB1b3/p3kp2/8/PP1r3p/7P/5PK1/1R6 b - - 0 1",
	"8/3P4/4K3/4p3/4k3/8/8/2b2B2 b - - 0 1",
	"rn1qkb1r/1p2pppp/2p2n2/p2p1b2/1P1P4/P1P2NP1/4PPBP/RNBQK2R w KQkq - 0 1",
	"7r/3R4/1p1p4/p3k1p1/P5p1/1P4P1/6K1/8 b - - 0 1",
	"8/8/1r1n4/3pk3/8/K7/8/8 b - - 0 1",
}

// tuner_test.go tests whether the gathering of coefficents and evaluating of positions
// is working correctly in the tuner. This is vital to ensuring the proper derivative
// is calculated and gradient descent will converge.

func TestTuner(t *testing.T) {
	weights := loadWeights()
	for _, fen := range TestFENs {
		pos := engine.Position{}
		pos.LoadFEN(fen)

		normalEval := float64(engine.EvaluatePos(&pos))
		if pos.SideToMove == engine.Black {
			normalEval = -normalEval
		}

		normalCoefficents, safetyCoefficents := getCoefficents(&pos)
		mgPhase := float64(256 - ((pos.Phase*256 + (engine.TotalPhase / 2)) / engine.TotalPhase))
		tunerEval := evaluate(weights, normalCoefficents, safetyCoefficents, mgPhase)

		if abs_float64(normalEval-tunerEval) > 1.5 {
			t.Errorf(
				"For position %s [%s], got %f for normal evaluation score, but %f for tuner evaluation.",
				pos.String(), fen, normalEval, tunerEval,
			)
		} else {
			t.Logf(
				"Pass for position %s: %f = ~%f",
				fen, normalEval, tunerEval,
			)
		}
	}
}

func abs_float64(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}
