//go:build js && wasm

package main

import (
	"bytes"
	"io"
	"syscall/js"

	"github.com/gopxl/beep/mp3"
)

func loadAudio(this js.Value, args []js.Value) interface{} {
	println("Received audio data")
	jsUint8Array := args[0]
	length := jsUint8Array.Length()

	data := make([]byte, length)

	js.CopyBytesToGo(data, jsUint8Array)

	reader := bytes.NewReader(data)
	readCloser := io.NopCloser(reader)

	_, format, err := mp3.Decode(readCloser)
	if err != nil {
		println("Error decoding audio: ", err)
		return false
	}

	println("Format: ", format.SampleRate, format.NumChannels, format.Precision)

	return true
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("loadAudio", js.FuncOf(loadAudio))

	println("WASM module initialized")

	<-c
}
