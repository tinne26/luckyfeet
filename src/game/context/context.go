package context

import "io/fs"

import "github.com/tinne26/luckyfeet/src/lib/input"
import "github.com/tinne26/luckyfeet/src/lib/audio"
import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/game/state"
import "github.com/tinne26/luckyfeet/src/game/graphics/gfxcore"
import "github.com/tinne26/luckyfeet/src/game/settings"
import "github.com/tinne26/luckyfeet/src/game/material/in"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/game/material/animations"
import "github.com/tinne26/luckyfeet/src/game/interfaces"

type Context struct {
	Input *input.KBGP
	Audio *audio.Soundscape
	Scenes *scene.Manager[*Context]
	State *state.State[*Context]
	Gfxcore *gfxcore.Graphics
	Settings *settings.Settings

	// --- extra random half hardcoded stuff ---
	Background interfaces.Background[*Context] // initialized on Start scene
	Animations *animations.Animations
}

func New(filesys fs.FS, sceneManager *scene.Manager[*Context]) (*Context, error) {
	var err error

	// set up audio context
	soundscape := audio.NewSoundscape()
	err = au.LoadAndConfigure(soundscape, filesys)
	if err != nil { return nil, err }
	
	// set up input context
	kbgp := input.NewKBGP()
	err = in.LoadAndConfigure(kbgp)
	if err != nil { return nil, err }

	// create new game state
	gameState := state.New[*Context]()

	// create new settings
	prefs := settings.New()

	// load graphics
	graphics, err := gfxcore.New(filesys)
	if err != nil { return nil, err }

	// load animations
	anims, err := animations.New(filesys)
	if err != nil { return nil, err }
	
	// return new context
	return &Context{
		Input: kbgp,
		State: gameState,
		Audio: soundscape,
		Scenes: sceneManager,
		Gfxcore: graphics,
		Settings: prefs,
		Animations: anims,
	}, nil
}

// Notice: scene manager update not included.
func (self *Context) UpdateSystems() error {
	var err error

	err = self.Audio.Update()
	if err != nil { return err }
	err = self.Input.Update()
	if err != nil { return err }

	return nil
}
