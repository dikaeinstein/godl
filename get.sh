#!/bin/sh

uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    aarch64) arch="arm64" ;;
  esac
  echo ${arch}
}

ARCH=$(uname_arch)
OS="darwin"
VERSION=$1
# $VERSION with v prefix stripped
VERSION_STRIPPED="${VERSION:1}"

echo "Downloading Godl CLI release ${VERSION} for ${OS}_${ARCH} ..."
echo ""

RELEASE_URL="https://github.com/dikaeinstein/godl/releases/download/${VERSION}/godl_${VERSION_STRIPPED}_${OS}_${ARCH}.tar.gz"
code=$(curl -w '%{http_code}' -L $RELEASE_URL -o /tmp/godl.tar.gz)

if [ $code != "200" ]; then
  echo ""
  echo "[error] Failed to download Godl release ${VERSION} for $OS $ARCH."
  echo "Received HTTP status code $code"
  echo ""
  echo "Supported version of the Godl CLI is:"
  echo " - darwin_amd64"
  echo " - darwin_arm64"
  echo ""
  exit 1
fi

tar -xzf /tmp/godl.tar.gz -C /tmp
sudo chmod +x /tmp/godl
sudo mv /tmp/godl /usr/local/bin/

echo ""
echo "Godl CLI ${VERSION} for ${OS}_${ARCH} installed."
