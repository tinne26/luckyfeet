package input

var _ KBGPLike = (*KBGP)(nil)

type KBGP struct {
	keyboard *Keyboard
	gamepad *Gamepad
	activityStatus int8 // 0 = idle, -1 = more recent gamepad use, 1 = more recent keyboard use
}

func NewKBGP() *KBGP {
	return &KBGP{
		keyboard: NewKeyboard(),
		gamepad: NewGamepad(),
	}
}

func (self *KBGP) Gamepad()  *Gamepad  { return self.gamepad  }
func (self *KBGP) Keyboard() *Keyboard { return self.keyboard }

func (self *KBGP) UsedGamepadMoreRecentlyThanKeyboard() bool {
	return self.activityStatus == -1
}
func (self *KBGP) UsedKeyboardMoreRecentlyThanGamepad() bool {
	return self.activityStatus == 1
}

func (self *KBGP) PressedTicks(action TriggerAction) int32 {
   kbpf := self.keyboard.PressedTicks(action)
	gppf := self.gamepad.PressedTicks(action)
	return max(kbpf, gppf)
}

func (self *KBGP) Pressed(action TriggerAction) bool {
   return self.keyboard.Pressed(action) || self.gamepad.Pressed(action)
}

func (self *KBGP) Trigger(action TriggerAction) bool {
   return self.PressedTicks(action) == 1
}

func (self *KBGP) Repeat(action TriggerAction) bool {
	kbpf := self.keyboard.PressedTicks(action)
	gppf := self.gamepad.PressedTicks(action)
	if kbpf >= gppf {
		return self.keyboard.isRepeatTickCount(kbpf)
	} else {
		return self.gamepad.isRepeatTickCount(gppf)
	}
}

// See also repeat detector.
func (self *KBGP) RepeatAs(action TriggerAction, repeatFirst, repeatNext int32) bool {
	tick := max(self.keyboard.PressedTicks(action), self.gamepad.PressedTicks(action))
	return isRepeatTickCount(tick, repeatFirst, repeatNext)
}

// See also repeat detector.
func (self *KBGP) RepeatDirAs(repeatFirst, repeatNext int32) Direction {
	if self.shouldPrioritizeKeyboardDir() {
		if isRepeatTickCount(self.keyboard.dirUnifiedTicks, repeatFirst, repeatNext) {
			return self.keyboard.cachedDir
		}
	} else {
		if isRepeatTickCount(self.gamepad.dirUnifiedTicks, repeatFirst, repeatNext) {
			return self.gamepad.cachedDir
		}
	}
	
	return DirNone
}

func (self *KBGP) RepeatDir() Direction {
	if self.shouldPrioritizeKeyboardDir() {
		return self.keyboard.RepeatDir()
	} else {
		return self.gamepad.RepeatDir()
	}
}

func (self *KBGP) TriggerDir() Direction {
	gpTrigDir := self.gamepad.TriggerDir()
	kbTrigDir := self.keyboard.TriggerDir()
	if kbTrigDir == DirNone { return gpTrigDir }
	return kbTrigDir // keyboard has preference on ties
}

func (self *KBGP) TriggerDir8() Direction {
	gpTrigDir := self.gamepad.TriggerDir8()
	kbTrigDir := self.keyboard.TriggerDir8()
	if kbTrigDir == DirNone { return gpTrigDir }
	return kbTrigDir // keyboard has preference on ties
}

func (self *KBGP) Dir() Direction {
	if self.shouldPrioritizeKeyboardDir() {
		return self.keyboard.Dir()
	} else {
		return self.gamepad.Dir()
	}
}

func (self *KBGP) Dir8() Direction {
	if self.shouldPrioritizeKeyboardDir() {
		return self.keyboard.Dir8()
	} else {
		return self.gamepad.Dir8()
	}
}

func (self *KBGP) HorzDir() Direction {
	if self.shouldPrioritizeKeyboardDir() {
		return self.keyboard.HorzDir()
	} else {
		return self.gamepad.HorzDir()
	}
}

func (self *KBGP) shouldPrioritizeKeyboardDir() bool {
	if self.gamepad.dirUnifiedTicks  <= 0 { return true  }
	if self.keyboard.dirUnifiedTicks <= 0 { return false }
	return self.keyboard.dirUnifiedTicks <= self.gamepad.dirUnifiedTicks
}

// All inputs are blocked (won't be triggered) until a
// subsequent update determines that actions are not pressed
// anymore. This is very helpful for scene transitions
// or context switches where you don't want previous state
// to affect the new context.
func (self *KBGP) Unwind() {
   self.gamepad.Unwind()
	self.keyboard.Unwind()
}

// Sets all inputs to zero. Notice that this can easily
// cause accidental retriggering (see also Unwind).
func (self *KBGP) Zero() {
	self.gamepad.Zero()
	self.keyboard.Zero()
}

func (self *KBGP) Update() error {
	err := self.keyboard.Update()
	if err != nil { self.activityStatus = 0 ; return err }
	err = self.gamepad.Update()
	if err != nil { self.activityStatus = 0 ; return err }

	// update activity status (whether the most recent use is for keyboard or gamepad)
	kbIdle, gpIdle := self.keyboard.IsIdle(), self.gamepad.IsIdle()
	if kbIdle == gpIdle { return nil } // keep activity status as it was
	if kbIdle { self.activityStatus = -1 } else { self.activityStatus = 1 }
	return nil
}

func (self *KBGP) SetRepeatRate(repeatFirst int32, repeatNext int32) {
   self.gamepad.Config().SetRepeatRate(repeatFirst, repeatNext)
	self.keyboard.Config().SetRepeatRate(repeatFirst, repeatNext)
}
