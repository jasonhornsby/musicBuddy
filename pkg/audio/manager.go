package audio

import (
	"fmt"
	"syscall/js"

	"parse_audio/pkg/audio/core"
)

type Manager struct {
	core.BaseNode

	rawMp3  *RawMp3Data
	decoded *core.DecodedAudioData
	loaded  bool
}

func NewManager() *Manager {
	return &Manager{
		BaseNode: core.BaseNode{
			ID: "source",
		},
		loaded: false,
	}
}

func (m *Manager) GetId() string {
	return "source"
}

func (m *Manager) GetData() (interface{}, error) {
	if !m.loaded {
		return nil, fmt.Errorf("Audio not laoded")
	}
	return m.decoded, nil
}

func (m *Manager) Load(msg js.Value) error {
	rawMp3SAB := msg.Get("rawMP3Buffer")
	rawMp3Size := msg.Get("rawMP3Size").Int()

	m.rawMp3 = NewRawMp3Data(rawMp3SAB, rawMp3Size)

	decodedJsBuffers := msg.Get("decodedBuffers")
	numChannels := msg.Get("numChannels").Int()
	numSamples := msg.Get("numSamples").Int()
	sampleRate := msg.Get("sampleRate").Int()

	channelSabs := make([]js.Value, numChannels)
	for i := 0; i < numChannels; i++ {
		channelSabs[i] = decodedJsBuffers.Index(i)
	}

	m.decoded = core.NewDecodedAudioData(channelSabs, numSamples, sampleRate)
	m.loaded = true

	m.Invalidate()

	return nil
}

func (m *Manager) IsLoaded() bool {
	return m.loaded
}

func (m *Manager) GetRawMp3() *RawMp3Data {
	return m.rawMp3
}

func (m *Manager) GetDecoded() *core.DecodedAudioData {
	return m.decoded
}

func (m *Manager) GetNumChannels() int {
	return len(m.decoded.Channels)
}

func (m *Manager) GetNumSamples() int {
	return m.decoded.NumSamples
}

func (m *Manager) GetSampleRate() int {
	return m.decoded.SampleRate
}
