package engine

// transposition.go contains an implementation of a transposition table (TT) to use
// in searching and perft.

const (
	// Default size of the transposition table, in MB.
	DefaultTTSize = 64

	// The number of buckets to have per transposition table index.
	NumTTBuckets = 2

	// Constant for the size of a transposition table search and perft entry, in bytes,
	// considering memory alignment.
	SearchEntrySize uint64 = 16
	PerftEntrySize  uint64 = 24

	// Constants representing the different flags for a transposition table entry,
	// which determine what kind of entry it is. If the entry has a score from
	// a fail-low node (alpha wasn't raised), it's an alpha entry. If the entry has
	// a score from a fail-high node (a beta cutoff occured), it's a beta entry. And
	// if the entry has an exact score (alpha was raised), it's an exact entry.
	AlphaFlag uint8 = iota
	BetaFlag
	ExactFlag

	// A constant representing the minimum or maximum value that a score from the search
	// must be below or above to be a checkmate score. The score assumes that the engine
	// will not find mate in 100.
	Checkmate = 9000
)

// A struct for a transposition table entry used in the search.
type SearchEntry struct {
	Hash  uint64
	Depth uint8
	Score int16
	Flag  uint8
	Best  Move
}

// A struct for a transposition table entry used in perft.
type PerftEntry struct {
	Hash  uint64
	Nodes uint64
	Depth uint8
}

func (entry *SearchEntry) Get(hash uint64, ply, depth uint8, alpha, beta int16, best *Move) (int16, bool) {
	adjustedScore := int16(0)
	shouldUse := false

	// Since index collisions can occur, test if the hash of the entry at this index
	// actually matches the hash for the current position.
	if entry.Hash == hash {

		// Even if we don't get a score we can use from the table, we can still
		// use the best move in this entry and put it first in our move ordering
		// scheme.
		*best = entry.Best

		// Return the score of the position to use as an estimate for various
		// pruning and extension techniques in the search.
		adjustedScore = entry.Score

		// To be able to get an accurate value from this entry, make sure the results of
		// this entry are from a search that is equal or greater than the current
		// depth of our search.
		if entry.Depth >= depth {
			score := entry.Score

			// If the score we get from the transposition table is a checkmate score, we need
			// to do a little extra work. This is because we store checkmates in the table using
			// their distance from the node they're found in, not their distance from the root.
			// So if we found a checkmate-in-8 in a node that was 5 plies from the root, we need
			// to store the score as a checkmate-in-3. Then, if we read the checkmate-in-3 from
			// the table in a node that's 4 plies from the root, we need to return the score as
			// checkmate-in-7.
			if score > Checkmate {
				score -= int16(ply)
			}

			if score < -Checkmate {
				score += int16(ply)
			}

			if entry.Flag == ExactFlag {
				// If we have an exact entry, we can use the saved score.
				adjustedScore = score
				shouldUse = true
			}

			if entry.Flag == AlphaFlag && score <= alpha {
				// If we have an alpha entry, and the entry's score is less than our
				// current alpha, then we know that our current alpha is the best score
				// we can get in this node, so we can stop searching and use alpha.
				adjustedScore = alpha
				shouldUse = true
			}

			if entry.Flag == BetaFlag && score >= beta {
				// If we have a beta entry, and the entry's score is greater than our
				// current beta, then we have a beta-cutoff, since while
				// searching this node previously, we found a value greater than the current
				// beta. so we can stop searching and use beta.
				adjustedScore = beta
				shouldUse = true
			}
		}
	}

	// Return the score
	return adjustedScore, shouldUse
}

func (entry *SearchEntry) Set(hash uint64, ply, depth uint8, score int16, flag uint8, best Move) {
	entry.Hash = hash
	entry.Depth = depth
	entry.Flag = flag
	entry.Best = best

	// If the score we get from the transposition table is a checkmate score, we need
	// to do a little extra work. This is because we store checkmates in the table using
	// their distance from the node they're found in, not their distance from the root.
	// So if we found a checkmate-in-8 in a node that was 5 plies from the root, we need
	// to store the score as a checkmate-in-3. Then, if we read the checkmate-in-3 from
	// the table in a node that's 4 plies from the root, we need to return the score as
	// checkmate-in-7.
	if score > Checkmate {
		score += int16(ply)
	}

	if score < -Checkmate {
		score -= int16(ply)
	}

	entry.Score = score
}

func (entry *PerftEntry) Get(hash uint64, depth uint8) (nodeCount uint64, ok bool) {
	if entry.Hash == hash && entry.Depth == depth {
		return entry.Nodes, true
	}
	return 0, false
}

func (entry *PerftEntry) Set(hash uint64, depth uint8, nodes uint64) {
	entry.Hash = hash
	entry.Depth = depth
	entry.Nodes = nodes
}

// A struct for a transposition table.
type TransTable[Entry SearchEntry | PerftEntry] struct {
	entries []Entry
	size    uint64
}

// Resize the transposition table given what the size should be in MB.
func (tt *TransTable[Entry]) Resize(sizeInMB uint64, entrySize uint64) {
	size := (sizeInMB * 1024 * 1024) / entrySize
	tt.entries = make([]Entry, size)
	tt.size = size
}

// Get an entry from the table.
func (tt *TransTable[Entry]) Probe(hash uint64) *Entry {
	// Get the entry from the table, calculating an index by modulo-ing the hash of
	// the position by the size of the table.
	return &tt.entries[hash%tt.size]
}

// Unitialize the memory used by the transposition table
func (tt *TransTable[Entry]) Unitialize() {
	tt.entries = nil
	tt.size = 0
}

// Clear the transposition table
func (tt *TransTable[Entry]) Clear() {
	for idx := uint64(0); idx < tt.size; idx++ {
		tt.entries[idx] = *new(Entry)
	}
}
