<#
.SYNOPSIS
    Go 项目自动化部署脚本 (Tar + SSH)

.DESCRIPTION
    1. 自动检测本地 SSH 密钥
    2. 检查并自动创建远程目标目录
    3. 使用 Git Bash 的 tar + ssh 管道流式上传文件
    4. 备份旧版本、重启 Supervisor 服务
    5. 健康检查，失败时自动回滚

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

# SSH 私钥路径 (留空则自动检测 ~/.ssh/id_rsa 或 id_ed25519)
$IdentityFile = ""

# 项目名称 (用于生成二进制文件名)
$ProjectName = "easy-qfnu-api-go"

# 目标操作系统
$TargetOS = "linux"

# 目标架构
$TargetArch = "amd64"

# Supervisor 服务名称
$SupervisorService = "easy-qfnu-api-go:easy-qfnu-api-go_00"

# 健康检查 URL
$HealthCheckUrl = "http://127.0.0.1:8141"

# 健康检查超时时间 (秒)
$HealthCheckTimeout = 30

# 健康检查重试次数
$HealthCheckRetries = 3

# 健康检查重试间隔 (秒)
$HealthCheckInterval = 5

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

# 辅助函数：执行远程命令
function Invoke-RemoteCommand {
    param([string]$Command)
    $cmd = $sshCmdPrefix + $Command
    & $cmd[0] $cmd[1..($cmd.Length - 1)]
    return $LASTEXITCODE
}

# 辅助函数：健康检查 (在远程服务器上执行)
function Test-ServiceHealth {
    param(
        [string]$Url,
        [int]$Timeout,
        [int]$Retries,
        [int]$Interval
    )

    for ($i = 1; $i -le $Retries; $i++) {
        Write-Host "[-] 健康检查 (第 $i/$Retries 次)..." -ForegroundColor Cyan

        # 在远程服务器上使用 curl 进行健康检查
        $healthCmd = "curl -s -o /dev/null -w '%{http_code}' --connect-timeout $Timeout '$Url'"
        $httpCode = & $sshCmdPrefix[0] $sshCmdPrefix[1..($sshCmdPrefix.Length - 1)] $healthCmd 2>$null

        if ($httpCode -eq "200") {
            Write-Host "[+] 健康检查通过! (HTTP $httpCode)" -ForegroundColor Green
            return $true
        }
        elseif ($httpCode) {
            Write-Host "[!] 健康检查返回非 200 状态码: $httpCode" -ForegroundColor Yellow
        }
        else {
            Write-Host "[!] 健康检查失败: 无法连接到服务" -ForegroundColor Yellow
        }

        if ($i -lt $Retries) {
            Write-Host "[-] 等待 $Interval 秒后重试..." -ForegroundColor DarkGray
            Start-Sleep -Seconds $Interval
        }
    }

    return $false
}

# 2. 交叉编译 (Windows -> Linux)
Write-Host "[-] 正在编译 $TargetOS ($TargetArch) 二进制文件..." -ForegroundColor Cyan
$BinaryName = "${ProjectName}-${TargetOS}-${TargetArch}"
$BackupBinaryName = "${BinaryName}.backup"

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
$exitCode = Invoke-RemoteCommand "mkdir -p $RemotePath"
if ($exitCode -ne 0) {
    Write-Error "无法创建远程目录，请检查连接或权限。"
}

# 4. 备份旧版本
Write-Host "[-] 正在备份旧版本..." -ForegroundColor Cyan
$backupCmd = "if [ -f '$RemotePath/$BinaryName' ]; then cp '$RemotePath/$BinaryName' '$RemotePath/$BackupBinaryName' && echo 'Backup created'; else echo 'No existing binary to backup'; fi"
$exitCode = Invoke-RemoteCommand $backupCmd
if ($exitCode -ne 0) {
    Write-Warning "备份旧版本失败，继续部署..."
}

# 5. 使用 Tar + SSH 上传二进制文件 (通过 Git Bash)
Write-Host "[-] 正在上传二进制文件..." -ForegroundColor Cyan

# 查找 Git Bash (优先使用 Git for Windows，避免 WSL)
$gitBashPaths = @(
    "$env:ProgramFiles\Git\bin\bash.exe",
    "${env:ProgramFiles(x86)}\Git\bin\bash.exe",
    "$env:LOCALAPPDATA\Programs\Git\bin\bash.exe",
    "$env:ProgramFiles\Git\usr\bin\bash.exe"
)

$bashExe = $null
foreach ($path in $gitBashPaths) {
    if (Test-Path $path) {
        $bashExe = $path
        break
    }
}

# 如果预设路径找不到，尝试从 PATH 中查找 Git 目录下的 bash
if (-not $bashExe) {
    $gitCmd = Get-Command git -ErrorAction SilentlyContinue
    if ($gitCmd) {
        # git.exe 通常在 Git\cmd 目录，bash 在 Git\bin 目录
        $gitDir = Split-Path (Split-Path $gitCmd.Source -Parent) -Parent
        $gitBash = Join-Path $gitDir "bin\bash.exe"
        if (Test-Path $gitBash) {
            $bashExe = $gitBash
        }
    }
}

