#!/usr/bin/env sh

set -e

export GOARCH=amd64

name=$1
version=$2

for os in linux darwin; do
  export GOOS=$os
  executable="./bin/${name}_v${version}_${GOOS}_${GOARCH}"

  mkdir -p ./artifacts

  go build -o "$executable"
  zip -j -r "./artifacts/${name}_${version}_${GOOS}_${GOARCH}.zip" "$executable"
done