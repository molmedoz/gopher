package progress

import (
	"os"

	"golang.org/x/term"
)

// getTerminalSize returns the width and height of the terminal
// Returns 0, 0, error if terminal size cannot be determined
func getTerminalSize() (width, height int, err error) {
	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	return width, height, err
}
