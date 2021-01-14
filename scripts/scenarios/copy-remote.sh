#!/bin/bash

TMPDIR=$1

source ./scripts/scenarios/lib/helpers.sh

mkdir -p $TMPDIR/remote

rm -rf $TMPDIR/local/huge $TMPDIR/local/big $TMPDIR/local/medium

prompt 'barf copy ~/local/small barf@barf:~/remote'
barf -w 132 copy $TMPDIR/local/small localhost:/$TMPDIR/remote/

listSizes $TMPDIR/local local
listSizes $TMPDIR/remote remote

finish