if (-not $bashExe) {
    Write-Error "未找到 Git Bash，请确保已安装 Git for Windows。"
}

Write-Host "[-] 使用 Git Bash: $bashExe" -ForegroundColor DarkGray

# 将 Windows 路径转换为 Unix 风格路径 (用于 Git Bash)
$unixIdentityFile = $IdentityFile -replace '\\', '/' -replace '^([A-Za-z]):', '/$1'

# 构造 bash 命令：
# 1. 本地 tar 打包二进制文件
# 2. SSH 传输
# 3. 远程 tar 解压
# 4. 远程 chmod +x 赋予执行权限
$bashCmd = "tar -c $BinaryName | ssh -i '$unixIdentityFile' -p $Port -o StrictHostKeyChecking=no $User@$Server 'tar -x -C $RemotePath && chmod +x $RemotePath/$BinaryName'"

Write-Host "Executing: Upload..." -ForegroundColor DarkGray
& $bashExe -c $bashCmd

if ($LASTEXITCODE -eq 0) {
    Write-Host "[+] 文件上传成功!" -ForegroundColor Green

    # 删除本地构建产物
    if (Test-Path $BinaryName) {
        Remove-Item $BinaryName -Force
        Write-Host "[-] 已清理本地构建产物: $BinaryName" -ForegroundColor DarkGray
    }
}
else {
    Write-Error "文件上传失败。"
}

# 6. 重启 Supervisor 服务
Write-Host "[-] 正在重启 Supervisor 服务: $SupervisorService" -ForegroundColor Cyan
$exitCode = Invoke-RemoteCommand "supervisorctl restart $SupervisorService"
if ($exitCode -ne 0) {
    Write-Warning "Supervisor 重启命令返回非零状态码"
}

# 等待服务启动
Write-Host "[-] 等待服务启动..." -ForegroundColor Cyan
Start-Sleep -Seconds 3

# 7. 健康检查
Write-Host "[-] 开始健康检查: $HealthCheckUrl" -ForegroundColor Cyan
$isHealthy = Test-ServiceHealth -Url $HealthCheckUrl -Timeout $HealthCheckTimeout -Retries $HealthCheckRetries -Interval $HealthCheckInterval

if ($isHealthy) {
    Write-Host "[+] 部署成功! 服务运行正常。" -ForegroundColor Green

    # 删除备份文件
    Write-Host "[-] 清理备份文件..." -ForegroundColor DarkGray
    Invoke-RemoteCommand "rm -f '$RemotePath/$BackupBinaryName'" | Out-Null
}
else {
    Write-Host "[!] 健康检查失败! 正在回滚..." -ForegroundColor Red

    # 检查备份是否存在
    $checkBackup = "test -f '$RemotePath/$BackupBinaryName' && echo 'EXISTS' || echo 'NOT_EXISTS'"
    $backupExists = & $sshCmdPrefix[0] $sshCmdPrefix[1..($sshCmdPrefix.Length - 1)] $checkBackup

    if ($backupExists -eq "EXISTS") {
        Write-Host "[-] 正在恢复旧版本..." -ForegroundColor Yellow

        # 恢复备份
        $restoreCmd = "cp '$RemotePath/$BackupBinaryName' '$RemotePath/$BinaryName' && chmod +x '$RemotePath/$BinaryName'"
        $exitCode = Invoke-RemoteCommand $restoreCmd

        if ($exitCode -eq 0) {
            Write-Host "[-] 旧版本已恢复，正在重启服务..." -ForegroundColor Yellow

            # 重启服务
            $exitCode = Invoke-RemoteCommand "supervisorctl restart $SupervisorService"

            # 等待服务启动
            Start-Sleep -Seconds 3

            # 再次健康检查
            $isRollbackHealthy = Test-ServiceHealth -Url $HealthCheckUrl -Timeout $HealthCheckTimeout -Retries 2 -Interval 3

            if ($isRollbackHealthy) {
                Write-Host "[+] 回滚成功! 服务已恢复到旧版本。" -ForegroundColor Yellow
            }
            else {
                Write-Host "[!] 回滚后健康检查仍然失败，请手动检查服务状态!" -ForegroundColor Red
            }
        }
        else {
            Write-Host "[!] 恢复旧版本失败，请手动检查!" -ForegroundColor Red
        }

        # 获取 Supervisor 日志
        Write-Host "`n[-] Supervisor 服务日志:" -ForegroundColor Cyan
        Invoke-RemoteCommand "supervisorctl tail $SupervisorService" | Out-Null
    }
    else {
        Write-Host "[!] 没有找到备份文件，无法回滚! 请手动检查服务状态。" -ForegroundColor Red

        # 获取 Supervisor 日志
        Write-Host "`n[-] Supervisor 服务日志:" -ForegroundColor Cyan
        Invoke-RemoteCommand "supervisorctl tail $SupervisorService" | Out-Null
    }

    Write-Error "部署失败，已尝试回滚。"
}

Write-Host "`n部署流程结束。" -ForegroundColor Cyan
