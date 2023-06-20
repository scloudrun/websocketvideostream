#!/bin/bash
# go build -tags "h264enc"

go build -ldflags '-linkmode external -w -extldflags "-static" ' -tags "h264enc"
