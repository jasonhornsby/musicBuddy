package audio

import "syscall/js"

type Manager struct {
	rawMp3  *RawMp3Data
	decoded *DecodedAudioData
	loaded  bool
}

func NewManager() *Manager {
	return &Manager{
		loaded: false,
	}
}

func (m *Manager) Load(msg js.Value) error {
	rawMp3SAB := msg.Get("rawMP3Buffer")
	rawMp3Size := msg.Get("rawMP3Size").Int()

	m.rawMp3 = NewRawMp3Data(rawMp3SAB, rawMp3Size)

	decodedJsBuffers := msg.Get("decodedBuffers")
	numChannels := msg.Get("numChannels").Int()
	numSamples := msg.Get("numSamples").Int()
	sampleRate := msg.Get("sampleRate").Int()
	duration := msg.Get("duration").Float()

	channelSabs := make([]js.Value, numChannels)
	for i := 0; i < numChannels; i++ {
		channelSabs[i] = decodedJsBuffers.Index(i)
	}

	m.decoded = NewDecodedAudioData(channelSabs, numSamples, sampleRate, duration)
	m.loaded = true

	return nil
}

func (m *Manager) IsLoaded() bool {
	return m.loaded
}

func (m *Manager) GetRawMp3() *RawMp3Data {
	return m.rawMp3
}

func (m *Manager) GetDecoded() *DecodedAudioData {
	return m.decoded
}
