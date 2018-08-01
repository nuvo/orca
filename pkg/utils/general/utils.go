package genutils

// AddIfNotContained adds a string to a slice if it is not contained in it and not empty
func AddIfNotContained(s []string, e string) (sout []string) {
	if (!contains(s, e)) && (e != "") {
		s = append(s, e)
	}

	return s
}

// contains checks if a slice contains a given value
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
