#!/usr/bin/env pwsh

[CmdletBinding()]
param(
    [Parameter(Position = 0, Mandatory = $true)]
    [ValidatePattern('^[a-zA-Z][a-zA-Z0-9_-]*$')]
    [string]$Name,

    [Parameter()]
    [string]$Path
)

function Write-Usage {
    Write-Host "Usage: new-cli.ps1 <name> [-Path <destination>]" -ForegroundColor Cyan
    Write-Host "Creates a new CLI project by cloning this template." -ForegroundColor Cyan
}

try {
    if ($PSBoundParameters.ContainsKey('Name') -eq $false) {
        Write-Usage
        throw 'Project name is required.'
    }

    $scriptDir = Split-Path -Parent $PSCommandPath
    $templateRoot = Split-Path -Parent $scriptDir

    if ([string]::IsNullOrWhiteSpace($Path)) {
        $parentDir = Split-Path -Parent $templateRoot
        $destination = Join-Path -Path $parentDir -ChildPath $Name
    } else {
        $destination = if ([System.IO.Path]::IsPathRooted($Path)) { $Path } else { Join-Path -Path (Get-Location) -ChildPath $Path }
    }

    if (Test-Path -LiteralPath $destination) {
        throw "Destination already exists: $destination"
    }

    New-Item -ItemType Directory -Path $destination | Out-Null

    $excluded = @('.git', 'target', '.DS_Store')

    function Copy-Template {
        param(
            [string]$Source,
            [string]$Dest
        )

        Get-ChildItem -LiteralPath $Source -Force | ForEach-Object {
            if ($excluded -contains $_.Name) {
                return
            }

            $targetPath = Join-Path -Path $Dest -ChildPath $_.Name

            if ($_.PSIsContainer) {
                if (-not (Test-Path -LiteralPath $targetPath)) {
                    New-Item -ItemType Directory -Path $targetPath | Out-Null
                }
                Copy-Template -Source $_.FullName -Dest $targetPath
            } else {
                Copy-Item -LiteralPath $_.FullName -Destination $targetPath -Force
            }
        }
    }

    Copy-Template -Source $templateRoot -Dest $destination

    $underscore = $Name -replace '-', '_'
    $upper = $Name.ToUpper() -replace '-', '_'

    $textExtensions = @(
        '.toml', '.md', '.rs', '.json', '.yml', '.yaml', '.sh', '.ps1',
        '.gitignore', '.txt', '.cfg', '.lock'
    )

    function Replace-InFile {
        param([string]$FilePath)

        try {
            $fullText = Get-Content -LiteralPath $FilePath -Raw -Encoding UTF8
        } catch {
            return
        }
        # .Replace is case-sensitive literal replacement (unlike -replace).
        $updated = $fullText
            .Replace('rust-cli', $Name)
            .Replace('rust_cli', $underscore)
            .Replace('RUST_CLI', $upper)
            .Replace('{{project_name}}', $Name)
            .Replace('your-binary-name', $Name)
        if ($updated -ne $fullText) {
            Set-Content -LiteralPath $FilePath -Value $updated -Encoding UTF8
        }
    }

    Get-ChildItem -LiteralPath $destination -Recurse -File -Force | ForEach-Object {
        # Keep the scaffolding/tooling scripts' own replacement tables intact.
        if ($_.Name -in @('new-cli.sh', 'new-cli.ps1', 'smoke-test.sh', 'drift-check.sh')) {
            return
        }
        if ($textExtensions -contains $_.Extension.ToLower()) {
            Replace-InFile -FilePath $_.FullName
        }
    }

    Write-Host "Created CLI project at $destination" -ForegroundColor Green
    exit 0
}
catch {
    Write-Error $_
    Write-Usage
    exit 1
}
