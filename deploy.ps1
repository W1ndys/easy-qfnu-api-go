<#
.SYNOPSIS
    Go 项目自动化部署脚本 (Tar + SSH)

.DESCRIPTION
    1. 自动检测本地 SSH 密钥
    2. 检查并自动创建远程目标目录
    3. 使用 tar + ssh 管道流式上传文件
    4. 远程执行重启命令

.EXAMPLE
    .\deploy.ps1
#>

# ============================================================
#                     配置区域 (请在此修改)
# ============================================================

# 服务器地址 (IP 或域名)
$Server = "my.server"

# SSH 端口
$Port = 22

# 登录用户名
$User = "root"

# 远程部署路径
$RemotePath = "/root/easy-qfnu-api-go"

# 部署完成后执行的远程命令 (留空则不执行)
$RestartCmd = "echo 'Deploy finished, no restart command specified.'"

# 本地项目路径
$LocalPath = "."

# SSH 私钥路径 (留空则自动检测 ~/.ssh/id_rsa 或 id_ed25519)
$IdentityFile = ""

# 项目名称 (用于生成二进制文件名)
$ProjectName = "easy-qfnu-api-go"

# 目标操作系统
$TargetOS = "linux"

# 目标架构
$TargetArch = "amd64"

# ============================================================
#                     配置区域结束
# ============================================================

$ErrorActionPreference = "Stop"

# 1. 自动检测 SSH 密钥
if ([string]::IsNullOrEmpty($IdentityFile)) {
    $sshDir = "$env:USERPROFILE\.ssh"
    $possibleKeys = @("id_rsa", "id_ed25519")

    foreach ($keyName in $possibleKeys) {
        $path = Join-Path $sshDir $keyName
        if (Test-Path $path) {
            $IdentityFile = $path
            Write-Host "[-] 自动检测到 SSH 密钥: $IdentityFile" -ForegroundColor Cyan
            break
        }
    }
}

if (-not (Test-Path $IdentityFile)) {
    Write-Error "未找到 SSH 密钥，请在配置区域指定 IdentityFile。"
}

# 构建基础 SSH 命令前缀
$sshCmdPrefix = @("ssh", "-i", "$IdentityFile", "-p", "$Port", "-o", "StrictHostKeyChecking=no", "$User@$Server")

# 2. 交叉编译 (Windows -> Linux)
Write-Host "[-] 正在编译 $TargetOS ($TargetArch) 二进制文件..." -ForegroundColor Cyan
$BinaryName = "${ProjectName}-${TargetOS}-${TargetArch}"

# 保存旧的环境变量
$OriginalGOOS = $env:GOOS
$OriginalGOARCH = $env:GOARCH

try {
    $env:CGO_ENABLED = "0"
    $env:GOOS = $TargetOS
    $env:GOARCH = $TargetArch

    go build -ldflags "-s -w" -o $BinaryName .

    if ($LASTEXITCODE -ne 0) {
        Write-Error "编译失败，请检查 Go 环境或代码错误。"
    }
    Write-Host "[-] 编译成功: $BinaryName" -ForegroundColor Green
}
finally {
    # 恢复环境变量
    $env:GOOS = $OriginalGOOS
    $env:GOARCH = $OriginalGOARCH
}

# 3. 检查并修复远程路径 (mkdir -p)
Write-Host "[-] 正在检查/创建远程目录: $RemotePath" -ForegroundColor Cyan
$mkdirCmd = $sshCmdPrefix + "mkdir -p $RemotePath"
& $mkdirCmd[0] $mkdirCmd[1..($mkdirCmd.Length - 1)]
if ($LASTEXITCODE -ne 0) {
    Write-Error "无法创建远程目录，请检查连接或权限。"
}

# 4. 使用 Tar + SSH 上传二进制文件
Write-Host "[-] 正在上传二进制文件..." -ForegroundColor Cyan

# 构造上传命令：
# 1. 本地 tar 打包二进制文件
# 2. SSH 传输
# 3. 远程 tar 解压
# 4. 远程 chmod +x 赋予执行权限
$uploadCmdString = "tar -c $BinaryName | ssh -i `"$IdentityFile`" -p $Port -o StrictHostKeyChecking=no $User@$Server `"tar -x -C $RemotePath && chmod +x $RemotePath/$BinaryName`""

Write-Host "Executing: Upload..." -ForegroundColor DarkGray
Invoke-Expression $uploadCmdString

if ($LASTEXITCODE -eq 0) {
    Write-Host "[+] 文件上传成功!" -ForegroundColor Green
}
else {
    Write-Error "文件上传失败。"
}

# 5. 执行远程重启命令
if (-not [string]::IsNullOrEmpty($RestartCmd)) {
    Write-Host "[-] 正在执行远程命令: $RestartCmd" -ForegroundColor Cyan
    $remoteExec = $sshCmdPrefix + $RestartCmd
    & $remoteExec[0] $remoteExec[1..($remoteExec.Length - 1)]

    if ($LASTEXITCODE -eq 0) {
        Write-Host "[+] 远程命令执行成功!" -ForegroundColor Green
    }
    else {
        Write-Warning "远程命令执行返回了非零状态码。"
    }
}

Write-Host "`n部署流程结束。" -ForegroundColor Cyan
