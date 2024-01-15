package carrot

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/utils"
import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"

type Variety uint8
const (
	None = iota
	Orange
	Yellow
	Purple
	numVarieties
)

type Carrot struct {
	Variety Variety
	OriginRow uint8
	OriginCol uint8
}

const carrotsCapacity = 3
type Inventory struct {
	Carrots [carrotsCapacity]Carrot
	ActiveIndex uint8 // 0, 1, 2
	FillLevels [carrotsCapacity]float64 // 0 means inactive and can eat and obtain
	SelectorOpacityBlinker utils.Blinker
}

func (self *Inventory) Initialize() {
	self.SelectorOpacityBlinker = *utils.NewBlinker(0.45, 0.6, 0.007)
}

func (self *Inventory) RemoveAll() {
	self.ActiveIndex = 0
	for i := 0; i < carrotsCapacity; i++ {
		self.Carrots[i] = Carrot{}
	}
}

var unfillSpeeds [numVarieties]float64 = [numVarieties]float64{0.0, 0.002, 0.004, 0.007}
func (self *Inventory) Update(ctx *context.Context) {
	self.SelectorOpacityBlinker.Update()
	for i, _ := range self.Carrots {
		variety := self.Carrots[i].Variety
		if variety == None { continue }
		if self.FillLevels[i] < 1.0 {
			self.FillLevels[i] -= unfillSpeeds[variety]
			if self.FillLevels[i] <= 0.0 {
				self.FillLevels[i] = 0.0
				self.Carrots[i] = Carrot{ Variety: None }
			}
		}
	}

	if ctx.Input.Trigger(in.ActionPrevCarrot) {
		if self.ActiveIndex == 0 {
			self.ActiveIndex = carrotsCapacity - 1
		} else {
			self.ActiveIndex -= 1
		}
		ctx.Audio.PlaySFX(au.SfxClick)
	} else if ctx.Input.Trigger(in.ActionNextCarrot) {
		self.ActiveIndex += 1
		if self.ActiveIndex >= carrotsCapacity {
			self.ActiveIndex = 0
		}
		ctx.Audio.PlaySFX(au.SfxClick)
	}

	if ctx.Input.Trigger(in.ActionUseCarrot) {
		_ = self.TryConsume(ctx)
	}
}

// Necessary to draw the carrot platform fills.
func (self *Inventory) GetFillOpacity(carrotVariety Variety) float32 {
	var level float64 = -1.0
	for i, _ := range self.Carrots {
		if self.Carrots[i].Variety == carrotVariety {
			varietyLevel := self.FillLevels[i]
			if varietyLevel < 1.0 {
				level = max(level, varietyLevel)
			}
		}
	}
	if level == -1.0 { return 0 }

	switch {
	case level > 0.66: return 1.0
	case level > 0.33: return 0.8
	default: return 0.65 // level <= 0.33
	}
}

func (self *Inventory) IsMapCarrotOn(col, row uint8) bool {
	for i, _ := range self.Carrots {
		if self.Carrots[i].Variety != None {
			if self.Carrots[i].OriginRow == row && self.Carrots[i].OriginCol == col {
				return false
			}
		}
	}
	return true
}

func (self *Inventory) TryConsume(ctx *context.Context) bool {
	if self.Carrots[self.ActiveIndex].Variety == None || self.FillLevels[self.ActiveIndex] < 1.0 {
		ctx.Audio.PlaySFX(au.SfxScratch)
		return false
	}
	ctx.Audio.PlaySFX(au.SfxCronch)
	self.FillLevels[self.ActiveIndex] = 0.9999
	return true
}

func (self *Inventory) TryAdd(ctx *context.Context, carrot Carrot) bool {
	if carrot.Variety == None { panic("precondition violation") }

	for i, _ := range self.Carrots {
		idx := (int(self.ActiveIndex) + i) % carrotsCapacity
		if self.Carrots[idx].Variety == None && self.IsMapCarrotOn(carrot.OriginCol, carrot.OriginRow) {
			self.Carrots[idx] = carrot
			self.FillLevels[idx] = 1.0
			return true
		}
	}

	return false
}

func (self *Inventory) Draw(canvas *ebiten.Image, ctx *context.Context) {
	bounds := ctx.Gfxcore.CarrotSelector.Bounds()
	const pad = 3
	x := float64(640 - 6 - bounds.Dx()*carrotsCapacity - pad*(carrotsCapacity - 1))
	y := float64(360 - 6 - bounds.Dy())
	// x := float64(640 - bounds.Dx()*carrotsCapacity - pad*(carrotsCapacity - 1))/2.0
	// y := float64(6)
	
	var opts ebiten.DrawImageOptions
	opts.GeoM.Translate(x, y)
	for i, _ := range self.Carrots {
		// draw selector if on active index
		if i == int(self.ActiveIndex) {
			a := float32(self.SelectorOpacityBlinker.Value())
			opts.ColorScale.Scale(a, a, a, a)
			canvas.DrawImage(ctx.Gfxcore.CarrotSelector, &opts)
			opts.ColorScale.Reset()
		}

		// draw carrot mask
		canvas.DrawImage(ctx.Gfxcore.CarrotNone, &opts)
		variety := self.Carrots[i].Variety
		if variety != None {
			a := float32(self.FillLevels[i])
			if a < 1.0 {
				opts.ColorScale.ScaleWithColor(color.RGBA{230, 30, 230, 255})
				canvas.DrawImage(ctx.Gfxcore.CarrotInUseMask, &opts)
				opts.ColorScale.Reset()
			}

			opts.ColorScale.Scale(a, a, a, a)
			switch variety {
			case Orange: canvas.DrawImage(ctx.Gfxcore.CarrotOrange, &opts)
			case Yellow: canvas.DrawImage(ctx.Gfxcore.CarrotYellow, &opts)
			case Purple: canvas.DrawImage(ctx.Gfxcore.CarrotPurple, &opts)
			default: panic("broken code")
			}
			opts.ColorScale.Reset()
		}

		opts.GeoM.Translate(float64(bounds.Dx()) + pad, 0)
	}
}
