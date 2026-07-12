#!/usr/bin/env python3
"""
cli-print: a test tool that prints all CLI input it receives.
Used as a target for Shim to verify argument forwarding, environment, etc.
Outputs structured key=value lines for easy assertion in tests.
Exit code = number of arguments (max 127), or custom via --exit-code N.
"""

import os
import sys


def quote(s):
    """Quote a string with C-style escaping for special characters.

    - None  -> <null>
    - Empty -> <empty>
    - Otherwise wraps in double quotes and escapes: \\ \" \\n \\r \\t
    """
    if s is None:
        return "<null>"
    if len(s) == 0:
        return "<empty>"
    result = ['"']
    for c in s:
        if c == '"':
            result.append('\\"')
        elif c == '\\':
            result.append('\\\\')
        elif c == '\n':
            result.append('\\n')
        elif c == '\r':
            result.append('\\r')
        elif c == '\t':
            result.append('\\t')
        else:
            result.append(c)
    result.append('"')
    return ''.join(result)


def try_parse_int(s):
    """Try to parse a string as a 32-bit signed integer.

    Mimics C# int.TryParse: allows optional leading sign and digits.
    Returns the integer value, or None if parsing fails.
    """
    if s is None:
        return None
    trimmed = s.strip()
    # Match optional sign followed by one or more digits
    if not trimmed:
        return None
    sign = 1
    idx = 0
    if trimmed[0] in ('+', '-'):
        if trimmed[0] == '-':
            sign = -1
        idx = 1
    if idx >= len(trimmed):
        return None
    if not trimmed[idx:].isdigit():
        return None
    try:
        val = sign * int(trimmed[idx:])
    except (ValueError, OverflowError):
        return None
    # Match C# int32 range
    if val < -2147483648 or val > 2147483647:
        return None
    return val


def main():
    # Get the full path of this script (equivalent to Environment.ProcessPath)
    script_path = os.path.abspath(sys.argv[0])

    # Arguments: sys.argv[0] is the script name, so actual args start at index 1
    args = sys.argv[1:]

    # Print process path
    print(f"cli-print: {script_path}")
    # Print current working directory
    print(f"cwd: {os.getcwd()}")
    # Print argument count
    print(f"argc: {len(args)}")
    # Print each argument with its index
    for i, arg in enumerate(args):
        print(f"argv[{i}]: {quote(arg)}")

    # Determine exit code: default = number of args, clamped to [0, 127]
    exit_code = min(len(args), 127)

    # Check for --exit-code N in arguments (search consecutive pairs)
    for i in range(len(args) - 1):
        if args[i] == "--exit-code":
            parsed = try_parse_int(args[i + 1])
            if parsed is not None:
                # Clamp custom exit code to [0, 255]
                exit_code = max(0, min(parsed, 255))
                break

    # Get environment variables and sort keys case-insensitively
    env = dict(os.environ)
    keys = sorted(env.keys(), key=lambda k: k.lower())

    # Print environment variable count
    print(f"env: {len(keys)} variables")
    # Print each environment variable, sorted case-insensitively
    for key in keys:
        value = env.get(key, "")
        print(f"env:{key}: {quote(value)}")

    print(f"exit_code: {exit_code}")
    return exit_code


if __name__ == "__main__":
    sys.exit(main())
