#!/bin/bash

function mkdata {
    mkdir -p $1

    fallocate -l 2GB $1/01_very_big_file.bin
    fallocate -l 1GB $1/02_big_file.bin
    fallocate -l 500MB $1/03_medium_file.bin
    fallocate -l 100MB $1/04_small_file.bin
}

function prompt {
    printf '\n\e[01;32mbarf\e[00m:\e[01;34m/home/barf\e[00m$ '
    sleep 0.2

    printf "$1\n"
}

function listSizes {
    prompt "du -ach * ~/$2"
    pushd $1 &> /dev/null
    du -ach *
    popd &> /dev/null
}

function finish {
    prompt ''
    sleep 6
    echo ""
}

function barf {
    ./barf.sh "$@"
}
