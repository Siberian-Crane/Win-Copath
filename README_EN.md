# Win-Copath

[README.md](https://github.com/Siberian-Crane/Win-Copath/blob/main/README.md)

A Windows right-click menu path copy tool. Adds a "Win Copath" option to the context menu of files and folders, supporting four path formats for copying.

## Features

Right-click on a file or folder to display the "Win Copath" submenu, which contains four options:

| Option                       | Format               | Example                       |
|------------------------------|----------------------|-------------------------------|
| `/ (Forward Slash)`          | Forward slash        | `C:/Users/test/file.txt`      |
| `\ (Backslash)`              | Backslash            | `C:\Users\test\file.txt`      |
| `"/" (Quoted Forward Slash)` | Quoted forward slash | `"C:/Users/test/file.txt"`    |
| `"\\" (Quoted Backslash)`    | Quoted backslash     | `"C:\\Users\\test\\file.txt"` |

## Installation

### Chocolatey

```powershell
choco install wincopath
```

### Manual Installation

1. Download `wcpath_windows_amd64.exe` from [Releases](https://github.com/Siberian-Crane/Win-Copath/releases)
2. Place `wcpath.exe` in a directory already added to PATH (e.g. `C:\Windows\System32`)

## Usage

After installation, **the right-click menu will NOT be modified by default**. You need to enable it manually:

```powershell
# Run PowerShell as Administrator
wcpath on     # Add Win Copath to the right-click menu
wcpath off    # Remove Win Copath from the right-click menu
wcpath help   # Show help
```

Once enabled, right-click on any file or folder to see the "Win Copath" submenu.

## Requirements

- Windows 10 or later
- Administrator privileges (required to modify the registry)

## Development

```bash
# Build
go build -o wcpath.exe .

# Test the copy command
.\wcpath.exe copy fwd "C:\test\file.txt"
```

## Registry Locations

After enabling the context menu (`wcpath on`), the tool creates registry entries at the following locations:

- `Computer\HKEY_CLASSES_ROOT\Directory\Background\shell\Win Copath`
- `Computer\HKEY_CLASSES_ROOT\Directory\ContextMenus\Win-Copath`

> Note: Two additional parent keys are also created — `HKEY_CLASSES_ROOT\Directory\shell\Win Copath` (right-click on folders) and `HKEY_CLASSES_ROOT\*\shell\Win Copath` (right-click on files). All three parents reference the same `Directory\ContextMenus\Win-Copath` submenu definition via `ExtendedSubCommandsKey`, so the submenu items are registered only once.

To inspect or clean up manually, open `regedit` and navigate to the paths above, or simply run `wcpath off` to remove them.

## Development Collaboration

This project was collaboratively developed by two models: `mimo-v2.5-pro` and `MiniMax-M3`.

## References

The context menu registry design of this tool references the implementation of PowerShell 7, using `ExtendedSubCommandsKey` to share the submenu definition and avoid duplicate registration.

## License

MIT
