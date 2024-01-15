package tile

import "slices"
import "errors"
import "image"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/components/tile/tcsts"
import "github.com/tinne26/luckyfeet/src/game/utils"
import "github.com/tinne26/luckyfeet/src/game/carrot"

// Raw map structure. For actual play, we reorganize data
// a bit, as most content can be predrawn, and collisions
// can be prepared into a few 32x18 arrays.
type Map struct {
	Layers [][]Tile // for indexing, see tcsts.Layer* constants
	
	ID uint8 // can't be zero
	TransferIDs [3]uint8 // 0 means undefined, not allowed as a map ID
	StartRow uint8
	StartCol uint8
}

func NewMap(id uint8) *Map {
	if id == 0 { panic("map ID can't be zero") }
	return &Map{ ID: id, Layers: make([][]Tile, tcsts.LayerCountSentinel) }
}

func (self *Map) SetTile(newTile Tile, layerIndex int) {
	// find insertion position
	layer := self.Layers[layerIndex]
	insertPos, found := slices.BinarySearchFunc(layer, newTile, func(tile, target Tile) int {
		return tile.Cmp(target)
	})
	if found { // replace
		layer[insertPos] = newTile
	} else { // insert
		self.Layers[layerIndex] = slices.Insert(layer, insertPos, newTile)
	}
}

func (self *Map) DeleteTile(row, col int, layerIndex int) {
	delTile := Tile{ Column: uint8(col), Row: uint8(row) }
	layer := self.Layers[layerIndex]
	index, found := slices.BinarySearchFunc(layer, delTile, func(tile, target Tile) int {
		return tile.Cmp(target)
	})
	if !found { return }
	self.Layers[layerIndex] = slices.Delete(layer, index, index + 1)
}

func (self *Map) DrawBackLogical(canvas *ebiten.Image, ctx *context.Context, carrots *carrot.Inventory) {
	for _, layer := range self.Layers[tcsts.LayerBack : tcsts.LayerMain] {
		for i, _ := range layer {
			layer[i].Draw(canvas, ctx, carrots)
		}
	}
}

func (self *Map) DrawMainLogical(canvas *ebiten.Image, ctx *context.Context, carrots *carrot.Inventory) {
	for _, layer := range self.Layers[tcsts.LayerMain : tcsts.LayerFront] {
		for i, _ := range layer {
			layer[i].Draw(canvas, ctx, carrots)
		}
	}
}

func (self *Map) DrawFrontLogical(canvas *ebiten.Image, ctx *context.Context, carrots *carrot.Inventory) {
	for _, layer := range self.Layers[tcsts.LayerFront : ] {
		for i, _ := range layer {
			layer[i].Draw(canvas, ctx, carrots)
		}
	}
	if ctx.State.Editing {
		tile := Tile{ ID: tcsts.StartPoint, Column: self.StartCol, Row: self.StartRow }
		tile.Draw(canvas, ctx, carrots)
	}
}

// Called like tilemap.Collides(ctx, rect, tcsts.LayerBack/LayerMain/LayerFront)
func (self *Map) Collides(ctx *context.Context, carrots *carrot.Inventory, rect image.Rectangle, layer int) bool {
	_, found := self.GetFirstCollision(ctx, carrots, rect, layer)
	return found
}

func (self *Map) GetFirstCollision(ctx *context.Context, carrots *carrot.Inventory, rect image.Rectangle, layer int) (Tile, bool) {
	tiles := self.Layers[layer]
	if len(tiles) == 0 { return Tile{}, false }
	
	var clamp = func(i int) uint8 { return uint8(min(max(i, 0), 255)) }
	minCol, minRow := clamp(rect.Min.X/20), clamp(rect.Min.Y/20)
	maxCol, maxRow := clamp(rect.Max.X/20), clamp(rect.Max.Y/20)
	minTile := Tile{ Row: minRow, Column: minCol }
	maxTile := Tile{ Row: maxRow, Column: maxCol }

	// get min/max indices
	minIndex, _ := slices.BinarySearchFunc(tiles, minTile, func(tile, target Tile) int {
		return tile.Cmp(target)
	})
	if minIndex >= len(tiles) { return Tile{}, false }
	maxIndex, found := slices.BinarySearchFunc(tiles[minIndex : ], maxTile, func(tile, target Tile) int {
		return tile.Cmp(target)
	})
	maxIndex += minIndex
	if found { maxIndex += 1 }

	for i, _ := range tiles[minIndex : maxIndex] {
		col := tiles[minIndex + i].Column
		if col < minCol || col > maxCol { continue }
		if tiles[minIndex + i].Collides(ctx, carrots, rect) {
			return tiles[minIndex + i], true
		}
	}
	return Tile{}, false
}

