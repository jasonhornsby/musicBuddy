package audio

import "syscall/js"

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

	outputLen := n.outputView.Get("length").Int()

	if len(n.renderBuf) != outputLen {
		n.renderBuf = make([]float32, outputLen)
	}

	step := len(samples) / outputLen

	step = max(step, 1)

	for i := range outputLen {
		srcIdx := int(float64(i) * float64(step))

		if srcIdx < len(samples) {
			n.renderBuf[i] = samples[srcIdx]
		}
		// Decimation
		// TODO: We can do way better than this
		n.renderBuf[i] = samples[srcIdx]
	}

	js.CopyBytesToJS(n.outputView, Float32ToBytes(n.renderBuf))

	n.isDirty = false
}
