package output

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

// getTerminalWidth returns the width of the terminal or a default fallback of 100.
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 100
	}
	return width
}

// RenderMarkdown renders the given markdown string formatted for terminal output using glamour.
// It uses auto style (detects light/dark mode) and sets a dynamic word wrap width based on terminal size.
func RenderMarkdown(markdown string) (string, error) {
	width := getTerminalWidth()

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create glamour renderer: %w", err)
	}

	out, err := r.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return out, nil
}

// PrintMarkdown renders markdown and prints it to stdout.
// Falls back to printing raw markdown if rendering fails.
func PrintMarkdown(markdown string) error {
	rendered, err := RenderMarkdown(markdown)
	if err != nil {
		fmt.Print(markdown)
		return err
	}
	fmt.Print(rendered)
	return nil
}
