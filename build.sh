#!/bin/bash

for GOOS in linux; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o bin/spegel-$GOOS-$GOARCH
  done
done
