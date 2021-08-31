#!/bin/bash
set -e
if [ $# -ne 1 ]; then
  echo "usage: $0 package" 1>&2
  exit 1
fi
if [ ! -d vendor/$1 ]; then
  echo "fatal: missing vendor/$1" 1>&2
  exit 1
fi
set -x
mkdir -p internal/$(dirname $1)
git mv vendor/$1 internal/$1
find . -type f -name \*.go -exec sed -i "s@$1@github.com/Psiphon-Labs/psiphon-tunnel-core/internal/$1@" {} \;
git commit -am "chore: move $1 to internal"
