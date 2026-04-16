$ErrorActionPreference = "Stop"

$Repo = "CollectCall/jira-plans-teams-dc-to-dc-migrator"
$Version = if ($env:VERSION) { $env:VERSION } else { "latest" }
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { Join-Path $HOME "bin" }
$LatestReleaseApi = "https://api.github.com/repos/$Repo/releases/latest"

$arch = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture) {
    "X64" { "amd64" }
    "Arm64" { "arm64" }
    default { throw "Unsupported Windows architecture: $([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture)" }
}

$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ("teams-migrator-" + [System.Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tmpDir | Out-Null

try {
    $ResolvedVersion = if ($Version -eq "latest") {
        (Invoke-RestMethod -Uri $LatestReleaseApi).tag_name
    } else {
        $Version
    }
    if ([string]::IsNullOrWhiteSpace($ResolvedVersion)) {
        throw "Failed to determine the latest release version."
    }

    $archive = "teams-migrator_${ResolvedVersion}_windows_${arch}.zip"
    $url = "https://github.com/$Repo/releases/download/$ResolvedVersion/$archive"
    $zipPath = Join-Path $tmpDir $archive
    Invoke-WebRequest -Uri $url -OutFile $zipPath

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    Expand-Archive -Path $zipPath -DestinationPath $tmpDir -Force
    Copy-Item (Join-Path $tmpDir "teams-migrator.exe") (Join-Path $InstallDir "teams-migrator.exe") -Force

    Write-Host "Installed teams-migrator.exe to $InstallDir"
    if (-not (($env:PATH -split ";") -contains $InstallDir)) {
        Write-Host "Add this directory to PATH if needed:"
        Write-Host "  $InstallDir"
    }
}
finally {
    Remove-Item -Recurse -Force $tmpDir
}
