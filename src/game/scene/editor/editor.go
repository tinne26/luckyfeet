package editor

import "strings"
import "image/color"
import "math/rand"
import "strconv"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/game/material/scene/keys"
import "github.com/tinne26/luckyfeet/src/game/components/menu"
import "github.com/tinne26/luckyfeet/src/game/components/info"
import "github.com/tinne26/luckyfeet/src/game/components/menuhint"
import "github.com/tinne26/luckyfeet/src/game/components/tile"
import "github.com/tinne26/luckyfeet/src/game/utils"

var _ scene.Scene[*context.Context] = (*Editor)(nil)

type Editor struct {
	controls *info.Layer
	menuActive bool
	menu menu.Menu
	
	tileX int
	tileY int
	tileBar TileBar
	maps []*tile.Map
	mapIndex int
	blinker *utils.Blinker
	menuOptsToRefreshOnMapChange []func()
	pendingTransition bool
}

const (
	keyMainMenu     menu.Key = menu.FirstKey
	keyEditor       menu.Key = menu.FirstKey + 1
	keyProject      menu.Key = menu.FirstKey + 2
	keyTransfers    menu.Key = menu.FirstKey + 3
	keySetSpawn     menu.Key = menu.FirstKey + 4
	keySetTransfers menu.Key = menu.FirstKey + 5
)

var menuTitles = []string{
	"TODAY ON THE MENU",
	"LET'S TAKE A BREAK",
	"TIRED OF EDITING?",
	"WHAT A NICE PLACE",
	"PAUSE AND RELAX",
	"MADE WITH EBITENGINE",
}

