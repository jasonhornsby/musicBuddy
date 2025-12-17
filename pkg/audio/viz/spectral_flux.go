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
		return nil
	}

	if n.Input.GetId() != fluxNode.GetId() {
		n.Input.RemoveDownstream(n)
		fluxNode.AddDownstream(n)
		n.Input = fluxNode
	}

	return nil
}
