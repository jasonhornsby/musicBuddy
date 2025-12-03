package parsers

type AudioMetadata struct {
	SampleRate int
	Channels   int
	DurationMs int
}

func GetAudioMetadata(audioData *AudioData) *AudioMetadata {
	return &AudioMetadata{
		SampleRate: int(audioData.Format.SampleRate),
		Channels:   audioData.Format.NumChannels,
		DurationMs: (int(audioData.Samples.Len()/int(audioData.Format.SampleRate)) * 1000),
	}
}
