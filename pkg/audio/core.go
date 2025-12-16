package audio

// Re-export core types for backward compatibility
import "parse_audio/pkg/audio/core"

type Node = core.Node
type VizNode = core.VizNode
type BaseNode = core.BaseNode
type ChannelMode = core.ChannelMode
type WindowingConfig = core.WindowingConfig
type WindowingMethod = core.WindowingMethod
type DecodedAudioData = core.DecodedAudioData

const (
	ChannelLeft  = core.ChannelLeft
	ChannelRight = core.ChannelRight
	ChannelMix   = core.ChannelMix
)

const WindowingMethodHann = core.WindowingMethodHann

var Float32ToBytes = core.Float32ToBytes
var BytesToFloat32Slice = core.BytesToFloat32Slice
var NewWindowingConfig = core.NewWindowingConfig
var NewDecodedAudioData = core.NewDecodedAudioData
