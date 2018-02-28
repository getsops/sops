#!/bin/bash

# THIS IS MEANT TO BE RUN BY CI

set -e
set +x

# Usage: retry MAX CMD...
# Retry CMD up to MAX times. If it fails MAX times, returns failure.
# Example: retry 3 docker push "mozilla/sops:$TAG"
function retry() {
    max=$1
    shift
    count=1
    until "$@"; do
        count=$((count + 1))
        if [[ $count -gt $max ]]; then
            return 1
        fi
        echo "$count / $max"
    done
    return 0
}

# docker tag and push git branch to dockerhub
if [ -n "$1" ]; then
    # configure docker creds
    retry 3  echo "$DOCKER_PASSWORD" | docker login -u="$DOCKER_USERNAME" --password-stdin

    [ "$1" == master ] && TAG=latest || TAG="$1"
    docker tag sops:build "mozilla/sops:$TAG" ||
        (echo "Couldn't tag sops:build as mozilla/sops:$TAG" && false)
    retry 3 docker push "mozilla/sops:$TAG" ||
        (echo "Couldn't push mozilla/sops:$TAG" && false)
    echo "Pushed mozilla/sops:$TAG"
fi
