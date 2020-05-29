#!/usr/bin/env bash
set -e

# WARN: non hermetic build (people must run this script inside docker to
# produce deterministic binaries).

CONTRACT_DIR=./build/contract/
GIT_RELEASE="https://github.com/bcbchain/contract/releases/download/"
CONTRACT_GIT="https://api.github.com/repos/bcbchain/contract/tags"
CONTRACT="genesis-smart-contract"

i=1
for _ in $(cat bcb.mod)
do
  NUM=$i
  TAG=$(awk 'NR=='$NUM' {print $1}' bcb.mod)
  VER=$(awk 'NR=='$NUM' {print $2}' bcb.mod)

  if [ "$TAG" == "$CONTRACT" ];then
    CONTRACT_LATEST_TAG=$VER
  fi

  if [ -n "$CONTRACT_LATEST_TAG" ];then
    echo "==> Downloading contract to ${CONTRACT_DIR}..."
    rm -rf "$CONTRACT_DIR"
    mkdir -p "$CONTRACT_DIR"
    pushd "$CONTRACT_DIR" >/dev/null

    curl -OL "${GIT_RELEASE}${CONTRACT_LATEST_TAG}/${CONTRACT}""_${CONTRACT_LATEST_TAG}.tar.gz"

    tar -zxf "${CONTRACT}""_${CONTRACT_LATEST_TAG}.tar.gz"
    rm -f "${CONTRACT}""_${CONTRACT_LATEST_TAG}.tar.gz"
    popd >/dev/null
  fi
  : $(( i++ ))
done

exit 0