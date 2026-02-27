# Claude Code Status Line - Windows Install Script
# Usage: irm https://raw.githubusercontent.com/EvanPluchart/claude-code-status-line/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "EvanPluchart/claude-code-status-line"
$Binary = "claude-code-status-line"

# Detect architecture
$Arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "Unsupported: 32-bit systems are not supported."
    exit 1
}

# Get latest version
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$Version = $Release.tag_name -replace "^v", ""

if (-not $Version) {
    Write-Error "Could not determine latest version."
    exit 1
}

$Url = "https://github.com/$Repo/releases/download/v$Version/${Binary}_${Version}_windows_${Arch}.zip"
$InstallDir = "$env:LOCALAPPDATA\Programs\$Binary"

Write-Host "Installing $Binary v$Version (windows/$Arch)..."

# Create temp directory
$TmpDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }

try {
    # Download and extract
    $ZipPath = Join-Path $TmpDir "archive.zip"
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath
    Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

    # Create install directory
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null

    # Copy binary
    Copy-Item -Path (Join-Path $TmpDir "$Binary.exe") -Destination (Join-Path $InstallDir "$Binary.exe") -Force

    # Add to PATH if not already there
    $UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
        Write-Host "Added $InstallDir to PATH."
    }

    Write-Host "Installed $Binary to $InstallDir\$Binary.exe"
    Write-Host ""
    Write-Host "Run '$Binary init' to configure your statusline."
    Write-Host "You may need to restart your terminal for PATH changes to take effect."
}
finally {
    Remove-Item -Path $TmpDir -Recurse -Force -ErrorAction SilentlyContinue
}
