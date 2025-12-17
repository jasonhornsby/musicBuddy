package nodes

import (
	"math/cmplx"
	"time"

	"parse_audio/pkg/audio/core"
)

type MagnitudeNode struct {
	core.BaseNode
	cache [][]float64
}

func NewMagnitudeNode(id string, input core.Node) *MagnitudeNode {
	n := &MagnitudeNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
			Input:   input,
		},
	}
	input.AddDownstream(n)
	return n
}

func (n *MagnitudeNode) GetData() (interface{}, error) {
	if !n.IsDirty {
		return n.cache, nil
	}

	dataIface, err := n.Input.GetData()
	if err != nil {
		return nil, err
	}
	startTime := time.Now()
	data := dataIface.([][]complex128)

	output := make([][]float64, len(data))
	for i, cFrame := range data {
		halfSize := (len(cFrame) / 2) + 1
		mags := make([]float64, halfSize)

		for j := range halfSize {
			mags[j] = cmplx.Abs(cFrame[j])
		}
		output[i] = mags
	}

	n.cache = output
	n.IsDirty = false

	dur := time.Since(startTime)
	n.SetComputeDurationMs(int(dur.Milliseconds()))

	return output, nil
}
