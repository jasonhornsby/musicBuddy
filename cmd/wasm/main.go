package main

import (
	"parse_audio/pkg/audio"
	"syscall/js"
)

var audioManager *audio.Manager

func onMessage(this js.Value, args []js.Value) interface{} {
	msg := args[0]
	msgData := msg.Get("data")
	msgType := msgData.Get("type").String()

	println("[Go] Message received: ", msgType)

	switch msgType {
	case "load_audio":
		println("[Go] Loading audio")
		err := audioManager.Load(msgData)
		if err != nil {
			println("Error loading audio: ", err)
			return nil
		}

		js.Global().Call("postMessage", js.ValueOf(map[string]interface{}{
			"type":        "audio_loaded",
			"numChannels": audioManager.GetDecoded().NumChannels(),
			"numSamples":  audioManager.GetDecoded().NumSamples(),
			"sampleRate":  audioManager.GetDecoded().SampleRate(),
		}))
		println("[Go] Audio loaded")
	default:
		println("[Go] Unknown message type: ", msgType)
		return nil
	}

	return nil
}

func main() {
	defer println("[Go] Wasm worker stopped")

	println("[Go] Initializing audio manager")
	audioManager = audio.NewManager()

	js.Global().Set("onmessage", js.FuncOf(onMessage))

	println("[Go] Wasm worker started")

	select {}
}
