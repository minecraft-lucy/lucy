package tools

func NormalizedLevenshteinDistance(s1, s2 string) float64 {
	maxLen := max(len(s1), len(s2))
	if maxLen == 0 {
		return 0
	}
	distance := LevenshteinDistance(s1, s2)
	return float64(distance) / float64(maxLen)
}

func LevenshteinDistance(s1, s2 string) int {
	// m, n being the lengths of the two strings
	m := len(s1)
	n := len(s2)

	// Create distance matrix
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	// Matrix initialization
	for i := 0; i <= m; i++ {
		d[i][0] = i
	}

	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	// Compute Levenshtein distance
	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			d[i][j] = ternaryMin(
				d[i-1][j]+1,      // del
				d[i][j-1]+1,      // insert
				d[i-1][j-1]+cost, // substitute
			)
		}
	}

	return d[m][n]
}

func ternaryMin(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
