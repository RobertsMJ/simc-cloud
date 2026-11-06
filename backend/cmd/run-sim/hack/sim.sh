#!/bin/sh
tmp="${TMPDIR:-/tmp}"
exec docker run --rm --init --user "$(id -u):$(id -g)" -v "$tmp:$tmp" \
  simulationcraftorg/simc:latest "$@"
