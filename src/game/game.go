package game

import "fmt"
import "math"
import "io/fs"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/utils"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/game/material/scene/registry"
import "github.com/tinne26/luckyfeet/src/game/material/scene/keys"
import "github.com/tinne26/luckyfeet/src/game/settings"

// asssert interface compliance
var _ ebiten.Game = (*Game)(nil)

// main type implementing ebiten.Game
type Game struct {
	ctx *context.Context
	canvas *ebiten.Image
	prevScreenFitMode settings.ScreenFitMode
	showFPS bool
}

func New(filesys fs.FS) (*Game, error) {
	// create context
	scenes := registry.NewSceneManager()
	ctx, err := context.New(filesys, scenes)
	if err != nil { return nil, err }
	ctx.Scenes.FirstLoad(keys.FirstSceneKey(), ctx)

	return &Game{
		ctx: ctx,
		canvas: ebiten.NewImage(640, 360), // hardcoding yaaaay
	}, nil
}

func (self *Game) Layout(logicWinWidth, logicWinHeight int) (int, int) {
	panic("using ebitengine >=v2.5.0 LayoutF()")
	// scale := ebiten.DeviceScaleFactor()
	// canvasWidth  := int(math.Ceil(float64(logicWinWidth)*scale))
	// canvasHeight := int(math.Ceil(float64(logicWinHeight)*scale))
	// return canvasWidth, canvasHeight
}

func (self *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	canvasWidth  := math.Ceil(logicWinWidth*scale)
	canvasHeight := math.Ceil(logicWinHeight*scale)
	return canvasWidth, canvasHeight
}

func (self *Game) Update() error {
	var err error
	err = self.ctx.UpdateSystems()
	if err != nil { return err }

	err = self.ctx.Scenes.Update(self.ctx)
	if err != nil { return err }

	// audio start (should be somewhere else...)
	if !self.ctx.Audio.IsActive(au.BgmMain) && au.IsContextReady() {
		self.ctx.Audio.FadeIn(au.BgmMain, 0, 0, 0)
	}

	// global actions
	if self.ctx.Input.Trigger(in.ActionFullscreen) {
		if ebiten.IsFullscreen() {
			fmt.Print("[Switching to windowed mode]\n")
			ebiten.SetFullscreen(false)
		} else {
			fmt.Print("[Switching to fullscreen mode]\n")
			ebiten.SetFullscreen(true)
		}
	}
	if self.ctx.Input.Trigger(in.ActionToggleFPS) {
		self.showFPS = !self.showFPS
	}

	// ...
	return nil
}

func (self *Game) Draw(canvas *ebiten.Image) {
	// clear canvas on changes
	if self.prevScreenFitMode != self.ctx.Settings.ScreenFit {
		self.prevScreenFitMode = self.ctx.Settings.ScreenFit
		canvas.Clear()
	}
	
	// draw on logical canvas first
	self.ctx.Scenes.DrawLogical(self.canvas, self.ctx)

	// fps debug
	if self.showFPS {
		fps := fmt.Sprintf("%.02f FPS", ebiten.ActualFPS())
		text.CenterDrawAt(self.canvas, 320, 10, []string{fps}, color.RGBA{0, 0, 0, 255}, 1)
	}
	
	// apply scaling (we are not using hi res for this game)
	switch self.ctx.Settings.ScreenFit {
	case settings.ScreenFitPixelPerfect:
		_ = utils.ProjectPixelPerfect(self.canvas, canvas)
	case settings.ScreenFitProportional:
		_ = utils.ProjectNearest(self.canvas, canvas)
	case settings.ScreenFitStretch:
		utils.ProjectStretched(self.canvas, canvas)
	default:
		panic("broken code")
	}

	// (if we were using hi res too)
	// self.ctx.Scenes.DrawHiRes(hiResCanvas, self.ctx)
}
