package audio

import "fmt"

type ChannelMode int

const (
	ChannelLeft  ChannelMode = 0
	ChannelRight ChannelMode = 1
	// Mix is an average of all channels
	ChannelMix ChannelMode = 99
)

type ChannelNode struct {
	BaseNode
	input Node
	mode  ChannelMode
	cache []float32
}

func NewChannelNode(id string, input Node, mode ChannelMode) *ChannelNode {
	n := &ChannelNode{
		BaseNode: BaseNode{
			id:      id,
			isDirty: true,
		},
		input: input,
		mode:  mode,
	}
	input.AddDownstream(n)
	return n
}

func (n *ChannelNode) GetId() string {
	return n.id
}

func (n *ChannelNode) GetData() (interface{}, error) {
	if !n.isDirty {
		return n.cache, nil
	}

	data, err := n.input.GetData()
	if err != nil {
		return nil, err
	}

	decoded := data.(*DecodedAudioData)

	if n.mode == ChannelMix {
		left := decoded.GetChannel(0)
		right := decoded.GetChannel(1)

		if len(left) != len(right) {
			return nil, fmt.Errorf("left and right channels have different lengths")
		}

		if len(n.cache) != len(left) {
			n.cache = make([]float32, len(left))
		}

		for i := 0; i < len(left); i++ {
			n.cache[i] = (left[i] + right[i]) / 2.0
		}
	} else {
		n.cache = decoded.GetChannel(int(n.mode))
	}

	n.isDirty = false

	return n.cache, nil
}
