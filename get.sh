#!/bin/sh

OS="darwin"
ARCH="amd64"
VERSION=$1
# $VERSION with v prefix stripped
VERSION_STRIPPED="${VERSION:1}"

echo "Downloading Godl CLI release ${VERSION} for ${OS}_${ARCH} ..."
echo ""

curl --fail -L "https://github.com/dikaeinstein/godl/releases/download/${VERSION}/godl_${VERSION_STRIPPED}_${OS}_${ARCH}.tar.gz" -o /tmp/godl.tar.gz
if ! [ $? -eq 0 ]; then
  echo ""
  echo "[error] Failed to download Godl release for $OS $ARCH."
  echo ""
  echo "Supported version of the Godl CLI is:"
  echo " - darwin_amd64"
  echo ""
  exit 1
fi

tar -xzf /tmp/godl.tar.gz -C /tmp
sudo chmod +x /tmp/godl
sudo mv /tmp/godl /usr/local/bin/

echo ""
echo "Godl CLI ${VERSION} for ${OS}_${ARCH} installed."
