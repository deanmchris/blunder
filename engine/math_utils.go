package engine

import (
	"golang.org/x/exp/constraints"
)

// Get the absolute value of an integer.
func Abs[Real constraints.Integer | constraints.Float](n Real) Real {
	if n < 0 {
		return -n
	}
	return n
}
