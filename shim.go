package wshim

import (
	"syscall/js"
)

var update map[string]func(any)

func Run(main func(), elements []InputElement) {
	seen := make(map[string]bool)

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

	main()
	<-make(chan bool)
}

type InputElement interface {
	Build() (label, key, sType string, elems []js.Value)
	Update(this js.Value, params []js.Value) any
}

type Slider struct {
	Name, Key            string
	Min, Max, Init, Step float64
	Val                  *float64
}

func (s *Slider) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	slider := document.Call("createElement", "input")
	slider.Call("setAttribute", "type", "range")
	slider.Call("setAttribute", "min", s.Min)
	slider.Call("setAttribute", "max", s.Max)
	slider.Call("setAttribute", "value", s.Init)
	slider.Call("setAttribute", "step", s.Step)
	slider.Call("setAttribute", "id", "s"+s.Key)
	slider.Call("setAttribute", "class", "optionSlider")

	number := document.Call("createElement", "input")
	number.Call("setAttribute", "type", "number")
	number.Call("setAttribute", "value", s.Init)
	number.Call("setAttribute", "id", "n"+s.Key)
	number.Call("setAttribute", "class", "optionNumber")

	update[s.Key] = func(v any) {
		f, _ := v.(float64)
		*s.Val = f
	}

	slider.Call("setAttribute", "oninput", "SliderUpdate(this.id, parseFloat(this.value))")
	number.Call("setAttribute", "onchange", "SliderUpdate(this.id, parseFloat(this.value))")

	return s.Name, s.Key, "Slider", []js.Value{number, slider}
}

func (s *Slider) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Float()

	update[id[1:]](val)

	document := js.Global().Get("parent").Get("document")

	oId := ""

	switch id[0] {
	case 'n':
		oId = "s" + id[1:]
	case 's':
		oId = "n" + id[1:]
	}

	other := document.Call("getElementById", oId)
	other.Set("value", val)

	return nil
}
