package input

type GamepadStatus uint8
const (
	gamepadStatusConnectedFlag    GamepadStatus = 0b0000_0001
	gamepadStatusChangedFlag      GamepadStatus = 0b0000_0010
	gamepadStatusConfiguredFlag   GamepadStatus = 0b0000_0100
)

func (self GamepadStatus) IsConnected() bool {
	return self & gamepadStatusConnectedFlag != 0
}

func (self GamepadStatus) IsDisconnected() bool {
	return self & gamepadStatusConnectedFlag == 0
}

func (self GamepadStatus) IsJustConnected() bool {
	return self.IsConnected() && (self & gamepadStatusChangedFlag != 0)
}

func (self GamepadStatus) IsJustDisconnected() bool {
	return self.IsDisconnected() && (self & gamepadStatusChangedFlag != 0)
}

func (self GamepadStatus) IsConfigured() bool {
	return self.IsConnected() && (self & gamepadStatusConfiguredFlag != 0)
}
