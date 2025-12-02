package parsers

import (
	"fmt"

	"github.com/gopxl/beep"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// AudioData holds the decoded audio samples and format information
type AudioData struct {
	Samples *beep.Buffer
	Format  beep.Format
}

func (a AudioData) String() string {
	return fmt.Sprintf("AudioData{Samples: %v, Format: %v}", a.Samples.Len(), a.Format)
}

// Parser is the interface that all audio parsers must implement
type Parser interface {
	Parse(data *AudioData) error
	Name() string
}

// PlotParser generates visual plots from audio data
type PlotParser struct {
	SamplesPerPoint int
}

func (pp *PlotParser) Name() string {
	return "Plot Parser"
}

func (pp *PlotParser) Parse(data *AudioData) error {

	fmt.Println("Plotting with ", pp.SamplesPerPoint, " samples per point")

	sampleDuration := data.Format.SampleRate.D(pp.SamplesPerPoint)
	p := plot.New()
	p.Title.Text = "Audio Plot"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Amplitude"

	streamer := data.Samples.Streamer(0, data.Samples.Len()-1)
	samples := make([][2]float64, pp.SamplesPerPoint)

	pointsLeft := make(plotter.XYs, 0)
	for {
		_, ok := streamer.Stream(samples[:])
		if !ok {
			break
		}
		average := averageSamples(samples)
		pointsLeft = append(
			pointsLeft, plotter.XY{
				X: float64(streamer.Position()) * sampleDuration.Seconds(),
				Y: average[0],
			},
		)
	}

	err := plotutil.AddLinePoints(p, "Left", pointsLeft)
	if err != nil {
		return err
	}

	err = p.Save(10*vg.Inch, 10*vg.Inch, "output/audio_plot.png")
	if err != nil {
		return err
	}

	return nil
}

func averageSamples(samples [][2]float64) [2]float64 {
	sumLeft := 0.0
	sumRight := 0.0
	for _, sample := range samples {
		sumLeft += sample[0]
		sumRight += sample[1]
	}
	return [2]float64{sumLeft / float64(len(samples)), sumRight / float64(len(samples))}
}
