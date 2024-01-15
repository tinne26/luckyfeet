package briefblackout

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/scene/keys"
import "github.com/tinne26/luckyfeet/src/lib/scene"

func SceneKey() scene.Key { return keys.BriefBlackout }
func RequestHook(ctx *context.Context) scene.Scene[*context.Context] {
	// ctx.Input.Unwind() // do not unwind input here, we don't want to be interfering
	black, err := New(ctx)
	if err != nil { panic(err) }
	return black
}
