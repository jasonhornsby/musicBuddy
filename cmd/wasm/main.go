//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"io"
	"parse_audio/pkg/parsers"
	"syscall/js"
	"time"
	"unsafe"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

var isInitialized = false
var audioData *parsers.AudioData

func loadParsedAudio(this js.Value, args []js.Value) interface{} {
	start := time.Now()
	println("Received parsed audio from JS")
	payload := args[0]
	parsedAudioJs := payload.Get("parsedAudio")
	rawDataJs := payload.Get("rawData")

	rawDataLength := rawDataJs.Length()
	rawData := make([]byte, rawDataLength)
	js.CopyBytesToGo(rawData, rawDataJs)
	reader := bytes.NewReader(rawData)
	readCloser := io.NopCloser(reader)
	_, fmt, _ := mp3.Decode(readCloser)

	// Get the underlying ArrayBuffer from the Float64Array and copy as bytes
	parsedAudioBuffer := parsedAudioJs.Get("buffer")
	parsedAudioUint8 := js.Global().Get("Uint8Array").New(parsedAudioBuffer)
	parsedAudioByteLength := parsedAudioUint8.Length()
	parsedAudioBytes := make([]byte, parsedAudioByteLength)
	js.CopyBytesToGo(parsedAudioBytes, parsedAudioUint8)

	// Reinterpret bytes as [][2]float64 without copying
	numStereoSamples := parsedAudioByteLength / 2
	parsedAudioData := unsafe.Slice((*[2]float64)(unsafe.Pointer(&parsedAudioBytes[0])), numStereoSamples)

	format := parsers.AudioFormat{
		SampleRate:  int(fmt.SampleRate),
		NumChannels: fmt.NumChannels,
		Precision:   fmt.Precision,
	}

	audioData = &parsers.AudioData{
		ParsedData: parsedAudioData,
		RawData:    rawData,
		Format:     format,
	}

	isInitialized = true
	println("Time taken to copy bytes to go: ", formatSeconds(time.Since(start).Seconds()))

	return js.ValueOf(true)
}

func loadAudio(this js.Value, args []js.Value) interface{} {
	println("Received audio data")
	start := time.Now()
	jsUint8Array := args[0]
	length := jsUint8Array.Length()

	data := make([]byte, length)

	js.CopyBytesToGo(data, jsUint8Array)
	println("Time taken to copy bytes to go: ", formatSeconds(time.Since(start).Seconds()))

	reader := bytes.NewReader(data)
	readCloser := io.NopCloser(reader)

	println("Time taken to create reader: ", formatSeconds(time.Since(start).Seconds()))

	streamer, format, err := mp3.Decode(readCloser)
	if err != nil {
		println("Error decoding audio: ", err)
		return false
	}

	println("Time taken to decode audio: ", formatSeconds(time.Since(start).Seconds()))
	streamerLength := streamer.Len()
	println("Streamer length: ", streamerLength)
	buffer := beep.NewBuffer(format)
	if streamerLength > 0 {
		println("Taking ", streamerLength, " samples from streamer")
		buffer.Append(beep.Take(streamerLength, streamer))
	} else {
		println("Appending full streamer")
		buffer.Append(streamer)
	}

	println("Time taken to create buffer: ", formatSeconds(time.Since(start).Seconds()))

	// Convert beep.Buffer to [][2]float64
	numSamples := buffer.Len()
	sampleStreamer := buffer.Streamer(0, numSamples)
	parsedAudioData := make([][2]float64, numSamples)
	n, ok := sampleStreamer.Stream(parsedAudioData)
	if !ok || n != numSamples {
		println("Error streaming samples, got ", n, " expected ", numSamples)
		return js.ValueOf(false)
	}

	myFormat := parsers.AudioFormat{
		SampleRate:  int(format.SampleRate),
		NumChannels: format.NumChannels,
		Precision:   2,
	}
	audioData = &parsers.AudioData{
		ParsedData: parsedAudioData,
		RawData:    data,
		Format:     myFormat,
	}

	isInitialized = true
	println("Time taken to load audio: ", formatSeconds(time.Since(start).Seconds()))

	// Bulk copy samples to JS using unsafe slice conversion
	byteLength := numSamples * 16 // 2 channels * 8 bytes per float64
	bytesSlice := unsafe.Slice((*byte)(unsafe.Pointer(&parsedAudioData[0])), byteLength)

	// Create Uint8Array and bulk copy all bytes at once
	samplesUint8Array := js.Global().Get("Uint8Array").New(byteLength)
	js.CopyBytesToJS(samplesUint8Array, bytesSlice)

	// Create Float64Array view over the same buffer (no additional copy)
	jsFloat64Array := js.Global().Get("Float64Array").New(samplesUint8Array.Get("buffer"))

	println("Time taken to create samples array: ", formatSeconds(time.Since(start).Seconds()))
	return jsFloat64Array
}

