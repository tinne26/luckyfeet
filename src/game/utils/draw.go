package utils

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

var vertices []ebiten.Vertex
var mask1x1 *ebiten.Image
var stdTriOpts ebiten.DrawTrianglesOptions

func init() {
	mask3x3 := ebiten.NewImage(3, 3)
	mask3x3.Fill(color.White)
	mask1x1 = mask3x3.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	vertices = make([]ebiten.Vertex, 4)
	for i := 0; i < 4; i++ {
		vertices[i].SrcX = 1.0
		vertices[i].SrcY = 1.0
	}
}

func FillOver(target *ebiten.Image, fillColor color.Color) {
	FillOverRect(target, target.Bounds(), fillColor)
}

func FillOverRect(target *ebiten.Image, bounds image.Rectangle, fillColor color.Color) {
	if bounds.Empty() { return }

	r, g, b, a := fillColor.RGBA()
	if a == 0 { return }
	fr, fg, fb, fa := float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535
	for i := 0; i < 4; i++ {
		vertices[i].ColorR = fr
		vertices[i].ColorG = fg
		vertices[i].ColorB = fb
		vertices[i].ColorA = fa
	}

	minX, minY := float32(bounds.Min.X), float32(bounds.Min.Y)
	maxX, maxY := float32(bounds.Max.X), float32(bounds.Max.Y)
	vertices[0].DstX = minX
	vertices[0].DstY = minY
	vertices[1].DstX = maxX
	vertices[1].DstY = minY
	vertices[2].DstX = maxX
	vertices[2].DstY = maxY
	vertices[3].DstX = minX
	vertices[3].DstY = maxY

	target.DrawTriangles(vertices[0 : 4], []uint16{0, 1, 2, 2, 3, 0}, mask1x1, &stdTriOpts)
}

func FillOverRectLighter(target *ebiten.Image, bounds image.Rectangle, fillColor color.Color) {
	if bounds.Empty() { return }

	r, g, b, a := fillColor.RGBA()
	if a == 0 { return }
	fr, fg, fb, fa := float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535
	for i := 0; i < 4; i++ {
		vertices[i].ColorR = fr
		vertices[i].ColorG = fg
		vertices[i].ColorB = fb
		vertices[i].ColorA = fa
	}

	minX, minY := float32(bounds.Min.X), float32(bounds.Min.Y)
	maxX, maxY := float32(bounds.Max.X), float32(bounds.Max.Y)
	vertices[0].DstX = minX
	vertices[0].DstY = minY
	vertices[1].DstX = maxX
	vertices[1].DstY = minY
	vertices[2].DstX = maxX
	vertices[2].DstY = maxY
	vertices[3].DstX = minX
	vertices[3].DstY = maxY

	stdTriOpts.Blend = ebiten.BlendLighter
	target.DrawTriangles(vertices[0 : 4], []uint16{0, 1, 2, 2, 3, 0}, mask1x1, &stdTriOpts)
	stdTriOpts.Blend = ebiten.BlendSourceOver
}
