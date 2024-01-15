package utils

import "image/color"

func WithAlpha(clr color.RGBA, newAlpha uint8) color.RGBA {
	scaling := float64(newAlpha)/float64(clr.A)
	return color.RGBA{
		R: uint8(float64(clr.R)*scaling),
		G: uint8(float64(clr.G)*scaling),
		B: uint8(float64(clr.B)*scaling),
		A: newAlpha,
	}
}

func RGBAf64(r, g, b, a float64) color.RGBA {
	if a < r || a < g || a < b { panic("expected alpha premultiplied color") }
	if r < 0 || g < 0 || b < 0 { panic("color channel must be in [0, 1]") }
	if r > 1 || g > 1 || b > 1 || a > 1 { panic("color channel must be in [0, 1]") }
	return color.RGBA{uint8(r*255), uint8(g*255), uint8(b*255), uint8(a*255)}
}
