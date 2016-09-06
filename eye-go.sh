#!/bin/bash
# eye script to be executed on file changes
# I.e. `eye -r -m '\.go$ -c './eye-dev.sh'`
# See: https://github.com/callerobertsson/eye

# Execute command
# First arg is a label, the rest will be executed
function execOrDie {
    echo "${@:1:1}..."
    res=$((${@:2}) 2>&1) # redir stderr to stdout and store in res

    if [ $? -ne 0 ]; then
        echo ""
        echo "---> ERROR!"
        echo "Command: $1" 
        echo "Output"
        echo "$res"
        exit
    fi
}

# Steps to execute
execOrDie "Go Build" go build -o /dev/null
execOrDie "Go Lint" golint
execOrDie "Go Vet" go vet

echo "...done"

