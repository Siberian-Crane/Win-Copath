$ErrorActionPreference = 'Stop'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$url      = 'https://github.com/Siberian-Crane/Win-Copath/releases/download/v1.0.0/wcpath_windows_amd64.exe'
$url64    = $url

$packageArgs = @{
  packageName    = $env:ChocolateyPackageName
  unzipLocation  = $toolsDir
  url            = $url
  url64bit       = $url64
  checksum       = 'B4FFC89AEC56C9BBB4BA4406F3D2A85E6EA43651AB6420C486748C41B1485D53'
  checksumType   = 'sha256'
  checksum64     = 'B4FFC89AEC56C9BBB4BA4406F3D2A85E6EA43651AB6420C486748C41B1485D53'
  checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs
