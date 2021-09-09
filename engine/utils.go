package engine

// utils.go contains various utility functions used throughout the engine.

// Convert a string board coordinate to its position
// number.
func coordinateToPos(coordinate string) uint8 {
	file := coordinate[0] - 'a'
	rank := int(coordinate[1]-'0') - 1
	return uint8(rank*8 + int(file))
}

// Convert a position number to a string board coordinate.
func posToCoordinate(pos uint8) string {
	file := FileOf(pos)
	rank := RankOf(pos)
	return string(rune('a'+file)) + string(rune('0'+rank+1))
}

// Given a board square, return it's file.
func FileOf(sq uint8) uint8 {
	return sq % 8
}

// Given a board square, return it's rank.
func RankOf(sq uint8) uint8 {
	return sq / 8
}

// Get the absolute value of a number.
func abs16(n int16) int16 {
	if n < 0 {
		return -n
	}
	return n
}
