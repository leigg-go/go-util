package _int

func Int8Each(s []int8, fn func(seq int, elem int8)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
