package text

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

var BackColor  = color.RGBA{ 20,  20,  20, 255}
var FrontColor = color.RGBA{255, 255, 255, 255}

const LineHeight = 7
const LineInterspace = 2
const BoxVertMargin = 3
const BoxHorzMargin = 10
const BoxBorderOffset = 1

// for long text and so on
func CenterDraw(canvas *ebiten.Image, lines []string, clr color.RGBA, scale int) {
	bounds := canvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	CenterDrawAt(canvas, bounds.Min.X + w/2, bounds.Min.Y + h/2, lines, clr, scale)
}

func CenterDrawAt(canvas *ebiten.Image, x, y int, lines []string, clr color.RGBA, scale int) {
	textHeight := len(lines)*LineHeight*scale + (len(lines) - 1)*LineInterspace*scale
	y -= textHeight/2

	for _, line := range lines {
		lineWidth := MeasureLineWidth(line, scale)
		DrawLine(canvas, line, x - lineWidth/2, y, clr, scale)
		y += LineHeight*scale + LineInterspace*scale
	}
}

func RightDrawAt(canvas *ebiten.Image, x, y int, lines []string, clr color.RGBA, scale int) {
	for _, line := range lines {
		lineWidth := MeasureLineWidth(line, scale)
		DrawLine(canvas, line, x - lineWidth, y, clr, scale)
		y += LineHeight*scale + LineInterspace*scale
	}
}

func DrawAt(canvas *ebiten.Image, x, y int, lines []string, clr color.RGBA, scale int) {
	for _, line := range lines {
		DrawLine(canvas, line, x, y, clr, scale)
		y += LineHeight*scale + LineInterspace*scale
	}
}

func drawBoxAt(canvas *ebiten.Image, x, y int, fillClr, borderClr color.RGBA, scale int, textHeight, maxWidth int) {
	rect := image.Rect(x, y, x + maxWidth + BoxHorzMargin*scale*2, y + textHeight + BoxVertMargin*scale*2)
	DrawRectBox(canvas, rect, BoxBorderOffset, fillClr, borderClr, scale)
}

func DrawRectBox(canvas *ebiten.Image, rect image.Rectangle, borderOffset int, fillClr, borderClr color.RGBA, scale int) {
	borderOffset *= scale

	// draw main fill
	fillOver(canvas, rect, fillClr)

	// draw horz borders
	x := rect.Min.X + borderOffset
	y := rect.Min.Y + borderOffset
	subRect := image.Rect(x, y, rect.Max.X - borderOffset, y + scale)
	fillOver(canvas, subRect, borderClr)
	fy := rect.Max.Y - borderOffset - scale
	subRect = image.Rect(x, fy, rect.Max.X - borderOffset, fy + scale)
	fillOver(canvas, subRect, borderClr)

	// draw vert borders
	y += scale
	subRect = image.Rect(x, y, x + scale, fy)
	fillOver(canvas, subRect, borderClr)
	x = rect.Max.X - borderOffset - scale
	subRect = image.Rect(x, y, x + scale, fy)
	fillOver(canvas, subRect, borderClr)
}

func DrawBoxAt(canvas *ebiten.Image, x, y int, lines []string, fillClr, borderClr color.RGBA, scale int) {
	textHeight := len(lines)*LineHeight*scale + (len(lines) - 1)*LineInterspace*scale
	var maxWidth int
	for _, line := range lines {
		maxWidth = max(MeasureLineWidth(line, scale), maxWidth)
	}

	drawBoxAt(canvas, x, y, fillClr, borderClr, scale, textHeight, maxWidth)
}

func DrawCenteredBoxAt(canvas *ebiten.Image, x, y int, lines []string, fillClr, borderClr color.RGBA, scale int) {
	textHeight := len(lines)*LineHeight*scale + (len(lines) - 1)*LineInterspace*scale
	y -= textHeight/2
	y -= BoxVertMargin*scale

	var maxWidth int
	for _, line := range lines {
		maxWidth = max(MeasureLineWidth(line, scale), maxWidth)
	}
	x -= maxWidth/2
	x -= BoxHorzMargin*scale

	drawBoxAt(canvas, x, y, fillClr, borderClr, scale, textHeight, maxWidth)
}

func DrawRightBoxAt(canvas *ebiten.Image, x, y int, lines []string, fillClr, borderClr color.RGBA, scale int) {
	textHeight := len(lines)*LineHeight*scale + (len(lines) - 1)*LineInterspace*scale
	var maxWidth int
	for _, line := range lines {
		maxWidth = max(MeasureLineWidth(line, scale), maxWidth)
	}
	x -= maxWidth
	x -= BoxHorzMargin*2*scale

	drawBoxAt(canvas, x, y, fillClr, borderClr, scale, textHeight, maxWidth)
}

func MeasureLineWidth(line string, scale int) int {
	var prevIsSpace bool
	width := 0
	for i, codePoint := range line {
		if codePoint == ' ' {
			width += 4*scale
			prevIsSpace = true
		} else {
			if i != 0 && !prevIsSpace { width += 1*scale }
			prevIsSpace = false
			bmp, found := pkgBitmaps[codePoint]
			if !found { panic("missing bitmap for '" + string(codePoint) + "'") }
			width += bmp.Bounds().Dx()*scale
		}
	}
	return width
}

func DrawLine(canvas *ebiten.Image, line string, ox, oy int, textColor color.RGBA, scale int) {
	var prevIsSpace bool
	x := 0
	opts := ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleWithColor(textColor)
	for i, codePoint := range line {
		if codePoint == ' ' {
			x += 4*scale
			prevIsSpace = true
		} else {
			if i != 0 && !prevIsSpace { x += 1*scale }
			prevIsSpace = false
			img, found := pkgBitmaps[codePoint]
			if !found { panic("missing bitmap for glyph '" + string(codePoint) + "'") }
			opts.GeoM.Scale(float64(scale), float64(scale))
			opts.GeoM.Translate(float64(ox + x), float64(oy))
			canvas.DrawImage(img, &opts)
			opts.GeoM.Reset()
			x += img.Bounds().Dx()*scale
		}
	}
}