func (self *Map) HasLandingFor(ctx *context.Context, carrots *carrot.Inventory, ox, fx, y int, layer int) bool {
	tiles := self.Layers[layer]
	if len(tiles) == 0 { return false }
	
	var clamp = func(i int) uint8 { return uint8(min(max(i, 0), 255)) }
	row := clamp(y/20)
	minCol, maxCol := clamp(ox/20), clamp(fx/20)
	minTile := Tile{ Row: row, Column: minCol }
	maxTile := Tile{ Row: row, Column: maxCol }

	// get min/max indices
	minIndex, _ := slices.BinarySearchFunc(tiles, minTile, func(tile, target Tile) int {
		return tile.Cmp(target)
	})
	if minIndex >= len(tiles) { return false }
	maxIndex, found := slices.BinarySearchFunc(tiles[minIndex : ], maxTile, func(tile, target Tile) int {
		return tile.Cmp(target)
	})
	maxIndex += minIndex
	if found { maxIndex += 1 }

	for i, _ := range tiles[minIndex : maxIndex] {
		col := tiles[minIndex + i].Column
		if col < minCol || col > maxCol { continue }
		if tiles[minIndex + i].IsLandingFor(ctx, carrots, ox, fx, y) { return true }
	}
	return false
}

func (self *Map) GetTileIDAt(row, col uint8, layer int) (uint8, bool) {
	target := Tile{ Row: row, Column: col }
	tiles := self.Layers[layer]
	index, found := slices.BinarySearchFunc(tiles, target, func(tile, target Tile) int {
		return tile.Cmp(target)
	})
	if !found { return tcsts.TileTypeMax, false }
	return tiles[index].ID, true
}

func (self *Map) ExportToString() (string, error) {
	data := make([]byte, 0, 1024)
	data = append(data, self.ID)
	data = append(data, self.TransferIDs[0])
	data = append(data, self.TransferIDs[1])
	data = append(data, self.TransferIDs[2])
	data = append(data, self.StartRow)
	data = append(data, self.StartCol)

	for i, _ := range self.Layers {
		numTiles := len(self.Layers[i])
		if numTiles == 0 { continue }
		data = append(data, uint8(i))
		data = append(data, uint8(uint16(numTiles) >> 8))
		data = append(data, uint8(uint16(numTiles) & 0xFF))
		for tileIndex := 0; tileIndex < numTiles; tileIndex++ {
			data = self.Layers[i][tileIndex].EncodeToBytes(data)
		}
	}
	
	return utils.GzipAndEncodeAsCh426(data)
}

func (self *Map) LoadFromString(data string) error {
	bytes, err := utils.DecodeFromCh426AndUngzip(data)
	if err != nil { return err }
	
	if len(bytes) < 7 { return errors.New("not enough data") }
	self.ID = bytes[0]
	self.TransferIDs[0] = bytes[1]
	self.TransferIDs[1] = bytes[2]
	self.TransferIDs[2] = bytes[3]
	self.StartRow = bytes[4]
	self.StartCol = bytes[5]

	var layerID uint8 = 255
	bytes = bytes[6 : ]
	for len(bytes) > 3 {
		newLayerID := bytes[0]
		if int(newLayerID) >= len(self.Layers) {
			return errors.New("too many layers encoded in the data")
		}
		if layerID >= newLayerID && layerID != 255 {
			return errors.New("invalid layer ordering")
		}
		layerID = newLayerID

		// read layer
		layerTiles := (uint16(bytes[1]) << 8) | uint16(bytes[2])
		if layerTiles == 0 {
			return errors.New("layers with zero tiles must not be encoded")
		}
		if layerTiles > 576 {
			return errors.New("layer contains too many tiles")
		}
		if len(bytes) < int(layerTiles*4 + 3) {
			return errors.New("not enough data for the declared tiles")
		}
		self.Layers[layerID] = make([]Tile, layerTiles)
		for i := uint16(0); i < layerTiles; i++ {
			startIndex := 3 + (i << 2)
			self.Layers[layerID][i] = DecodeTileFromBytes(bytes[startIndex : startIndex + 4])
		}
		bytes = bytes[3 + layerTiles*4 : ]
	}

	if len(bytes) != 0 { return errors.New("truncated data end") }


	// verify proper order of tiles
	for _, layer := range self.Layers {
		var prevTile Tile
		for tileIndex, tile := range layer {
			if tileIndex > 0 {
				if tile.Cmp(prevTile) != 1 { panic("broken code") }
			}
			prevTile = tile
		}
	}

	return nil
}

func LoadMapFromString(data string) (*Map, error) {
	tilemap := NewMap(1)
	err := tilemap.LoadFromString(data)
	return tilemap, err
}
