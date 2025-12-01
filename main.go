package main

import (
	"fmt"
	"log"
	"os"
	"parse_audio/parsers"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

func main() {

	audioData, err := loadAudio("examples/jazz.mp3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registering parsers")
	parserList := []parsers.Parser{
		&parsers.PlotParser{SamplesPerPoint: 500},
	}

	for _, parser := range parserList {
		fmt.Println("Parsing with ", parser.Name())
		err := parser.Parse(audioData)
		if err != nil {
			log.Fatal("Failed to parse: ", err)
		}
	}

	fmt.Println("Done")
}

func loadAudio(filename string) (*parsers.AudioData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	return &parsers.AudioData{Samples: buffer, Format: format}, nil
}

/*
func plotPoints(streamer beep.Streamer, format beep.Format, samplesPerPoint int) {
	sampleDuration := format.SampleRate.D(samplesPerPoint)

	fmt.Println("Sample duration: ", sampleDuration)

	p := plot.New()
	p.Title.Text = "Audio Plot"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Amplitude"

	sample := make([][2]float64, samplesPerPoint)

	pointsLeft := make(plotter.XYs, int(format.SampleRate)/samplesPerPoint)

	for i := 0; i < int(format.SampleRate)/samplesPerPoint; i++ {
		_, ok := streamer.Stream(sample[:])
		if !ok {
			break
		}
		average := averageSamples(sample)
		pointsLeft[i].X = float64(i) * sampleDuration.Seconds()
		pointsLeft[i].Y = average[0]
	}

	err := plotutil.AddLinePoints(p, "Left", pointsLeft)
	if err != nil {
		log.Fatal("Failed to add line points: ", err)
	}

	if err := p.Save(10*vg.Inch, 10*vg.Inch, "audio_plot.png"); err != nil {
		log.Fatal("Failed to save plot: ", err)
	}
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
*/
