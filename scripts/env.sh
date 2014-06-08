#!/bin/bash

# This script is to be *sourced* by the other scripts to set up env variables

export GO_DIST_DIR=$HOME/go_dist
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GO_DIST_DIR/go/bin:$PATH
