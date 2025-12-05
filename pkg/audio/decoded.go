package audio

import "syscall/js"

type DecodedAudioData struct {
	channels   []*Channel
	sampleRate int
	numSamples int
	duration   float64
}

func NewDecodedAudioData(channelSABs []js.Value, numSamples int, sampleRate int, duration float64) *DecodedAudioData {
	channels := make([]*Channel, len(channelSABs))

	for i, sab := range channelSABs {
		channels[i] = NewChannel(sab)
	}

	return &DecodedAudioData{
		channels:   channels,
		sampleRate: sampleRate,
		numSamples: numSamples,
		duration:   duration,
	}
}

func (d *DecodedAudioData) NumChannels() int {
	return len(d.channels)
}

func (d *DecodedAudioData) SampleRate() int {
	return d.sampleRate
}

func (d *DecodedAudioData) NumSamples() int {
	return d.numSamples
}

func (d *DecodedAudioData) Duration() float64 {
	return d.duration
}

func (d *DecodedAudioData) GetSample(channelIndex int, sampleIndex int) float32 {
	if channelIndex < 0 || channelIndex >= len(d.channels) {
		return 0
	}
	if sampleIndex < 0 || sampleIndex >= d.numSamples {
		return 0
	}
	return d.channels[channelIndex].GetSample(sampleIndex)
}

func (d *DecodedAudioData) SetSample(channelIndex int, sampleIndex int, value float32) {
	if channelIndex < 0 || channelIndex >= len(d.channels) {
		return
	}
	if sampleIndex < 0 || sampleIndex >= d.numSamples {
		return
	}
	d.channels[channelIndex].SetSample(sampleIndex, value)
}

func (d *DecodedAudioData) GetFrame(sampleIndex int) []float32 {
	if sampleIndex < 0 || sampleIndex >= d.numSamples {
		return nil
	}
	frame := make([]float32, len(d.channels))
	for i, channel := range d.channels {
		frame[i] = channel.GetSample(sampleIndex)
	}
	return frame
}

func (d *DecodedAudioData) GetChannelView(channelIndex int) js.Value {
	if channelIndex < 0 || channelIndex >= len(d.channels) {
		return js.Undefined()
	}
	return d.channels[channelIndex].View()
}

func (d *DecodedAudioData) GetChannel(channelIndex int) *Channel {
	if channelIndex < 0 || channelIndex >= len(d.channels) {
		return nil
	}
	return d.channels[channelIndex]
}

func (d *DecodedAudioData) ProcessWindow(channelIndex, start, length int, fn func(sample float32, index int) float32) {
	if channelIndex < 0 || channelIndex >= len(d.channels) {
		return
	}
	if start < 0 || start+length > d.numSamples {
		return
	}

	channel := d.channels[channelIndex]
	for i := 0; i < length; i++ {
		idx := start + i
		sample := channel.GetSample(idx)
		result := fn(sample, i)
		channel.SetSample(idx, result)
	}
}
