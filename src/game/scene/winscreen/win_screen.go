package winscreen

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/utils"
import "github.com/tinne26/luckyfeet/src/game/components/info"
import "github.com/tinne26/luckyfeet/src/game/components/menu"
import "github.com/tinne26/luckyfeet/src/game/material/scene/keys"

var _ scene.Scene[*context.Context] = (*WinScreen)(nil)

type WinScreen struct {
	credits *info.Layer
	controls *info.Layer
	menu menu.Menu
	pendingTransition bool
}

const keyMainMenu menu.Key = menu.FirstKey

func New(ctx *context.Context) (*WinScreen, error) {
	win := &WinScreen{ pendingTransition: true }
	
	win.credits = info.New([]string{
		"CODE, ART AND MUSIC BY TINNE",
		"FOR THE \"EBITENGINE HOLIDAY HACK 2023\"",
		"",
		"ALL CODE IS OPEN SOURCE",
		"ALL ASSETS AVAILABLE UNDER CC LICENSES",
		"GITHUB.COM/TINNE26/LUCKYFEET",
	})
	win.controls = info.New(info.ControlsKB)

	opts := win.menu.NewOptionList(keyMainMenu)
	opts.Add(win.credits.NewOption("CREDITS"))
	opts.Add(&menu.NavOption{ Label: "OPTIONS", To: menu.Options })
	opts.Add(&menu.SceneChangeOption{ Label: "BACK TO TITLE", Change: *scene.Pop() })
	win.menu.NewGameOptionsOptionList(ctx)
	win.menu.JumpTo(keyMainMenu)
	
	return win, nil
}

func (self *WinScreen) Update(ctx *context.Context) (*scene.Change, error) {
	if ctx.Scenes.Current() != self { return nil, nil }
	if self.pendingTransition {
		self.pendingTransition = false
		return scene.PushTo(keys.BriefBlackout), nil
	}

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

func (self *WinScreen) DrawLogical(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	if !foremost { return }
	ctx.State.Editing = false

	ctx.Background.DrawLogical(canvas, ctx)
	
	// draw clear time (white part behind, black part in front)
	white := color.RGBA{244, 244, 244, 244} // slightly translucid
	black := color.RGBA{ 16,  16,  16, 255}
	bounds := canvas.Bounds()
	x := bounds.Dx()/2
	y := bounds.Dy()/4
	strs := []string{ "CLEARED IN " + utils.FmtTicksToTimeStrCents(ctx.State.LastClearTicks) }
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

func (self *WinScreen) DrawHiRes(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	// ...
}
