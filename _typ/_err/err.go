package _err

func ErrEach(s []error, fn func(seq int, elem error)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
