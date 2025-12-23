param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$Name,
    [Parameter(Position = 1)]
    [string]$Path,
    [string]$Module
)

function Show-Usage {
    @'
Usage: new-cli.ps1 -Name <name> [-Path <dir>] [-Module <module>]

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

if (-not $Module) {
    $Module = $Name
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$TemplateRoot = Split-Path -Parent $ScriptDir

if (-not $Path -or $Path -eq '') {
    $Path = Join-Path (Split-Path -Parent $TemplateRoot) $Name
}

$Destination = Resolve-Path -Path $Path -ErrorAction SilentlyContinue
if ($Destination) {
    Write-Error "Destination already exists: $Destination"
    exit 1
}

$Destination = $Path

$ignore = @('.git', 'dist', 'target', '.DS_Store', '__pycache__')

Copy-Item -Path $TemplateRoot -Destination $Destination -Recurse -Force -Container -Exclude $ignore

Get-ChildItem -Path $Destination -Recurse -Filter 'new-cli.sh' | ForEach-Object {
    $_.Attributes = $_.Attributes -bor [System.IO.FileAttributes]::Normal
}

$upper = $Name.ToUpper().Replace('-', '_')

$replacements = @{
    'go-cli' = $Name
    'GO_CLI' = $upper
    'gitlab.cc-asp.fraunhofer.de/templates/go-cli' = $Module
}

Get-ChildItem -Path $Destination -Recurse -File | ForEach-Object {
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

Write-Output "Created CLI project at $Destination"
Write-Output "Next steps:"
Write-Output "  1. Set-Location $Destination"
Write-Output "  2. go mod tidy"
Write-Output "  3. go run . -- --help"
