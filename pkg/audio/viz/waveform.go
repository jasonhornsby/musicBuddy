package viz

import (
	"fmt"
	"syscall/js"
	"time"

	"parse_audio/pkg/audio/core"
)

type WaveVizNode struct {
	core.BaseNode
	Input      core.Node
	outputView js.Value

	renderBuf    []float32
	cachedBufLen int
}

func NewWaveVizNode(id string, input core.Node) *WaveVizNode {
	n := &WaveVizNode{
		BaseNode: core.BaseNode{
			ID:      id,
			IsDirty: true,
		},
		Input: input,
	}
	input.AddDownstream(n)
	return n
}

func (n *WaveVizNode) GetData() (interface{}, error) {
	return nil, nil
}

func (n *WaveVizNode) BindOutput(bufferJs js.Value) {
	n.outputView = bufferJs
	n.cachedBufLen = bufferJs.Get("length").Int()
}

func (n *WaveVizNode) GetInput() core.Node {
	return n.Input
}

func (n *WaveVizNode) SetInput(input core.Node) {
	n.Input = input
	n.IsDirty = true
}

func (n *WaveVizNode) Update() {
	startTime := time.Now()

	val, _ := n.Input.GetData()
	samples := val.([]float32)
	totalSamples := len(samples)

	// Use cached length
	outputCount := n.cachedBufLen / 2
	outputLen := n.cachedBufLen

	if len(n.renderBuf) != outputLen {
		n.renderBuf = make([]float32, outputLen)
	}

	if outputCount == 0 {
		return
	}

	// Calculate how many raw samples fit into one visual pixel column
	step := totalSamples / outputCount

	// We will never check more than this many samples per pixel column.
	// Lower limit = faster but might miss interesting peak
	const maxSamplesPerPixel = 50

	// Calculate the Stride (Skip Rate)
	stride := 1
	if step > maxSamplesPerPixel {
		stride = step / maxSamplesPerPixel
	}

	readIdx := 0
	writeIdx := 0

	var globalMin float32 = 2.0
	var globalMax float32 = -2.0

	for range outputCount {
		end := min(readIdx+step, totalSamples)
		chunk := samples[readIdx:end]

		var min float32 = 2.0
		var max float32 = -2.0

		foundAny := false

		for k := 0; k < len(chunk); k += stride {
			val := chunk[k]
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
			foundAny = true
		}

		// Fallback for silence/empty chunks
		if !foundAny {
			min = 0
			max = 0
		}

		n.renderBuf[writeIdx] = min
		n.renderBuf[writeIdx+1] = max
		writeIdx += 2

		if min < globalMin {
			globalMin = min
		}
		if max > globalMax {
			globalMax = max
		}

		readIdx = end
	}

	n.IsDirty = false

	// Fast cast to bytes
	js.CopyBytesToJS(n.outputView, core.Float32ToBytes(n.renderBuf))

	// Performance logging
	dur := time.Since(startTime)
	actualScanned := outputCount * (step / stride)
	println("[GO] Waveform Update took:", dur.Milliseconds(), "ms. Samples scanned:", actualScanned)
	fmt.Printf("[GO] Output buffer range - smallest: %f, largest: %f\n", globalMin, globalMax)
}
