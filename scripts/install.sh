#!/bin/bash -e

case $(uname -sm) in
"Linux x86_64") target="linux-amd64" ;;
"Linux armv5tel") target="linux-arm5" ;;
"Linux armv7l") target="linux-arm7" ;;
*) target="unknown" ;;
esac

if [ "$target" == "unknown" ]; then
    echo "Unsupported target platform"
    exit 1
fi

curl --fail --progress-bar -L https://github.com/mattiasrunge/barf/releases/latest/download/barf-$target.tar.gz | tar xvz -C /usr/local/bin

echo "barf was downloaded and installed at /usr/local/bin"

barf -v
