#!/bin/bash
mkdir release

GOARCH=amd64 GOOS=linux go build -o release/server -tags=server

GOARCH=amd64 GOOS=linux go build -o release/client -tags=client

GOARCH=amd64 GOOS=windows go build -o release/client.exe -tags=client
