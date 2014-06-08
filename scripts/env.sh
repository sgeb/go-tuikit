#!/bin/bash

# This script is to be *sourced* by the other scripts to set up env variables

export GO_DIST_DIR=$HOME/go_dist
export GOROOT=$GO_DIST_DIR/go/
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
