package audio

import "github.com/hajimehoshi/ebiten/v2/audio"

// Can only be called once.
func Initialize(sampleRate int) {
	if audio.CurrentContext() != nil {
		panic("audio context already exists")
	}
	audio.NewContext(sampleRate)
}

func SampleRate() int {
	return audio.CurrentContext().SampleRate()
}

func IsContextReady() bool {
	ctx := audio.CurrentContext()
	if ctx == nil { return false }
	return ctx.IsReady()
}

func CurrentContext() *audio.Context {
	return audio.CurrentContext()
}
