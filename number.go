package wshim

import (
	"fmt"
	"github.com/charmbracelet/log"
	"syscall/js"
)

// FloatSlider creates a new FloatSliderElement
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

// FloatSliderElement is a slider element that stores a float value
type FloatSliderElement struct {
	Name, Key      string
	Min, Max, Step float64
	Val            *float64
}

// Build returns the []js.Value for the float slider input element as well as some identifiers.
func (s *FloatSliderElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	slider := document.Call("createElement", "input")
	slider.Call("setAttribute", "type", "range")
	slider.Call("setAttribute", "min", s.Min)
	slider.Call("setAttribute", "max", s.Max)

	log.Debug("building float slider element with parameters:", "name", s.Name, "key", s.Key, "min", s.Min, "max", s.Max, "step", s.Step)

	if *s.Val > s.Max || *s.Val < s.Min {
		clamped := clamp(*s.Val, s.Min, s.Max)

		log.Debug("initial value is out of range, clamping:", "initial", *s.Val, "clamped", clamped)

		*s.Val = clamped
	} else {
		log.Debug("initial value registered:", "initial", *s.Val)
	}

	slider.Call("setAttribute", "value", *s.Val)
	slider.Call("setAttribute", "step", s.Step)
	slider.Call("setAttribute", "id", s.Key)
	slider.Call("setAttribute", "class", "optionSlider")

	number := document.Call("createElement", "output")
	number.Call("setAttribute", "type", "number")
	number.Call("setAttribute", "value", *s.Val)
	number.Call("setAttribute", "id", s.Key+"-num")
	number.Call("setAttribute", "class", "optionNumber")

	number.Call("appendChild", document.Call("createTextNode", fmt.Sprintf("%.2f", *s.Val)))

	update[s.Key] = func(v any) {
		f, _ := v.(float64)
		*s.Val = f
	}

	slider.Call("setAttribute", "oninput", "FloatSliderElementUpdate(this.id, parseFloat(this.value))")

	return s.Name, s.Key, "FloatSliderElement", []js.Value{slider, number}
}

// Update updates the value of the float slider input element.
func (s *FloatSliderElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Float()

	update[id](val)

	document := js.Global().Get("parent").Get("document")

	oId := id + "-num"
	other := document.Call("getElementById", oId)
	other.Set("value", val)

	return nil
}

// IntSlider creates a new IntSliderElement
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

// IntSliderElement is a slider element that stores an integer value
type IntSliderElement struct {
	Name, Key      string
	Min, Max, Step int
	Val            *int
	onChange       func(oldVal, newVal int)
}

// OnChange registers a callback function that is called when the value of the int slider changes.
func (s *IntSliderElement) OnChange(f func(oldVal, newVal int)) *IntSliderElement {
	s.onChange = f
	return s
}

// Build returns the []js.Value for the int slider input element as well as some identifiers.
func (s *IntSliderElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	slider := document.Call("createElement", "input")
	slider.Call("setAttribute", "type", "range")
	slider.Call("setAttribute", "min", s.Min)
	slider.Call("setAttribute", "max", s.Max)

	log.Debug("building integer slider element with parameters:", "name", s.Name, "key", "min", s.Min, "max", s.Max, "step", s.Step)

	if *s.Val > s.Max || *s.Val < s.Min {
		clamped := clamp(*s.Val, s.Min, s.Max)

		log.Debug("initial value is out of range, clamping:", "initial", *s.Val, "clamped", clamped)

		*s.Val = clamp(*s.Val, s.Min, s.Max)
	} else {
		log.Debug("initial value registered:", "initial", *s.Val)
	}

	slider.Call("setAttribute", "value", *s.Val)
	slider.Call("setAttribute", "step", s.Step)
	slider.Call("setAttribute", "id", s.Key)
	slider.Call("setAttribute", "class", "optionSlider")

	number := document.Call("createElement", "output")
	number.Call("setAttribute", "type", "number")
	number.Call("setAttribute", "value", *s.Val)
	number.Call("setAttribute", "id", s.Key+"-num")
	number.Call("setAttribute", "class", "optionNumber")

	number.Call("appendChild", document.Call("createTextNode", fmt.Sprintf("%d", *s.Val)))

	update[s.Key] = func(v any) {
		if s.onChange != nil {
			s.onChange(*s.Val, v.(int))
		}

		vi, _ := v.(int)
		*s.Val = vi
	}

	slider.Call("setAttribute", "oninput", "IntSliderElementUpdate(this.id, parseInt(this.value))")

	return s.Name, s.Key, "IntSliderElement", []js.Value{slider, number}
}

// Update updates the value of the integer slider input element.
func (s *IntSliderElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Int()

	update[id](val)

	document := js.Global().Get("parent").Get("document")

	oId := id + "-num"
	other := document.Call("getElementById", oId)
	other.Set("value", val)

	return nil
}
