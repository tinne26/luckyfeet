package au

//import "time"
import "io"
import "io/fs"
import "strings"

import "github.com/tinne26/edau"
import "github.com/hajimehoshi/ebiten/v2/audio/vorbis"
import "github.com/hajimehoshi/ebiten/v2/audio/wav"

import "github.com/tinne26/luckyfeet/src/lib/audio"

// Prioritarily called from main.go so we make sure we are only doing it once.
// If you need to refer to the sample rate ever again, use audio.SampleRate().
func Initialize() {
	audio.Initialize(44100)
}

func IsContextReady() bool {
	return audio.IsContextReady()
}

func LoadAndConfigure(soundscape *audio.Soundscape, filesys fs.FS) error {
	var err error
	err = LoadSFXs(soundscape, filesys)
	if err != nil { return err }
	err = LoadBGMs(soundscape, filesys)
	if err != nil { return err }

	// set volumes after having created bgms and sfxs
	// so the changes can take effect
	soundscape.SetUserBGMVolume(0.5)
	soundscape.SetUserSFXVolume(0.5)

	return nil
}

func LoadSFXs(soundscape *audio.Soundscape, filesys fs.FS) error {
	var err error
	var sfx *audio.SfxPlayer
	
	// load sfxs
	sfx, err = loadWavMultiSFX(filesys, "assets/audio/sfx/step*.wav", '1', '4')
	if err != nil { return err }
	//sfx.SetBackoff(time.Millisecond*85)
	sfx.SetVolumeCorrectorFactor(0.23)
	SfxStep = soundscape.RegisterSFX(sfx)

	sfx.SetVolumeCorrectorFactor(0.1)
	SfxLowStep = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavMultiSFX(filesys, "assets/audio/sfx/click*.wav", '1', '3')
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.6)
	SfxClick = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavMultiSFX(filesys, "assets/audio/sfx/jump*.wav", '1', '4')
	if err != nil { return err }
	SfxJump = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSFX(filesys, "assets/audio/sfx/confirm.wav")
	if err != nil { return err }
	SfxConfirm = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSFX(filesys, "assets/audio/sfx/scratch.wav")
	if err != nil { return err }
	SfxScratch = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSFX(filesys, "assets/audio/sfx/cronch.wav")
	if err != nil { return err }
	SfxCronch = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSFX(filesys, "assets/audio/sfx/back.wav")
	if err != nil { return err }
	SfxBack = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSFX(filesys, "assets/audio/sfx/land.wav")
	if err != nil { return err }
	SfxLand = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSFX(filesys, "assets/audio/sfx/tictac.wav")
	if err != nil { return err }
	SfxTicTac = soundscape.RegisterSFX(sfx)

	return nil
}

func LoadBGMs(soundscape *audio.Soundscape, filesys fs.FS) error {
	// load and set up bgms
	loop, err := loadLooper(filesys, "assets/audio/bgm/unlucky.ogg", 35822, 6769917)
	if err != nil { return err }
	bgm := audio.NewBgmFromLooper(loop)
	bgm.SetVolumeCorrectorFactor(0.5)
	BgmMain = soundscape.RegisterBGM(bgm)
	return nil
}

// ---- helper methods ----

func loadLooper(filesys fs.FS, filename string, loopStartSample, loopEndSample int64) (*edau.Looper, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	file, err := filesys.Open(filename)
	if err != nil { return nil, err }
	stream, err := vorbis.DecodeWithSampleRate(sampleRate, file)
	if err != nil { return nil, err }
	return edau.NewLooper(stream, loopStartSample*4, loopEndSample*4), nil
}

func loadWavMultiSFX(filesys fs.FS, filename string, r1, r2 byte) (*audio.SfxPlayer, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	if r2 <= r1 { panic("r1 <= r2") }
	if r2 - r1 > 26 { panic("r2 too far away from r1") }
	index := strings.IndexRune(filename, '*')
	if index == -1 { panic("wildcard '*' not found") }
	bytes := make([][]byte, 0, 1 + r2 - r1)

	filenameBytes := []byte(filename)
	for r := r1; r <= r2; r++ {
		filenameBytes[index] = r
		file, err := filesys.Open(string(filenameBytes))
		if err != nil { return nil, err }
		stream, err := wav.DecodeWithSampleRate(sampleRate, file)
		if err != nil { return nil, err }
		audioBytes, err := io.ReadAll(stream)
		if err != nil { return nil, err }
		bytes = append(bytes, audioBytes)
	}
	return audio.NewSfxPlayer(bytes...), nil
}

func loadWavSFX(filesys fs.FS, filename string) (*audio.SfxPlayer, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	file, err := filesys.Open(filename)
	if err != nil { return nil, err }
	stream, err := wav.DecodeWithSampleRate(sampleRate, file)
	if err != nil { return nil, err }
	audioBytes, err := io.ReadAll(stream)
	if err != nil { return nil, err }
	return audio.NewSfxPlayer(audioBytes), nil
}
