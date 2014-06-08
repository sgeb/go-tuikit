#!/usr/bin/env bash

set -x
set -e

if [ "$CI" != "true" ]; then
    echo "Run only in CI"
    exit 1
fi

source $(dirname $0)/env.sh

if [ "$TUIKIT_CLEAN_HOME" == "true" ]; then
    echo -n "Cleaning ${HOME}... "
    rm -rf $GO_DIST_DIR
    rm -rf $GOPATH
    rm -rf $HOME/dl
    echo "done"
fi

if [ ! -d $GO_DIST_DIR ]; then
    echo -n "Downloading and preparing go... "
    (
        # For debugging
        tree $HOME

        mkdir $GO_DIST_DIR
        cd $GO_DIST_DIR
        GO_ARCHIVE=go1.2.2.linux-amd64.tar.gz
        wget -q https://storage.googleapis.com/golang/$GO_ARCHIVE
        tar xzf $GO_ARCHIVE
    )
    echo "done"
fi
