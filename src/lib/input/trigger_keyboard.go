package input

import "strings"

import "github.com/hajimehoshi/ebiten/v2"

type KeyboardTrigger interface {
	Pressed() bool
	Inputs() []KeyInput
	String() string
}

type KeyInput struct {
	Key ebiten.Key
	Mode InputMode
}

// ---- UnassignedKey ----
var _ KeyboardTrigger = UnassignedKey{}
type UnassignedKey struct{}

func (self UnassignedKey) Pressed() bool {
	return false
}
func (self UnassignedKey) String() string {
	return "UnassignedKey"
}
func (self UnassignedKey) Inputs() []KeyInput {
	return []KeyInput{}
}

// ---- SingleKey ----
var _ KeyboardTrigger = SingleKey(0)
type SingleKey ebiten.Key

func (self SingleKey) Pressed() bool {
   return ebiten.IsKeyPressed(ebiten.Key(self))
}
func (self SingleKey) String() string {
	return strings.ToTitle(ebiten.KeyName(ebiten.Key(self)))
}
func (self SingleKey) Inputs() []KeyInput {
	return []KeyInput{
		KeyInput{
			Key: ebiten.Key(self),
			Mode: InputModeNormal,
		},
	}
}

// ---- KeyList ----
var _ KeyboardTrigger = KeyList{}
type KeyList []ebiten.Key

func NewKeyList(keys ...ebiten.Key) KeyList {
	return KeyList(keys)
}

func (self KeyList) Pressed() bool {
   for _, key := range self {
		if ebiten.IsKeyPressed(key) { return true }
	}
	return false
}
func (self KeyList) String() string {
	var str strings.Builder
	for i, key := range self {
		str.WriteString(strings.ToTitle(ebiten.KeyName(key)))
		if i != 0 { str.Write([]byte{'|'}) }
	}
	return str.String()
}
func (self KeyList) Inputs() []KeyInput {
	var inputs []KeyInput = make([]KeyInput, len(self))
	for i, key := range self {
		inputs[i] = KeyInput{ Key: key, Mode: InputModeOr }
	}
	return inputs
}

// ---- MultiKey ----
var _ KeyboardTrigger = MultiKey{}
type MultiKey struct {
	KeyInputs []KeyInput
}

func NewMultiKey(keys ...ebiten.Key) MultiKey {
	multi := MultiKey{ KeyInputs: make([]KeyInput, len(keys)) }
	for i, key := range keys {
		multi.KeyInputs[i] = KeyInput{
			Key: key,
			Mode: InputModeNormal,
		}
	}
	return multi
}

func (self *MultiKey) AddNormalKey(key ebiten.Key) {
	self.KeyInputs = append(self.KeyInputs, KeyInput{ Key: key, Mode: InputModeNormal })
}
func (self *MultiKey) AddBreakerKey(key ebiten.Key) {
	self.KeyInputs = append(self.KeyInputs, KeyInput{ Key: key, Mode: InputModeBreaker })
}

func (self MultiKey) Pressed() bool {
	for i := 0; i < len(self.KeyInputs); i++ {
		switch self.KeyInputs[i].Mode {
		case InputModeNormal:
			if !ebiten.IsKeyPressed(self.KeyInputs[i].Key) {
				return false
			}
		case InputModeBreaker:
			if ebiten.IsKeyPressed(self.KeyInputs[i].Key) {
				return false
			}
		default:
			panic("broken code")
		}
	}
	return true
}
func (self MultiKey) String() string {
	var str strings.Builder
	for i, in := range self.KeyInputs {
		if i != 0 { str.Write([]byte{' ', '+', ' '}) }
		if in.Mode == InputModeBreaker { str.Write([]byte{'!'}) }
		keyStr := ebiten.KeyName(in.Key)
		str.WriteString(strings.ToTitle(keyStr))
	}
	return str.String()
}
func (self MultiKey) Inputs() []KeyInput {
	return self.KeyInputs
}
