#!/usr/bin/env bash
set -e

# TODO ARM support. figure out mapping to uname -m output.

tag="${1}"

if [ "${tag}" == "" ]; then
  echo "tag argument required"
  exit 1
fi

rm -rf dist
GOOS=darwin GOARCH=amd64 go build -o "dist/darwin-x86_64-${tag}"
GOOS=linux GOARCH=386 go build -o "dist/linux-i386-${tag}"
GOOS=linux GOARCH=amd64 go build -o "dist/linux-x86_64-${tag}"
GOOS=windows GOARCH=386 go build -o "dist/windows-i386-${tag}"
GOOS=windows GOARCH=amd64 go build -o "dist/windows-x86_64-${tag}"

gh release create $tag ./dist/* --title="${tag}" --notes "${tag}"
