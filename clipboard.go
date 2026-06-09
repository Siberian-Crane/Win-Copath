package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateClipboardScript creates the PowerShell helper script for clipboard operations.
// It reads the file path from %PATHVAR% and format from %FORMATVAR% environment variables.
// This avoids complex quoting issues when the command is invoked from the Windows Registry.
func generateClipboardScript() string {
	return `$path = $env:PATHVAR
$fmt = $env:FORMATVAR
switch ($fmt) {
  "fwd"  { $path = $path -replace '\\', '/' }
  "bck"  { $path = $path -replace '/', '\' }
  "qfwd" { $path = '"' + ($path -replace '\\', '/') + '"' }
  "qbck" { $path = '"' + ($path -replace '/', '\') + '"' }
}
$path | clip.exe`
}

// getClipboardScriptPath returns the full path to the clipboard helper script.
func getClipboardScriptPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Join(filepath.Dir(exePath), clipboardFile), nil
}

// writeClipboardScript writes the clipboard helper script to disk.
func writeClipboardScript() error {
	scriptPath, err := getClipboardScriptPath()
	if err != nil {
		return err
	}
	return os.WriteFile(scriptPath, []byte(generateClipboardScript()), 0644)
}

// removeClipboardScript deletes the clipboard helper script from disk.
func removeClipboardScript() error {
	scriptPath, err := getClipboardScriptPath()
	if err != nil {
		return err
	}
	if err := os.Remove(scriptPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove script: %w", err)
	}
	return nil
}

// buildClipCmd builds the command line string for the registry "Command" value.
//
// The registry command is executed by cmd.exe, which passes it to powershell.exe.
// We use -Command with escaped inner quotes so that:
// 1. cmd.exe strips the outer quotes
// 2. PowerShell receives the inner-quoted command string and executes it
// 3. $env:PATHVAR and $env:FORMATVAR are set in the new PowerShell session
// 4. The clipboard script dot-sources and uses these environment variables
//
// Format (as stored in registry):
//
//	powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "$env:PATHVAR='%V'; $env:FORMATVAR='<flag>'; . '<script>'"
func buildClipCmd(flag string) string {
	exePath, _ := os.Executable()
	scriptPath := filepath.Join(filepath.Dir(exePath), clipboardFile)
	// \" in Go raw string produces literal backslash-double-quote in the output,
	// which cmd.exe interprets as an escaped quote (literal quote passed to PowerShell).
	return fmt.Sprintf(
		`powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "$env:PATHVAR='%%V'; $env:FORMATVAR='%s'; . '%s'"`,
		flag, scriptPath,
	)
}
