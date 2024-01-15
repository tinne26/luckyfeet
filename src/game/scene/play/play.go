package play

import "strings"
import "math/rand"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/scene"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/game/material/scene/keys"
import "github.com/tinne26/luckyfeet/src/game/components/menu"
import "github.com/tinne26/luckyfeet/src/game/components/info"
import "github.com/tinne26/luckyfeet/src/game/components/menuhint"
import "github.com/tinne26/luckyfeet/src/game/components/racetimer"
import "github.com/tinne26/luckyfeet/src/game/components/tile"
import "github.com/tinne26/luckyfeet/src/game/components/tile/tcsts"
import "github.com/tinne26/luckyfeet/src/game/player"
import "github.com/tinne26/luckyfeet/src/game/utils"
import "github.com/tinne26/luckyfeet/src/game/carrot"

var _ scene.Scene[*context.Context] = (*Play)(nil)

type Play struct {
	player *player.Player
	maps []*tile.Map // map 0 must be nil
	mapIndex int

	controls *info.Layer
	menuActive bool
	menu menu.Menu
	carrots carrot.Inventory
	
	smallLightBlinker *utils.Blinker
	bigLightBlinker *utils.Blinker
	lightScaleBlinker *utils.Blinker
	ticksStopwatch int
	pendingTransition bool
}

const (
	keyMainMenu menu.Key = menu.FirstKey
	keyStopIt   menu.Key = menu.FirstKey + 1
)

var menuTitles = []string{
	"TODAY ON THE MENU",
	"MADE WITH EBITENGINE",
	"ALL THE VITAMIN A YOU NEED",
	"DON'T SPEEDRUN DINNER",
	"EAT YOUR VEGGIES",
	"JUMP ME TO THE MOON",
	"NO COST TOO GREAT",
	"CARROT AND STICK",
}

func New(ctx *context.Context) (*Play, error) {
	var controls info.Layer
	play := &Play{ controls: &controls, mapIndex: 0, pendingTransition: true }
	play.carrots.Initialize()
	play.player = player.New(ctx)

	// load maps
	var mapsData string
	if ctx.State.PlaytestData != "" {
		mapsData = ctx.State.PlaytestData
		play.mapIndex = int(ctx.State.PlaytestMapID) - 1 // unsafe but let's assume it works
	} else if ctx.State.LoadMapDataFromClipboard {
		ctx.State.LoadMapDataFromClipboard = false
		mapsData = utils.ReadClipboard()
	} else {
		mapsData = defaultMapsData
	}

	var mapDatas []string = strings.Split(strings.TrimSpace(mapsData), ".")
	play.maps = make([]*tile.Map, len(mapDatas))
	for i, str := range mapDatas {
		var err error
		play.maps[i], err = tile.LoadMapFromString(str)
		if err != nil { return play, err }
	}

	// create menu
	var mainMenu menu.Menu
	mainMenu.Title = menuTitles[rand.Intn(len(menuTitles))]
	opts := mainMenu.NewOptionList(keyMainMenu)
	opts.AddBackOption(&menu.BasicOption{
		Label: "CONTINUE",
		Func: func(*context.Context) menu.Key {
			play.menuActive = false
			return menu.NoChange
		},
	})
	opts.Add(&menu.NavOption{ Label: "OPTIONS", To: menu.Options })
	opts.Add(controls.NewOption("CONTROLS"))
	opts.Add(&menu.NavOption{ Label: "STOP IT", To: keyStopIt })
	
	mainMenu.NewGameOptionsOptionList(ctx)

	opts = mainMenu.NewOptionList(keyStopIt)
	opts.Add(&menu.EffectOption{
		Label: "RESPAWN",
		OnConfirm: func(fnCtx *context.Context) error {
			play.carrots.RemoveAll()
			play.respawnPlayer(fnCtx)
			play.menu.JumpTo(keyMainMenu)
			play.menuActive = false
			return nil
		},
	})
	opts.Add(&menu.SceneChangeEffectOption{
		Label: "EXIT RACE",
		Change: *scene.Pop(),
		OnConfirm: func(fnCtx *context.Context) error {
			fnCtx.State.PlaytestData = ""
			return nil
		},
	})
	opts.AddBackOption(&menu.NavOption{ Label: "BACK", To: keyMainMenu })

	mainMenu.JumpTo(keyMainMenu)

	play.menu = mainMenu
	play.smallLightBlinker = utils.NewBlinker(0.8, 1.0, 0.0055)
	play.bigLightBlinker = utils.NewBlinker(0.8, 0.9, 0.00035)
	play.lightScaleBlinker = utils.NewBlinker(1.1, 1.3, 0.002)
	play.respawnPlayer(ctx)

	return play, nil
}

