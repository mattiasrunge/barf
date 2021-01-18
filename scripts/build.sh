#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
V=$1

set -x

function build {
    NAME="$GOOS-$GOARCH$GOARM"
    echo "Building barf ${V} with: GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM"
    bash -c "\
        go build -a \
        -ldflags=\"\
            -s -w \
            -X 'barf/internal/config.production=yes' \
            -X 'barf/internal/config.Version=$V' \
            -X 'barf/internal/config.BuildName=$NAME' \
            -X 'barf/internal/config.BuildTime=$(date)' \
            -X 'barf/internal/config.BuildChecksum=$(git rev-parse HEAD)'\
        \" \
        -o \"$DIR/../build/$NAME/barf\" $DIR/../cmd/barf/main.go"
    bash -c "\
        pushd $DIR/../build/$NAME &> /dev/null && \
        tar -czf ../barf-$NAME.tar.gz * && \
        popd &> /dev/null\
    "
    echo "$DIR/../build/barf-$NAME.tar.gz"
    SIZE=$(stat -c%s "$DIR/../build/barf-$NAME.tar.gz")
    echo "barf-$NAME.tar.gz ($SIZE bytes) complete!"
    echo ""
}

GOOS=linux GOARCH=arm GOARM=5 build
GOOS=linux GOARCH=arm GOARM=7 build
GOOS=linux GOARCH=amd64 build
