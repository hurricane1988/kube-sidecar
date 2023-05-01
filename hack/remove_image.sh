#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

IMAGE_TAR="kube-sidecar.tar"
OUTPUT="_output"

if [ -f ${OUTPUT}/${IMAGE_TAR} ] ; then
  rm -f ${OUTPUT}/${IMAGE_TAR}
fi