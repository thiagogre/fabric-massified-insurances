#!/bin/bash

SOURCE=${1:-"internal/domain"}
DESTINATION=${2:-"internal/domain/mocks"}

mkdir -p $DESTINATION

for file in $SOURCE/*.go; do
    base_name=$(basename $file .go)

    mockgen -source=$file -destination="${DESTINATION}/${base_name}_mock.go" -package=mocks

    echo "Generated mock for $file"
done
