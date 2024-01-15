package input

// The [GamepadConfig] allows the user to modify the mappings, triggers,
// repeat rates and direction buttons for a [Gamepad]. This object can
// be obtained through the [Gamepad.Config]() method.
type GamepadConfig Gamepad

func (self *GamepadConfig) SetLayouts(knownLayouts map[GamepadGUID]GamepadLayout) {
	self.knownLayouts = knownLayouts

	// refresh gamepad status
	if len(self.gamepadIds) == 0 { return }
	(*Gamepad)(self).tryConfigureNewID(self.gamepadIds[0])
}

// Equivalent to [*GamepadConfig.MapTriggerAction](action, nil).
func (self *GamepadConfig) DeleteTriggerAction(action TriggerAction) {
	self.MapTriggerAction(action, nil)
}

func (self *GamepadConfig) DeleteAllTriggerActions() {
	for action, _ := range self.mapping {
		delete(self.mapping, action)
		delete(self.actionAccTicks, action)
	}
}

// Returns nil if not mapped.
func (self *GamepadConfig) GetMapping(action TriggerAction) GamepadTrigger {
	return self.mapping[action]
}

// Sets a trigger.
func (self *GamepadConfig) MapTriggerAction(action TriggerAction, trigger GamepadTrigger) {
	if trigger == nil {
		delete(self.mapping, action)
		delete(self.actionAccTicks, action)
	} else {
		self.mapping[action] = trigger
	}
}

func (self *GamepadConfig) MapTriggerActionToButton(action TriggerAction, button GamepadStandardInput) {
	self.MapTriggerAction(action, SingleButton(button))
}

func (self *GamepadConfig) MapTriggerActionToButtons(action TriggerAction, buttons ...GamepadStandardInput) {
	// create multi-button trigger
	if len(buttons) == 0 { panic("len(buttons) == 0") }
	trigger := MultiButton{
		ButtonInputs: make([]ButtonInput, len(buttons)),
	}
	for i, button := range buttons {
		trigger.ButtonInputs[i] = ButtonInput{
			Button: button,
			Mode: InputModeNormal,
		}
	}

	// map trigger action
	self.MapTriggerAction(action, trigger)
}

func (self *GamepadConfig) SetRepeatRate(repeatFirst int32, repeatNext int32) {
   self.repeatFirst = repeatFirst
   self.repeatNext  = repeatNext
}

func (self *GamepadConfig) SetDirButtons(dirButtons GamepadDirButtons) {
	self.dirButtons = dirButtons
}
