package back

import "math/rand"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/utils"

type Particle struct {
	X float64
	Y float64
	Speed float64
	Dir uint8 // 0 == NE, 1 == NW, 2 == SE, 3 == SW
	LifeTicksLeft uint16
	TransitionTicks uint16
	TransitionElapsed uint16
}

func (self *Particle) Update() {
	// update movement
	switch self.Dir {
	case 0: // NE
		self.X += self.Speed
		self.Y -= self.Speed
	case 1: // NW
		self.X -= self.Speed
		self.Y -= self.Speed
	case 2: // SE
		self.X += self.Speed
		self.Y += self.Speed
	case 3: // SW
		self.X -= self.Speed
		self.Y += self.Speed
	default:
		panic("broken code")
	}

	// update state
	if self.LifeTicksLeft > 0 {
		if self.TransitionElapsed < self.TransitionTicks {
			self.TransitionElapsed += 1 // fade in
		} else {
			self.LifeTicksLeft -= 1 // life
			if self.LifeTicksLeft == 0 {
				self.TransitionElapsed = 0
			}
		}
	} else { // fade out
		self.TransitionElapsed += 1
	}

	// if dead, reroll
	if self.LifeTicksLeft > 0 || self.TransitionElapsed < self.TransitionTicks {
		return
	}
	self.X = rand.Float64()*640
	self.Y = rand.Float64()*360
	self.Speed = 0.02 + rand.Float64()*0.05
	self.Dir = uint8(rand.Intn(4))
	self.LifeTicksLeft = 60 + uint16(rand.Intn(160))
	self.TransitionTicks = 30 + uint16(rand.Intn(90))
	self.TransitionElapsed = 0
}

func (self *Particle) DrawLogical(canvas *ebiten.Image, ctx *context.Context, opts *ebiten.DrawImageOptions) {
	var opacity float32
	if self.LifeTicksLeft > 0 {
		if self.TransitionElapsed >= self.TransitionTicks {
			opacity = 1.0
		} else { // fade in
			opacity = float32(self.TransitionElapsed)/float32(self.TransitionTicks)
		}
	} else { // fade out
		opacity = 1.0 - float32(self.TransitionElapsed)/float32(self.TransitionTicks)
	}

	//clr := color.RGBA{255, 184, 228, 255}
	clr := color.RGBA{255, 196, 219, 255}
	yFactor := 0.6 + 0.4*(max(min(self.Y, 360), 0)/360.0)
	opts.GeoM.Translate(self.X, self.Y)
	opts.ColorScale.ScaleWithColor(utils.WithAlpha(clr, uint8(0.9*opacity*255*float32(yFactor))))
	canvas.DrawImage(ctx.Gfxcore.BackParticles[self.Dir], opts)
	opts.GeoM.Translate(-self.X, -self.Y)
	opts.ColorScale.Reset()
}

