package input

import "strings"

import "github.com/hajimehoshi/ebiten/v2"

type GamepadTrigger interface {
	Pressed(id ebiten.GamepadID, layout *GamepadLayout) bool
	Inputs() []ButtonInput
	String() string
}

type ButtonInput struct {
	Button GamepadStandardInput
	Mode InputMode
}

// ---- UnassignedButton ----
var _ GamepadTrigger = UnassignedButton{}
type UnassignedButton struct{}

func (self UnassignedButton) Pressed(id ebiten.GamepadID, layout *GamepadLayout) bool {
   return false
}
func (self UnassignedButton) String() string {
	return "GamepadUnassignedInput"
}
func (self UnassignedButton) Inputs() []ButtonInput {
	return []ButtonInput{}
}

// ---- SingleButton ----
var _ GamepadTrigger = SingleButton(0)
type SingleButton GamepadStandardInput

func (self SingleButton) Pressed(id ebiten.GamepadID, layout *GamepadLayout) bool {
	if layout != nil { return layout[self].Pressed(id) }
	return GamepadStandardInput(self).StdEquivalentPressed(id)
}
func (self SingleButton) String() string {
	return self.String()
}
func (self SingleButton) Inputs() []ButtonInput {
	return []ButtonInput{
		ButtonInput{
			Button: GamepadStandardInput(self),
			Mode: InputModeNormal,
		},
	}
}

// ---- MultiButton ----
var _ GamepadTrigger = MultiButton{}
type MultiButton struct {
	ButtonInputs []ButtonInput
}

func NewMultiButton(buttons ...GamepadStandardInput) MultiButton {
	multi := MultiButton{ ButtonInputs: make([]ButtonInput, len(buttons)) }
	for i, button := range buttons {
		multi.ButtonInputs[i] = ButtonInput{
			Button: button,
			Mode: InputModeNormal,
		}
	}
	return multi
}

func (self *MultiButton) Add(button GamepadStandardInput) {
	self.ButtonInputs = append(self.ButtonInputs, ButtonInput{ Button: button, Mode: InputModeNormal })
}
func (self *MultiButton) AddBreaker(button GamepadStandardInput) {
	self.ButtonInputs = append(self.ButtonInputs, ButtonInput{ Button: button, Mode: InputModeBreaker })
}

func (self MultiButton) Pressed(id ebiten.GamepadID, layout *GamepadLayout) bool {
	if layout != nil {
		for i := 0; i < len(self.ButtonInputs); i++ {
			pressed := layout[self.ButtonInputs[i].Button].Pressed(id)
			switch self.ButtonInputs[i].Mode {
			case InputModeNormal  : if !pressed { return false }
			case InputModeBreaker : if  pressed { return false }
			default:
				panic("broken code")
			}
		}
	} else {
		for i := 0; i < len(self.ButtonInputs); i++ {
			pressed := self.ButtonInputs[i].Button.StdEquivalentPressed(id)
			switch self.ButtonInputs[i].Mode {
			case InputModeNormal  : if !pressed { return false }
			case InputModeBreaker : if  pressed { return false }
			default:
				panic("broken code")
			}
		}
	}
	
	return true
}
func (self MultiButton) String() string {
	var str strings.Builder
	for i, in := range self.ButtonInputs {
		if i != 0 { str.Write([]byte{' ', '+', ' '}) }
		if in.Mode == InputModeBreaker { str.Write([]byte{'!'}) }
		str.WriteString(in.Button.String())
	}
	return str.String()
}
func (self MultiButton) Inputs() []ButtonInput {
	return self.ButtonInputs
}
