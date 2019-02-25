package util

// FilterString removes any instances of the needle from haystack.  It
// returns a new slice with all instances of needle removed, and a
// count of the number instances encountered.
func FilterString(haystack []string, needle string) ([]string, int) {
	newSlice := haystack[:0] // Share the backing array.
	found := 0

	for _, x := range haystack {
		if x != needle {
			newSlice = append(newSlice, x)
		} else {
			found++
		}
	}

	return newSlice, found
}
