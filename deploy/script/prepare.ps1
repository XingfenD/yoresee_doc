#Requires -Version 5.1
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ScriptDir      = Split-Path -Parent $MyInvocation.MyCommand.Definition
$BasePath       = (Resolve-Path "$ScriptDir\..\.." ).Path
$DeployDir      = "$BasePath\deploy"
$EnvExampleFile = "$DeployDir\.env.example"
$EnvFile        = "$DeployDir\.env"

if (-not (Test-Path $EnvExampleFile)) {
    Write-Error "Missing required file: $EnvExampleFile"
    exit 1
}

if (-not (Test-Path $EnvFile)) {
    Write-Host "deploy\.env not found, creating from .env.example"
    Copy-Item $EnvExampleFile $EnvFile
}

# Load key=value pairs from a .env file into a hashtable (last value wins)
function Read-EnvFile([string]$Path) {
    $map = @{}
    foreach ($line in (Get-Content $Path)) {
        if ($line -match '^\s*#' -or $line.Trim() -eq '') { continue }
        $idx = $line.IndexOf('=')
        if ($idx -lt 1) { continue }
        $map[$line.Substring(0, $idx).Trim()] = $line.Substring($idx + 1)
    }
    return $map
}

# Expand ${VAR} placeholders in a template string using the provided env hashtable
function Expand-Template([string]$Text, [hashtable]$Env) {
    return [regex]::Replace($Text, '\$\{([A-Z0-9_]+)\}', {
        param($m)
        $key = $m.Groups[1].Value
        if ($Env.ContainsKey($key)) { $Env[$key] } else { '' }
    })
}

# Render a template file to an output path, substituting ${VAR} with env values
function Render-Template([string]$Template, [string]$Output, [hashtable]$Env) {
    if (-not (Test-Path $Template)) {
        Write-Error "Template not found: $Template"
        return
    }

    $outDir = Split-Path -Parent $Output
    if (-not (Test-Path $outDir)) { New-Item -ItemType Directory -Force -Path $outDir | Out-Null }

    # If a directory exists at the output path, back it up
    if (Test-Path $Output -PathType Container) {
        $backup = "$Output.backup.$([DateTimeOffset]::UtcNow.ToUnixTimeSeconds())"
        Move-Item $Output $backup
        Write-Host "Detected directory at $Output, moved to $backup"
    }

    $content = Get-Content $Template -Raw
    $rendered = Expand-Template $content $Env
    [System.IO.File]::WriteAllText($Output, $rendered)
}

# Merge example defaults then overlay with actual .env values
$env_map = Read-EnvFile $EnvExampleFile
foreach ($kv in (Read-EnvFile $EnvFile).GetEnumerator()) {
    $env_map[$kv.Key] = $kv.Value
}

$templateMappings = @(
    @{ Template = "$BasePath\backend\config.toml.tmpl";                    Output = "$BasePath\backend\config.toml" }
    @{ Template = "$BasePath\frontend\nginx.conf.tmpl";                    Output = "$BasePath\frontend\nginx.conf" }
    @{ Template = "$DeployDir\nginx\nginx.conf.tmpl";                      Output = "$DeployDir\nginx\nginx.conf" }
    @{ Template = "$DeployDir\nginx\conf.d\default.conf.tmpl";             Output = "$DeployDir\nginx\conf.d\default.conf" }
    @{ Template = "$DeployDir\redis\redis.conf.tmpl";                      Output = "$DeployDir\redis\redis.conf" }
    @{ Template = "$DeployDir\rabbitmq\rabbitmq.conf.tmpl";                Output = "$DeployDir\rabbitmq\rabbitmq.conf" }
)

Write-Host "Rendering configuration files from templates..."
foreach ($m in $templateMappings) {
    Render-Template $m.Template $m.Output $env_map
    $rel = $m.Output.Replace("$BasePath\", '')
    Write-Host "  - $rel"
}

Write-Host "Configuration preparation completed!"
