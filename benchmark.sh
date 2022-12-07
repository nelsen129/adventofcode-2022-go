#!/bin/bash

directories=$(find . -maxdepth 1 -type d -regex '\./.*[0-9]' | sort)

for directory in $directories ; do
    echo "$directory"
    cd $directory
    go run . input.txt
    cd ..
done
