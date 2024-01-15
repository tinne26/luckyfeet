package scene

import "errors"

type Registry[Context any] struct {
	scenes map[Key]func(Context) Scene[Context]
}

func NewRegistry[Context any]() *Registry[Context] {
	return &Registry[Context]{
		scenes: make(map[Key]func(Context) Scene[Context], 4),
	}
}

func (self *Registry[Context]) Register(request func(Context) Scene[Context], key Key) error {
	_, alreadyRegistered := self.scenes[key]
	if alreadyRegistered {
		return errors.New("scene with key '" + key.String() + "' already registered")
	}
	if request == nil {
		panic("can't register nil scene request function")
	}
	self.scenes[key] = request
	return nil
}

// May return nil if the given key is not registered.
func (self *Registry[Context]) Load(key Key, ctx Context) Scene[Context] {
	fn, found := self.scenes[key]
	if !found { return nil }
	return fn(ctx)
}

// Panics if scene is not found.
func (self *Registry[Context]) MustLoad(key Key, state Context) Scene[Context] {
	scene := self.Load(key, state)
	if scene == nil {
		panic("scene with key '" + key.String() + "' not found")
	}
	return scene
}
