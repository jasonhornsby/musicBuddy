package nodes

import (
	"parse_audio/pkg/audio/core"
	"time"

	"github.com/mjibson/go-dsp/fft"
)

type STFTNode struct {
	core.BaseNode
	cache [][]complex128
}

func NewSTFTNode(id string, input core.Node) *STFTNode {
	n := &STFTNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
			Input:   input,
		},
	}
	input.AddDownstream(n)
	return n
}

func (n *STFTNode) GetData() (interface{}, error) {

	if !n.IsDirty {
		return n.cache, nil
	}

	dataIface, err := n.Input.GetData()
	startTime := time.Now()

	if err != nil {
		return nil, err
	}
	data := dataIface.([][]float64)
	frames := make([][]complex128, len(data))
	for i, frame := range data {
		fftComplex := fft.FFTReal(frame)
		frames[i] = fftComplex
	}

	n.cache = frames
	n.IsDirty = false
	dur := time.Since(startTime)
	n.SetComputeDurationMs(int(dur.Milliseconds()))

	return frames, nil
}
