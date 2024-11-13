package hiddify_extension

import (
	"context"
	"time"

	"github.com/hiddify/hiddify-core/config"
	"github.com/sagernet/sing-box/option"

	"github.com/fatih/color"
	ex "github.com/hiddify/hiddify-core/extension"
	ui "github.com/hiddify/hiddify-core/extension/ui"
)

// Color definitions for console output
var (
	red    = color.New(color.FgRed).Add(color.Bold)
	green  = color.New(color.FgGreen).Add(color.Underline)
	yellow = color.New(color.FgYellow)
)

// V2hiddifyExtensionData holds the data specific to V2hiddifyExtension
type V2hiddifyExtensionData struct {
	Count int `json:"count"` // Number of counts for the extension
}

// V2hiddifyExtension represents the core functionality of the extension
type V2hiddifyExtension struct {
	ex.Base[V2hiddifyExtensionData]                    // Embedding base extension functionality
	cancel                        context.CancelFunc // Function to cancel background tasks
	console                       string             // Stores console output
}

// backgroundTask runs a task in the background, updating the console at intervals
func (e *V2hiddifyExtension) backgroundTask(ctx context.Context) {
	for count := 1; count <= e.Base.Data.Count; count++ {
		select {
		case <-ctx.Done(): // If context is done (cancel is pressed), exit the task
			e.cancel = nil
			e.addAndUpdateConsole(red.Sprint("Background Task Canceled")) // Notify cancellation
			return
		case <-time.After(1 * time.Second): // Wait for a second before the next iteration
			e.addAndUpdateConsole(red.Sprint(count), yellow.Sprint(" Background task ", count, " working..."))
		}
	}
	e.cancel = nil
	e.addAndUpdateConsole(green.Sprint("Background Task Finished Successfully")) // Task completion message
}

// SubmitData processes and validates form submission data
func (e *V2hiddifyExtension) SubmitData(button string, data map[string]string) error {
	switch button {
	case ui.ButtonDialogOk, ui.ButtonDialogClose:
		return nil
	case ui.ButtonCancel:
		return e.stop()
	case ui.ButtonSubmit:
		if err := e.setFormData(data); err != nil {
			e.ShowMessage("Invalid data", err.Error()) // Show error message if data is invalid
			return err                                 // Return the error
		}

		// stop any ongoing background task
		if e.cancel != nil {
			e.cancel()
		}

		// Create a new context for the task and store the cancel function
		ctx, cancel := context.WithCancel(context.Background())
		e.cancel = cancel

		// Run the background task concurrently
		go e.backgroundTask(ctx)

		return nil

	default:
		// Show message for undefined button actions
		return e.ShowMessage("Button "+button+" is pressed", "No action is defined for this button")
	}
}

// Cancel stops the ongoing background task if it exists
func (e *V2hiddifyExtension) stop() error {
	if e.cancel != nil {
		e.cancel()     // Cancel the task
		e.cancel = nil // Clear the cancel function
	}
	return nil // Return nil after cancellation
}

// Stop is called when the extension is closed
func (e *V2hiddifyExtension) Close() error {
	return e.stop() // Simply delegate to stop
}

// To Modify user's config before connecting, you can use this function
func (e *V2hiddifyExtension) BeforeAppConnect(hiddifySettings *config.HiddifyOptions, singconfig *option.Options) error {
	return nil
}

// NewV2hiddifyExtension initializes a new instance of V2hiddifyExtension with default values
func NewV2hiddifyExtension() ex.Extension {
	return &V2hiddifyExtension{
		Base: ex.Base[V2hiddifyExtensionData]{
			Data: V2hiddifyExtensionData{ // Set default data
				Count: 4, // Default count value
			},
		},
		console: yellow.Sprint("Welcome to ") + green.Sprint("V2hiddifyExtension\n"), // Default message
	}
}

// init registers the extension with the provided metadata
func init() {
	ex.RegisterExtension(
		ex.ExtensionFactory{
			Id:          "github.com/V2hiddify/v2hiddifyextension/hiddify_extension", // Package identifier
			Title:       "V2hiddifyExtension",                                                         // Display title of the extension
			Description: "Awesome Extension v2hiddifyextension created by V2hiddify",                                                     // Brief description of the extension
			Builder:     NewV2hiddifyExtension,                                                       // Function to create a new instance
		},
	)
}
