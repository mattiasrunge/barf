#!/bin/bash

TMPDIR=$1

source ./scripts/scenarios/lib/helpers.sh

mkdir -p $TMPDIR/remote

prompt 'barf copy ~/local/* ~/remote/'
timeout 4 barf -w 132 copy $TMPDIR/local/* $TMPDIR/remote/
echo "^C"

prompt 'cat ~/.config/barf/journal/active/*.json | jq'
cat ~/.config/barf/journal/active/*.json | sed "s|/*$TMPDIR|~|g" | jq

finish
