@echo off
REM cli-print: a test tool that prints all CLI input it receives.
REM Used as a target for Shim to verify argument forwarding, environment, etc.
REM Outputs structured key=value lines for easy assertion in tests.
REM Exit code = number of arguments (max 127), or custom via --exit-code N.
REM
REM Limitations: Control characters (newline, tab, etc.) are not escaped.
REM Exclamation marks (!) in values may cause issues with delayed expansion.

REM Save environment variables to temp file BEFORE setting any local variables.
REM This ensures we capture only the original environment, not script-local vars.
set "_CLI_PRINT_TMP=%TEMP%\_cli_print_env_%RANDOM%.tmp"
set > "%_CLI_PRINT_TMP%"

setlocal enabledelayedexpansion

REM Get the full path of this script (equivalent to Environment.ProcessPath)
echo cli-print: %~f0

REM Print current working directory
echo cwd: %CD%

REM Collect arguments into array and count them.
REM Use shift loop since batch has no direct argc equivalent.
set argc=0
:collect_args
if "%~1"=="" goto args_done
set "argv[!argc!]=%~1"
set /a argc+=1
shift
goto collect_args
:args_done

REM Print argument count
echo argc: !argc!

REM Print each argument with its index
if !argc! gtr 0 (
    set /a last_idx=argc-1
    for /l %%i in (0,1,!last_idx!) do (
        call :quote "!argv[%%i]!"
        echo argv[%%i]: !quoted!
    )
)

REM Determine exit code: default = argc clamped to [0, 127]
set exit_code=!argc!
if !exit_code! gtr 127 set exit_code=127

REM Check for --exit-code N in arguments (search consecutive pairs)
set /a search_end=argc-2
for /l %%i in (0,1,!search_end!) do (
    if "!argv[%%i]!"=="--exit-code" (
        set /a next_idx=%%i+1
        REM Use call set for double expansion of array index
        call set next_val=%%argv[!next_idx!]%%
        call :check_int "!next_val!"
        if "!is_valid_int!"=="1" (
            set /a code=!next_val!
            if !code! lss 0 set code=0
            if !code! gtr 255 set code=255
            set exit_code=!code!
            goto exit_code_found
        )
    )
)
:exit_code_found

REM Count environment variables from temp file
set env_count=0
for /f "tokens=1,* delims==" %%a in (%_CLI_PRINT_TMP%) do set /a env_count+=1
echo env: !env_count! variables

REM Print each environment variable from temp file.
REM set command outputs variables sorted alphabetically (case-insensitive on Windows).
for /f "tokens=1,* delims==" %%a in (%_CLI_PRINT_TMP%) do (
    set "envkey=%%a"
    set "envval=%%b"
    call :quote "!envval!"
    echo env:!envkey!: !quoted!
)

REM Clean up temp file
del "%_CLI_PRINT_TMP%" 2>nul

REM Exit with the determined code
echo exit_code: !exit_code!
exit /b !exit_code!

REM ===== Subroutines =====

:quote
REM Quote a string with C-style escaping.
REM Handles: \ -> \\, " -> \"
REM Note: Control characters (newline, tab, etc.) are not escaped in batch.
set "input=%~1"
if "!input!"=="" (
    set "quoted=<empty>"
    goto :eof
)
REM Escape backslash first, then double-quote (order matters)
set "output=!input:\=\\!"
set "output=!output:"=\"!"
REM Wrap in quotes using a helper variable to avoid parser issues
set "q=""
set "quoted=!q!!output!!q!"
goto :eof

:check_int
REM Check if a string is a valid integer (optional sign + digits).
REM Sets is_valid_int=1 if valid, 0 otherwise.
set "is_valid_int=1"
set "check_str=%~1"
if "!check_str!"=="" set "is_valid_int=0" & goto :eof
REM Remove optional leading sign
if "!check_str:~0,1!"=="+" set "check_str=!check_str:~1!"
if "!check_str:~0,1!"=="-" set "check_str=!check_str:~1!"
if "!check_str!"=="" set "is_valid_int=0" & goto :eof
REM Check if remaining characters are all digits.
REM for /f with delims=0123456789 extracts non-digit parts; if any found, not a number.
set "non_digit="
for /f "delims=0123456789" %%d in ("!check_str!") do set "non_digit=%%d"
if defined non_digit set "is_valid_int=0"
goto :eof
