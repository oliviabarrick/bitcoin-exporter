#!/bin/bash
ls ./*.go ./**/*.go test_data/*.xml |entr -c docker-compose -f docker/docker-compose.yml up &
ls docker/Dockerfile docker/docker-compose.yml |entr -c docker-compose -f docker/docker-compose.yml up --build
