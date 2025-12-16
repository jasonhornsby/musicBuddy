package audio

import (
	"fmt"
	"math"
)

type WindowingNode struct {
	BaseNode
	Input Node
	cfg   WindowingConfig
	cache [][]float64
}

func NewWindowingNode(id string, input Node, cfg WindowingConfig) *WindowingNode {
	n := &WindowingNode{
		BaseNode: BaseNode{
			ID:      id,
			IsDirty: true,
		},
		Input: input,
		cfg:   cfg,
	}
	input.AddDownstream(n)
	return n
}

func hannWindow(size int) []float64 {
	window := make([]float64, size)
	for i := 0; i < size; i++ {
		// TODO: Fix all this float64/float32 conversion
		window[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(size-1)))
	}
	return window
}

func (n *WindowingNode) GetData() (interface{}, error) {
	if !n.IsDirty {
		return n.cache, nil
	}

	dataIface, err := n.Input.GetData()
	data := dataIface.([]float32)
	if err != nil {
		return nil, err
	}

	var windowWeights []float64

	switch n.cfg.Method {
	case WindowingMethodHann:
		windowWeights = hannWindow(n.cfg.WindowSize)
	default:
		return nil, fmt.Errorf("invalid windowing method: %s", n.cfg.Method)
	}

	var frames [][]float64

	for i := 0; i < len(data)-n.cfg.WindowSize; i += n.cfg.HopSize {
		rawFrame := data[i : i+n.cfg.WindowSize]
		processedFrame := make([]float64, n.cfg.WindowSize)
		for j := 0; j < n.cfg.WindowSize; j++ {
			processedFrame[j] = float64(rawFrame[j]) * windowWeights[j]
		}
		frames = append(frames, processedFrame)
	}

	n.cache = frames
	n.IsDirty = false

	return frames, nil
}
