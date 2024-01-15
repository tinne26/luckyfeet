package input

import "github.com/hajimehoshi/ebiten/v2"

// Struct for use with [*KeyboardConfig.SetDirTriggers]()
type KeyboardDirTriggers struct {
	Up KeyboardTrigger
	Down KeyboardTrigger
	Right KeyboardTrigger
	Left KeyboardTrigger
}

// Get an instance of [KeyboardDirTriggers] corresponding to WASD.
func DirKeysWASD() KeyboardDirTriggers {
	return KeyboardDirTriggers{
		Up: SingleKey(ebiten.KeyW),
		Down: SingleKey(ebiten.KeyS),
		Right: SingleKey(ebiten.KeyD),
		Left: SingleKey(ebiten.KeyA),
	}
}

// Get an instance of [KeyboardDirTriggers] corresponding to directional arrows.
func DirKeysArrows() KeyboardDirTriggers {
	return KeyboardDirTriggers{
		Up: SingleKey(ebiten.KeyArrowUp),
		Down: SingleKey(ebiten.KeyArrowDown),
		Right: SingleKey(ebiten.KeyArrowRight),
		Left: SingleKey(ebiten.KeyArrowLeft),
	}
}

// Get an instance of [KeyboardDirTriggers] corresponding to both WASD and directional arrows.
func DirKeysArrowsAndWASD() KeyboardDirTriggers {
	return KeyboardDirTriggers{
		Up: KeyList([]ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyW}),
		Down: KeyList([]ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyS}),
		Right: KeyList([]ebiten.Key{ebiten.KeyArrowRight, ebiten.KeyD}),
		Left: KeyList([]ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyA}),
	}
}
