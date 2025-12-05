package audio

import "syscall/js"

type Channel struct {
	sab  js.Value
	view js.Value
}

func NewChannel(sab js.Value) *Channel {
	return &Channel{
		sab:  sab,
		view: js.Global().Get("Float32Array").New(sab),
	}
}

func (c *Channel) GetSample(index int) float32 {
	return float32(c.view.Index(index).Float())
}

func (c *Channel) SetSample(index int, value float32) {
	c.view.SetIndex(index, value)
}

func (c *Channel) View() js.Value {
	return c.view
}
