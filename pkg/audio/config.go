package audio

// NodeConfig marker interface
type NodeConfig interface {
	nodeConfig()
}

// WaveformConfig for waveform visualizer
type WaveformConfig struct {
	Channel ChannelMode
}

func (WaveformConfig) nodeConfig() {}
