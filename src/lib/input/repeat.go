package input

const pkgRepeatFirst = 17
const pkgRepeatNext  = 6

func isRepeatTickCount(tick, repeatFirst, repeatNext int32) bool {
	if tick <= 0 { return false }
	if tick == 1 || tick == repeatFirst { return true }
	if tick < repeatFirst { return false }
	return ((tick - repeatFirst) % repeatNext) == 0
}

// NOTE: I don't think generics are relevant for this right now.
type RepeatDetectorKBGP struct {
	First int32
	Next int32
}

func (self RepeatDetectorKBGP) Repeat(input *KBGP, action TriggerAction) bool {
	return input.RepeatAs(action, self.First, self.Next)
}

func (self RepeatDetectorKBGP) RepeatDir(input *KBGP) Direction {
	return input.RepeatDirAs(self.First, self.Next)
}
