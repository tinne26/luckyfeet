package input

import "io"
import "fmt"
import "errors"

import "github.com/hajimehoshi/ebiten/v2"

// there are two types of values: discrete and analog/continuous.
// let's call it binary and continuous? boolean or floating?
// that's better. floating values are axes. boolean values are buttons.
// problem with floating values is that they can go from -1 to 1 or
// from 0 to 1.

type GamepadStandardInput uint8
const (
	GamepadUp GamepadStandardInput = iota
	GamepadRight
	GamepadLeft
	GamepadDown

	GamepadButtonTop
	GamepadButtonRight
	GamepadButtonLeft
	GamepadButtonBottom

	GamepadShoulderLeft
	GamepadShoulderRight

	GamepadTriggerLeft
	GamepadTriggerRight

	GamepadLeftStickButton
	GamepadRightStickButton
	GamepadLeftStickHorzAxis
	GamepadLeftStickVertAxis
	GamepadRightStickHorzAxis
	GamepadRightStickVertAxis

	// TODO: maybe use GamepadFuncRight or OptLeft? weird though. it's also "options" on dualshock4
	// also missing dualShock4 touchpad button. could determine by num buttons.
	// also main stick vs secondary stick instead of left and right?
	GamepadStart
	GamepadSelect
	GamepadMeta

	gamepadNumStandardInputs // must always be last
)

// (must be kept in sync with the GamepadStandardInput constants above)
var ebitengineStdEquivalents = [gamepadNumStandardInputs]ebiten.StandardGamepadButton{
	ebiten.StandardGamepadButtonLeftTop, // GamepadUp GamepadStandardInput = iota
	ebiten.StandardGamepadButtonLeftRight, // GamepadRight
	ebiten.StandardGamepadButtonLeftLeft, // GamepadLeft
	ebiten.StandardGamepadButtonLeftBottom, // GamepadDown

	ebiten.StandardGamepadButtonRightTop, // GamepadButtonTop
	ebiten.StandardGamepadButtonRightRight, // GamepadButtonRight
	ebiten.StandardGamepadButtonRightLeft, // GamepadButtonLeft
	ebiten.StandardGamepadButtonRightBottom, // GamepadButtonBottom

	ebiten.StandardGamepadButtonFrontTopLeft, // GamepadShoulderLeft
	ebiten.StandardGamepadButtonFrontTopRight, // GamepadShoulderRight

	ebiten.StandardGamepadButtonFrontBottomLeft, // GamepadTriggerLeft
	ebiten.StandardGamepadButtonFrontBottomRight, // GamepadTriggerRight

	ebiten.StandardGamepadButtonLeftStick, // GamepadLeftStickButton
	ebiten.StandardGamepadButtonRightStick, // GamepadRightStickButton
	-1, // GamepadLeftStickHorzAxis
	-1, // GamepadLeftStickVertAxis
	-1, // GamepadRightStickHorzAxis
	-1, // GamepadRightStickVertAxis

	ebiten.StandardGamepadButtonCenterRight, // GamepadStart
	ebiten.StandardGamepadButtonCenterLeft, // GamepadSelect
	ebiten.StandardGamepadButtonCenterCenter, // GamepadMeta
}

const GamepadNumStandardInputs = int(gamepadNumStandardInputs)

func (self GamepadStandardInput) String() string {
	switch self {
	case GamepadUp: return "GamepadUp"
	case GamepadRight: return "GamepadRight"
	case GamepadLeft: return "GamepadLeft"
	case GamepadDown: return "GamepadDown"
	case GamepadButtonTop: return "GamepadButtonTop"
	case GamepadButtonRight: return "GamepadButtonRight"
	case GamepadButtonLeft: return "GamepadButtonLeft"
	case GamepadButtonBottom: return "GamepadButtonBottom"
	case GamepadShoulderLeft: return "GamepadShoulderLeft"
	case GamepadShoulderRight: return "GamepadShoulderRight"
	case GamepadTriggerLeft: return "GamepadTriggerLeft"
	case GamepadTriggerRight: return "GamepadTriggerRight"
	case GamepadLeftStickButton: return "GamepadLeftStickButton"
	case GamepadRightStickButton: return "GamepadRightStickButton"
	case GamepadLeftStickHorzAxis: return "GamepadLeftStickHorzAxis"
	case GamepadLeftStickVertAxis: return "GamepadLeftStickVertAxis"
	case GamepadRightStickHorzAxis: return "GamepadRightStickHorzAxis"
	case GamepadRightStickVertAxis: return "GamepadRightStickVertAxis"
	case GamepadStart: return "GamepadStart"
	case GamepadSelect: return "GamepadSelect"
	case GamepadMeta: return "GamepadMeta"
	default:
		return fmt.Sprintf("GamepadUnknown#%d", self)
	}
}

