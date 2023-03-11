package wshim

import (
	"fmt"
	"syscall/js"
)

func FloatSlider(name string, min, max, step float64, val *float64) *FloatSliderElement {
	return &FloatSliderElement{
		Name: name,
		Key:  findValidId(name, "FloatSlider"),
		Min:  min,
		Max:  max,
		Step: step,
		Val:  val,
	}
}

type FloatSliderElement struct {
	Name, Key      string
	Min, Max, Step float64
	Val            *float64
}

func (s *FloatSliderElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	slider := document.Call("createElement", "input")
	slider.Call("setAttribute", "type", "range")
	slider.Call("setAttribute", "min", s.Min)
	slider.Call("setAttribute", "max", s.Max)
	slider.Call("setAttribute", "value", *s.Val)
	slider.Call("setAttribute", "step", s.Step)
	slider.Call("setAttribute", "id", s.Key+"s")
	slider.Call("setAttribute", "class", "optionSlider")

	number := document.Call("createElement", "output")
	number.Call("setAttribute", "type", "number")
	number.Call("setAttribute", "value", *s.Val)
	number.Call("setAttribute", "id", s.Key+"n")
	number.Call("setAttribute", "class", "optionNumber")

	update[s.Key] = func(v any) {
		f, _ := v.(float64)
		*s.Val = f
	}

	slider.Call("setAttribute", "oninput", "FloatSliderElementUpdate(this.id, parseFloat(this.value))")

	return s.Name, s.Key, "FloatSliderElement", []js.Value{slider, number}
}

func (s *FloatSliderElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Float()

	update[id[:len(id)-1]](val)

	document := js.Global().Get("parent").Get("document")

	oId := ""

	switch id[0] {
	case 'n':
		oId = id[:len(id)-1] + "s"
	case 's':
		oId = id[:len(id)-1] + "n"
	}

	other := document.Call("getElementById", oId)

	if other.IsNull() {
		fmt.Println("WARNING: Failed to update float slider values:")
		fmt.Println(id, id[:len(id)-1], oId)
		return nil
	}

	other.Set("value", val)

	return nil
}

func IntSlider(name string, min, max, step int, val *int) *IntSliderElement {
	return &IntSliderElement{
		Name: name,
		Key:  findValidId(name, "IntSlider"),
		Min:  min,
		Max:  max,
		Step: step,
		Val:  val,
	}
}

type IntSliderElement struct {
	Name, Key      string
	Min, Max, Step int
	Val            *int
}

func (i *IntSliderElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	slider := document.Call("createElement", "input")
	slider.Call("setAttribute", "type", "range")
	slider.Call("setAttribute", "min", i.Min)
	slider.Call("setAttribute", "max", i.Max)
	slider.Call("setAttribute", "value", *i.Val)
	slider.Call("setAttribute", "step", i.Step)
	slider.Call("setAttribute", "id", i.Key+"s")
	slider.Call("setAttribute", "class", "optionSlider")

	number := document.Call("createElement", "output")
	number.Call("setAttribute", "type", "number")
	number.Call("setAttribute", "value", *i.Val)
	number.Call("setAttribute", "id", i.Key+"n")
	number.Call("setAttribute", "class", "optionNumber")

	update[i.Key] = func(v any) {
		vi, _ := v.(int)
		*i.Val = vi
	}

	slider.Call("setAttribute", "oninput", "IntSliderElementUpdate(this.id, parseInt(this.value))")

	return i.Name, i.Key, "IntSliderElement", []js.Value{slider, number}
}

func (i *IntSliderElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Int()

	update[id[:len(id)-1]](val)

	document := js.Global().Get("parent").Get("document")

	oId := ""

	switch id[0] {
	case 'n':
		oId = id[:len(id)-1] + "s"
	case 's':
		oId = id[:len(id)-1] + "n"
	}

	other := document.Call("getElementById", oId)

	if other.IsNull() {
		fmt.Println("WARNING: Failed to update float slider values:")
		fmt.Println(id, id[:len(id)-1], oId)
		return nil
	}

	other.Set("value", val)

	return nil
}
