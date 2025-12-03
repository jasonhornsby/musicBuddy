//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"io"
	"parse_audio/pkg/parsers"
	"syscall/js"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

var isInitialized = false
var rawData []byte
var audioData *parsers.AudioData

func loadAudio(this js.Value, args []js.Value) interface{} {
	println("Received audio data")
	start := time.Now()
	jsUint8Array := args[0]
	length := jsUint8Array.Length()

	rawData = make([]byte, length)

	js.CopyBytesToGo(rawData, jsUint8Array)
	println("Time taken to copy bytes to go: ", formatSeconds(time.Since(start).Seconds()))

	reader := bytes.NewReader(rawData)
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

	audioData = &parsers.AudioData{
		Samples: buffer,
		Format:  format,
	}

	isInitialized = true
	println("Time taken to load audio: ", formatSeconds(time.Since(start).Seconds()))
	return true
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

	metadata := parsers.GetAudioMetadata(audioData)

	println("Audio metadata: ", metadata)

	println("Returning audio metadata")

	return map[string]interface{}{
		"sampleRate":     metadata.SampleRate,
		"channels":       metadata.Channels,
		"durationMs":     metadata.DurationMs,
		"decodedBitrate": metadata.DecodedBitrate,
	}
}

func main() {
	c := make(chan struct{}, 0)

	// Prep functions
	js.Global().Set("loadAudio", js.FuncOf(loadAudio))
	js.Global().Set("unloadAudio", js.FuncOf(unloadAudio))

	// Audio metadata functions
	js.Global().Set("getAudioMetadata", js.FuncOf(getAudioMetadata))

	<-c
}
