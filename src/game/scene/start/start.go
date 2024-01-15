package start

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/components/back"
import "github.com/tinne26/luckyfeet/src/game/components/menu"
import "github.com/tinne26/luckyfeet/src/game/components/info"
import "github.com/tinne26/luckyfeet/src/game/components/tile"
import "github.com/tinne26/luckyfeet/src/game/material/scene/keys"

var _ scene.Scene[*context.Context] = (*Start)(nil)

type Start struct {
	credits *info.Layer
	controls *info.Layer
	menu menu.Menu
	backMap *tile.Map
}

const (
	keyMainMenu menu.Key = menu.FirstKey
	keyEditor   menu.Key = menu.FirstKey + 1
	keyWonder   menu.Key = menu.FirstKey + 2
)

func New(ctx *context.Context) (*Start, error) {
	if ctx.Background == nil {
		ctx.Background = back.New()
	}

	credits := info.New([]string{
		"CODE, ART AND MUSIC BY TINNE",
		"FOR THE \"EBITENGINE HOLIDAY HACK 2023\"",
		"",
		"ALL CODE IS OPEN SOURCE",
		"ALL ASSETS AVAILABLE UNDER CC LICENSES",
		"GITHUB.COM/TINNE26/LUCKYFEET",
	})
	controls := info.New(info.ControlsKB)
	
	var mainMenu menu.Menu
	opts := mainMenu.NewOptionList(keyMainMenu)
	opts.Add(&menu.SceneChangeOption{ Label: "JUMP", Change: *scene.PushTo(keys.Play) })
	opts.Add(&menu.NavOption{ Label: "DIG", To: keyEditor })
	opts.Add(&menu.NavOption{ Label: "WONDER", To: keyWonder })
	opts = mainMenu.NewOptionList(keyWonder)
	opts.Add(&menu.NavOption{ Label: "OPTIONS", To: menu.Options })
	opts.Add(controls.NewOption("CONTROLS"))
	opts.Add(credits.NewOption("CREDITS"))
	opts.AddBackOption(&menu.NavOption{ Label: "BACK", To: keyMainMenu })
	opts = mainMenu.NewOptionList(keyEditor)
	opts.Add(&menu.SceneChangeOption{ Label: "NEW PROJECT", Change: *scene.PushTo(keys.Editor) })
	opts.Add(&menu.SceneChangeEffectOption{
		Label: "EDIT FROM CLIPBOARD",
		Change: *scene.PushTo(keys.Editor),
		OnConfirm: func(fnCtx *context.Context) error {
			fnCtx.State.LoadMapDataFromClipboard = true
			return nil
		}})
		opts.Add(&menu.SceneChangeEffectOption{
			Label: "PLAY FROM CLIPBOARD",
			Change: *scene.PushTo(keys.Play),
			OnConfirm: func(fnCtx *context.Context) error {
				fnCtx.State.LoadMapDataFromClipboard = true
				return nil
			}})
	opts.AddBackOption(&menu.NavOption{ Label: "BACK", To: menu.Back })
	mainMenu.NewGameOptionsOptionList(ctx)
	mainMenu.JumpTo(keyMainMenu)

	var backMap *tile.Map = tile.NewMap(1)
	if backMapData != "" {
		var err error
		backMap, err = tile.LoadMapFromString(backMapData)
		if err != nil { return nil, err }
	}
	return &Start{
		credits: credits,
		controls: controls,
		menu: mainMenu,
		backMap: backMap,
	}, nil
}

func (self *Start) Update(ctx *context.Context) (*scene.Change, error) {
	if ctx.Scenes.Current() != self { return nil, nil }

	// update background animation
	ctx.Background.Update()
	
	// update menu or info layers
	if self.credits.IsVisible() {
		self.credits.Update(ctx)
	} else if self.controls.IsVisible() {
		self.controls.Update(ctx)
	} else {
		change, err := self.menu.Update(ctx)
		if change != nil { self.menu.JumpTo(keyMainMenu) }
		return change, err
	}

	return nil, nil
}

func (self *Start) DrawLogical(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	if !foremost { return }
	ctx.State.Editing = false

	ctx.Background.DrawLogical(canvas, ctx)
	self.backMap.DrawBackLogical(canvas, ctx, nil)
	self.backMap.DrawMainLogical(canvas, ctx, nil)
	self.backMap.DrawFrontLogical(canvas, ctx, nil)
	
	// draw game title (white part behind, black part in front)
	white := color.RGBA{244, 244, 244, 244} // slightly translucid
	black := color.RGBA{ 16,  16,  16, 255}
	bounds := canvas.Bounds()
	x := bounds.Dx()/2
	y := bounds.Dy()/4
	strs := []string{"LUCKY FEET"}
	text.CenterDrawAt(canvas, x, y - 4, strs, white, 4)
	text.CenterDrawAt(canvas, x, y - 0, strs, black, 4)
	
	// draw menu or info layer
	if self.credits.IsVisible() {
		self.credits.Draw(canvas)
	} else if self.controls.IsVisible() {
		if ctx.Input.UsedGamepadMoreRecentlyThanKeyboard() {
			self.controls.SetContent(info.ControlsGP)
		} else {
			self.controls.SetContent(info.ControlsKB)
		}
		self.controls.Draw(canvas)
	} else {
		self.menu.DrawLogical(canvas, ctx)
	}
}

func (self *Start) DrawHiRes(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	// ...
}
