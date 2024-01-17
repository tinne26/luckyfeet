package audio

import "time"
import "slices"

type Soundscape struct {
	automationPanel *AutomationPanel

	userVolumeMaster float32
	userVolumeSFX float32
	userVolumeBGM float32
	muteMaster bool
	muteSFX bool
	muteBGM bool

	sfxs []*SfxPlayer
	bgms []*BGM
	activeBGM *BGM
	fadingOutBGMs []*BGM
}

func NewSoundscape() *Soundscape {
	return &Soundscape{
		userVolumeMaster: 0.7,
		userVolumeBGM: 0.7,
		userVolumeSFX: 0.7,
		sfxs: make([]*SfxPlayer, 0, 16),
		bgms: make([]*BGM, 0, 8),
		automationPanel: NewAutomationPanel(),
	}
}

func (self *Soundscape) SetMasterMuted(muted bool) {
	if muted == self.muteMaster { return }
	self.muteMaster = muted
	if self.activeBGM != nil {
		self.activeBGM.SetMuted(self.muteBGM || self.muteMaster)
	}
	for _, bgm := range self.fadingOutBGMs {
		bgm.SetMuted(self.muteBGM || self.muteMaster)
	}
}
func (self *Soundscape) SetBGMMuted(muted bool) {
	if muted == self.muteBGM { return }
	self.muteBGM = muted
	if self.activeBGM != nil {
		self.activeBGM.SetMuted(self.muteBGM || self.muteMaster)
	}
	for _, bgm := range self.fadingOutBGMs {
		bgm.SetMuted(self.muteBGM || self.muteMaster)
	}
}
func (self *Soundscape) SetSFXMuted(muted bool) { self.muteSFX = muted }
func (self *Soundscape) GetMasterMuted() bool { return self.muteMaster }
func (self *Soundscape) GetBGMMuted() bool { return self.muteBGM }
func (self *Soundscape) GetSFXMuted() bool { return self.muteSFX }

// Notice: SFXs volume will not change until next play if already playing.
func (self *Soundscape) SetMasterVolume(volume float32) {
	if volume < 0 { panic("volume < 0") }
	if volume > 1 { panic("volume > 1") }
	self.userVolumeMaster = volume

	if self.activeBGM != nil {
		self.activeBGM.SetUserVolume(self.userVolumeMaster)
	}
	for _, bgm := range self.fadingOutBGMs {
		bgm.SetMasterVolume(self.userVolumeMaster)
	}
}

func (self *Soundscape) SetUserSFXVolume(volume float32) {
	if volume < 0 { panic("volume < 0") }
	if volume > 1 { panic("volume > 1") }
	self.userVolumeSFX = volume
}

func (self *Soundscape) GetMasterVolume() float32 {
	return self.userVolumeMaster
}

func (self *Soundscape) GetUserSFXVolume() float32 {
	return self.userVolumeSFX
}

func (self *Soundscape) SetUserBGMVolume(volume float32) {
	if volume < 0 { panic("volume < 0") }
	if volume > 1 { panic("volume > 1") }
	self.userVolumeBGM = volume

	if self.activeBGM != nil {
		self.activeBGM.SetUserVolume(self.userVolumeBGM)
	}
	for _, bgm := range self.fadingOutBGMs {
		bgm.SetUserVolume(self.userVolumeBGM)
	}
}

func (self *Soundscape) GetUserBGMVolume() float32 {
	return self.userVolumeBGM
}

func (self *Soundscape) RegisterSFX(sfx *SfxPlayer) SfxKey {
	key := SfxKey(len(self.sfxs))
	self.sfxs = append(self.sfxs, sfx)
	return key
}

func (self *Soundscape) PlaySFX(key SfxKey) {
	if self.muteSFX || self.muteMaster { return }
	self.sfxs[key].PlayWithVolume(self.userVolumeMaster*self.userVolumeSFX)
}

func (self *Soundscape) RegisterBGM(bgm *BGM) BgmKey {
	self.refreshVolumeBGM(bgm)
	key := BgmKey(len(self.bgms))
	self.bgms = append(self.bgms, bgm)
	return key
}

func (self *Soundscape) IsActive(key BgmKey) bool {
	return self.activeBGM == self.bgms[key]
}

func (self *Soundscape) refreshVolumeBGM(bgm *BGM) {
	bgm.SetMasterVolume(self.userVolumeMaster)
	bgm.SetUserVolume(self.userVolumeBGM)
	bgm.SetMuted(self.muteBGM || self.muteMaster)
}

func (self *Soundscape) addToFadingOutBGMs(fadingBgm *BGM) {
	for _, bgm := range self.fadingOutBGMs {
		if bgm == fadingBgm { return }
	}
	self.refreshVolumeBGM(fadingBgm)
	self.fadingOutBGMs = append(self.fadingOutBGMs)
}

func (self *Soundscape) removeFromFadingOutBGM(fadingBgm *BGM) {
	for i, bgm := range self.fadingOutBGMs {
		if bgm == fadingBgm {
			last := len(self.fadingOutBGMs) - 1
			self.fadingOutBGMs[i], self.fadingOutBGMs[last] = self.fadingOutBGMs[last], self.fadingOutBGMs[i]
			self.fadingOutBGMs = self.fadingOutBGMs[0 : last]
			return
		}
	}
}

func (self *Soundscape) FadeOut(fadeOut time.Duration) {
	if self.activeBGM != nil {
		self.activeBGM.FadeOut(fadeOut)
		self.addToFadingOutBGMs(self.activeBGM)
	}
	self.activeBGM = nil
}

func (self *Soundscape) FadeIn(key BgmKey, fadeOut, wait, fadeIn time.Duration) {
	if self.activeBGM != nil {
		self.activeBGM.FadeOut(fadeOut)
		self.addToFadingOutBGMs(self.activeBGM)
	}
	self.activeBGM = self.bgms[key]
	self.refreshVolumeBGM(self.activeBGM)
	self.removeFromFadingOutBGM(self.activeBGM)
	self.activeBGM.FadeIn(fadeOut + wait, fadeIn)
}

func (self *Soundscape) Crossfade(key BgmKey, fadeOut, inWait, fadeIn time.Duration) {
	if self.activeBGM != nil {
		self.activeBGM.FadeOut(fadeOut)
	}
	self.activeBGM = self.bgms[key]
	self.refreshVolumeBGM(self.activeBGM)
	self.removeFromFadingOutBGM(self.activeBGM)
	self.activeBGM.FadeIn(inWait, fadeIn)
}

func (self *Soundscape) AutomationPanel() *AutomationPanel {
	return self.automationPanel
}

func (self *Soundscape) Update() error {
	// clean up faded out bgms
	self.fadingOutBGMs = slices.DeleteFunc(
		self.fadingOutBGMs,
		func(bgm *BGM) bool { return bgm.FullyFadedOut() },
	)

	// ...

	return nil
}
