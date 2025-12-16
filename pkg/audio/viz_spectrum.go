package audio

import "syscall/js"

type SpectrumVizNode struct {
	BaseNode
	input      Node
	outputView js.Value

	renderBuf    []float32
	cachedBufLen int
}

func NewSpectrumVizNode(id string, input Node) *SpectrumVizNode {
	n := &SpectrumVizNode{
		BaseNode: BaseNode{
			id:      id,
			isDirty: true,
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

func (n *SpectrumVizNode) GetInput() Node {
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

	js.CopyBytesToJS(n.outputView, Float32ToBytes(n.renderBuf))

	n.isDirty = false

	println("SpectrumVizNode: We have ", len(data), " data points")
}
