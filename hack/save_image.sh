#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

TAG=${TAG:-latest}
OUTPUT=${OUTPUT:-_output}
TYPE=${TYPE:-tar}
IMAGE=${IMAGE:-kube-sidecar}

# support other container tools. e.g. podman
CONTAINER_CLI=${CONTAINER_CLI:-docker}
CONTAINER_BUILDER=${CONTAINER_BUILDER:-build}
CONTAINER_SAVE=${CONTAINER_SAVE:-save}

mkdir -p "${OUTPUT}"
# shellcheck disable=SC2086
${CONTAINER_CLI} "${CONTAINER_SAVE}" ${IMAGE}:${TAG} --output "${OUTPUT}"/"${IMAGE}".${TYPE}