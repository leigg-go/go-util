package _int

func Int32Each(s []int32, fn func(seq int, elem int32)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
