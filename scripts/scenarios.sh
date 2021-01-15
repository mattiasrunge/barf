#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

SDIR="$DIR/scenarios"
OUTDIR="$DIR/../docs/svg"
WIDTH=132
TMPDIR=$(mktemp -d -t barf-XXXXXXXXXX)

SCENARIOS="$SDIR/*.sh"

function gendata {
    mkdir -p $TMPDIR/local

    fallocate -l 1500MB $TMPDIR/local/huge
    fallocate -l 1000MB $TMPDIR/local/big
    fallocate -l 500MB $TMPDIR/local/medium
    fallocate -l 100MB $TMPDIR/local/small
}

function scenario {
    echo "Scenario: $1"

    echo " - Generating data..."
    gendata

    echo " - Running scenario and creating SVG..."
    svg-term --out $OUTDIR/$1.svg --command="bash $SDIR/$1.sh $TMPDIR" --window --height $2 --width $WIDTH --no-cursor

    echo " - Removing data..."
    rm -rf $TMPDIR/*

    echo " - Done"
    echo ""
}

# for f in $SCENARIOS; do
#     filename=$(basename -- "$f")
#     name="${filename%.*}"

#     scenario $name
# done
scenario copy-normal 23
scenario copy-monitor 10
scenario copy-monitor-many 22
scenario copy-remote 16
scenario daemon-journal 14

rm -rf $TMPDIR
