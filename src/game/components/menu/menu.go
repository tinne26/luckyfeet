package menu

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/settings"

type Key uint8
const (
	NoChange Key = iota
	NoConfirm
	Back
	Options
	OptsAudio
	OptsWindow
	OptsScaling
	FirstKey
)

type Menu struct {
	Title string
	key Key
	options map[Key]*OptionList
}

func (self *Menu) JumpTo(key Key) {
	self.key = key
	self.options[key].index = 0
	// TODO: would probably need to adjust the 'from' key on the target option list too
}

func (self *Menu) Update(ctx *context.Context) (*scene.Change, error) {
	key, change, err := self.options[self.key].Update(ctx)
	if key != NoChange {
		if key == Back { panic("OptionList can't return 'Back' key directly") }
		optList, found := self.options[key]
		if !found { panic(key) }
		if self.options[self.key].from != key {
			if self.key != key { optList.from = self.key }
			optList.index = 0
		}
		self.key = key
	}
	return change, err
}

func (self *Menu) DrawLogical(canvas *ebiten.Image, ctx *context.Context) {
	if self.Title != "" {
		scale := 3
		text.CenterDrawAt(canvas, 320, 124        , []string{self.Title}, color.RGBA{64, 64, 64, 255}, scale)
		text.CenterDrawAt(canvas, 320, 124 - scale, []string{self.Title}, color.RGBA{255, 255, 255, 255}, scale)
	}
	self.options[self.key].Draw(canvas, ctx)
}

func (self *Menu) NewOptionList(key Key) *OptionList {
	_, alreadyDefined := self.options[key]
	if alreadyDefined {
		panic("option list already defined for the given key")
	}
	optList := &OptionList{}
	if self.options == nil {
		self.options = make(map[Key]*OptionList, 4)
	}
	self.options[key] = optList
	return optList
}

func (self *Menu) NewGameOptionsOptionList(ctx *context.Context) {
	opts := self.NewOptionList(Options)
	opts.Add(&NavOption{ Label: "AUDIO", To: OptsAudio })
	opts.Add(&NavOption{ Label: "WINDOW", To: OptsWindow })
	opts.Add(&NavOption{ Label: "SCALING", To: OptsScaling })
	opts.AddBackOption(&NavOption{ Label: "BACK", To: Back })

	opts = self.NewOptionList(OptsAudio)
	opts.Add(&AudioOption{
		BaseLabel: "MUSIC",
		GetLevel: func() float32 { return ctx.Audio.GetUserBGMVolume() },
		SetLevel: func(value float32) error {
			ctx.Audio.SetUserBGMVolume(value)
			return nil
		}})
	opts.Add(&AudioOption{
		BaseLabel: "SFX",
		GetLevel: func() float32 { return ctx.Audio.GetUserSFXVolume() },
		SetLevel: func(value float32) error {
			ctx.Audio.SetUserSFXVolume(value)
			return nil
		}})
	opts.AddBackOption(&NavOption{ Label: "BACK", To: Back })

	opts = self.NewOptionList(OptsWindow)
	opts.Add(&EffectOptionWithHighlight{
		Label: "FULLSCREEN",
		OnConfirm: func(*context.Context) error {
			ebiten.SetFullscreen(true)
			return nil
		},
		HighlightFunc: func(*context.Context) bool {
			return ebiten.IsFullscreen()
		},
	})
	opts.Add(&EffectOptionWithHighlight{
		Label: "WINDOWED",
		OnConfirm: func(*context.Context) error {
			ebiten.SetFullscreen(false)
			return nil
		},
		HighlightFunc: func(*context.Context) bool {
			return !ebiten.IsFullscreen()
		},
	})
	opts.AddBackOption(&NavOption{ Label: "BACK", To: Back })

	opts = self.NewOptionList(OptsScaling)
	opts.Add(&EffectOptionWithHighlight{
		Label: "PIXEL PERFECT",
		OnConfirm: func(ctx *context.Context) error {
			ctx.Settings.ScreenFit = settings.ScreenFitPixelPerfect
			return nil
		},
		HighlightFunc: func(ctx *context.Context) bool {
			return ctx.Settings.ScreenFit == settings.ScreenFitPixelPerfect
		},
	})
	opts.Add(&EffectOptionWithHighlight{
		Label: "PROPORTIONAL",
		OnConfirm: func(ctx *context.Context) error {
			ctx.Settings.ScreenFit = settings.ScreenFitProportional
			return nil
		},
		HighlightFunc: func(ctx *context.Context) bool {
			return ctx.Settings.ScreenFit == settings.ScreenFitProportional
		},
	})
	opts.Add(&EffectOptionWithHighlight{
		Label: "STRETCHED",
		OnConfirm: func(ctx *context.Context) error {
			ctx.Settings.ScreenFit = settings.ScreenFitStretch
			return nil
		},
		HighlightFunc: func(ctx *context.Context) bool {
			return ctx.Settings.ScreenFit == settings.ScreenFitStretch
		},
	})
	opts.AddBackOption(&NavOption{ Label: "BACK", To: Back })
}
