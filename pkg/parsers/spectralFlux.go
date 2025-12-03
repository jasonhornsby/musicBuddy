package parsers

import (
	"errors"
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

type SpectralFluxParser struct {
	WindowSize int
	HopSize    int
}

func (sfp *SpectralFluxParser) Name() string {
	return "Spectral Flux Parser"
}

func (sfp *SpectralFluxParser) Parse(data *AudioData) ([]float64, error) {
	samples := make([]float64, data.Samples.Len())
	streamer := data.Samples.Streamer(0, data.Samples.Len())

	// Temporary buffer to store the samples
	tempBuffer := make([][2]float64, sfp.WindowSize)
	sampleIdx := 0

	for {
		n, ok := streamer.Stream(tempBuffer[:])
		if !ok || n == 0 {
			break
		}
		for i := 0; i < n; i++ {
			// Mix down to single channel, (Left + Right) / 2
			samples[sampleIdx] = (tempBuffer[i][0] + tempBuffer[i][1]) / 2.0
			sampleIdx++
		}
	}

	win := window.Hann(sfp.WindowSize)

	var previousSpectrum []float64
	var spectralFlux []float64

	for i := 0; i < len(samples)-sfp.WindowSize; i += sfp.HopSize {
		frame := make([]float64, sfp.WindowSize)
		copy(frame, samples[i:i+sfp.WindowSize])

		for j := 0; j < sfp.WindowSize; j++ {
			frame[j] *= win[j]
		}

		fftComplex := fft.FFTReal(frame)

		currSpectrum := make([]float64, len(fftComplex))
		for j, c := range fftComplex {
			currSpectrum[j] = cmplx.Abs(c)
		}
		flux := 0.0
		if previousSpectrum != nil {
			for j := 0; j < len(currSpectrum); j++ {
				diff := currSpectrum[j] - previousSpectrum[j]
				if diff > 0 {
					flux += diff
				}
			}
			flux = math.Sqrt(flux)
		}
		spectralFlux = append(spectralFlux, flux)
		previousSpectrum = currSpectrum
	}

	return smooth(spectralFlux, 50), nil
}

func smooth(input []float64, window int) []float64 {
	output := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		sum := 0.0
		count := 0
		for j := i - window/2; j <= i+window/2; j++ {
			if j >= 0 && j < len(input) {
				sum += input[j]
				count++
			}
		}
		output[i] = sum / float64(count)
	}
	return output
}

func EstimateBPM(spectralFlux []float64, sampleRate int, hopSize int) (float64, error) {
	fluxFPS := float64(sampleRate) / float64(hopSize)
	minBPM := 60.0
	maxBPM := 180.0

	maxLag := int(fluxFPS * 60.0 / minBPM)
	minLag := int(fluxFPS * 60.0 / maxBPM)

	bestLag := 0
	maxCorrelation := -1.0

	// Average flux to remove DC offset
	totalFlux := 0.0
	for _, f := range spectralFlux {
		totalFlux += f
	}
	avgFlux := totalFlux / float64(len(spectralFlux))

	for lag := minLag; lag <= maxLag; lag++ {
		correlation := 0.0

		// Correlate the signal with itself shifted by 'lag'
		for i := 0; i < len(spectralFlux)-lag; i++ {
			val1 := spectralFlux[i] - avgFlux
			val2 := spectralFlux[i+lag] - avgFlux
			correlation += val1 * val2
		}

		// Check if this is the strongest recurring pattern
		if correlation > maxCorrelation {
			maxCorrelation = correlation
			bestLag = lag
		}
	}
	if bestLag == 0 {
		return 0.0, errors.New("No valid BPM found. Try adjusting the min/max BPM range.")
	}

	estimatedBPM := 60.0 / (float64(bestLag) / fluxFPS)

	return estimatedBPM, nil
}
