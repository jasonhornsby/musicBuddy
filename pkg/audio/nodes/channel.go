package nodes

import (
	"fmt"
	"time"

	"parse_audio/pkg/audio/core"
)

type ChannelNode struct {
	core.BaseNode
	mode  core.ChannelMode
	cache []float32
}

func NewChannelNode(id string, input core.Node, mode core.ChannelMode) *ChannelNode {
	n := &ChannelNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
			Input:   input,
		},
		mode: mode,
	}
	input.AddDownstream(n)
	return n
}

func (n *ChannelNode) GetId() string {
	return n.ID
}

func (n *ChannelNode) GetData() (interface{}, error) {
	if !n.IsDirty {
		return n.cache, nil
	}

	dataIface, err := n.Input.GetData()
	startTime := time.Now()
	if err != nil {
		return nil, err
	}

	decoded := dataIface.(*core.DecodedAudioData)

	if n.mode == core.ChannelMix {
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

	n.IsDirty = false
	dur := time.Since(startTime)
	n.SetComputeDurationMs(int(dur.Milliseconds()))
	return n.cache, nil
}

func GetChannelSchema() []core.ParamDef {
	return []core.ParamDef{
		{
			Key:     "channel",
			Label:   "Channel",
			Type:    core.ParamSelect,
			Default: core.ChannelMix,
			Options: []string{"mix", "left", "right"},
		},
	}
}

func ParseChannelConfig(cfg core.ConfigMap) core.ChannelMode {
	channel := cfg.GetString("channel", "mix")
	switch channel {
	case "mix":
		return core.ChannelMix
	case "left":
		return core.ChannelLeft
	case "right":
		return core.ChannelRight
	}
	return core.ChannelMix
}
