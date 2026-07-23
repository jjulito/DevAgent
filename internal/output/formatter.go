package output

import (
	"github.com/fatih/color"
)

var (
	accent  = color.New(color.FgCyan, color.Bold)
	success = color.New(color.FgGreen, color.Bold)
	warn    = color.New(color.FgYellow)
	errC    = color.New(color.FgRed, color.Bold)
	dim     = color.New(color.FgHiBlack)
)

func Banner() {
	accent.Println("╔══════════════════════════════════╗")
	accent.Println("║           DevAgent CLI           ║")
	accent.Println("╚══════════════════════════════════╝")
}
