#!/bin/bash -e

V=$1

function build {
    echo "Building GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM"
    bash -c "go build -ldflags=\"-X 'rft/internal/config.Version=v$V' -X 'rft/internal/config.BuildTime=$(date)' -X 'rft/internal/config.BuildChecksum=$(git rev-parse HEAD)'\" -o \"build/$GOOS-$GOARCH$GOARM/rft\" cmd/rft/main.go"
    bash -c 'pushd build/$GOOS-$GOARCH$GOARM &> /dev/null && tar -czvf ../rtf-$GOOS-$GOARCH$GOARM.tar.gz * && popd &> /dev/null '
}

GOOS=linux GOARCH=arm GOARM=5 build
GOOS=linux GOARCH=arm GOARM=7 build
GOOS=linux GOARCH=amd64 build
