package scene

import "github.com/hajimehoshi/ebiten/v2"

type Manager[Context any] struct {
	Scenes []Scene[Context]
	Registry *Registry[Context]
}

func NewManager[Context any](registry *Registry[Context]) *Manager[Context] {
	manager := &Manager[Context]{
		Scenes: make([]Scene[Context], 0, 4),
		Registry: registry,
	}
	return manager
}

func (self *Manager[Context]) FirstLoad(sceneKey Key, ctx Context) {
	if len(self.Scenes) != 0 { panic("scene manager already has scenes") }
	newScene := self.Registry.MustLoad(sceneKey, ctx)
	self.Scenes = append(self.Scenes, newScene)
	if len(self.Scenes) != 1 { panic("broken code") }
}

func (self *Manager[Context]) Current() Scene[Context] {
	if len(self.Scenes) == 0 { return nil }
	return self.Scenes[len(self.Scenes) - 1]
}

func (self *Manager[Context]) Update(ctx Context) error {
	// update non-active scenes
	for _, scene := range self.Scenes[0 : len(self.Scenes) - 1] {
		change, err := scene.Update(ctx)
		if err != nil { return err }
		if change != nil {
			panic("scene change requested from non-active scene")
		}
	}

	// update active scene
	activeIndex := len(self.Scenes) - 1
	activeScene := self.Scenes[activeIndex]
	change, err := activeScene.Update(ctx)
	if err != nil { return err }
	if change != nil {
		switch change.Operation {
		case ChangePush:
			newScene := self.Registry.MustLoad(change.SceneKey, ctx)
			self.Scenes = append(self.Scenes, newScene)
		case ChangePop:
			if activeIndex == 0 { panic("can't pop last remaining scene") }
			self.Scenes = self.Scenes[0 : activeIndex]
		case ChangeReplace:
			newScene := self.Registry.MustLoad(change.SceneKey, ctx)
			self.Scenes[activeIndex] = newScene
		default:
			panic(change.Operation)
		}
	}
	
	return nil
}

func (self *Manager[Context]) DrawLogical(canvas *ebiten.Image, ctx Context) {
	foremostIndex := len(self.Scenes) - 1
	for i, scene := range self.Scenes {
		scene.DrawLogical(canvas, (i == foremostIndex), ctx)
	}
}

func (self *Manager[Context]) DrawHiRes(canvas *ebiten.Image, ctx Context) {
	foremostIndex := len(self.Scenes) - 1
	for i, scene := range self.Scenes {
		scene.DrawHiRes(canvas, (i == foremostIndex), ctx)
	}
}
