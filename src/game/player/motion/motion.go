package motion

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/lib/audio"

import "github.com/tinne26/luckyfeet/src/game/material/au"

type SfxKey uint8
const (
	SfxNone SfxKey = iota
	SfxStep
	SfxLowStep
	//SfxJump
	//SfxDeath
)

type Animation struct {
	name string
	frameIndex uint8
	loopIndex uint8
	frames []*ebiten.Image
	sfxs []SfxKey
	frameDurations []uint8
	frameDurationLeft uint8
}

func NewAnimation(name string) *Animation {
	return &Animation{ name: name }
}

func (self *Animation) Name() string {
	return self.name
}

func (self *Animation) AddFrame(frame *ebiten.Image, frameTicks uint8) {
	self.AddFrameWithSfx(frame, frameTicks, SfxNone)
}

func (self *Animation) AddFrameWithSfx(frame *ebiten.Image, frameTicks uint8, sfxKey SfxKey) {
	self.frames = append(self.frames, frame)
	self.sfxs = append(self.sfxs, sfxKey)
	if frameTicks == 0 { panic("can't add frame with duration of 0 ticks") }
	if len(self.frameDurations) == 0 { self.frameDurationLeft = frameTicks }
	self.frameDurations = append(self.frameDurations, frameTicks)
}

func (self *Animation) FrameTicksElapsed() uint8 {
	return self.frameDurations[self.frameIndex] - self.frameDurationLeft
}

func (self *Animation) GetCurrentFrame() *ebiten.Image {
	return self.frames[self.frameIndex]
}

func (self *Animation) InPreLoopPhase() bool {
	return self.frameIndex < self.loopIndex
}

func (self *Animation) SkipIntro(soundscape *audio.Soundscape) {
	self.frameIndex = self.loopIndex
	self.frameDurationLeft = self.frameDurations[self.loopIndex]
	self.playSfx(soundscape)
}

func (self *Animation) Rewind(soundscape *audio.Soundscape) {
	self.frameIndex = 0
	self.frameDurationLeft = self.frameDurations[0]
	self.playSfx(soundscape)
}

func (self *Animation) RewindToLoop(soundscape *audio.Soundscape) {
	self.frameIndex = self.loopIndex
	self.frameDurationLeft = self.frameDurations[self.loopIndex]
	self.playSfx(soundscape)
}

func (self *Animation) Update(soundscape *audio.Soundscape) {
	self.frameDurationLeft -= 1
	if self.frameDurationLeft == 0 {
		if self.frameIndex == uint8(len(self.frames) - 1) {
			self.frameIndex = self.loopIndex
		} else {
			self.frameIndex += 1
		}
		self.playSfx(soundscape)
		self.frameDurationLeft = self.frameDurations[self.frameIndex]
	}
}

func (self *Animation) SetLoopStart(index uint8) {
	self.loopIndex = index
}

func (self *Animation) playSfx(soundscape *audio.Soundscape) {
	sfxKey := self.sfxs[self.frameIndex]
	switch sfxKey {
	case SfxNone:
		// nothing
	case SfxStep:
	 	soundscape.PlaySFX(au.SfxStep)
	case SfxLowStep:
		soundscape.PlaySFX(au.SfxLowStep)
	// case SfxJump:
	// 	soundscape.PlaySFX(au.SfxJump)
	// case SfxDeath:
	// 	soundscape.PlaySFX(au.SfxDeath)
	default:
		panic(sfxKey)
	}
}
