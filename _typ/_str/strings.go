package _str

func Each(s []string, fn func(i int, elem string)) {
	for i, elem := range s {
		fn(i, elem)
	}
}

func Find(ss []string, s string) int {
	idx := -1 // not found
	Each(ss, func(i int, elem string) {
		if elem == s {
			idx = i
		}
	})
	return idx
}

func Contains(ss []string, s string) bool {
	return Find(ss, s) != -1
}
