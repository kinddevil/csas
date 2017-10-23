#!/bin/bash

function err()
{
    echo "\033[31m $1 \033[0m"
    exit 1
}

basepath=$(cd `dirname $0`; pwd)
GOSRC=$basepath/src
SRCDIR=$basepath/..
mkdir -p $GOSRC || err "make directory of go source error"
export GOPATH=$basepath
cp -R $SRCDIR/* $GOSRC/
pushd $GOSRC/ads/
go get || err "go get error" 
go build main.go || err "go build error"
popd
