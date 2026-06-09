package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// Only support Windows
	if runtime.GOOS != "windows" {
		_, _ = fmt.Fprintln(os.Stderr, "Error: Win Copath only supports Windows")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	// The "copy" subcommand is invoked internally by the context menu entry.
	// It should not be advertised in the help text.
	switch os.Args[1] {
	case "on":
		handleOn()
	case "off":
		handleOff()
	case "copy":
		handleCopy(os.Args[2:])
	case "help", "--help", "-h":
		printUsage()
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Win Copath - Windows path copy tool for the context menu

Usage:
  wcpath on     Add Win Copath to the right-click context menu
  wcpath off    Remove Win Copath from the right-click context menu
  wcpath help   Show this help message

After running "wcpath on", you can right-click on any file or folder
to see the "Win Copath" submenu with four path format options:
  / (正斜杠)            e.g. C:/Users/test/file.txt
  \ (反斜杠)            e.g. C:\Users\test\file.txt
  "/\" (引号正斜杠)     e.g. "C:/Users/test/file.txt"
  "\\" (引号反斜杠)      e.g. "C:\\Users\\test\\file.txt"

Requires Administrator privileges to modify the registry.`)
}
