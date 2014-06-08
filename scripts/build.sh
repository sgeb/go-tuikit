#!/usr/bin/env bash

set -e

if [ "$CI" == "true" ]; then
    set -x
    source $(dirname $0)/env.sh
fi

echo "Go version: $(go version)"

FAILED_ALL=0
for p in examples/databinding examples/widgets examples/tktop; do
    FAILED=0
    echo -n "=> Building ${p}... "
    cd $p
    go build || FAILED=1
    cd -

    if [ "$FAILED" == "1" ]; then
        FAILED_ALL=1
        echo "FAILED"
    else
        echo "done"
    fi
done

if [ "$FAILED_ALL" == "1" ]; then
    exit 1
fi
