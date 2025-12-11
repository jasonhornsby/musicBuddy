package audio

import (
	"reflect"
	"syscall/js"
	"unsafe"
)

func Float32ToBytes(floats []float32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&floats))
	header.Len *= 4
	header.Cap *= 4
	return *(*[]byte)(unsafe.Pointer(&header))
}

type Node interface {
	GetId() string
	Invalidate()
	AddDownstream(n Node)
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