func (self *Play) respawnPlayer(ctx *context.Context) {
	tilemap := self.maps[self.mapIndex]
	self.player.Respawn(ctx, tilemap)
}

func (self *Play) Update(ctx *context.Context) (*scene.Change, error) {
	if ctx.Scenes.Current() != self { return nil, nil }
	if self.pendingTransition {
		self.pendingTransition = false
		return scene.PushTo(keys.BriefBlackout), nil
	}

	// helper variables
	var change *scene.Change
	var err error
	
	// update background animation
	ctx.Background.Update()

	// update info layer / menu
	if self.controls.IsVisible() {
		self.controls.Update(ctx)
	} else {
		// detect menu opening / closing
		if ctx.Input.Trigger(in.ActionMenu) {
			ctx.Audio.PlaySFX(au.SfxConfirm)
			self.menu.JumpTo(keyMainMenu)
			self.menuActive = !self.menuActive
			self.menu.Title = menuTitles[rand.Intn(len(menuTitles))]
		}

		// update menu
		if self.menuActive {
			change, err = self.menu.Update(ctx)
		} else {
			change, err = self.mainUpdate(ctx)
			if change != nil || err != nil {
				return change, err
			}
		}
	}
	
	return change, err
}

func (self *Play) mainUpdate(ctx *context.Context) (*scene.Change, error) {
	var err error

	self.ticksStopwatch += 1

	err = self.player.Update(ctx, &self.carrots, self.maps[self.mapIndex])
	if err != nil { return nil, err }

	if self.player.HasFallen() {
		ctx.Audio.PlaySFX(au.SfxBack)
		self.carrots.RemoveAll()
		self.respawnPlayer(ctx)
	} else {
		change, err := self.playerSpecialUpdate(ctx)
		if change != nil || err != nil { return change, err }
	}

	self.carrots.Update(ctx)
	self.smallLightBlinker.Update()
	self.bigLightBlinker.Update()
	self.lightScaleBlinker.Update()

	return nil, nil
}

func (self *Play) playerSpecialUpdate(ctx *context.Context) (*scene.Change, error) {
	rect := self.player.GetSpecialRect()
	tilemap := self.maps[self.mapIndex]
	tile, found := tilemap.GetFirstCollision(ctx, &self.carrots, rect, tcsts.LayerSpecial)
	if !found { return nil, nil }
	
	switch tile.ID {
	case tcsts.RaceGoal:
		ctx.State.LastClearTicks = self.ticksStopwatch
		ctx.Audio.PlaySFX(au.SfxClick)
		return scene.ReplaceTo(keys.WinScreen), nil
	case tcsts.CarrotOrange:
		carr := carrot.Carrot{ Variety: carrot.Orange, OriginCol: tile.Column, OriginRow: tile.Row }
		if self.carrots.TryAdd(ctx, carr) { ctx.Audio.PlaySFX(au.SfxClick) }
	case tcsts.CarrotYellow:
		carr := carrot.Carrot{ Variety: carrot.Yellow, OriginCol: tile.Column, OriginRow: tile.Row }
		if self.carrots.TryAdd(ctx, carr) { ctx.Audio.PlaySFX(au.SfxClick) }
	case tcsts.CarrotPurple:
		carr := carrot.Carrot{ Variety: carrot.Purple, OriginCol: tile.Column, OriginRow: tile.Row }
		if self.carrots.TryAdd(ctx, carr) { ctx.Audio.PlaySFX(au.SfxClick) }
	default:
		if tile.ID >= tcsts.TransferUp {
			var targetMapID uint8
			switch tile.ID & 0b11 { // could be shortened with [(tile.ID & 0b11) - 1]
			case 1: targetMapID = tilemap.TransferIDs[0] // A
			case 2: targetMapID = tilemap.TransferIDs[1] // B
			case 3: targetMapID = tilemap.TransferIDs[2] // C
			default:
				panic("broken code")
			}
			if targetMapID == 0 {
				panic("map transfer doesn't have a target defined")
			}

			self.mapIndex = int(targetMapID - 1) // this is not safe, but I have bigger problems in my life
			self.respawnPlayer(ctx)
			ctx.Audio.PlaySFX(au.SfxClick)
			return scene.PushTo(keys.BriefBlackout), nil
		}
	}

	return nil, nil
}

