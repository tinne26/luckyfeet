package input

// A common interface implemented by Gamepad, Keyboard and KBGP.
type KBGPLike interface {
	Update() error
	
	Zero()
	Unwind()
	
	Pressed(TriggerAction) bool
	Trigger(TriggerAction) bool
	Repeat(TriggerAction) bool
	RepeatAs(TriggerAction, int32, int32) bool
	PressedTicks(TriggerAction) int32
	
	Dir() Direction
	TriggerDir() Direction
	Dir8() Direction
	TriggerDir8() Direction
	RepeatDir() Direction
	RepeatDirAs(int32, int32) Direction
}
