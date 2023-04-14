#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DEFAULT_TAG=${DEFAULT_TAG:-latest}

CONTAINER_CLI=${CONTAINER_CLI:-docker}
CONTAINER_BUILDER=${CONTAINER_BUILDER:-build}

${CONTAINER_CLI} "${CONTAINER_BUILDER}" \
  -f build/Dockerfile \
  -t kube-sidecar:"${DEFAULT_TAG}" .
