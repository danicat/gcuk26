package audio

import (
	"bytes"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SampleRate     = 44100
	BytesPerSample = 4 // 16-bit stereo = 2 bytes * 2 channels
)

// AudioSystem manages runtime audio playback and procedurally synthesized sound effects.
type AudioSystem struct {
	context     *audio.Context
	bgmPlayer   *audio.Player
	jumpPCM     []byte
	iceCreamPCM []byte
	bottlePCM   []byte
	bellPCM     []byte
	boostPCM    []byte
	crashPCM    []byte
	bgmPCM      []byte
	muted       bool
}

// NewAudioSystem initializes Ebitengine audio context and synthesizes game sounds.
func NewAudioSystem() *AudioSystem {
	ctx := audio.NewContext(SampleRate)
	sys := &AudioSystem{
		context: ctx,
	}
	sys.initSounds()
	return sys
}

// initSounds synthesizes all sound effects and background music into memory buffers.
func (a *AudioSystem) initSounds() {
	a.jumpPCM = synthesizeJumpSound()
	a.iceCreamPCM = synthesizeIceCreamSound()
	a.bottlePCM = synthesizeBottleSound()
	a.bellPCM = synthesizeBellSound()
	a.boostPCM = synthesizeBoostSound()
	a.crashPCM = synthesizeCrashSound()
	a.bgmPCM = synthesizeLondonBGM()
}

// PlayJump plays the bunny-hop jump sound effect.
func (a *AudioSystem) PlayJump() {
	a.playSFX(a.jumpPCM)
}

// PlayIceCream plays the ice cream pick-up sound effect.
func (a *AudioSystem) PlayIceCream() {
	a.playSFX(a.iceCreamPCM)
}

// PlayBottle plays the plastic bottle recycling sound effect.
func (a *AudioSystem) PlayBottle() {
	a.playSFX(a.bottlePCM)
}

// PlayBell plays the bicycle bell ring sound effect.
func (a *AudioSystem) PlayBell() {
	a.playSFX(a.bellPCM)
}

// PlayBoost plays the turbo speed boost sound effect.
func (a *AudioSystem) PlayBoost() {
	a.playSFX(a.boostPCM)
}

// PlayCrash plays the hazard collision crash sound effect.
func (a *AudioSystem) PlayCrash() {
	a.playSFX(a.crashPCM)
}

// StartBGM begins playing the upbeat London background music loop.
func (a *AudioSystem) StartBGM() {
	if a.muted || a.context == nil || len(a.bgmPCM) == 0 {
		return
	}
	if a.bgmPlayer != nil && a.bgmPlayer.IsPlaying() {
		return
	}
	loopStream := audio.NewInfiniteLoop(bytes.NewReader(a.bgmPCM), int64(len(a.bgmPCM)))
	player, err := a.context.NewPlayer(loopStream)
	if err == nil {
		a.bgmPlayer = player
		a.bgmPlayer.SetVolume(0.35)
		a.bgmPlayer.Play()
	}
}

// StopBGM pauses background music playback.
func (a *AudioSystem) StopBGM() {
	if a.bgmPlayer != nil && a.bgmPlayer.IsPlaying() {
		a.bgmPlayer.Pause()
	}
}

func (a *AudioSystem) playSFX(pcm []byte) {
	if a.muted || a.context == nil || len(pcm) == 0 {
		return
	}
	player := a.context.NewPlayerFromBytes(pcm)
	player.SetVolume(0.6)
	player.Play()
}

// --- Procedural DSP Synthesizers ---

