package _int

/*
uint is a alias of uint32
*/

func UInt32Each(s []uint32, fn func(seq int, elem uint32)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}

func UIntEach(s []uint, fn func(seq int, elem uint)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
