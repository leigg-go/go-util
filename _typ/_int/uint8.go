package _int

/*
uint is a alias of uint32
*/

func UInt8Each(s []uint8, fn func(seq int, elem uint8)) {
	for seq, elem := range s {
		fn(seq, elem)
	}
}
