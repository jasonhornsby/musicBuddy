package viz

import (
	"syscall/js"

	"parse_audio/pkg/audio/core"
)

type SpectrumVizNode struct {
	core.BaseNode
	input      core.Node
	outputView js.Value

	renderBuf    []float32
	cachedBufLen int
}

func NewSpectrumVizNode(id string, input core.Node) *SpectrumVizNode {
	n := &SpectrumVizNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
		},
		input: input,
	}
	input.AddDownstream(n)
	return n
}

func (n *SpectrumVizNode) GetData() (interface{}, error) {
	return nil, nil
}

func (n *SpectrumVizNode) BindOutput(bufferJs js.Value) {
	n.outputView = bufferJs
	n.cachedBufLen = bufferJs.Get("length").Int()
}

func (n *SpectrumVizNode) GetInput() core.Node {
	return n.input
}

func (n *SpectrumVizNode) Update() {
	val, _ := n.input.GetData()

	data := val.([]float64)

	if len(n.renderBuf) != n.cachedBufLen {
		n.renderBuf = make([]float32, n.cachedBufLen)
	}

	for i := 0; i < n.cachedBufLen; i++ {
		n.renderBuf[i] = float32(data[i])
	}

	js.CopyBytesToJS(n.outputView, core.Float32ToBytes(n.renderBuf))

	n.IsDirty = false

	println("SpectrumVizNode: We have ", len(data), " data points")
}
