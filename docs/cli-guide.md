# VSP CLI Guide

**vsp** provides a complete ABAP development toolchain from the terminal. Every feature works as a single-binary CLI — no SAP GUI, no Eclipse, no IDE required.

## Quick Start

```bash
# Connect to SAP
export SAP_URL=https://host:44300 SAP_USER=dev SAP_PASSWORD=secret

# Or use saved profiles
vsp -s dev search "ZCL_*"
```

## Command Reference

### Source Code

```bash
# Read source
vsp source read CLAS ZCL_MY_CLASS
vsp source read PROG ZTEST_REPORT

# Read with compressed dependency context
vsp context CLAS ZCL_MY_CLASS
vsp context CLAS ZCL_MY_CLASS --max-deps 30
vsp context CLAS ZCL_DEEP --depth 2             # deps of deps
vsp context CLAS ZCL_COMPLEX --depth 3           # 3 levels deep

# Write source
vsp source write CLAS ZCL_MY_CLASS < new_source.abap
cat source.abap | vsp source write PROG ZTEST

# Surgical edit (find & replace)
vsp source edit CLAS ZCL_MY_CLASS --old "old_code" --new "new_code"
```

**Requirements:** Standard ADT. No ZADT_VSP needed.

### Search & Discovery

```bash
# Search objects by name
vsp search "ZCL_ORDER*"
vsp search "Z*" --type CLAS --max 50

# Search source code across packages
vsp grep "SELECT.*FROM.*mara" --package '$TMP'
vsp grep "TYPE REF TO" --package 'ZFINANCE' -i
vsp grep "cl_abap_unit" --package '$ZADT' --type CLAS

# System information
vsp system info
```

**Requirements:** Standard ADT. No ZADT_VSP needed.

### Call Graph & Dependency Analysis

```bash
# Show what an object uses (callees)
vsp graph CLAS ZCL_MY_CLASS
vsp graph CLAS ZCL_MY_CLASS --depth 2

# Show what uses an object (callers / where-used)
vsp graph CLAS ZCL_MY_CLASS --direction callers

# Both directions
vsp graph CLAS ZCL_MY_CLASS --direction both

# Works for all object types
vsp graph INTF ZIF_MY_INTERFACE --direction callers
vsp graph PROG ZREPORT
vsp graph TRAN SE80           # resolves transaction → program automatically

# Package dependency analysis + transport readiness
vsp deps '$ZADT_VSP'
vsp deps '$ZADT_VSP' --format summary
vsp deps '$ZFINANCE' --include-subpackages
vsp deps '$TMP' --object ZCL_MY_CLASS
```

