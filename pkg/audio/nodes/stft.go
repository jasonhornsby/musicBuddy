package nodes

import (
	"parse_audio/pkg/audio/core"
	"time"

	"gonum.org/v1/gonum/dsp/fourier"
)

type STFTNode struct {
	core.BaseNode
	cache   [][]complex128
	fftPlan *fourier.FFT
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
	if len(data) == 0 {
		return [][]complex128{}, nil
	}

	// Create or reuse FFT plan (reusing is faster)
	frameSize := len(data[0])
	if n.fftPlan == nil || n.fftPlan.Len() != frameSize {
		n.fftPlan = fourier.NewFFT(frameSize)
	}

	frames := make([][]complex128, len(data))
	for i, frame := range data {
		frames[i] = n.fftPlan.Coefficients(nil, frame)
	}

	n.cache = frames
	n.IsDirty = false
	dur := time.Since(startTime)
	n.SetComputeDurationMs(int(dur.Milliseconds()))

	return frames, nil
}
