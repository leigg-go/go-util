package _int

func UInt64Each(s []uint64, fn func(seq int, elem uint64)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
