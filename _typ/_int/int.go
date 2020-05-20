package _int

func IntEach(s []int, fn func(seq int, elem int)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
