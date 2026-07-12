param(
    [switch]$Help,
    [string]$OutputFile = "block_calls.json",
    [string[]]$ScanPaths = @("udts", "blocks")
)

if ($Help) {
    @"
Usage: .\scan_calls.ps1 [-OutputFile <path>] [-ScanPaths <dir[]>] [-Help]

Description:
  Scans files under the specified directories, detects cross-references
  between blocks, and outputs a JSON call graph.

Options:
  -OutputFile   Path to the output JSON file (default: block_calls.json)
  -ScanPaths    Array of directories to scan (default: udts, blocks)
  -Help         Show this help message

Example:
  .\scan_calls.ps1
  .\scan_calls.ps1 -OutputFile "result.json"
  .\scan_calls.ps1 -ScanPaths @("udts", "blocks\Function")
  .\scan_calls.ps1 -ScanPaths @("myDir") -OutputFile "out.json"
"@
    exit 0
}

$ErrorActionPreference = "Stop"

$baseDir = Split-Path -Parent $PSCommandPath

function Get-BlockType {
    param([string]$FilePath, [string]$Content)

    $firstLine = if ($Content) { 
        ($Content -split "`r?`n")[0] 
    } else { 
        Get-Content -LiteralPath $FilePath -TotalCount 1 -Encoding UTF8 
    }

    if ($firstLine -match '^ORGANIZATION_BLOCK\s+"')  { return "OB" }
    if ($firstLine -match '^FUNCTION_BLOCK\s+"')      { return "FB" }
    if ($firstLine -match '^FUNCTION\s+"')            { return "FC" }
    if ($firstLine -match '^DATA_BLOCK\s+"')          {
        if (-not $Content) {
            $Content = Get-Content -LiteralPath $FilePath -Raw -Encoding UTF8
        }
        $hasVar    = ($Content -match '\bVAR\b') -and ($Content -match '\bEND_VAR\b')
        $hasStruct = ($Content -match '\bSTRUCT\b') -and ($Content -match '\bEND_STRUCT\b')

        if (-not $hasVar -and -not $hasStruct) {
            $fileName = [System.IO.Path]::GetFileNameWithoutExtension($FilePath)
            $lines = if ($Content) {
                ($Content -split "`r?`n")
            } else {
                @(Get-Content -LiteralPath $FilePath -Encoding UTF8)
            }
            $line5 = if ($lines.Count -ge 5) { $lines[4] } else { "" }

            if (($fileName -match '^DI_') -or ($line5 -match '^"FB')) {
                return "DI"
            }
        }
        return "DB"
    }
    if ($firstLine -match '^TYPE\s+"')                { return "UDT" }

    return "UNKNOWN"
}

function Find-Files {
    param([string[]]$Paths)

    $files = @()
    foreach ($scanPath in $Paths) {
        $fullPath = Join-Path $baseDir $scanPath
        if (-not (Test-Path -LiteralPath $fullPath -PathType Container)) {
            Write-Warning "Directory not found, skipped: $fullPath"
            continue
        }

        $items = Get-ChildItem -LiteralPath $fullPath -Recurse -File
        $files += @($items)
    }
    return $files
}

$allFiles = Find-Files $ScanPaths

$blockNames = @($allFiles | ForEach-Object { [System.IO.Path]::GetFileNameWithoutExtension($_.FullName) })

$total = $allFiles.Count
Write-Host "Found $total files to scan"

$results = [System.Collections.ArrayList]::new()

$escapedNames = $blockNames | ForEach-Object { [regex]::Escape($_) }
$regexPattern = '\b(' + ($escapedNames -join '|') + ')\b'
$regex = [regex]::new($regexPattern)

for ($i = 0; $i -lt $allFiles.Count; $i++) {
    $file = $allFiles[$i]
    $fileName = [System.IO.Path]::GetFileNameWithoutExtension($file.FullName)

    Write-Progress -Activity "Scanning files" -Status "$fileName" -PercentComplete (($i / $total) * 100)

    $content = Get-Content -LiteralPath $file.FullName -Raw -Encoding UTF8
    if (-not $content) { continue }

    $fileType = Get-BlockType -FilePath $file.FullName -Content $content

    $foundRefs = [System.Collections.ArrayList]::new()
    $matches = $regex.Matches($content)

    foreach ($match in $matches) {
        $refName = $match.Value
        if ($refName -eq $fileName) { continue }
        if ($foundRefs -notcontains $refName) {
            [void]$foundRefs.Add($refName)
        }
    }

    $useList = @()
    foreach ($refName in $foundRefs) {
        $refFile = $allFiles | Where-Object { [System.IO.Path]::GetFileNameWithoutExtension($_.FullName) -eq $refName } | Select-Object -First 1
        $refType = if ($refFile) { Get-BlockType -FilePath $refFile.FullName } else { "UNKNOWN" }
        $useList += [PSCustomObject]@{ name = $refName; type = $refType }
    }

    [void]$results.Add([PSCustomObject]@{
        name = $fileName
        type = $fileType
        use  = $useList
    })
}

Write-Host "Formatting and writing to $OutputFile ..."

$sb = [System.Text.StringBuilder]::new()
[void]$sb.AppendLine('[')

$sortedResults = @($results | Sort-Object { $_.name })
$total = $sortedResults.Count

for ($i = 0; $i -lt $total; $i++) {
    $block = $sortedResults[$i]
    $comma = if ($i -lt $total - 1) { ',' } else { '' }

    [void]$sb.AppendLine("  {")
    [void]$sb.AppendLine("    ""name"": ""$($block.name)"",")
    [void]$sb.AppendLine("    ""type"": ""$($block.type)"",")

    $useItems = $block.use
    $useCount = $useItems.Count

    if ($useCount -eq 0) {
        [void]$sb.AppendLine("    ""use"": []")
    } else {
        [void]$sb.AppendLine("    ""use"": [")
        for ($j = 0; $j -lt $useCount; $j++) {
            $item = $useItems[$j]
            $itemComma = if ($j -lt $useCount - 1) { ',' } else { '' }
            [void]$sb.AppendLine("      { ""name"": ""$($item.name)"", ""type"": ""$($item.type)"" }$itemComma")
        }
        [void]$sb.AppendLine("    ]")
    }

    [void]$sb.AppendLine("  }$comma")
}

[void]$sb.AppendLine(']')

Set-Content -LiteralPath (Join-Path $baseDir $OutputFile) -Value $sb.ToString() -Encoding UTF8

Write-Host "Done. Processed $($results.Count) blocks."
