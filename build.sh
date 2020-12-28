#!/bin/bash -e

function build {
    bash -c 'go build -o "build/$GOOS-$GOARCH$GOARM/rft"'
    bash -c 'pushd build/$GOOS-$GOARCH$GOARM &> /dev/null && tar -czvf ../rtf-$GOOS-$GOARCH$GOARM.tar.gz * && popd &> /dev/null '
}

GOOS=linux GOARCH=arm GOARM=5 build
GOOS=linux GOARCH=arm GOARM=7 build
GOOS=linux GOARCH=amd64 build
