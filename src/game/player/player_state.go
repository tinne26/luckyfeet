package player

type State uint8
const (
	StIdle State = iota
	StRunning
	StJumpingHold
	StJumpingInertial
	StTicTacHold
	StTicTacInertial
	StFalling
)

func (self State) String() string {
	switch self {
	case StIdle: return "Idle"
	case StRunning: return "Running"
	case StJumpingHold: return "Jumping (Hold)"
	case StJumpingInertial: return "Jumping (Inertial)"
	case StTicTacHold: return "Tic-Tac (Hold)"
	case StTicTacInertial: return "Tic-Tac (Inertial)"
	case StFalling: return "Falling"
	default:
		return "Unknown State"
	}
}
