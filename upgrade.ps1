#Requires -Version 5.1

<#
.SYNOPSIS
    Recursively upgrade all Go module dependencies in the monorepo.

.DESCRIPTION
    Traverses every directory containing a go.mod file and runs:
      go get -u ./...
      go mod tidy
#>

$ErrorActionPreference = "Stop"

# Resolve root directory from script location.
$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Write-Status {
    param(
        [string]$Message,
        [ConsoleColor]$Color = [ConsoleColor]::White
    )
    Write-Host $Message -ForegroundColor $Color
}

function Upgrade-Module {
    param([string]$ModDir)

    # Relative path for display.
    $relPath = $ModDir.Substring($RootDir.Length).TrimStart('\', '/')
    if (-not $relPath) { $relPath = "(root)" }

    Write-Status "[upgrade] $relPath" ([ConsoleColor]::Cyan)

    # go get -u ./...
    Push-Location $ModDir
    $prevPref = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    try {
        $output = & go get -u ./... 2>&1 | Out-String
        $getCode = $LASTEXITCODE
        if ($getCode -ne 0) {
            $ErrorActionPreference = $prevPref
            Write-Status "  X go get -u failed:" ([ConsoleColor]::Red)
            Write-Host $output
            return @{ Success = $false; FailedAt = "go get" }
        }

        # go mod tidy
        $output = & go mod tidy 2>&1 | Out-String
        $tidyCode = $LASTEXITCODE
        if ($tidyCode -ne 0) {
            $ErrorActionPreference = $prevPref
            Write-Status "  X go mod tidy failed:" ([ConsoleColor]::Red)
            Write-Host $output
            return @{ Success = $false; FailedAt = "go mod tidy" }
        }
    }
    finally {
        $ErrorActionPreference = $prevPref
        Pop-Location
    }

    Write-Status "  V done" ([ConsoleColor]::Green)
    return @{ Success = $true }
}

# Find all go.mod files, excluding .git and .idea directories.
$modFiles = Get-ChildItem -Path $RootDir -Filter "go.mod" -Recurse -File |
    Where-Object { $_.FullName -notmatch '[\\/]\.git[\\/]' -and $_.FullName -notmatch '[\\/]\.idea[\\/]' }

$Total = $modFiles.Count
$Success = 0
$Failed = @()

foreach ($modFile in $modFiles) {
    $modDir = Split-Path -Parent $modFile.FullName
    $result = Upgrade-Module -ModDir $modDir
    if ($result.Success) {
        $Success++
    } else {
        $relPath = $modDir.Substring($RootDir.Length).TrimStart('\', '/')
        if (-not $relPath) { $relPath = "(root)" }
        $Failed += "$relPath ($($result.FailedAt))"
    }
}

# Summary.
Write-Host ""
Write-Status "=== Summary ===" ([ConsoleColor]::Cyan)
Write-Host "Total modules: $Total"
Write-Status "Succeeded:     $Success" ([ConsoleColor]::Green)

if ($Failed.Count -gt 0) {
    Write-Status "Failed:        $($Failed.Count)" ([ConsoleColor]::Red)
    Write-Host ""
    Write-Status "Failed modules:" ([ConsoleColor]::Red)
    foreach ($f in $Failed) {
        Write-Host "  X $f"
    }
    exit 1
} else {
    Write-Status "All modules upgraded successfully!" ([ConsoleColor]::Green)
}