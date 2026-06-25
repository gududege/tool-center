# TIA Exporter

A command-line tool that exports content from **Siemens TIA Portal** projects using the official Openness API. Supports TIA Portal V17 - V20.

## Features

- **Batch export** — scan a directory for TIA Portal project files (`.ap17`-`.ap20`, `.zap17`-`.zap20`) and export them all
- **Multiple export modes** — hardware, blocks, tag tables, UDTs, watch tables, cross-references, or any combination
- **Parallel processing** — configurable worker count for concurrent project export
- **Flexible output formats** — SimaticML, ExternalSource (SCL/AWL/DB), SimaticSD
- **Folder structure** — optionally preserve the original TIA Portal group hierarchy
- **Incremental update** — re-export only specific types without clearing other output (`--incremental-update`)
- **Export statistics** — per-PLC, per-type success/failure counts recorded in `manifest.json`
- **Enum validation** — invalid option values produce a clear error and exit code 1
- **Graceful shutdown** — Ctrl+C immediately releases TIA Portal and aborts remaining operations
- **UMAC authentication** — support for password-protected projects
- **Safety blocks** — dedicated format options for Safety blocks and UDTs

## Prerequisites

- **Windows** (required by the TIA Portal Openness API)
- **.NET Framework 4.8** runtime
- **Siemens TIA Portal** V17, V18, V19, or V20 installed (the Openness API DLLs are loaded at runtime from the TIA Portal installation)

> [!NOTE]
> The tool must be run on a machine where TIA Portal is installed. It detects the installed version automatically and loads the matching `Siemens.Engineering.dll` at runtime.

## Installation

Download `tia-export.exe` from the [latest release](../../releases/latest) and place it anywhere on your PATH.

## Usage

```text
tia-export.exe --indir <PATH> --outdir <PATH> [options]

Required:
  --indir <PATH>          Input directory containing TIA Portal projects
  --outdir <PATH>         Output directory for exported files

Options:
  --logfile <PATH>        Log file path
  --loglevel <LEVEL>      Minimum log level: debug, info, warn, error
  --portal-mode <MODE>    WithoutUserInterface (default) or WithUserInterface
  --keep-folder-structure  Preserve TIA Portal group hierarchy in output
  --incremental-update     Do not delete other export types' output; only overwrite matching types
  --export-mode <MODE>    Comma-separated: hardware, block, tag, udt, watchtable, reference, all (default)
  --max-workers <N>       Maximum parallel workers (default: auto)
  --project-filter <REGEX> Regex filter for project names
  --umac-user <USER>      UMAC local user name
  --umac-password <PASS>  UMAC local user password

Block format options:
  --scl-format <FORMAT>     ExternalSource (default), SimaticML, SimaticSD
  --stl-format <FORMAT>     ExternalSource (default), SimaticML
  --lad-format <FORMAT>    SimaticML (default), SimaticSD
  --db-format <FORMAT>     ExternalSource (default), SimaticML, SimaticSD

UDT format options:
  --udt-format <FORMAT>            ExternalSource (default), SimaticML, SimaticSD
  --safety-db-format <FORMAT>       SimaticML (default), SimaticSD
  --safety-udt-format <FORMAT>      SimaticML (default), SimaticSD
```

### Examples

Export everything from all projects in a directory:

```shell
tia-export.exe --indir D:\Projects --outdir D:\Exports
```

Export only blocks and tags, keeping folder structure:

```shell
tia-export.exe --indir D:\Projects --outdir D:\Exports --export-mode block,tag --keep-folder-structure
```

Export with debug logging to a file:

```shell
tia-export.exe --indir D:\Projects --outdir D:\Exports --loglevel debug --logfile export.log
```

## Output Structure

```
outdir/
  <ProjectName>/
    manifest.json
    hardware.aml
    <PlcName>/
      blocks/
        Organization_Block/   (*.xml, *.scl, *.awl, *.db, …)
        Function_Block/
        Function/
        Data_Block/
      tagtables/              (*.xml)
      udts/                   (*.xml, *.udt)
      watchtables/            (*.xml)
      references/             (*.json)
```

With `--keep-folder-structure`, the TIA Portal group hierarchy is preserved as nested directories inside each block type folder.

### manifest.json

Each project folder contains a `manifest.json` with export metadata and per-PLC, per-type statistics:

```json
{
  "ProjectName": "MyProject",
  "TiaVersion": "V19",
  "ExportTime": "2026-06-24T12:00:00.0000000Z",
  "ExporterVersion": "1.0.0",
  "Plcs": ["PLC_1", "PLC_2"],
  "HardwareExport": {
    "ExportType": "hardware",
    "ExportTime": "2026-06-24T12:00:00.0000000Z",
    "SuccessCount": 1,
    "FailureCount": 0
  },
  "PlcExports": [
    {
      "PlcName": "PLC_1",
      "Results": [
        { "ExportType": "blocks", "ExportTime": "...", "SuccessCount": 42, "FailureCount": 0 },
        { "ExportType": "tags", "ExportTime": "...", "SuccessCount": 5, "FailureCount": 0 },
        { "ExportType": "udts", "ExportTime": "...", "SuccessCount": 12, "FailureCount": 1 },
        { "ExportType": "watchtables", "ExportTime": "...", "SuccessCount": 2, "FailureCount": 0 },
        { "ExportType": "references", "ExportTime": "...", "SuccessCount": 38, "FailureCount": 0 }
      ]
    }
  ]
}
```

## Use with --incremental-update

Re-export specific types without clearing previously exported content:

```shell
# Initial full export
tia-export.exe --indir D:\Projects --outdir D:\Exports

# Later: re-export only UDTs, keep blocks/tags/etc. intact
tia-export.exe --indir D:\Projects --outdir D:\Exports --export-mode udt --incremental-update
```

In incremental mode, only the folders/files corresponding to the selected `--export-mode` are replaced; all other export types' output is preserved.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Invalid arguments |
| 2 | Project scan failed |
| 3 | TIA Portal not found |
| 4 | Siemens.Engineering.dll load failed |
| 5 | Project open failed |
| 6 | Export failed |
| 7 | Worker failed |
| 8 | Unexpected error |

## Building from Source

```shell
git clone https://github.com/<your-repo>/tia-exporter.git
cd tia-exporter
dotnet build src/TiaExporter.Cli/TiaExporter.Cli.csproj -c Release
```

The built executable is in `src/TiaExporter.Cli/bin/Release/net48/TiaExporter.Cli.exe`.

## License

This project uses the Siemens TIA Portal Openness API. Ensure you comply with the Siemens license terms when using this tool.