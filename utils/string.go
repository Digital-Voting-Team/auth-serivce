package utils

func Hint(s string, size int) string {
	if len([]rune(s)) < size {
		return s
	}
	return string([]rune(s)[0:size])
}
