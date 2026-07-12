#!/usr/bin/env node
'use strict';

/**
 * cli-print: a test tool that prints all CLI input it receives.
 * Used as a target for Shim to verify argument forwarding, environment, etc.
 * Outputs structured key=value lines for easy assertion in tests.
 * Exit code = number of arguments (max 127), or custom via --exit-code N.
 */

const path = require('path');

/**
 * Quote a string with C-style escaping for special characters.
 *
 * - null/undefined -> <null>
 * - Empty string   -> <empty>
 * - Otherwise wraps in double quotes and escapes: \\ \" \n \r \t
 */
function quote(s) {
    if (s === null || s === undefined) return "<null>";
    if (s.length === 0) return "<empty>";
    let result = '"';
    for (const c of s) {
        switch (c) {
            case '"':  result += '\\"';  break;
            case '\\': result += '\\\\'; break;
            case '\n': result += '\\n';  break;
            case '\r': result += '\\r';  break;
            case '\t': result += '\\t';  break;
            default:   result += c;      break;
        }
    }
    result += '"';
    return result;
}

/**
 * Try to parse a string as a 32-bit signed integer.
 * Mimics C# int.TryParse: allows optional leading sign and digits.
 * Returns the integer value, or NaN if parsing fails.
 */
function tryParseInt(s) {
    if (s === null || s === undefined) return NaN;
    const trimmed = String(s).trim();
    // Match optional sign followed by one or more digits only
    if (!/^[+-]?\d+$/.test(trimmed)) return NaN;
    const n = parseInt(trimmed, 10);
    // Match C# int32 range
    if (n < -2147483648 || n > 2147483647) return NaN;
    return n;
}

function main() {
    // Get the full path of this script (equivalent to Environment.ProcessPath)
    // process.argv[0] = node executable, process.argv[1] = script path
    const scriptPath = path.resolve(process.argv[1]);

    // Arguments: skip node (index 0) and script path (index 1)
    const args = process.argv.slice(2);

    // Print process path
    console.log(`cli-print: ${scriptPath}`);
    // Print current working directory
    console.log(`cwd: ${process.cwd()}`);
    // Print argument count
    console.log(`argc: ${args.length}`);
    // Print each argument with its index
    for (let i = 0; i < args.length; i++) {
        console.log(`argv[${i}]: ${quote(args[i])}`);
    }

    // Determine exit code: default = number of args, clamped to [0, 127]
    let exitCode = Math.min(args.length, 127);

    // Check for --exit-code N in arguments (search consecutive pairs)
    for (let i = 0; i < args.length - 1; i++) {
        if (args[i] === "--exit-code") {
            const parsed = tryParseInt(args[i + 1]);
            if (!isNaN(parsed)) {
                // Clamp custom exit code to [0, 255]
                exitCode = Math.max(0, Math.min(parsed, 255));
                break;
            }
        }
    }

    // Get environment variable keys and sort case-insensitively
    const keys = Object.keys(process.env).sort((a, b) => {
        const la = a.toLowerCase();
        const lb = b.toLowerCase();
        if (la < lb) return -1;
        if (la > lb) return 1;
        return 0;
    });

    // Print environment variable count
    console.log(`env: ${keys.length} variables`);
    // Print each environment variable, sorted case-insensitively
    for (const key of keys) {
        const value = process.env[key] || "";
        console.log(`env:${key}: ${quote(value)}`);
    }

    console.log(`exit_code: ${exitCode}`);
    return exitCode;
}

process.exit(main());
