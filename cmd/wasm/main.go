//go:build js && wasm

package main

import (
	"strconv"
	"syscall/js"
)

func add(this js.Value, args []js.Value) interface{} {
	// args holds the arguments passed from JS
	// We should validate length, but for this example we assume correct input
	value1Str := args[0].String()
	value2Str := args[1].String()

	// Convert strings to integers
	v1, _ := strconv.Atoi(value1Str)
	v2, _ := strconv.Atoi(value2Str)

	// Return the result
	return v1 + v2
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("add", js.FuncOf(add))

	println("WASM module initialized")

	<-c
}
