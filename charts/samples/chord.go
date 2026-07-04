package samples

// Chord returns a small directed flow matrix between five cities and the
// matching key labels, ready to assign to ChordProps.Data and ChordProps.Keys.
func Chord() (matrix [][]float64, keys []string) {
	keys = []string{"Tokyo", "Osaka", "Kyoto", "Nagoya", "Sapporo"}
	matrix = [][]float64{
		{0, 15834, 6987, 4211, 1893},
		{12345, 0, 5432, 3210, 987},
		{5678, 4321, 0, 2109, 654},
		{3456, 2345, 1876, 0, 432},
		{1234, 876, 543, 321, 0},
	}
	return matrix, keys
}
