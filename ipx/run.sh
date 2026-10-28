#!/usr/bin/env bash
# Builds the ipx image (unless IPX_IMAGE is set), starts it in Docker, runs
# the latency benchmark against it, and stops the container afterwards.
# Extra arguments are passed to the benchmark.
#
#   ipx/run.sh                      # defaults
#   ipx/run.sh -n 500 -c 4 -json results.json
#   IPX_VERSION=3.1.2 ipx/run.sh    # benchmark another ipx release
#   IPX_DOCKER_ARGS="--cpus=2" ipx/run.sh
set -euo pipefail

cd "$(dirname "$0")/.."

VERSION="${IPX_VERSION:-4.0.0}"
IMAGE="${IPX_IMAGE:-proxy-benchmark/ipx:${VERSION}}"
PORT="${IPX_PORT:-8080}"
NAME="ipx-latency-bench"

if [ -z "${IPX_IMAGE:-}" ]; then
  docker build -q --build-arg "IPX_VERSION=${VERSION}" -t "$IMAGE" ipx >/dev/null
fi

docker rm -f "$NAME" >/dev/null 2>&1 || true
docker run -d --name "$NAME" \
  -p "127.0.0.1:${PORT}:8080" \
  --add-host=host.docker.internal:host-gateway \
  -e IPX_HTTP_DOMAINS=host.docker.internal \
  ${IPX_DOCKER_ARGS:-} \
  "$IMAGE" >/dev/null
trap 'docker rm -f "$NAME" >/dev/null 2>&1 || true' EXIT

# ipx has no health endpoint; any HTTP response means it is listening.
for _ in $(seq 1 30); do
  curl -s -o /dev/null "http://127.0.0.1:${PORT}/" && break
  sleep 1
done

go run . -proxy ipx -proxy-url "http://127.0.0.1:${PORT}" "$@"
