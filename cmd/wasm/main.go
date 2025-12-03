//go:build js && wasm

package main

import (
	"bytes"
	"io"
	"parse_audio/pkg/parsers"
	"syscall/js"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

var audioData *parsers.AudioData

func loadAudio(this js.Value, args []js.Value) interface{} {
	println("Received audio data")
	jsUint8Array := args[0]
	length := jsUint8Array.Length()

	data := make([]byte, length)

	js.CopyBytesToGo(data, jsUint8Array)

	reader := bytes.NewReader(data)
	readCloser := io.NopCloser(reader)

	streamer, format, err := mp3.Decode(readCloser)
	if err != nil {
		println("Error decoding audio: ", err)
		return false
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	audioData = &parsers.AudioData{
		Samples: buffer,
		Format:  format,
	}
	println("Finished loading audio")
	return true
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("loadAudio", js.FuncOf(loadAudio))

	println("WASM module initialized")

	<-c
}
