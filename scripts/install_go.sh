#!/usr/bin/env bash

set -x
set -e

if [ "$CI" != "true" ]; then
    echo "Run only in CI"
    exit 1
fi

source $(dirname $0)/env.sh

echo -n "Downloading and preparing go... "
(
    # For debugging
    pwd
    ls -l
    ls -l $HOME

    GO_ARCHIVE=go1.2.2.linux-amd64.tar.gz
    mkdir $GO_DIST_DIR
    cd $GO_DIST_DIR
    wget -q https://storage.googleapis.com/golang/$GO_ARCHIVE
    tar xzf $GO_ARCHIVE
)
