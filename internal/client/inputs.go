package client

import "github.com/pterm/pterm"

// inputWithReslult creates a new input in command line
// and returns result string and error.
func inputWithReslult(label string) (string, error) {
	input := pterm.DefaultInteractiveTextInput
	input.DefaultText = label
	result, err := input.Show()
	if err != nil {
		return result, err
	}

	pterm.Println()

	return result, nil
}
