package client

import "github.com/pterm/pterm"

// inputWithResult creates a new input in command line
// and returns result string and error.
func inputWithResult(label string) (string, error) {
	input := pterm.DefaultInteractiveTextInput
	input.DefaultText = label
	result, err := input.Show()
	if err != nil {
		return result, err
	}

	pterm.Println()

	return result, nil
}
