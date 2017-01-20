#!/bin/bash

# Execute Golang toolchain commands
# Preferably used as command when running eye file watcher

function execOrDie {
    local step=${@:1:1}
    local cmd=${@:2}

    echo "Executing: $step"
    $cmd
    if [ $? -ne 0 ]; then
        echo ""
        echo "Terminating in $step step"
        echo ""
        echo "^^^^^^^^ FAILURE ^^^^^^^^"
        echo -ne '\007'
        exit
    fi
}

echo ""
echo "=== TOOL CHAIN INITIATED ==="
echo "Directory:" `pwd`
execOrDie "Go Build" go build -o /dev/null
execOrDie "Go Test" go test ./...
execOrDie "Go Vet" go tool vet .
execOrDie "Go Lint" golint .
echo "Success"

