package parsers

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type SpectralFluxParser struct {
	WindowSize int
	HopSize    int
}

func (sfp *SpectralFluxParser) Name() string {
	return "Spectral Flux Parser"
}

func (sfp *SpectralFluxParser) Parse(data *AudioData) error {

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
			samples[i] = (tempBuffer[i][0] + tempBuffer[i][1]) / 2.0
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
	fmt.Printf("Calculated %d flux points\n", len(spectralFlux))

	p := plot.New()
	p.Title.Text = "Spectral Flux"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Flux"

	points := make(plotter.XYs, len(spectralFlux))
	for i, flux := range spectralFlux {
		points[i] = plotter.XY{X: float64(i), Y: flux}
	}

	err := plotutil.AddLinePoints(p, "Spectral Flux", points)
	if err != nil {
		return err
	}

	err = p.Save(10*vg.Inch, 10*vg.Inch, "output/spectral_flux.png")
	if err != nil {
		return err
	}

	fmt.Println("Spectral Flux Parser")
	return nil
}
