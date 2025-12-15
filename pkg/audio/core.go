package audio

import (
	"encoding/binary"
	"math"
	"syscall/js"
	"unsafe"
)

func Float32ToBytes(floats []float32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&floats[0])), len(floats)*4)
}

func BytesToFloat32Slice(b []byte) []float32 {
	floats := make([]float32, len(b)/4)
	for i := range floats {
		bits := binary.LittleEndian.Uint32(b[i*4:])
		floats[i] = math.Float32frombits(bits)
	}
	return floats
}

type Node interface {
	GetId() string
	Invalidate()
	AddDownstream(n Node)
	RemoveDownstream(n Node)
	GetData() (interface{}, error)
}

type VizNode interface {
	Node
	BindOutput(bufferJs js.Value)
	Update()
}

type BaseNode struct {
	id         string
	isDirty    bool
	downstream []Node
}

func (n *BaseNode) GetId() string {
	return n.id
}

func (n *BaseNode) Invalidate() {
	if n.isDirty {
		return
	}
	n.isDirty = true
	for _, child := range n.downstream {
		child.Invalidate()
	}
}

func (bn *BaseNode) AddDownstream(n Node) {
	bn.downstream = append(bn.downstream, n)
}

func (bn *BaseNode) RemoveDownstream(n Node) {
	for i, child := range bn.downstream {
		if child.GetId() == n.GetId() {
			bn.downstream = append(bn.downstream[:i], bn.downstream[i+1:]...)
			return
		}
	}
}
