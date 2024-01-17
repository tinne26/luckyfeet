package audio

import "time"
import "math/rand"
import "github.com/hajimehoshi/ebiten/v2/audio"

type SfxPlayer struct {
	sources [][]byte
	volumeCorrectorFactor float32
	minBackoff time.Duration
	lastPlayed time.Time
}

func NewSfxPlayer(bytes ...[]byte) *SfxPlayer {
	return &SfxPlayer{ sources: bytes, volumeCorrectorFactor: 1.0 }
}

func (self *SfxPlayer) SetBackoff(duration time.Duration) {
	self.minBackoff = duration
}

func (self *SfxPlayer) SetVolumeCorrectorFactor(factor float32) {
	if factor < 0 { panic("factor < 0") }
	if factor > 1 { panic("factor > 1") }
	self.volumeCorrectorFactor = factor
}

func (self *SfxPlayer) Play() {
	self.PlayWithVolume(1.0)
}

// Typically the volume is masterVolume*sfxVolume.
func (self *SfxPlayer) PlayWithVolume(volume float32) {
	// backoff logic
	now := time.Now()
	if self.minBackoff > now.Sub(self.lastPlayed) { return }
	self.lastPlayed = now

	// play from pool of sfxs
	index := 0
	if len(self.sources) > 1 { index = rand.Intn(len(self.sources)) }
	player := audio.CurrentContext().NewPlayerFromBytes(self.sources[index])
	player.SetVolume(float64(volume*self.volumeCorrectorFactor))
	player.Play()
}
