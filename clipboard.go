package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateClipboardScript creates the PowerShell helper script for clipboard operations.
// It reads the file path from %PATHVAR% and format from %FORMATVAR% environment variables.
// This avoids complex quoting issues when the command is invoked from the Windows Registry.
//
// Format notes (PowerShell -replace with single-quoted strings):
//   - pattern: '\\'    -> regex "match one literal backslash"
//   - replacement: '\\' -> two literal backslash chars, inserted as-is per char
//     So '-replace "\\","\\"' doubles every single backslash, which is what we want
//     for the "qbck" (quoted-backslash) format so the output can be pasted directly
//     into a code/JSON string literal.
func generateClipboardScript() string {
	return `$path = $env:PATHVAR
$fmt = $env:FORMATVAR
switch ($fmt) {
  "fwd"  { $path = $path -replace '\\', '/' }
  "bck"  { $path = $path -replace '/', '\' }
  "qfwd" { $path = '"' + ($path -replace '\\', '/') + '"' }
  "qbck" { $path = '"' + ($path -replace '\\', '\\') + '"' }
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

// getClipboardVBSPath returns the full path to the silent-launcher VBS file.
func getClipboardVBSPath() string {
	exePath, _ := os.Executable()
	return filepath.Join(filepath.Dir(exePath), clipboardVBSFile)
}

// generateClipboardVBS returns the source of the silent launcher. It is invoked
// from the context menu via wscript.exe (GUI subsystem, no console). The VBS
// then launches PowerShell with WindowStyle=0, so no console window appears at
// any point. Arguments: <filePath> <formatFlag>.
func generateClipboardVBS() string {
	return `' wcpath-clip.vbs - hidden launcher for wcpath-clip.ps1
' Invoked as: wscript.exe wcpath-clip.vbs <filePath> <formatFlag>
Set args = WScript.Arguments
If args.Count < 2 Then WScript.Quit 1

Set fso = CreateObject("Scripting.FileSystemObject")
psScript = fso.BuildPath(fso.GetParentFolderName(WScript.ScriptFullName), "wcpath-clip.ps1")

q = Chr(34)
' Escape single quotes in the path so PowerShell's single-quoted string stays valid.
safePath = Replace(args(0), "'", "''")
cmd = "powershell.exe -NoProfile -ExecutionPolicy Bypass -Command " & q & _
      "$env:PATHVAR='" & safePath & "'; " & _
      "$env:FORMATVAR='" & args(1) & "'; " & _
      "& '" & psScript & "'" & q

' WindowStyle=0 (hidden), bWaitOnReturn=False (fire-and-forget so the context
' menu returns immediately even if PowerShell takes a beat to start).
CreateObject("WScript.Shell").Run cmd, 0, False
`
}

// writeClipboardScript writes the clipboard helper script (ps1) and the silent
// launcher (vbs) to disk, next to the Win-Copath executable. The vbs is needed
// only for the context-menu path (to avoid a console flash); the ps1 is also
// used directly by `wcpath copy <fmt> <path>` from the command line.
func writeClipboardScript() error {
	scriptPath, err := getClipboardScriptPath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(scriptPath, []byte(generateClipboardScript()), 0644); err != nil {
		return err
	}
	vbsPath := getClipboardVBSPath()
	return os.WriteFile(vbsPath, []byte(generateClipboardVBS()), 0644)
}

// removeClipboardScript deletes the clipboard helper ps1 and the silent vbs
// launcher from disk. Missing files are not treated as errors.
func removeClipboardScript() error {
	scriptPath, err := getClipboardScriptPath()
	if err != nil {
		return err
	}
	if err := os.Remove(scriptPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove script: %w", err)
	}
	vbsPath := getClipboardVBSPath()
	if err := os.Remove(vbsPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove vbs: %w", err)
	}
	return nil
}

// buildClipCmd builds the command line string for the registry "Command" value.
//
// To avoid the cmd.exe console flash that -WindowStyle Hidden cannot suppress,
// the registry value launches wscript.exe (a GUI-subsystem script host, no
// console) which runs wcpath-clip.vbs. The VBS in turn uses WScript.Shell.Run
// with WindowStyle=0 (hidden) to invoke PowerShell, so the entire chain runs
// silently. The PowerShell command sets PATHVAR and FORMATVAR env vars and
// dot-sources the sibling wcpath-clip.ps1 script.
//
// Format (as stored in registry):
//
//	wscript.exe "<vbs>" "%V" "<flag>"
func buildClipCmd(flag string) string {
	vbsPath := getClipboardVBSPath()
	return fmt.Sprintf(`wscript.exe "%s" "%%V" "%s"`, vbsPath, flag)
}
