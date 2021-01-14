#!/bin/bash

TMPDIR=$1

function prompt {
    printf '\e[01;32mbarf\e[00m:\e[01;34m/home/barf\e[00m$ '
    sleep 0.1

    printf "$1\n"
}

prompt 'ls -sh -1 ~/from'
ls -sh -1 $TMPDIR/from

sleep 1
echo ""

prompt 'barf copy ~/from/* ~/to/'
timeout 4 ./barf.sh -w 132 copy $TMPDIR/from/* $TMPDIR/to/
echo "^C"

sleep 1
echo ""

prompt 'barf monitor'
./barf.sh -w 132 monitor
echo ""

sleep 1
echo ""

prompt 'ls -sh -1 ~/to'
ls -sh -1 $TMPDIR/to

echo ""

prompt ''

sleep 5
echo ""
