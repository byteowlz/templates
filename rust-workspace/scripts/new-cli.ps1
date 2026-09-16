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
    Write-Host "Creates a new workspace project by cloning this template." -ForegroundColor Cyan
    Write-Host "Renames all crates from rust-* to <name>-* pattern." -ForegroundColor Cyan
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

    # Replacement values.
    $underscore = $Name -replace '-', '_'
    $upper = $Name.ToUpper() -replace '-', '_'

    # Exact names only: `rust-magic-linter` (lint preset) and `dtolnay/rust-toolchain`
    # (CI action) must stay untouched, so do NOT do a generic `rust-` rewrite.
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
        # .Replace is case-sensitive literal replacement.
        $updated = $fullText
            .Replace('rust-workspace', $Name)
            .Replace('RUST_WORKSPACE', $upper)
            .Replace('{{project_name}}-core', "$Name-core")
            .Replace('{{project_name}}-cli', "$Name-cli")
            .Replace('{{project_name}}-tui', "$Name-tui")
            .Replace('{{project_name}}-mcp', "$Name-mcp")
            .Replace('{{project_name}}-api', "$Name-api")
            .Replace('{{project_name}}_core', "$underscore`_core")
            .Replace('{{project_name}}_cli', "$underscore`_cli")
            .Replace('{{project_name}}_tui', "$underscore`_tui")
            .Replace('{{project_name}}_mcp', "$underscore`_mcp")
            .Replace('{{project_name}}_api', "$underscore`_api")
            .Replace('{{project_name}}', $Name)
            .Replace('your-binary-name', $Name)
        if ($updated -ne $fullText) {
            Set-Content -LiteralPath $FilePath -Value $updated -Encoding UTF8
        }
    }

    # Update file contents across the whole tree, keeping the scaffolding
    # scripts' own replacement table intact.
    Get-ChildItem -LiteralPath $destination -Recurse -File -Force | ForEach-Object {
        if ($_.Name -in @('new-cli.sh', 'new-cli.ps1', 'smoke-test.sh', 'drift-check.sh', 'privilege-boundary-test.sh')) {
            return
        }
        if ($textExtensions -contains $_.Extension.ToLower()) {
            Replace-InFile -FilePath $_.FullName
        }
    }

    # Rename crate directories ({{project_name}}-core -> <name>-core, etc.; older
    # templates carried rust-* -> <name>-core).
    $cratesDir = Join-Path -Path $destination -ChildPath 'crates'
    if (Test-Path -LiteralPath $cratesDir) {
        Get-ChildItem -LiteralPath $cratesDir -Directory | ForEach-Object {
            $n = $_.Name
            if ($n.Contains('{{project_name}}')) {
                $newName = $n.Replace('{{project_name}}', $Name)
            } elseif ($n.StartsWith('rust-')) {
                $newName = "$Name-" + $n.Substring('rust-'.Length)
            } else {
                return
            }
            Rename-Item -LiteralPath $_.FullName -NewName $newName
        }
    }

    Write-Host "Created workspace project at $destination" -ForegroundColor Green
    Write-Host "Crates renamed to: $Name-core, $Name-cli, $Name-tui, $Name-mcp, $Name-api" -ForegroundColor Green
    exit 0
}
catch {
    Write-Error $_
    Write-Usage
    exit 1
}
