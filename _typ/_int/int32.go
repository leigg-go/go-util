package _int

/*
int is a alias of int32
*/

func Int32Each(s []int32, fn func(seq int, elem int32)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}

func IntEach(s []int, fn func(seq int, elem int)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
