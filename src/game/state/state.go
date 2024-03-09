package state

import "github.com/tinne26/luckyfeet/src/game/material/level"

type State[Context any] struct {
	LoadMapDataFromClipboard bool
	PlaytestData string
	PlaytestMapID uint8
	LevelKey level.Key
	Editing bool
	LastClearTicks int
}

func New[Context any]() *State[Context] {
	return &State[Context]{}
}
