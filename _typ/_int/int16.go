package _int

func Int16Each(s []int16, fn func(seq int, elem int16)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
