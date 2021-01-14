#!/bin/bash

TMPDIR=$1

source ./scripts/scenarios/lib/helpers.sh

mkdir -p $TMPDIR/remote1
mkdir -p $TMPDIR/remote2
mkdir -p $TMPDIR/remote3

prompt 'barf copy ~/local/* ~/remote1/'
timeout 3 barf -w 132 copy $TMPDIR/local/* $TMPDIR/remote1/
echo "^C"
prompt 'barf copy ~/local/* ~/remote2/'
timeout 3 barf -w 132 copy $TMPDIR/local/* $TMPDIR/remote2/
echo "^C"
prompt 'barf copy ~/local/* ~/remote3/'
timeout 3 barf -w 132 copy $TMPDIR/local/* $TMPDIR/remote3/
echo "^C"

sleep 0.2

prompt 'barf monitor'
barf -w 132 monitor

finish
