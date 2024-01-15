package input

// TODO: I think "prioritary directions" are not well integrated with button combos.
//       But we can work around that to some extent, as long as we only allow customizing
//       directions and button hold combos separately...?

type Direction uint8
const (
	DirNone   Direction = 0b00000000
	DirUp     Direction = 0b00000001
	DirRight  Direction = 0b00000010
	DirDown   Direction = 0b00000100
	DirLeft   Direction = 0b00001000
	dirUnused Direction = 0b11110000
)

func (self Direction) Vert() Direction {
	return self & (DirUp | DirDown)
}

func (self Direction) Horz() Direction {
	return self & (DirRight | DirLeft)
}

func (self Direction) String() string {
	if self == DirNone { return "DirNone" }
	var str string = "Dir"
	if self & DirUp    != 0 { str += "Up"    }
	if self & DirDown  != 0 { str += "Down"  }
	if self & DirRight != 0 { str += "Right" }
	if self & DirLeft  != 0 { str += "Left"  }
	if self & dirUnused != 0 {
		str += "+" + uint8toBinaryStr(uint8(self & dirUnused))
	}
	return str
}

// ---- helpers ----

func ticksToDir8(ticks [4]int32) Direction {
	var dir Direction
	
	// vertical
	upTicks, downTicks := ticks[0], ticks[2]
	if upTicks > 0 || downTicks > 0 {
		if upTicks <= downTicks {
			dir |= DirUp
		} else {
			dir |= DirDown
		}
	}

	// horizontal
	rightTicks, leftTicks := ticks[1], ticks[3]
	if rightTicks > 0 || leftTicks > 0 {
		if rightTicks <= leftTicks {
			dir |= DirRight
		} else {
			dir |= DirLeft
		}
	}

	return dir
}
