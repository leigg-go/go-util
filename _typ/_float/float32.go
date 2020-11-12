package _float

func Float32Each(s []float32, fn func(seq int, elem float32)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
