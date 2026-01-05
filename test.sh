#!/bin/bash

URL='http://localhost:9696'
USER='admin'
PASS='prowlarr'
COOKIE_JAR=$(mktemp)

curl 'http://localhost:9696/login?returnUrl=%2F' \
  --data-raw 'username=admin&password=prowlarr' \
  -L \
  -s \
  -c "$COOKIE_JAR"

curl "http://localhost:9696/initialize.json" \
  -s \
  -b "$COOKIE_JAR" | \
    jq '.apiKey'