package wshim

import (
	"github.com/charmbracelet/log"
	"syscall/js"
	"time"
)

const (
	timeout = 10 * time.Second
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

// SetLogLevel sets the log level of the program.
func SetLogLevel(level log.Level) {
	log.SetLevel(level)
}

// Run runs the main function with the given elements.
func Run(mainFunc func(), elements ...InputElement) {
	log.Debug("program started...")

	// Store the time that the program started.
	programStart := time.Now()

	// Create a map of known elements and update callbacks
	seen := make(map[string]bool)
	update = make(map[string]func(any))

	// get the parent document object, and from it the options panel for the page
	parentDocument := js.Global().Get("parent").Get("document")

	// If the parent document is not defined, the page may not have been loaded yet.
	retryStart := time.Now()
	for checkUndefined(parentDocument) && time.Since(retryStart) <= timeout {
		log.Debug("parent document not found, retrying...")
		parentDocument = js.Global().Get("parent").Get("document")
		time.Sleep(100 * time.Millisecond)
	}

	// If the parent document is null then the page is likely not inside an iframe.
	// In this case the program should run, but the elements should not be added.
	if checkUndefined(parentDocument) {
		log.Error("parent document not found, running without options...")

		// Run the main function.
		mainFunc()

		// Wait for the program to exit.
		<-make(chan bool)

		return
	}

	// Get the options panel, if the page is inside an iframe.
	optionsPanel := parentDocument.Call("getElementById", "options")

	// If the optionsPanel is null, then the page has not been loaded yet, wait for it to load and try again.
	retryStart = time.Now()
	for checkUndefined(optionsPanel) && time.Since(retryStart) <= timeout {
		log.Debug("options panel not found, retrying...")
		optionsPanel = parentDocument.Call("getElementById", "options")
		time.Sleep(100 * time.Millisecond)
	}

	// If the optionsPanel is still null, throw an error.
	if checkUndefined(optionsPanel) {
		log.Error("options panel not found, running without options...")

		// Run the main function.
		mainFunc()

		// Wait for the program to exit.
		<-make(chan bool)

		return
	}

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
		label := parentDocument.Call("createElement", "label")
		label.Call("setAttribute", "for", k)
		label.Call("setAttribute", "class", "optionLabel")
		label.Call("appendChild", parentDocument.Call("createTextNode", l))

		// Create the option element
		option := parentDocument.Call("createElement", "div")
		option.Call("setAttribute", "class", "option")

		// Create the input box
		inputBox := parentDocument.Call("createElement", "div")
		inputBox.Call("setAttribute", "class", "inputBox")

		// Put all the elements into the input box
		for _, elem := range elems {
			inputBox.Call("appendChild", elem)
		}

		// Put the label and input box into the option
		option.Call("appendChild", label)
		option.Call("appendChild", inputBox)

		// Put the option into the options panel
		optionsPanel.Call("appendChild", option)
	}

	log.Debug("done building interface...", "time", time.Since(programStart).String())

	// Run the main function
	mainFunc()

	// Wait for the program to exit.
	<-make(chan bool)
}
