package parsers

type AudioMetadata struct {
	SampleRate     int
	Channels       int
	DurationMs     int
	DecodedBitrate int
}

func GetAudioMetadata(audioData *AudioData) *AudioMetadata {

	bitrate := int(audioData.Format.SampleRate) * audioData.Format.NumChannels * audioData.Format.Precision * 8

	return &AudioMetadata{
		SampleRate:     int(audioData.Format.SampleRate),
		Channels:       audioData.Format.NumChannels,
		DurationMs:     (int(audioData.Samples.Len()/int(audioData.Format.SampleRate)) * 1000),
		DecodedBitrate: bitrate,
	}
}
