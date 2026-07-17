<#
.SYNOPSIS
批量向 IndexNow 提交新增、更新或已删除的 URL。

.DESCRIPTION
支持通过 -Url 或 -UrlFile 提供 URL，自动去重并按 IndexNow 的单批上限分批提交。
已删除的 URL 应先在生产环境返回 404 或 410，再通过本脚本提交。

.EXAMPLE
.\scripts\submit-indexnow.ps1 -UrlFile .\scripts\indexnow-urls.example.txt -DryRun

.EXAMPLE
.\scripts\submit-indexnow.ps1 -UrlFile .\scripts\indexnow-urls.example.txt

.EXAMPLE
.\scripts\submit-indexnow.ps1 -Url 'https://www.northstarvisa.com/', 'https://www.northstarvisa.com/web/about/'
#>
[CmdletBinding()]
param(
    [string[]]$Url,

    [string]$UrlFile,

    [string]$Key,

    [string]$KeyFile,

    [string]$HostName = 'www.northstarvisa.com',

    [string]$KeyLocation,

    [ValidateRange(1, 10000)]
    [int]$BatchSize = 10000,

    [string]$Endpoint = 'https://api.indexnow.org/indexnow',

    [switch]$SkipKeyCheck,

    [switch]$DryRun
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($KeyFile)) {
    $KeyFile = Join-Path $PSScriptRoot '..\frontend\public\indexnow-key.txt'
}

function Get-IndexNowStatusMessage {
    param([int]$StatusCode)

    switch ($StatusCode) {
        200 { return '提交成功' }
        202 { return '已接收，密钥验证尚未完成' }
        400 { return '请求格式无效' }
        403 { return '密钥无效或密钥文件不可访问' }
        422 { return 'URL、主机名或密钥信息不匹配' }
        429 { return '请求过于频繁或触发反滥用限制' }
        default { return '未预期的响应状态' }
    }
}

function Get-HttpStatusCodeFromError {
    param([System.Management.Automation.ErrorRecord]$ErrorRecord)

    $response = $ErrorRecord.Exception.Response
    if ($null -eq $response) {
        return $null
    }

    try {
        return [int]$response.StatusCode
    }
    catch {
        return $null
    }
}

if ($HostName -notmatch '^[A-Za-z0-9.-]+$') {
    throw "主机名格式无效：$HostName"
}

if ([string]::IsNullOrWhiteSpace($Key)) {
    if (-not (Test-Path -LiteralPath $KeyFile -PathType Leaf)) {
        throw "未提供 IndexNow 密钥，且找不到密钥文件：$KeyFile"
    }

    $Key = (Get-Content -Raw -Encoding utf8 -LiteralPath $KeyFile).Trim()
}

if ($Key -notmatch '^[A-Za-z0-9-]{8,128}$') {
    throw 'IndexNow 密钥必须为 8 到 128 位，且只能包含字母、数字和连字符。'
}

if ([string]::IsNullOrWhiteSpace($KeyLocation)) {
    $KeyLocation = "https://$HostName/indexnow-key.txt"
}

$endpointUri = $null
if (-not [Uri]::TryCreate($Endpoint, [UriKind]::Absolute, [ref]$endpointUri) -or $endpointUri.Scheme -ne 'https') {
    throw "IndexNow 接口必须是有效的 HTTPS 地址：$Endpoint"
}

$keyLocationUri = $null
if (-not [Uri]::TryCreate($KeyLocation, [UriKind]::Absolute, [ref]$keyLocationUri) -or
    $keyLocationUri.Scheme -notin @('http', 'https') -or
    $keyLocationUri.Host -ine $HostName) {
    throw "密钥地址必须是当前主机上的 HTTP 或 HTTPS 地址：$KeyLocation"
}

$candidates = [System.Collections.Generic.List[string]]::new()
foreach ($item in @($Url)) {
    if ($null -ne $item) {
        $candidates.Add($item)
    }
}

if (-not [string]::IsNullOrWhiteSpace($UrlFile)) {
    if (-not (Test-Path -LiteralPath $UrlFile -PathType Leaf)) {
        throw "URL 文件不存在：$UrlFile"
    }

    foreach ($line in Get-Content -Encoding utf8 -LiteralPath $UrlFile) {
        $candidates.Add($line)
    }
}

