package briefblackout

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/scene"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/utils"

var _ scene.Scene[*context.Context] = (*BriefBlackout)(nil)

// Look, there are some things it's better not to ask about.

type BriefBlackout struct {
	ticksLeft int
}

func New(ctx *context.Context) (*BriefBlackout, error) {
	return &BriefBlackout{50}, nil
}

func (self *BriefBlackout) Update(ctx *context.Context) (*scene.Change, error) {
	self.ticksLeft -= 1
	if self.ticksLeft <= 0 { return scene.Pop(), nil }
	return nil, nil
}

func (self *BriefBlackout) DrawLogical(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	utils.FillOver(canvas, color.RGBA{0, 0, 0, 255})
}

func (self *BriefBlackout) DrawHiRes(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	// ...
}
