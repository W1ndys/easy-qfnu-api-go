# Tailwind CSS Standalone CLI 构建脚本
# 无需 Node.js/npm

param(
    [Parameter(Position=0)]
    [ValidateSet("build", "dev", "download", "clean", "help")]
    [string]$Command = "help"
)

$TAILWIND_VERSION = "4.1.18"
$TAILWIND_EXE = "tailwindcss.exe"
$TAILWIND_URL = "https://github.com/tailwindlabs/tailwindcss/releases/download/v$TAILWIND_VERSION/tailwindcss-windows-x64.exe"

function Download-Tailwind {
    if (Test-Path $TAILWIND_EXE) {
        Write-Host "$TAILWIND_EXE 已存在" -ForegroundColor Green
    } else {
        Write-Host "正在下载 Tailwind CSS CLI v$TAILWIND_VERSION..." -ForegroundColor Cyan
        Invoke-WebRequest -Uri $TAILWIND_URL -OutFile $TAILWIND_EXE
        Write-Host "下载完成: $TAILWIND_EXE" -ForegroundColor Green
    }
}

function Build-Css {
    Download-Tailwind
    Write-Host "正在构建 Tailwind CSS..." -ForegroundColor Cyan
    & ".\$TAILWIND_EXE" -i .\static\css\input.css -o .\static\css\tailwind.css --minify
    Write-Host "构建完成! 输出: static\css\tailwind.css" -ForegroundColor Green
}

function Start-Dev {
    Download-Tailwind
    Write-Host "启动 Tailwind CSS 监听模式..." -ForegroundColor Cyan
    & ".\$TAILWIND_EXE" -i .\static\css\input.css -o .\static\css\tailwind.css --watch
}

function Clean-Files {
    if (Test-Path $TAILWIND_EXE) {
        Remove-Item $TAILWIND_EXE -Force
        Write-Host "已删除 $TAILWIND_EXE" -ForegroundColor Yellow
    }
    if (Test-Path "static\css\tailwind.css") {
        Remove-Item "static\css\tailwind.css" -Force
        Write-Host "已删除 static\css\tailwind.css" -ForegroundColor Yellow
    }
    Write-Host "清理完成!" -ForegroundColor Green
}

function Show-Help {
    Write-Host ""
    Write-Host "Tailwind CSS Standalone CLI 构建工具" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "用法: .\build.ps1 <命令>" -ForegroundColor White
    Write-Host ""
    Write-Host "命令:" -ForegroundColor White
    Write-Host "  download  - 下载 Tailwind CLI"
    Write-Host "  build     - 生产构建 (压缩)"
    Write-Host "  dev       - 开发模式 (监听变化)"
    Write-Host "  clean     - 清理 CLI 和输出文件"
    Write-Host "  help      - 显示此帮助"
    Write-Host ""
    Write-Host "无需 Node.js/npm!" -ForegroundColor Green
    Write-Host ""
}

switch ($Command) {
    "download" { Download-Tailwind }
    "build"    { Build-Css }
    "dev"      { Start-Dev }
    "clean"    { Clean-Files }
    "help"     { Show-Help }
    default    { Show-Help }
}