func New(ctx *context.Context) (*Editor, error) {
	var controls info.Layer
	editor := &Editor{
		controls: &controls,
		mapIndex: 0,
		blinker: utils.NewBlinker(0.0, 0.22, 0.02), // min, max, speedOverOne
		menuOptsToRefreshOnMapChange: make([]func(), 0, 5),
		pendingTransition: true,
	}
	editor.tileBar.ResetVariations()

	// load maps
	if ctx.State.LoadMapDataFromClipboard {
		ctx.State.LoadMapDataFromClipboard = false
		mapStrs := strings.Split(strings.TrimSpace(utils.ReadClipboard()), ".")
		editor.maps = make([]*tile.Map, len(mapStrs))
		for i, str := range mapStrs {
			var err error
			editor.maps[i], err = tile.LoadMapFromString(str)
			if err != nil { return editor, err }
		}
	} else {
		editor.maps = make([]*tile.Map, 1)
		editor.maps[0] = tile.NewMap(1)
	}

	var mainMenu menu.Menu
	mainMenu.Title = menuTitles[rand.Intn(len(menuTitles))]
	opts := mainMenu.NewOptionList(keyMainMenu)
	opts.AddBackOption(&menu.BasicOption{
		Label: "CONTINUE",
		Func: func(*context.Context) menu.Key {
			editor.menuActive = false
			return menu.NoChange
		},
	})
	opts.Add(&menu.NavOption{ Label: "TRANSFERS", To: keyTransfers })
	opts.Add(&menu.NavOption{ Label: "OPTIONS", To: menu.Options })
	opts.Add(controls.NewOption("CONTROLS"))
	opts.Add(&menu.NavOption{ Label: "PROJECT", To: keyProject })
	
	opts = mainMenu.NewOptionList(keyTransfers)
	opts.Add(&menu.NavOption{ Label: "SET SPAWN", To: keySetSpawn })
	opts.Add(&menu.NavOption{ Label: "SET TRANSFERS", To: keySetTransfers })
	opts.Add(&JumpToOption{ Editor: editor })
	opts.AddBackOption(&menu.NavOption{ Label: "BACK", To: keyMainMenu })

	opts = mainMenu.NewOptionList(keySetSpawn)
	opt := &TileOption{
		Label: "SPAWN COLUMN",
		MinTile: 0,
		MaxTile: 31,
		NotifyChange: editor.notifyMapStartColChange,
	}
	editor.menuOptsToRefreshOnMapChange = append(editor.menuOptsToRefreshOnMapChange, func() {
		opt.CurrentTile = int(editor.maps[editor.mapIndex].StartCol)
	})
	opts.Add(opt)
	opt = &TileOption{
		Label: "SPAWN ROW",
		MinTile: 0,
		MaxTile: 17,
		NotifyChange: editor.notifyMapStartRowChange,
	}
	editor.menuOptsToRefreshOnMapChange = append(editor.menuOptsToRefreshOnMapChange, func() {
		opt.CurrentTile = int(editor.maps[editor.mapIndex].StartRow)
	})
	opts.Add(opt)
	
	opts.AddBackOption(&menu.NavOption{ Label: "BACK", To: keyTransfers })

	opts = mainMenu.NewOptionList(keySetTransfers)
	trOpt := &TransferOption{ Label: "TRANSFER A", Editor: editor, TransferIndex: 0 }
	editor.menuOptsToRefreshOnMapChange = append(editor.menuOptsToRefreshOnMapChange, trOpt.refreshFunc)
	opts.Add(trOpt)
	trOpt = &TransferOption{ Label: "TRANSFER B", Editor: editor, TransferIndex: 1 }
	editor.menuOptsToRefreshOnMapChange = append(editor.menuOptsToRefreshOnMapChange, trOpt.refreshFunc)
	opts.Add(trOpt)
	trOpt = &TransferOption{ Label: "TRANSFER C", Editor: editor, TransferIndex: 2 }
	editor.menuOptsToRefreshOnMapChange = append(editor.menuOptsToRefreshOnMapChange, trOpt.refreshFunc)
	opts.Add(trOpt)
	opts.AddBackOption(&menu.NavOption{ Label: "BACK", To: keyTransfers })

	opts = mainMenu.NewOptionList(keyProject)
	opts.Add(&menu.SceneChangeEffectOption{
		Label: "PLAYTEST",
		Change: *scene.PushTo(keys.Play),
		OnConfirm: func(fnCtx *context.Context) error {
			data, err := editor.mapsToString()
			if err != nil { return err }
			fnCtx.State.LoadMapDataFromClipboard = false
			fnCtx.State.PlaytestData = data
			fnCtx.State.PlaytestMapID = uint8(editor.mapIndex + 1) // unsafe but let's assume it works
			return nil
		}})
	opts.Add(&menu.EffectOption{ Label: "SAVE TO CLIPBOARD", OnConfirm: func(*context.Context) error {
		data, err := editor.mapsToString()
		if err != nil { return err }
		utils.WriteClipboard(data)
		return nil
	}})
	opts.Add(&menu.SceneChangeOption{ Label: "EXIT WITHOUT SAVING", Change: *scene.Pop() })
	opts.Add(&menu.EffectOption{ Label: "ADD NEW MAP", OnConfirm: func(*context.Context) error {
		// safety check
		for i, tilemap := range editor.maps {
			if int(tilemap.ID) != i + 1 { panic("broken code") }
		}

		// add map and jump to i
		id := len(editor.maps) + 1
		if id > 255 { panic("can't exceed 256 maps") }
		editor.maps = append(editor.maps, tile.NewMap(uint8(id)))
		editor.mapIndex = id - 1
		editor.menuActive = false
		editor.menu.JumpTo(keyMainMenu)
		editor.mapChangeRefresh()
		return nil
	}})
	opts.AddBackOption(&menu.NavOption{ Label: "BACK", To: keyMainMenu})
	mainMenu.NewGameOptionsOptionList(ctx)
	mainMenu.JumpTo(keyMainMenu)

	editor.menu = mainMenu
	editor.tileX = 60
	editor.tileY = 240
	editor.mapChangeRefresh()
	return editor, nil
}

func (self *Editor) mapsToString() (string, error) {
	var strs []string
	for _, tilemap := range self.maps {
		encodedData, err := tilemap.ExportToString()
		if err != nil { return "", err }
		strs = append(strs, encodedData)
	}
	return strings.Join(strs, "."), nil
}

func (self *Editor) mapChangeRefresh() {
	for i, _ := range self.menuOptsToRefreshOnMapChange {
		self.menuOptsToRefreshOnMapChange[i]()
	}
}

func (self *Editor) notifyMapStartColChange(ctx *context.Context, newStartColumn int) {
	self.maps[self.mapIndex].StartCol = uint8(newStartColumn)
}

func (self *Editor) notifyMapStartRowChange(ctx *context.Context, newStartRow int) {
	self.maps[self.mapIndex].StartRow = uint8(newStartRow)
}

