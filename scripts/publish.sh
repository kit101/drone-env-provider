#!/bin/sh

registry=${REGISTRY:-docker.io}
image=${IMAGE:-kit101z/drone-ext-envs}

tag="$1"

ocifile=".oci.`date +%s`"


clean="1"
if printf "%s\n" "$@" | grep -q -- '--no-clean'; then
  clean="0"
fi

set -ex


docker buildx b -t $image:$tag --platform linux/amd64,linux/arm64 -o type=oci,dest=$ocifile,tar=false . --no-cache

skopeo copy oci:$ocifile docker://$registry/$image:$tag --all

if [ "$clean" == "1" ]; then
  rm -rf .oci.*
fi