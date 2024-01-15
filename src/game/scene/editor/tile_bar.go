package editor

import "fmt"
import "math/rand"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/text"

import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/components/tile"
import "github.com/tinne26/luckyfeet/src/game/components/tile/tcsts"
import "github.com/tinne26/luckyfeet/src/game/utils"

// The structure of the tile bar groups is hardcoded as f*ck
var tileGroups = [][]uint8{
	{tcsts.BackGround, tcsts.BackGroundSide, tcsts.BackGroundCorner}, // back layer ground
	{tcsts.BackGroundMark, tcsts.BackGroundMarkCorner}, // back layer ground marks
	{tcsts.MainGround, tcsts.MainGroundRaiser, tcsts.MainGroundSide, tcsts.MainGroundCorner, tcsts.MainSinglePlatform}, // main layer ground
	{tcsts.MainGroundMark, tcsts.MainGroundMarkCorner}, // main layer ground marks
	{tcsts.MainGrassSide, tcsts.MainGrassSideFull, tcsts.MainGrassCorner, tcsts.MainGrassCornerFull}, // main layer grass
	{ // main layer plats
		tcsts.MainOrangePlatSingle, tcsts.MainOrangePlatLeft, tcsts.MainOrangePlatRight,
		tcsts.MainYellowPlatSingle, tcsts.MainYellowPlatLeft, tcsts.MainYellowPlatRight,
		tcsts.MainPurplePlatSingle, tcsts.MainPurplePlatLeft, tcsts.MainPurplePlatRight,
	}, 
	{tcsts.FrontGround, tcsts.FrontGroundRaiser, tcsts.FrontGroundSide, tcsts.FrontGroundCorner, tcsts.FrontSinglePlatform}, // front layer ground
	{tcsts.FrontGroundMark, tcsts.FrontGroundMarkCorner}, // front layer ground marks
	{tcsts.FrontGrassSide, tcsts.FrontGrassSideFull, tcsts.FrontGrassCorner, tcsts.FrontGrassCornerFull}, // front layer grass
	{tcsts.CarrotOrange, tcsts.CarrotYellow, tcsts.CarrotPurple, tcsts.RaceGoal }, // carrots
	{ // transfers
		tcsts.TransferRightA, tcsts.TransferRightB, tcsts.TransferRightC,
		tcsts.TransferLeftA, tcsts.TransferLeftB, tcsts.TransferLeftC,
		tcsts.TransferUpA, tcsts.TransferUpB, tcsts.TransferUpC, 
		tcsts.TransferDownA, tcsts.TransferDownB, tcsts.TransferDownC, 
	},
}
var tileGroupLayers = []int{
	tcsts.LayerBack,
	tcsts.LayerBackDecor,
	tcsts.LayerMain,
	tcsts.LayerMainDecor,
	tcsts.LayerMain,
	tcsts.LayerMain,
	tcsts.LayerFront,
	tcsts.LayerFrontDecor,
	tcsts.LayerFront,
	tcsts.LayerSpecial,
	tcsts.LayerSpecial,
}

type TileBar struct {
	GroupIndex int
	CursorIndex int
	Orientation tile.Orientation
	Variations []uint8
}

func (self *TileBar) ResetVariations() {
	desiredLen := len(tileGroups[self.GroupIndex])
	if cap(self.Variations) >= desiredLen {
		self.Variations = self.Variations[0 : desiredLen]
		for i := 0; i < desiredLen; i++ {
			self.Variations[i] = 0
		}
	} else {
		self.Variations = make([]uint8, desiredLen)
	}
	self.CursorIndex = 0
}

func (self *TileBar) Update(ctx *context.Context) {
	const RF, RN = in.RFDefault, in.RNDefault

	// TODO: sfxs?
	if ctx.Input.RepeatAs(in.ActionNextTile, RF, RN) {
		if self.CursorIndex < len(tileGroups[self.GroupIndex]) - 1 {
			self.CursorIndex += 1
		} else {
			self.CursorIndex = 0
		}
	} else if ctx.Input.RepeatAs(in.ActionPrevTile, RF, RN) {
		if self.CursorIndex > 0 {
			self.CursorIndex -= 1
		} else {
			self.CursorIndex = len(tileGroups[self.GroupIndex]) - 1
		}
	} else if ctx.Input.Trigger(in.ActionPrevTileGroup) {
		if self.GroupIndex > 0 {
			self.GroupIndex -= 1
			self.ResetVariations()
		}
	} else if ctx.Input.Trigger(in.ActionNextTileGroup) {
		if self.GroupIndex < len(tileGroups) - 1 {
			self.GroupIndex += 1
			self.ResetVariations()
		}
	}

	if ctx.Input.Trigger(in.ActionTileVariation) {
		numVariations := uint8(len(ctx.Gfxcore.Tiles[self.CurrentTileID()]))
		if numVariations > 1 {
			variation := uint8(rand.Intn(int(numVariations)))
			if variation != self.Variations[self.CursorIndex] {
				self.Variations[self.CursorIndex] = variation
			} else {
				self.Variations[self.CursorIndex] = (self.Variations[self.CursorIndex] + 1) % numVariations
			}
		}
	}

	if ctx.Input.Pressed(in.ActionModKey) {
		if ctx.Input.Trigger(in.ActionRotateTileLeft) || ctx.Input.Trigger(in.ActionRotateTileRight) {
			self.Orientation = self.Orientation.Mirrored()
		}
	} else if ctx.Input.Trigger(in.ActionRotateTileRight) {
		self.Orientation = self.Orientation.RotatedRight()
	} else if ctx.Input.Trigger(in.ActionRotateTileLeft) {
		self.Orientation = self.Orientation.RotatedLeft()
	}
}

func (self *TileBar) DrawLogical(canvas *ebiten.Image, ctx *context.Context) {
	numTiles := len(tileGroups[self.GroupIndex])
	
	utils.FillOverRect(canvas, utils.Rect(8, 8, 8 + numTiles*20 + numTiles + 1, 8 + 22), color.RGBA{0, 0, 0, 64})
	for i, tileID := range tileGroups[self.GroupIndex] {
		if i == self.CursorIndex {
			ox := 8 + i*20 + i
			utils.FillOverRect(canvas, utils.Rect(ox, 8, ox + 22, 8 + 22), color.RGBA{32, 0, 32, 64})
		}
		tile.DrawAt(canvas, ctx, nil, 9 + i*20 + i, 9, tileID, self.Variations[i], self.Orientation)
	}

	numVariations := uint8(len(ctx.Gfxcore.Tiles[self.CurrentTileID()]))
	info := fmt.Sprintf("VARIATION %d/%d", self.Variations[self.CursorIndex] + 1, numVariations)
	if self.Orientation.IsMirrored() { info += " [MIRRORED]" }
	text.DrawAt(canvas, 9, 8 + 22, []string{info}, color.RGBA{0, 0, 0, 128}, 1)
}

// Notice: Row and Column fields must be manually set.
func (self *TileBar) CurrentTile() tile.Tile {
	return tile.Tile{
		ID: self.CurrentTileID(),
		Variation: self.Variations[self.CursorIndex],
		Orientation: self.Orientation,
	}
}

func (self *TileBar) CurrentTileID() uint8 {
	return tileGroups[self.GroupIndex][self.CursorIndex]
}

func (self *TileBar) CurrentLayer() int {
	return tileGroupLayers[self.GroupIndex]
}
