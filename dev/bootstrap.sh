#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$DIR/.."

log() {
  echo "$(basename ${BASH_SOURCE[0]}): $@"
}

install_hooks() {
  git config core.hooksPath \
    || git config core.hooksPath ./dev/hooks
}

log 're/configuring hooks...'
install_hooks
