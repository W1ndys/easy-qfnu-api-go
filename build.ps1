# Easy QFNU API Go 构建脚本
# 用于交叉编译 Windows 和 Linux 平台的二进制文件

$ErrorActionPreference = "Stop"

$APP_NAME = "easy-qfnu-api-go"
$OUTPUT_DIR = "build"

# 创建输出目录
if (-not (Test-Path $OUTPUT_DIR)) {
    New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
}

Write-Host "开始构建 $APP_NAME ..." -ForegroundColor Cyan

# 构建 Linux amd64
Write-Host "正在构建 Linux (amd64) ..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$OUTPUT_DIR/${APP_NAME}-linux-amd64" .
Write-Host "√ Linux (amd64) 构建完成: $OUTPUT_DIR/${APP_NAME}-linux-amd64" -ForegroundColor Green

# 构建 Linux arm64
Write-Host "正在构建 Linux (arm64) ..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -o "$OUTPUT_DIR/${APP_NAME}-linux-arm64" .
Write-Host "√ Linux (arm64) 构建完成: $OUTPUT_DIR/${APP_NAME}-linux-arm64" -ForegroundColor Green

# 构建 Windows amd64
Write-Host "正在构建 Windows (amd64) ..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$OUTPUT_DIR/${APP_NAME}-windows-amd64.exe" .
Write-Host "√ Windows (amd64) 构建完成: $OUTPUT_DIR/${APP_NAME}-windows-amd64.exe" -ForegroundColor Green

# 构建 Windows arm64
Write-Host "正在构建 Windows (arm64) ..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "arm64"
go build -o "$OUTPUT_DIR/${APP_NAME}-windows-arm64.exe" .
Write-Host "√ Windows (arm64) 构建完成: $OUTPUT_DIR/${APP_NAME}-windows-arm64.exe" -ForegroundColor Green

# 清理环境变量
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "所有平台构建完成！输出目录: $OUTPUT_DIR/" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Get-ChildItem $OUTPUT_DIR | Format-Table Name, Length, LastWriteTime
