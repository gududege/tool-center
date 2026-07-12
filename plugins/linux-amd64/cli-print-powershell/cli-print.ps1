<#
.SYNOPSIS
    cli-print: a test tool that prints all CLI input it receives.
.DESCRIPTION
    Used as a target for Shim to verify argument forwarding, environment, etc.
    Outputs structured key=value lines for easy assertion in tests.
    Exit code = number of arguments (max 127), or custom via --exit-code N.
.NOTES
    To run this script, you may need to set the execution policy:
    Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
    Or run with bypass:
    powershell -ExecutionPolicy Bypass -File cli-print.ps1
#>

# Quote a string with C-style escaping for special characters.
# - null  -> <null>
# - Empty -> <empty>
# - Otherwise wraps in double quotes and escapes: \ " newline carriage-return tab
function Quote-String {
    param([string]$s)

    if ($null -eq $s) { return "<null>" }
    if ($s.Length -eq 0) { return "<empty>" }

    $sb = [System.Text.StringBuilder]::new($s.Length + 2)
    [void]$sb.Append('"')
    foreach ($c in $s.ToCharArray()) {
        switch ($c) {
            '"'  { [void]$sb.Append('\"') }
            '\'  { [void]$sb.Append('\\') }
            "`n" { [void]$sb.Append('\n') }
            "`r" { [void]$sb.Append('\r') }
            "`t" { [void]$sb.Append('\t') }
            default { [void]$sb.Append($c) }
        }
    }
    [void]$sb.Append('"')
    return $sb.ToString()
}

# Get the full path of this script (equivalent to Environment.ProcessPath)
$scriptPath = $PSCommandPath
if (-not $scriptPath) { $scriptPath = $MyInvocation.MyCommand.Path }

# Arguments: $args is an automatic variable containing all script arguments
$argList = @($args)

# Print process path
Write-Output "cli-print: $scriptPath"
# Print current working directory (use .NET for exact match with C# behavior)
Write-Output "cwd: $([Environment]::CurrentDirectory)"
# Print argument count
Write-Output "argc: $($argList.Count)"

# Print each argument with its index
for ($i = 0; $i -lt $argList.Count; $i++) {
    Write-Output "argv[$i]: $(Quote-String $argList[$i])"
}

# Determine exit code: default = number of args, clamped to [0, 127]
$exitCode = [Math]::Min($argList.Count, 127)

# Check for --exit-code N in arguments (search consecutive pairs)
for ($i = 0; $i -lt $argList.Count - 1; $i++) {
    if ($argList[$i] -eq "--exit-code") {
        $parsed = 0
        if ([int]::TryParse($argList[$i + 1], [ref]$parsed)) {
            # Clamp custom exit code to [0, 255]
            $exitCode = [Math]::Max(0, [Math]::Min($parsed, 255))
            break
        }
    }
}

# Get environment variables and sort keys case-insensitively.
# Uses the same StringComparer as the C# version (OrdinalIgnoreCase).
$envDict = [Environment]::GetEnvironmentVariables()
$envKeys = [System.Collections.ArrayList]::new()
foreach ($key in $envDict.Keys) {
    [void]$envKeys.Add([string]$key)
}
$envKeys.Sort([System.StringComparer]::OrdinalIgnoreCase)

# Print environment variable count
Write-Output "env: $($envKeys.Count) variables"

# Print each environment variable, sorted case-insensitively
foreach ($key in $envKeys) {
    $value = [Environment]::GetEnvironmentVariable($key)
    if ($null -eq $value) { $value = "" }
    Write-Output "env:${key}: $(Quote-String $value)"
}

# Print exit code
Write-Output "exit_code: $exitCode"

# Exit with the determined code
exit $exitCode
