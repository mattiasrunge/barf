#!/bin/bash

TMPDIR=$1

source ./scripts/scenarios/lib/helpers.sh

mkdir -p $TMPDIR/remote

prompt 'barf copy ~/local/* ~/remote/'
timeout 4 barf copy $TMPDIR/local/* $TMPDIR/remote/
echo "^C"

sleep 0.2

prompt 'barf monitor'
barf monitor

finish
