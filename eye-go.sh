#!/bin/bash
# eye script to be executed on file changes
# I.e. `eye -r -m '\.go$ -c './eye-dev.sh'`
# See: https://github.com/callerobertsson/eye

# Execute Golang toolchain commands
# Preferably used as command when running eye file watcher

BACK="\033[0m"
FAIL="\033[31m"
PASS="\033[32m"
INFO="\033[33m"

function execOrDie {
    local step=${@:1:1}
    local cmd=${@:2}

    $cmd

    if [ $? -ne 0 ]; then
        echo -ne '\007'
        printf "\n${FAIL}$step: FAIL${BACK}\n"
        exit
    fi

    printf "${PASS}$step: PASS${BACK}\n"
}

function execIgnore {
    local step=${@:1:1}
    local cmd=${@:2}

    $cmd
    printf "${INFO}$step DONE${BACK}"
}

clear
printf "${INFO}Golang Tool Chain${BACK}\n"
printf "Directory: %s\n\n" `pwd`
execOrDie "Go Build" go build -o /dev/null
execOrDie "Go Test" go test ./...
execOrDie "Go Vet" go tool vet .
execOrDie "Go Lint" golint .
execIgnore "Ack TODO" ack --go TODO