func synthesizeJumpSound() []byte {
	duration := 0.18
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		progress := t / duration
		// Rising frequency sweep 180Hz -> 650Hz
		freq := 180.0 + progress*470.0
		amp := math.Sin(2.0*math.Pi*freq*t) * (1.0 - progress) * 0.4

		sample := int16(amp * 32760)
		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func synthesizeIceCreamSound() []byte {
	// Happy 3-note ascending arpeggio (C5 - E5 - G5)
	notes := []float64{523.25, 659.25, 783.99}
	duration := 0.08
	var fullPCM []byte

	for _, freq := range notes {
		numSamples := int(duration * SampleRate)
		buf := make([]byte, numSamples*BytesPerSample)
		for i := 0; i < numSamples; i++ {
			t := float64(i) / SampleRate
			progress := t / duration
			// Triangle / sine blend
			amp := (0.7*math.Sin(2.0*math.Pi*freq*t) + 0.3*(2.0*math.Abs(2.0*(t*freq-math.Floor(t*freq+0.5)))-1.0)) * (1.0 - progress) * 0.35
			sample := int16(amp * 32760)
			idx := i * 4
			buf[idx] = byte(sample)
			buf[idx+1] = byte(sample >> 8)
			buf[idx+2] = byte(sample)
			buf[idx+3] = byte(sample >> 8)
		}
		fullPCM = append(fullPCM, buf...)
	}
	return fullPCM
}

func synthesizeBottleSound() []byte {
	// High crisp recycling ding (B5 -> E6)
	duration := 0.22
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		progress := t / duration
		freq := 987.77 + math.Sin(progress*math.Pi)*300.0
		amp := math.Sin(2.0*math.Pi*freq*t) * math.Exp(-progress*6.0) * 0.45

		sample := int16(amp * 32760)
		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func synthesizeBellSound() []byte {
	// Double bicycle bell "ring ring"
	duration := 0.25
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		subT := math.Mod(t, 0.12)
		freq := 1760.0 // A6 high bell
		amp := math.Sin(2.0*math.Pi*freq*t) * math.Exp(-subT*25.0) * 0.3

		sample := int16(amp * 32760)
		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func synthesizeBoostSound() []byte {
	duration := 0.3
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		progress := t / duration
		freq := 300.0 + math.Pow(progress, 2.0)*900.0
		// Square wave pulse width boost
		phase := t * freq
		phase = phase - math.Floor(phase)
		var sq float64
		if phase < 0.3 {
			sq = 0.3
		} else {
			sq = -0.3
		}
		amp := sq * (1.0 - progress*0.5) * 0.25

		sample := int16(amp * 32760)
		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func synthesizeCrashSound() []byte {
	duration := 0.35
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		progress := t / duration
		noise := (rand.Float64()*2.0 - 1.0) * (1.0 - progress) * 0.5
		lowFreq := math.Sin(2.0*math.Pi*110.0*t) * (1.0 - progress) * 0.3
		amp := noise + lowFreq

		sample := int16(math.Max(-1.0, math.Min(1.0, amp)) * 32760)
		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}

func synthesizeLondonBGM() []byte {
	// Cheerful 4-bar London city chiptune loop (C major 120 BPM)
	bpm := 120.0
	beatLen := 60.0 / bpm
	totalBeats := 16
	duration := float64(totalBeats) * beatLen
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	melody := []float64{
		261.63, 329.63, 392.00, 523.25, // C4 E4 G4 C5
		349.23, 440.00, 523.25, 698.46, // F4 A4 C5 F5
		392.00, 493.88, 587.33, 783.99, // G4 B4 D5 G5
		523.25, 392.00, 329.63, 261.63, // C5 G4 E4 C4
	}

	bassline := []float64{
		130.81, 130.81, 174.61, 174.61,
		196.00, 196.00, 130.81, 130.81,
	}

	for i := 0; i < numSamples; i++ {
		t := float64(i) / SampleRate
		beatIndex := int(t / beatLen) % len(melody)
		bassIndex := int(t / (beatLen * 2)) % len(bassline)

		// Lead wave
		leadFreq := melody[beatIndex]
		leadAmp := math.Sin(2.0*math.Pi*leadFreq*t) * 0.15

		// Bass wave
		bassFreq := bassline[bassIndex]
		bassT := math.Mod(t, beatLen*2)
		bassAmp := (2.0*(bassT*bassFreq-math.Floor(bassT*bassFreq+0.5))) * 0.12

		amp := leadAmp + bassAmp
		sample := int16(math.Max(-1.0, math.Min(1.0, amp)) * 32760)

		idx := i * 4
		buf[idx] = byte(sample)
		buf[idx+1] = byte(sample >> 8)
		buf[idx+2] = byte(sample)
		buf[idx+3] = byte(sample >> 8)
	}
	return buf
}
