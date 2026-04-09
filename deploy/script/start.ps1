#Requires -Version 5.1
param(
    [Parameter(Mandatory)][ValidateSet('dev','release')][string]$Environment,
    [Parameter(Mandatory)][ValidateSet('rebuild','clear','restart','up','build')][string]$Action
)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$BasePath  = (Resolve-Path "$ScriptDir\..\.." ).Path
$DeployDir = "$BasePath\deploy"

Write-Host "Environment: $Environment"
Write-Host "Action:      $Action"
Write-Host "Base path:   $BasePath"

$ComposeFile = if ($Environment -eq 'dev') { 'docker-compose.dev.yml' } else { 'docker-compose.yml' }
Write-Host "Compose file: $ComposeFile"

if ($Environment -eq 'dev') {
    Write-Host "Generating gRPC code for dev..."
    & powershell -ExecutionPolicy Bypass -File "$ScriptDir\gen_proto.ps1"
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}

Push-Location $DeployDir
try {
    switch ($Action) {
        'rebuild' {
            Write-Host "Rebuilding and starting services..."
            docker compose -f $ComposeFile down -v
            docker compose -f $ComposeFile up -d --build
        }
        'clear' {
            Write-Host "Clearing containers and volumes..."
            docker compose -f $ComposeFile down -v
        }
        'restart' {
            Write-Host "Restarting services..."
            docker compose -f $ComposeFile down
            docker compose -f $ComposeFile up -d
        }
        'up' {
            Write-Host "Starting services..."
            docker compose -f $ComposeFile up -d
        }
        'build' {
            Write-Host "Building images..."
            docker compose -f $ComposeFile build
        }
    }
} finally {
    Pop-Location
}

if ($LASTEXITCODE -ne 0) {
    Write-Error "Operation failed."
    exit $LASTEXITCODE
}

Write-Host "Operation completed successfully."
