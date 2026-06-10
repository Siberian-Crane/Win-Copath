$ErrorActionPreference = 'Stop'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

# Remove the context menu entries if the user previously enabled them
# (only runs if the binary is still on disk; otherwise the user can rerun `wcpath off` manually).
$wcpathCmd = Get-Command wcpath -ErrorAction SilentlyContinue
if ($wcpathCmd) {
  try {
    & wcpath off | Out-Null
  } catch {
    Write-Warning "Failed to remove context menu entries during uninstall: $_"
  }
}

Uninstall-ChocolateyZipPackage $env:ChocolateyPackageName
