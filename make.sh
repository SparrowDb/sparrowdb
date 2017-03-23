#!/bin/bash

if ! [ -x "$(command -v go)" ]; then
    echo "Could not find Golang"
    exit 1    
fi

echo "Generating SparrowDb binaries"
rm -rf dist
mkdir dist

echo "Building ..."
go build -o dist/sparrow .
go build -o dist/commander tools/commander/commander.go
go build -o dist/datafile tools/datafile/datafile.go

echo "Copying ..."
cp -r scripts dist/scripts
cp -r config dist/config

echo "Done !"