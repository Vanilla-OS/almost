#!/bin/sh

# sudo check
if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root" 1>&2
    exit 1
fi

# Essential tools
if ! [ -x "$(command -v go)" ]; then
  echo 'Error: go is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v apx)" ]; then
  echo 'Error: apx is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v distrobox)" ]; then
  echo 'Error: distrobox is not installed (required by apx).' >&2
  exit 1
fi

# Compile almost exit if failed
cd almost
go build -o almost main.go 
if [ $? -ne 0 ]; then
  echo 'Error: almost build failed.' >&2
  exit 1
fi

# Install binary
chmod +x almost
ln -s $(pwd)/almost /usr/local/bin/almost

# Done
printf "Almost installed successfully!"
