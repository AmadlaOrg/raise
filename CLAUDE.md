# CLAUDE.md

## AI Skills

Follow the practices defined in `~/Projects/SiteNetSoft/ai-skills/`:
- `dev-practices/golang/` — Go style, error handling, functions, testing, linting
- `dev-practices/git/` — Git authorship rules, multi-repo workspace patterns

## Project Overview

Raise is the first step in the Amadla tool pipeline (`raise` -> `lay` -> `enjoin` -> `weaver` -> `waiter`). It provisions infrastructure (VMs, cloud instances) by discovering `raise-*` plugins on PATH and delegating operations (up, halt, destroy, ssh, status) to them. Each plugin handles a specific provider (libvirt, virtualbox, aws, hetzner, etc.).

**Entity type (core):** `amadla.org/entity/infrastructure@^v1.0.0` (schemas in `Entities/Infrastructure/`). Schema URN IDs use no "Entity" prefix. Raise core only declares its own entity type; plugins declare their own sub-types independently.

## Build Commands

```bash
make build    # Build for current platform
make test     # Run tests
make clean    # Remove build artifacts
```

## Architecture

**UNIX Plugin Protocol:**
```
raise info [-o table|json|yaml]                    ->  outputs raise metadata (supports -o flag)
raise up [name] --provider <provider> -f <file>    ->  raise-<provider> up [name] -f <file>
raise up [name] -f <file>                          ->  reads provider from entity, delegates to raise-<provider>
raise halt [name] --provider <provider>            ->  raise-<provider> halt [name]
raise destroy [name] --provider <provider>         ->  raise-<provider> destroy [name]
raise ssh [name] --provider <provider>             ->  raise-<provider> ssh [name]  (exec replaces process)
raise status [--provider <provider>]               ->  raise-<provider> status (or all plugins)
raise plugins [-o table|json|yaml]                 ->  scans PATH for raise-* binaries
```

**Entity type routing:**
- raise core declares only `amadla.org/entity/infrastructure@^v1.0.0`
- Each plugin declares its own supported entity types (e.g., infrastructure/vm, infrastructure/cloud)
- amadla discovers plugin entity types by calling `<plugin> info` directly
- raise reads `provider` from Infrastructure entity to dispatch to the correct plugin

**Package Structure:**
- `main.go` - CLI entry point (Cobra)
- `entity/` - Entity file parsing (reads provider field from Infrastructure entities)
- `plugin/` - Plugin discovery and execution (PATH scanning, subprocess delegation)
- `cmd/` - CLI commands (info, up, halt, destroy, ssh, status, plugins)
- `cmd/output.go` - Shared output formatting (table/json/yaml via -o flag)

**Plugin Protocol (raise-* binaries):**
- `info` subcommand -> metadata conforming to Tools/Info schema (name, version, engine, supports, description)
- All commands support `-o table|json|yaml` (default: table)
- `up [name] -f <file>` subcommand -> provisions infrastructure
- `halt [name]` subcommand -> stops infrastructure
- `destroy [name]` subcommand -> destroys infrastructure
- `ssh [name]` subcommand -> opens SSH session
- `status` subcommand -> shows infrastructure status
- Exit codes: 0 success, 1 failure, 2 usage error
- Data to stdout, diagnostics to stderr
