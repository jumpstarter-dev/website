#!/bin/sh

set -ex

ARCH=$(uname -m)

case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

URL=$(curl https://api.github.com/repos/gohugoio/hugo/releases/latest -s | \
  jq -r ".assets[] | select(.name | test(\"linux\")) | .browser_download_url" | \
  grep .tar.gz | \
  grep extended | \
  grep "${ARCH}")

echo $URL
curl -L $URL | tar xvfz - -C /usr/bin hugo

