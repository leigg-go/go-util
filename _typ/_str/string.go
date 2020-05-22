package _str

func SetDefault(s *string, new string, def string) {
	if new != "" {
		*s = new
	} else {
		*s = def
	}
}

func If(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
}
