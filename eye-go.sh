#!/bin/bash

function execOrDie {
    echo text $0
    local args=$@
    echo cmd $@:1
    local status=$?
    if [ $status -ne 0 ]; then
        echo "error with $1" >&2
        exit
    fi
    return $status
}

echo directory: `pwd`
execOrDie "Go Build" go build -o /dev/null
execOrDie "Go Lint" golint
execOrDie "Go Vet" go vet
echo Done

