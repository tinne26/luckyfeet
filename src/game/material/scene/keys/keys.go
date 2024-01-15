package keys

import "github.com/tinne26/luckyfeet/src/lib/scene"

const (
	Start scene.Key = iota + 1
	Editor
	Play
	BriefBlackout
	WinScreen
	// ...
)

func FirstSceneKey() scene.Key {
	// note: this is a method so it's easy to add
	//       program arg detection and so on for
	//       quickstarts and testing
	return Start
}
