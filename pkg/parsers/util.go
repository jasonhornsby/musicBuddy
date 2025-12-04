package parsers

import (
	"fmt"
)

type AudioFormat struct {
	SampleRate  int
	NumChannels int
	Precision   int
}

// AudioData holds the decoded audio samples and format information
type AudioData struct {
	ParsedData [][2]float64
	RawData    []byte
	Format     AudioFormat
}

func (a AudioData) String() string {
	return fmt.Sprintf("AudioData{ParsedData: %v, RawData: %v, Format: %v}", a.ParsedData, a.RawData, a.Format)
}
