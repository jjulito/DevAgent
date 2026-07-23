package output

import (
	"fmt"
	"os"
	"strings"

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

func Info(msg string) {
	accent.Printf("▸ %s\n", msg)
}

func Success(msg string) {
	success.Printf("✔ %s\n", msg)
}

func Warn(msg string) {
	warn.Printf("⚠ %s\n", msg)
}

func Error(msg string) {
	errC.Fprintf(os.Stderr, "✖ %s\n", msg)
}

func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

func Dim(msg string) {
	dim.Println(msg)
}

func Divider() {
	dim.Println(strings.Repeat("─", 50))
}

func ModelTag(provider, model string) {
	dim.Printf("[%s/%s]\n", provider, model)
}
