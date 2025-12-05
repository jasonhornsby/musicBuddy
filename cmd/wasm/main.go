package main

import (
	"parse_audio/pkg/audio"
	"syscall/js"
)

var audioManager *audio.Manager

func onMessage(this js.Value, args []js.Value) interface{} {
	msg := args[0]
	msgType := msg.Get("type").String()

	switch msgType {
	case "loadAudio":
		err := audioManager.Load(msg)
		if err != nil {
			println("Error loading audio: ", err)
			return nil
		}

		js.Global().Call("postMessage", js.ValueOf(map[string]interface{}{
			"type":        "audioLoaded",
			"numChannels": audioManager.GetDecoded().NumChannels(),
			"numSamples":  audioManager.GetDecoded().NumSamples(),
			"sampleRate":  audioManager.GetDecoded().SampleRate(),
		}))
	}

	return nil
}

func main() {
	audioManager = audio.NewManager()

	js.Global().Set("onmessage", js.FuncOf(onMessage))

	println("Wasm worker started")

	select {}
}
