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
	key := fmt.Sprintf("channel_%d", mode)
	if node, ok := pm.sharedNodes[key]; ok {
		return node
	}

	node := NewChannelNode(key, pm.source, mode)
	pm.sharedNodes[key] = node
	return node
}

func (pm *PipelineManager) CreateVisualizer(id string, vizType string, cfg NodeConfig) {
	if vizType == "waveform" {
		waveCfg, ok := cfg.(WaveformConfig)
		if !ok {
			panic("invalid waveform configuration")
		}
		audioSrc := pm.GetChannelNode(waveCfg.Channel)
		node := NewWaveVizNode(id, audioSrc)
		pm.nodes[id] = node
	}
}

// TODO: The cfg should not be a WaveformConfig, but a NodeConfig
// This way we can use the same function to configure any visualizer node
func (pm *PipelineManager) ConfigureVizNode(id string, cfg WaveformConfig) error {
	node, ok := pm.nodes[id]
	if !ok {
		return fmt.Errorf("visualizer not found: %s", id)
	}

	waveNode, ok := node.(*WaveVizNode)
	if !ok {
		return fmt.Errorf("node is not a WaveVizNode: %s", id)
	}

	// Get the new channel node
	newInput := pm.GetChannelNode(cfg.Channel)
	oldInput := waveNode.GetInput()

	// Skip if same input
	if oldInput.GetId() == newInput.GetId() {
		return nil
	}

	// Rewire: disconnect from old, connect to new
	oldInput.RemoveDownstream(waveNode)
	newInput.AddDownstream(waveNode)
	waveNode.SetInput(newInput)

	return nil
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
