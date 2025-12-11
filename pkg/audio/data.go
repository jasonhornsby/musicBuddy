package audio

import "syscall/js"

type DecodedAudioData struct {
	channels   [][]float32
	numSamples int
	sampleRate int
}

func NewDecodedAudioData(channelSabs []js.Value, numSamples int, sampleRate int) *DecodedAudioData {
	goChannels := make([][]float32, len(channelSabs))

	for i, sabView := range channelSabs {
		buf := make([]float32, numSamples)
		uInt8Array := js.Global().Get("Uint8Array").New(
			sabView.Get("buffer"),
			sabView.Get("byteOffset"),
			sabView.Get("byteLength"),
		)
		js.CopyBytesToGo(Float32ToBytes(buf), uInt8Array)
		goChannels[i] = buf
	}

	return &DecodedAudioData{
		channels:   goChannels,
		numSamples: numSamples,
		sampleRate: sampleRate,
	}
}

func (d *DecodedAudioData) GetChannel(channelIndex int) []float32 {
	if channelIndex < 0 || channelIndex >= len(d.channels) {
		return nil
	}
	return d.channels[channelIndex]
}

func (d *DecodedAudioData) GetNumSamples() int {
	return d.numSamples
}
