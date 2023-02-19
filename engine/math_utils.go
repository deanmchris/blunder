package engine

import (
	"golang.org/x/exp/constraints"
)

// Get the absolute value of an integer.
func abs[Int constraints.Integer](n Int) Int {
	if n < 0 {
		return -n
	}
	return n
}
