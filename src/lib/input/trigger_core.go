package input

type TriggerAction uint16
//const NoTriggerAction TriggerAction = 0

type InputMode uint8
const (
	InputModeNormal  InputMode = 0
	InputModeBreaker InputMode = 1 // like a negation
	InputModeOr      InputMode = 2 // one affirmative evaluation is enough
	// TODO: InputModeModifier could be used depending on the order it's evaluated
)
