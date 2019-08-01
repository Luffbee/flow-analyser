#!/bin/zsh

ADDR=$1
ID=$2
DUR=$3

if [ "${ADDR}x" = "x" ]; then
  echo "please give the address"
  exit 1
fi

if [ "${ID}x" = "x" ]; then
  echo "please give the region id"
  exit 1
fi

if [ "${DUR}x" = "x" ]; then
  DUR="10m"
fi

curl -gsSL "http://${ADDR}/api/v1/query?query=pd_hotcache_region_flow{type=%22write-bytes%22,%20name%20=%22region-${ID}%22}[${DUR}]" -o region-${ID}.json
