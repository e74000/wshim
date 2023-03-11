package wshim

import (
	"fmt"
	"syscall/js"
)

type InputElement interface {
	Build() (label, key, sType string, elems []js.Value)
	Update(this js.Value, params []js.Value) any
}

var update map[string]func(any)

var ids = map[string]bool{}

func Run(mainFunc func(), elements ...InputElement) {
	fmt.Println("Running webshim")

	seen := make(map[string]bool)
	update = make(map[string]func(any))

	document := js.Global().Get("parent").Get("document")
	optionPanel := document.Call("getElementById", "options")

	for _, element := range elements {
		l, k, t, elems := element.Build()
		if !seen[t] {
			seen[t] = true
			js.Global().Get("parent").Set(t+"Update", js.FuncOf(element.Update))
		}

		label := document.Call("createElement", "label")
		label.Call("setAttribute", "for", k)
		label.Call("setAttribute", "class", "optionLabel")
		label.Call("appendChild", document.Call("createTextNode", l))

		option := document.Call("createElement", "div")
		option.Call("setAttribute", "class", "option")

		inputBox := document.Call("createElement", "div")
		inputBox.Call("setAttribute", "class", "inputBox")

		for _, elem := range elems {
			inputBox.Call("appendChild", elem)
		}

		option.Call("appendChild", label)
		option.Call("appendChild", inputBox)
		optionPanel.Call("appendChild", option)
	}

	mainFunc()
	<-make(chan bool)
}
