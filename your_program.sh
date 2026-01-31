#!/bin/sh
# Build and run gosh locally.

set -e

cd "$(dirname "$0")"
go build -o gosh app/*.go
exec ./gosh "$@"
