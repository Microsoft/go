# Copyright (c) Microsoft Corporation.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

<#
.DESCRIPTION
This script builds and runs a tool defined in 'eng/_util'.

To run a tool:
  run.ps1 <tool> [arguments...]

For example, to build the repository:
  run.ps1 build

To list all possible tools:
  run.ps1

Builds 'eng/_util/cmd/<tool>/<tool>.go' and runs it using the list of arguments.

This command automatically installs a known version of Microsoft Go that will be
used to build the tools. The known version of Go will also be used to build the
Go source code, if it's built. Set environment variable "MS_USE_PATH_GO" to 1 to
your own Go from PATH instead.

Every tool accepts a '-h' argument to show tool usage help.
#>

$ErrorActionPreference = 'Stop'

# Take the first arg as the tool name, and "shift" the rest so we can pass "@args" to the tool.
#
# "param($tool)" would make PowerShell eagerly bind partial matches ("-too", "-to", "-t"). These
# matches overlap with tool args and cause unexpected behavior. ("sync" uses "-to", for example.)
$tool = $args[0]
$args = $args[1..$args.Count]

# Import utilities. May throw if our version of PowerShell is too old.
. (Join-Path $PSScriptRoot "utilities.ps1")
if ($LASTEXITCODE) {
  throw "Dot-sourcing utilities.ps1 failed with non-null/non-zero exit code. ($LASTEXITCODE)"
}

function Write-ToolList() {
  Write-Host "Possible tools:"
  foreach ($tool in Get-ChildItem (Join-Path $PSScriptRoot "_util" "cmd" "*")) {
    Write-Host "  $($tool.Name)"
  }
  Write-Host ""
}

if (-not $tool) {
  Write-Host "No tool specified. Showing help and listing available tools:"
  Write-Host ""
  ((Get-Help $PSCommandPath).DESCRIPTION | Out-String).Trim() | Write-Host
  Write-Host ""
  Write-ToolList
  exit 0
}

# Find tool script file based on the name given.
$tool_search = Join-Path $PSScriptRoot "_util" "cmd" "$tool" "$tool.go"
# Find matches, and force the result to be an array.
$tool_matches = @(Get-Item $tool_search)

if ($tool_matches.Count -gt 1) {
  Write-Host "Error: Multiple tools match '$tool_search'. Found:"
  $tool_matches | Write-Host
  Write-Host "This is a most likely a repository infrastructure issue. Every name should be unique."
  exit 1
} elseif ($tool_matches.Count -lt 1) {
  Write-Host "Error: No tools found matching '$tool_search'."
  Write-ToolList
  exit 1
}

$tool_source = $tool_matches[0]
if (-not ($tool_source -is [System.IO.FileInfo])) {
  Write-Host "Found tool source code, but it is not a file: $tool_source"
  exit 1
}

# Now that we have a single result, navigate upwards to see which module it's in.
$tool_module = $tool_source.Directory.Parent.Parent.FullName

# Download a consistent stage 0 version of Go unless opted out.
if ($env:MS_USE_PATH_GO -eq "1") {
  try {
    Write-Host "Using $(go version) from '$(go env GOROOT)'. Results may differ from CI environment."
  } catch {
    Write-Host "Error: 'go' is most likely not in PATH. To download the known version, set 'MS_USE_PATH_GO' to '0' or unset it, then try again."
    Write-Host "Exception: $_"
    exit 1
  }
} else {
  Download-Stage0
}

$stage0_goroot = & go env GOROOT

# The tool may need to know where our copy of Go is located. Save it in env to give it access. Don't
# pass it to the tool as an arg, becuase that would complicate arg handling in each tool.
$env:STAGE_0_GOROOT = $stage0_goroot

# Decide where to place the compiled tool.
$tool_output = Join-Path $PSScriptRoot "artifacts" "toolbin" $tool
if ($IsWindows) {
  $tool_output += ".exe"
}

try {
  # Move into module so "go build" detects it and fetches dependencies.
  Push-Location $tool_module
  # Use a module-local path so Go resolves imports correctly.
  $module_local_script_path = Join-Path "." "cmd" "$tool"

  # The caller may have passed in GOOS/GOARCH to cross-compile Go. We can't use those values here:
  # we need to be able to run the tool on the host, so we must always target the host OS/ARCH. Clear
  # out the GOOS/GOARCH values (empty string) to detect host OS/ARCH automatically for the tools.
  Invoke-CrossGoBlock "" "" {
    Write-Host "In '$tool_module', building '$module_local_script_path' -> $tool_output"
    & (Join-Path $stage0_goroot "bin" "go") build -o $tool_output $module_local_script_path
    if ($LASTEXITCODE) {
      Write-Host "Failed to build tool."
      exit 1
    }
  }

  Write-Host "Built '$tool'. Running from repo root..."
} finally {
  Pop-Location
}

try {
  # Run tool from the root of the repo.
  Push-Location (Join-Path $PSScriptRoot "..")
  & "$tool_output" @args
  if ($LASTEXITCODE) {
    Write-Host "Failed to run tool."
    exit 1
  }
} finally {
  Pop-Location
}
