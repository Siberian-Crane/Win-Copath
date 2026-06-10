# Win-Copath

A Windows right-click menu path copy tool. Adds a "Win Copath" option to the context menu of files and folders, supporting four path formats for copying.

## Features

Right-click on a file or folder to display the "Win Copath" submenu, which contains four options:

| Option | Format | Example |
|----------------|---------|-------------------------------|
| `/ (Forward Slash)` | Forward slash | `C:/Users/test/file.txt` |
| `\ (Backslash)` | Backslash | `C:\Users\test\file.txt` |
| `"/" (Quoted Forward Slash)` | Quoted forward slash | `"C:/Users/test/file.txt"` |
| `"\\" (Quoted Backslash)` | Quoted backslash | `"C:\\Users\\test\\file.txt"` |

## Installation

### WinGet (Recommended)

```powershell
winget install Siberian-Crane.WinCopath
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

## Development Collaboration

This project was collaboratively developed by two models: `mimo-v2.5-pro` and `MiniMax-M3`.

## License

MIT