func (self GamepadStandardInput) StdEquivalent() ebiten.StandardGamepadButton {
	return ebitengineStdEquivalents[self]
}

func (self GamepadStandardInput) StdEquivalentPressed(id ebiten.GamepadID) bool {
	return ebiten.IsStandardGamepadButtonPressed(id, self.StdEquivalent())
}

func (self GamepadStandardInput) IsAxis() bool {
	switch self {
	case GamepadLeftStickHorzAxis: return true
	case GamepadLeftStickVertAxis: return true
	case GamepadRightStickHorzAxis: return true
	case GamepadRightStickVertAxis: return true
	default:
		return false
	}
}

type GamepadLayoutInput uint8
const (
	gamepadLayoutInputButtonBit GamepadLayoutInput = 0b1000_0000
	gamepadLayoutInputAnalogBit GamepadLayoutInput = 0b0100_0000
	gamepadLayoutInputIndexBits GamepadLayoutInput = 0b0011_1111
)

func (self GamepadLayoutInput) String() string {
	index := self.Index()
	if index == -1 { return "Undefined" }

	if gamepadLayoutInputButtonBit & self != 0 {
		if gamepadLayoutInputAnalogBit & self != 0 {
			return fmt.Sprintf("AnalogButton#%d", index)
		} else {
			return fmt.Sprintf("BinaryButton#%d", index)
		}
	} else {
		return fmt.Sprintf("Axis#%d", index)
	}
}

// The index will be -1 if undefined.
func (self GamepadLayoutInput) Index() int {
	return int(self & gamepadLayoutInputIndexBits) - 1
}

func NewGamepadBinaryButtonLayoutInput(index int) GamepadLayoutInput {
	if index < 0 || index >= 63 { panic("invalid index") }
	return gamepadLayoutInputButtonBit | GamepadLayoutInput(index + 1)
}

func NewGamepadAnalogButtonLayoutInput(index int) GamepadLayoutInput {
	return NewGamepadBinaryButtonLayoutInput(index) | gamepadLayoutInputAnalogBit
}

func NewGamepadAxisLayoutInput(index int) GamepadLayoutInput {
	if index < 0 || index >= 63 { panic("invalid index") }
	return gamepadLayoutInputAnalogBit | GamepadLayoutInput(index + 1)
}

func (self GamepadLayoutInput) Value(id ebiten.GamepadID) float64 {
	index := self.Index()
	if index == -1 { return 0.0 }
	if self.IsAnalog() {
		return ebiten.GamepadAxisValue(id, index)
	} else {
		if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(index)) { return 1.0 }
		return 0.0
	}
}

func (self GamepadLayoutInput) Pressed(id ebiten.GamepadID) bool {
	const PressThreshold = 0.2

	index := self.Index()
	if index == -1 { return false }
	if self.IsAnalog() {
		value := ebiten.GamepadAxisValue(id, index)
		return value >= PressThreshold || value <= -PressThreshold
	} else {
		return ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(index))
	}
}

func (self GamepadLayoutInput) IsButton() bool {
	return self & gamepadLayoutInputButtonBit != 0
}
func (self GamepadLayoutInput) IsBinary() bool {
	return self & gamepadLayoutInputAnalogBit == 0
}
func (self GamepadLayoutInput) IsAnalog() bool {
	return self & gamepadLayoutInputAnalogBit != 0
}

type GamepadLayout [GamepadNumStandardInputs]GamepadLayoutInput

func (self *GamepadLayout) Reset() {
	for i := 0; i < len(self); i++ {
		self[i] = 0
	}
}

// Slow.
func (self *GamepadLayout) AxisCorrespondence(index int) (GamepadStandardInput, bool) {
	for i := 0; i < len(self); i++ {
		in := self[i]
		if in.Index() == index && in.IsAnalog() {
			return GamepadStandardInput(i), true
		}
	}
	return 0, false
}

// Slow.
func (self *GamepadLayout) ButtonCorrespondence(index int) (GamepadStandardInput, bool) {
	for i := 0; i < len(self); i++ {
		in := self[i]
		if in.Index() == index && in.IsBinary() {
			return GamepadStandardInput(i), true
		}
	}
	return 0, false
}

