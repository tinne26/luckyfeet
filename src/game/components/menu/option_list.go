package menu

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/input"
import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"

type OptionList struct {
	from Key // for "Back" operations
	options []Option
	backOption Option
	index int
}

func (self *OptionList) Add(opt Option) {
	self.options = append(self.options, opt)
}

func (self *OptionList) AddBackOption(opt Option) {
	if self.backOption != nil { panic("can't have multiple back options") }
	self.Add(opt)
	self.backOption = opt
}

func (self *OptionList) Update(ctx *context.Context) (Key, *scene.Change, error) {
	// handle back action
	if self.backOption != nil {
		if ctx.Input.Trigger(in.ActionBack) {
			ctx.Audio.PlaySFX(au.SfxBack)
			key, change, err := self.backOption.Confirm(ctx)
			if key == Back { key = self.from }
			return key, change, err
		}
	}

	// handle confirmation
	var key Key = NoChange
	var change *scene.Change
	var err error
	if ctx.Input.Trigger(in.ActionConfirm) {
		key, change, err = self.options[self.index].Confirm(ctx)
		if key == NoConfirm {
			key = NoChange
		} else {
			// play sound effect
			ctx.Audio.PlaySFX(au.SfxConfirm)
		}
		if key == Back { key = self.from }
		return key, change, err
	}

	// hover update
	self.options[self.index].HoverUpdate(ctx)

	// handle regular menu navigation
	dir := ctx.Input.RepeatDirAs(in.RFDefault, in.RNSlow)
	if dir == in.DirUp {
		self.index -= 1
		if self.index < 0 {
			self.index = len(self.options) - 1
		}
		ctx.Audio.PlaySFX(au.SfxClick)
	} else if dir == in.DirDown {
		self.index += 1
		if self.index >= len(self.options) {
			self.index = 0
		}
		ctx.Audio.PlaySFX(au.SfxClick)
	}

	return NoChange, nil, nil
}

func (self *OptionList) navigate(inputHandler input.KBGPLike) bool {
	dir := inputHandler.RepeatDirAs(in.RFDefault, in.RNSlow)
	if dir == in.DirUp {
		self.index -= 1
		if self.index < 0 {
			self.index = len(self.options) - 1
		}
		return true
	} else if dir == in.DirDown {
		self.index += 1
		if self.index >= len(self.options) {
			self.index = 0
		}
		return true
	}
	return false
}

func (self *OptionList) Draw(canvas *ebiten.Image, ctx *context.Context) {
	white  := color.RGBA{244, 244, 244, 244} // slightly translucid
	black  := color.RGBA{ 16,  16,  16, 255}
	orange := color.RGBA{255, 71, 20, 255}
	highlt := color.RGBA{250, 212, 226, 250}
	
	var longestOptName string
	var longestOptNameWidth int
	for _, opt := range self.options {
		name := opt.Name()
		maxNameOpt, hasMaxName := opt.(OptionWithMaxName)
		if hasMaxName { name = maxNameOpt.MaxName() }
		width := text.MeasureLineWidth(name, 1)
		if width > longestOptNameWidth {
			longestOptName = name
			longestOptNameWidth = width
		}
	}

	x := canvas.Bounds().Dx()/2
	y := 177 // hardcoding to death, of course
	for i, opt := range self.options {
		lineClr, fillClr := black, white
		if i == self.index { lineClr = orange }
		if opt.SoftHighlight(ctx) { fillClr = highlt }
		text.DrawCenteredBoxAt(canvas, x, y, []string{ longestOptName }, fillClr, lineClr, 2)
		text.CenterDrawAt(canvas, x, y, []string{ opt.Name() }, lineClr, 2)
		y += 34
	}
}
