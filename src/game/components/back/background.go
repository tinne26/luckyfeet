package back

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/interfaces"

var _ interfaces.Background[*context.Context] = (*Background)(nil)

type Background struct {
	particles [64]Particle
	palette [5]color.RGBA
	patternOffsets [5]int
	animTicks int
}

func New() *Background {
	return &Background{
		palette: [5]color.RGBA{
			{209, 235, 255, 255},
			{198, 213, 255, 255},
			{186, 191, 255, 255},
			{175, 168, 255, 255},
			{165, 145, 255, 255},
		},
		patternOffsets: [5]int{0, 0, 0, 0, 0},
	}
}

func (self *Background) Update() {
	// update background animation
	self.animTicks += 1
	if self.animTicks % 128 == 0 {
		self.patternOffsets[0] -= 1
		if self.patternOffsets[0] == -8 {
			self.patternOffsets[0] = 0
		}
	}
	if self.animTicks % 64 == 0 {
		self.patternOffsets[1] -= 1
		if self.patternOffsets[1] == -8 {
			self.patternOffsets[1] = 0
		}
	}
	if self.animTicks % 32 == 0 {
		self.patternOffsets[2] -= 1
		if self.patternOffsets[2] == -8 {
			self.patternOffsets[2] = 0
		}
	}
	if self.animTicks % 16 == 0 {
		self.patternOffsets[3] -= 1
		if self.patternOffsets[3] == -8 {
			self.patternOffsets[3] = 0
		}
	}

	// update particles
	for i, _ := range self.particles {
		self.particles[i].Update()
	}
}

func (self *Background) DrawLogical(canvas *ebiten.Image, ctx *context.Context) {
	// draw background pattern
	bounds := canvas.Bounds()
	if bounds.Dx() != 640 { panic("expected canvas to be 640px wide") }
	if bounds.Dy() != 360 { panic("expected canvas to be 360px tall") }

	heights := []int{125, 90, 64, 48, 33}

	y := 0
	for i, height := range heights {
		rect := image.Rect(0, y, 640, y + height)
		canvas.SubImage(rect).(*ebiten.Image).Fill(self.palette[i])
		y += height
	}

	patternBounds := ctx.Gfxcore.BackPatternMaskX32.Bounds()
	patternWidth, patternHeight := patternBounds.Dx(), patternBounds.Dy()
	
	var opts ebiten.DrawImageOptions
	opts.GeoM.Translate(0, -float64(patternHeight))
	for i, height := range heights {
		opts.GeoM.Translate(float64(self.patternOffsets[i]), float64(height))
		if i == len(heights) - 1 { break }
		
		opts.ColorScale.Reset()
		opts.ColorScale.ScaleWithColor(self.palette[i + 1])
		canvas.DrawImage(ctx.Gfxcore.BackPatternMaskX32, &opts)
		opts.GeoM.Translate(float64(patternWidth), 0)
		canvas.DrawImage(ctx.Gfxcore.BackPatternMaskX32, &opts)
		opts.GeoM.Translate(float64(patternWidth), 0)
		canvas.DrawImage(ctx.Gfxcore.BackPatternMaskX32, &opts)
		opts.GeoM.SetElement(0, 2, 0) // reset X to 0
	}

	// draw particles
	opts.GeoM.Reset()
	opts.ColorScale.Reset()
	for i, _ := range self.particles {
		self.particles[i].DrawLogical(canvas, ctx, &opts)
	}
}
