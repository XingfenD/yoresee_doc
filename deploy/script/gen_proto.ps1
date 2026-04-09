#Requires -Version 5.1
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Definition
$RootDir    = (Resolve-Path "$ScriptDir\..\.." ).Path
$ProtoDir   = "$RootDir\proto"
$GoOutDir   = "$RootDir\backend\pkg\gen"
$GoOutDir2  = "$RootDir\collab-go\pkg\gen"
$FrontendOutDir = "$RootDir\frontend\src\gen"
$CollabOutDir   = "$RootDir\collab\src\gen"
$ProtoFile  = "$ProtoDir\yoresee_doc\v1\yoresee_doc.proto"

# Add Go bin to PATH so protoc-gen-go / protoc-gen-go-grpc are found
$GoPath = (go env GOPATH)
# Add grpc-tools bin to PATH so protoc is found
$NpmGlobal = (npm root -g)
$env:PATH = "$GoPath\bin;$NpmGlobal\grpc-tools\bin;$env:PATH"

# --- Go code generation (backend + collab-go) ---

New-Item -ItemType Directory -Force -Path $GoOutDir | Out-Null
protoc -I $ProtoDir `
  --go_out=$GoOutDir --go_opt=paths=source_relative `
  --go-grpc_out=$GoOutDir --go-grpc_opt=paths=source_relative `
  $ProtoFile

New-Item -ItemType Directory -Force -Path $GoOutDir2 | Out-Null
protoc -I $ProtoDir `
  --go_out=$GoOutDir2 --go_opt=paths=source_relative `
  --go-grpc_out=$GoOutDir2 --go-grpc_opt=paths=source_relative `
  $ProtoFile

# --- Frontend ES / Connect-ES code generation ---

$EsBin      = "$RootDir\frontend\node_modules\.bin\protoc-gen-es.cmd"
$ConnectBin = "$RootDir\frontend\node_modules\.bin\protoc-gen-connect-es.cmd"

if ((Test-Path $EsBin) -and (Test-Path $ConnectBin)) {
  New-Item -ItemType Directory -Force -Path $FrontendOutDir | Out-Null
  protoc -I $ProtoDir `
    "--plugin=protoc-gen-es=$EsBin" `
    "--plugin=protoc-gen-connect-es=$ConnectBin" `
    "--es_out=$FrontendOutDir" --es_opt=target=js,import_extension=.js `
    "--connect-es_out=$FrontendOutDir" --connect-es_opt=target=js,import_extension=.js `
    $ProtoFile
} else {
  Write-Warning "protoc-gen-es or protoc-gen-connect-es not found; skipping frontend connect-es code generation"
}

# --- Collab Node.js gRPC code generation ---
# protoc requires the actual .exe plugin, not the .cmd wrapper.

$NpmGlobal  = (npm root -g)
$GrpcPlugin = "$NpmGlobal\grpc-tools\bin\grpc_node_plugin.exe"

if (Test-Path $GrpcPlugin) {
  New-Item -ItemType Directory -Force -Path $CollabOutDir | Out-Null
  protoc -I $ProtoDir `
    "--plugin=protoc-gen-grpc=$GrpcPlugin" `
    "--grpc_out=grpc_js:$CollabOutDir" `
    "--js_out=import_style=commonjs,binary:$CollabOutDir" `
    $ProtoFile
  Write-Host "Generated gRPC code for collab at $CollabOutDir"
} else {
  Write-Error "grpc_tools_node_protoc_plugin not found at $GrpcPlugin`nInstall with: npm install -g grpc-tools"
  exit 1
}
