package audio

import "syscall/js"

type DecodedAudioData struct {
	channels   [][]float32
	numSamples int
	sampleRate int
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
