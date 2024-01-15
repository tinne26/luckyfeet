package editor

import "strconv"

import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/game/components/menu"

// extra option types for unique menus

type TileOption struct {
	Label string
	MinTile int
	MaxTile int
	CurrentTile int
	NotifyChange func(*context.Context, int)
}
func (self *TileOption) Name() string {
	tile := strconv.Itoa(self.CurrentTile)
	return self.Label + " " + string(text.TriangleLeftWithPad) + tile + string(text.TriangleRightWithPad)
}
func (self *TileOption) MaxName() string {
	tile := strconv.Itoa(self.MaxTile*10)
	return self.Label + " " + string(text.TriangleLeftWithPad) + tile + string(text.TriangleRightWithPad)
}
func (self *TileOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *TileOption) HoverUpdate(ctx *context.Context) {
	dir := ctx.Input.RepeatDirAs(in.RFDefault, in.RNDefault)
	if dir == in.DirNone { return }

	prevTile := self.CurrentTile
	if dir == in.DirRight {
		self.CurrentTile = min(self.MaxTile, self.CurrentTile + 1)
	} else if dir == in.DirLeft {
		self.CurrentTile = max(self.MinTile, self.CurrentTile - 1)
	}

	if self.CurrentTile != prevTile {
		self.NotifyChange(ctx, self.CurrentTile)
		ctx.Audio.PlaySFX(au.SfxClick)
	} else {
		ctx.Audio.PlaySFX(au.SfxBack)
	}
}
func (self *TileOption) Confirm(ctx *context.Context) (menu.Key, *scene.Change, error) {
	return menu.NoChange, nil, nil
}

// --- "jump to" option (for changing maps) ---

type JumpToOption struct {
	Editor *Editor
	JumpMapID int
}
func (self *JumpToOption) Name() string {
	if self.JumpMapID == 0 || self.JumpMapID == self.Editor.mapIndex + 1 {
		if len(self.Editor.maps) <= 1 { return "(NO MAPS TO JUMP TO)" }
		self.JumpMapID = int(self.Editor.mapIndex + 2)
		if self.JumpMapID >= len(self.Editor.maps) + 1 {
			self.JumpMapID = 1
		}
	}
	
	id := strconv.Itoa(int(self.JumpMapID))
	return "JUMP TO " + string(text.TriangleLeftWithPad) + id + string(text.TriangleRightWithPad)
}
func (self *JumpToOption) MaxName() string {
	return "(NO MAPS TO JUMP TO)"
}
func (self *JumpToOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *JumpToOption) HoverUpdate(ctx *context.Context) {
	dir := ctx.Input.RepeatDirAs(in.RFDefault, in.RNDefault)
	if dir == in.DirNone || len(self.Editor.maps) <= 1 { return }
	
	ctx.Audio.PlaySFX(au.SfxClick)
	for {
		if dir == in.DirRight {
			self.JumpMapID += 1
			if self.JumpMapID > len(self.Editor.maps) { self.JumpMapID = 1 }
		} else if dir == in.DirLeft {
			self.JumpMapID -= 1
			if self.JumpMapID <= 0 { self.JumpMapID = len(self.Editor.maps) }
		}
		if int(self.JumpMapID) != self.Editor.mapIndex + 1 { break }
	}
}
func (self *JumpToOption) Confirm(ctx *context.Context) (menu.Key, *scene.Change, error) {
	if self.JumpMapID == 0 { panic("broken code") }
	self.Editor.mapIndex = self.JumpMapID - 1
	self.Editor.menu.JumpTo(keyMainMenu)
	self.Editor.menuActive = false
	self.Editor.mapChangeRefresh()
	return menu.NoChange, nil, nil
}

// --- "set transfer" option ---

type TransferOption struct {
	Label string
	Editor *Editor
	AssignedID uint8
	TransferIndex uint8
}
func (self *TransferOption) refreshFunc() {
	tilemap := self.Editor.maps[self.Editor.mapIndex]
	self.AssignedID = tilemap.TransferIDs[self.TransferIndex]
}
func (self *TransferOption) Name() string {
	if self.AssignedID == 0 {
		return self.Label + " " + string(text.TriangleLeftWithPad) + "(NONE)" + string(text.TriangleRightWithPad)
	}
	
	id := strconv.Itoa(int(self.AssignedID))
	return self.Label + " " + string(text.TriangleLeftWithPad) + id + string(text.TriangleRightWithPad)
}
func (self *TransferOption) MaxName() string {
	return self.Label + " " + string(text.TriangleLeftWithPad) + "255" + string(text.TriangleRightWithPad)
}
func (self *TransferOption) SoftHighlight(ctx *context.Context) bool { return false }
func (self *TransferOption) HoverUpdate(ctx *context.Context) {
	dir := ctx.Input.RepeatDirAs(in.RFDefault, in.RNDefault)
	if dir == in.DirNone { return }

	ctx.Audio.PlaySFX(au.SfxClick)
	for {
		if dir == in.DirRight {
			if self.AssignedID == 255 {
				self.AssignedID = 0
			} else {
				self.AssignedID += 1
			}
		} else if dir == in.DirLeft {
			if self.AssignedID == 0 {
				self.AssignedID = 255
			} else {
				self.AssignedID -= 1
			}
		}
		if self.AssignedID != self.Editor.maps[self.Editor.mapIndex].ID { break }
	}
	self.Editor.maps[self.Editor.mapIndex].TransferIDs[self.TransferIndex] = self.AssignedID
}
func (self *TransferOption) Confirm(ctx *context.Context) (menu.Key, *scene.Change, error) {	
	return menu.NoConfirm, nil, nil
}
