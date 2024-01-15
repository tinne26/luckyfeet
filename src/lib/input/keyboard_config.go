package input

import "github.com/hajimehoshi/ebiten/v2"

type KeyboardConfig Keyboard

// TODO: instead of deleting, disabling may be better. this way I'd keep
//       maps more stable, though size isn't deallocated anyway, so it may
//       end up being the same on full remappings.

// Equivalent to [*KeyboardConfig.MapTriggerAction](action, nil).
func (self *KeyboardConfig) DeleteTriggerAction(action TriggerAction) {
	self.MapTriggerAction(action, nil)
}

func (self *KeyboardConfig) DeleteAllTriggerActions() {
	for action, _ := range self.mapping {
		delete(self.mapping, action)
		delete(self.actionAccTicks, action)
	}
}

// Returns nil if not mapped.
func (self *KeyboardConfig) GetMapping(action TriggerAction) KeyboardTrigger {
	return self.mapping[action]
}

// Sets a trigger.
func (self *KeyboardConfig) MapTriggerAction(action TriggerAction, trigger KeyboardTrigger) {
	if trigger == nil {
		delete(self.mapping, action)
		delete(self.actionAccTicks, action)
	} else {
		self.mapping[action] = trigger
	}
}

func (self *KeyboardConfig) MapTriggerActionToKey(action TriggerAction, key ebiten.Key) {
	self.MapTriggerAction(action, SingleKey(key))
}

func (self *KeyboardConfig) MapTriggerActionToKeys(action TriggerAction, keys []ebiten.Key) {
	// create multi-key trigger
	trigger := MultiKey{
		KeyInputs: make([]KeyInput, len(keys)),
	}
	for i, key := range keys {
		trigger.KeyInputs[i] = KeyInput{
			Key: key,
			Mode: InputModeNormal,
		}
	}

	// map trigger action
	self.MapTriggerAction(action, trigger)
}

func (self *KeyboardConfig) SetRepeatRate(repeatFirst int32, repeatNext int32) {
   self.repeatFirst = repeatFirst
   self.repeatNext  = repeatNext
}

func (self *KeyboardConfig) SetDirTriggers(dirTriggers KeyboardDirTriggers) {
	self.dirTriggers = dirTriggers
}
