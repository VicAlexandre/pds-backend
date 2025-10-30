#!/usr/bin/env bash
# exit on error
set -o errexit

STORAGE_DIR=/opt/render/project/.render

if [[ ! -d $STORAGE_DIR/chromium ]]; then
  echo "...Installing Chromium"
  mkdir -p $STORAGE_DIR
  # Install Chromium using apt
  apt-get update
  apt-get install -y chromium chromium-l10n
else
  echo "...Using Chromium from cache"
fi

go mod download
go build -o bin/server ./cmd/api