func (self *Play) DrawLogical(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	if !foremost { return }

	ctx.Background.DrawLogical(canvas, ctx)
	
	// draw main content
	self.mainDraw(canvas, ctx)

	// draw controls or menu
	if self.controls.IsVisible() {
		if ctx.Input.UsedGamepadMoreRecentlyThanKeyboard() {
			self.controls.SetContent(info.ControlsGP)
		} else {
			self.controls.SetContent(info.ControlsKB)
		}
		self.controls.Draw(canvas)
	} else if self.menuActive {
		utils.FillOver(canvas, color.RGBA{0, 0, 0, 48})
		self.menu.DrawLogical(canvas, ctx)
	} else {
		// draw menu hint
		menuhint.Draw(canvas, ctx)
	}
}

func (self *Play) mainDraw(canvas *ebiten.Image, ctx *context.Context) {
	// draw back lighting to slightly improve contrast
	x, y := self.player.GetLightCenterPoint()
	var opts ebiten.DrawImageOptions
	opts.Filter = ebiten.FilterLinear
	lightScale := self.lightScaleBlinker.Value()

	imgBounds := ctx.Gfxcore.BackLightingBig.Bounds()
	opts.GeoM.Translate(float64(-imgBounds.Dx()/2), float64(-imgBounds.Dy()/2))
	opts.GeoM.Scale(lightScale, lightScale)
	opts.GeoM.Translate(float64(x), float64(y))
	a := float32(self.bigLightBlinker.Value())
	opts.ColorScale.Scale(a, a, a, a)
	canvas.DrawImage(ctx.Gfxcore.BackLightingBig, &opts)

	opts.GeoM.Reset()
	opts.ColorScale.Reset()

	imgBounds = ctx.Gfxcore.BackLightingSmall.Bounds()
	opts.GeoM.Translate(float64(-imgBounds.Dx()/2), float64(-imgBounds.Dy()/2))
	opts.GeoM.Scale(lightScale, lightScale)
	opts.GeoM.Translate(float64(x), float64(y))
	a = float32(self.bigLightBlinker.Value())
	opts.ColorScale.Scale(a, a, a, a)
	canvas.DrawImage(ctx.Gfxcore.BackLightingSmall, &opts)
	
	// draw map and player
	self.maps[self.mapIndex].DrawBackLogical(canvas, ctx, &self.carrots)
	if self.player.BehindMain() { self.player.Draw(canvas, ctx) }
	self.maps[self.mapIndex].DrawMainLogical(canvas, ctx, &self.carrots)
	if self.player.InFrontMain() { self.player.Draw(canvas, ctx) }
	self.maps[self.mapIndex].DrawFrontLogical(canvas, ctx, &self.carrots)

	// draw timer
	racetimer.Draw(canvas, ctx, self.ticksStopwatch)

	// draw carrots inventory
	self.carrots.Draw(canvas, ctx)
}

func (self *Play) DrawHiRes(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	// ...
}
