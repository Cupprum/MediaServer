#!/bin/bash

set -e
set -u
set -o pipefail

case ${1:-} in
  build)
    docker build -t gemini-jail .
    ;;
  run)
    cd ..
    docker run -it\
      --rm \
      --userns=keep-id \
      --security-opt label=disable \
      -v "$(pwd):/home/gemini/workspace" \
      -v "$HOME/.gemini:/home/gemini/.gemini" \
      gemini-jail bash
    ;;
  *)
    echo "Usage: ./docker.sh [build|run]"
    exit 1
    ;;
esac
