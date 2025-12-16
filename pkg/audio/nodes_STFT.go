package audio

import "github.com/mjibson/go-dsp/fft"

type STFTNode struct {
	BaseNode
	input Node
	cache [][]complex128
}

func NewSTFTNode(id string, input Node) *STFTNode {
	n := &STFTNode{
		BaseNode: BaseNode{
			id:      id,
			isDirty: true,
		},
		input: input,
	}
	input.AddDownstream(n)
	return n
}

func (n *STFTNode) GetData() (interface{}, error) {
	if !n.isDirty {
		return n.cache, nil
	}

	dataIface, err := n.input.GetData()
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
	n.isDirty = false

	return frames, nil
}
