#!/bin/bash

TMPDIR=$1

source ./scripts/scenarios/lib/helpers.sh

mkdir -p $TMPDIR/remote

rm -rf $TMPDIR/local/huge $TMPDIR/local/big

prompt 'barf copy ~/local/* barf@barf:~/remote'
./barf.sh -w 132 copy $TMPDIR/local/* localhost:/$TMPDIR/remote/

listSizes $TMPDIR/local local
listSizes $TMPDIR/remote remote

finish
