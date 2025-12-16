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

type WindowingMethod string

const (
	WindowingMethodHann WindowingMethod = "hann"
)

type WindowingConfig struct {
	WindowSize int
	HopSize    int
	Method     WindowingMethod
}

func NewWindowingConfig() *WindowingConfig {
	return &WindowingConfig{
		WindowSize: 1024,
		HopSize:    512,
		Method:     WindowingMethodHann,
	}
}

func (WindowingConfig) nodeConfig() {}
