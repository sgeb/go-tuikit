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
    echo "done"
fi

if [ ! -d $GO_DIST_DIR ]; then
    echo -n "Downloading and preparing go... "
    (
        mkdir $GO_DIST_DIR
        cd $GO_DIST_DIR
        GO_ARCHIVE=go1.2.2.linux-amd64.tar.gz
        wget -q https://storage.googleapis.com/golang/$GO_ARCHIVE
        tar xzf $GO_ARCHIVE
    )
    echo "done"
fi

echo -n "Installing dependencies... "
go get -u "github.com/nsf/termbox-go"
go get -u "github.com/nsf/tulib"
go get -u "github.com/sgeb/go-sglog"
go get -u "github.com/cloudfoundry/gosigar"
echo "done"

echo -n "Symlinking repo into GOPATH... "
LINK_DEST=$GOPATH/src/github.com/sgeb/go-tuikit
rm -rf $LINK_DEST
mkdir -p $(dirname $LINK_DEST)
ln -s $(pwd) $LINK_DEST
echo "done"
