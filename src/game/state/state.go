package state

type State[Context any] struct {
	LoadMapDataFromClipboard bool
	PlaytestData string
	PlaytestMapID uint8
	Editing bool
	LastClearTicks int
}

func New[Context any]() *State[Context] {
	return &State[Context]{}
}
