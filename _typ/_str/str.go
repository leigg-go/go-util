package _str

func Each(s []string, fn func(i int, elem string)) {
	for i, elem := range s {
		fn(i, elem)
	}
}
