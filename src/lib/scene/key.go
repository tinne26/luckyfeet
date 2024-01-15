package scene

import "strconv"

type Key uint16
func (self Key) String() string {
	return strconv.Itoa(int(self))
}
