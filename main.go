package main

import (
	"fmt"
	"log"
	"os"
	"parse_audio/parsers"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

func main() {
	startTime := time.Now()

	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <path to audio file>")
	}

	path := os.Args[1]
	audioData, err := loadAudio(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registering parsers")
	parserList := []parsers.Parser{
		&parsers.PlotParser{SamplesPerPoint: 500},
		&parsers.SpectralFluxParser{WindowSize: 1024, HopSize: 512},
	}

	for _, parser := range parserList {
		parserStartTime := time.Now()
		fmt.Println("Parsing with ", parser.Name())
		err := parser.Parse(audioData)
		if err != nil {
			log.Fatal("Failed to parse: ", err)
		}
		fmt.Println("Time taken: ", time.Since(parserStartTime))
		fmt.Println("--------------------------------")
	}

	fmt.Println("Done")
	fmt.Println("Time taken: ", time.Since(startTime))
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
