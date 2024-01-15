package registry

import "github.com/tinne26/luckyfeet/src/lib/scene"
import "github.com/tinne26/luckyfeet/src/game/context"

import "github.com/tinne26/luckyfeet/src/game/scene/start"
import "github.com/tinne26/luckyfeet/src/game/scene/editor"
import "github.com/tinne26/luckyfeet/src/game/scene/play"
import "github.com/tinne26/luckyfeet/src/game/scene/briefblackout"
import "github.com/tinne26/luckyfeet/src/game/scene/winscreen"

func NewSceneManager() *scene.Manager[*context.Context] {
	registry := scene.NewRegistry[*context.Context]()
	
	// Register all scenes so they can be used with the scene manager.
	registry.Register(start.RequestHook, start.SceneKey())
	registry.Register(editor.RequestHook, editor.SceneKey())
	registry.Register(play.RequestHook, play.SceneKey())
	registry.Register(briefblackout.RequestHook, briefblackout.SceneKey())
	registry.Register(winscreen.RequestHook, winscreen.SceneKey())

	return scene.NewManager(registry)
}