// formatSeconds prints seconds with up to 3 decimal places and no scientific notation.
func formatSeconds(seconds float64) string {
	return fmt.Sprintf("%.3fs", seconds)
}

func unloadAudio(this js.Value, args []js.Value) interface{} {
	if !isInitialized {
		return js.ValueOf(false)
	}

	audioData = nil
	isInitialized = false
	println("Unloaded audio", audioData, isInitialized)
	return js.ValueOf(true)
}

func getAudioMetadata(this js.Value, args []js.Value) interface{} {
	println("Getting audio metadata", isInitialized)
	if !isInitialized {
		return js.ValueOf(false)
	}

	metadata, err := parsers.GetAudioMetadata(audioData)

	if err != nil {
		println("Error getting audio metadata: ", err)
		return js.ValueOf(false)
	}

	return map[string]interface{}{
		"sampleRate":     metadata.SampleRate,
		"channels":       metadata.Channels,
		"durationMs":     metadata.DurationMs,
		"decodedBitrate": metadata.DecodedBitrate,
		"metadata": map[string]interface{}{
			"name":   metadata.Metadata.Name,
			"artist": metadata.Metadata.Artist,
			"album":  metadata.Metadata.Album,
			"year":   metadata.Metadata.Year,
			"format": metadata.Metadata.Format,
		},
	}
}

func getSpectralFlux(this js.Value, args []js.Value) interface{} {

	println("Getting spectral flux", isInitialized)
	if !isInitialized {
		return js.ValueOf(false)
	}

	parser := &parsers.SpectralFluxParser{
		WindowSize: 1024,
		HopSize:    512,
	}

	spectralFlux, err := parser.Parse(audioData)
	println("Spectral flux: ", spectralFlux)
	if err != nil {
		println("Error getting spectral flux: ", err)
		return js.ValueOf(false)
	}

	println("Something new")

	// Bulk copy to JS using unsafe slice conversion
	byteLength := len(spectralFlux) * 8 // 8 bytes per float64
	bytesSlice := unsafe.Slice((*byte)(unsafe.Pointer(&spectralFlux[0])), byteLength)

	// Create Uint8Array and bulk copy all bytes at once
	uint8Array := js.Global().Get("Uint8Array").New(byteLength)
	js.CopyBytesToJS(uint8Array, bytesSlice)

	// Create Float64Array view over the same buffer (no additional copy)
	return js.Global().Get("Float64Array").New(uint8Array.Get("buffer"))
}

func main() {
	c := make(chan struct{}, 0)

	// Prep functions
	js.Global().Set("loadAudio", js.FuncOf(loadAudio))
	js.Global().Set("loadParsedAudio", js.FuncOf(loadParsedAudio))
	js.Global().Set("unloadAudio", js.FuncOf(unloadAudio))

	// Audio metadata functions
	js.Global().Set("getAudioMetadata", js.FuncOf(getAudioMetadata))
	// Spectral flux functions
	js.Global().Set("getSpectralFlux", js.FuncOf(getSpectralFlux))

	<-c
}
