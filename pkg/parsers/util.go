package parsers

import (
	"fmt"

	"github.com/gopxl/beep"
)

// AudioData holds the decoded audio samples and format information
type AudioData struct {
	Samples *beep.Buffer
	RawData []byte
	Format  beep.Format
}

func (a AudioData) String() string {
	return fmt.Sprintf("AudioData{Samples: %v, Format: %v}", a.Samples.Len(), a.Format)
}
