package _int

func Int64Each(s []int64, fn func(seq int, elem int64)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
