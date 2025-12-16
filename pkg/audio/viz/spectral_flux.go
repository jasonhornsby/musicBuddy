package viz

import (
	"syscall/js"
	"time"

	"parse_audio/pkg/audio/core"
)

type SpectalFluxVizNode struct {
	core.BaseNode
	outputView js.Value

	renderBuf    []float32
	cachedBufLen int
}

func NewSpectalFluxVizNode(id string, input core.Node) *SpectalFluxVizNode {
	n := &SpectalFluxVizNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
			Input:   input,
		},
	}
	input.AddDownstream(n)
	return n
}

func (n *SpectalFluxVizNode) GetData() (interface{}, error) {
	return nil, nil
}

func (n *SpectalFluxVizNode) BindOutput(bufferJs js.Value) {
	n.outputView = bufferJs
	n.cachedBufLen = bufferJs.Get("length").Int()
}

func (n *SpectalFluxVizNode) Update() {
	startTime := time.Now()
	val, _ := n.Input.GetData()

	data := val.([]float64)

	if len(n.renderBuf) != n.cachedBufLen {
		n.renderBuf = make([]float32, n.cachedBufLen)
	}

	for i := 0; i < n.cachedBufLen; i++ {
		n.renderBuf[i] = float32(data[i])
	}

	js.CopyBytesToJS(n.outputView, core.Float32ToBytes(n.renderBuf))

	n.IsDirty = false
	dur := time.Since(startTime)
	n.SetComputeDurationMs(int(dur.Milliseconds()))

	println("SpectrumVizNode: We have ", len(data), " data points")
}
