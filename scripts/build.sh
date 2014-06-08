#!/usr/bin/env bash

set -x
set -e

if [ "$CI" != "true" ]; then
    echo "Run only in CI"
    exit 1
fi

source $(dirname $0)/env.sh

echo "Go version: $(go version)"

OK=1
for p in examples/databinding examples/widgets; do
    echo -n "Building ${p}... "
    (
        cd $p
        go build || OK=1
    )
    echo "done"
done

if [ "$OK" == "1" ]; then
    exit 1
fi
