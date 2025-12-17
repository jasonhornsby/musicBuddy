package viz

import (
	"syscall/js"
	"time"

	"parse_audio/pkg/audio/core"
	"parse_audio/pkg/audio/nodes"
)

type SpectalFluxVizNode struct {
	core.BaseNode
	outputView js.Value

	renderBuf    []float32
	cachedBufLen int
}

func NewSpectalFluxVizNode(id string) *SpectalFluxVizNode {
	n := &SpectalFluxVizNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
		},
	}
	return n
}

func (n *SpectalFluxVizNode) GetData() (interface{}, error) {
	return nil, nil
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

func (n *SpectalFluxVizNode) GetSchema() []core.ParamDef {
	s := []core.ParamDef{}

	s = append(s, nodes.GetChannelSchema()...)
	s = append(s, nodes.GetWindowingSchema()...)

	return s
}

func (n *SpectalFluxVizNode) Reconfigure(cfg core.ConfigMap, provider core.NodeProvider) error {
	channelConfig := nodes.ParseChannelConfig(cfg)
	winSize := nodes.ParseWindowingConfig(cfg)

	fluxNode := provider.GetFluxNode(channelConfig, winSize)

	if n.Input == nil {
		n.Input = fluxNode
		fluxNode.AddDownstream(n)
	} else if n.Input.GetId() != fluxNode.GetId() {
		n.Input.RemoveDownstream(n)
		fluxNode.AddDownstream(n)
		n.Input = fluxNode
	}

	// Make sure the output buffer is the correct size
	numFloats := winSize / 2
	sizeBytes := numFloats * 4

	n.ensureBufferSize(sizeBytes)

	return nil
}

func (n *SpectalFluxVizNode) ensureBufferSize(size int) {
	println("[GO] Checking buffer size: ", size, n.ID)
	if n.outputView.Truthy() && n.outputView.Get("length").Int() == size {
		return
	}

	println("[GO] Requesting new buffer size: ", size, n.ID)

	buffer := js.Global().Call("allocateVizBuffer", js.ValueOf(n.ID), js.ValueOf(size))

	if !buffer.Truthy() {
		panic("[GO] Failed to allocate buffer")
	}

	n.outputView = buffer
	n.cachedBufLen = size
}
