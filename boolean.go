package wshim

import (
	"fmt"
	"syscall/js"
)

func Toggle(name string, val *bool) *ToggleElement {
	return &ToggleElement{
		Name: name,
		Key:  findValidId(name, "Toggle"),
		Val:  val,
	}
}

type ToggleElement struct {
	Name, Key string
	Val       *bool
	onChange  func(oldVal, newVal bool)
}

func (t *ToggleElement) OnChange(f func(oldVal, newVal bool)) *ToggleElement {
	t.onChange = f
	return t
}

func (t *ToggleElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	toggle := document.Call("createElement", "input")
	toggle.Call("setAttribute", "type", "checkbox")
	toggle.Call("setAttribute", "id", t.Key)
	toggle.Call("setAttribute", "class", "optionToggle")

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

func (t *ToggleElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].Bool()

	update[id](val)

	return nil
}

func Radio(name string, vals []string, selected *string) *RadioElement {
	return &RadioElement{
		Name: name,
		Key:  findValidId(name, "Radio"),
		Vals: vals,
		Val:  selected,
	}
}

func (r *RadioElement) OnChange(f func(oldVal, newVal string)) *RadioElement {
	r.onChange = f
	return r
}

type RadioElement struct {
	Name, Key string
	Vals      []string
	Val       *string
	onChange  func(oldVal, newVal string)
}

func (r *RadioElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	startIndex := -1

	for i, val := range r.Vals {
		if val == *r.Val {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		startIndex = 0
		*r.Val = r.Vals[0]
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

func (r *RadioElement) Update(this js.Value, params []js.Value) any {
	id := params[0].String()
	val := params[1].String()

	update[id](val)

	return nil
}
