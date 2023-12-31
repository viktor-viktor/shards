#!/bin/bash
go test -coverprofile=c.out ./...
go tool cover -html=c.out -o=c.html

verbose=false
while getopts 'v' flag; do
    case "${flag}" in
        v) verbose=true
    esac
done

if [ $verbose == 'true' ]; then
  open c.html
fi
