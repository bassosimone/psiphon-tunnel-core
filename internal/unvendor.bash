#!/bin/bash
set -e
if [ $# -ne 1 ]; then
  echo "usage: $0 package" 1>&2
  exit 1
fi
if [ ! -d internal/$1 ]; then
  echo "fatal: missing internal/$1" 1>&2
  exit 1
fi
set -x
git rm -rf internal/$1
find . -type f -name \*.go -exec sed -i "s@github.com/bassosimone/psiphon-tunnel-core/internal/$1@$1@" {} \;
go mod tidy
git commit -am "chore: stop vendoring $1 in internal"
