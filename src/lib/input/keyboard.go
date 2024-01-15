package input

var _ KBGPLike = (*Keyboard)(nil)

type Keyboard struct {
	dirTriggers KeyboardDirTriggers // up/right/down/left indexing
	mapping map[TriggerAction]KeyboardTrigger
	actionAccTicks map[TriggerAction]int32
	dirAccTicks [4]int32
	dirUnifiedTicks int32
	dirUnifiedTrigger bool
	isIdle bool
	repeatFirst int32
	repeatNext int32
	cachedDir Direction
	cachedDir8 Direction
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		dirTriggers: DirKeysArrows(),
		mapping: make(map[TriggerAction]KeyboardTrigger, 8),
		actionAccTicks: make(map[TriggerAction]int32, 8),
		repeatFirst: pkgRepeatFirst,
		repeatNext: pkgRepeatNext,
	}
}

func (self *Keyboard) IsIdle() bool {
	return self.isIdle
}

func (self *Keyboard) PressedTicks(action TriggerAction) int32 {
   return self.actionAccTicks[action]
}

func (self *Keyboard) Pressed(action TriggerAction) bool {
   return self.PressedTicks(action) > 0
}

func (self *Keyboard) Trigger(action TriggerAction) bool {
   return self.PressedTicks(action) == 1
}

func (self *Keyboard) Repeat(action TriggerAction) bool {
	return self.isRepeatTickCount(self.PressedTicks(action))
}

func (self *Keyboard) RepeatAs(action TriggerAction, repeatFirst, repeatNext int32) bool {
	return isRepeatTickCount(self.PressedTicks(action), repeatFirst, repeatNext)
}

func (self *Keyboard) RepeatDir() Direction {
	if self.isRepeatTickCount(self.dirUnifiedTicks) {
		return self.cachedDir
	}
	return DirNone
}

func (self *Keyboard) RepeatDirAs(repeatFirst, repeatNext int32) Direction {
	if isRepeatTickCount(self.dirUnifiedTicks, repeatFirst, repeatNext) {
		return self.cachedDir
	}
	return DirNone
}

func (self *Keyboard) isRepeatTickCount(count int32) bool {
	return isRepeatTickCount(count, self.repeatFirst, self.repeatNext)
}

func (self *Keyboard) TriggerDir() Direction {
	if self.dirUnifiedTrigger {
		return self.cachedDir
	}
	return DirNone
}

func (self *Keyboard) TriggerDir8() Direction {
	if self.dirUnifiedTrigger {
		return self.cachedDir8
	}
	return DirNone
}

func (self *Keyboard) Dir() Direction {
	return self.cachedDir
}

func (self *Keyboard) Dir8() Direction {
	return self.cachedDir8
}

func (self *Keyboard) HorzDir() Direction {
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
func (self *Keyboard) Unwind() {
   for i, _ := range self.actionAccTicks {
      self.actionAccTicks[i] = -1
   }
   for i, _ := range self.dirAccTicks {
      self.dirAccTicks[i] = -1
   }
	self.dirUnifiedTicks = 0
	self.dirUnifiedTrigger = false
}

func (self *Keyboard) Zero() {
	for i, _ := range self.actionAccTicks {
		self.actionAccTicks[i] = 0
	}
	for i, _ := range self.dirAccTicks {
		self.dirAccTicks[i] = 0
	}
	self.dirUnifiedTicks = 0
	self.dirUnifiedTrigger = false
}

func (self *Keyboard) Update() error {
	self.isIdle = true

	// update actions
	for action, trigger := range self.mapping {
		if trigger.Pressed() {
			if self.actionAccTicks[action] < 0 { continue } // blocked
			self.actionAccTicks[action] += 1
			self.isIdle = false
		} else {
			if self.actionAccTicks[action] != 0 {
				self.isIdle = false
				self.actionAccTicks[action] = 0
			}
		}
	}

	// update directions
	newestDirIndex := -1
	newestDirValue := int32(0x7FFFFFFF)
	newestDirIndex, newestDirValue = self.updateDirTrigger(0, self.dirTriggers.Up   , newestDirIndex, newestDirValue)
	newestDirIndex, newestDirValue = self.updateDirTrigger(1, self.dirTriggers.Right, newestDirIndex, newestDirValue)
	newestDirIndex, newestDirValue = self.updateDirTrigger(2, self.dirTriggers.Down , newestDirIndex, newestDirValue)
	newestDirIndex, newestDirValue = self.updateDirTrigger(3, self.dirTriggers.Left , newestDirIndex, newestDirValue)

	if newestDirIndex != -1 {
		self.cachedDir  = Direction(1 << newestDirIndex)
		prevDir := self.cachedDir8
		self.cachedDir8 = ticksToDir8(self.dirAccTicks)
		self.dirUnifiedTicks += 1
		self.dirUnifiedTrigger = (prevDir != self.cachedDir8)
	} else {
		self.cachedDir  = 0
		self.cachedDir8 = 0
		self.dirUnifiedTicks = 0
		self.dirUnifiedTrigger = false
	}
	
	return nil
}

func (self *Keyboard) updateDirTrigger(i int, trigger KeyboardTrigger, newestDirIndex int, newestDirValue int32) (int, int32) {
	if trigger.Pressed() {
		if self.dirAccTicks[i] < 0 { return newestDirIndex, newestDirValue } // blocked
		self.isIdle = false
		self.dirAccTicks[i] += 1
		if self.dirAccTicks[i] < newestDirValue {
			newestDirIndex = i
			newestDirValue = self.dirAccTicks[i]
		}
	} else if self.dirAccTicks[i] != 0 {
		self.isIdle = false
		self.dirAccTicks[i] = 0
	}
	return newestDirIndex, newestDirValue
}

func (self *Keyboard) Config() *KeyboardConfig {
	return (*KeyboardConfig)(self)
}
