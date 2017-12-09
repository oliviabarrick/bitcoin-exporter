#!/bin/bash
find -path ./vendor -prune -o -name '*.go' -print |entr -c docker-compose -f docker/docker-compose.yml up &
ls docker/Dockerfile docker/docker-compose.yml Gopkg.toml |entr -c docker-compose -f docker/docker-compose.yml up --build
