package info

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/game/components/menu"
import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/utils"

type Layer struct {
	content []string
	visible bool
}

func New(content []string) *Layer {
	return &Layer{
		content: content,
		visible: false,
	}
}

func (self *Layer) SetContent(content []string) {
	self.content = content
}

func (self *Layer) IsVisible() bool {
	return self.visible
}

func (self *Layer) Show() {
	self.visible = true
}

func (self *Layer) Update(ctx *context.Context) {
	if !self.visible { return }
	if ctx.Input.Trigger(in.ActionBack) || ctx.Input.Trigger(in.ActionConfirm) {
		self.visible = false
		ctx.Audio.PlaySFX(au.SfxBack)
	}
}

func (self *Layer) Draw(canvas *ebiten.Image) {
	utils.FillOver(canvas, color.RGBA{0, 0, 0, 138})
	text.CenterDraw(canvas, self.content, color.RGBA{255, 255, 255, 255}, 2)
}

func (self *Layer) NewOption(label string) menu.Option {
	return &menu.EffectOption{
		Label: label,
		OnConfirm: self.optButtonOnConfirm,
	}
}

func (self *Layer) optButtonOnConfirm(*context.Context) error {
	self.visible = true
	return nil
}
