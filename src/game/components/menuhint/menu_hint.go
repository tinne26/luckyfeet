package menuhint

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"

func Draw(canvas *ebiten.Image, ctx *context.Context) {
	// draw [TAB] to enter menu
	white := color.RGBA{244, 244, 244, 244} // slightly translucid
	black := color.RGBA{ 16,  16,  16, 255}
	iconWidth := ctx.Gfxcore.MenuMask.Bounds().Dx()
	boxPad := 4
	offset := 5
	scale  := 2
	
	var txt []string
	if ctx.Input.UsedGamepadMoreRecentlyThanKeyboard() {
		txt = []string{ string(text.HalfSpace) + "STR" }
	} else {
		txt = []string{ string(text.HalfSpace) + "TAB" }
	}
	
	txtWidth := text.MeasureLineWidth(txt[0], scale)
	rect := image.Rect(
		640 - offset - txtWidth - iconWidth - boxPad*2 - scale*2,
		offset,
		640 - offset,
		offset + text.LineHeight*scale + boxPad*scale,
	)
	text.DrawRectBox(canvas, rect, 1, white, black, scale)
	text.RightDrawAt(canvas, 640 - offset - boxPad - scale, offset + boxPad, txt, black, scale)
	var opts ebiten.DrawImageOptions
	opts.GeoM.Translate(float64(640 - offset - boxPad - txtWidth - scale - iconWidth), float64(offset + boxPad + scale))
	opts.ColorScale.ScaleWithColor(black)
	canvas.DrawImage(ctx.Gfxcore.MenuMask, &opts)
}