$urls = [System.Collections.Generic.List[string]]::new()
$seen = [System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::Ordinal)

foreach ($candidate in $candidates) {
    $value = $candidate.Trim()
    if ([string]::IsNullOrWhiteSpace($value) -or $value.StartsWith('#')) {
        continue
    }

    $parsedUri = $null
    if (-not [Uri]::TryCreate($value, [UriKind]::Absolute, [ref]$parsedUri) -or
        $parsedUri.Scheme -notin @('http', 'https')) {
        throw "URL 格式无效：$value"
    }

    if ($parsedUri.Host -ine $HostName) {
        throw "URL 主机名必须是 $HostName：$value"
    }

    if (-not [string]::IsNullOrEmpty($parsedUri.Fragment)) {
        throw "URL 不能包含片段标识：$value"
    }

    if ($seen.Add($parsedUri.AbsoluteUri)) {
        $urls.Add($parsedUri.AbsoluteUri)
    }
}

if ($urls.Count -eq 0) {
    throw '没有可提交的 URL，请通过 -Url 或 -UrlFile 提供至少一个 URL。'
}

$batchCount = [Math]::Ceiling($urls.Count / [double]$BatchSize)
Write-Host "已加载 $($urls.Count) 个唯一 URL，共 $batchCount 批。"

if (-not $DryRun -and -not $SkipKeyCheck) {
    Write-Host "正在校验 IndexNow 密钥文件：$KeyLocation"
    try {
        $keyResponse = Invoke-WebRequest -Uri $KeyLocation -Method Get -UseBasicParsing
    }
    catch {
        throw "无法访问 IndexNow 密钥文件，请先部署前端静态文件：$KeyLocation。$($_.Exception.Message)"
    }

    if ($keyResponse.Content.Trim() -cne $Key) {
        throw "线上密钥文件内容与提交密钥不一致：$KeyLocation"
    }
}

if ($PSVersionTable.PSVersion.Major -le 5) {
    [Net.ServicePointManager]::SecurityProtocol =
        [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12
}

for ($start = 0; $start -lt $urls.Count; $start += $BatchSize) {
    $end = [Math]::Min($start + $BatchSize - 1, $urls.Count - 1)
    $batch = @($urls[$start..$end])
    $batchNumber = [int]($start / $BatchSize) + 1

    if ($DryRun) {
        Write-Host "[预览] 第 $batchNumber/$batchCount 批，共 $($batch.Count) 条："
        foreach ($item in $batch) {
            Write-Host "  $item"
        }
        continue
    }

    $payload = [ordered]@{
        host        = $HostName
        key         = $Key
        keyLocation = $KeyLocation
        urlList     = $batch
    }
    $json = $payload | ConvertTo-Json -Depth 3 -Compress
    $body = [Text.Encoding]::UTF8.GetBytes($json)

    Write-Host "正在提交第 $batchNumber/$batchCount 批，共 $($batch.Count) 条..."
    try {
        $response = Invoke-WebRequest `
            -Uri $Endpoint `
            -Method Post `
            -ContentType 'application/json; charset=utf-8' `
            -Body $body `
            -UseBasicParsing
        $statusCode = [int]$response.StatusCode
    }
    catch {
        $statusCode = Get-HttpStatusCodeFromError -ErrorRecord $_
        if ($null -eq $statusCode) {
            throw "第 $batchNumber 批提交失败：$($_.Exception.Message)"
        }

        $statusMessage = Get-IndexNowStatusMessage -StatusCode $statusCode
        throw "第 $batchNumber 批提交失败：HTTP $statusCode，$statusMessage。"
    }

    if ($statusCode -notin @(200, 202)) {
        $statusMessage = Get-IndexNowStatusMessage -StatusCode $statusCode
        throw "第 $batchNumber 批提交失败：HTTP $statusCode，$statusMessage。"
    }

    $statusMessage = Get-IndexNowStatusMessage -StatusCode $statusCode
    Write-Host "第 $batchNumber/$batchCount 批完成：HTTP $statusCode，$statusMessage。"
}

if ($DryRun) {
    Write-Host '预览完成，未向 IndexNow 发送请求。'
}
else {
    Write-Host "全部完成，共提交 $($urls.Count) 个 URL。"
}
