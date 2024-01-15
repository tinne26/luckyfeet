package utils

import "image"

func Rect(ox, oy, fx, fy int) image.Rectangle {
	return image.Rect(ox, oy, fx, fy)
}
