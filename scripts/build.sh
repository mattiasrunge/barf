#!/bin/bash -e

V=$1

function build {
    echo "Building barf ${V} with: GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM"
    bash -c "\
        go build -a \
        -ldflags=\"\
            -s -w \
            -X 'barf/internal/config.production=yes' \
            -X 'barf/internal/config.Version=$V' \
            -X 'barf/internal/config.BuildTime=$(date)' \
            -X 'barf/internal/config.BuildChecksum=$(git rev-parse HEAD)'\
        \" \
        -o \"../build/$GOOS-$GOARCH$GOARM/barf\" ../cmd/barf/main.go"
    bash -c "\
        pushd ../build/$GOOS-$GOARCH$GOARM &> /dev/null && \
        tar -czf ../barf-$GOOS-$GOARCH$GOARM.tar.gz * && \
        popd &> /dev/null\
    "
    SIZE=$(stat -c%s "../build/barf-$GOOS-$GOARCH$GOARM.tar.gz")
    echo "barf-$GOOS-$GOARCH$GOARM.tar.gz ($SIZE bytes) complete!"
    echo ""
}

GOOS=linux GOARCH=arm GOARM=5 build
GOOS=linux GOARCH=arm GOARM=7 build
GOOS=linux GOARCH=amd64 build
