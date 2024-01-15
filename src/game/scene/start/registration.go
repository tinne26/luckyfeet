package start

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/material/scene/keys"
import "github.com/tinne26/luckyfeet/src/lib/scene"

func SceneKey() scene.Key { return keys.Start }
func RequestHook(ctx *context.Context) scene.Scene[*context.Context] {
	ctx.Input.Unwind()
	start, err := New(ctx)
	if err != nil { panic(err) }
	return start
}
