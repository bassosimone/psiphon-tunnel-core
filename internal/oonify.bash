#!/bin/bash
set -e
if [ $# -ne 2 ]; then
  echo "usage: $0 from to" 1>&2
  exit 1
fi
set -x
find . -type f -name \*.go -exec sed -i "s@github.com/$1/psiphon-tunnel-core@github.com/$2/psiphon-tunnel-core@" {} \;
git commit -am "chore: rename github.com/{$1 => $2}/psiphon-tunnel-core"
