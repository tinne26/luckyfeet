package state

import "fmt"

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

func (self *State[Context]) GetClearTimeStr() string {
	return self.FmtTicksToTimeStr(self.LastClearTicks)
}

// This shouldn't be here, but I'm late.
func (self *State[Context]) FmtTicksToTimeStr(ticks int) string {
	secs := ticks/120
	mins := secs/60
	secs -= mins*60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}
