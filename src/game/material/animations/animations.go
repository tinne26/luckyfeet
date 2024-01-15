package animations

import "io/fs"
import "image"
import "image/png"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game/player/motion"

type Animations struct {
	Idle *motion.Animation
	Running *motion.Animation
	InAir *motion.Animation
}

func New(filesys fs.FS) (*Animations, error) {
	// load mc animations image
	const CreaturesPath = "assets/graphics/creatures/"
	mc, err := loadImage(filesys, CreaturesPath + "mc.png")
	if err != nil { return nil, err }

	anims := &Animations{}
	anims.Idle = motion.NewAnimation("idle")
	idle1 := frame(mc, 0, 0)
	idle2 := frame(mc, 2, 0)
	idleFeet := frame(mc, 1, 0)
	anims.Idle.AddFrame(idle1, 160)
	anims.Idle.AddFrame(idle2, 160)
	anims.Idle.AddFrame(idle1, 160)
	anims.Idle.AddFrame(idle2, 160)
	anims.Idle.AddFrame(idle1, 160)
	anims.Idle.AddFrame(idle2, 160)
	anims.Idle.AddFrame(idle1, 160)
	anims.Idle.AddFrame(idle2, 160)
	anims.Idle.AddFrame(idle1, 160)
	anims.Idle.AddFrame(idle2, 160)
	anims.Idle.AddFrame(idle1, 160)
	anims.Idle.AddFrame(idleFeet, 100)
	anims.Idle.AddFrame(idle1, 160)
	anims.Idle.AddFrame(idle2, 160)

	anims.Running = motion.NewAnimation("running")
	walk2Run := frame(mc, 3, 0)
	run1 := frame(mc, 0, 1)
	run2 := frame(mc, 1, 1)
	run3 := frame(mc, 2, 1)
	run4 := frame(mc, 3, 1)
	run5 := frame(mc, 4, 1)
	anims.Running.AddFrameWithSfx(walk2Run, 10, motion.SfxLowStep)
	anims.Running.AddFrameWithSfx(run1, 18, motion.SfxStep)
	anims.Running.AddFrameWithSfx(run2, 18, motion.SfxLowStep)
	anims.Running.AddFrameWithSfx(run3, 18, motion.SfxStep)
	anims.Running.AddFrameWithSfx(run4, 18, motion.SfxLowStep)
	anims.Running.AddFrameWithSfx(run5, 18, motion.SfxStep)
	anims.Running.SetLoopStart(1)

	anims.InAir = motion.NewAnimation("air")
	air1 := frame(mc, 0, 2)
	anims.InAir.AddFrame(air1, 255)

	return anims, nil
}

func loadImage(filesys fs.FS, filepath string) (*ebiten.Image, error) {
	file, err := filesys.Open(filepath)
	if err != nil { return nil, err }
	img, err := png.Decode(file)
	if err != nil { return nil, err }
	return ebiten.NewImageFromImage(img), nil
}

func frame(img *ebiten.Image, col, row int) *ebiten.Image {
	const FrameWidth, FrameHeight = 15, 35
	ox, oy := FrameWidth*col, FrameHeight*row
	rect := image.Rect(ox, oy, ox + FrameWidth, oy + FrameHeight)
	return img.SubImage(rect).(*ebiten.Image)
}
