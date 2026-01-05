#!/bin/bash

URL='http://localhost:9696'
USER=''
PASS=''
COOKIE_JAR=$(mktemp)

curl 'http://localhost:9696/login?returnUrl=%2F' \
  --data-raw "username=$USER&password=$PASS" \
  -L \
  -s \
  -c "$COOKIE_JAR"

curl "http://localhost:9696/initialize.json" \
  -s \
  -b "$COOKIE_JAR" | \
    jq '.apiKey'