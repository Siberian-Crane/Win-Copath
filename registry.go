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

// contextMenusKey is the absolute HKCR path that holds the shared submenu items
// for the Win-Copath cascade. All three top-level parent menu keys
// (HKCR\Directory\shell\Win-Copath, HKCR\Directory\Background\shell\Win-Copath,
// HKCR\*\shell\Win-Copath) point here via their ExtendedSubCommandsKey value.
//
// This is the same pattern PowerShell7, VS Code, and Git for Windows use: a
// single HKCR-only tree, no HKLM\CommandStore involvement, and one set of
// submenu items shared across all parent contexts.
//
// Docs: https://learn.microsoft.com/en-us/windows/win32/shell/how-to--create-cascading-menus-with-the-subcommands-registry-entry
const contextMenusKey = `HKCR\Directory\ContextMenus\Win-Copath`

// contextMenusRelative is the path written to each parent's ExtendedSubCommandsKey
// value. Windows resolves it relative to HKCR root, so all three parent
// contexts (Directory, Directory\Background, *) share the same contextMenusKey.
const contextMenusRelative = `Directory\ContextMenus\Win-Copath`

// menuExists checks if the Win Copath context menu entries exist.
func menuExists() bool {
	return regKeyExists(`HKCR\Directory\shell\` + menuKeyName)
}

// addMenuEntries creates the Win Copath context menu entries in the Windows Registry.
// It adds entries under:
//   - HKCR\Directory\shell (right-click on folders)
//   - HKCR\Directory\Background\shell (right-click on empty space in folders)
//   - HKCR\*\shell (right-click on files)
//
// The actual submenu items are defined once under HKCR\Directory\ContextMenus\Win-Copath
// and all three parent contexts reference it via ExtendedSubCommandsKey.
func addMenuEntries() error {
	// Write the clipboard helper script first
	if err := writeClipboardScript(); err != nil {
		return fmt.Errorf("failed to write clipboard script: %w", err)
	}

	// Create the shared submenu items once. Done before the parents reference it.
	if err := addContextMenuItems(); err != nil {
		return fmt.Errorf("failed to add context menu items: %w", err)
	}

	// Add parent entries that point to the shared ContextMenus block.
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

// addContextMenuItems populates the shared HKCR\Directory\ContextMenus\Win-Copath
// location with one submenu entry per path format. Each entry has its display
// label as the (Default) value and its command line under \command.
func addContextMenuItems() error {
	for _, pf := range pathFormats {
		subKey := fmt.Sprintf(`%s\shell\%s`, contextMenusKey, pf.flag)
		if err := setRegDefaultValue(subKey, pf.label); err != nil {
			return err
		}
		if err := setRegDefaultValue(subKey+`\command`, buildClipCmd(pf.flag)); err != nil {
			return err
		}
	}
	return nil
}

// addMenuEntry adds the Win-Copath cascade parent to a specific registry path.
//
// Pattern (matches PowerShell7 / VSCode / Git for Windows):
//   - MUIVerb              = "Win Copath"   (display label)
//   - ExtendedSubCommandsKey = "Directory\ContextMenus\Win-Copath"
//
// Windows resolves ExtendedSubCommandsKey relative to HKCR root, so this single
// value works for all three parent contexts (Directory, Directory\Background, *).
func addMenuEntry(basePath string) error {
	if err := setRegStringValue(basePath, "MUIVerb", appName); err != nil {
		return err
	}
	return setRegStringValue(basePath, "ExtendedSubCommandsKey", contextMenusRelative)
}

// setRegDefaultValue writes the default (unnamed) REG_SZ value to a registry key.
// Uses /ve flag to set the unnamed (default) value.
func setRegDefaultValue(keyPath, value string) error {
	cmd := exec.Command("reg", "add", keyPath, "/ve", "/t", "REG_SZ", "/d", value, "/f")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reg add failed: %s - %w", string(output), err)
	}
	return nil
}

// setRegStringValue writes a named REG_SZ value to a registry key.
func setRegStringValue(keyPath, name, value string) error {
	cmd := exec.Command("reg", "add", keyPath, "/v", name, "/t", "REG_SZ", "/d", value, "/f")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reg add failed: %s - %w", string(output), err)
	}
	return nil
}

// removeMenuEntries removes the Win-Copath context menu entries from the Windows Registry.
func removeMenuEntries() error {
	// Remove the parent entry from each context
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

	// Remove the shared submenu items block (one reg delete cleans up all
	// pathFormat children at once).
	if err := regDeleteKey(contextMenusKey); err != nil {
		fmt.Printf("Note: Could not remove %s (may not exist)\n", contextMenusKey)
	}

	// Remove the clipboard script
	return removeClipboardScript()
}
