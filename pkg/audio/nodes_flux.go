package audio

type FluxNode struct {
	BaseNode
	input Node
	cache []float64
}

func NewFluxNode(id string, input Node) *FluxNode {
	n := &FluxNode{
		BaseNode: BaseNode{
			id:      id,
			isDirty: true,
		},
		input: input,
	}
	input.AddDownstream(n)
	return n
}

func (n *FluxNode) GetData() (interface{}, error) {
	if !n.isDirty {
		return n.cache, nil
	}

	dataIface, err := n.input.GetData()
	if err != nil {
		return nil, err
	}
	// Input is [Frames][Magnitudes]float64
	magnitudeFrames := dataIface.([][]float64)
	numFrames := len(magnitudeFrames)

	flux := make([]float64, numFrames)
	for i := 1; i < numFrames; i++ {
		currFrame := magnitudeFrames[i]
		prevFrame := magnitudeFrames[i-1]

		var sumDiffs float64

		for j := 0; j < len(currFrame); j++ {
			diff := currFrame[j] - prevFrame[j]
			if diff > 0 {
				sumDiffs += diff * diff
			}
		}

		// Skipping sqrt for performance
		flux[i] = sumDiffs
	}

	n.cache = flux
	n.isDirty = false

	return flux, nil
}
