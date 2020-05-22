package _err

func Each(s []error, fn func(seq int, elem error)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
