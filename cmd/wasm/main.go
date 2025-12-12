package main

import (
	"parse_audio/pkg/audio"
	"syscall/js"
)

var (
	audioManager    = audio.NewManager()
	pipelineManager = audio.NewPipelineManager(audioManager)
)

func onMessage(this js.Value, args []js.Value) interface{} {
	event := args[0]
	data := event.Get("data")

	if data.Type() == js.TypeUndefined || data.Type() == js.TypeNull {
		println("[Go] Invalid message data")
		return nil
	}

	msgType := data.Get("type").String()
	println("[Go] Message received: ", msgType)

	switch msgType {
	case "load_audio":
		handleLoadAudio(data)
	case "create_viz":
		handleCreateViz(data)
	case "update_viz":
		handleUpdateViz(data)

	default:
		println("[Go] Unknown message type: ", msgType)
		return nil
	}

	return nil
}

func handleLoadAudio(data js.Value) {
	println("[Go] Loading audio")
	err := audioManager.Load(data)
	if err != nil {
		println("[Go] Error loading audio: ", err.Error())
		postMessage("error", map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	postMessage("audio_loaded", map[string]interface{}{
		"numChannels": audioManager.GetNumChannels(),
		"numSamples":  audioManager.GetNumSamples(),
		"sampleRate":  audioManager.GetSampleRate(),
	})
}

func handleUpdateViz(data js.Value) {
	println("[Go] Updating visualizer")
	id := data.Get("id").String()
	pipelineManager.UpdateViz(id)
}

func handleCreateViz(data js.Value) {
	println("[Go] Creating visualizer")
	id := data.Get("id").String()
	vizType := data.Get("vizType").String()
	buffer := data.Get("buffer")

	pipelineManager.CreateVisualizer(id, vizType)
	pipelineManager.BindVizBuffer(id, buffer)
	pipelineManager.UpdateViz(id)
}

func postMessage(msgType string, payload map[string]interface{}) {
	obj := js.Global().Get("Object").New()
	obj.Set("type", msgType)

	for k, v := range payload {
		obj.Set(k, v)
	}

	js.Global().Call("postMessage", obj)
	println("[Go] Message sent: ", msgType)
}

func main() {
	defer println("[Go] Wasm worker stopped")

	println("[Go] Initializing audio manager")
	audioManager = audio.NewManager()
	pipelineManager = audio.NewPipelineManager(audioManager)

	js.Global().Set("onmessage", js.FuncOf(onMessage))

	println("[Go] Wasm worker started")

	select {}
}
