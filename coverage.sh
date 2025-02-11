#!/usr/bin/env bash
# Copyright (c) 2023 FajarLaksono. All Rights Reserved.

set -x
set -e
set -o pipefail

rm coverage.txt || true 2> /dev/null

go install github.com/axw/gocov/gocov@v1.1.0

IS_TESTING_ON_CICD=true go test -v -p=1 -coverprofile=profile.out ./... | tee -a >(go-junit-report > test.xml) coverage.txt
gocov convert profile.out | gocov-xml > coverage.xml
