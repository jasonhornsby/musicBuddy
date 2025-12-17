package nodes

import (
	"math"
	"time"

	"parse_audio/pkg/audio/core"
)

type WindowingNode struct {
	core.BaseNode
	winSize int
	cache   [][]float64
}

func NewWindowingNode(id string, input core.Node, winSize int) *WindowingNode {
	n := &WindowingNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
			Input:   input,
		},
		winSize: winSize,
	}
	input.AddDownstream(n)
	return n
}

func GetWindowingSchema() []core.ParamDef {
	return []core.ParamDef{
		{
			Key:     "win_size",
			Label:   "Window size",
			Type:    core.ParamSelect,
			Default: 1024,
			Options: []string{"128", "256", "512", "1024", "2048", "4096"},
		},
	}
}

func ParseWindowingConfig(cfg core.ConfigMap) int {
	winSize := cfg.GetInt("win_size", 1024)
	println("[Go] Windowing config: ", winSize)
	return winSize
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
	startTime := time.Now()

	var windowWeights []float64

	windowWeights = hannWindow(n.winSize)

	var frames [][]float64

	for i := 0; i < len(data)-n.winSize; i += n.winSize / 2 {
		rawFrame := data[i : i+n.winSize]
		processedFrame := make([]float64, n.winSize)
		for j := 0; j < n.winSize; j++ {
			processedFrame[j] = float64(rawFrame[j]) * windowWeights[j]
		}
		frames = append(frames, processedFrame)
	}

	n.cache = frames
	n.IsDirty = false

	dur := time.Since(startTime)
	n.SetComputeDurationMs(int(dur.Milliseconds()))

	return frames, nil
}
