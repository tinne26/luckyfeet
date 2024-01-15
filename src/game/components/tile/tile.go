package tile

import "fmt"
import "math"
import "image"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/components/tile/tcsts"
import "github.com/tinne26/luckyfeet/src/game/carrot"

var tileDrawOpts ebiten.DrawImageOptions
var matrices []ebiten.GeoM // orientations follow Orientation values, which are consecutive from 0 - 7
func init() {
	// helper functions
	var rotate = func(matrix *ebiten.GeoM, angle int) {
		matrix.Translate(-10, -10)
		matrix.Rotate(float64(angle)*math.Pi/180)
		matrix.Translate(10, 10)
	}
	var mirror = func(matrix *ebiten.GeoM) {
		matrix.Scale(-1, 1)
		matrix.Translate(20, 0)
	}

	// matrices global var init
	matrices = make([]ebiten.GeoM, 8)

	// rotations
	rotate(&matrices[1],  90)
	rotate(&matrices[2], 180)
	rotate(&matrices[3], 270)
	
	// mirrors
	matrices[5] = matrices[1]
	matrices[6] = matrices[2]
	matrices[7] = matrices[3]
	mirror(&matrices[4])
	mirror(&matrices[5])
	mirror(&matrices[6])
	mirror(&matrices[7])
}

type Tile struct {
	ID uint8
	Variation uint8
	Orientation Orientation // only uses lowest 3 bits
	Row uint8
	Column uint8
}

func (self *Tile) String() string {
	mirrored := ""
	if self.Orientation.IsMirrored() { mirrored = " mirrored" }
	return fmt.Sprintf("(T%d/V%d, %s%s, [%s])", self.ID, self.Variation, self.Orientation.RotationStr(), mirrored, self.RawRect())
}

func (self *Tile) Cmp(other Tile) int {
	if self.Row < other.Row { return -1 }
	if self.Row > other.Row { return  1 }
	if self.Column < other.Column { return -1 }
	if self.Column > other.Column { return  1 }
	return 0
}

func (self *Tile) EncodeToBytes(buffer []byte) []byte {
	return append(buffer, self.ID, (self.Variation << 3) | uint8(self.Orientation), self.Row, self.Column)
}

func DecodeTileFromBytes(bytes []byte) Tile {
	if len(bytes) != 4 { panic("expected 4 bytes") }
	return Tile{
		ID: bytes[0],
		Variation: bytes[1] >> 3,
		Orientation: Orientation(bytes[1] & 0x07),
		Row: bytes[2],
		Column: bytes[3],
	}
}

func (self *Tile) Draw(canvas *ebiten.Image, ctx *context.Context, carrots *carrot.Inventory) {
	self.DrawAt(canvas, ctx, carrots, int(self.Column)*20, int(self.Row)*20)
}

func (self *Tile) DrawAt(canvas *ebiten.Image, ctx *context.Context, carrots *carrot.Inventory, x, y int) {
	DrawAt(canvas, ctx, carrots, x, y, self.ID, self.Variation, self.Orientation)
}

func DrawAt(canvas *ebiten.Image, ctx *context.Context, carrots *carrot.Inventory, x, y int, id uint8, variation uint8, orientation Orientation) {	
	tileDrawOpts.GeoM = matrices[orientation]
	tileDrawOpts.GeoM.Translate(float64(x), float64(y))
	if id <= tcsts.RaceGoal {	
		canvas.DrawImage(ctx.Gfxcore.Tiles[id][variation], &tileDrawOpts)
	} else {
		drawSpecialAt(canvas, ctx, carrots, x, y, id, variation, orientation)
	}
}

// Notice: GeoM translation is already applied.
func drawSpecialAt(canvas *ebiten.Image, ctx *context.Context, carrots *carrot.Inventory, x, y int, id uint8, variation uint8, orientation Orientation) {
	if id < tcsts.MainOrangePlatSingle {
		// carrot
		col, row := x/20, y/20
		if ctx.State.Editing || (carrots != nil && carrots.IsMapCarrotOn(uint8(col), uint8(row))) {
			canvas.DrawImage(ctx.Gfxcore.Tiles[id][variation], &tileDrawOpts)
		} else {
			canvas.DrawImage(ctx.Gfxcore.Tiles[tcsts.CarrotMissing][0], &tileDrawOpts)
		}
	} else if id < tcsts.TransferUp {
		// carrot platform
		var variety carrot.Variety
		switch {
		case id < tcsts.MainYellowPlatSingle: variety = carrot.Orange
		case id < tcsts.MainPurplePlatSingle: variety = carrot.Yellow
		default: variety = carrot.Purple
		}

		var fillOpacity float32
		if carrots != nil {
			fillOpacity = carrots.GetFillOpacity(variety)
		}
		tileDrawOpts.ColorScale.Scale(fillOpacity, fillOpacity, fillOpacity, fillOpacity)
		canvas.DrawImage(ctx.Gfxcore.Tiles[id + 1][0], &tileDrawOpts) // draw filler
		tileDrawOpts.ColorScale.Reset()
		canvas.DrawImage(ctx.Gfxcore.Tiles[id][variation], &tileDrawOpts)
	} else {
		// transfer
		if !ctx.State.Editing {
			id -= (id - tcsts.TransferUp) & 0b011
		}
		canvas.DrawImage(ctx.Gfxcore.Tiles[id][0], &tileDrawOpts)
	}
}

func (self *Tile) RawRect() image.Rectangle {
	tox, toy := int(self.Column)*20, int(self.Row)*20
	return image.Rect(tox, toy, tox + 20, toy + 20)
}

func (self *Tile) Collides(ctx *context.Context, carrots *carrot.Inventory, rect image.Rectangle) bool {
	if CollisionFuncs[tcsts.GeometryTable[self.ID]](ctx, self.Orientation, self.RawRect(), rect) == false {
		return false
	}
	if self.ID >= tcsts.MainOrangePlatSingle && self.ID <= tcsts.MainPurplePlatRightFill {
		return carrots.GetFillOpacity(carrotVarietyTable[self.ID - tcsts.MainOrangePlatSingle]) != 0
	}
	return true
}

func (self *Tile) IsLandingFor(ctx *context.Context, carrots *carrot.Inventory, ox, fx, y int) bool {
	if LandingFuncs[tcsts.GeometryTable[self.ID]](ctx, self.Orientation, self.RawRect(), ox, fx, y) == false {
		return false
	}
	if self.ID >= tcsts.MainOrangePlatSingle && self.ID <= tcsts.MainPurplePlatRightFill {
		variety := carrotVarietyTable[self.ID - tcsts.MainOrangePlatSingle]
		opacity := carrots.GetFillOpacity(variety)
		return opacity != 0
	}
	return true
}

// crazy hardcoding hacks
var carrotVarietyTable [1 + tcsts.MainPurplePlatRightFill - tcsts.MainOrangePlatSingle]carrot.Variety
func init() {
	i := 0
	for j := i + 6; i < j; i++ {
		carrotVarietyTable[i] = carrot.Orange
	}
	for j := i + 6; i < j; i++ {
		carrotVarietyTable[i] = carrot.Yellow
	}
	for j := i + 6; i < j; i++ {
		carrotVarietyTable[i] = carrot.Purple
	}
}

