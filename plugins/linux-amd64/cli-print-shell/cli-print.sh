#!/usr/bin/env bash
# cli-print: a test tool that prints all CLI input it receives.
# Used as a target for Shim to verify argument forwarding, environment, etc.
# Outputs structured key=value lines for easy assertion in tests.
# Exit code = number of arguments (max 127), or custom via --exit-code N.
#
# Note: This script requires bash 4.0+ for mapfile builtin.

# Quote a string with C-style escaping for special characters.
# Uses awk for reliable handling of all edge cases including embedded newlines.
# Output: <empty> for empty string, or "..." with escaped \ " \n \r \t
quote() {
    local s="$1"
    if [ -z "$s" ]; then
        printf '<empty>'
        return
    fi
    # awk processes line by line; NR>1 means we had a newline between lines
    printf '%s' "$s" | awk '
    BEGIN { printf "\"" }
    {
        gsub(/\\/, "\\\\")
        gsub(/"/, "\\\"")
        gsub(/\t/, "\\t")
        gsub(/\r/, "\\r")
        if (NR > 1) printf "\\n"
        printf "%s", $0
    }
    END { printf "\"" }
    '
}

# Get the full path of this script (equivalent to Environment.ProcessPath)
script_path=$(readlink -f "$0" 2>/dev/null || echo "$0")

# Print process path
echo "cli-print: ${script_path}"
# Print current working directory
echo "cwd: $(pwd)"
# Print argument count ($# does not include $0)
echo "argc: $#"

# Print each argument with its index
i=0
for arg in "$@"; do
    echo "argv[${i}]: $(quote "$arg")"
    i=$((i + 1))
done

# Determine exit code: default = number of args, clamped to [0, 127]
argc=$#
if [ "$argc" -gt 127 ]; then
    exit_code=127
else
    exit_code=$argc
fi

# Check for --exit-code N in arguments (search consecutive pairs)
args=("$@")
for ((i = 0; i < ${#args[@]} - 1; i++)); do
    if [ "${args[i]}" = "--exit-code" ]; then
        next="${args[i+1]}"
        # Check if next argument is a valid integer (optional sign + digits)
        if [[ "$next" =~ ^[+-]?[0-9]+$ ]]; then
            code=$((next))
            # Clamp custom exit code to [0, 255]
            if [ "$code" -lt 0 ]; then
                code=0
            elif [ "$code" -gt 255 ]; then
                code=255
            fi
            exit_code=$code
            break
        fi
    fi
done

# Get environment variable keys, sorted case-insensitively.
# Use env -0 for null-separated output to handle multi-line values correctly.
# This includes special variables like _ that compgen -e may miss.
# Fallback to env | cut for systems without env -0 support.
if env -0 >/dev/null 2>&1; then
    # GNU coreutils: env -0 uses null separators
    mapfile -d '' -t env_entries < <(env -0)
    env_keys=()
    for entry in "${env_entries[@]}"; do
        # Extract key (everything before the first =)
        env_keys+=("${entry%%=*}")
    done
    # Sort case-insensitively
    mapfile -t env_keys < <(printf '%s\n' "${env_keys[@]}" | sort -f)
else
    # Fallback: use env with cut (may not handle multi-line values correctly)
    mapfile -t env_keys < <(env | cut -d= -f1 | sort -f)
fi

# Print environment variable count
echo "env: ${#env_keys[@]} variables"

# Print each environment variable, sorted case-insensitively
for key in "${env_keys[@]}"; do
    # Use printenv to get the value directly from the environment
    # This handles special variables like _ that ${!key} may not resolve
    value=$(printenv "$key" 2>/dev/null || echo "")
    echo "env:${key}: $(quote "$value")"
done

echo "exit_code: ${exit_code}"
exit "$exit_code"
