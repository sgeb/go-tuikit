#!/usr/bin/env bash

set -e

if [ "$CI" == "true" ]; then
    set -x
    source $(dirname $0)/env.sh
fi

echo "Go version: $(go version)"

FAILED=0
for p in examples/databinding examples/widgets; do
    echo -n "Building ${p}... "
    cd $p
    go build || FAILED=1
    cd -
    echo "done"
done

if [ "$FAILED" == "1" ]; then
    exit 1
fi
