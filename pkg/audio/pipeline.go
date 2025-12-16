package audio

import (
	"fmt"
	"syscall/js"

	"parse_audio/pkg/audio/core"
	"parse_audio/pkg/audio/nodes"
	"parse_audio/pkg/audio/viz"
)

type PipelineManager struct {
	source      *Manager
	vizNodes    map[string]core.VizNode
	sharedNodes map[string]core.Node
}

func NewPipelineManager(source *Manager) *PipelineManager {
	return &PipelineManager{
		source:      source,
		vizNodes:    make(map[string]core.VizNode),
		sharedNodes: make(map[string]core.Node),
	}
}

func (pm *PipelineManager) getOrBuildNode(key string, builder func() core.Node) core.Node {
	if node, ok := pm.sharedNodes[key]; ok {
		return node
	}

	node := builder()
	pm.sharedNodes[key] = node
	return node
}

type VizCfg struct {
	Channel    core.ChannelMode
	WindowSize int
}

// Dependency builders

func (pm *PipelineManager) GetChannelNode(mode core.ChannelMode) core.Node {
	key := fmt.Sprintf("channel_%d", mode)
	if node, ok := pm.sharedNodes[key]; ok {
		return node
	}

	return pm.getOrBuildNode(key, func() core.Node {
		return nodes.NewChannelNode(key, pm.source, mode)
	})
}

func (pm *PipelineManager) GetWindowingNode(mode core.ChannelMode, size int) core.Node {
	key := fmt.Sprintf("windowing_%d_%d", mode, size)

	return pm.getOrBuildNode(key, func() core.Node {
		input := pm.GetChannelNode(mode)
		return nodes.NewWindowingNode(key, input, *core.NewWindowingConfig())
	})
}

func (pm *PipelineManager) GetSTFTNode(mode core.ChannelMode, size int) core.Node {
	key := fmt.Sprintf("stft_%d_%d", mode, size)

	return pm.getOrBuildNode(key, func() core.Node {
		input := pm.GetWindowingNode(mode, size)
		return nodes.NewSTFTNode(key, input)
	})
}

func (pm *PipelineManager) GetMagnitudeNode(mode core.ChannelMode, size int) core.Node {
	key := fmt.Sprintf("magnitude_%d_%d", mode, size)

	return pm.getOrBuildNode(key, func() core.Node {
		input := pm.GetSTFTNode(mode, size)
		return nodes.NewMagnitudeNode(key, input)
	})
}

func (pm *PipelineManager) GetFluxNode(mode core.ChannelMode, size int) core.Node {
	key := fmt.Sprintf("flux_%d_%d", mode, size)

	return pm.getOrBuildNode(key, func() core.Node {
		input := pm.GetMagnitudeNode(mode, size)
		return nodes.NewFluxNode(key, input)
	})
}

func (pm *PipelineManager) CreateVisualizer(id string, vizType string, cfg VizCfg) {
	var vizNode core.VizNode

	switch vizType {
	case "waveform":
		audioSrc := pm.GetChannelNode(cfg.Channel)
		vizNode = viz.NewWaveVizNode(id, audioSrc)
	case "spectral_flux":
		input := pm.GetFluxNode(cfg.Channel, cfg.WindowSize)
		vizNode = viz.NewSpectalFluxVizNode(id, input)
	default:
		panic("invalid visualizer type: " + vizType)
	}
	pm.vizNodes[id] = vizNode
}

// TODO: The cfg should not be a WaveformConfig, but a NodeConfig
// This way we can use the same function to configure any visualizer node
func (pm *PipelineManager) ConfigureVizNode(id string, cfg WaveformConfig) error {
	node, ok := pm.vizNodes[id]
	if !ok {
		return fmt.Errorf("visualizer not found: %s", id)
	}

	waveNode, ok := node.(*viz.WaveVizNode)
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
	if node, ok := pm.vizNodes[id]; ok {
		node.BindOutput(buffer)
	}
}

func (pm *PipelineManager) UpdateViz(id string) {
	if node, ok := pm.vizNodes[id]; ok {
		node.Update()
	} else {
		println("[Go] Visualizer not found: ", id)
	}
	pm.PrintPipeline()
}

func (pm *PipelineManager) PrintPipeline() {
	println("[Go] Pipeline:")
	println("[Go] - ", len(pm.vizNodes), " viz nodes")
	for id, node := range pm.vizNodes {
		println("[Go] - ", id, node.GetId())
	}
	println("[Go] Shared nodes:")
	for id, node := range pm.sharedNodes {
		println("[Go] - ", id, node.GetId())
	}
}

func (pm *PipelineManager) GetRenderTree() core.RenderTree {
	var outputEdges []core.GraphEdge
	var outputNodes []core.GraphNode

	addNode := func(n Node, category string) {
		nodeType := fmt.Sprintf("%T", n)

		outputNodes = append(outputNodes, core.GraphNode{
			ID:       n.GetId(),
			Type:     nodeType,
			Category: category,
			Label:    n.GetId(),
		})
	}

	resolveInput := func(n Node) {
		var input Node

		switch n.(type) {
		case *nodes.WindowingNode:
			input = n.(*nodes.WindowingNode).Input
		case *nodes.STFTNode:
			input = n.(*nodes.STFTNode).Input
		case *nodes.MagnitudeNode:
			input = n.(*nodes.MagnitudeNode).Input
		case *nodes.FluxNode:
			input = n.(*nodes.FluxNode).Input
		case *nodes.ChannelNode:
			return
		case *viz.WaveVizNode:
			input = n.(*viz.WaveVizNode).Input
		case *viz.SpectalFluxVizNode:
			input = n.(*viz.SpectalFluxVizNode).Input
		default:
			panic("unknown node type: " + fmt.Sprintf("%T", n))
		}

		if input != nil {
			outputEdges = append(outputEdges, core.GraphEdge{
				From: input.GetId(),
				To:   n.GetId(),
			})
		}
	}

	for _, node := range pm.sharedNodes {
		addNode(node, "compute")
		resolveInput(node)
	}

	for _, node := range pm.vizNodes {
		addNode(node, "visualizer")
		resolveInput(node)
	}

	return core.RenderTree{
		Nodes: outputNodes,
		Edges: outputEdges,
	}
}
