package audio

import "math/cmplx"

type MagnitudeNode struct {
	BaseNode
	input Node
	cache [][]float64
}

func NewMagnitudeNode(id string, input Node) *MagnitudeNode {
	n := &MagnitudeNode{
		BaseNode: BaseNode{
			id:      id,
			isDirty: true,
		},
		input: input,
	}
	input.AddDownstream(n)
	return n
}

func (n *MagnitudeNode) GetData() (interface{}, error) {
	if !n.isDirty {
		return n.cache, nil
	}

	dataIface, err := n.input.GetData()
	if err != nil {
		return nil, err
	}
	data := dataIface.([][]complex128)

	output := make([][]float64, len(data))
	for i, cFrame := range data {
		halfSize := (len(cFrame) / 2) + 1
		mags := make([]float64, halfSize)

		for j := 0; j < halfSize; j++ {
			mags[j] = cmplx.Abs(cFrame[j])
		}
		output[i] = mags
	}

	n.cache = output
	n.isDirty = false

	return output, nil
}
