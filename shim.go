package wshim

import (
	"github.com/charmbracelet/log"
	"syscall/js"
	"time"
)

// InputElement is an interface that represents an input element.
type InputElement interface {
	Build() (label, key, sType string, elems []js.Value)
	Update(this js.Value, params []js.Value) any
}

// update is a map of functions that update the state of the elements
var update map[string]func(any)

// ids is a map of ids that are currently in use.
var ids = map[string]bool{}

// Run runs the main function with the given elements.
func Run(mainFunc func(), elements ...InputElement) {
	log.Debug("program started...")

	// Store the time that the program started.
	start := time.Now()

	// Create a map of known elements and update callbacks
	seen := make(map[string]bool)
	update = make(map[string]func(any))

	// get the document object, and from it the options panel for the page
	document := js.Global().Get("parent").Get("document")
	optionPanel := document.Call("getElementById", "options")

	log.Debug("begin building interface...")

	// For each of the specified elements...
	for _, element := range elements {
		// Build the element
		l, k, t, elems := element.Build()

		log.Debug("got element:", "label", l, "key", k, "type", t)

		// If the element type has not been seen before...
		if !seen[t] {
			log.Debug("discovered element, adding update function...", "type", t)

			// Flag the element type as being seen.
			seen[t] = true
			// Add the update function to the docu
			//ment
			js.Global().Get("parent").Set(t+"Update", js.FuncOf(element.Update))
		}

		// Create the label for the element
		label := document.Call("createElement", "label")
		label.Call("setAttribute", "for", k)
		label.Call("setAttribute", "class", "optionLabel")
		label.Call("appendChild", document.Call("createTextNode", l))

		// Create the option element
		option := document.Call("createElement", "div")
		option.Call("setAttribute", "class", "option")

		// Create the input box
		inputBox := document.Call("createElement", "div")
		inputBox.Call("setAttribute", "class", "inputBox")

		// Put all the elements into the input box
		for _, elem := range elems {
			inputBox.Call("appendChild", elem)
		}

		// Put the label and input box into the option
		option.Call("appendChild", label)
		option.Call("appendChild", inputBox)

		// Put the option into the options panel
		optionPanel.Call("appendChild", option)
	}

	log.Debug("done building interface...", "time", time.Since(start).String())

	// Run the main function
	mainFunc()

	// Wait for the program to exit.
	<-make(chan bool)
}
