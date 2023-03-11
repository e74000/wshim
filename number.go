package wshim

import "syscall/js"

type FloatSlider struct {
	Name, Key            string
	Min, Max, Init, Step float64
	Val                  *float64
}

func (s *FloatSlider) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	slider := document.Call("createElement", "input")
	slider.Call("setAttribute", "type", "range")
	slider.Call("setAttribute", "min", s.Min)
	slider.Call("setAttribute", "max", s.Max)
	slider.Call("setAttribute", "value", s.Init)
	slider.Call("setAttribute", "step", s.Step)
	slider.Call("setAttribute", "id", "s"+s.Key)
	slider.Call("setAttribute", "class", "optionSlider")

	number := document.Call("createElement", "output")
	number.Call("setAttribute", "type", "number")
	number.Call("setAttribute", "value", s.Init)
	number.Call("setAttribute", "id", "n"+s.Key)
	number.Call("setAttribute", "class", "optionNumber")

	update[s.Key] = func(v any) {
		f, _ := v.(float64)
		*s.Val = f
	}

	slider.Call("setAttribute", "oninput", "FloatSliderUpdate(this.id, parseFloat(this.value))")

	return s.Name, s.Key, "FloatSlider", []js.Value{slider, number}
}

func (s *FloatSlider) Update(this js.Value, params []js.Value) any {
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

type IntSlider struct {
	Name, Key            string
	Min, Max, Init, Step int
	Val                  *int
}

func (i *IntSlider) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	slider := document.Call("createElement", "input")
	slider.Call("setAttribute", "type", "range")
	slider.Call("setAttribute", "min", i.Min)
	slider.Call("setAttribute", "max", i.Max)
	slider.Call("setAttribute", "value", i.Init)
	slider.Call("setAttribute", "step", i.Step)
	slider.Call("setAttribute", "id", "i"+i.Key)
	slider.Call("setAttribute", "class", "optionSlider")

	number := document.Call("createElement", "output")
	number.Call("setAttribute", "type", "number")
	number.Call("setAttribute", "value", i.Init)
	number.Call("setAttribute", "id", "n"+i.Key)
	number.Call("setAttribute", "class", "optionNumber")

	update[i.Key] = func(v any) {
		vi, _ := v.(int)
		*i.Val = vi
	}

	slider.Call("setAttribute", "oninput", "IntSlider(this.id, parseInt(this.value))")

	return i.Name, i.Key, "IntSlider", []js.Value{slider, number}
}

func (i *IntSlider) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Int()

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
