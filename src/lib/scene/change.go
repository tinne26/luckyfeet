package scene

type ChangeType uint8
const (
	ChangePush    ChangeType = 1
	ChangePop     ChangeType = 2
	ChangeReplace ChangeType = 3
	// TODO: what do we do when we reached a scene that was
	//       already in the stack and we kinda want to "bubble it"?
	//       May happen in complex menus, missing that case.
	//       Should we always bubble? Have a method to ask whether
	//       bubbling is ok? Etc.
)

type Change struct {
	SceneKey Key
	Operation ChangeType
	// TODO: I may need input params for the new scene. tricky
}

func PushTo(key Key) *Change {
	return &Change{ SceneKey: key, Operation: ChangePush }
}

func Pop() *Change {
	return &Change{ Operation: ChangePop }
}

func ReplaceTo(key Key) *Change {
	return &Change{ SceneKey: key, Operation: ChangeReplace }
}
