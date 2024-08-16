#!/usr/bin/env bash
set -e

echo "Running fuzz tests..."
# manually specify the tests to fuzz since go toolchain doesn't support
# fuzzing multiple packages with multiple fuzz tests
cd x/blob/types && go test -fuzz=FuzzPFBGasEstimation -fuzztime 3m
