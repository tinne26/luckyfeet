package racetimer

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"

func Draw(canvas *ebiten.Image, ctx *context.Context, ticks int) {
	white := color.RGBA{244, 244, 244, 244} // slightly translucid
	black := color.RGBA{ 16,  16,  16, 255}
	scale  := 2
	
	pad := 4
	vertBoxOffset := 5
	horzBoxOffset := 8

	var txt []string = []string{ ctx.State.FmtTicksToTimeStr(ticks) }
	txtWidth := text.MeasureLineWidth(txt[0], scale)
	rect := image.Rect(pad, pad, pad + txtWidth + horzBoxOffset*2, pad + text.LineHeight*scale + vertBoxOffset*2)
	text.DrawRectBox(canvas, rect, 1, white, black, scale)
	text.DrawAt(canvas, pad + horzBoxOffset, pad + vertBoxOffset, txt, black, scale)
}
