#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o exporter
scp exporter root@apps.test3.test:~/
ssh root@apps.test3.test ./exporter
