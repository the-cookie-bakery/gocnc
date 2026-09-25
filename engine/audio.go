package engine

import (
	"bytes"
	"io"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type AudioManager struct {
	context *audio.Context
	samples map[string][]byte
	mu      sync.Mutex
	players []*audio.Player
	volume  float64
	muted   bool
}

func NewAudioManager(fsys fs.FS, sampleRate int) *AudioManager {
	am := &AudioManager{
		context: audio.NewContext(sampleRate),
		samples: make(map[string][]byte),
		volume:  1.0,
	}
	return am
}

func (am *AudioManager) Load(name string, path string, fsys fs.FS) error {
	f, err := fsys.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		return err
	}

	am.mu.Lock()
	am.samples[name] = buf.Bytes()
	am.mu.Unlock()
	return nil
}

func (am *AudioManager) Play(name string) {
	if am.muted {
		return
	}

	am.mu.Lock()
	data, ok := am.samples[name]
	am.mu.Unlock()
	if !ok {
		return
	}

	p, err := am.context.NewPlayer(bytes.NewReader(data))
	if err != nil {
		return
	}

	p.SetVolume(am.volume)
	p.Play()

	am.players = append(am.players, p)

	active := make([]*audio.Player, 0)
	for _, pl := range am.players {
		if pl.IsPlaying() {
			active = append(active, pl)
		} else {
			pl.Close()
		}
	}
	am.players = active
}

func (am *AudioManager) SetVolume(v float64) {
	am.volume = v
}

func (am *AudioManager) ToggleMute() {
	am.muted = !am.muted
}

func (am *AudioManager) Close() {
	for _, p := range am.players {
		p.Close()
	}
	am.players = nil
}
