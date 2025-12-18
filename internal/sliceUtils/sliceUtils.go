package sliceutils

func Remove[S ~[]E, E any](s S, i int) S {
	if i >= 0 && i < len(s) {
		s = append(s[:i], s[i+1:]...)
	}
	return s
}
