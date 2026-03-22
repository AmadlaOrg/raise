# CLAUDE.md

## Project Overview

Raise is the first step in the Amadla tool pipeline (`raise` -> `lay` -> `enjoin` -> `weaver` -> `waiter`). It provisions infrastructure (VMs, cloud instances) by discovering `raise-*` plugins on PATH and delegating operations (up, halt, destroy, ssh, status) to them. Each plugin handles a specific provider (libvirt, virtualbox, aws, hetzner, etc.).

**Entity types:** Infrastructure, Infrastructure/VM (schemas in `Entities/Infrastructure/`). Schema URN IDs use no "Entity" prefix.

## Build Commands

```bash
make build    # Build for current platform
make test     # Run tests
make clean    # Remove build artifacts
```

## Architecture

**UNIX Plugin Protocol:**
```
raise up [name] --from <provider> -f <file>    ->  raise-<provider> up [name] -f <file>
raise halt [name] --from <provider>            ->  raise-<provider> halt [name]
raise destroy [name] --from <provider>         ->  raise-<provider> destroy [name]
raise ssh [name] --from <provider>             ->  raise-<provider> ssh [name]  (exec replaces process)
raise status [--from <provider>]               ->  raise-<provider> status (or all plugins)
raise plugins                                  ->  scans PATH for raise-* binaries
```

**Package Structure:**
- `main.go` - CLI entry point (Cobra)
- `plugin/` - Plugin discovery and execution (PATH scanning, subprocess delegation)
- `cmd/` - CLI commands (up, halt, destroy, ssh, status, plugins)

**Plugin Protocol (raise-* binaries):**
- `info` subcommand -> JSON metadata (name, version, engine, description, supports)
- `up [name] -f <file>` subcommand -> provisions infrastructure
- `halt [name]` subcommand -> stops infrastructure
- `destroy [name]` subcommand -> destroys infrastructure
- `ssh [name]` subcommand -> opens SSH session
- `status` subcommand -> shows infrastructure status
- Exit codes: 0 success, 1 failure, 2 usage error
- Data to stdout, diagnostics to stderr
