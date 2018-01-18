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
if [[ "$DOCKER_DEPLOY" == "true" ]]; then 
    # configure docker creds
    retry 3 docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
    # docker tag and push git branch to dockerhub
    if [ -n "$1" ]; then
        retry 3 docker push "mozilla/sops:$1" ||
            (echo "Couldn't push mozilla/sops:$1" && false)
        echo "Pushed mozilla/sops:$1"
    fi
fi
