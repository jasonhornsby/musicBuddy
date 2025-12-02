//go:build js && wasm

package main

import (
	"syscall/js"
)

func add(this js.Value, args []js.Value) interface{} {
	println("add function called")
	// args holds the arguments passed from JS
	// We should validate length, but for this example we assume correct input
	v1 := args[0].Int()
	v2 := args[1].Int()
	println("value1Str: ", v1)
	println("value2Str: ", v2)

	// Return the result
	return v1 + v2
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("add", js.FuncOf(add))

	println("WASM module initialized")

	<-c
}
