package audio

import "syscall/js"

type RawMp3Data struct {
	sab  js.Value
	view js.Value
	size int
}

func NewRawMp3Data(sab js.Value, size int) *RawMp3Data {
	return &RawMp3Data{
		sab:  sab,
		view: js.Global().Get("Uint8Array").New(sab),
		size: size,
	}
}

func (r *RawMp3Data) GetByte(index int) byte {
	if index < 0 || index >= r.size {
		return 0
	}
	return byte(r.view.Index(index).Int())
}

func (r *RawMp3Data) Size() int {
	return r.size
}
