package input

import "github.com/hajimehoshi/ebiten/v2"

var _ KBGPLike = (*Gamepad)(nil)

// Limitations:
// - Only the first / "oldest" gamepad is used. The code could easily be
//   adapted, but the UIs and the model are not so trivial to adjust.
// - We are not really making use of analog values and joysticks yet.
// - Vibration not exposed through anywhere, though that feels like it
//   should be on a separate place, as it's not really input but feedback.
// - Haven't created a scene for new gamepad interruption + config yet.
// - Would be great to be able to export all the gamepad stuff to an
//   external library, constants, screens, helpers, dat format, etc.

type Gamepad struct {
	gamepadIds []ebiten.GamepadID
	dirButtons GamepadDirButtons // up/right/down/left indexing
	mapping map[TriggerAction]GamepadTrigger
	status GamepadStatus
	knownLayouts map[GamepadGUID]GamepadLayout
	currentLayout GamepadLayout
	
	actionAccTicks map[TriggerAction]int32
	dirAccTicks [4]int32
	dirUnifiedTicks int32 // so changing dir while repeating doesn't alter the rhythm
	repeatFirst int32
	repeatNext int32
	cachedDir Direction
	cachedDir8 Direction
	dirUnifiedTrigger bool
	isLayoutSet bool
	isIdle bool
}

func NewGamepad() *Gamepad {
	return &Gamepad{
		dirButtons: DirButtonsDPad(),
		mapping: make(map[TriggerAction]GamepadTrigger, 8),
		actionAccTicks: make(map[TriggerAction]int32, 8),
		repeatFirst: pkgRepeatFirst,
		repeatNext: pkgRepeatNext,
	}
}

func (self *Gamepad) IsIdle() bool {
	return self.isIdle
}

func (self *Gamepad) Status() GamepadStatus {
	return self.status
}

func (self *Gamepad) PressedTicks(action TriggerAction) int32 {
   return self.actionAccTicks[action]
}

func (self *Gamepad) Pressed(action TriggerAction) bool {
   return self.PressedTicks(action) > 0
}

func (self *Gamepad) Trigger(action TriggerAction) bool {
   return self.PressedTicks(action) == 1
}

func (self *Gamepad) Repeat(action TriggerAction) bool {
	return self.isRepeatTickCount(self.PressedTicks(action))
}

func (self *Gamepad) RepeatAs(action TriggerAction, repeatFirst, repeatNext int32) bool {
	return isRepeatTickCount(self.PressedTicks(action), repeatFirst, repeatNext)
}

func (self *Gamepad) RepeatDir() Direction {
	if self.isRepeatTickCount(self.dirUnifiedTicks) {
		return self.cachedDir
	}
	return DirNone
}

func (self *Gamepad) RepeatDirAs(repeatFirst, repeatNext int32) Direction {
	if isRepeatTickCount(self.dirUnifiedTicks, repeatFirst, repeatNext) {
		return self.cachedDir
	}
	return DirNone
}

func (self *Gamepad) isRepeatTickCount(count int32) bool {
	return isRepeatTickCount(count, self.repeatFirst, self.repeatNext)
}

func (self *Gamepad) TriggerDir() Direction {
	if self.dirUnifiedTrigger {
		return self.cachedDir
	}
	return DirNone
}

func (self *Gamepad) TriggerDir8() Direction {
	if self.dirUnifiedTrigger {
		return self.cachedDir8
	}
	return DirNone
}

func (self *Gamepad) Dir() Direction {
	return self.cachedDir
}

func (self *Gamepad) Dir8() Direction {
	return self.cachedDir8
}

func (self *Gamepad) HorzDir() Direction {
	right, left := self.dirAccTicks[1], self.dirAccTicks[3]
	if right <= 0 {
		if left > 0 { return DirLeft }
		return DirNone
	} else if left <= 0 {
		return DirRight
	} else { // both are positive
		if right <= left {
			return DirRight
		}
		return DirLeft
	}
}

// All inputs are blocked (won't be triggered) until a
// subsequent update determines that actions are not pressed
// anymore. This is very helpful for scene transitions
// or context switches where you don't want previous state
// to affect the new context.
func (self *Gamepad) Unwind() {
   for i, _ := range self.actionAccTicks {
      self.actionAccTicks[i] = -1
   }
   for i, _ := range self.dirAccTicks {
      self.dirAccTicks[i] = -1
   }
	self.dirUnifiedTicks = 0
	self.dirUnifiedTrigger = false
}

func (self *Gamepad) Zero() {
	for i, _ := range self.actionAccTicks {
      self.actionAccTicks[i] = 0
   }
   for i, _ := range self.dirAccTicks {
      self.dirAccTicks[i] = 0
   }
	self.dirUnifiedTicks = 0
	self.dirUnifiedTrigger = false
}

