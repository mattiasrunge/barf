#!/bin/bash

TMPDIR=$1

source ./scripts/scenarios/lib/helpers.sh

mkdir -p $TMPDIR/remote

prompt 'barf copy ~/local/* ~/remote/'
barf -w 132 copy $TMPDIR/local/* $TMPDIR/remote/

listSizes $TMPDIR/local local
listSizes $TMPDIR/remote remote

prompt 'barf copy ~/local/* ~/remote/'
barf -w 132 copy $TMPDIR/local/* $TMPDIR/remote/

finish