var gamepadLayoutMagicBytes = []byte{'o', 'b', 'g', 'p', 'd', 'l', 'y', 't'}
const lenGamepadMagicBytes = 8 // len(gamepadLayoutMagicBytes)
const layoutDatLen = lenGamepadMagicBytes + 16 + GamepadNumStandardInputs
const GamepadLayoutDatSize = layoutDatLen

func init() {
	if lenGamepadMagicBytes != len(gamepadLayoutMagicBytes) { panic("broken constant") }
}

// TODO: could consider adding the gamepad model / type too, but general xbox-like
//       seems fine, or it could be detected from the layout defined and undefined
//       inputs.
func (self *GamepadLayout) Export(writer io.Writer, guid GamepadGUID) error {
	bytes := make([]byte, layoutDatLen)
	copy(bytes, gamepadLayoutMagicBytes)
	guid.WriteToBuffer(bytes[lenGamepadMagicBytes : ])
	mappingBytes := bytes[lenGamepadMagicBytes + 16 : ]
	for i := 0; i < GamepadNumStandardInputs; i++ {
		mappingBytes[i] = byte(self[i])
	}
	n, err := writer.Write(bytes)
	if err == nil && n != len(bytes) { err = io.ErrShortWrite }
	return err
}

func ParseGamepadLayout(data []byte) (GamepadGUID, GamepadLayout, error) {
	var layout GamepadLayout
	if len(data) != layoutDatLen {
		err := fmt.Errorf("gamepad layout data must be %d bytes", layoutDatLen)
		return GamepadGUID{}, layout, err
	}

	// parse signature
	for i := 0; i < lenGamepadMagicBytes; i++ {
		if data[i] != gamepadLayoutMagicBytes[i] {
			return GamepadGUID{}, layout, errors.New("invalid signature")
		}
	}
	data = data[lenGamepadMagicBytes : ]

	// parse gamepad guid
	guid, err := BytesToGamepadGUID(data[ : 16])
	if err != nil { return guid, layout, errors.New("invalid guid: " + err.Error()) }
	data = data[16 : ]

	// parse actual layout mapping
	for i := 0; i < GamepadNumStandardInputs; i++ {
		layout[i] = GamepadLayoutInput(data[i])
	}
	return guid, layout, nil
}

func StandardGamepadLayout() GamepadLayout {
	return GamepadLayout([GamepadNumStandardInputs]GamepadLayoutInput{
		gamepadLayoutInputButtonBit | 12, // GamepadUp
		gamepadLayoutInputButtonBit | 15, // GamepadRight
		gamepadLayoutInputButtonBit | 14, // GamepadLeft
		gamepadLayoutInputButtonBit | 13, // GamepadDown
		gamepadLayoutInputButtonBit |  3, // GamepadButtonTop
		gamepadLayoutInputButtonBit |  1, // GamepadButtonRight
		gamepadLayoutInputButtonBit |  2, // GamepadButtonLeft
		gamepadLayoutInputButtonBit |  0, // GamepadButtonBottom

		gamepadLayoutInputButtonBit |  4, // GamepadShoulderLeft
		gamepadLayoutInputButtonBit |  5, // GamepadShoulderRight

		gamepadLayoutInputButtonBit |  6, // GamepadTriggerLeft
		gamepadLayoutInputButtonBit |  7, // GamepadTriggerRight

		gamepadLayoutInputButtonBit | 10, // GamepadLeftStickButton
		gamepadLayoutInputButtonBit | 11, // GamepadRightStickButton
		gamepadLayoutInputAnalogBit |  0, // GamepadLeftStickHorzAxis
		gamepadLayoutInputAnalogBit |  1, // GamepadLeftStickVertAxis
		gamepadLayoutInputAnalogBit |  2, // GamepadRightStickHorzAxis
		gamepadLayoutInputAnalogBit |  3, // GamepadRightStickVertAxis
		
		gamepadLayoutInputButtonBit |  9, // GamepadStart
		gamepadLayoutInputButtonBit |  8, // GamepadSelect
		gamepadLayoutInputButtonBit | 16, // GamepadMeta
	})
}

func (self *Gamepad) tryConfigureNewID(id ebiten.GamepadID) {
	self.isLayoutSet = false
	self.status &= ^gamepadStatusConfiguredFlag
	sdlid := ebiten.GamepadSDLID(id)
	guid, err := StringToGamepadGUID(sdlid)
	if err != nil { panic(err) } // TODO: maybe link a logger... or return as update() err?

	if self.knownLayouts == nil { return }
	layout, hasLayout := self.knownLayouts[guid]
	if !hasLayout { return }
	self.currentLayout = layout
	self.status |= gamepadStatusConfiguredFlag
}