func (self *Gamepad) Update() error {
	self.isIdle = true

	// detect active id change
	var hadGamepadID bool = len(self.gamepadIds) > 0
	var prevGamepadID ebiten.GamepadID
	if hadGamepadID {
		prevGamepadID = self.gamepadIds[0]
	}
	self.gamepadIds = ebiten.AppendGamepadIDs(self.gamepadIds)

	// update gamepad status
	if len(self.gamepadIds) == 0 {
		self.status &= ^gamepadStatusConnectedFlag // mark disconnected
		if hadGamepadID { self.Zero() }
		return nil // no gamepads available, halt update
	} else {
		self.status |= gamepadStatusConnectedFlag // mark connected
		if !hadGamepadID || prevGamepadID != self.gamepadIds[0] {
			self.isIdle = false
			self.status |= gamepadStatusChangedFlag
			self.status &= ^gamepadStatusConfiguredFlag
			self.tryConfigureNewID(self.gamepadIds[0])
		} else {
			self.status &= ^gamepadStatusChangedFlag
		}
	}

	// get current layout (nil should use std layouts as fallback)
	layout := &self.currentLayout
	if !self.isLayoutSet { layout = nil }

	// update actions
	id := self.gamepadIds[0]
	for action, trigger := range self.mapping {
		if trigger.Pressed(id, layout) {
			if self.actionAccTicks[action] < 0 { continue } // blocked
			self.actionAccTicks[action] += 1
			self.isIdle = false
		} else if self.actionAccTicks[action] != 0 {
			self.isIdle = false
			self.actionAccTicks[action] = 0
		}
	}

	// update directions
	newestDirIndex := -1
	newestDirValue := int32(0x7FFFFFFF)
	if self.isLayoutSet {
		newestDirIndex, newestDirValue = self.updateDirWithLayout(0, id, self.dirButtons.Up   , newestDirIndex, newestDirValue)
		newestDirIndex, newestDirValue = self.updateDirWithLayout(1, id, self.dirButtons.Right, newestDirIndex, newestDirValue)
		newestDirIndex, newestDirValue = self.updateDirWithLayout(2, id, self.dirButtons.Down , newestDirIndex, newestDirValue)
		newestDirIndex, newestDirValue = self.updateDirWithLayout(3, id, self.dirButtons.Left , newestDirIndex, newestDirValue)
	} else {
		newestDirIndex, newestDirValue = self.updateDirWithoutLayout(0, id, self.dirButtons.Up   , newestDirIndex, newestDirValue)
		newestDirIndex, newestDirValue = self.updateDirWithoutLayout(1, id, self.dirButtons.Right, newestDirIndex, newestDirValue)
		newestDirIndex, newestDirValue = self.updateDirWithoutLayout(2, id, self.dirButtons.Down , newestDirIndex, newestDirValue)
		newestDirIndex, newestDirValue = self.updateDirWithoutLayout(3, id, self.dirButtons.Left , newestDirIndex, newestDirValue)
	}
	
	if newestDirIndex != -1 {
		self.cachedDir = Direction(1 << newestDirIndex)
		prevDir := self.cachedDir8
		self.cachedDir8 = ticksToDir8(self.dirAccTicks)
		self.dirUnifiedTicks += 1
		self.dirUnifiedTrigger = (prevDir != self.cachedDir8)
	} else {
		self.cachedDir = 0
		self.cachedDir8 = 0
		self.dirUnifiedTicks = 0
		self.dirUnifiedTrigger = false
	}
	
	return nil
}

func (self *Gamepad) updateDirWithLayout(i int, id ebiten.GamepadID, button GamepadStandardInput, newestDirIndex int, newestDirValue int32) (int, int32) {
	if self.currentLayout[button].Pressed(id) {
		if self.dirAccTicks[i] < 0 { return newestDirIndex, newestDirValue } // blocked
		self.dirAccTicks[i] += 1
		self.isIdle = false
		if self.dirAccTicks[i] < newestDirValue {
			return i, self.dirAccTicks[i]
		}
	} else if self.dirAccTicks[i] != 0 {
		self.isIdle = false
		self.dirAccTicks[i] = 0
	}
	return newestDirIndex, newestDirValue
}

// std layouts fallback
func (self *Gamepad) updateDirWithoutLayout(i int, id ebiten.GamepadID, button GamepadStandardInput, newestDirIndex int, newestDirValue int32) (int, int32) {
	if button.StdEquivalentPressed(id) {
		if self.dirAccTicks[i] < 0 { return newestDirIndex, newestDirValue } // blocked
		self.dirAccTicks[i] += 1
		self.isIdle = false
		if self.dirAccTicks[i] < newestDirValue {
			return i, self.dirAccTicks[i]
		}
	} else if self.dirAccTicks[i] != 0 {
		self.isIdle = false
		self.dirAccTicks[i] = 0
	}
	return newestDirIndex, newestDirValue
}

func (self *Gamepad) Config() *GamepadConfig {
	return (*GamepadConfig)(self)
}
