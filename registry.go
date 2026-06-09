package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// checkAdmin checks if the current process has administrator privileges.
func checkAdmin() bool {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command",
		`([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)`)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "True"
}

// regKeyExists checks if a Windows registry key exists using reg query.
func regKeyExists(keyPath string) bool {
	cmd := exec.Command("reg", "query", keyPath)
	err := cmd.Run()
	return err == nil
}

// regDeleteKey removes a registry key and all its subkeys using reg delete.
func regDeleteKey(keyPath string) error {
	cmd := exec.Command("reg", "delete", keyPath, "/f")
	return cmd.Run()
}

// menuExists checks if the Win-Copath context menu entries exist.
func menuExists() bool {
	return regKeyExists(`HKCR\Directory\shell\` + menuKeyName)
}

// addMenuEntries creates the Win-Copath context menu entries in the Windows Registry.
// It adds entries under:
//   - HKCR\Directory\shell (right-click on folders)
//   - HKCR\Directory\Background\shell (right-click on empty space in folders)
//   - HKCR\*\shell (right-click on files)
func addMenuEntries() error {
	// Write the clipboard helper script first
	if err := writeClipboardScript(); err != nil {
		return fmt.Errorf("failed to write clipboard script: %w", err)
	}

	// Add entries for each supported context
	contexts := []string{
		`HKCR\Directory\shell\` + menuKeyName,
		`HKCR\Directory\Background\shell\` + menuKeyName,
		`HKCR\*\shell\` + menuKeyName,
	}

	for _, contextPath := range contexts {
		if err := addMenuEntry(contextPath); err != nil {
			return fmt.Errorf("failed to add menu entry for %s: %w", contextPath, err)
		}
	}

	return nil
}

// addMenuEntry adds the Win-Copath menu and all its sub-entries to a specific registry path.
func addMenuEntry(basePath string) error {
	// Create the main "Win Copath" key with SubCommands
	// (default) = display name
	// SubCommands = indicates this is a submenu
	entries := []struct {
		key   string
		name  string
		value interface{}
		typ   uint32
	}{
		{basePath, "", appName, 1},       // REG_SZ: display name
		{basePath, "SubCommands", "", 1}, // REG_SZ: marks as submenu
	}

	for _, e := range entries {
		if err := setRegStringValue(e.key, e.name, e.value.(string)); err != nil {
			return err
		}
	}

	// Add each path format as a sub-command
	for _, pf := range pathFormats {
		subKey := fmt.Sprintf(`%s\shell\%s`, basePath, pf.flag)

		// Set display name and command
		if err := setRegStringValue(subKey, "", pf.label); err != nil {
			return err
		}
		if err := setRegStringValue(subKey, "Command", buildClipCmd(pf.flag)); err != nil {
			return err
		}
	}

	return nil
}

// removeMenuEntries removes the Win-Copath context menu entries from the Windows Registry.
func removeMenuEntries() error {
	// Remove from each context
	contexts := []string{
		`HKCR\Directory\shell\` + menuKeyName,
		`HKCR\Directory\Background\shell\` + menuKeyName,
		`HKCR\*\shell\` + menuKeyName,
	}

	for _, contextPath := range contexts {
		if err := regDeleteKey(contextPath); err != nil {
			// Ignore errors if key doesn't exist
			fmt.Printf("Note: Could not remove %s (may not exist)\n", contextPath)
		}
	}

	// Remove the clipboard script
	return removeClipboardScript()
}

// setRegStringValue writes a REG_SZ string value to the Windows Registry.
func setRegStringValue(keyPath, name, value string) error {
	cmd := exec.Command("reg", "add", keyPath, "/v", name, "/t", "REG_SZ", "/d", value, "/f")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reg add failed: %s - %w", string(output), err)
	}
	return nil
}
