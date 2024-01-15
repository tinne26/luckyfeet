package settings

// Game settings that the user can change. They may affect 
// the user, a savefile, apply globally, etc. It depends.

type Settings struct {
	// audio
	MusicLevel uint8 // uses 0 - 100 range
	SfxLevel uint8   // uses 0 - 100 range
	MusicMuted bool
	SfxMuted bool

	// screen and resolution
	ScreenFit ScreenFitMode
	AllowWinResize bool // ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// debug and performance
	ShowFPS bool
}

func New() *Settings {
	return &Settings{
		MusicLevel: 50,
		SfxLevel: 50,
	}
}
