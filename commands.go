package main

import (
	"fmt"
	"os"
	"os/exec"
)

// handleOn adds the Win Copath context menu entries to the Windows Registry.
// It checks if entries already exist to avoid duplicates.
func handleOn() {
	fmt.Println("Adding Win Copath to context menu...")

	// Check admin privileges
	if !checkAdmin() {
		_, _ = fmt.Fprintln(os.Stderr, "Error: Administrator privileges required. Please run PowerShell as Administrator.")
		os.Exit(1)
	}

	// Check if already added
	if menuExists() {
		fmt.Println("Win Copath is already in the context menu. Nothing to do.")
		return
	}

	// Add menu entries
	if err := addMenuEntries(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done! Right-click on any file or folder to use Win Copath.")
}

// handleOff removes the Win Copath context menu entries from the Windows Registry.
func handleOff() {
	fmt.Println("Removing Win Copath from context menu...")

	// Check admin privileges
	if !checkAdmin() {
		_, _ = fmt.Fprintln(os.Stderr, "Error: Administrator privileges required. Please run PowerShell as Administrator.")
		os.Exit(1)
	}

	if !menuExists() {
		fmt.Println("Win Copath is not in the context menu. Nothing to do.")
		return
	}

	if err := removeMenuEntries(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done! Win Copath has been removed from the context menu.")
}

// handleCopy executes the clipboard copy operation using the PowerShell helper script.
// It sets environment variables for the file path and format, then runs the script.
func handleCopy(args []string) {
	if len(args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: wcpath copy <format> <filepath>")
		_, _ = fmt.Fprintln(os.Stderr, "Formats: fwd, bck, qfwd, qbck")
		os.Exit(1)
	}

	format := args[0]
	filePath := args[1]

	// Validate format
	validFormats := map[string]bool{"fwd": true, "bck": true, "qfwd": true, "qbck": true}
	if !validFormats[format] {
		_, _ = fmt.Fprintf(os.Stderr, "Error: invalid format '%s'. Valid formats: fwd, bck, qfwd, qbck\n", format)
		os.Exit(1)
	}

	scriptPath, err := getClipboardScriptPath()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Ensure the script exists (create it if missing)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		if err := writeClipboardScript(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating clipboard script: %v\n", err)
			os.Exit(1)
		}
	}

	// Build and execute the PowerShell command
	// Use environment variables to pass the path (avoids quoting issues)
	psArgs := []string{
		"-NoProfile",
		"-ExecutionPolicy", "Bypass",
		"-File", scriptPath,
	}

	cmd := exec.Command("powershell.exe", psArgs...)
	cmd.Env = append(os.Environ(),
		"PATHVAR="+filePath,
		"FORMATVAR="+format,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
