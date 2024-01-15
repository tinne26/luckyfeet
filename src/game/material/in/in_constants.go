package in

import "github.com/tinne26/luckyfeet/src/lib/input"

// Trigger repeat constants (RF = repeat first, RN = repeat next)
const RFDefault = 17*2
const RFFast = 11*2
const RNDefault = 6*2
const RNSlow = 10*2
const RNFast = 4*2

// Basically, aliases to provide better access to 'input' constants.
type Direction = input.Direction
const (
	DirNone  input.Direction = input.DirNone
	DirUp    input.Direction = input.DirUp
	DirRight input.Direction = input.DirRight
	DirDown  input.Direction = input.DirDown
	DirLeft  input.Direction = input.DirLeft
)
