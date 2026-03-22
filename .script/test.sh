#!/bin/bash
set -euo pipefail
mkdir -p .reports
go test -coverprofile=.reports/coverage.out ./...
