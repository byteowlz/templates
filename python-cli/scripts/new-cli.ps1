param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$Name,
    [string]$Path,
    [string]$Package
)

function Show-Usage {
    @'
Usage: new-cli.ps1 -Name <name> [-Path <dir>] [-Package <module>]

Create a new CLI project by cloning the current template into DIR (defaults to <name>).
'@
}

if (-not $Name) {
    Show-Usage
    exit 1
}

if ($Name -notmatch '^[a-zA-Z][a-zA-Z0-9_-]*$') {
    Write-Error "Project name must start with a letter and contain only letters, numbers, '_' or '-'"
    exit 1
}

if (-not $Package -or $Package -eq '') {
    $Package = $Name.Replace('-', '_')
}

if ($Package -notmatch '^[a-zA-Z_][a-zA-Z0-9_]*$') {
    Write-Error "Package name must start with a letter or underscore and contain only letters, numbers, or '_'"
    exit 1
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$TemplateRoot = Split-Path -Parent $ScriptDir

if (-not $Path -or $Path -eq '') {
    $Path = Join-Path (Split-Path -Parent $TemplateRoot) $Name
}

if (Test-Path -Path $Path) {
    Write-Error "Destination already exists: $Path"
    exit 1
}

$ignore = @('.git', '.ruff_cache', '.mypy_cache', '.pytest_cache', '.venv', 'dist', 'build', '__pycache__')

Copy-Item -Path $TemplateRoot -Destination $Path -Recurse -Force -Exclude $ignore

Get-ChildItem -Path $Path -Recurse -Filter '*.sh' | ForEach-Object {
    $_.Attributes = $_.Attributes -band (-bnot [System.IO.FileAttributes]::ReadOnly)
}

$upper = $Name.ToUpper().Replace('-', '_')

$replacements = @{
    'python-cli' = $Name
    'PYTHON_CLI' = $upper
    'python_cli' = $Package
}

Get-ChildItem -Path $Path -Recurse -File | ForEach-Object {
    try {
        $content = Get-Content $_.FullName -Raw -ErrorAction Stop
    } catch {
        return
    }
    $updated = $content
    foreach ($pair in $replacements.GetEnumerator()) {
        $updated = $updated.Replace($pair.Key, $pair.Value)
    }
    if ($updated -ne $content) {
        Set-Content $_.FullName -Value $updated
    }
}

$packageDir = Join-Path $Path 'python_cli'
if (Test-Path $packageDir) {
    Rename-Item -Path $packageDir -NewName $Package -Force
}

Write-Output "Created CLI project at $Path"
Write-Output "Next steps:"
Write-Output "  1. Set-Location $Path"
Write-Output "  2. uv sync"
Write-Output "  3. uv run $Name -- --help"
