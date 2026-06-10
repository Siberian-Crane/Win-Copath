$ErrorActionPreference = 'Stop'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

# Try to remove the context menu entries before deleting the binary.
# Prefer the binary co-located with this script (tools\wcpath.exe) so the call
# works even when the user has not refreshed their PATH in the current
# PowerShell session. Fall back to whatever is on PATH.
$wcpathExe = Join-Path $toolsDir 'wcpath.exe'
if (-not (Test-Path $wcpathExe)) { $wcpathExe = (Get-Command wcpath -ErrorAction SilentlyContinue).Source }
if ($wcpathExe) {
  try {
    & $wcpathExe off | Out-Null
  } catch {
    Write-Warning "Failed to remove context menu entries during uninstall: $_"
  }
} else {
  Write-Warning "wcpath.exe not found; registry entries may need to be removed manually (run 'wcpath off' or delete HKCR\Directory\ContextMenus\Win-Copath)."
}

Uninstall-ChocolateyZipPackage $env:ChocolateyPackageName
