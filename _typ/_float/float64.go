package _float

func Float64Each(s []float64, fn func(seq int, elem float64)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
