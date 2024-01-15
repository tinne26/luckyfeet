package menu

import "strconv"

import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"

type Option interface {
	Name() string
	HoverUpdate(*context.Context) // for extra stuff on hover
	Confirm(*context.Context) (Key, *scene.Change, error)
	SoftHighlight(*context.Context) bool
}

type OptionWithMaxName interface {
	Option
	MaxName() string
}

type NavOption struct {
	Label string
	To Key
}
func (self *NavOption) Name() string { return self.Label }
func (self *NavOption) HoverUpdate(ctx *context.Context) {}
func (self *NavOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *NavOption) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return self.To, nil, nil
}

type SceneChangeOption struct {
	Label string
	Change scene.Change
}
func (self *SceneChangeOption) Name() string { return self.Label }
func (self *SceneChangeOption) HoverUpdate(ctx *context.Context) {}
func (self *SceneChangeOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *SceneChangeOption) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return NoChange, &self.Change, nil
}

type SceneChangeEffectOption struct {
	Label string
	Change scene.Change
	OnConfirm func(*context.Context) error
}
func (self *SceneChangeEffectOption) Name() string { return self.Label }
func (self *SceneChangeEffectOption) HoverUpdate(ctx *context.Context) {}
func (self *SceneChangeEffectOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *SceneChangeEffectOption) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return NoChange, &self.Change, self.OnConfirm(ctx)
}

type BasicOption struct {
	Label string
	Func func(*context.Context) Key
}
func (self *BasicOption) Name() string { return self.Label }
func (self *BasicOption) HoverUpdate(ctx *context.Context) {}
func (self *BasicOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *BasicOption) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return self.Func(ctx), nil, nil
}

type EffectOption struct {
	Label string
	OnConfirm func(*context.Context) error
}
func (self *EffectOption) Name() string { return self.Label }
func (self *EffectOption) HoverUpdate(ctx *context.Context) {}
func (self *EffectOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *EffectOption) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return NoChange, nil, self.OnConfirm(ctx)
}

type EffectOptionWithHighlight struct {
	Label string
	OnConfirm func(*context.Context) error
	HighlightFunc func(*context.Context) bool
}
func (self *EffectOptionWithHighlight) Name() string { return self.Label }
func (self *EffectOptionWithHighlight) HoverUpdate(ctx *context.Context) {}
func (self *EffectOptionWithHighlight) SoftHighlight(ctx *context.Context) bool {
	return self.HighlightFunc(ctx)
}
func (self *EffectOptionWithHighlight) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return NoChange, nil, self.OnConfirm(ctx)
}

type PercentOption struct {
	BaseLabel string
	NotifyPercent func(uint8)
	Percent uint8
}
func (self *PercentOption) Name() string {
	perc := strconv.Itoa(int(self.Percent))
	return self.BaseLabel + " " + string(text.TriangleLeftWithPad) + perc + string(text.TriangleRightWithPad)
}
func (self *PercentOption) MaxName() string {
	return self.BaseLabel + " " + string(text.TriangleLeftWithPad) + "100" + string(text.TriangleRightWithPad)
}
func (self *PercentOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *PercentOption) HoverUpdate(ctx *context.Context) {
	dir := ctx.Input.RepeatDirAs(in.RFDefault, in.RNDefault)
	if dir == in.DirRight {
		if self.Percent < 100 { self.Percent += 1 }
		self.NotifyPercent(self.Percent)
		ctx.Audio.PlaySFX(au.SfxClick)
	} else if dir == in.DirLeft {
		if self.Percent > 0 { self.Percent -= 1 }
		self.NotifyPercent(self.Percent)
		ctx.Audio.PlaySFX(au.SfxClick)
	}
}
func (self *PercentOption) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return NoConfirm, nil, nil
}

type AudioOption struct {
	BaseLabel string
	GetLevel func() float32 // in [0, 1] range
	SetLevel func(value float32) error
}
func (self *AudioOption) Name() string {
	perc := strconv.Itoa(int(max(min(self.GetLevel(), 1.0), 0.0)*100))
	return self.BaseLabel + " " + string(text.TriangleLeftWithPad) + perc + string(text.TriangleRightWithPad)
}
func (self *AudioOption) MaxName() string {
	return self.BaseLabel + " " + string(text.TriangleLeftWithPad) + "100" + string(text.TriangleRightWithPad)
}
func (self *AudioOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *AudioOption) HoverUpdate(ctx *context.Context) {
	dir := ctx.Input.RepeatDirAs(in.RFDefault, in.RNDefault).Horz()
	if dir == in.DirNone { return }

	level := self.GetLevel()
	if dir == in.DirRight {
		level += 0.01
		if level > 1.0 { level = 1.0 }
	} else if dir == in.DirLeft {
		level -= 0.01
		if level < 0 { level = 0 }
	}
	ctx.Audio.PlaySFX(au.SfxClick)
	self.SetLevel(level)
}
func (self *AudioOption) Confirm(ctx *context.Context) (Key, *scene.Change, error) {
	return NoConfirm, nil, nil
}