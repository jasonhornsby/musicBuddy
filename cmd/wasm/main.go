package main

import (
	"encoding/json"
	"syscall/js"

	"parse_audio/pkg/audio"
	"parse_audio/pkg/audio/core"
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
	case "configure_viz":
		handleConfigureViz(data)

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
	postMessage("viz_updated", map[string]interface{}{
		"id": id,
	})
}

func handleCreateViz(data js.Value) {
	println("[Go] Creating visualizer")

	id := data.Get("id").String()
	vizType := data.Get("vizType").String()
	buffer := data.Get("buffer")

	// We use default values for config until the config schema is shown to the user

	err := pipelineManager.CreateVisualizer(id, vizType)
	if err != nil {
		postMessage("error", map[string]any{
			"message": err.Error(),
		})
		return
	}

	// Bind the visualizer buffer. This needs to be done after the visualizer is created and configured
	// TODO: Make the Visualizer request a buffer of a specific size
	pipelineManager.BindVizBuffer(id, buffer)

	// Get the configuration schema for the visualizer
	schema, err := pipelineManager.GetVizSchema(id)
	if err != nil {
		postMessage("error", map[string]any{
			"message": err.Error(),
		})
		return
	}
	schemaJSON, _ := json.Marshal(schema)
	treeJSON, _ := json.Marshal(pipelineManager.GetRenderTree())
	postMessage("viz_created", map[string]any{
		"id":     id,
		"schema": string(schemaJSON),
		"tree":   string(treeJSON),
	})

	pipelineManager.UpdateViz(id)
	postMessage("viz_updated", map[string]any{
		"id": id,
	})
	postMessage("viz_ready", map[string]any{
		"id": id,
	})
}

func handleConfigureViz(data js.Value) {
	println("[Go] Configuring visualizer")
	id := data.Get("id").String()
	configJS := data.Get("config").String()

	cfg := map[string]interface{}{}
	err := json.Unmarshal([]byte(configJS), &cfg)
	if err != nil {
		postMessage("error", map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	err = pipelineManager.ConfigureVizNode(id, cfg)
	if err != nil {
		postMessage("error", map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	postMessage("viz_configured", map[string]interface{}{
		"id": id,
	})
}

func parseChannelMode(s string) core.ChannelMode {
	switch s {
	case "left":
		return core.ChannelLeft
	case "right":
		return core.ChannelRight
	default:
		return core.ChannelMix
	}
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
