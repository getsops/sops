#!/usr/bin/env bash

set -e
echo "" > coverage.txt

failed=0

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d && true
    rc=$?
    if [ $rc != 0 ]; then
      failed=$rc
    fi
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

exit ${failed}
