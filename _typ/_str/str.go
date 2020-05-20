package _str

func StrEach(s []string, fn func(seq int, elem string)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
