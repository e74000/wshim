package wshim

import (
	"fmt"
	"github.com/charmbracelet/log"
	"syscall/js"
)

// Toggle creates a new ToggleElement with the given name and value.
func Toggle(name string, val *bool) *ToggleElement {
	return &ToggleElement{
		Name: name,
		Key:  findValidId(name, "Toggle"),
		Val:  val,
	}
}

// ToggleElement represents a checkbox input element.
type ToggleElement struct {
	Name, Key string
	Val       *bool
	onChange  func(oldVal, newVal bool)
}

// OnChange registers a callback function to be called when the value of the checkbox changes.
func (t *ToggleElement) OnChange(f func(oldVal, newVal bool)) *ToggleElement {
	t.onChange = f
	return t
}

// Build returns the []js.Value for the checkbox input element as well as some identifiers.
func (t *ToggleElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	toggle := document.Call("createElement", "input")
	toggle.Call("setAttribute", "type", "checkbox")
	toggle.Call("setAttribute", "id", t.Key)
	toggle.Call("setAttribute", "class", "optionToggle")

	log.Debug("building toggle element", "name", t.Name, "key", t.Key)
	log.Debug("initial value registered", "initial", *t.Val)

	if *t.Val {
		toggle.Call("setAttribute", "defaultChecked", checked(*t.Val))
		toggle.Call("setAttribute", "checked", checked(*t.Val))
	}

	update[t.Key] = func(v any) {
		if t.onChange != nil {
			t.onChange(*t.Val, v.(bool))
		}
		b, _ := v.(bool)
		*t.Val = b
	}

	toggle.Call("setAttribute", "onclick", "ToggleElementUpdate(this.id, this.checked)")

	return t.Name, t.Key, "ToggleElement", []js.Value{toggle}
}

// Update updates the value of the checkbox input element.
func (t *ToggleElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Bool()

	update[id](val)

	return nil
}

// Radio creates a new RadioElement with the given name and values.
func Radio(name string, vals []string, selected *string) *RadioElement {
	return &RadioElement{
		Name: name,
		Key:  findValidId(name, "Radio"),
		Vals: vals,
		Val:  selected,
	}
}

// RadioElement represents a radio button input element.
type RadioElement struct {
	Name, Key string
	Vals      []string
	Val       *string
	onChange  func(oldVal, newVal string)
}

// OnChange registers a callback function to be called when the value of the radio button changes.
func (r *RadioElement) OnChange(f func(oldVal, newVal string)) *RadioElement {
	r.onChange = f
	return r
}

// Build returns the []js.Value for the radio button input element as well as some identifiers.
func (r *RadioElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	startIndex := -1

	for i, val := range r.Vals {
		if val == *r.Val {
			startIndex = i
			break
		}
	}

	log.Debug("Building radio element with parameters:", "name", r.Name, "key", r.Key, "values", r.Vals)

	if startIndex == -1 {
		log.Debug("invalid initial value", "initial", *r.Val, "fallback", r.Vals[0])

		startIndex = 0
		*r.Val = r.Vals[0]
	} else {
		log.Debug("initial value registered", "initial", *r.Val)
	}

	elems = make([]js.Value, 3*len(r.Vals))

	for i, s := range r.Vals {
		id := fmt.Sprintf("%s-%d", r.Key, i)

		radio := document.Call("createElement", "input")
		radio.Call("setAttribute", "type", "radio")
		radio.Call("setAttribute", "id", id)
		radio.Call("setAttribute", "name", r.Key)
		radio.Call("setAttribute", "value", s)

		radio.Call("setAttribute", "onchange", "RadioElementUpdate(this.name, this.value)")

		if i == startIndex {
			radio.Call("setAttribute", "checked", "checked")
		}

		radioLabel := document.Call("createElement", "label")
		radioLabel.Call("setAttribute", "for", id)
		radioLabel.Call("appendChild", document.Call("createTextNode", s))

		br := document.Call("createElement", "br")

		elems[3*i+0], elems[3*i+1], elems[3*i+2] = radio, radioLabel, br
	}

	update[r.Key] = func(v any) {
		if r.onChange != nil {
			r.onChange(*r.Val, v.(string))
		}

		b, _ := v.(string)
		*r.Val = b
	}

	return r.Name, r.Key, "RadioElement", elems
}

// Update updates the value of the radio button input element.
func (r *RadioElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].String()

	update[id](val)

	return nil
}