`graph` uses ADT call graph API when available, falls back to **WBCROSSGT + CROSS** tables automatically (same approach as SAP's where-used list).

`deps` analyzes all objects in a package and classifies references as:
- **Internal** — within the same package (safe)
- **External custom** — Z/Y objects in other packages (must transport first)
- **SAP standard** — always available on target system

**Requirements:** Standard ADT. No ZADT_VSP needed.

### Database Queries

```bash
# Query tables
vsp query T000
vsp query T000 --top 5
vsp query USR02 --where "BNAME = 'DEVELOPER'" --top 10
vsp query DD03L --where "TABNAME = 'MARA'" --fields "FIELDNAME,DATATYPE,LENG"
vsp query TADIR --where "DEVCLASS = '\$TMP'" --top 20 --order "OBJ_NAME"
```

**Requirements:** Standard ADT. No ZADT_VSP needed.
Safety: Use `--block-free-sql` to prevent arbitrary SQL execution.

### Testing & Quality

```bash
# Unit tests
vsp test CLAS ZCL_MY_CLASS
vsp test --package '$TMP'

# ATC checks
vsp atc CLAS ZCL_MY_CLASS
vsp atc PROG ZTEST --variant MY_VARIANT

# ABAP Lint (offline — no SAP needed!)
vsp lint CLAS ZCL_MY_CLASS              # fetch from SAP
vsp lint --file myclass.clas.abap       # local file
echo "DATA x." | vsp lint --stdin       # piped input
vsp lint --file src.abap --max-length 100
```

Lint rules: `line_length`, `empty_statement`, `obsolete_statement`, `max_one_statement`, `preferred_compare_operator`, `colon_missing_space`, `local_variable_names`.

**Requirements:** `lint` works fully offline. `test`/`atc` need standard ADT.

### Compile & Transpile

```bash
# WASM → ABAP (offline)
vsp compile wasm program.wasm
vsp compile wasm program.wasm --class ZCL_MY_WASM
vsp compile wasm program.wasm -o ./src/

# TypeScript → ABAP (needs Node.js)
vsp compile ts lexer.ts --prefix zcl_
vsp compile ts lexer.ts -o ./src/

# Parse ABAP (offline)
vsp parse --file myclass.clas.abap --format summary
echo "DATA lv_x TYPE i." | vsp parse --stdin
vsp parse CLAS ZCL_TEST --format json
```

**Requirements:** `compile wasm` and `parse` are fully offline. `compile ts` needs Node.js.

### Deploy & Transport

```bash
# Deploy source files
vsp deploy zcl_test.clas.abap '$TMP'
vsp deploy zreport.prog.abap '$TMP' --transport A4HK900001

# Transport management
vsp transport list
vsp transport list --user DEVELOPER
vsp transport get A4HK900001

# Install components
vsp install zadt-vsp          # deploy ZADT_VSP service classes
vsp install abapgit           # deploy abapGit
vsp install list              # check what's installed
```

**Requirements:** Standard ADT for deploy. `install` creates objects in `$TMP` or specified package.

### Lua Scripting

```bash
# Interactive REPL
vsp -s dev lua
# lua> objs = searchObject("ZCL_*")
# lua> rows = query("SELECT MANDT, MTEXT FROM T000")
# lua> issues = lint(getSource("CLAS", "ZCL_TEST"))

# Run scripts
vsp -s dev lua examples/scripts/package-audit.lua
vsp -s dev lua examples/scripts/table-explorer.lua
vsp -s dev lua examples/scripts/dependency-check.lua
```

50+ Lua functions: search, query, grep, source, lint, parse, context, debug (breakpoints, stepping, recording, replay), diagnostics (dumps, traces), JSON handling.

**Requirements:** Standard ADT for SAP access. `lint()` and `parse()` work offline within scripts.

### YAML Workflows

```bash
# CI pipeline: discover → syntax check → test → fail on errors
vsp -s dev workflow run examples/workflows/ci-pipeline.yaml

# Quality gate with variables
vsp -s dev workflow run examples/workflows/quality-gate.yaml --var PACKAGE='$ZADT_VSP'
```

9 built-in actions: `search`, `test`, `syntax_check`, `transform`, `save`, `activate`, `print`, `fail_if`, `foreach`. Variables, step chaining, conditional execution, error handling.

**Requirements:** Standard ADT.

### Execute ABAP

```bash
# Run code on SAP
vsp execute "WRITE sy-datum."
vsp execute --file script.abap
echo "WRITE 'hello'." | vsp execute --stdin
```

**Requirements:** Write permissions. Uses ExecuteABAP (unit test wrapper).
If blocked: `vsp install zadt-vsp` for WebSocket-based execution.

### Export & Import

```bash
# Export packages to ZIP (abapGit format)
vsp export '$ZPACKAGE' -o backup.zip
vsp export '$ZORK' '$ZLLM' -o combined.zip --subpackages

# Import from ZIP
vsp copy backup.zip '$TMP'
```

**Requirements:** Export needs ZADT_VSP WebSocket. Standard ADT for import via deploy.

## Feature Requirements Matrix

| Command | Standard ADT | ZADT_VSP | Node.js | Offline |
|---------|:---:|:---:|:---:|:---:|
| `source read/write/edit` | ✅ | — | — | — |
| `context` (+ `--depth`) | ✅ | — | — | — |
| `graph` | ✅ | — | — | — |
| `deps` | ✅ | — | — | — |
| `search` | ✅ | — | — | — |
| `query` | ✅ | — | — | — |
| `grep` | ✅ | — | — | — |
| `system info` | ✅ | — | — | — |
| `test` | ✅ | — | — | — |
| `atc` | ✅ | — | — | — |
| `deploy` | ✅ | — | — | — |
| `transport` | ✅ | — | — | — |
| `lint` | — | — | — | ✅ |
| `parse` | — | — | — | ✅ |
| `compile wasm` | — | — | — | ✅ |
| `compile ts` | — | — | ✅ | — |
| `lua` (REPL/scripts) | ✅ | — | — | — |
| `workflow` | ✅ | — | — | — |
| `execute` | ✅ | optional | — | — |
| `export` | — | ✅ | — | — |
| `install` | ✅ | — | — | — |

**Legend:**
- **Standard ADT** — works with any SAP system that has ADT enabled (default since 7.50)
- **ZADT_VSP** — enhanced features via `vsp install zadt-vsp` (WebSocket, RFC, Git export)
- **Node.js** — required for TypeScript parsing only
- **Offline** — no SAP connection needed at all

## Fallback Behavior

vsp is designed to work with what's available:

1. **No SAP connection?** → `lint`, `parse`, `compile wasm` work fully offline
2. **Standard ADT only?** → `source`, `search`, `query`, `grep`, `graph`, `deps`, `test`, `atc`, `deploy` all work
3. **ZADT_VSP installed?** → `export`, `execute` (via WebSocket), `debug` (via RFC) become available
4. **Missing component?** → Clear error messages tell you what to install and how
5. **ADT call graph unavailable?** → `graph` falls back to WBCROSSGT/CROSS tables automatically

```
$ vsp execute "WRITE 'hello'."
Error: ExecuteABAP requires write permissions.
Check --read-only and --allowed-ops settings.

$ vsp export '$TMP'
Error: WebSocket connect failed.
Ensure ZADT_VSP is deployed: vsp install zadt-vsp
```

## Multi-System Profiles

Save system configs in `.vsp.json`:

```json
{
  "systems": {
    "dev": {
      "url": "https://dev-host:44300",
      "user": "DEVELOPER",
      "client": "001"
    },
    "prod": {
      "url": "https://prod-host:44300",
      "user": "READER",
      "client": "100"
    }
  }
}
```

```bash
vsp -s dev query T000
vsp -s prod search "ZCL_*"
vsp -s dev deploy myclass.clas.abap '$TMP'
```

Passwords via env vars: `VSP_DEV_PASSWORD`, `VSP_PROD_PASSWORD`.

## Pipeline Integration

```bash
# CI/CD: test all custom code
vsp -s dev test --package '$ZCUSTOM' || exit 1

# Lint local files before commit
find src/ -name "*.abap" -exec vsp lint --file {} \;

# Export for backup
vsp -s prod export '$ZPRODUCTION' -o "backup-$(date +%F).zip"

# Compile WASM and deploy
vsp compile wasm calculator.wasm -o ./build/
vsp -s dev deploy ./build/zcl_wasm_calculator.clas.abap '$TMP'

# Query and filter
vsp -s dev query TADIR --where "DEVCLASS = '\$TMP'" --top 50 | grep CLAS

# Check transport readiness before release
vsp -s dev deps '$ZFINANCE' --format summary

# Graph: who uses our interface?
vsp -s dev graph INTF ZIF_ORDER_SERVICE --direction callers
```
