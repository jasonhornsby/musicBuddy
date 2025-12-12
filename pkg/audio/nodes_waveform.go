package audio

import (
	"math"
	"syscall/js"
)

type WaveVizNode struct {
	BaseNode
	input      Node
	outputView js.Value

	renderBuf []float32
}

func NewWaveVizNode(id string, input Node) *WaveVizNode {
	n := &WaveVizNode{
		BaseNode: BaseNode{
			id:      id,
			isDirty: true,
		},
		input: input,
	}
	input.AddDownstream(n)
	return n
}

func (n *WaveVizNode) GetData() (interface{}, error) {
	return nil, nil
}

func (n *WaveVizNode) BindOutput(bufferJs js.Value) {
	n.outputView = bufferJs
}

func (n *WaveVizNode) Update() {
	// Get data from channel selector node
	val, _ := n.input.GetData()
	samples := val.([]float32)

	// We have space for two points per pixel
	outputCount := n.outputView.Get("length").Int() / 2

	// Amount of points to render * 2 for min and max. Data is interleaved
	outputLen := n.outputView.Get("length").Int()

	if len(n.renderBuf) != outputLen {
		n.renderBuf = make([]float32, outputLen)
	}

	step := len(samples) / outputCount
	for i := 0; i < outputCount; i++ {
		min := math.Inf(1)
		max := math.Inf(-1)
		for j := 0; j < step; j++ {
			idx := i*step + j
			if idx < len(samples) {
				min = math.Min(float64(min), float64(samples[idx]))
				max = math.Max(float64(max), float64(samples[idx]))
			}
		}
		n.renderBuf[i*2] = float32(min)
		n.renderBuf[i*2+1] = float32(max)
	}

	js.CopyBytesToJS(n.outputView, Float32ToBytes(n.renderBuf))

	n.isDirty = false
}
