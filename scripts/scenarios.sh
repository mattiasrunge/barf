#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

SDIR="$DIR/scenarios"
OUTDIR="$DIR/../docs/svg"
WIDTH=132
TMPDIR=$(mktemp -d -t barf-XXXXXXXXXX)

SCENARIOS="$SDIR/*.sh"


function gendata {
    mkdir -p $TMPDIR/from
    mkdir -p $TMPDIR/to

    fallocate -l 2GB $TMPDIR/from/01_very_big_file.bin
    fallocate -l 1GB $TMPDIR/from/02_big_file.bin
    fallocate -l 500MB $TMPDIR/from/03_medium_file.bin
    fallocate -l 100MB $TMPDIR/from/04_small_file.bin
}

function scenario {
    echo "Scenario: $1"

    echo " - Generating data..."
    gendata

    echo " - Running script and creating SVG..."
    svg-term --out $OUTDIR/$1.svg --command="bash $SDIR/$1.sh $TMPDIR" --window --height 20 --width $WIDTH --no-cursor

    echo " - Removing data..."
    rm -rf $TMPDIR/*

    echo " - Done"
    echo ""
}

for f in $SCENARIOS; do
    filename=$(basename -- "$f")
    name="${filename%.*}"

    scenario $name
done

rm -rf $TMPDIR
