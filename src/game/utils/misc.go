package utils

import "image"

func Rect(ox, oy, fx, fy int) image.Rectangle {
	return image.Rect(ox, oy, fx, fy)
}

func ClipToLastN[T any](slice []T, n int) []T {
	if len(slice) <= n { return slice }
	copy(slice[ : n], slice[len(slice) - n : ])
	return slice[ : n]
}

func EqualRunes(slice []rune, target string) bool {
	var i int
	for _, codePoint := range target {
		if i >= len(slice) { return false }
		if codePoint != slice[i] { return false }
		i += 1
	}
	return true
}
