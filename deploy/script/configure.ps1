#Requires -Version 5.1
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ScriptDir      = Split-Path -Parent $MyInvocation.MyCommand.Definition
$BasePath       = (Resolve-Path "$ScriptDir\..\.." ).Path
$DeployDir      = "$BasePath\deploy"
$EnvFile        = "$DeployDir\.env"
$EnvExampleFile = "$DeployDir\.env.example"

if (-not (Test-Path $EnvExampleFile)) {
    Write-Error "Missing required file: $EnvExampleFile"
    exit 1
}

if (-not (Test-Path $EnvFile)) {
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

# Write or update a single key in the .env file
function Update-EnvKey([string]$Key, [string]$Value) {
    $lines   = Get-Content $EnvFile
    $updated = $false
    $result  = foreach ($line in $lines) {
        if ($line -match "^$Key=") { "$Key=$Value"; $updated = $true }
        else                       { $line }
    }
    if (-not $updated) { $result += "$Key=$Value" }
    $result | Set-Content $EnvFile
}

# Prompt the user for a value, showing the current default in brackets
function Prompt-Value([string]$Key, [string]$Label, [string]$Fallback) {
    $current = if ($script:cfg.ContainsKey($Key) -and $script:cfg[$Key] -ne '') {
        $script:cfg[$Key]
    } else { $Fallback }

    $input = Read-Host "$Label [$current]"
    if ([string]::IsNullOrEmpty($input)) { $input = $current }

    $script:cfg[$Key] = $input
    Update-EnvKey $Key $input
}

# Merge example defaults then overlay with actual .env values
$script:cfg = Read-EnvFile $EnvExampleFile
foreach ($kv in (Read-EnvFile $EnvFile).GetEnumerator()) {
    $script:cfg[$kv.Key] = $kv.Value
}

Write-Host "== Yoresee Deploy Config Init =="
Write-Host "Press Enter to keep the default value in brackets."

# Host exposed ports
Prompt-Value 'NGINX_HTTP_PORT'            'Host port for Nginx HTTP'               '8080'
Prompt-Value 'NGINX_HTTPS_PORT'           'Host port for Nginx HTTPS'              '8443'
Prompt-Value 'BACKEND_GRPC_HOST_PORT'     'Host port for backend gRPC'             '9090'
Prompt-Value 'BACKEND_DEBUG_HOST_PORT'    'Host port for backend debug (dev only)' '2345'
Prompt-Value 'POSTGRES_HOST_PORT'         'Host port for Postgres (dev only)'      '5432'
Prompt-Value 'REDIS_HOST_PORT'            'Host port for Redis (dev only)'         '6379'
Prompt-Value 'RABBITMQ_AMQP_HOST_PORT'    'Host port for RabbitMQ AMQP (dev only)' '5672'
Prompt-Value 'CONSUL_HOST_PORT'           'Host port for Consul UI/API'            '8500'
Prompt-Value 'ELASTICSEARCH_HOST_PORT'    'Host port for Elasticsearch (dev only)' '9200'

# Credentials and app secrets
Prompt-Value 'POSTGRES_USER'              'Postgres user'                          'root'
Prompt-Value 'POSTGRES_PASSWORD'          'Postgres password'                      'your_password'
Prompt-Value 'POSTGRES_DB'               'Postgres database'                      'yoresee_doc_db'
Prompt-Value 'REDIS_PASSWORD'             'Redis password'                         'your_redis_password'
Prompt-Value 'RABBITMQ_DEFAULT_USER'      'RabbitMQ user'                          'guest'
Prompt-Value 'RABBITMQ_DEFAULT_PASS'      'RabbitMQ password'                      'guest'
Prompt-Value 'MINIO_ROOT_USER'            'MinIO root user'                        'minioadmin'
Prompt-Value 'MINIO_ROOT_PASSWORD'        'MinIO root password'                    'minioadmin'
Prompt-Value 'CONSUL_ROOT_TOKEN'          'Consul root token'                      'yoresee_doc_root_token'
Prompt-Value 'BACKEND_INTERNAL_RPC_KEY'   'Internal RPC key'                       'yoresee_doc_internal_key'
Prompt-Value 'JWT_SECRET'                 'JWT secret'                             'yoresee_doc_jwt_secret_key'

# App behaviour
$defaultApiUrl = "http://localhost:$($script:cfg['NGINX_HTTP_PORT'])"
Prompt-Value 'VITE_API_BASE_URL'          'Frontend API base URL'                  $defaultApiUrl
Prompt-Value 'VITE_GRPC_WEB_ENDPOINT'     'Frontend gRPC-web endpoint'             '/grpc'
Prompt-Value 'DIRTY_DOC_NOTIFY_THRESHOLD' 'Dirty doc notify threshold'             '5'

# Derived value — keep MinIO redirect bound to API base URL
$minioRedirect = "$($script:cfg['VITE_API_BASE_URL'].TrimEnd('/'))/minio"
Update-EnvKey 'MINIO_BROWSER_REDIRECT_URL' $minioRedirect
Write-Host "MinIO browser redirect URL (derived): $minioRedirect"

Prompt-Value 'ELASTICSEARCH_ENABLED'      'Enable Elasticsearch'                   'true'
Prompt-Value 'ELASTICSEARCH_INDEX_PREFIX' 'Elasticsearch index prefix'             'yoresee_doc'
Prompt-Value 'ELASTICSEARCH_USERNAME'     'Elasticsearch username (optional)'      ''
Prompt-Value 'ELASTICSEARCH_PASSWORD'     'Elasticsearch password (optional)'      ''

Write-Host "Written: $EnvFile"

& powershell -ExecutionPolicy Bypass -File "$ScriptDir\prepare.ps1"

Write-Host "Initialization completed."
Write-Host "You can start services with: powershell -ExecutionPolicy Bypass -File deploy/script/start.ps1 dev up"
