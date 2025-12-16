package core

import (
	"encoding/binary"
	"math"
	"syscall/js"
	"unsafe"
)

func Float32ToBytes(floats []float32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&floats[0])), len(floats)*4)
}

func BytesToFloat32Slice(b []byte) []float32 {
	floats := make([]float32, len(b)/4)
	for i := range floats {
		bits := binary.LittleEndian.Uint32(b[i*4:])
		floats[i] = math.Float32frombits(bits)
	}
	return floats
}

type Node interface {
	GetId() string
	Invalidate()
	AddDownstream(n Node)
	RemoveDownstream(n Node)
	GetData() (interface{}, error)
}

type VizNode interface {
	Node
	BindOutput(bufferJs js.Value)
	Update()
}

type BaseNode struct {
	ID         string
	IsDirty    bool
	Downstream []Node
}

func (n *BaseNode) GetId() string {
	return n.ID
}

func (n *BaseNode) Invalidate() {
	if n.IsDirty {
		return
	}
	n.IsDirty = true
	for _, child := range n.Downstream {
		child.Invalidate()
	}
}

func (bn *BaseNode) AddDownstream(n Node) {
	bn.Downstream = append(bn.Downstream, n)
}

func (bn *BaseNode) RemoveDownstream(n Node) {
	for i, child := range bn.Downstream {
		if child.GetId() == n.GetId() {
			bn.Downstream = append(bn.Downstream[:i], bn.Downstream[i+1:]...)
			return
		}
	}
}

// ChannelMode represents how audio channels should be processed
type ChannelMode int

const (
	ChannelLeft  ChannelMode = 0
	ChannelRight ChannelMode = 1
	// ChannelMix is an average of all channels
	ChannelMix ChannelMode = 99
)

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

// DecodedAudioData represents decoded audio samples
type DecodedAudioData struct {
	Channels   [][]float32
	NumSamples int
	SampleRate int
}

func NewDecodedAudioData(channelSabs []js.Value, numSamples int, sampleRate int) *DecodedAudioData {
	println("[Go] New decoded audio data: ", len(channelSabs), " channels, ", numSamples, " samples, ", sampleRate, " sample rate")
	goChannels := make([][]float32, len(channelSabs))

	for i, sabView := range channelSabs {
		bytesBuf := make([]byte, numSamples*4)
		uInt8Array := js.Global().Get("Uint8Array").New(
			sabView.Get("buffer"),
			sabView.Get("byteOffset"),
			sabView.Get("byteLength"),
		)
		js.CopyBytesToGo(bytesBuf, uInt8Array)
		buf := BytesToFloat32Slice(bytesBuf)
		println("[Go] Channel ", i, " has ", len(buf), " samples")
		println("[Go] Channel ", i, " first 10 samples: ", buf[:10])
		// Check the range to ensure normalization between -1 and 1
		var minVal, maxVal float32 = buf[0], buf[0]
		for _, v := range buf {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
		println("[Go] Channel", i, "buffer min:", minVal, "max:", maxVal)
		if minVal < -1.0 || maxVal > 1.0 {
			println("[Go][Warning] Channel", i, " buffer not normalized between -1 and 1!")
		}
		goChannels[i] = buf
	}
	return &DecodedAudioData{
		Channels:   goChannels,
		NumSamples: numSamples,
		SampleRate: sampleRate,
	}
}

func (d *DecodedAudioData) GetChannel(channelIndex int) []float32 {
	if channelIndex < 0 || channelIndex >= len(d.Channels) {
		return nil
	}
	return d.Channels[channelIndex]
}

func (d *DecodedAudioData) GetNumSamples() int {
	return d.NumSamples
}
