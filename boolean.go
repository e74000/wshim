package wshim

import (
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
}

func (t *ToggleElement) Build() (label, key, sType string, elems []js.Value) {
	document := js.Global().Get("parent").Get("document")

	toggle := document.Call("createElement", "input")
	toggle.Call("setAttribute", "type", "checkbox")
	toggle.Call("setAttribute", "checked", checked(*t.Val))
	toggle.Call("setAttribute", "id", t.Key)
	toggle.Call("setAttribute", "class", "optionToggle")

	update[t.Key] = func(v any) {
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
