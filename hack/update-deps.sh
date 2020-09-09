#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

export GO111MODULE=on

# Prune modules.
go mod tidy
go mod vendor

# Clean up the vendor area to remove OWNERS and tests.
rm -rf $(find vendor/ -name 'OWNERS')
rm -rf $(find vendor/ -name '*_test.go')

