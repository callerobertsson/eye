#!/bin/bash
# eye script to be executed on file changes
# I.e. `eye -r -m '\.go$ -c './eye-dev.sh'`
# See: https://github.com/callerobertsson/eye

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

function execIgnore {
    local step=${@:1:1}
    local cmd=${@:2}

    echo "Executing: $step"
    $cmd
}

echo ""
echo "=== TOOL CHAIN INITIATED ==="
echo "Directory:" `pwd`
execOrDie "Go Build" go build -o /dev/null
execOrDie "Go Test" go test ./...
execOrDie "Go Vet" go tool vet .
execOrDie "Go Lint" golint .
execIgnore "Ack TODO" ack --go TODO
echo "Success"
