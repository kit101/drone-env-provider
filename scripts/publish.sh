#!/bin/sh

registry=${REGISTRY:-docker.io}
image=${IMAGE:-kit101z/drone-ext-envs}

tag="$1"

ocifile=".oci.`date +%s`"


noclean="0"
if printf "%s\n" "$@" | grep -q -- '--no-clean'; then
  noclean="1"
fi
skopeo="0"
if printf "%s\n" "$@" | grep -q -- '--skopeo'; then
  skopeo="1"
fi
push="0"
if printf "%s\n" "$@" | grep -q -- '--push'; then
  push="1"
fi

run() {
  local cmd="$1"
  echo "+ $cmd"
  eval "$cmd"
  echo ""
}


build_cmd="docker buildx b -t $image:$tag --platform linux/amd64,linux/arm64 --no-cache ."

#set -xe

if [ "$skopeo" == "1" ]; then
  final_build_cmd=$( [ "$push" == "1" ] && echo "$build_cmd -o type=oci,dest=$ocifile,tar=false" || echo "$build_cmd" )
  run "$final_build_cmd"
  if [ "$push" == "1" ]; then
    run "skopeo copy oci:$ocifile docker://$registry/$image:$tag --all"
  fi
  if [ "$noclean" == "0" ]; then
    run "rm -rf .oci.*"
  fi
else
  final_build_cmd=$( [ "$push" == "1" ] && echo "$build_cmd --push" || echo "$build_cmd" )
  run "$final_build_cmd"
fi
