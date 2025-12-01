package parsers

import (
	"fmt"
)

type SpectralFluxParser struct {
}

func (sfp *SpectralFluxParser) Name() string {
	return "Spectral Flux Parser"
}

func (sfp *SpectralFluxParser) Parse(data *AudioData) error {
	fmt.Println("Spectral Flux Parser")
	return nil
}
