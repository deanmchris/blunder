package engine

// transposition.go contains an implementation of a transposition table (TT) to use
// in search.

const (
	// Default size of the transposition table, in MB.
	DefaultTTSize = 64

	// Constant for the size of a transposition table entry, in bytes,
	// considering memory alignment.
	TTEntrySize = 16

	// Constants representing the different flags for a transposition table entry,
	// which determine what kind of entry it is. If the entry has a score from
	// a fail-low node (alpha wasn't raised), it's an alpha entry. If the entry has
	// a score from a fail-high node (a beta cutoff occured), it's a beta entry. And
	// if the entry has an exact score (alpha was raised), it's an exact entry.
	AlphaFlag uint8 = iota
	BetaFlag
	ExactFlag

	// A constant representing an invalid score from probing the transposition table.
	// this constant's value doesn't matter as long as it's not in the range of possible
	// score values. So it must be outside of the range (-Inf, Inf).
	Invalid int16 = 20000

	// A constant representing the minimum or maximum value that a score from the search
	// must be below or above to be a checkmate score. The score assumes that the engine
	// will not find mate in 100.
	Checkmate = 9000
)

// A struct for a transposition table entry.
type TT_Entry struct {
	Hash  uint64
	Depth uint8
	Score int16
	Flag  uint8
	Best  Move
}

// A struct for a transposition table.
type TransTable struct {
	entries []TT_Entry
	size    uint64
}

// Resize the transposition table given what the size should be in MB.
func (tt *TransTable) Resize(sizeInMB uint64) {
	size := (sizeInMB * 1024 * 1024) / TTEntrySize
	tt.entries = make([]TT_Entry, size)
	tt.size = size
}

// Get an entry from the table.
func (tt *TransTable) Probe(hash uint64, ply, depth uint8, alpha, beta int16, best *Move) int16 {
	// Get the entry from the table, calculating an index by modulo-ing the hash of
	// the position by the size of the table.
	entry := tt.entries[hash%tt.size]

	adjustedScore := Invalid

	// Since index collisions can occur, test if the hash of the entry at this index
	// actually matches the hash for the current position.
	if entry.Hash == hash {

		// Even if we don't get a score we can use from the table, we can still
		// use the best move in this entry and put it first in our move ordering
		// scheme.
		*best = entry.Best

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
			}

			if entry.Flag == AlphaFlag && score <= alpha {
				// If we have an alpha entry, and the entry's score is less than our
				// current alpha, then we know that our current alpha is the best score
				// we can get in this node, so we can stop searching and use alpha.
				adjustedScore = alpha
			}

			if entry.Flag == BetaFlag && score >= beta {
				// If we have a beta entry, and the entry's score is greater than our
				// current beta, then we have a beta-cutoff, since while
				// searching this node previously, we found a value greater than the current
				// beta. so we can stop searching and use beta.
				adjustedScore = beta
			}
		}
	}

	// Return the score
	return adjustedScore
}

// Store an entry in the table.
func (tt *TransTable) Store(hash uint64, ply, depth uint8, score int16, flag uint8, best Move) {
	entry := &tt.entries[hash%tt.size]
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

// Unitialize the memory used by the transposition table
func (tt *TransTable) Unitialize() {
	tt.entries = nil
	tt.size = 0
}

// Clear the transposition table
func (tt *TransTable) Clear() {
	var idx uint64
	for idx = 0; idx < tt.size; idx++ {
		tt.entries[idx] = TT_Entry{}
	}
}
