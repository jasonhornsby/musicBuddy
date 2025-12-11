package audio

import (
	"fmt"
	"syscall/js"
)

type PipelineManager struct {
	source      *Manager
	nodes       map[string]VizNode
	sharedNodes map[string]Node
}

func NewPipelineManager(source *Manager) *PipelineManager {
	return &PipelineManager{
		source:      source,
		nodes:       make(map[string]VizNode),
		sharedNodes: make(map[string]Node),
	}
}

func (pm *PipelineManager) GetChannelNode(mode ChannelMode) Node {
	key := fmt.Sprintf("channel_%s", mode)
	if node, ok := pm.sharedNodes[key]; ok {
		return node
	}

	node := NewChannelNode(key, pm.source, mode)
	pm.sharedNodes[key] = node
	return node
}

func (pm *PipelineManager) CreateVisualizer(id string, vizType string) {
	audioSrc := pm.GetChannelNode(ChannelMix)

	if vizType == "waveform" {
		node := NewWaveVizNode(id, audioSrc)
		pm.nodes[id] = node
	}
}

func (pm *PipelineManager) BindVizBuffer(id string, buffer js.Value) {
	if node, ok := pm.nodes[id]; ok {
		node.BindOutput(buffer)
	}
}

func (pm *PipelineManager) UpdateViz(id string) {
	if node, ok := pm.nodes[id]; ok {
		node.Update()
	} else {
		println("[Go] Visualizer not found: ", id)
	}
	pm.PrintPipeline()
}

func (pm *PipelineManager) PrintPipeline() {
	println("[Go] Pipeline:")
	println("[Go] - ", len(pm.nodes), " nodes")
	for id, node := range pm.nodes {
		println("[Go] - ", id, node.GetId())
	}
	println("[Go] Shared nodes:")
	for id, node := range pm.sharedNodes {
		println("[Go] - ", id, node.GetId())
	}
}
