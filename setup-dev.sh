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

if ! [ -x "$(command -v systemctl)" ]; then
  echo 'Error: systemctl is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v distrobox)" ]; then
  echo 'Error: distrobox is not installed (required by apx).' >&2
  exit 1
fi

# Compile almost exit if failed
go build -o almost main.go 
if [ $? -ne 0 ]; then
  echo 'Error: almost build failed.' >&2
  exit 1
fi

# Install binary
chmod +x almost
if [ -f /usr/bin/almost ]; then
  rm /usr/bin/almost
fi
ln -s $(pwd)/almost /usr/bin/almost

# Systemd unit
sudo cp systemd/almost.service /usr/lib/systemd/system/almost.service
sudo systemctl daemon-reload
sudo systemctl enable almost.service
sudo systemctl start almost.service

# Done
printf "Almost installed successfully!"
