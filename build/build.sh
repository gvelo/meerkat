#!/usr/bin/env bash

VERSION=$(git describe --always --tags --abbrev=0)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
TAG=$(git describe --always --tags --abbrev=0)
BUILDUSER=$USER
BUILDDATE=$(date -u)
COMMIT=$(git rev-parse HEAD)

LDFLAGS="-X meerkat/internal/build.Version=$VERSION \
         -X meerkat/internal/build.Branch=$BRANCH \
         -X meerkat/internal/build.Tag=$TAG \
         -X meerkat/internal/build.BuildUser=$BUILDUSER \
         -X \"meerkat/internal/build.BuildDate=$BUILDDATE\" \
         -X meerkat/internal/build.Commit=$COMMIT"


go build -ldflags "$LDFLAGS" meerkat/cmd/meerkat



