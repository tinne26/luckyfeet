package input

// Struct for use with [*GamepadConfig.SetDirButtons]()
type GamepadDirButtons struct {
	Up GamepadStandardInput
	Down GamepadStandardInput
	Right GamepadStandardInput
	Left GamepadStandardInput
}

func DirButtonsDPad() GamepadDirButtons {
	return GamepadDirButtons{
		Up: GamepadUp,
		Down: GamepadDown,
		Right: GamepadRight,
		Left: GamepadLeft,
	}
}