func (self *Editor) Update(ctx *context.Context) (*scene.Change, error) {
	ctx.State.Editing = (ctx.Scenes.Current() == self)
	if !ctx.State.Editing { return nil, nil }
	if self.pendingTransition {
		self.pendingTransition = false
		return scene.PushTo(keys.BriefBlackout), nil
	}

	// helper variables
	var change *scene.Change
	var err error
	
	self.blinker.Update()

	// update background animation
	ctx.Background.Update()

	// update info layer / menu
	if self.controls.IsVisible() {
		self.controls.Update(ctx)
	} else {
		// detect menu opening / closing
		if ctx.Input.Trigger(in.ActionMenu) || ctx.Input.Trigger(in.ActionMenuBrowserAlt) {
			ctx.Audio.PlaySFX(au.SfxConfirm)
			self.menu.JumpTo(keyMainMenu)
			self.menuActive = !self.menuActive
			self.menu.Title = menuTitles[rand.Intn(len(menuTitles))]
		}

		// update menu
		if self.menuActive {
			change, err = self.menu.Update(ctx)
		} else {
			err = self.updateEditor(ctx)
		}
	}
	
	return change, err
}

func (self *Editor) updateEditor(ctx *context.Context) error {
	// update tile bar
	self.tileBar.Update(ctx)
	
	// update tile positioning
	tileMoveDir := self.getTileMoveDir(ctx)
	switch tileMoveDir {
	case in.DirUp    : if self.tileY >= 20 { self.tileY -= 20 }
	case in.DirLeft  : if self.tileX >= 20 { self.tileX -= 20 }
	case in.DirDown  : if self.tileY < 340 { self.tileY += 20 }	
	case in.DirRight : if self.tileX < 620 { self.tileX += 20 }
	}

	// remove/set tile
	if ctx.Input.Trigger(in.ActionBack) {
		ctx.Audio.PlaySFX(au.SfxBack)
		col, row := self.tileX/20, self.tileY/20
		self.maps[self.mapIndex].DeleteTile(row, col, self.tileBar.CurrentLayer())
	} else if ctx.Input.Trigger(in.ActionConfirm) {
		ctx.Audio.PlaySFX(au.SfxClick)
		col, row := self.tileX/20, self.tileY/20
		tile := self.tileBar.CurrentTile()
		tile.Column = uint8(col)
		tile.Row = uint8(row)
		self.maps[self.mapIndex].SetTile(tile, self.tileBar.CurrentLayer())
	}
	
	return nil
}

func (self *Editor) getTileMoveDir(ctx *context.Context) in.Direction {
	const RF, RN = in.RFFast, in.RNFast

	if ctx.Input.Pressed(in.ActionModKey) { return in.DirNone }

	kb := ctx.Input.Keyboard()
	if !kb.Pressed(in.ActionModKey) {
		dir := kb.RepeatDirAs(RF, RN)
		if dir != in.DirNone { return dir }
	}
	return ctx.Input.Gamepad().RepeatDirAs(RF, RN)
}

func (self *Editor) DrawLogical(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	if !foremost { return }

	ctx.Background.DrawLogical(canvas, ctx)
	
	// draw editor
	self.drawEditor(canvas, ctx)

	// draw menu hint
	menuhint.Draw(canvas, ctx)

	// draw controls or menu
	if self.controls.IsVisible() {
		if ctx.Input.UsedGamepadMoreRecentlyThanKeyboard() {
			self.controls.SetContent(info.EditorControlsGP)
		} else {
			self.controls.SetContent(info.EditorControlsKB)
		}
		self.controls.Draw(canvas)
	} else if self.menuActive {
		utils.FillOver(canvas, color.RGBA{0, 0, 0, 48})
		self.menu.DrawLogical(canvas, ctx)
	}
}

func (self *Editor) drawEditor(canvas *ebiten.Image, ctx *context.Context) {
	// draw tiles
	self.maps[self.mapIndex].DrawBackLogical(canvas, ctx, nil)
	self.maps[self.mapIndex].DrawMainLogical(canvas, ctx, nil)
	self.maps[self.mapIndex].DrawFrontLogical(canvas, ctx, nil)
	
	// draw tile bar
	self.tileBar.DrawLogical(canvas, ctx)
	tile := self.tileBar.CurrentTile()
	tile.Column = uint8(self.tileX/20)
	tile.Row = uint8(self.tileY/20)
	tile.Draw(canvas, ctx, nil)

	// draw map ID
	info := "MAP ID #" + strconv.Itoa(int(self.maps[self.mapIndex].ID))
	text.DrawAt(canvas, 4, 360 - 4 - text.LineHeight, []string{info}, color.RGBA{0, 0, 0, 128}, 1)

	// draw tile rect
	a := self.blinker.Value()
	rect := utils.Rect(self.tileX, self.tileY, self.tileX + 20, self.tileY + 20)
	utils.FillOverRectLighter(canvas, rect, utils.RGBAf64(a, a, a, a))
}

func (self *Editor) DrawHiRes(canvas *ebiten.Image, foremost bool, ctx *context.Context) {
	// ...
}
